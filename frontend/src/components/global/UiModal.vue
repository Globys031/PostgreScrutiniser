<template>
  <div class="modal">
    <div class="modal-content">
      <div>
        <h2 class="modal-title" v-text="title" />
        <div class="header-line"></div>
      </div>
      <div class="modal-body">
        <p v-text="content" />
      </div>
      <div class="button-container">
        <UiButton @click="confirmEmits" type="warning" :text="buttonText" />
        <UiButton @click="$emit('close')" type="undo" text="Cancel" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  title: string;
  content: string;
  buttonText: string;
  functionToRun: () => Promise<void>;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "confirm", functionToRun: Promise<void>): void;
}>();

function confirmEmits() {
  emit("confirm", props.functionToRun());
  emit("close");
}
</script>

<style scoped>
.modal {
  display: block;
  position: fixed;
  z-index: 100;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;
  overflow: auto;
  background-color: rgba(0, 0, 0, 0.4);
}

.modal-content {
  background-color: #fefefe;
  margin: 15% auto;
  border: 1px solid #888;
}

@media (min-width: 1024px) {
  .modal-content {
    width: 45%;
  }
}

@media (max-width: 1024px) {
  .modal-content {
    width: 65%;
  }
}

@media (max-width: 720px) {
  .modal-content {
    width: 90%;
  }
}

.modal-title {
  padding: 5px;
  font-size: 20px;
}

.header-line {
  height: 2px;
  background-color: #f44336;
  flex-grow: 1;
}

.modal-body {
  padding: 10px 0px 10px 10px;
}

.button-container {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
  padding: 10px;
}
</style>
