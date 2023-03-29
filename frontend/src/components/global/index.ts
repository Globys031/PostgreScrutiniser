import type { Component, ComputedOptions, MethodOptions } from "vue";

import UiSpinner from "./UiSpinner.vue";
import UiButton from "./UiButton.vue";
import UiCollapsible from "./UiCollapsible.vue";
import UiCollapsibleItem from "./UiCollapsibleItem.vue";
import UiModal from "./UiModal.vue";
import UiNotification from "./UiNotification.vue";

interface ComponentMap {
  [key: string]: Component<any, any, any, ComputedOptions, MethodOptions>;
}

export const globalComponents: ComponentMap = {
  'UiSpinner': UiSpinner,
  'UiButton': UiButton,
  'UiCollapsible': UiCollapsible,
  'UiModal': UiModal,
  'UiNotification': UiNotification,
  'UiCollapsibleItem': UiCollapsibleItem,
};
