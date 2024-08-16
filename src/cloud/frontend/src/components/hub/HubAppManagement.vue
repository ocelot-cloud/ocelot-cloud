<template>
  <div v-if="!isEditingTags">
    <h4 >App Management</h4>
    <div class="d-flex justify-content-center mb-3">
      <div class="col-6">
        <input id="input-app" v-model="newAppToCreate" class="form-control" placeholder="Name of New App" required />
      </div>
    </div>
    <button id="button-create-app" @click="createApp" class="btn btn-primary">Create App</button>
    <br>
    <br>
    <h4>App List</h4>
    <p v-if="appList == null"> (no apps created yet) </p>
    <div class="d-flex justify-content-center">
      <ul id="app-list" class="list-group">
        <li
            v-for="(app, index) in appList"
            :key="app"
            class="list-group-item"
            :class="{ active: selectedApp === app }"
            @click="selectApp(app)"
            style="cursor: pointer;"
        >
          {{ index + 1 }}) {{ app }}
        </li>
      </ul>
    </div>
    <br>
    <div v-if="appList != null && selectedApp != ''">
      <h4>App Operations</h4>
      <button id="button-edit-tags" @click="goToTagManagement()" class="btn btn-warning me-2">Edit Tags</button>
      <!-- TODO There should be a confirmation dialog previously -->
      <button id="button-delete-app" @click="showDeleteConfirmation = true" class="btn btn-danger ms-2">Delete</button>
    </div>
    <p>TODO: {{showDeleteConfirmation}}</p>
  </div>

  <HubDeletionConfirmationDialog
      v-model:visible="showDeleteConfirmation"
      :on-confirm="deleteApp"
      messageSuffix="this app?"
  ></HubDeletionConfirmationDialog>
</template>

<script lang="ts">
import {defineComponent, onMounted, ref} from "vue";
import router from "@/router";
import {doRequest, session} from "@/components/hub/shared";
import axios from "axios";
import HubDeletionConfirmationDialog from "@/components/hub/HubDeletionConfirmationDialog.vue";

export default defineComponent({
  name: 'HubAppManagement',
  components: {HubDeletionConfirmationDialog},

  setup() {
    const user = ref("");
    const showDeleteConfirmation = ref(false);
    const newAppToCreate = ref('');
    const appList = ref<string[]>([]);
    const selectedApp = ref("");
    const isEditingTags = ref(false);

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
      await doRequest("/apps", { value: newAppToCreate.value })
      await getApps()
      newAppToCreate.value = ""
    };

    const getApps = async () => {
      const response = await doRequest("/apps/get-list", null)
      if (response != null) {
        appList.value = response.data as string[];
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
      showDeleteConfirmation
    }
  },
})
</script>