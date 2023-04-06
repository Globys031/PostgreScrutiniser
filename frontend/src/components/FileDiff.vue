<template>
  <div class="file-diff">
    <div
      v-for="(diffItem, index) in diff"
      :key="index"
      :class="diffItem.type.toLowerCase()"
    >
      <pre v-text="diffItem.line"></pre>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { FileDiffLine } from "@/openapi/api/file";

defineProps<{
  diff: Array<FileDiffLine>;
}>();
</script>

<style scoped>
.equal {
  color: #24292e;
}

.delete {
  background-color: #ffeef0;
  color: #cb2431;
}

.insert {
  background-color: #e6ffed;
  color: #22863a;
}

.file-diff pre {
  margin: 0;
  padding: 0.5em;
  white-space: pre-wrap;
  word-wrap: break-word;
  font-family: monospace;
}
</style>
