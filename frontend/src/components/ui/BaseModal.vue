<template>
  <div v-if="isOpen" class="fixed inset-0 bg-black bg-opacity-20 flex items-center justify-center z-50">
    <div class="bg-bg-muted p-8 rounded-lg shadow-lg w-full" :class="maxWidthClass">
      <div class="flex justify-between items-center mb-6">
        <h3 class="text-2xl font-bold text-text-base">{{ title }}</h3>
        <button @click="$emit('close')" class="text-text-muted hover:text-text-base">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
        </button>
      </div>
      <div class="modal-content">
        <slot></slot>
      </div>
      <div v-if="$slots.footer" class="modal-footer mt-6 flex justify-end space-x-4">
        <slot name="footer"></slot>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  isOpen: {
    type: Boolean,
    required: true,
  },
  title: {
    type: String,
    default: '',
  },
  maxWidth: {
    type: String,
    default: 'md', // sm, md, lg, xl, 2xl, etc.
  },
});

const emit = defineEmits(['close']);

const maxWidthClass = computed(() => {
  switch (props.maxWidth) {
    case 'sm': return 'max-w-sm';
    case 'md': return 'max-w-md';
    case 'lg': return 'max-w-lg';
    case 'xl': return 'max-w-xl';
    case '2xl': return 'max-w-2xl';
    case '3xl': return 'max-w-3xl';
    case '4xl': return 'max-w-4xl';
    case '5xl': return 'max-w-5xl';
    case '6xl': return 'max-w-6xl';
    case '7xl': return 'max-w-7xl';
    default: return 'max-w-md';
  }
});
</script>

<style scoped>
/* No specific styles needed, Tailwind handles it */
</style>
