<template>
  <div class="container my-5">
    <div class="row justify-content-center">
      <div class="col-lg-5 col-md-7 col-sm-9">
        <div class="hub-management-container p-4 shadow-sm bg-dark rounded">
          <h3 class="text-center mb-4">Login</h3>
          <form @submit.prevent="login" class="p-4">
            <div class="mb-3">
              <input
                  v-model="user"
                  id="input-username"
                  type="text"
                  class="form-control"
                  :class="{'is-invalid': submitted && usernameError}"
                  placeholder="Username"
                  required
              />
              <div v-if="submitted && usernameError" class="invalid-feedback">
                {{ usernameErrorMessage }}
              </div>
            </div>
            <div class="mb-3">
              <input
                  v-model="password"
                  id="input-password"
                  type="password"
                  class="form-control"
                  :class="{'is-invalid': submitted && passwordError}"
                  placeholder="Password"
                  required
              />
              <div v-if="submitted && passwordError" class="invalid-feedback">
                {{ passwordErrorMessage }}
              </div>
            </div>
            <div class="d-grid">
              <button
                  id="button-login"
                  type="submit"
                  class="btn btn-primary"
              >
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
                class="text-primary"
            >
              here
            </a>.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>


<script lang="ts">
import { computed, defineComponent, ref } from 'vue';
import {
  doRequest, globalPasswordErrorMessage, globalUsernameErrorMessage,
  goToHubPage,
  passwordPattern,
  usernamePattern
} from "@/components/hub/shared";

export default defineComponent({
  name: 'HubLogin',
  setup() {
    const user = ref('');
    const password = ref('');
    const submitted = ref(false);

    const usernameError = computed(() => !usernamePattern.test(user.value));
    const passwordError = computed(() => !passwordPattern.test(password.value));
    const usernameErrorMessage = globalUsernameErrorMessage
    const passwordErrorMessage = globalPasswordErrorMessage

    const login = async () => {
      submitted.value = true;
      if (!usernameError.value && !passwordError.value) {
        const loginForm = { user: user.value, password: password.value, origin: window.origin };
        await doRequest("/login", loginForm);
        goToHubPage("");
      }
    };

    const redirectToRegistration = () => {
      goToHubPage("/registration");
    };

    return {
      user,
      password,
      login,
      redirectToRegistration,
      usernameError,
      passwordError,
      usernameErrorMessage,
      passwordErrorMessage,
      submitted,
    };
  },
});
</script>
