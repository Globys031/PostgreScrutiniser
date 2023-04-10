<template>
  <component :is="vButtonComponent" />
</template>

<script setup lang="ts">
import { h, useCssModule, computed, ref, watch } from "vue";
import {
  IconBxUndo,
  IconBxInfoCircle,
  IconBxError,
  IconBxSend,
} from "@iconify-prerendered/vue-bx";
import type { VNode, VNodeArrayChildren } from "vue";

const props = withDefaults(
  defineProps<{
    disabled?: boolean; // make button unclickable
    countdown?: number; // countdown till button becomes clickable
    type: string; // undo/error/submit/info
    text: string; // text inside the button
  }>(),
  {
    disabled: false,
    countdown: 3,
  }
);

const styles = useCssModule();

// Using a computed ref for actual display to avoid side effects
type VButtonType = ReturnType<typeof createVButton>;
const vButtonRef = ref<VButtonType>(createVButton());
const vButtonComponent = computed(() => vButtonRef.value);

watch(
  () => props.disabled,
  () => {
    if (!props.disabled) return;

    if (
      !vButtonRef.value.children ||
      !isVNodeArrayChildren(vButtonRef.value.children) ||
      !vButtonRef.value.children[1]
    ) {
      console.error(
        "Could not start countdown. VButton children array was possibly null"
      );
      return;
    }
    // @ts-ignore: this is causing too many issues that make code overcomplicated.
    // Making an exception for this line and disabling typescript type checking.
    const buttonText = vButtonRef.value.children[1].children[0].el;

    let countdown = props.countdown;
    buttonText.data = `Available in ${countdown}`;
    const intervalId = setInterval(() => {
      countdown -= 1;
      buttonText.data = `Available in ${countdown}`;
      if (countdown === 0) {
        clearInterval(intervalId);
        buttonText.data = props.text;
      }
    }, 1000);
  }
);

// Cannot use slots with button vnode so using a watcher instead
watch(
  () => props.text,
  (newText, oldText) => {
    if (newText === oldText) return;

    if (
      !vButtonRef.value.children ||
      !isVNodeArrayChildren(vButtonRef.value.children) ||
      !vButtonRef.value.children[1]
    ) {
      console.error("Could not replace button text");
      return;
    }
    // @ts-ignore
    const buttonText = vButtonRef.value.children[1].children[0].el;
    buttonText.data = newText;
  }
);

// Create a dynamic vNode button based on props.
function createVButton() {
  /* Equivalent to:
  ```
  <button class="btn btn-warning btn-sep">
    <div class="icon-container">
      <IconBxUndo />
    </div>
    <div class="text-container"></slot></div>
  </button>
  ``` */
  const iconNode = h("div", { class: styles["icon-container"] }, []);
  // const textSlot = h("slot", {}, instance?.slots);
  const textNode = h("div", { class: styles["text-container"] }, [props.text]);
  const vButton = h(
    "button",
    {
      class: [styles.btn, styles[`btn-${props.type}`]],
    },
    [iconNode, textNode]
  );
  addButtonIcon(vButton);
  return vButton;
}

function addButtonIcon(vButton: VButtonType) {
  const vButtonChildren = (vButton.children as VNode[])[0];
  if (!vButtonChildren.children) {
    console.error(
      "Could not add button icon because vButton children was null"
    );
    return;
  }

  if (isVNodeArrayChildren(vButtonChildren.children)) {
    switch (props.type) {
      case "submit":
        vButtonChildren.children.push(h(IconBxSend));
        break;
      case "warning":
        vButtonChildren.children.push(h(IconBxError));
        break;
      case "info":
        vButtonChildren.children.push(h(IconBxInfoCircle));
        break;
      case "undo":
        vButtonChildren.children.push(h(IconBxUndo));
        break;
      default:
        console.error("Wrong button type submitted");
    }
  }
}

// Define a type guard for VNodeArrayChildren
function isVNodeArrayChildren(
  children: unknown
): children is VNodeArrayChildren {
  return Array.isArray(children);
}
</script>

<style module>
.btn {
  border: none;
  font-family: "Lato";
  font-size: inherit;
  color: inherit;
  background: none;
  cursor: pointer;
  display: inline-block;
  margin-bottom: 1em;
  text-transform: uppercase;
  letter-spacing: 1px;
  outline: none;
  position: relative;
  transition: all 0.3s;

  display: flex;
  padding: 0;
}

.icon-container {
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(0, 0, 0, 0.1 /* Adjust this value to change shade */);
  padding: 13px;
  height: 100%;
  font-size: 20px;
}

.text-container {
  padding: 20px 15px;
}

/* Button 1 */
.btn-info {
  background: #3498db;
}

.btn-info:hover {
  background: #2980b9;
}

.btn-info:active {
  top: 2px; /* creates effect that button is moving when clicked */
}

/* Button 2 */
.btn-submit {
  background: #2ecc71;
  color: #fff;
}

.btn-submit:hover {
  background: #27ae60;
}

.btn-submit:active {
  top: 2px;
}

/* Button 3 */
.btn-warning {
  background: #e74c3c;
  color: #fff;
}

.btn-warning:hover {
  background: #c0392b;
}

.btn-warning:active {
  top: 2px;
}

/* Button 4 */
.btn-undo {
  background: #34495e;
  color: #fff;
}

.btn-undo:hover {
  background: #2c3e50;
}

.btn-undo:active {
  top: 2px;
}
</style>
