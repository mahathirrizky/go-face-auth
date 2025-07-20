<template>
  <FloatLabel>
    <template v-if="type === 'password'">
      <Password
        :id="id"
        :modelValue="modelValue"
        @update:modelValue="$emit('update:modelValue', $event)"
        @focus="handleFocus"
        @blur="handleBlur"
        :class="{ 'p-input-icon-right': hasIcon }"
        class="w-full"
        :required="required"
        :toggleMask="toggleMask"
        :feedback="feedback"
        :invalid="invalid"
        :name="name"
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
        @focus="handleFocus"
        @blur="handleBlur"
        :class="{ 'p-input-icon-right': hasIcon }"
        class="w-full"
        :required="required"
        :invalid="invalid"
        :name="name"
      />
    </template>
    <slot name="icon" v-if="hasIcon"></slot>
    <label :for="id" v-if="label">{{ label }}</label>
  </FloatLabel>
  <template v-if="invalid && errors.length && isFocused">
    <Message v-for="(error, index) of errors" :key="index" severity="error" size="small">{{ error.message }}</Message>
  </template>
</template>

<script setup>
import { ref, useSlots } from 'vue';
import InputText from 'primevue/inputtext';
import Password from 'primevue/password';
import FloatLabel from 'primevue/floatlabel';
import Message from 'primevue/message';

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
  errors: {
    type: Array,
    default: () => [],
  },
  name: {
    type: String,
    default: '',
  }
});

const isFocused = ref(false);

const handleFocus = () => {
  isFocused.value = true;
};

const handleBlur = () => {
  isFocused.value = false;
};

defineEmits(['update:modelValue']);
</script>
