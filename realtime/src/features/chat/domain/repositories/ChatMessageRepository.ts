import { ChatMessage } from "../entities/ChatMessage";

export type ChatMessageHandler = (message: ChatMessage) => void;

export interface ChatMessageRepository {
  onNewChatMessage(handler: ChatMessageHandler): Promise<void>;
}
