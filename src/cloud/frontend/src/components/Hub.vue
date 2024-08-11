<template>
  <div id="app" class="container mt-5 col-lg-6 col-md-8 col-sm-10">
    <h3>Ocelot Hub</h3>
    <div class="d-flex justify-content-end align-items-center mb-3">
      <span class="me-2">Logged in as: {{ user }}</span>
      <button style="margin-right: 5px" type="button" class="btn btn-primary" @click="visitCloud">Back to Cloud</button>
      <div id="dropdown" class="dropdown">
        <button class="btn btn-primary dropdown-toggle" type="button" id="settingsDropdown" data-bs-toggle="dropdown" aria-expanded="false">
          <i class="fas fa-cog"></i>
        </button>
        <ul class="dropdown-menu" aria-labelledby="settingsDropdown">
          <li id="button-logout"><a class="dropdown-item" href="#" @click="logout">Logout</a></li>
          <li id="button-change-password"><a class="dropdown-item" href="#" @click="redirectToChangePassword">Change Password</a></li>
          <li id="button-delete-account"><a class="dropdown-item text-danger" href="#" @click="showDeleteConfirmation = true">Delete Account</a></li>
        </ul>
      </div>
    </div>
    <br>
    <h4>App and Tag Management</h4>
    <div class="d-flex justify-content-center mb-3">
      <div class="col-6">
        <input id="input-app" v-model="newApp" class="form-control" placeholder="App" required />
      </div>
    </div>
    <button id="button-create-app" @click="createApp" class="btn btn-primary">Create App</button>
    <button id="button-delete-app" @click="deleteApp" class="btn btn-danger ms-2">Delete App</button>

    <h5>App List:</h5>
    <div>
      <ul id="app-list">
        <li v-for="app in appList" :key="app">{{ app }}</li>
      </ul>
    </div>
    <router-view />

    <div v-if="showDeleteConfirmation" class="modal fade show" style="display: block;" tabindex="-1">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Confirm Account Deletion</h5>
            <button type="button" class="btn-close" @click="showDeleteConfirmation = false" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <p>Are you sure you want to delete your account?</p>
          </div>
          <div class="modal-footer">
            <button id="button-delete-cancel" type="button" class="btn btn-secondary" @click="showDeleteConfirmation = false">Cancel</button>
            <button id="button-delete-confirmation" type="button" class="btn btn-danger" @click="confirmDeleteAccount">Confirm</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>


<script lang="ts">
import { defineComponent, ref, onMounted } from 'vue';
import axios from "axios";
import router from "@/router";

// TODO If app list is empty, then show an according message
export default defineComponent({
  name: 'HubComponent',
  setup() {
    const user = ref<string | null>(null);
    const showDeleteConfirmation = ref(false);
    const app = ref('');
    const appList = ref<string[]>([]);

    const checkAuth = async () => {
      try {
        const url = 'http://localhost:8082';
        const response = await axios.get(url + "/auth-check");
        if (response.status === 200) {
          user.value = response.data.value;
        }
      } catch (error) {
        redirectToLogin();
      }
      // TODO add error?
    };

    const logout = async () => {
      try {
        const url = 'http://localhost:8082';
        await axios.get(url + "/logout");
        user.value = "";
        redirectToLogin();
      } catch (error) {
        alert(error);
      }
    };

    const deleteAccount = async () => {
      try {
        const url = 'http://localhost:8082';
        await axios.delete(url + "/user");
        user.value = "";
        redirectToLogin();
      } catch (error) {
        alert(error);
      }
    };

    const confirmDeleteAccount = async () => {
      showDeleteConfirmation.value = false;
      await deleteAccount();
    };

    const redirectToLogin = () => {
      router.push('/hub/login');
    };

    const redirectToChangePassword = () => {
      router.push('/hub/change-password');
    };

    const createApp = async () => {
      const url = 'http://localhost:8082';
      try {
        await axios.post(url + '/apps', { value: app.value });
      } catch (error) {
        // TODO correctly interpret error, so that backend message is displayed.
        alert("app creation error: " + error)
      }
      app.value = ""
      getApps()
    };

    const deleteApp = async () => {
      const url = 'http://localhost:8082';
      try {
        await axios.delete(url + '/apps', {data: { value: app.value }});
      } catch (error) {
        // TODO correctly interpret error, so that backend message is displayed.
        alert("app deletion error: " + error)
      }
      app.value = ""
      getApps()
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
      app.value = ""
    };

    const visitCloud = () => {
      router.push('/');
    };

    onMounted(() => {
      checkAuth();
      getApps();
    });

    return {
      user,
      logout,
      showDeleteConfirmation,
      deleteAccount,
      confirmDeleteAccount,
      redirectToChangePassword,
      visitCloud,
      newApp: app,
      createApp,
      appList,
      deleteApp,
    };
  },
});
</script>


TODO: Hub
* CreateApp
* GetTags
* DeleteApp
* UploadTag
* DownloadTag
* DeleteTag

TODO: Cloud
* FindApps
* DownloadTag

TODO: Input validation in frontend, so that user know why their input is not accepted