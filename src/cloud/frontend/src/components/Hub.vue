<template>
  <div id="app" class="container mt-5">
    <h3>Ocelot Hub</h3>
    <div class="d-flex justify-content-end align-items-center mb-3">
      <span class="me-2">Logged in as: {{ user }}</span>
      <button class="btn btn-primary" @click="logout">Logout</button>
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
        alert(error);
        redirectToLogin();
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
</script>
