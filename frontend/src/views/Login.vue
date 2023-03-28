<template>
  <div id="login-form-wrap">
    <h2>Login Page</h2>
    <form id="login-form" @submit.prevent="submitForm">
      <p>
        <input
          type="text"
          id="username"
          name="username"
          placeholder="Username"
          v-model="username"
          required
        /><i class="validation"><span></span><span></span></i>
      </p>
      <p>
        <input
          type="password"
          id="password"
          name="password"
          placeholder="Password"
          v-model="password"
          required
        />
        <i class="validation"><span></span><span></span></i>
      </p>
      <p class="submit-area">
        <input v-if="!isLoading" type="submit" id="login" value="Login" />
        <UiSpinner v-else />
      </p>
    </form>

    <span
      class="response-msg"
      :class="{
        'error-msg': gotError,
        'success-msg': !gotError && apiResponse,
      }"
      v-text="apiResponse"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import axios from "axios";
import { useSessionStore } from "../stores/session";
import { AuthApiFp } from "../openapi/auth/api";
import { Configuration } from "../openapi/auth/configuration";

import type { ErrorMessage } from "../openapi/auth/api";

const sessionStore = useSessionStore();

const postLogin = AuthApiFp(
  new Configuration({
    basePath: "http://192.168.56.102:9090/api",
  })
).postLogin;

const username = ref<string>("");
const password = ref<string>("");
const apiResponse = ref<string>(""); // login request response message
const gotError = ref<Boolean>(false);

const isLoading = ref<Boolean>(false);

async function submitForm() {
  try {
    isLoading.value = true;

    // Prepare login API request
    const postLoginRequest = postLogin({
      name: username.value,
      password: password.value,
    });
    const createRequest = await postLoginRequest;

    // Execute API request and save data in our session store
    const { data } = await createRequest();
    sessionStore.$patch({
      token: data.token,
      username: username.value,
    });

    apiResponse.value = "Request successful";
    gotError.value = false;
    isLoading.value = false;
  } catch (error) {
    if (axios.isAxiosError(error)) {
      apiResponse.value = `Server response: ${
        (error.response?.data as ErrorMessage)?.error_message
      }`;
      gotError.value = true;
    } else {
      apiResponse.value = "Something went wrong. See console tab for more info";
      console.error("Login Failed:", error);
    }
    isLoading.value = false;
  }
}
</script>

<style scoped>
h2 {
  font-weight: 300;
  text-align: center;
}

p {
  position: relative;
  padding: 10px;
}

.response-msg {
  transition: all 0.5s ease;
  display: inline-block;
  padding: 0.5rem;
  margin: 1rem 0;
  font-weight: bold;
  font-size: 1rem;
  width: 100%;
}

.error-msg {
  background-color: var(--color-error-background);
  color: var(--vt-error-color);
}

.success-msg {
  background-color: var(--color-success-background);
  color: var(--vt-success-color);
}

.submit-area,
#login-form {
  padding: 0 0px;
}

#login-form-wrap {
  font-size: 1.6rem;
  font-family: "Open Sans", sans-serif;
  background-color: #fff;
  margin: 30px auto;
  text-align: center;
  padding: 20px 0 20px 0;
  border-radius: 4px;
  box-shadow: 0px 30px 50px 0px rgba(0, 0, 0, 0.2);
}

@media (min-width: 1024px) {
  #login-form-wrap {
    width: 45%;
  }
}

@media (max-width: 1024px) {
  #login-form-wrap {
    width: 65%;
  }
}

@media (max-width: 720px) {
  #login-form-wrap {
    width: 100%;
  }
}

.submit-area,
input {
  display: block;
  box-sizing: border-box;
  width: 100%;
  outline: none;
  height: 60px;
  line-height: 60px;
  border-radius: 4px;
}

input[type="text"],
input[type="password"] {
  width: 100%;
  padding: 0 0 0 10px;
  margin: 0;
  color: #8a8b8e;
  border: 1px solid #c2c0ca;
  font-style: normal;
  font-size: 16px;
  -webkit-appearance: none;
  -moz-appearance: none;
  appearance: none;
  position: relative;
  display: inline-block;
  background: none;
}
input[type="text"]:focus,
input[type="password"]:focus {
  border-color: #3ca9e2;
}
input[type="text"]:focus:invalid,
input[type="password"]:focus:invalid {
  color: #cc1e2b;
  border-color: #cc1e2b;
}

.validation {
  display: none;
  position: absolute;
  content: " ";
  height: 60px;
  width: 30px;
  right: 15px;
  top: 0px;
}

/* Draw ticks marks next to valid elements */
input[type="text"]:valid ~ .validation,
input[type="password"]:valid ~ .validation {
  display: block;
}
input[type="text"]:valid ~ .validation span,
input[type="password"]:valid ~ .validation span {
  background: #0c0;
  position: absolute;
  border-radius: 6px;
}
input[type="text"]:valid ~ .validation span:first-child,
input[type="password"]:valid ~ .validation span:first-child {
  top: 40px;
  left: 14px;
  width: 20px;
  height: 3px;
  -webkit-transform: rotate(-45deg);
  transform: rotate(-45deg);
}
input[type="text"]:valid ~ .validation span:last-child,
input[type="password"]:valid ~ .validation span:last-child {
  top: 45px;
  left: 8px;
  width: 11px;
  height: 3px;
  -webkit-transform: rotate(45deg);
  transform: rotate(45deg);
}

.submit-area,
input[type="submit"] {
  border: none;
  display: block;
  /* background-color: #3ca9e2; */
  background-color: var(--color-button);
  color: #fff;
  font-weight: bold;
  text-transform: uppercase;
  cursor: pointer;
  -webkit-transition: all 0.2s ease;
  transition: all 0.2s ease;
  font-size: 18px;
  position: relative;
  display: inline-block;
  cursor: pointer;
  text-align: center;
}
input[type="submit"]:hover {
  /* background-color: #329dd5; */
  /* background-color: #329dd5; */
  background-color: var(--vt-button-hover);
  -webkit-transition: all 0.2s ease;
  transition: all 0.2s ease;
}
</style>
