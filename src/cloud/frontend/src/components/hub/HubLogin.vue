<template>
  <div class="container my-5">
    <div class="row justify-content-center">
      <div class="col-lg-5 col-md-7 col-sm-9">
        <div class="hub-management-container p-4 shadow-sm bg-dark rounded">
          <h3 class="text-center mb-4">Ocelot Hub Login</h3>
          <form @submit.prevent="login" class="p-4">
            <div class="mb-3">
              <input
                  v-model="user"
                  id="input-username"
                  type="text"
                  class="form-control"
                  placeholder="Username"
                  required
              />
            </div>
            <div class="mb-3">
              <input
                  v-model="password"
                  id="input-password"
                  type="password"
                  class="form-control"
                  placeholder="Password"
                  required
              />
            </div>
            <div class="d-grid">
              <button
                  id="button-login"
                  type="submit"
                  class="btn btn-primary">
                Login
              </button>
            </div>
          </form>
          <p class="text-center mt-3">
            Or create an account
            <a
                id="registration-redirect"
                @click.prevent="redirectToRegistration"
                href="#"
                class="text-primary">
              here
            </a>.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

// TODO Hub: I need a backend endpoint for "isCookieValid" and "isOriginValid". Both executed at page load.
// TODO frontend: At loading page check if user info (cookie/origin) is found, else register/login? If so, check if cookie is okay, else show login form. if so, check if origin is okay, else change it. Show remote repo contents, like list of apps/tags.
// TODO Hub: "invalid input" in insufficient info for users
// TODO frontend: hub/, hub/login, hub/registration

<script lang="ts">
import { defineComponent, ref } from 'vue';
import { useRouter } from 'vue-router';
import {doRequest, goToHubPage} from "@/components/hub/shared";

export default defineComponent({
  name: 'HubLogin',
  setup() {
    const user = ref('');
    const password = ref('');
    const router = useRouter();

    const login = async () => {
      const loginForm = { user: user.value, password: password.value, origin: window.origin };
      await doRequest("/login", loginForm)
      goToHubPage("")
    };

    const redirectToRegistration = () => {
      goToHubPage("/registration")
    };

    return {
      user,
      password,
      login,
      redirectToRegistration,
    };
  },
});
</script>
