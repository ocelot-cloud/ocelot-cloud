<template>
  <div>
    <p>this is the hub page</p>
    <HubLoginPopup v-if="!isAuthenticated" @authenticated="showWelcomeMessage" />
    <div v-else>
      <p>Welcome in the hub!</p>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import HubLoginPopup from "@/components/HubLoginPopup.vue";

export default defineComponent({
  name: 'HubComponent',
  components: {
    HubLoginPopup
  },
  data() {
    return {
      isAuthenticated: false
    };
  },
  mounted() {
    this.checkAuth();
  },
  methods: {
    checkAuth() {
      const authCookie = document.cookie.split('; ').find(row => row.startsWith('auth='));
      if (authCookie && authCookie.split('=')[1] === 'true') {
        this.isAuthenticated = true;
      }
    },
    showWelcomeMessage() {
      this.isAuthenticated = true;
    }
  }
});
</script>