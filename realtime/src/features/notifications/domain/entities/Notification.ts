export interface Notification {
  id: string;
  userId: string;
  type: string;
  subtype: string;
  title: string;
  body: string;
  data: unknown;
  read: boolean;
  createdAt: string;
}
