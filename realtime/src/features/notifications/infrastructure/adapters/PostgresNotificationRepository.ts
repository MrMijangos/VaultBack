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

export class PostgresNotificationRepository implements NotificationRepository {
  constructor(private readonly databaseUrl: string) {}

  async onNewNotification(handler: NotificationHandler): Promise<void> {
    const client = new Client({ connectionString: this.databaseUrl });
    await client.connect();
    await client.query("LISTEN new_notification");

    client.on("notification", (message) => {
      if (!message.payload) return;
      const row = JSON.parse(message.payload) as NotificationRow;
      handler(toNotification(row));
    });

    client.on("error", (err) => {
      console.error("error en la conexion de escucha de postgres:", err);
    });

    console.log("escuchando notificaciones nuevas en postgres (canal new_notification)");
  }
}
