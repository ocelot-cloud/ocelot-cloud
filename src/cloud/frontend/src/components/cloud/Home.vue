<template>
  <div class="container-fluid table-container">
    <div class="mt-5">
      <h3>Ocelot Cloud</h3>
      <button type="button" class="btn btn-primary" @click="visitHub">Visit Ocelot Hub</button>
      <br>
      <br>
      <table class="table table-dark table-hover" id="stack-table">
        <thead>
        <tr>
          <th scope="col">Name</th>
          <th scope="col" style="width: 250px;">State</th>
          <th scope="col">Link</th>
          <th scope="col">Actions</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="stack in stacks" :key="stack.name">
          <td>{{ stack.name }}</td>
          <td :class="getBootstrapBackgroundClass(stack.state)">
            <div class="d-flex align-items-center justify-content-center">
              <span class="me-2">{{ stack.state }}</span>
                <span v-if="stack.state === 'Starting' || stack.state === 'Downloading' || stack.state === 'Stopping'">
                  <span class="spinner-border" role="status" style="width: 1rem; height: 1rem;"></span>
                </span>
            </div>
          </td>
          <td>
            <button class="btn btn-primary" :id="'open-button-' + stack.name" :data-stack-url="getUrlFromStack(stack)" @click="openNewTab(stack)" :disabled="stack.state !== 'Available'">Open</button>
          </td>
          <td>
            <button @click="start(stack.name)" class="btn btn-success start-button" :disabled="stack.state !== 'Uninitialized'">Start</button>
            <button @click="stop(stack.name)" class="btn btn-danger stop-button" :disabled="stack.state !== 'Available'">Stop</button>
          </td>
        </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>


<script lang="ts">

import {defineComponent} from 'vue';
import {backendClient, baseDomain, scheme, stackUrl, waitTimeInMillis} from "@/components/cloud/Config";
import {Stack} from "@/components/cloud/Shared";

function getUrlFromStack(stack: Stack) {
  return `${scheme}://${stack.name}.${baseDomain}${stack.urlPath}`;
}

export default defineComponent({
  name: 'home-component',
  data() {
    return {
      stacks: [] as Stack[]
    };
  },
  created() {
    this.fetchData()
    setInterval(this.fetchData, waitTimeInMillis);
  },
  methods: {
    async fetchData() {
      // TODO When I load "home" and then cloud, the fetching continue. Maybe add a condition, do this "only on the cloud home page"
      try {
        const response = await backendClient.getResponsePromise(stackUrl);
        if (!response.ok) {
          throw new Error('Network response was not ok');
        } else if (response.status === 302 || response.status === 301) {
          console.log("Redirecting to keycloak: " + response.headers.get("X-Redirect-URL"))
          window.location.href = response.headers.get("X-Redirect-URL") || response.url;
        } else {
          this.stacks = await response.json();
          this.stacks.sort((a, b) => a.name.localeCompare(b.name));
        }
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    },
    start(name: string) {
      console.log('Starting:', name);
      backendClient.postRequest(name, stackUrl, "deploy")
    },
    stop(name: string) {
      console.log('Deleting:', name);
      backendClient.postRequest(name, stackUrl, "stop")
    },
    openNewTab(stack: Stack) {
      window.open(getUrlFromStack(stack), '_blank');
    },
    getUrlFromStack(stack: Stack) {
      return getUrlFromStack(stack)
    },
    getBootstrapBackgroundClass(state: string) {
      switch(state) {
        case 'Available': return 'bg-success text-white state-column';
        case 'Starting': return 'bg-warning text-dark state-column';
        case 'Downloading': return 'bg-warning text-dark state-column';
        case 'Stopping': return 'bg-warning text-dark state-column';
        case 'Uninitialized': return 'bg-dark text-white state-column';
        default: return '';
      }
    },
    visitHub() {
      this.$router.push('/hub');
    },
  }
});
</script>

<style scoped lang="sass">
.table-container
  @media (min-width: 576px)
    max-width: 75%
    margin: auto
</style>
