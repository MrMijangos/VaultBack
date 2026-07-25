import { ConnectionRegistry } from "../../../core/websocket/ConnectionRegistry";
import { ChatMessageRepository } from "../domain/repositories/ChatMessageRepository";

export class BroadcastChatMessageUseCase {
  constructor(
    private readonly chatMessageRepository: ChatMessageRepository,
    private readonly connectionRegistry: ConnectionRegistry
  ) {}

  async execute(): Promise<void> {
    await this.chatMessageRepository.onNewChatMessage((message) => {
      // El payload va en snake_case a propósito: es el mismo contrato que
      // ChatMessageResponse en api/ (Go), que es lo que el cliente Flutter
      // ya sabe parsear (ChatMessageModel.fromJson) para las respuestas
      // REST. Mandar los nombres de campo camelCase de esta entidad TS tal
      // cual rompía ese parser en el primer mensaje que llegaba por
      // WebSocket.
      this.connectionRegistry.sendToUser(message.recipientId, {
        event: "chat_message",
        id: message.id,
        sender_id: message.senderId,
        recipient_id: message.recipientId,
        cipher_text: message.cipherText,
        encrypted_aes_key: message.encryptedAesKey,
        encrypted_aes_key_sender: message.encryptedAesKeySender,
        iv: message.iv,
        status: message.status,
        created_at: message.createdAt,
      });
    });
  }
}
