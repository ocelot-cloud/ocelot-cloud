<template>
  <div class="container my-5">
    <div class="row justify-content-center">
      <div class="col-lg-5 col-md-7 col-sm-9">
        <div class="hub-management-container p-4 shadow-sm bg-dark rounded">
          <h3 class="text-center mb-4">Registration</h3>
          <form @submit.prevent="register" class="p-4">
            <div class="mb-3">
              <input v-model="user" id="input-username" type="text" class="form-control" placeholder="Username" required />
            </div>
            <div class="mb-3">
              <input v-model="password" id="input-password" type="password" class="form-control" placeholder="Password" required />
            </div>
            <div class="mb-3">
              <input v-model="email" id="input-email" type="email" class="form-control" placeholder="E-Mail" required />
            </div>
            <div class="d-grid">
              <button id="button-register" type="submit" class="btn btn-primary">Register</button>
            </div>
          </form>
          <p class="text-center mt-3">
            Back to <a @click.prevent="redirectToLogin" href="#" class="text-primary">login</a>.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue';
import axios from 'axios';
import { useRouter } from 'vue-router';
import {doRequest, goToHubPage} from "@/components/hub/shared";

export default defineComponent({
  name: 'HubRegistration',
  setup() {
    const user = ref('');
    const password = ref('');
    const email = ref('');

    const register = async () => {
      const registrationForm = { user: user.value, password: password.value, origin: window.location.origin, email: email.value };
      await doRequest("/registration", registrationForm)
      goToHubPage("/login")
    };

    const redirectToLogin = () => {
      goToHubPage("/login")
    };

    return {
      user,
      password,
      email,
      register,
      redirectToLogin,
    };
  },
});
</script>
