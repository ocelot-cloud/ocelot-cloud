<template>
  <div class="mb-3">
    <input
        :id="computedId"
        :type="inputType"
        :class="['form-control', {'is-invalid': submitted && hasError}]"
        :value="modelValue"
        @input="updateValue"
        :placeholder="placeholderText"
        required
    />
    <div v-if="submitted && hasError" class="invalid-feedback">
      {{ errorMessageText }}
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, computed } from 'vue';
type ValidationType = 'username' | 'password';
export default defineComponent({
  name: 'ValidatedInput',
  props: {
    modelValue: {
      type: String,
      required: true,
    },
    validationType: {
      type: String as () => ValidationType,
      required: true,
    },
    submitted: {
      type: Boolean,
      required: true,
    },
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {

    interface ValidationConfig {
      id: string;
      type: string;
      pattern: RegExp;
      errorMessage: string;
      placeholder: string;
    }

    const validationConfig: Record<ValidationType, ValidationConfig> = {
      username: {
        id: 'input-username',
        type: 'text',
        pattern: /^[0-9a-zA-Z]{3,20}$/,
        errorMessage: 'Invalid input, allowed symbols are [0-9a-zA-Z-] and the length must be between 3 and 20.',
        placeholder: 'Username',
      },
      password: {
        id: 'input-password',
        type: 'password',
        pattern: /^[0-9a-zA-Z]{8,20}$/,
        errorMessage: 'Password must be between 8 and 20 characters and contain only [0-9a-zA-Z-].',
        placeholder: 'Password',
      },
    };

    const config = validationConfig[props.validationType as ValidationType];

    const hasError = computed(() => !config.pattern.test(props.modelValue));
    const computedId = computed(() => config.id);
    const inputType = computed(() => config.type);
    const errorMessageText = computed(() => config.errorMessage);
    const placeholderText = computed(() => config.placeholder);

    const updateValue = (event: Event) => {
      emit('update:modelValue', (event.target as HTMLInputElement).value);
    };

    return {
      hasError,
      computedId,
      inputType,
      errorMessageText,
      placeholderText,
      updateValue,
    };
  },
});
</script>
