<template>
  <div class="mb-4">
    <FloatLabel variant="on">
      <span :class="{ 'p-input-icon-right': hasIcon }" class="w-full block">
        <template v-if="type === 'password'">
          <Password
            :id="id"
            :modelValue="modelValue"
            @update:modelValue="$emit('update:modelValue', $event)"
            class="w-full"
            inputClass="w-full"
            :required="required"
            :placeholder="placeholder"
            :toggleMask="toggleMask"
            :feedback="feedback"
            :invalid="invalid"
          >
            <template #header v-if="$slots.header">
              <slot name="header"></slot>
            </template>
            <template #footer v-if="$slots.footer">
              <slot name="footer"></slot>
            </template>
          </Password>
        </template>
        <template v-else>
          <InputText
            :type="type"
            :id="id"
            :modelValue="modelValue"
            @update:modelValue="$emit('update:modelValue', $event)"
            class="w-full"
            :required="required"
            :placeholder="placeholder"
            :invalid="invalid"
          />
        </template>
        <slot name="icon" v-if="hasIcon"></slot>
      </span>
      <label :for="id" v-if="label">{{ label }}</label>
    </FloatLabel>
    <small v-if="invalid" class="p-error">{{ errorMessage }}</small>
  </div>
</template>

<script setup>
import { useSlots, computed } from 'vue';
import InputText from 'primevue/inputtext';
import Password from 'primevue/password';
import FloatLabel from 'primevue/floatlabel';

const slots = useSlots();
const hasIcon = !!slots.icon;

const props = defineProps({
  id: {
    type: String,
    required: true,
  },
  label: {
    type: String,
    default: '',
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
  // Props khusus untuk komponen Password
  toggleMask: {
    type: Boolean,
    default: false,
  },
  feedback: {
    type: Boolean,
    default: false,
  },
  invalid: {
    type: Boolean,
    default: false,
  },
  errorMessage: {
    type: String,
    default: '',
  },
});

defineEmits(['update:modelValue']);
</script>


<style scoped>
/* Tailwind handles styling */
</style>
