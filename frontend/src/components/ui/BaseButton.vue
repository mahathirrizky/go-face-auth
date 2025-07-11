<template>
  <button
    :type="type"
    class="btn btn-secondary"
    :class="{ 'w-full': fullWidth }"
  >
    <slot></slot>
  </button>
</template>

<script setup>
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
