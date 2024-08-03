<template>
  <div id="app" class="app-container">
    <h3>Ocelot Hub</h3>
    <div class="header">
      <span>Logged in as: {{ user }}</span>
      <button @click="logout">Logout</button>
    </div>
    <router-view />
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

    const checkAuth = async () => {
      const url = 'http://localhost:8082';
      try {
        const response = await axios.get(url + "/auth-check");
        user.value = response.data.value;
      } catch (error) {
        alert(error)
        // TODO redirectToLogin()
      }
    };

    const logout = () => {
      // TODO: Implement actual logout logic (e.g., API call to logout)
      user.value = "";
      redirectToLogin();
    };
    const redirectToLogin = () => {
      router.push('/hub/login');
    };

    onMounted(() => {
      checkAuth();
    });

    return {
      user,
      logout,
    };
  },
});

// TODO
/*
Make username visible, Logout
FindApps
DownloadApp
ChangePassword
DeleteUser
CreateApp
DeleteApp
UploadTag
DeleteTag
GetTags
 */
</script>

<style lang="sass">
.app-container
  width: 66.66%
  margin: 0 auto
  padding: 20px
  background-color: #fff
  border-radius: 8px

.header
  display: flex
  justify-content: flex-end
  margin-bottom: 20px

.header span
  margin-top: 3px
  margin-right: 10px
</style>
