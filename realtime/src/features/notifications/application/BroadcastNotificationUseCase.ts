import { ConnectionRegistry } from "../../../core/websocket/ConnectionRegistry";
import { NotificationRepository } from "../domain/repositories/NotificationRepository";

export class BroadcastNotificationUseCase {
  constructor(
    private readonly notificationRepository: NotificationRepository,
    private readonly connectionRegistry: ConnectionRegistry
  ) {}

  async execute(): Promise<void> {
    await this.notificationRepository.onNewNotification((notification) => {
      // snake_case a propósito -- mismo contrato que NotificationResponse
      // en api/ (Go), que es lo que el cliente Flutter ya sabe parsear
      // (NotificationModel.fromJson). Ver el comentario equivalente en
      // BroadcastChatMessageUseCase.
      this.connectionRegistry.sendToUser(notification.userId, {
        event: "notification",
        id: notification.id,
        user_id: notification.userId,
        type: notification.type,
        subtype: notification.subtype,
        title: notification.title,
        body: notification.body,
        data: notification.data,
        read: notification.read,
        created_at: notification.createdAt,
      });
    });
  }
}
