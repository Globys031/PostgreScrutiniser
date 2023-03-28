<template>
  <TheSidebar v-if="sessionStore.token" />
  <div :class="{ 'main-body': sessionStore.token }">
    <RouterView />
  </div>
</template>

<script setup lang="ts">
import { RouterView } from "vue-router";
import { useSessionStore } from "./stores/session";
import TheSidebar from "./components/TheSidebar.vue";

const sessionStore = useSessionStore();
</script>

<style scoped>
.main-body {
  position: relative;
  min-height: 100vh;
  top: 0;
  left: 78px;
  width: calc(100% - 78px);
  transition: all 0.5s ease;
  z-index: 2;
}

.sidebar.open + .main-body {
  left: 310px;
  width: calc(100% - 250px);
}
</style>
