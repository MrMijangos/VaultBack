import { Client } from "pg";

import { ChatMessage } from "../../domain/entities/ChatMessage";
import {
  ChatMessageHandler,
  ChatMessageRepository,
} from "../../domain/repositories/ChatMessageRepository";

interface ChatMessageRow {
  id: string;
  sender_id: string;
  recipient_id: string;
  cipher_text: string;
  encrypted_aes_key: string;
  encrypted_aes_key_sender: string;
  iv: string;
  status: string;
  created_at: string;
}

function toChatMessage(row: ChatMessageRow): ChatMessage {
  return {
    id: row.id,
    senderId: row.sender_id,
    recipientId: row.recipient_id,
    cipherText: row.cipher_text,
    encryptedAesKey: row.encrypted_aes_key,
    encryptedAesKeySender: row.encrypted_aes_key_sender,
    iv: row.iv,
    status: row.status,
    createdAt: row.created_at,
  };
}

const BASE_RECONNECT_DELAY_MS = 1_000;
const MAX_RECONNECT_DELAY_MS = 30_000;

export class PostgresChatMessageRepository implements ChatMessageRepository {
  constructor(private readonly databaseUrl: string) {}

  async onNewChatMessage(handler: ChatMessageHandler): Promise<void> {
    let delayMs = BASE_RECONNECT_DELAY_MS;

    // Esta es la ÚNICA conexión que escucha el canal `new_chat_message` --
    // si se cae (blip de red, el proveedor la recicla, reinicio de la BD)
    // y no se reconecta, el proceso sigue corriendo pero ningún mensaje de
    // chat vuelve a llegar por WebSocket para NADIE hasta reiniciarlo
    // entero. Backoff exponencial acotado, y se reinicia a la base en
    // cuanto un LISTEN vuelve a quedar establecido.
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
        console.warn(`conexion LISTEN de chat perdida, reintentando en ${delayMs}ms...`);
        setTimeout(() => {
          void connect();
        }, delayMs);
        delayMs = Math.min(delayMs * 2, MAX_RECONNECT_DELAY_MS);
      };

      client.on("notification", (message) => {
        if (!message.payload) return;
        const row = JSON.parse(message.payload) as ChatMessageRow;
        handler(toChatMessage(row));
      });

      client.on("error", (err) => {
        console.error("error en la conexion de escucha de postgres (chat):", err);
        scheduleReconnect();
      });

      client.on("end", scheduleReconnect);

      try {
        await client.connect();
        await client.query("LISTEN new_chat_message");
        delayMs = BASE_RECONNECT_DELAY_MS;
        console.log("escuchando mensajes de chat nuevos en postgres (canal new_chat_message)");
      } catch (err) {
        console.error("no se pudo conectar el LISTEN de chat:", err);
        scheduleReconnect();
      }
    };

    await connect();
  }
}
