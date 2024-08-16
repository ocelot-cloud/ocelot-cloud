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
// Define color variables for the dark theme
$background-color: #282c34
$secondary-background-color: #3a3f47
$primary-color: #007bff
$primary-color-hover: lighten($primary-color, 10%)
$accent-color: #006400
$text-color: #ffffff
$text-muted: #aaaaaa
$border-color: #ffffff
$hover-background-color: #404040
$drop-zone-background-color: darken($secondary-background-color, 5%)
$drop-zone-border-color: $border-color
$drop-zone-hover-background-color: darken($secondary-background-color, 10%)
$delete-button-hover-color: darken(#dc3545, 10%)
$box-shadow-color: rgba(0, 0, 0, 0.5)

// Apply styles
html, body, input, .form-control, .modal-content, .dropdown-menu
  background-color: $background-color !important
  color: $text-color !important

#app
  font-family: Avenir, Helvetica, Arial, sans-serif
  -webkit-font-smoothing: antialiased
  -moz-osx-font-smoothing: grayscale
  text-align: center
  color: $text-color
  margin-top: 60px

p, h1, h2, h3, h4, h5, h6, label, input, span, .dropdown-item, li
  color: $text-color !important

input::placeholder
  color: $text-muted !important

input:-webkit-autofill,
input:-webkit-autofill:hover,
input:-webkit-autofill:focus,
input:-webkit-autofill:active
  -webkit-box-shadow: 0 0 0 30px $background-color inset !important
  -webkit-text-fill-color: $text-color !important

.dropdown-menu, .modal-content
  border: 1px solid $border-color !important

.dropdown-item:hover
  background-color: $hover-background-color !important

.btn-close
  background-color: $text-color !important

#drag-and-drop-area
  background-color: $drop-zone-background-color
  border: 2px dashed $drop-zone-border-color
  color: $text-color
  padding: 30px
  text-align: center
  transition: background-color 0.3s ease

  &:hover
    background-color: $drop-zone-hover-background-color

.list-group-item
  background-color: $secondary-background-color !important
  color: $text-color !important
  border: none !important
  padding: 0.5rem 1rem !important
  text-align: left !important
  transition: background-color 0.3s ease

  &:hover
    background-color: lighten($secondary-background-color, 5%) !important

  &.active
    background-color: $accent-color !important
    border-color: $accent-color !important

.tag-management-container
  background-color: $background-color !important
  border-radius: 8px
  box-shadow: 0 2px 10px $box-shadow-color
  padding: 20px
  margin: 20px

.app-info
  font-size: 1.1rem
  color: $text-color

.file-upload-area
  .drop-zone
    width: 100%
    height: 150px
    border-radius: 8px
    transition: background-color 0.3s ease

    &:hover
      background-color: $drop-zone-hover-background-color

.tag-list-section
  .list-group-item
    font-weight: normal
    transition: background-color 0.3s ease, color 0.3s ease

    &.active
      background-color: $primary-color
      color: $text-color
      font-weight: bold

.app-operations
  button
    min-width: 100px
    color: $text-color

  .btn-danger:hover
    background-color: $delete-button-hover-color !important

</style>

