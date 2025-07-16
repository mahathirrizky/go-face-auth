<template>
  <Dialog
    :visible="isOpen"
    @update:visible="$emit('update:isOpen', $event)"
    :header="title"
    modal
    :style="{ width: widthValue }"
  >
    <template #default>
      <slot></slot>
    </template>
    <template #footer v-if="$slots.footer">
      <slot name="footer"></slot>
    </template>
  </Dialog>
</template>

<script setup>
import { computed } from 'vue';
import Dialog from 'primevue/dialog';

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
    default: 'md', // sm, md, lg, xl, etc.
  },
});

defineEmits(['update:isOpen']);

const widthValue = computed(() => {
  switch (props.maxWidth) {
    case 'sm': return '24rem';
    case 'md': return '30rem';
    case 'lg': return '40rem';
    case 'xl': return '50rem';
    case '2xl': return '60rem';
    case '3xl': return '70rem';
    case '4xl': return '80rem';
    case '5xl': return '90rem';
    default: return '30rem';
  }
});
</script>