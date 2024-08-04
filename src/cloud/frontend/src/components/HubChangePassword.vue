<template>
  <div class="popup-content container my-4">
    <div class="row justify-content-center">
      <div class="col-lg-6 col-md-8 col-sm-10">
        <h3>Ocelot Hub</h3>
        <form @submit.prevent="changePassword" class="p-4 border rounded shadow-sm">
          <div class="mb-3">
            <input v-model="user" id="username" type="text" class="form-control" placeholder="Username" required />
          </div>
          <div class="mb-3">
            <input v-model="oldPassword" id="old_password" type="password" class="form-control" placeholder="Old Password" required />
          </div>
          <div class="mb-3">
            <input v-model="newPassword" id="new_password" type="password" class="form-control" placeholder="New Password" required />
          </div>
          <button type="submit" class="btn btn-primary">Change Password</button>
        </form>
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
import axios from 'axios';
import { useRouter } from 'vue-router';
import router from "@/router";

export default defineComponent({
  name: 'HubChangePassword',
  setup() {
    const user = ref('');
    const oldPassword = ref('');
    const newPassword = ref('');

    const changePassword = async () => {
      const url = 'http://localhost:8082';
      try {
        const changePasswordForm = { user: user.value, old_password: oldPassword.value, new_password: newPassword.value };
        const response = await axios.post(url + '/user/password', changePasswordForm);
        if (response.status === 200) {
          alert("Password was changed.")
          await router.push('/hub');
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

    return {
      user,
      oldPassword,
      newPassword,
      changePassword,
    };
  },
});
</script>
