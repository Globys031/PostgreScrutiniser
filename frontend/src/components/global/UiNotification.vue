<template>
  <transition name="slide-fade">
    <div
      v-if="visible"
      class="notification"
      :class="type"
      @click="closeNotification"
    >
      <div class="icon-container">
        <IconBxCheckCircle v-if="type === 'success'" />
        <IconBxError v-else />
      </div>
      <div class="notification-content">
        <span class="title" v-text="title" />
        <span class="message" v-text="message" />
      </div>
    </div>
  </transition>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { IconBxCheckCircle, IconBxError } from "@iconify-prerendered/vue-bx";

// Declare props with default 3second duration
const props = withDefaults(
  defineProps<{
    msg?: string;
    type: string; // success/error
    title: string;
    message: string; // notification message
    duration?: number;
  }>(),
  {
    duration: 3000,
  }
);

const visible = ref<Boolean>(true);

onMounted(() => {
  setTimeout(() => {
    closeNotification();
  }, props.duration);
});

function closeNotification() {
  visible.value = false;
}
</script>

<style scoped>
.notification {
  border: none;
  font-family: "Lato";
  font-size: inherit;
  color: inherit;
  background: none;
  cursor: pointer;
  display: flex;
  margin-bottom: 1em;
  bottom: 20px;
  right: 20px;
  border-radius: 5px;
  color: white;
  cursor: pointer;
  max-width: 70%;
}

.notification.success {
  background-color: #4caf50;
}

.notification.success:hover {
  background: #27ae60;
}

.notification.error {
  background-color: #f44336;
}

.notification.error:hover {
  background: #c0392b;
}

.notification-content {
  padding: 15px;
  display: flex;
  flex-direction: column;
  word-break: break-word;
}

.title {
  font-size: 1.2em;
  font-weight: bold;
  margin-bottom: 0.5em;
}

.message {
  font-size: 1em;
  margin: 0;
}

.icon-container {
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(0, 0, 0, 0.1 /* Adjust this value to change shade */);
  padding: 20px;
  font-size: 20px;
}

.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: all 0.5s;
}
.slide-fade-enter,
.slide-fade-leave-to {
  transform: translateY(30px);
  opacity: 0;
}
</style>
