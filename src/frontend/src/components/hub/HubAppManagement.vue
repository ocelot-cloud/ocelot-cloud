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
              :key="app.name"
              class="list-group-item d-flex justify-content-between align-items-center bg-secondary bg-opacity-25 text-white"
              :class="{ active: selectedApp.id === app.id }"
              @click="selectApp(app)"
              style="cursor: pointer;"
          >
            <span>{{ app.name }}</span>
            <i v-if="selectedApp.id === app.id" class="bi bi-check-circle-fill text-success"></i>
          </li>
        </ul>
      </div>
    </div>

    <div v-if="appList && selectedApp.id != -1" class="app-operations d-flex justify-content-end">
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
import router, {hubSession} from "@/router";
import HubDeletionConfirmationDialog from "@/components/shared/HubDeletionConfirmationDialog.vue";
import ValidatedInput from "@/components/shared/ValidatedInput.vue";
import {doHubRequest} from "@/components/shared/shared";

class App {
  name: string;
  id: number;

  constructor(name: string, id: number) {
    this.name = name;
    this.id = id;
  }
}

export default defineComponent({
  name: 'HubAppManagement',
  components: {ValidatedInput, HubDeletionConfirmationDialog},

  setup() {
    const user = ref("");
    const showDeleteConfirmation = ref(false);
    const newAppToCreate = ref('');
    const appList = ref<App[]>([]);
    const selectedApp = ref<App>(new App("", -1));
    const isEditingTags = ref(false);
    const submitted = ref(false);

    const selectApp = (app: App) => {
      if (selectedApp.value.id == app.id) {
        selectedApp.value.id = -1
      } else {
        selectedApp.value.id = app.id;
        selectedApp.value.name = app.name;
      }
    };

    const goToTagManagement = () => {
      router.push({ path: '/hub/tag-management', query: { user: user.value, appName: selectedApp.value.name, appId: selectedApp.value.id } });
    }

    const createApp = async () => {
      submitted.value = true
      const response = await doHubRequest("/apps/create", { value: newAppToCreate.value })
      if (response) {
        await getApps()
        newAppToCreate.value = ""
        submitted.value = false
      }
    };

    const getApps = async () => {
      const response = await doHubRequest("/apps/get-list", null)
      if (response != null && response.data) {
        appList.value = response.data as App[];
        if (appList.value.length > 0) {
          appList.value.sort((a: App, b: App) => a.name.localeCompare(b.name));
        }
      }
    };

    const deleteApp = async () => {
      const response = await doHubRequest("/apps/delete", { value: selectedApp.value.id })
      if (response != null) {
        await getApps()
        appList.value = appList.value.filter(app => app.id !== selectedApp.value.id);
        selectedApp.value.id = -1
        showDeleteConfirmation.value = false
      }
    };

    const confirmDeleteAccount = async () => {
      showDeleteConfirmation.value = false;
      await deleteApp();
    };

    onMounted(() => {
      user.value = hubSession.user
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