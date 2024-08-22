<template>
  <div id="app">
    <router-view></router-view>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from 'vue-class-component';
import Home from './components/cloud/Home.vue';
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
@import 'bootstrap/scss/bootstrap'

$primary-color: $primary
$secondary-color: $secondary
$background-color: #282c34
$secondary-background-color: #3a3f47
$text-light: white

html, body, input, .form-control, .modal-content, .dropdown-menu
  background-color: $background-color !important
  color: $text-light !important

#app
  font-family: Avenir, Helvetica, Arial, sans-serif
  -webkit-font-smoothing: antialiased
  -moz-osx-font-smoothing: grayscale
  text-align: left
  color: $text-light

p, h1, h2, h3, h4, h5, h6, label, input, span, .dropdown-item, li, a
  color: $text-light !important

input::placeholder
  color: $secondary !important

input:-webkit-autofill,
input:-webkit-autofill:hover,
input:-webkit-autofill:focus,
input:-webkit-autofill:active
  -webkit-box-shadow: 0 0 0 30px $background-color inset !important
  -webkit-text-fill-color: $text-light !important

.dropdown-menu li
  cursor: pointer

.dropdown-menu, .modal-content
  border: 1px solid $secondary-color !important

.dropdown-item:hover
  background-color: darken($secondary-background-color, 5%) !important

.list-group-item
  border: none !important
  transition: background-color 0.3s ease

  &:hover
    background-color: lighten($secondary-background-color, 5%) !important

  &.active
    background-color: $primary-color !important
    color: $text-light !important

.entity-management-container
  padding: 20px
  margin: 20px
  background-color: $secondary-background-color
  border-radius: 8px
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.5)

.modal-header, .modal-footer, .modal-body
  border-bottom: none !important
  border-top: none !important

.invalid-feedback
  display: block

.file-upload-area
  .drop-zone
    width: 100%
    height: 150px
    border: 2px dashed $primary-color
    border-radius: 8px
    color: lighten($secondary-color, 20%)
    transition: background-color 0.3s ease

    &:hover
      background-color: lighten($secondary-background-color, 5%)
</style>
