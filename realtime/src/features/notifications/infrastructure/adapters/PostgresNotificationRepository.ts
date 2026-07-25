import { Client } from "pg";

import { Notification } from "../../domain/entities/Notification";
import {
  NotificationHandler,
  NotificationRepository,
} from "../../domain/repositories/NotificationRepository";

interface NotificationRow {
  id: string;
  user_id: string;
  type: string;
  subtype: string;
  title: string;
  body: string;
  data: unknown;
  read: boolean;
  created_at: string;
}

function toNotification(row: NotificationRow): Notification {
  return {
    id: row.id,
    userId: row.user_id,
    type: row.type,
    subtype: row.subtype,
    title: row.title,
    body: row.body,
    data: row.data,
    read: row.read,
    createdAt: row.created_at,
  };
}

const BASE_RECONNECT_DELAY_MS = 1_000;
const MAX_RECONNECT_DELAY_MS = 30_000;

export class PostgresNotificationRepository implements NotificationRepository {
  constructor(private readonly databaseUrl: string) {}

  async onNewNotification(handler: NotificationHandler): Promise<void> {
    let delayMs = BASE_RECONNECT_DELAY_MS;

    // Esta es la ÚNICA conexión que escucha el canal `new_notification` --
    // si se cae (blip de red, el proveedor la recicla, reinicio de la BD)
    // y no se reconecta, el proceso sigue corriendo pero ninguna
    // notificación vuelve a llegar por WebSocket para NADIE hasta
    // reiniciarlo entero. Backoff exponencial acotado, y se reinicia a la
    // base en cuanto un LISTEN vuelve a quedar establecido.
    const connect = async (): Promise<void> => {
      const client = new Client({
        connectionString: this.databaseUrl,
        ssl: { rejectUnauthorized: false },
      });

      let reconnecting = false;
      const scheduleReconnect = () => {
        if (reconnecting) return;
        reconnecting = true;
        client.removeAllListeners();
        client.end().catch(() => {});
        console.warn(`conexion LISTEN de notificaciones perdida, reintentando en ${delayMs}ms...`);
        setTimeout(() => {
          void connect();
        }, delayMs);
        delayMs = Math.min(delayMs * 2, MAX_RECONNECT_DELAY_MS);
      };

      client.on("notification", (message) => {
        if (!message.payload) return;
        const row = JSON.parse(message.payload) as NotificationRow;
        handler(toNotification(row));
      });

      client.on("error", (err) => {
        console.error("error en la conexion de escucha de postgres:", err);
        scheduleReconnect();
      });

      client.on("end", scheduleReconnect);

      try {
        await client.connect();
        await client.query("LISTEN new_notification");
        delayMs = BASE_RECONNECT_DELAY_MS;
        console.log("escuchando notificaciones nuevas en postgres (canal new_notification)");
      } catch (err) {
        console.error("no se pudo conectar el LISTEN de notificaciones:", err);
        scheduleReconnect();
      }
    };

    await connect();
  }
}
