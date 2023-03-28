import type { Component, ComputedOptions, MethodOptions } from "vue";

import UiSpinner from "./UiSpinner.vue";
import UiButton from "./UiButton.vue";

interface ComponentMap {
  [key: string]: Component<any, any, any, ComputedOptions, MethodOptions>;
}

export const globalComponents: ComponentMap = {
  'UiSpinner': UiSpinner,
  'UiButton': UiButton,
};
