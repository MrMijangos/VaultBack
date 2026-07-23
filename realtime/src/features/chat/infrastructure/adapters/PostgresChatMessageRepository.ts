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
    iv: row.iv,
    status: row.status,
    createdAt: row.created_at,
  };
}

export class PostgresChatMessageRepository implements ChatMessageRepository {
  constructor(private readonly databaseUrl: string) {}

  async onNewChatMessage(handler: ChatMessageHandler): Promise<void> {
    const client = new Client({ connectionString: this.databaseUrl });
    await client.connect();
    await client.query("LISTEN new_chat_message");

    client.on("notification", (message) => {
      if (!message.payload) return;
      const row = JSON.parse(message.payload) as ChatMessageRow;
      handler(toChatMessage(row));
    });

    client.on("error", (err) => {
      console.error("error en la conexion de escucha de postgres (chat):", err);
    });

    console.log("escuchando mensajes de chat nuevos en postgres (canal new_chat_message)");
  }
}
