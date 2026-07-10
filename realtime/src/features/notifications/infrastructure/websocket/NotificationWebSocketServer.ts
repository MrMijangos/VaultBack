import { Server as HTTPServer, IncomingMessage } from "http";
import { URL } from "url";
import { WebSocketServer } from "ws";

import { verifyToken } from "../../../../core/security/JWT";
import { ConnectionRegistry } from "../../domain/repositories/ConnectionRegistry";

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

export class NotificationWebSocketServer {
  private readonly wss: WebSocketServer;

  constructor(
    private readonly httpServer: HTTPServer,
    private readonly connectionRegistry: ConnectionRegistry,
    private readonly jwtSecret: string
  ) {
    this.wss = new WebSocketServer({ noServer: true });
    this.registerUpgradeHandler();
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
        this.connectionRegistry.add(userId, ws);
        ws.on("close", () => this.connectionRegistry.remove(userId, ws));
        ws.on("error", () => this.connectionRegistry.remove(userId, ws));
      });
    });
  }
}
