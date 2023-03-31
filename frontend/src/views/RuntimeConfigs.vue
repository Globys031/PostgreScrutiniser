<template>
  <div class="button-container">
    <UiButton type="submit" text="Get checks" />
    <UiButton type="info" text="submit" />
    <UiButton type="undo" text="submit" />
    <UiButton type="warning" text="submit" />
  </div>

  <!-- <UiNotification type="success" message="example message" duration="1000"/> -->
  <UiNotificationContainer
    ref="notificationContainer"
    type="success"
    message="example message"
  />
  <button @click="addNotification">Add notification</button>

  <UiCollapsible ref="collapsible">
    <template #title> Suggestions (%n) </template>
    <template #content>
      <UiCollapsibleItem @resize="(size: string) => setChildSize(size)">
        test
      </UiCollapsibleItem>
      <UiCollapsibleItem @resize="(size: string) => setChildSize(size)">
        test1
      </UiCollapsibleItem>
      <UiCollapsibleItem @resize="(size: string) => setChildSize(size)">
        test2
      </UiCollapsibleItem>
      <UiCollapsibleItem @resize="(size: string) => setChildSize(size)">
        test3
      </UiCollapsibleItem>
      <UiCollapsibleItem @resize="(size: string) => setChildSize(size)">
        test4
      </UiCollapsibleItem>
    </template>
  </UiCollapsible>

  <UiCollapsible ref="collapsible">
    <template #title> Checks that passed (%n) </template>
    <template #content>
      <UiCollapsibleItem @resize="(size: string) => setChildSize(size)">
        test
      </UiCollapsibleItem>
      <UiCollapsibleItem @resize="(size: string) => setChildSize(size)">
        test1
      </UiCollapsibleItem>
      <UiCollapsibleItem @resize="(size: string) => setChildSize(size)">
        test2
      </UiCollapsibleItem>
      <UiCollapsibleItem @resize="(size: string) => setChildSize(size)">
        test3
      </UiCollapsibleItem>
      <UiCollapsibleItem @resize="(size: string) => setChildSize(size)">
        test4
      </UiCollapsibleItem>
    </template>
  </UiCollapsible>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { ResourceApiFp } from "@/openapi/api/resource-config";
import { Configuration } from "@/openapi/configuration";
import { useSessionStore } from "@/stores/session";
import type { UiNotificationContainer } from "@/types/notification";
import type { UiCollapsibleComponent } from "@/types/collapsibleComponent";
import type { ErrorMessage } from "@/openapi/api/auth";

const sessionStore = useSessionStore();

// const postLogin = ResourceApiFp().getResourceConfigById;
const getConfigs = ResourceApiFp(
  new Configuration({
    accessToken: sessionStore.token,
  })
).getResourceConfigs;

getSuggestions();

const collapsible = ref<UiCollapsibleComponent | null>(null);
const notificationContainer = ref<UiNotificationContainer | null>(null);

function addNotification() {
  console.log("notificationContainer.value: ", notificationContainer.value);
  if (notificationContainer.value) {
    notificationContainer.value.addNotification(
      "success",
      "Operation successful",
      "New notification",
      3000
    );

    notificationContainer.value.addNotification(
      "error",
      "Operation successful",
      "New notification",
      3000
    );
  }
}

async function getSuggestions() {
  try {
    // isLoading.value = true;

    // Prepare login API request
    const getConfigsRequest = getConfigs();
    const getRequest = await getConfigsRequest;

    // Execute API request and save data in our session store
    const { data } = await getRequest();

    console.log("data: ", data);

    // apiResponse.value = "Request successful";
    // gotError.value = false;
    // isLoading.value = false;
  } catch (error) {
    console.error(error);
    // if (axios.isAxiosError(error)) {
    //   apiResponse.value = `Server response: ${
    //     (error.response?.data as ErrorMessage)?.error_message
    //   }`;
    //   gotError.value = true;
    // } else {
    //   apiResponse.value = "Something went wrong. See console tab for more info";
    //   console.error("Login Failed:", error);
    // }
    // isLoading.value = false;
  }
}

// After `UiCollapsibleItem` is resized, resize `UiCollapsible` as well
function setChildSize(size: string) {
  if (collapsible.value) {
    collapsible.value.resizeContentMaxHeight(size);
  }
}
</script>

<style scoped>
.button-container {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
}
</style>
