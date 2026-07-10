import { WebSocket } from "ws";

export interface ConnectionRegistry {
  add(userId: string, socket: WebSocket): void;
  remove(userId: string, socket: WebSocket): void;
  sendToUser(userId: string, payload: unknown): void;
}
