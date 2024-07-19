<template>
  <div v-if="showPopup" class="popup">
    <div class="popup-content">
      <h2>{{ isRegistering ? 'Register' : 'Login' }}</h2>
      <form @submit.prevent="submit">
        <input v-model="username" type="text" placeholder="Username" required />
        <input v-model="password" type="password" placeholder="Password" required />
        <button type="submit">{{ isRegistering ? 'Register' : 'Login' }}</button>
      </form>
      <button @click="toggleRegister">{{ isRegistering ? 'Already have an account? Login' : 'Create a new account' }}</button>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import axios from 'axios';

export default defineComponent({
  name: 'LoginPopup',
  data() {
    return {
      showPopup: true,
      isRegistering: false,
      username: '',
      password: ''
    };
  },
  methods: {
    async submit() {
      const url = this.isRegistering ? '/registration' : '/login';
      try {
        await axios.post(`http://localhost:8082${url}`, {
          user: this.username,
          password: this.password
        });
        if (this.isRegistering) {
          // Automatically login after registration
          await axios.post('http://localhost:8082/login', {
            user: this.username,
            password: this.password
          });
        }
        this.showPopup = false;
        this.$emit('authenticated');
      } catch (error) {
        console.error('Error during authentication', error);
      }
    },
    toggleRegister() {
      this.isRegistering = !this.isRegistering;
    }
  }
});
</script>
