<template>
  <div id="app" class="container mt-5">
    <h3>Ocelot Hub</h3>
    <div class="d-flex justify-content-end align-items-center mb-3">
      <span class="me-2">Logged in as: {{ user }}</span>
      <button class="btn btn-primary" @click="logout">Logout</button>
      <button class="btn btn-primary" @click="redirectToChangePassword">Change Password</button>
      <button class="btn btn-danger" @click="showDeleteConfirmation = true">Delete Account</button>
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
            <button type="button" class="btn btn-secondary" @click="showDeleteConfirmation = false">Cancel</button>
            <button type="button" class="btn btn-danger" @click="confirmDeleteAccount">Confirm</button>
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

export default defineComponent({
  name: 'HubComponent',
  setup() {
    const user = ref<string | null>(null);
    const showDeleteConfirmation = ref(false);

    const checkAuth = async () => {
      const url = 'http://localhost:8082';
      try {
        const response = await axios.get(url + "/auth-check");
        user.value = response.data.value;
      } catch (error) {
        redirectToLogin();
      }
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

    const changePassword = async () => {
      const url = 'http://localhost:8082';
      const data = {
        user: "yourUsername",
        old_password: "yourOldPassword",
        new_password: "yourNewPassword"
      };
      await axios.post(url + "/password", data);
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

    onMounted(() => {
      checkAuth();
    });

    return {
      user,
      logout,
      showDeleteConfirmation,
      deleteAccount,
      confirmDeleteAccount,
      redirectToChangePassword,
    };
  },
});
</script>


TODO: Hub
* DeleteUser
* Write automated cypress tests
* ChangePassword
* CreateApp
* GetTags
* DeleteApp
* UploadTag
* DownloadTag
* DeleteTag

TODO: Cloud
* FindApps
* DownloadTag