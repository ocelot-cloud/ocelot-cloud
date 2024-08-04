<template>
  <div id="app">
    <router-view></router-view>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from 'vue-class-component';
import Home from './components/Home.vue';
import axios from "axios";

@Options({
  components: {
    Home,
  },
})
export default class App extends Vue {}

axios.defaults.withCredentials = true;
</script>

<style lang="sass">
html body, input, .form-control, .modal-content, .dropdown-menu
  background-color: #282c34 !important

.modal-content, .dropdown-menu
  border: 1px solid white !important

#app
  font-family: Avenir, Helvetica, Arial, sans-serif
  -webkit-font-smoothing: antialiased
  -moz-osx-font-smoothing: grayscale
  text-align: center
  color: #2c3e50
  margin-top: 60px

p, h1, h2, h3, h4, h5, h6, label, input, span, .dropdown-item
  color: white !important

.dropdown-item:hover
  background-color: #404040 !important

.btn-close
  background-color: white !important

input::placeholder
  color: gray !important

input:-webkit-autofill,
input:-webkit-autofill:hover,
input:-webkit-autofill:focus,
input:-webkit-autofill:active
  -webkit-box-shadow: 0 0 0 30px #282c34 inset !important
  -webkit-text-fill-color: white !important
</style>

