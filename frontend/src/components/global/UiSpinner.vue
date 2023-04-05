<template>
  <div class="loading-container">
    <div class="loading-spinner" />
    <div class="loading-text">
      Loading
      <span
        class="dot"
        v-for="n in 3"
        :key="n"
        :style="{ animationDelay: `${n * 0.3}s` }"
        >.</span
      >
    </div>
  </div>
</template>

<script setup lang="ts"></script>

<style scoped>
.loading-container {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin: 20px;
}

.loading-spinner {
  border: 4px solid #f3f3f3;
  border-top: 4px solid #3498db;
  border-radius: 50%;
  width: 20px;
  height: 20px;
  animation: spin 2s linear infinite;
}

.loading-text {
  font-size: 14px;
  color: #333;
}

.dot {
  animation: fadeInOut 1.5s infinite;
  opacity: 0;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

@keyframes fadeInOut {
  0%,
  100% {
    opacity: 0;
  }
  50% {
    opacity: 1;
  }
}
</style>
