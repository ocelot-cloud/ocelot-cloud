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
      message="Are you sure you want to delete this app?"
  ></HubDeletionConfirmationDialog>
</template>

<script lang="ts">
import {defineComponent, onMounted, ref} from "vue";
import router from "@/router";
import {session} from "@/components/hub/shared";
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
      const url = 'http://localhost:8082';
      try {
        await axios.post(url + '/apps', { value: newAppToCreate.value });
      } catch (error) {
        // TODO correctly interpret error, so that backend message is displayed.
        alert("app creation error: " + error)
      }
      newAppToCreate.value = ""
      await getApps()
    };

    const getApps = async () => {
      const url = 'http://localhost:8082';
      try {
        const response = await axios.get(url + '/apps');
        if (response.status === 200) {
          appList.value = response.data as string[];
          console.log("received apps: ", appList.value)
        }
      } catch (error) {
        console.log("todo")
      }
    };

    const deleteApp = async () => {
      try {
        const url = 'http://localhost:8082';
        await axios.delete(url + '/apps', {data: { value: selectedApp.value }});
        user.value = "";
      } catch (error) {
        alert(error);
      }
      selectedApp.value = ""
      await getApps()
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