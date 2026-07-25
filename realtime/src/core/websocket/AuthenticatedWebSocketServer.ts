import { Server as HTTPServer, IncomingMessage } from "http";
import { URL } from "url";
import { WebSocket, WebSocketServer } from "ws";

import { verifyToken } from "../security/JWT";
import { ConnectionRegistry } from "./ConnectionRegistry";

// Un socket que perdió la conexión de golpe (cambio de red, app en segundo
// plano, NAT/proxy que corta conexiones inactivas) no siempre manda un
// close/FIN limpio -- para `ws`, `readyState` se queda en OPEN aunque ya
// nadie esté del otro lado, así que sendToUser() manda al vacío en
// silencio. El ping/pong de abajo detecta y cierra esos sockets zombis.
type TrackedSocket = WebSocket & { isAlive?: boolean };

const HEARTBEAT_INTERVAL_MS = 30_000;

function extractToken(request: IncomingMessage): string | null {
  const url = new URL(request.url ?? "", "http://localhost");
  const queryToken = url.searchParams.get("token");
  if (queryToken) return queryToken;

  const cookieHeader = request.headers.cookie;
  if (!cookieHeader) return null;

  const match = cookieHeader
    .split(";")
    .map((c) => c.trim())
    .find((c) => c.startsWith("vault_token="));

  return match ? decodeURIComponent(match.split("=")[1]) : null;
}

export class AuthenticatedWebSocketServer {
  private readonly wss: WebSocketServer;

  constructor(
    private readonly httpServer: HTTPServer,
    private readonly connectionRegistry: ConnectionRegistry,
    private readonly jwtSecret: string
  ) {
    this.wss = new WebSocketServer({ noServer: true });
    this.registerUpgradeHandler();
    this.startHeartbeat();
  }

  private startHeartbeat(): void {
    const interval = setInterval(() => {
      for (const ws of this.wss.clients) {
        const tracked = ws as TrackedSocket;
        if (tracked.isAlive === false) {
          tracked.terminate();
          continue;
        }
        tracked.isAlive = false;
        tracked.ping();
      }
    }, HEARTBEAT_INTERVAL_MS);
    this.wss.on("close", () => clearInterval(interval));
  }

  private registerUpgradeHandler(): void {
    this.httpServer.on("upgrade", (request, socket, head) => {
      const url = new URL(request.url ?? "", "http://localhost");
      if (url.pathname !== "/ws") {
        socket.write("HTTP/1.1 404 Not Found\r\n\r\n");
        socket.destroy();
        return;
      }

      const token = extractToken(request);
      if (!token) {
        socket.write("HTTP/1.1 401 Unauthorized\r\n\r\n");
        socket.destroy();
        return;
      }

      let userId: string;
      try {
        userId = verifyToken(token, this.jwtSecret).userId;
      } catch {
        socket.write("HTTP/1.1 401 Unauthorized\r\n\r\n");
        socket.destroy();
        return;
      }

      this.wss.handleUpgrade(request, socket, head, (ws) => {
        const tracked = ws as TrackedSocket;
        tracked.isAlive = true;
        ws.on("pong", () => {
          tracked.isAlive = true;
        });

        this.connectionRegistry.add(userId, ws);
        ws.on("close", () => this.connectionRegistry.remove(userId, ws));
        ws.on("error", () => this.connectionRegistry.remove(userId, ws));
      });
    });
  }
}
