<template>
  <div class="container my-5">
    <div class="row justify-content-center">
      <div class="col-lg-5 col-md-7 col-sm-9">
        <div class="entity-management-container p-4 shadow-sm bg-dark rounded">
          <h3 class="text-center mb-4">Registration</h3>
          <form @submit.prevent="register" class="p-4">
            <ValidatedInput :submitted="submitted" validation-type="username" v-model="user"></ValidatedInput>
            <ValidatedInput :submitted="submitted" validation-type="password" v-model="password"></ValidatedInput>
            <ValidatedInput :submitted="submitted" validation-type="email" v-model="email"></ValidatedInput>
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
import {doRequest, goToHubPage} from "@/components/hub/shared";
import ValidatedInput from "@/components/hub/ValidatedInput.vue";

export default defineComponent({
  name: 'HubRegistration',
  components: {ValidatedInput},
  setup() {
    const user = ref('');
    const password = ref('');
    const email = ref('');
    const submitted = ref(false);

    const register = async () => {
      submitted.value = true
      const registrationForm = { user: user.value, password: password.value, origin: window.location.origin, email: email.value };
      const response = await doRequest("/registration", registrationForm)
      if (response) {
        goToHubPage("/login");
      }
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
      submitted,
    };
  },
});
</script>
