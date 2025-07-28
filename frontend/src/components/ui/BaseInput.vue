<template>
  <FloatLabel variant="on" :class="$attrs.class">
    <template v-if="type === 'password'">
      <Password
        :id="id"
        :modelValue="modelValue"
        @update:modelValue="$emit('update:modelValue', $event)"
        @focus="handleFocus"
        @blur="handleBlur"
        @keydown.enter="$emit('keydown.enter', $event)"
        :class="{ 'p-input-icon-right': hasIcon }"
        class="w-full"
        :required="required"
        :toggleMask="toggleMask"
        :feedback="feedback"
        :invalid="invalid"
        :name="name"
        fluid
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
        @keydown.enter="$emit('keydown.enter', $event)"
        :class="{ 'p-input-icon-right': hasIcon }"
        class="w-full"
        :required="required"
        :invalid="invalid"
        :name="name"
        autofocus
        v-bind="$attrs"
      />
    </template>
    <slot name="icon" v-if="hasIcon"></slot>
    <label :for="id" v-if="label">{{ label }}</label>
  </FloatLabel>
  <template v-if="invalid && errors.length">
    <div v-for="(error, index) of errors" :key="index" class="text-red-500 text-xs mt-1">{{ typeof error === 'object' && error !== null && 'message' in error ? error.message : error }}</div>
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
  },
  fluid: {
    type: Boolean,
    default: false,
  }
});



const isFocused = ref(false);

const handleFocus = () => {
  isFocused.value = true;
};

const emit = defineEmits(['update:modelValue', 'blur', 'keydown.enter']);

const handleBlur = (event) => {
  isFocused.value = false;
  emit('blur', event);
};
</script>