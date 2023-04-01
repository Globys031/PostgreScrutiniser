<template>
  <div
    @click="toggleCollapse"
    role="button"
    class="collapsible"
    :class="{ active: isActive }"
  >
    <slot name="title"></slot>
  </div>
  <div ref="content" class="content">
    <slot name="content" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import {
  toggleCollapseContent,
  resizeContentIfOpen,
} from "../../composables/collapsible";

const isActive = ref<Boolean>(false);
const content = ref<HTMLElement | null>(null);

const emit = defineEmits<{
  (e: "resize", childSize: string): void;
}>();

// When `content` slot content changes, the collapsible should be resized
// to take into account new content size
let observer: MutationObserver;
onMounted(() => {
  // observer = new MutationObserver((mutations) => {
  observer = new MutationObserver(() => {
    // This callback will be called whenever the slot content changes
    resizeContentIfOpen(content.value);
  });

  if (content.value) {
    observer.observe(content.value, {
      childList: true,
      subtree: true,
    });
  }
});

onUnmounted(() => {
  if (observer) {
    observer.disconnect();
  }
});

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
  background-color: #f1f1f1;
  font-weight: bold;
  cursor: pointer;
  padding: 18px;
  width: 100%;
  border: 1px solid rgba(0, 0, 0, 0.24);

  text-align: left;
  outline: none;
  font-size: 15px;
  box-shadow: rgba(0, 0, 0, 0.24) 0px 3px 8px;
}

.content {
  padding: 0 10px;
  overflow: hidden;
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
