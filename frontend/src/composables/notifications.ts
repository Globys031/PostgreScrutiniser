import axios from "axios";
import type { Ref } from "vue";
import type { UiNotificationContainer } from "@/types/notification";
import type { AxiosError } from "axios";
import type { ErrorMessage } from "@/openapi/api/auth";

// Display either a success or failure message for api
export const displayNotification = (
  notificationContainer: Ref<UiNotificationContainer | null>,
  notificationType: string,
  message: string
) => {
  const titleMessage =
    notificationType === "success" ? "successful" : "unsuccessful";

  if (notificationContainer.value) {
    notificationContainer.value.addNotification(
      notificationType,
      `Operation ${titleMessage}`,
      message,
      3000
    );
  }
};

export const displayNotificationError = (
  notificationContainer: Ref<UiNotificationContainer | null>,
  error: AxiosError | Error
) => {
  console.error(error);

  axios.isAxiosError(error)
    ? displayNotification(
        notificationContainer,
        "error",
        (error.response?.data as ErrorMessage)?.error_message
      )
    : displayNotification(
        notificationContainer,
        "error",
        "Something went wrong. See console tab for more info"
      );
};
