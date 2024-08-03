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
        <br>
        <p>Back to <a @click.prevent="redirectToLogin" href="#">login</a>.</p>
      </div>
    </div>
  </div>
</template>


<script lang="ts">
import { defineComponent, ref } from 'vue';
import axios from 'axios';
import { useRouter } from 'vue-router';

export default defineComponent({
  name: 'HubRegistration',
  setup() {
    const user = ref('');
    const password = ref('');
    const email = ref('');
    const router = useRouter();

    const register = async () => {
      const url = 'http://localhost:8082';
      const registrationForm = { user: user.value, password: password.value, origin: window.location.origin, email: email.value };

      try {
        const response = await axios.post(url + '/registration', registrationForm);
        if (response.status === 200) {
          alert('Registration successful. Please click "OK" and login.');
          await router.push('/hub/login');
        } else {
          alert(response.data);
        }
      } catch (error) {
        if (axios.isAxiosError(error) && error.response) {
          const errorMessage = error.response.data || 'An unknown error occurred';
          showErrorPopup(`An error occurred: ${errorMessage}`);
        } else {
          showErrorPopup('An unknown error occurred');
        }
      }
    };

    const showErrorPopup = (message: string) => {
      alert(message);
    };

    const redirectToLogin = () => {
      router.push('/hub/login');
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
