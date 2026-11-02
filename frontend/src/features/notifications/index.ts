// Feature: Notifications
export interface NotificationItem {
  id: string;
  type: string;
  title: string;
  message: string;
  reference_id?: string;
  is_read: boolean;
  created_at: string;
}
