import { Config } from "../../../core/config/Config";
import { ConnectionRegistry } from "../../../core/websocket/ConnectionRegistry";
import { BroadcastChatMessageUseCase } from "../application/BroadcastChatMessageUseCase";
import { PostgresChatMessageRepository } from "./adapters/PostgresChatMessageRepository";

export function buildBroadcastChatMessageUseCase(
  connectionRegistry: ConnectionRegistry,
  config: Config
): BroadcastChatMessageUseCase {
  const chatMessageRepository = new PostgresChatMessageRepository(config.databaseUrl);

  return new BroadcastChatMessageUseCase(chatMessageRepository, connectionRegistry);
}
