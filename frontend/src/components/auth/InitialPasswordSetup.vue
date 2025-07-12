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
      <form class="mt-8 space-y-6" @submit.prevent="handlePasswordSetup">
        <PasswordInput
          id="password"
          label="Kata Sandi Baru"
          v-model="password"
          placeholder="Masukkan kata sandi baru Anda"
          required
        />
        <PasswordInput
          id="confirm-password"
          label="Konfirmasi Kata Sandi"
          v-model="confirmPassword"
          placeholder="Konfirmasi kata sandi baru Anda"
          required
        />

        <div class="mt-6">
          <BaseButton :fullWidth="true">
            Atur Kata Sandi
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
import { useToast } from 'vue-toastification';
import PasswordInput from '../ui/PasswordInput.vue';
import BaseButton from '../ui/BaseButton.vue';

const route = useRoute();
const router = useRouter();
const toast = useToast();

const password = ref('');
const confirmPassword = ref('');
const token = ref('');

onMounted(() => {
  token.value = route.query.token || '';
  if (!token.value) {
    toast.error('Token tidak ditemukan. Link tidak valid.');
    router.push('/'); // Redirect to home or login
  }
});

const handlePasswordSetup = async () => {
  if (password.value !== confirmPassword.value) {
    toast.error('Kata sandi dan konfirmasi kata sandi tidak cocok.');
    return;
  }

  try {
    const response = await axios.post('/api/initial-password-setup', {
      token: token.value,
      password: password.value,
    });

    if (response.data && response.data.status === 'success') {
      toast.success('Kata sandi berhasil diatur! Silakan login.');
      router.push('/login/employee'); // Redirect to employee login
    } else {
      toast.error(response.data.meta.message || 'Gagal mengatur kata sandi.');
    }
  } catch (error) {
    console.error('Password setup error:', error);
    toast.error(error.response?.data?.meta?.message || 'Terjadi kesalahan saat mengatur kata sandi.');
  }
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>
