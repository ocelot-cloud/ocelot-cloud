<template>
  <div class="container my-5">
    <div class="row justify-content-center">
      <div class="col-lg-5 col-md-7 col-sm-9">
        <div class="entity-management-container p-4 shadow-sm bg-dark rounded">
          <h3 class="text-center mb-4">Login</h3>
          <form @submit.prevent="login" class="p-4">
            <ValidatedInput
                validationType="username"
                v-model="user"
                :submitted="submitted"
            />
            <ValidatedInput
                validationType="password"
                v-model="password"
                :submitted="submitted"
            />
            <div class="d-grid">
              <button id="button-login" type="submit" class="btn btn-primary">
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
import { defineComponent, ref } from 'vue';
import { doRequest, goToHubPage } from "@/components/hub/shared";
import ValidatedInput from "@/components/hub/ValidatedInput.vue";

export default defineComponent({
  name: 'HubLogin',
  components: { ValidatedInput },
  setup() {
    const user = ref('');
    const password = ref('');
    const submitted = ref(false);

    const login = async () => {
      submitted.value = true;
      if (user.value && password.value) {
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
      submitted,
    };
  },
});
</script>
