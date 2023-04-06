<template>
  <div class="button-container">
    <UiButton @click="getSuggestions" type="info" text="Get checks" />
    <UiButton
      :disabled="buttonsDisabled"
      :class="{ disabled: buttonsDisabled }"
      @click="
        buttonsDisabled
          ? disabledButtonsNotification()
          : modal.triggerModal(
              'Reset Configurations',
              `Are you sure you want to reset configurations? A backup will be created of the current postgresql.auto.conf file and current configurations within that file will be wiped out`,
              'Restore',
              resetConfigs
            )
      "
      type="warning"
      text="Reset Configurations"
    />
    <UiButton
      v-show="suggestions.length !== 0"
      :disabled="buttonsDisabled"
      :class="{ disabled: buttonsDisabled }"
      @click="
        buttonsDisabled
          ? disabledButtonsNotification()
          : applySuggestions(suggestions)
      "
      type="submit"
      text="Apply all suggestions"
    />
    <UiSpinner v-if="isLoadingConfigs" />
  </div>

  <template v-if="configChecks.length !== 0">
    <UiCollapsible ref="collapsibleSuggestions" :isSuggestions="true">
      <template #title>
        <span v-text="`Suggestions (${suggestions.length})`" />
      </template>
      <template #content>
        <UiCollapsibleItem
          v-for="(suggestion, index) in suggestions"
          :key="index"
          @resize="setChildSize(collapsibleSuggestions)"
        >
          <template #title>
            <span v-text="suggestion.Name" />
          </template>

          <template #content>
            <div v-if="suggestion.GotError" class="error-container">
              An error occurred when trying to check this configuration
              parameter. See application error logs for more information
            </div>
            <div v-else>
              <div class="suggestion-container">
                <div>
                  <span class="value-tag">Current value: </span>
                  <span v-text="`${suggestion.Value} ${suggestion.Unit}`" />
                </div>
                <div>
                  <span class="value-tag">Suggested value: </span>
                  <span
                    v-text="`${suggestion.SuggestedValue} ${suggestion.Unit}`"
                  />
                </div>
                <UiButton
                  :disabled="buttonsDisabled"
                  :class="{ disabled: buttonsDisabled }"
                  @click="
                    buttonsDisabled
                      ? disabledButtonsNotification()
                      : applySuggestions([suggestion])
                  "
                  type="submit"
                  text="Apply suggestion"
                />
              </div>
              <div class="details-container">
                <span v-text="suggestion.Details" />
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
      <template #title>
        <span v-text="`Checks that passed  (${passedChecks.length})`" />
      </template>
      <template #content>
        <UiCollapsibleItem
          v-for="(passedCheck, index) in passedChecks"
          :key="index"
          @resize="(size: string) => setChildSize(collapsiblePassedChecks)"
        >
          <template #title>
            <span v-text="passedCheck.Name" />
          </template>

          <template #content>
            <div v-if="passedCheck.GotError" class="error-container">
              An error occurred when trying to check this configuration
              parameter. See application error logs for more information
            </div>
            <div v-else>
              <div class="suggestion-container">
                <div>
                  <span class="value-tag">Current value: </span>
                  <span v-text="`${passedCheck.Value} ${passedCheck.Unit}`" />
                </div>
              </div>
              <div class="details-container">
                <span v-text="passedCheck.Details" />
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

  <UiModal
    v-if="modal.showModal.value"
    :title="modal.modalTitle.value"
    :content="modal.modalContent.value"
    :buttonText="modal.modalButtonText.value"
    :functionToRun="modal.modalFunction.value"
    @close="modal.showModal.value = !modal.showModal"
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
import { useModal } from "@/composables/modal";
import type { ResourceConfigPatchSchema } from "@/openapi/api/resource-config";
import type { ResourceConfigPascalCase } from "@/openapi/typeInference";
import type { UiNotificationContainer } from "@/types/notification";
import type { UiCollapsibleComponent } from "@/types/collapsibleComponent";

const sessionStore = useSessionStore();
const modal = useModal();

const resourceApi = ResourceApiFp(
  new Configuration({
    accessToken: sessionStore.token,
  })
);

const timeout = 3000; // 3 seconds
const buttonsDisabled = ref<boolean>(false); // disable for 3 seconds after applying configs

// Html ref tag references:
const collapsibleSuggestions = ref<UiCollapsibleComponent | null>(null);
const collapsiblePassedChecks = ref<UiCollapsibleComponent | null>(null);
const notificationContainer = ref<UiNotificationContainer | null>(null);

// Data to be displayed
const configChecks = ref<ResourceConfigPascalCase[]>([]);

const isLoadingConfigs = ref<boolean>(false);

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
    isLoadingConfigs.value = true;
    const { data } = await getRequest();
    displayNotification(
      notificationContainer,
      "success",
      "Suggestions will be displayed inside collapsible tables"
    );

    configChecks.value = data as unknown as ResourceConfigPascalCase[];
    isLoadingConfigs.value = false;
  } catch (error) {
    configChecks.value = [];
    isLoadingConfigs.value = false;
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
    refreshChecks();
    disableButtons();
  } catch (error) {
    displayError(error);
    configChecks.value = []; // reset table to empty
  }
}

//  TO DO: needs a modal window for confirmation
async function resetConfigs() {
  console.log("resetConfigs()");
  try {
    const deleteConfigsRequest = resourceApi.deleteResourceConfigs();
    const deleteRequest = await deleteConfigsRequest;
    await deleteRequest();

    // After patching requests, renew table (suggestion) data.
    refreshChecks();
    disableButtons();
  } catch (error) {
    displayError(error);
  }
}

// After need configurations settings are applied, refresh settings table
function refreshChecks() {
  getSuggestions();
  setChildSize(collapsibleSuggestions.value);
  setChildSize(collapsiblePassedChecks.value);
}

function displayError(error: unknown) {
  error instanceof Error
    ? displayNotificationError(notificationContainer, error)
    : console.error(error);
}

// After `UiCollapsibleItem` is resized, resize `UiCollapsible` as well
function setChildSize(collapsible: UiCollapsibleComponent | null) {
  collapsible
    ? collapsible.resizeContentMaxHeight()
    : console.error("collapsible component is null");
}

// After applying configs, it is necessary that user cannot apply suggestions for a
// short period of time to avoid systemd errors caused by too many restart of psotgresql
function disableButtons() {
  buttonsDisabled.value = true;
  setTimeout(() => {
    buttonsDisabled.value = false;
  }, timeout);
}

function disabledButtonsNotification() {
  displayNotification(
    notificationContainer,
    "error",
    "Cannot click buttons that would cause PostgreSql to reload. Some configurations require restarting PostgreSql. To avoid crashing systemd instance, the amount of requests is limited to one every 3 seconds.",
    10000
  );
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

/**Disabled button**/
.disabled {
  background: #ccc;
  cursor: not-allowed;
}
.disabled:hover {
  color: #999;
}
</style>
