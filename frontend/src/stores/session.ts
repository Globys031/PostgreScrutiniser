import { defineStore } from "pinia";
import { ref, watch } from "vue";

export const useSessionStore = defineStore("session", () => {
  const token = ref<string>(localStorage.getItem("token") || "");
  const username = ref<string>(localStorage.getItem("username") || "");

  // Use local storage to persist data upon page reload
  watch(token, (newToken: string) => {
    localStorage.setItem("token", newToken);
  });

  watch(username, (newUsername: string) => {
    localStorage.setItem("username", newUsername);
  });

  // Watchers will clear localStorage
  function clearSession() {
    token.value = "";
    username.value = "";
  }

  return { token, username, clearSession };
});
