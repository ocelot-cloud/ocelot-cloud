<template>
  <div>
    <h3>Ocelot Hub</h3>
    <button @click="redirectToLogin">Login</button>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import axios from "axios";

export default defineComponent({
  name: 'HubComponent',
  components: {

  },
  mounted() {
    this.checkAuth();
  },
  methods: {
    redirectToLogin() {
      this.$router.push('/hub/login');
    },
    async checkAuth() {
      const url = 'http://localhost:8082';
      try {
        await axios.get(url + "/auth-check");
      } catch (error) {
        this.redirectToLogin()
      }
    },
  }
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