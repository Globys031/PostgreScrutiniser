<template>
  <component :is="vButtonComponent" />
</template>

<script setup lang="ts">
import { h, useCssModule, computed } from "vue";
import {
  IconBxUndo,
  IconBxInfoCircle,
  IconBxError,
  IconBxSend,
} from "@iconify-prerendered/vue-bx";
import type { VNode, VNodeArrayChildren } from "vue";

const props = defineProps<{
  type: string; // undo/error/submit/info
  text: string; // text inside the button
}>();

const styles = useCssModule();

// ... Create a virtual node button
type VButtonType = ReturnType<typeof createVButton>;
const vButton: VButtonType = createVButton();

if (isVNodeArrayChildren(vButton.children)) {
  switch (props.type) {
    case "submit":
      ((vButton.children as VNode[])[0].children as VNode[]).push(
        h(IconBxSend)
      );
      break;
    case "warning":
      ((vButton.children as VNode[])[0].children as VNode[]).push(
        h(IconBxError)
      );
      break;
    case "info":
      ((vButton.children as VNode[])[0].children as VNode[]).push(
        h(IconBxInfoCircle)
      );
      break;
    case "undo":
      ((vButton.children as VNode[])[0].children as VNode[]).push(
        h(IconBxUndo)
      );
      break;
    default:
      console.error("Wrong button type submitted");
  }
}

// Using a computed render function instead of calling `vButton`
// directly to avoid typescript errors
const vButtonComponent = computed(() => ({
  render() {
    return vButton;
  },
}));

/* Create a dynamic vNode button based on props.
Equivalent to:
```
  <button class="btn btn-warning btn-sep">
    <div class="icon-container">
      <IconBxUndo />
    </div>
    <div class="text-container">Test</div>
  </button>
```
*/
function createVButton() {
  return h(
    "button",
    {
      class: [styles.btn, styles[`btn-${props.type}`]],
    },
    [
      h("div", { class: styles["icon-container"] }, []),
      h("div", { class: styles["text-container"] }, [props.text]),
    ]
  );
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
  font-weight: 700;
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
  padding: 20px;
  height: 100%;
}

.text-container {
  padding: 25px 60px;
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
