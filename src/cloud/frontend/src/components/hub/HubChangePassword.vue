<template>
  <div class="container my-5">
    <div class="row justify-content-center">
      <div class="col-lg-5 col-md-7 col-sm-9">
        <div class="entity-management-container p-4 shadow-sm bg-dark rounded">
          <h3 class="text-center mb-4">Change Password</h3>
          <form @submit.prevent="changePassword" class="p-4">
            <div class="mb-3">
              <input
                  v-model="oldPassword"
                  id="old_password"
                  type="password"
                  class="form-control"
                  placeholder="Old Password"
                  required
              />
            </div>
            <div class="mb-3">
              <input
                  v-model="newPassword"
                  id="new_password"
                  type="password"
                  class="form-control"
                  placeholder="New Password"
                  required
              />
            </div>
            <div class="d-grid">
              <button type="submit" class="btn btn-primary">Change Password</button>
            </div>
          </form>
          <p class="text-center mt-3">
            Back to
            <a @click.prevent="redirectToHubHomePage" href="#" class="text-primary">
              Hub home page
            </a>.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue';
import {doRequest, goToHubPage} from "@/components/hub/shared";

export default defineComponent({
  name: 'HubChangePassword',
  setup() {
    const user = ref('');
    const oldPassword = ref('');
    const newPassword = ref('');

    const changePassword = async () => {
      const changePasswordForm = { old_password: oldPassword.value, new_password: newPassword.value };
      await doRequest("/user/password", changePasswordForm)
      goToHubPage("")
    };

    const redirectToHubHomePage = () => {
      goToHubPage("")
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
