<template>
  <div class="mb-4">
    <label :for="id" class="block text-text-muted text-sm font-bold mb-2">{{ label }}</label>
    <span :class="{ 'p-input-icon-right': hasIcon }">
      <InputText
        :type="type"
        :id="id"
        :modelValue="modelValue"
        @update:modelValue="$emit('update:modelValue', $event)"
        class="w-full"
        :required="required"
        :placeholder="placeholder"
      />
      <slot name="icon" v-if="hasIcon"></slot>
    </span>
  </div>
</template>

<script setup>
import { useSlots } from 'vue';
import InputText from 'primevue/inputtext';

const slots = useSlots();
const hasIcon = !!slots.icon;

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
    type: [String, Number],
    default: '',
  },
  type: {
    type: String,
    default: 'text',
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