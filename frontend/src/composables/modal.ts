import { ref } from "vue";

// Encapsulate all modal logic for better reuse and logic encapsulation
export const useModal = () => {
  const showModal = ref<boolean>(false);
  const modalTitle = ref<string>("");
  const modalContent = ref<string>("");
  const modalButtonText = ref<string>("");
  const modalFunction = ref<() => Promise<void>>();

  function triggerModal(
    title: string,
    content: string,
    buttonText: string,
    functionToRun: () => Promise<void>
  ) {
    showModal.value = true;
    modalTitle.value = title;
    modalContent.value = content;
    modalButtonText.value = buttonText;
    modalFunction.value = () => functionToRun();
  }
  return {
    showModal,
    modalTitle,
    modalContent,
    modalButtonText,
    modalFunction,
    triggerModal,
  };
};
