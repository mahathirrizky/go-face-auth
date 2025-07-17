<template>
  <div class="reset-password-page flex flex-col items-center justify-center min-h-screen bg-bg-base p-4">
    <div class="bg-bg-muted p-8 rounded-lg shadow-md w-full max-w-md">
      <h1 class="text-3xl font-bold text-center text-text-base mb-6">Reset Kata Sandi Admin</h1>

      <p v-if="!tokenValid" class="text-danger text-center mb-4">Tautan reset kata sandi tidak valid atau sudah kedaluwarsa.</p>
      <p v-else class="text-text-muted text-center mb-4">Masukkan kata sandi baru Anda.</p>

      <BaseForm :resolver="resolver" :initialValues="initialValues" @submit="handleResetPassword" v-if="tokenValid">
        <BaseInput
          id="newPassword"
          name="newPassword"
          label="Kata Sandi Baru:"
          type="password"
          :required="true"
          :toggleMask="true"
          :feedback="true"
          :invalid="$form.newPassword?.invalid"
          :errorMessage="$form.newPassword?.error?.message"
        />

        <BaseInput
          id="confirmPassword"
          name="confirmPassword"
          label="Konfirmasi Kata Sandi Baru:"
          type="password"
          :required="true"
          :toggleMask="true"
          :feedback="false"
          :invalid="$form.confirmPassword?.invalid"
          :errorMessage="$form.confirmPassword?.error?.message"
        />

        <div class="flex items-center justify-between mt-6">
          <BaseButton :fullWidth="true" type="submit">
            <i class="pi pi-refresh"></i> Reset Kata Sandi
          </BaseButton>
        </div>
      </BaseForm>

      <div class="text-center mt-4">
        <router-link to="/admin" class="inline-block align-baseline font-bold text-sm text-accent hover:opacity-90">
          Kembali ke Login
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import BaseButton from '../ui/BaseButton.vue';
import BaseInput from '../ui/BaseInput.vue';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import BaseForm from '../ui/BaseForm.vue'; // Import BaseForm
import { zodResolver } from '@primevue/forms/resolvers/zod';
import { z } from 'zod';

const route = useRoute();
const router = useRouter();
const toast = useToast();
const token = ref('');
const tokenValid = ref(true); // Assume valid until checked

onMounted(() => {
  token.value = route.query.token || '';

  if (!token.value) {
    tokenValid.value = false;
    router.replace({ name: 'TokenInvalid' }); // Redirect to new page
    return;
  }
  // In a real app, you might want to make an API call here to validate the token immediately
  // For now, we'll rely on the backend validation during the reset attempt.
});

const passwordSchema = z.object({
  newPassword: z.string()
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
}).refine((data) => data.newPassword === data.confirmPassword, {
  message: 'Kata sandi baru dan konfirmasi kata sandi tidak cocok.',
  path: ['confirmPassword'],
});

const resolver = zodResolver(passwordSchema);

const initialValues = ref({
  newPassword: '',
  confirmPassword: '',
});

const handleResetPassword = async (event) => {
  const { valid, data } = event;

  if (!valid) {
    toast.add({ severity: 'error', summary: 'Validasi Gagal', detail: 'Silakan periksa kembali input Anda.', life: 3000 });
    return;
  }

  console.log('Reset password attempt for:', token.value, data.newPassword);
  // Implement API call here
  try {
    const response = await axios.post('/api/reset-password', {
      token: token.value,
      new_password: data.newPassword,
    });
    toast.add({ severity: 'success', summary: 'Success', detail: response.data.message || 'Kata sandi berhasil direset.', life: 3000 });
    router.push({ name: 'AuthPage' }); // Redirect to login after successful reset
  } catch (error) {
    console.error('Reset password error:', error);
    if (error.response && error.response.status === 400) { // Assuming 400 for invalid/expired token
      router.replace({ name: 'TokenInvalid' }); // Redirect to new page
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data?.message || 'Terjadi kesalahan saat mereset kata sandi.', life: 3000 });
    }
    tokenValid.value = false; // Mark token as invalid on backend error
  }
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>
