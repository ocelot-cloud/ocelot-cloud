<template>
  <div class="popup-content container my-4">
    <div class="row justify-content-center">
      <div class="col-lg-6 col-md-8 col-sm-10">
        <form @submit.prevent="login" class="p-4 border rounded shadow-sm">
          <div class="mb-3">
            <input v-model="user" id="username" type="text" class="form-control" placeholder="Username" required />
          </div>
          <div class="mb-3">
            <input v-model="password" id="password" type="password" class="form-control" placeholder="Password" required />
          </div>
          <button type="submit" class="btn btn-primary">Login</button>
        </form>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import axios from "axios";

export class LoginForm {
  constructor(
      public user: string,
      public password: string,
      public origin: string
  ) {}
}

// TODO Hub: I need a backend endpoint for "isCookieValid" and "isOriginValid". Both executed at page load.
// TODO frontend: At loading page check if user info (cookie/origin) is found, else register/login? If so, check if cookie is okay, else show login form. if so, check if origin is okay, else change it. Show remote repo contents, like list of apps/tags.
// TODO Hub: "invalid input" in insufficient info for users
// TODO frontend: hub/, hub/login, hub/registration

export default defineComponent({
  name: 'HubLogin',
  data() {
    return {
      user: '',
      password: '',
      email: ''
    };
  },
  methods: {
    async login() {
      const url = 'http://localhost:8082';
      try {
        const loginForm = new LoginForm(this.user, this.password, window.origin);
        const response = await axios.post(url + "/login", loginForm);
        console.log("Got a response")
        if (response.status === 200) {
          console.log("Status is okay")
          this.$router.push('/hub');
        } else {
          alert(response.data)
        }
      } catch (error) {
        if (axios.isAxiosError(error) && error.response) {
          const errorMessage = error.response.data || 'An unknown error occurred';
          this.showErrorPopup(`An error occurred: ${errorMessage}`);
        } else {
          this.showErrorPopup('An unknown error occurred');
        }
      }
    },
    showErrorPopup(message: string) {
      alert(message);
    }
  }
});
</script>