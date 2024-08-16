<template>
  <div class="popup-content container my-4">
    <div class="row justify-content-center">
      <div class="col-lg-6 col-md-8 col-sm-10">
        <h3>Ocelot Hub</h3>
        <form @submit.prevent="changePassword" class="p-4 border rounded shadow-sm">
          <div class="mb-3">
            <input v-model="oldPassword" id="old_password" type="password" class="form-control" placeholder="Old Password" required />
          </div>
          <div class="mb-3">
            <input v-model="newPassword" id="new_password" type="password" class="form-control" placeholder="New Password" required />
          </div>
          <button type="submit" class="btn btn-primary">Change Password</button>
        </form>
        <br>
        <p>Back to <a @click.prevent="redirectToHubHomePage" href="#">Hub home page</a>.</p>
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
import router from "@/router";
import {doRequest} from "@/components/hub/shared";

export default defineComponent({
  name: 'HubChangePassword',
  setup() {
    const user = ref('');
    const oldPassword = ref('');
    const newPassword = ref('');

    const changePassword = async () => {
      const changePasswordForm = { old_password: oldPassword.value, new_password: newPassword.value };
      await doRequest("/user/password", changePasswordForm)
      await router.push('/hub');
    };

    const redirectToHubHomePage = () => {
      router.push("/hub")
    }

    return {
      user,
      oldPassword,
      newPassword,
      changePassword,
      redirectToHubHomePage,
    };
  },
});
</script>
