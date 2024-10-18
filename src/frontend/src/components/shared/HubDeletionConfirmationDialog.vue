<template>
  <div v-if="visible" class="modal fade show" style="display: block;" tabindex="-1">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header bg-dark text-light">
          <h5 class="modal-title">Deletion Confirmation</h5>
        </div>
        <div class="modal-body bg-dark text-light">
          <p>Are you sure you want to delete {{ messageSuffix }}</p>
        </div>
        <div class="modal-footer bg-dark">
          <button id="button-delete-confirmation" type="button" class="btn btn-danger" @click="onConfirm">Confirm</button>
          <button id="button-delete-cancel" type="button" class="btn btn-secondary" @click="closeModal">Cancel</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import {defineComponent, PropType} from 'vue';

export default defineComponent({
  props: {
    visible: {
      type: Boolean,
      required: true,
    },
    messageSuffix: {
      type: String,
      required: true,
    },
    onConfirm: {
      type: Function as PropType<(event: MouseEvent) => void>,
      required: true,
    },
  },
  methods: {
    closeModal() {
      this.$emit('update:visible', false);
    },
  },
});
</script>