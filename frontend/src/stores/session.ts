import { defineStore } from "pinia";
import { ref, watch, computed } from "vue";
import { isValidIp } from "@/composables/ips";

export const useSessionStore = defineStore("session", () => {
  const token = ref<string>(localStorage.getItem("token") || "");
  const username = ref<string>(localStorage.getItem("username") || "");
  const hostname = ref<string>(localStorage.getItem("hostname") || "");

  const baseAPIPath = computed(() => {
    // If using a domain name instead of ip address, should attempt to connect via https
    const protocol = isValidIp(hostname.value) ? "http" : "https";

    return `${protocol}://${hostname.value}:${
      import.meta.env.VITE_BACKEND_PORT
    }/api`;
  });

  // Use local storage to persist data upon page reload
  watch(token, (newToken: string) => {
    localStorage.setItem("token", newToken);
  });
  watch(username, (newUsername: string) => {
    localStorage.setItem("username", newUsername);
  });
  watch(hostname, (newHostname: string) => {
    localStorage.setItem("hostname", newHostname);
  });

  // Watchers will clear localStorage
  function clearSession() {
    localStorage.removeItem("token");
    localStorage.removeItem("username");
    localStorage.removeItem("hostname");
    token.value = "";
    username.value = "";
    hostname.value = "";
  }

  return {
    token,
    username,
    hostname,
    baseAPIPath,
    clearSession,
  };
});
