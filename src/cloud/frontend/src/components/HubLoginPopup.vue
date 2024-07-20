<template>
    <div class="popup-content">
      <form @submit.prevent="register">
        <input v-model="username" type="text" placeholder="Username" required />
        <input v-model="password" type="password" placeholder="Password" required />
        <button @click="register()">Register</button>
      </form>
    </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import axios from "axios";

export class RegistrationForm {
  user: string;
  password: string;
  origin: string;
  email: string;

  constructor(user: string, password: string, origin: string, email: string) {
    this.user = user;
    this.password = password;
    this.origin = origin;
    this.email = email;
  }
}

export default defineComponent({
  name: 'HubLoginPopup',
  data() {
    return {
      username: '',
      password: ''
    };
  },
  methods: {
    async register() {
      const url = 'http://localhost:8082/registration';
      const payload = new RegistrationForm("admin", "admin", "http://localhost:8081", "asdf@asdf.com")
      try {
        const response = await axios.post(url, payload);
        if (response.status === 200) {
          this.$emit('authenticated');
        }
      } catch (error) {
        alert('An error occurred: ' + error);
      }
    },
  }
});
</script>