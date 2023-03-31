<template>
  <div @click="toggleCollapse" role="button">
    <div class="collapsible" :class="{ active: isActive }">
      <slot></slot>
    </div>
    <div ref="content" class="content">
      <p>
        Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod
        tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim
        veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea
        commodo consequat.
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { toggleCollapseContent } from "../../composables/collapsible";

const isActive = ref<Boolean>(false);
const content = ref<HTMLElement | null>(null);

const emit = defineEmits<{
  (e: "resize", childSize: string): void;
}>();

// Collapse the item itself and emit an event
// so the parent that contains multiple items can resize itself as well
function toggleCollapse() {
  isActive.value = !isActive.value;
  toggleCollapseContent(content.value);
  if (content.value) emit("resize", content.value.style.maxHeight);
}
</script>

<style scoped>
.collapsible {
  background-color: #777;
  color: white;
  cursor: pointer;
  padding: 18px;
  width: 100%;
  border: none;
  text-align: left;
  outline: none;
  font-size: 15px;
  box-shadow: rgba(0, 0, 0, 0.24) 0px 3px 8px;
}

.content {
  padding: 0 10px;
  overflow: hidden;
  background-color: #f1f1f1;
  max-height: 0;
  transition: max-height 0.2s ease-out;
}

/* ::::: */
.collapsible:after {
  content: "\002B";
  font-weight: bold;
  float: right;
  margin-left: 5px;
}

.collapsible.active:after {
  content: "\2212";
}
</style>