<template>
  <UiCollapsible ref="collapsible">
    <template #title> Recommendations (%n) </template>
    <template #content>
      <UiCollapsibleItem @resize="(size: string) => setChildSize(size)">
        test
      </UiCollapsibleItem>
      <UiCollapsibleItem @resize="(size: string) => setChildSize(size)">
        test1
      </UiCollapsibleItem>
      <UiCollapsibleItem @resize="(size: string) => setChildSize(size)">
        test2
      </UiCollapsibleItem>
      <UiCollapsibleItem @resize="(size: string) => setChildSize(size)">
        test3
      </UiCollapsibleItem>
      <UiCollapsibleItem @resize="(size: string) => setChildSize(size)">
        test4
      </UiCollapsibleItem>
    </template>
  </UiCollapsible>
</template>

<script setup lang="ts">
import { ref } from "vue";
import type { UiCollapsibleComponent } from "../openapi/types/collapsibleComponent";

const collapsible = ref<UiCollapsibleComponent | null>(null);

// After `UiCollapsibleItem` is resized, resize `UiCollapsible` as well
function setChildSize(size: string) {
  if (collapsible.value) {
    collapsible.value.resizeContentMaxHeight(size);
  }
}
</script>

<style scoped></style>