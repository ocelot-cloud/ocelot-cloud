<template>
  <div>
    <h3>Ocelot Hub</h3>
    <button>Login</button>
    <button>Register</button>
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

      alert("Cookie expired. Please reload page? TODO")
    },
  }
});
</script>