<template>
  <div @click="toggleCollapse" class="header">
    <IconBxCommentCheck v-if="isSuggestions" class="icon" />
    <IconBxCommentX v-else class="icon" />
    <div class="text" :class="{ active: isActive }">
      <slot name="title" />
    </div>
  </div>

  <div ref="content" class="content">
    <slot name="content" />
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { toggleCollapseContent } from "../../composables/collapsible";
import {
  IconBxCommentCheck,
  IconBxCommentX,
} from "@iconify-prerendered/vue-bx";

defineExpose({
  resizeContentMaxHeight,
});

defineProps<{
  isSuggestions?: boolean;
}>();

const isActive = ref<boolean>(false);
const content = ref<HTMLElement | null>(null);

function toggleCollapse() {
  isActive.value = !isActive.value;
  toggleCollapseContent(content.value);
}

// Function for resizing the maximum height based on contents inside element
// @childSize {*} [options] Override http request option.
// @param {childSize} - the max height of a child collapsible.
// function resizeContentMaxHeight(childSize: string) {
function resizeContentMaxHeight() {
  if (!content.value) {
    console.error("Content inside collapsible was null");
    return;
  }
  const childArray = Array.from(content.value.children);
  const totalHeight = childArray.reduce(
    (totalHeight, child) => totalHeight + child.scrollHeight,
    0
  );

  content.value
    ? (content.value.style.maxHeight = totalHeight + "px")
    : console.error("content ref is undefined");
}
</script>

<style scoped>
.header {
  cursor: pointer;
  position: relative;
  font-size: 18px;
  background-color: var(--vt-primary-background);
  color: var(--vt-c-white);
  padding: 20px 50px 20px;
}

.header:before {
  content: "";
  position: absolute;
  top: 0;
  left: 0;
  width: 50px;
  height: 100%;
}

.header .icon {
  position: absolute;
  left: 10px;
  font-size: 28px;
}

.content {
  max-height: 0;
  overflow: hidden;

  /* force content to not disappear when viewport too small */
  min-width: 100px;

  box-shadow: rgba(0, 0, 0, 0.24) 0px 3px 8px;
  transition: max-height 0.2s ease-out;
}

/* ::::: */
.text:after {
  content: "\002B";
  font-weight: bold;
  float: right;
  margin-left: 5px;
}

.text.active:after {
  content: "\2212";
}

/* Temporary hack to force header and content to resize at the same rate */
@media screen and (max-width: 414px) {
  .header,
  .content {
    width: 100%;
  }
}
</style>
