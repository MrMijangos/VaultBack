import { ConnectionRegistry } from "../../../core/websocket/ConnectionRegistry";
import { ChatMessageRepository } from "../domain/repositories/ChatMessageRepository";

export class BroadcastChatMessageUseCase {
  constructor(
    private readonly chatMessageRepository: ChatMessageRepository,
    private readonly connectionRegistry: ConnectionRegistry
  ) {}

  async execute(): Promise<void> {
    await this.chatMessageRepository.onNewChatMessage((message) => {
      this.connectionRegistry.sendToUser(message.recipientId, { event: "chat_message", ...message });
    });
  }
}
