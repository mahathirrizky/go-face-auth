<template>
  <div class="min-h-screen flex items-center justify-center bg-bg-base py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8 p-10 bg-bg-muted rounded-lg shadow-xl">
      <div class="flex justify-center mb-6">
        <img class="h-20 w-auto" src="/vite.svg" alt="Workflow" />
      </div>
      <div>
        <h2 class="mt-2 text-center text-3xl font-extrabold text-text-base">
          Atur Kata Sandi Awal
        </h2>
        <p class="mt-2 text-center text-sm text-text-muted">
          Silakan atur kata sandi Anda untuk pertama kali.
        </p>
      </div>
      <BaseForm :resolver="resolver" :initialValues="initialValues" @submit="handlePasswordSetup" v-slot="{ $form }" class="mt-8 space-y-6">
        <BaseInput
          id="password"
          name="password"
          label="Kata Sandi Baru:"
          type="password"
          :required="true"
          :toggleMask="true"
          :feedback="false"
          :invalid="$form.password?.invalid"
          :fluid="true"
          :errors="$form.password?.errors"
        />

        <BaseInput
          id="confirm-password"
          name="confirmPassword"
          label="Konfirmasi Kata Sandi:"
          type="password"
          :required="true"
          :toggleMask="true"
          :feedback="false"
          :invalid="$form.confirmPassword?.invalid"
          :errors="$form.confirmPassword?.errors"
        />

        <div class="mt-6">
          <BaseButton :fullWidth="true" type="submit">
            Atur Kata Sandi
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

const handlePasswordSetup = async (event) => {
  const { valid, values: data } = event;

  if (!valid) {
    toast.add({ severity: 'error', summary: 'Validasi Gagal', detail: 'Silakan periksa kembali input Anda.', life: 3000 });
    return;
  }

  try {
    const response = await axios.post('/api/initial-password-setup', {
      token: token.value,
      password: data.password,
    });

    if (response.data && response.data.status === 'success') {
      toast.add({ severity: 'success', summary: 'Success', detail: 'Kata sandi berhasil diatur!', life: 3000 });
      router.push('/initial-password-success'); // Redirect to success page
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data.meta.message || 'Gagal mengatur kata sandi.', life: 3000 });
    }
  } catch (error) {
    console.error('Password setup error:', error);
    toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data?.meta?.message || 'Terjadi kesalahan saat mengatur kata sandi.', life: 3000 });
  }
};
</script>