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
      <form class="mt-8 space-y-6" @submit.prevent="handleResetPassword">
        <BaseInput
          id="password"
          label="Kata Sandi Baru:"
          v-model="password"
          type="password"
          placeholder="Masukkan kata sandi baru Anda"
          :required="true"
          :toggleMask="true"
          :feedback="true"
          :passwordFeedback="true"
        >
        </BaseInput>

        <BaseInput
          id="confirm-password"
          label="Konfirmasi Kata Sandi:"
          v-model="confirmPassword"
          type="password"
          placeholder="Konfirmasi kata sandi baru Anda"
          :required="true"
          :toggleMask="true"
          :feedback="false"
        />

        <div class="mt-6">
          <BaseButton :fullWidth="true">
            Reset Kata Sandi
          </BaseButton>
        </div>
      </form>
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

const route = useRoute();
const router = useRouter();
const toast = useToast();

const password = ref('');
const confirmPassword = ref('');
const token = ref('');

onMounted(() => {
  token.value = route.query.token || '';
  if (!token.value) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Token tidak ditemukan. Link tidak valid.', life: 3000 });
    router.push('/'); // Redirect to home or login
  }
});

const handleResetPassword = async () => {
  if (password.value !== confirmPassword.value) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Kata sandi dan konfirmasi kata sandi tidak cocok.', life: 3000 });
    return;
  }

  try {
    const response = await axios.post('/api/reset-password', {
      token: token.value,
      password: password.value,
      password_confirmation: confirmPassword.value,
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
    toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data?.meta?.message || 'Terjadi kesalahan saat mereset kata sandi.', life: 3000 });
  }
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>
