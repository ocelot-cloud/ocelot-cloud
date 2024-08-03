<template>
  <div>
    <h3>Ocelot Hub</h3>
    <button @click="redirectToLogin">Login</button>
    <button @click="redirectToRegistration">Register</button>
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
    redirectToRegistration() {
      this.$router.push('/hub/registration');
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
</script>