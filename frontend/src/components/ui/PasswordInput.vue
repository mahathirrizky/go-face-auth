<template>
  <BaseInput
    :id="id"
    :label="label"
    :type="passwordFieldType"
    :modelValue="modelValue"
    @update:modelValue="$emit('update:modelValue', $event)"
    :required="required"
    :placeholder="placeholder"
  >
    <template #icon>
      <span class="cursor-pointer" @click="togglePasswordVisibility">
        <font-awesome-icon :icon="showPassword ? ['far', 'eye-slash'] : ['far', 'eye']" class="h-5 w-5 text-gray-400 hover:text-gray-600" />
      </span>
    </template>
  </BaseInput>
</template>

<script setup>
import { ref } from 'vue';
import BaseInput from './BaseInput.vue';

defineProps({
  id: {
    type: String,
    required: true,
  },
  label: {
    type: String,
    required: true,
  },
  modelValue: {
    type: String,
    default: '',
  },
  required: {
    type: Boolean,
    default: false,
  },
  placeholder: {
    type: String,
    default: '',
  },
});

defineEmits(['update:modelValue']);

const passwordFieldType = ref('password');
const showPassword = ref(false);

const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value;
  passwordFieldType.value = showPassword.value ? 'text' : 'password';
};
</script>
