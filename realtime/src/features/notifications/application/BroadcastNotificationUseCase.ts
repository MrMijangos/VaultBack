import { ConnectionRegistry } from "../domain/repositories/ConnectionRegistry";
import { NotificationRepository } from "../domain/repositories/NotificationRepository";

export class BroadcastNotificationUseCase {
  constructor(
    private readonly notificationRepository: NotificationRepository,
    private readonly connectionRegistry: ConnectionRegistry
  ) {}

  async execute(): Promise<void> {
    await this.notificationRepository.onNewNotification((notification) => {
      this.connectionRegistry.sendToUser(notification.userId, notification);
    });
  }
}
