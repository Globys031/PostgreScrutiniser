<template>
  <div class="notifications-container">
    <UiNotification
      v-for="(item, index) in notifications"
      :key="index"
      :type="item.type"
      :title="item.title"
      :message="item.message"
      :duration="item.duration"
      @closed="onNotificationClosed(index)"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import UiNotification from "./UiNotification.vue";
import type { UiNotificationType } from "@/types/notification";

// Expose addNotification to be used outside of this component
defineExpose({
  addNotification,
});

const notifications = ref<UiNotificationType[]>([]);

function addNotification(
  type: string,
  title: string,
  message: string,
  duration: number = 3000
) {
  notifications.value.push({ type, title, message, duration });
}

function onNotificationClosed(index: number) {
  notifications.value.splice(index, 1);
}
</script>

<style scoped>
.notifications-container {
  z-index: 1;
  position: fixed;
  bottom: 20px;
  right: 20px;
  display: flex;
  flex-direction: column-reverse;
  align-items: flex-end;
  gap: 10px;
}
</style>
