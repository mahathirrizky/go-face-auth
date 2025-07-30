<template>
  <div class="forgot-password-page flex flex-col items-center justify-center min-h-screen bg-bg-base p-4">
    <Card class="w-full max-w-md shadow-xl">
      <template #title>
        <h1 class="text-3xl font-bold text-center text-text-base">Lupa Kata Sandi Admin</h1>
      </template>
      <template #content>
        <p class="text-text-muted text-center mb-4">Masukkan email admin Anda untuk menerima tautan reset kata sandi.</p>

        <form @submit.prevent="handleForgotPassword" class="p-fluid mt-4">
          <div class="field mb-4">
            <label for="email">Email</label>
            <InputText id="email" v-model="email" type="email" placeholder="Masukkan email Anda" required fluid />
          </div>

          <Button type="submit" label="Kirim Tautan Reset" icon="pi pi-envelope" class="w-full" />
        </form>

        <div class="text-center mt-4">
          <router-link to="/admin" class="inline-block align-baseline font-bold text-sm text-accent hover:opacity-90">
            Kembali ke Login
          </router-link>
        </div>
      </template>
    </Card>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import InputText from 'primevue/inputtext';
import Button from 'primevue/button';
import Card from 'primevue/card';

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
.field > label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 600;
}
</style>