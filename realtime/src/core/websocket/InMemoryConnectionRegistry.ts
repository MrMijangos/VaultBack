import { WebSocket } from "ws";

import { ConnectionRegistry } from "./ConnectionRegistry";

export class InMemoryConnectionRegistry implements ConnectionRegistry {
  private connectionsByUser: Map<string, Set<WebSocket>> = new Map();

  add(userId: string, socket: WebSocket): void {
    const sockets = this.connectionsByUser.get(userId) ?? new Set<WebSocket>();
    sockets.add(socket);
    this.connectionsByUser.set(userId, sockets);
  }

  remove(userId: string, socket: WebSocket): void {
    const sockets = this.connectionsByUser.get(userId);
    if (!sockets) return;

    sockets.delete(socket);
    if (sockets.size === 0) {
      this.connectionsByUser.delete(userId);
    }
  }

  sendToUser(userId: string, payload: unknown): void {
    const sockets = this.connectionsByUser.get(userId);
    if (!sockets) return;

    const message = JSON.stringify(payload);
    for (const socket of sockets) {
      if (socket.readyState === WebSocket.OPEN) {
        socket.send(message);
      }
    }
  }
}
