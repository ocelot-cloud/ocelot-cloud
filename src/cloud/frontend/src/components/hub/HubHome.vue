<template>
  <div id="app" class="container mt-5 col-lg-6 col-md-8 col-sm-10">
    <h3>Ocelot Hub</h3>
    <div class="d-flex justify-content-end align-items-center mb-3">
      <span class="me-2" id="user-label">Logged in as: {{ user }}</span>
      <button style="margin-right: 5px" type="button" class="btn btn-primary" @click="visitCloud">Back to Cloud</button>
      <div id="dropdown" class="dropdown">
        <button class="btn btn-primary dropdown-toggle" type="button" id="settingsDropdown" data-bs-toggle="dropdown" aria-expanded="false">
          <i class="fas fa-cog"></i>
        </button>
        <ul class="dropdown-menu" aria-labelledby="settingsDropdown">
          <li id="button-logout" class="dropdown-item" @click="logout">Logout</li>
          <li id="button-change-password" class="dropdown-item" @click="redirectToChangePassword">Change Password</li>
          <li id="button-delete-account" class="dropdown-item text-danger" @click="showDeleteConfirmation = true">Delete Account</li>
        </ul>
      </div>
    </div>
    <br>

    <HubAppManagement></HubAppManagement>

    <HubDeletionConfirmationDialog
        v-model:visible="showDeleteConfirmation"
        :on-confirm="deleteAccount"
        message="Are you sure you want to delete this app?"
    ></HubDeletionConfirmationDialog>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted } from 'vue';
import axios from "axios";
import router from "@/router";
import {session} from "@/components/hub/shared";
import HubAppManagement from "@/components/hub/HubAppManagement.vue";
import HubDeletionConfirmationDialog from "@/components/hub/HubDeletionConfirmationDialog.vue";

export default defineComponent({
  name: 'HubComponent',
  components: {HubDeletionConfirmationDialog, HubAppManagement},

  setup() {
    const user = ref<string>("");
    const showDeleteConfirmation = ref(false);

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

    const visitCloud = () => {
      router.push('/');
    };

    onMounted(() => {
      user.value = session.user
    });

    return {
      user,
      logout,
      showDeleteConfirmation,
      deleteAccount,
      confirmDeleteAccount,
      redirectToChangePassword,
      visitCloud,
    };
  },
});
</script>


TODO: Hub
* DownloadTag

TODO: Cloud
* FindApps
* DownloadTag