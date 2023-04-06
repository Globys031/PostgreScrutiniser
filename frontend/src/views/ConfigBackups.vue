<template>
  <div class="button-container">
    <UiButton
      :class="{ disabled: isLoadingList }"
      @click="isLoadingList ? disabledButtonsNotification() : getBackups()"
      type="info"
      text="List backups"
    />
    <UiButton
      :class="{ disabled: isLoadingList }"
      @click="
        isLoadingList
          ? disabledButtonsNotification()
          : modal.triggerModal(
              'Delete backups',
              'Are you sure you want to delete all backups? This action is irreversible.',
              'Delete',
              deleteBackups
            )
      "
      type="warning"
      text="Delete backups"
    />
    <UiSpinner v-if="isLoadingList" />
  </div>

  <template v-if="backups.length !== 0">
    <div class="info-message">
      <span class="green-text">Green text</span> is what would be added should
      backup be applied. <span class="red-text">Red</span> is what would be
      removed.
    </div>
    <div class="content">
      <!-- <template v-if="backups.length !== 0"> -->
      <UiCollapsibleItem v-for="(backup, index) in backups" :key="index">
        <template #title>
          <div class="title-container">
            <span v-text="backup.name" />
            <span v-text="backup.time" />
          </div>
        </template>

        <template #content>
          <FileDiff v-show="fileDiff" :diff="backup.diff" />
          <div class="suggestion-container">
            <UiButton
              @click="
                isLoadingList
                  ? disabledButtonsNotification()
                  : modal.triggerModal(
                      'Restore backup',
                      `Are you sure you want to restore ${backup.name}? A backup will be created of the current postgresql.auto.conf file and current configuration will be replaced with content from ${backup.name}`,
                      'Restore',
                      () => restoreBackup(backup.name)
                    )
              "
              type="submit"
              text="Restore backup"
            />
            <UiButton
              @click="
                isLoadingList
                  ? disabledButtonsNotification()
                  : modal.triggerModal(
                      'Delete backup file',
                      `Are you sure you want to delete ${backup.name}? This action is irreversible.`,
                      'Delete',
                      () => deleteBackup(backup.name)
                    )
              "
              type="warning"
              text="Delete backup"
            />
          </div>
        </template>
      </UiCollapsibleItem>
    </div>
  </template>
  <div v-else class="no-data-message">
    No data yet. Click "List backups" to get started
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
import { ref } from "vue";
import { Configuration } from "@/openapi/configuration";
import { useSessionStore } from "@/stores/session";
import {
  displayNotification,
  displayNotificationError,
} from "@/composables/notifications";
import { useModal } from "@/composables/modal";
import { BackupApiFp } from "@/openapi/api/file";
import FileDiff from "@/components/FileDiff.vue";
import type { BackupFile, FileDiffLine } from "@/openapi/api/file";
import type { UiNotificationContainer } from "@/types/notification";

const sessionStore = useSessionStore();
const modal = useModal();

const backupApi = BackupApiFp(
  new Configuration({
    accessToken: sessionStore.token,
  })
);

// Html ref tag references:
const notificationContainer = ref<UiNotificationContainer | null>(null);

// Data to be displayed
const backups = ref<BackupFile[]>([]);
const fileDiff = ref<FileDiffLine[]>([]);

const isLoadingList = ref<boolean>(false);

async function getBackups() {
  try {
    const getRequest = await backupApi.getBackups();
    isLoadingList.value = true;
    const { data } = await getRequest();
    isLoadingList.value = false;

    // could be null after deleting all backups
    if (!data) {
      displayNotification(
        notificationContainer,
        "success",
        "No backups were found"
      );
      backups.value = [];
      return;
    }

    transformDate(data as unknown as BackupFile[]); // change data.time to format of YYYY-MM-DD

    displayNotification(
      notificationContainer,
      "success",
      "Backups will be displayed on the screeen"
    );

    backups.value = data as unknown as BackupFile[];
  } catch (error) {
    backups.value = [];
    isLoadingList.value = false;
    displayError(error);
  }
}

async function restoreBackup(backupName: string) {
  if (!backupName) {
    displayError("Could not get backup name during restore");
    return;
  }
  try {
    const restoreRequest = await backupApi.putBackup(backupName);
    await restoreRequest();
    displayNotification(
      notificationContainer,
      "success",
      "Backups has been restored"
    );

    // After restoring a backup, get new list of backups and resize table
    getBackups();
  } catch (error) {
    displayError(error);
  }
}

// Function for removing all backup files (TO DO: needs a modal window for confirmation)
async function deleteBackups() {
  try {
    const deleteAllRequest = await backupApi.deleteBackups();
    await deleteAllRequest();
    displayNotification(
      notificationContainer,
      "success",
      "Backups have been deleted"
    );

    // After deleting backups, get new list of backups and resize table
    getBackups();
  } catch (error) {
    displayError(error);
  }
}

// Function for removing a single backup file (TO DO: needs a modal window for confirmation)
async function deleteBackup(backupName: string) {
  if (!backupName) {
    displayError("Could not get backup name during removal");
    return;
  }
  try {
    const deleteRequest = await backupApi.deleteBackup(backupName);
    await deleteRequest();
    displayNotification(
      notificationContainer,
      "success",
      "Backups has been deleted"
    );

    // After deleting a backup, get new list of backups and resize table
    getBackups();
  } catch (error) {
    displayError(error);
  }
}

// REUSE (This means that RuntimeConfigs.vue also has same exact code)
function displayError(error: unknown) {
  error instanceof Error
    ? displayNotificationError(notificationContainer, error)
    : console.error(error);
}

// Function for transforming data.time to a format of YYYY-MM-DD
function transformDate(data: BackupFile[]) {
  data.map((backup) => {
    if (backup.time) {
      const datetime = backup.time.split("+")[0];
      backup.time = datetime.replace("T", " ");
    } else {
      console.error("Backup file does not contain time when it was created");
    }
  });
}

function disabledButtonsNotification() {
  displayNotification(
    notificationContainer,
    "error",
    "Fetching latest list of backups, please wait"
  );
}
</script>

<style scoped>
/* Iskelt i atskira kad butu reusable (iskelt i base.css) */
.button-container {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
}

/* iskelt i atskira kad butu reusable */
.empty-table-message {
  padding: 10px;
}

/* Iskelt i atskira kad butu reusable */
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

/* Iskelt i atskira kad butu reusable (pervadint i "collapsible-children-container") */
.suggestion-container {
  padding-top: 10px;
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  margin-bottom: 1rem;
}

.title-container {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
}

/* iskelt i atskira kad butu reusable */
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

.content {
  box-shadow: rgba(0, 0, 0, 0.24) 0px 3px 8px;
}

.info-message {
  font-size: 14px;
  color: #333;
}

.green-text {
  color: #2ecc71;
  font-weight: bold;
}

.red-text {
  color: #e74c3c;
  font-weight: bold;
}

/* TO DO: make reusable */
/**Disabled button**/
.disabled {
  background: #ccc;
  cursor: not-allowed;
}
.disabled:hover {
  color: #999;
}
</style>
