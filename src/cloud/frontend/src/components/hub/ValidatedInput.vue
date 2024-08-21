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
import {
  defaultAllowedSymbols,
  generateInvalidInputMessage, getDefaultValidationRegex, maxLengthPassword,
  defaultMaxLength, minLengthPassword,
  defaultMinLength
} from "@/components/hub/shared";

type ValidationType = 'username' | 'password' | 'email' | 'app';

const validationConfig = {
  username: {
    id: 'input-username',
    type: 'text',
    pattern: getDefaultValidationRegex(),
    errorMessage: generateInvalidInputMessage('username', defaultAllowedSymbols, defaultMinLength, defaultMaxLength),
    placeholder: 'Username',
  },
  password: {
    id: 'input-password',
    type: 'password',
    pattern: new RegExp(`^${defaultAllowedSymbols}{${minLengthPassword},${maxLengthPassword}}$`),
    errorMessage: generateInvalidInputMessage('password', defaultAllowedSymbols, minLengthPassword, maxLengthPassword),
    placeholder: 'Password',
  },
  email: {
    id: 'input-email',
    type: 'email',
    pattern: new RegExp(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+.[a-zA-Z]{2,}$`),
    errorMessage: "Invalid email, must have this format: x@x.x where x is a placeholder",
    placeholder: 'Email',
  },
  app: {
    id: 'input-app',
    type: 'text',
    pattern: getDefaultValidationRegex(),
    errorMessage: generateInvalidInputMessage('app', defaultAllowedSymbols, minLengthPassword, maxLengthPassword),
    placeholder: 'New app to create',
  },
};

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
