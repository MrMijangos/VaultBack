import { Server as HTTPServer } from "http";

import { Config } from "../../../core/config/Config";
import { BroadcastNotificationUseCase } from "../application/BroadcastNotificationUseCase";
import { InMemoryConnectionRegistry } from "./adapters/InMemoryConnectionRegistry";
import { PostgresNotificationRepository } from "./adapters/PostgresNotificationRepository";
import { NotificationWebSocketServer } from "./websocket/NotificationWebSocketServer";

export function buildBroadcastNotificationUseCase(
  httpServer: HTTPServer,
  config: Config
): BroadcastNotificationUseCase {
  const connectionRegistry = new InMemoryConnectionRegistry();
  const notificationRepository = new PostgresNotificationRepository(config.databaseUrl);

  new NotificationWebSocketServer(httpServer, connectionRegistry, config.jwtSecret);

  return new BroadcastNotificationUseCase(notificationRepository, connectionRegistry);
}
