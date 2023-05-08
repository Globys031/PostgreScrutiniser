<template>
  <div class="g-button-container">
    <UiButton
      :class="{ disabled: isLoading }"
      @click="isLoading ? disabledButtonsNotification() : getBackups()"
      type="info"
      text="List backups"
    />
    <UiButton
      :class="{ disabled: isLoading }"
      @click="
        isLoading
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
    <UiSpinner v-if="isLoading" />
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
          <div class="content-container">
            <UiButton
              @click="
                isLoading
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
                isLoading
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
  <div v-else class="g-no-data-message">
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
import { displayNotification, displayError } from "@/composables/notifications";
import { useModal } from "@/composables/modal";
import { BackupApiFp } from "@/openapi/api/file";
import FileDiff from "@/components/FileDiff.vue";
import type { BackupFile, FileDiffLine } from "@/openapi/api/file";
import type { UiNotificationContainer } from "@/types/notification";

const sessionStore = useSessionStore();
const modal = useModal();

const backupApi = BackupApiFp(
  new Configuration({
    basePath: sessionStore.baseAPIPath,
    accessToken: sessionStore.token,
  })
);

// Html ref tag references:
const notificationContainer = ref<UiNotificationContainer | null>(null);

// Data to be displayed
const backups = ref<BackupFile[]>([]);
const fileDiff = ref<FileDiffLine[]>([]);

const isLoading = ref<boolean>(false);

async function getBackups() {
  try {
    isLoading.value = true;
    const getRequest = await backupApi.getBackups();
    const { data } = await getRequest();

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

    backups.value = (data as unknown as BackupFile[]).reverse();
  } catch (error) {
    backups.value = [];
    displayError(notificationContainer, error);
  }
  isLoading.value = false;
}

async function restoreBackup(backupName: string) {
  if (!backupName) {
    displayError(
      notificationContainer,
      "Could not get backup name during restore"
    );
    return;
  }
  try {
    isLoading.value = true;
    const restoreRequest = await backupApi.putBackup(backupName);
    await restoreRequest();
    displayNotification(
      notificationContainer,
      "success",
      "Backups have been restored"
    );

    // After restoring a backup, get new list of backups and resize table
    getBackups();
  } catch (error) {
    displayError(notificationContainer, error);
  }
  isLoading.value = false;
}

async function deleteBackups() {
  try {
    isLoading.value = true;
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
    displayError(notificationContainer, error);
  }
  isLoading.value = false;
}

async function deleteBackup(backupName: string) {
  if (!backupName) {
    displayError(
      notificationContainer,
      "Could not get backup name during removal"
    );
    return;
  }
  try {
    isLoading.value = true;
    const deleteRequest = await backupApi.deleteBackup(backupName);
    await deleteRequest();
    displayNotification(
      notificationContainer,
      "success",
      "Backups have been deleted"
    );

    // After deleting a backup, get new list of backups and resize table
    getBackups();
  } catch (error) {
    displayError(notificationContainer, error);
  }
  isLoading.value = false;
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
.content-container {
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

.disabled {
  background: #ccc;
  cursor: not-allowed;
}
.disabled:hover {
  color: #999;
}
</style>
