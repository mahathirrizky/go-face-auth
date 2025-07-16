<template>
  <Button
    :type="type"
    :class="{ 'w-full': fullWidth }"
  >
    <slot></slot>
  </Button>
</template>

<script setup>
import Button from 'primevue/button';

defineProps({
  type: {
    type: String,
    default: 'submit',
  },
  fullWidth: {
    type: Boolean,
    default: false,
  }
});
</script>