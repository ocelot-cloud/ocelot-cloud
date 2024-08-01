<template>
  <div>
    <h3>Ocelot Hub</h3>
    <button @click="redirectToLogin">Login</button>
    <button @click="redirectToRegistration">Register</button>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';

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
    checkAuth() {
      const authCookie = document.cookie.split('; ').find(row => row.startsWith('auth='));

      if (authCookie) {
        const cookieParts = authCookie.split('; ');
        let expirationDate = null;

        for (let part of cookieParts) {
          if (part.startsWith('expires=')) {
            expirationDate = new Date(part.split('=')[1]);
          }
        }

        if (expirationDate) {
          const currentTime = new Date();
          const tenSecondsFromNow = new Date(currentTime.getTime() + 10000);

          if (expirationDate > tenSecondsFromNow) {
            return;
          }
        }
      }

      // this.redirectToLogin()
    },
  }
});
</script>