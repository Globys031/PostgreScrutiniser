<template>
  <div class="button-container">
    <UiButton @click="getSuggestions" type="info" text="Get checks" />
    <UiButton
      @click="resetConfigs"
      type="warning"
      text="Reset Configurations"
    />
    <UiButton
      v-if="suggestions.length !== 0"
      @click="applySuggestions(suggestions)"
      type="submit"
      text="Apply all suggestions"
    />
  </div>

  <template v-if="configChecks.length !== 0">
    <UiCollapsible ref="collapsibleSuggestions">
      <template #title> Suggestions (%n) </template>
      <template #content>
        <UiCollapsibleItem
          v-for="(item, index) in suggestions"
          :key="index"
          @resize="(size: string) => setChildSize(collapsibleSuggestions, size)"
        >
          <template #title>
            <span v-text="item.Name" />
          </template>

          <template #content>
            <div v-if="item.GotError" class="error-container">
              An error occurred when trying to check this configuration
              parameter. See application error logs for more information
            </div>
            <div v-else>
              <div class="suggestion-container">
                <div>
                  <span class="value-tag">Current value: </span>
                  <span v-text="`${item.Value} ${item.Unit}`" />
                </div>
                <div>
                  <span class="value-tag">Suggested value: </span>
                  <span v-text="`${item.SuggestedValue} ${item.Unit}`" />
                </div>
                <UiButton
                  @click="applySuggestions([item])"
                  type="submit"
                  text="Apply suggestion"
                />
              </div>
              <div class="details-container">
                <span v-text="item.Details" />
              </div>
            </div>
          </template>
        </UiCollapsibleItem>

        <span v-if="suggestions.length === 0" class="empty-table-message">
          No suggestions were made
        </span>
      </template>
    </UiCollapsible>

    <UiCollapsible ref="collapsiblePassedChecks">
      <template #title> Checks that passed (%n) </template>
      <template #content>
        <UiCollapsibleItem
          v-for="(item, index) in passedChecks"
          :key="index"
          @resize="(size: string) => setChildSize(collapsiblePassedChecks,size)"
        >
          <template #title>
            <span v-text="item.Name" />
          </template>

          <template #content>
            <div v-if="item.GotError" class="error-container">
              An error occurred when trying to check this configuration
              parameter. See application error logs for more information
            </div>
            <div v-else>
              <div class="suggestion-container">
                <div>
                  <span class="value-tag">Current value: </span>
                  <span v-text="`${item.Value} ${item.Unit}`" />
                </div>
              </div>
              <div class="details-container">
                <span v-text="item.Details" />
              </div>
            </div>
          </template>
        </UiCollapsibleItem>
        <span v-if="passedChecks.length === 0" class="empty-table-message">
          No checks have passed
        </span>
      </template>
    </UiCollapsible>
  </template>
  <div v-else class="no-data-message">
    No data yet. Click "Get Checks" to get started
  </div>

  <UiNotificationContainer
    ref="notificationContainer"
    type="success"
    message="example message"
  />
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import { ResourceApiFp } from "@/openapi/api/resource-config";
import { Configuration } from "@/openapi/configuration";
import { useSessionStore } from "@/stores/session";
import {
  displayNotification,
  displayNotificationError,
} from "@/composables/notifications";
import type { ResourceConfigPatchSchema } from "@/openapi/api/resource-config";
import type { ResourceConfigPascalCase } from "@/openapi/typeInference";
import type { UiNotificationContainer } from "@/types/notification";
import type { UiCollapsibleComponent } from "@/types/collapsibleComponent";

const sessionStore = useSessionStore();

const resourceApi = ResourceApiFp(
  new Configuration({
    accessToken: sessionStore.token,
  })
);

// Html ref tag references:
const collapsibleSuggestions = ref<UiCollapsibleComponent | null>(null);
const collapsiblePassedChecks = ref<UiCollapsibleComponent | null>(null);
const notificationContainer = ref<UiNotificationContainer | null>(null);

// Data to be displayed
const configChecks = ref<ResourceConfigPascalCase[]>([]);

const suggestions = computed(() => {
  const values = Object.values(configChecks.value);
  return values.filter(
    (suggestion: ResourceConfigPascalCase) => suggestion.SuggestedValue
  );
});

const passedChecks = computed(() => {
  const values = Object.values(configChecks.value);
  return values.filter((suggestion) => !suggestion.SuggestedValue);
});

async function getSuggestions() {
  try {
    const getRequest = await resourceApi.getResourceConfigs();
    const { data } = await getRequest();
    displayNotification(
      notificationContainer,
      "success",
      "Suggestions will be displayed inside collapsible tables"
    );

    configChecks.value = data as ResourceConfigPascalCase[];
  } catch (error) {
    configChecks.value = [];
    displayError(error);
  }
}

async function applySuggestions(suggestions: ResourceConfigPascalCase[]) {
  try {
    // Prepare suggestions for patch request
    const formattedSuggestions = suggestions.map(
      (suggestion: ResourceConfigPascalCase) => ({
        name: suggestion.Name,
        suggested_value: suggestion.SuggestedValue,
      })
    );

    // Execute patch api request
    const patchRequest = await resourceApi.patchResourceConfigs(
      formattedSuggestions as ResourceConfigPatchSchema[]
    );
    await patchRequest();
    displayNotification(
      notificationContainer,
      "success",
      "Suggestions applied"
    );

    // After patching requests, renew table (suggestion) data.
    getSuggestions();
  } catch (error) {
    displayError(error);
    configChecks.value = []; // reset table to empty
  }
}

async function resetConfigs() {
  try {
    const deleteConfigsRequest = resourceApi.deleteResourceConfigs();
    const deleteRequest = await deleteConfigsRequest;
    await deleteRequest();

    // After patching requests, renew table (suggestion) data.
    getSuggestions();
  } catch (error) {
    displayError(error);
  }
}

function displayError(error: unknown) {
  error instanceof Error
    ? displayNotificationError(notificationContainer, error)
    : console.error(error);
}

// After `UiCollapsibleItem` is resized, resize `UiCollapsible` as well
function setChildSize(
  collapsible: UiCollapsibleComponent | null,
  size: string
) {
  collapsible
    ? collapsible.resizeContentMaxHeight(size)
    : console.error("collapsible component is null");
}
</script>

<style scoped>
.button-container {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
}

.empty-table-message {
  padding: 10px;
}

.suggestion-container {
  padding-top: 10px;
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
}

.error-container {
  color: var(--vt-error-color);
  background-color: var(--color-error-background);
}

.error-container,
.suggestion-container,
.details-container {
  margin-bottom: 1rem;
}
.value-tag {
  font-weight: bold;
}

.no-data-message {
  padding: 1rem;
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-radius: 4px;
  font-size: 1.1em;
  font-weight: bold;
  text-align: center;
  color: #444;
  background-color: #f0f0f0;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  margin: 2rem 0;
}
</style>
