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

type ValidationType = 'username' | 'password' | 'email';

const allowedSymbols = '[0-9a-zA-Z]';
const minLengthUsername = 3;
const maxLengthUsername = 20;
const minLengthPassword = 8;
const maxLengthPassword = 30;

function generateInvalidInputMessage(fieldName: string, allowedSymbols: string, minLength: number, maxLength: number): string {
  return `Invalid ${fieldName}, allowed symbols are ${allowedSymbols} and the length must be between ${minLength} and ${maxLength}.`;
}

const validationConfig = {
  username: {
    id: 'input-username',
    type: 'text',
    pattern: new RegExp(`^${allowedSymbols}{${minLengthUsername},${maxLengthUsername}}$`),
    errorMessage: generateInvalidInputMessage('username', allowedSymbols, minLengthUsername, maxLengthUsername),
    placeholder: 'Username',
  },
  password: {
    id: 'input-password',
    type: 'password',
    pattern: new RegExp(`^${allowedSymbols}{${minLengthPassword},${maxLengthPassword}}$`),
    errorMessage: generateInvalidInputMessage('password', allowedSymbols, minLengthPassword, maxLengthPassword),
    placeholder: 'Password',
  },
  email: {
    id: 'input-email',
    type: 'email',
    pattern: new RegExp(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+.[a-zA-Z]{2,}$`),
    errorMessage: generateInvalidInputMessage('email', allowedSymbols, minLengthPassword, maxLengthPassword),
    placeholder: 'Email',
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
