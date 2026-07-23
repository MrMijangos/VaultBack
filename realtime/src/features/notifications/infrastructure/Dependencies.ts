import { Config } from "../../../core/config/Config";
import { ConnectionRegistry } from "../../../core/websocket/ConnectionRegistry";
import { BroadcastNotificationUseCase } from "../application/BroadcastNotificationUseCase";
import { PostgresNotificationRepository } from "./adapters/PostgresNotificationRepository";

export function buildBroadcastNotificationUseCase(
  connectionRegistry: ConnectionRegistry,
  config: Config
): BroadcastNotificationUseCase {
  const notificationRepository = new PostgresNotificationRepository(config.databaseUrl);

  return new BroadcastNotificationUseCase(notificationRepository, connectionRegistry);
}
