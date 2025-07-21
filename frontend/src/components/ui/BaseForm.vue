<template>
  <Form class="space-y-5" v-slot="$form" :initialValues="initialValues" :resolver="resolver" @submit="(...args) => $emit('submit', ...args)">
    <slot :$form="$form"></slot>
  </Form>
</template>

<script setup>
import { Form } from '@primevue/forms';

defineProps({
  initialValues: {
    type: Object,
    required: true,
  },
  resolver: {
    type: Function,
    required: true,
  },
});

defineEmits(['submit']);
</script>