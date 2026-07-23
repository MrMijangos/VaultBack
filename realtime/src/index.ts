import { createServer } from "http";

import { loadConfig } from "./core/config/LoadConfig";
import { AuthenticatedWebSocketServer } from "./core/websocket/AuthenticatedWebSocketServer";
import { InMemoryConnectionRegistry } from "./core/websocket/InMemoryConnectionRegistry";
import { buildBroadcastChatMessageUseCase } from "./features/chat/infrastructure/Dependencies";
import { buildBroadcastNotificationUseCase } from "./features/notifications/infrastructure/Dependencies";

async function main(): Promise<void> {
  const config = loadConfig();

  const httpServer = createServer((_req, res) => {
    res.writeHead(200, { "Content-Type": "text/plain" });
    res.end("vault-realtime up");
  });

  const connectionRegistry = new InMemoryConnectionRegistry();
  new AuthenticatedWebSocketServer(httpServer, connectionRegistry, config.jwtSecret);

  const broadcastNotifications = buildBroadcastNotificationUseCase(connectionRegistry, config);
  await broadcastNotifications.execute();

  const broadcastChatMessages = buildBroadcastChatMessageUseCase(connectionRegistry, config);
  await broadcastChatMessages.execute();

  httpServer.listen(config.port, () => {
    console.log(`vault-realtime escuchando en el puerto ${config.port}`);
  });
}

main().catch((err) => {
  console.error("error fatal al iniciar vault-realtime:", err);
  process.exit(1);
});
