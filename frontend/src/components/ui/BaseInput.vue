<template>
  <div class="mb-4">
    <label :for="id" class="block text-text-muted text-sm font-bold mb-2">{{ label }}</label>
    <div class="relative">
      <input
        :type="type"
        :id="id"
        :value="modelValue"
        @input="$emit('update:modelValue', $event.target.value)"
        class="shadow appearance-none border rounded w-full py-2 px-3 text-text-base bg-bg-base leading-tight focus:outline-none focus:shadow-outline"
        :class="{ 'pr-10': hasIcon }"
        :required="required"
        :placeholder="placeholder"
      />
      <div class="absolute top-1/2 right-0 -translate-y-1/2 pr-3" v-if="hasIcon">
        <slot name="icon"></slot>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useSlots } from 'vue';

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
