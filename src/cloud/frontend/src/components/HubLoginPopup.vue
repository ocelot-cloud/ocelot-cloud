<template>
  <div class="popup-content container my-4">
    <div class="row justify-content-center">
      <div class="col-lg-6 col-md-8 col-sm-10">
        <form @submit.prevent="register" class="p-4 border rounded shadow-sm">
          <div class="mb-3">
            <input v-model="user" id="username" type="text" class="form-control" placeholder="Username" required />
          </div>

          <div class="mb-3">
            <input v-model="password" id="password" type="password" class="form-control" placeholder="Password" required />
          </div>

          <div class="mb-3">
            <input v-model="email" id="email" type="email" class="form-control" placeholder="E-Mail" required />
          </div>

          <button type="submit" class="btn btn-primary">Register</button>
        </form>
      </div>
    </div>
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

// TODO At loading page: check if user info (cookie/origin) is found, else register/login? If so, check if cookie is okay, else show login form. if so, check if origin is okay, else change it. Show remote repo contents, like list of apps/tags.

export default defineComponent({
  name: 'HubLoginPopup',
  data() {
    return {
      user: '',
      password: '',
      email: ''
    };
  },
  methods: {
    async register() {
      const url = 'http://localhost:8082/registration';
      const payload = new RegistrationForm(this.user, this.password, window.location.origin, this.email)
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