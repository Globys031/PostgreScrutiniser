export interface UiNotificationContainer {
  addNotification: (
    type: string,
    title: string,
    message: string,
    duration: number
  ) => void;
}

export type UiNotificationType = {
  type: string;
  title: string;
  message: string;
  duration: number;
};
