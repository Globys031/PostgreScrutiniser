import { defineStore } from "pinia";
import { ref, watch } from "vue";

export const useSessionStore = defineStore("session", () => {
  const token = ref<string>(localStorage.getItem("token") || "");
  const username = ref<string>(localStorage.getItem("username") || "");

  // Use local storage to persist data upon page reload
  watch(token, (oldToken: string) => {
    localStorage.setItem("token", oldToken);
  });

  watch(username, (oldUsername: string) => {
    localStorage.setItem("username", oldUsername);
  });

  // Watchers will clear localStorage
  function clearSession() {
    token.value = "";
    username.value = "";
  }

  return { token, username, clearSession };
});
