import { defineConfig } from "cypress";

export default defineConfig({
  e2e: {},
});

module.exports = {
  defaultCommandTimeout: 10000, // == 10 seconds
}