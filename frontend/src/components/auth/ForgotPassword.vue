<template>
  <div class="forgot-password-page flex flex-col items-center justify-center min-h-screen bg-bg-base p-4">
    <div class="bg-bg-muted p-8 rounded-lg shadow-md w-full max-w-md">
      <h1 class="text-3xl font-bold text-center text-text-base mb-6">Lupa Kata Sandi Admin</h1>

      <p class="text-text-muted text-center mb-4">Masukkan email admin Anda untuk menerima tautan reset kata sandi.</p>

      <form @submit.prevent="handleForgotPassword">
        <BaseInput
          id="email"
          label="Email:"
          v-model="email"
          type="email"
          placeholder="Masukkan email Anda"
          required
        />

        <div class="flex items-center justify-between">
          <BaseButton :fullWidth="true">
            <i class="pi pi-envelope"></i> Kirim Tautan Reset
          </BaseButton>
        </div>
      </form>

      <div class="text-center mt-4">
        <router-link to="/admin" class="inline-block align-baseline font-bold text-sm text-accent hover:opacity-90">
          Kembali ke Login
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import BaseInput from '../ui/BaseInput.vue';
import BaseButton from '../ui/BaseButton.vue';

const email = ref('');
const toast = useToast();

const handleForgotPassword = async () => {
  try {
    const response = await axios.post('/api/forgot-password', {
      email: email.value,
    });
    toast.add({ severity: 'success', summary: 'Berhasil', detail: response.data.message || 'Jika akun dengan email tersebut terdaftar, tautan reset kata sandi telah dikirim.', life: 5000 });
  } catch (error) {
    console.error('Forgot password error:', error);
    toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data?.message || 'Terjadi kesalahan saat mengirim tautan reset.', life: 3000 });
  }
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>
