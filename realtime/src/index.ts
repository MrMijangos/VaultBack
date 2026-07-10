import { createServer } from "http";

import { loadConfig } from "./core/config/LoadConfig";
import { buildBroadcastNotificationUseCase } from "./features/notifications/infrastructure/Dependencies";

async function main(): Promise<void> {
  const config = loadConfig();

  const httpServer = createServer((_req, res) => {
    res.writeHead(200, { "Content-Type": "text/plain" });
    res.end("vault-realtime up");
  });

  const broadcastNotifications = buildBroadcastNotificationUseCase(httpServer, config);
  await broadcastNotifications.execute();

  httpServer.listen(config.port, () => {
    console.log(`vault-realtime escuchando en el puerto ${config.port}`);
  });
}

main().catch((err) => {
  console.error("error fatal al iniciar vault-realtime:", err);
  process.exit(1);
});
