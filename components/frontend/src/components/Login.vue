<template>
  <div class="container mt-5">
    <div class="row justify-content-center">
      <div class="col-md-3">
        <h1 class="text-center mb-4">Login</h1>
        <form @submit.prevent="login" class="card p-4">
          <div class="mb-3">
            <label for="username-field" class="form-label">Username</label>
            <input type="text" class="form-control" id="username-field" v-model="username" placeholder="Enter username">
          </div>
          <div class="mb-3">
            <label for="password-field" class="form-label">Password</label>
            <input type="password" class="form-control" id="password-field" v-model="password" placeholder="Enter password">
          </div>
          <div class="d-grid">
            <button type="submit" class="btn btn-primary" id="login-button">Login</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue';
import { useRouter } from 'vue-router';

export default defineComponent({
  name: 'login-component',
  setup() {
    const username = ref('');
    const password = ref('');
    const router = useRouter();

    const login = async () => {
      try {
        const response = await fetch('/api/login', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            username: username.value,
            password: password.value,
          }),
        });

        if (response.ok) {
          router.push('/');
        } else {
          alert('Login failed!');
        }
      } catch (error) {
        console.error('Login error:', error);
      }
    };

    return {
      username,
      password,
      login,
    };
  },
});
</script>

