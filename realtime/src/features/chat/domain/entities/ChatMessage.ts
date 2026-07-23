export interface ChatMessage {
  id: string;
  senderId: string;
  recipientId: string;
  cipherText: string;
  encryptedAesKey: string;
  encryptedAesKeySender: string;
  iv: string;
  status: string;
  createdAt: string;
}
