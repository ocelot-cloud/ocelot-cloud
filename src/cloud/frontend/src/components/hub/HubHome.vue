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
        messageSuffix="your account?"
    ></HubDeletionConfirmationDialog>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted } from 'vue';
import router from "@/router";
import {doRequest, goToHubPage, session} from "@/components/hub/shared";
import HubAppManagement from "@/components/hub/HubAppManagement.vue";
import HubDeletionConfirmationDialog from "@/components/hub/HubDeletionConfirmationDialog.vue";

export default defineComponent({
  name: 'HubComponent',
  components: {HubDeletionConfirmationDialog, HubAppManagement},

  setup() {
    const user = ref<string>("");
    const showDeleteConfirmation = ref(false);

    const logout = async () => {
      await doRequest("/logout", null)
      user.value = "";
      redirectToLogin();
    };

    const deleteAccount = async () => {
      await doRequest("/user", null)
      user.value = "";
      redirectToLogin();
    };

    const confirmDeleteAccount = async () => {
      showDeleteConfirmation.value = false;
      await deleteAccount();
    };

    const redirectToLogin = () => {
      goToHubPage("/login")
    };

    const redirectToChangePassword = () => {
      goToHubPage("/change-password")
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