<template>
  <div class="entity-management-container p-4 shadow-sm bg-dark rounded">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h4>App Management</h4>
      <button id="button-create-app" @click="createApp" class="btn btn-primary">Create App</button>
    </div>

    <ValidatedInput :submitted="submitted" validation-type="app" v-model="newAppToCreate"></ValidatedInput>

    <div class="app-list-section mb-4">
      <h4>App List</h4>
      <p v-if="!appList || appList.length === 0">(No apps created yet)</p>

      <div class="d-flex justify-content-center">
        <ul id="app-list" class="list-group w-100">
          <li
              v-for="app in appList"
              :key="app"
              class="list-group-item d-flex justify-content-between align-items-center bg-secondary bg-opacity-25 text-white"
              :class="{ active: selectedApp === app }"
              @click="selectApp(app)"
              style="cursor: pointer;"
          >
            <span>{{ app }}</span>
            <i v-if="selectedApp === app" class="bi bi-check-circle-fill text-success"></i>
          </li>
        </ul>
      </div>
    </div>

    <div v-if="appList && selectedApp" class="app-operations d-flex justify-content-end">
      <button id="button-edit-tags" @click="goToTagManagement()" class="btn btn-primary me-2">Edit Tags</button>
      <button id="button-delete-app" @click="showDeleteConfirmation = true" class="btn btn-danger">Delete</button>
    </div>
  </div>

  <HubDeletionConfirmationDialog
      v-model:visible="showDeleteConfirmation"
      :on-confirm="deleteApp"
      messageSuffix="this app?"
  />
</template>

<script lang="ts">
import {defineComponent, onMounted, ref} from "vue";
import router from "@/router";
import {doRequest, session} from "@/components/hub/shared";
import HubDeletionConfirmationDialog from "@/components/hub/HubDeletionConfirmationDialog.vue";
import ValidatedInput from "@/components/hub/ValidatedInput.vue";

export default defineComponent({
  name: 'HubAppManagement',
  components: {ValidatedInput, HubDeletionConfirmationDialog},

  setup() {
    const user = ref("");
    const showDeleteConfirmation = ref(false);
    const newAppToCreate = ref('');
    const appList = ref<string[]>([]);
    const selectedApp = ref("");
    const isEditingTags = ref(false);
    const submitted = ref(false);

    const selectApp = (app: string) => {
      if (selectedApp.value == app) {
        selectedApp.value = ""
      } else {
        selectedApp.value = app;
      }
    };

    const goToTagManagement = () => {
      router.push({ path: '/hub/tag-management', query: { user: user.value, app: selectedApp.value } });
    }

    const createApp = async () => {
      submitted.value = true
      const response = await doRequest("/apps/create", { value: newAppToCreate.value })
      if (response) {
        await getApps()
        newAppToCreate.value = ""
      }
    };

    const getApps = async () => {
      const response = await doRequest("/apps/get-list", null)
      if (response != null) {
        appList.value = response.data as string[];
        if (appList.value != null) {
          appList.value.sort()
        }
      }
    };

    const deleteApp = async () => {
      await doRequest("/apps/delete", { value: selectedApp.value })
      await getApps()
      selectedApp.value = ""
      showDeleteConfirmation.value = false
    };

    const confirmDeleteAccount = async () => {
      showDeleteConfirmation.value = false;
      await deleteApp();
    };

    onMounted(() => {
      user.value = session.user
      getApps();
    });

    return {
      isEditingTags,
      appList,
      selectedApp,
      newAppToCreate,
      selectApp,
      goToTagManagement,
      createApp,
      deleteApp,
      confirmDeleteAccount,
      showDeleteConfirmation,
      submitted
    }
  },
})
</script>