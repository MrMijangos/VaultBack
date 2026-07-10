import { Notification } from "../entities/Notification";

export type NotificationHandler = (notification: Notification) => void;

export interface NotificationRepository {
  onNewNotification(handler: NotificationHandler): Promise<void>;
}
