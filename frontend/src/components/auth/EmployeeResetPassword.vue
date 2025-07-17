<template>
  <div class="min-h-screen flex items-center justify-center bg-bg-base py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8 p-10 bg-bg-muted rounded-lg shadow-xl">
      <div class="flex justify-center mb-6">
        <img class="h-20 w-auto" src="/vite.svg" alt="Workflow" />
      </div>
      <div>
        <h2 class="mt-2 text-center text-3xl font-extrabold text-text-base">
          Reset Kata Sandi Karyawan
        </h2>
        <p class="mt-2 text-center text-sm text-text-muted">
          Masukkan kata sandi baru Anda.
        </p>
      </div>
      <BaseForm :resolver="resolver" :initialValues="initialValues" @submit="handleResetPassword" v-slot="{ $form }" class="mt-8 space-y-6">
        <div class="flex flex-col gap-1">
          <BaseInput
            id="password"
            name="password"
            label="Kata Sandi Baru:"
            type="password"
            placeholder="Masukkan kata sandi baru Anda"
            :required="true"
            :toggleMask="true"
            :feedback="true"
            :invalid="$form.password?.invalid"
          />
          <template v-if="$form.password?.invalid">
              <Message v-for="(error, index) of $form.password.errors" :key="index" severity="error" size="small" variant="simple">{{ error.message }}</Message>
          </template>
        </div>

        <BaseInput
          id="confirm-password"
          name="confirmPassword"
          label="Konfirmasi Kata Sandi:"
          type="password"
          placeholder="Konfirmasi kata sandi baru Anda"
          :required="true"
          :toggleMask="true"
          :feedback="false"
          :invalid="$form.confirmPassword?.invalid"
          :errorMessage="$form.confirmPassword?.error?.message"
        />

        <div class="mt-6">
          <BaseButton :fullWidth="true" type="submit">
            Reset Kata Sandi
          </BaseButton>
        </div>
      </BaseForm>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import BaseButton from '../ui/BaseButton.vue';
import BaseInput from '../ui/BaseInput.vue';
import BaseForm from '../ui/BaseForm.vue'; // Import BaseForm
import { zodResolver } from '@primevue/forms/resolvers/zod';
import { z } from 'zod';
import Message from 'primevue/message';


const route = useRoute();
const router = useRouter();
const toast = useToast();

const token = ref('');

onMounted(() => {
  token.value = route.query.token || '';
  if (!token.value) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Token tidak ditemukan. Link tidak valid.', life: 3000 });
    router.push('/'); // Redirect to home or login
  }
});

const passwordSchema = z.object({
  password: z.string()
    .min(8, { message: 'Minimal 8 karakter.' })
    .refine((value) => /[a-z]/.test(value), {
      message: 'Minimal satu huruf kecil.'
    })
    .refine((value) => /[A-Z]/.test(value), {
      message: 'Minimal satu huruf besar.'
    })
    .refine((value) => /\d/.test(value), {
      message: 'Minimal satu angka.'
    }),
  confirmPassword: z.string(),
}).refine((data) => data.password === data.confirmPassword, {
  message: 'Kata sandi dan konfirmasi kata sandi tidak cocok.',
  path: ['confirmPassword'],
});

const resolver = zodResolver(passwordSchema);

const initialValues = ref({
  password: '',
  confirmPassword: '',
});

const handleResetPassword = async ({ valid, data }) => {
  if (!valid) {
    toast.add({ severity: 'error', summary: 'Validasi Gagal', detail: 'Silakan periksa kembali input Anda.', life: 3000 });
    return;
  }

  try {
    const response = await axios.post('/api/reset-password', {
      token: token.value,
      password: data.password,
      password_confirmation: data.confirmPassword,
      token_type: 'employee_password_reset' // Specify token type
    });

    if (response.data && response.data.status === 'success') {
      toast.add({ severity: 'success', summary: 'Success', detail: 'Kata sandi berhasil direset! Silakan login.', life: 3000 });
      router.push('/login/employee'); // Redirect to employee login
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data.meta.message || 'Gagal mereset kata sandi.', life: 3000 });
    }
  } catch (error) {
    console.error('Password reset error:', error);
    toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data?.message || 'Terjadi kesalahan saat mereset kata sandi.', life: 3000 });
  }
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>
