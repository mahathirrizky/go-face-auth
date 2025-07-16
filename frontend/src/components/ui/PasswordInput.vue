<template>
  <div class="mb-4">
    <label :for="id" class="block text-text-muted text-sm font-bold mb-2">{{ label }}</label>
    <Password
      :id="id"
      :modelValue="modelValue"
      @update:modelValue="$emit('update:modelValue', $event)"
      :required="required"
      :placeholder="placeholder"
      toggleMask
      :feedback="false"
      class="w-full"
    />
  </div>
</template>

<script setup>
import Password from 'primevue/password';

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
</script>