<template>
  <div class="forgot-password-page flex flex-col items-center justify-center min-h-screen bg-bg-base p-4">
    <div class="bg-bg-muted p-8 rounded-lg shadow-md w-full max-w-md">
      <h1 class="text-3xl font-bold text-center text-text-base mb-6">Lupa Kata Sandi Admin</h1>

      <p class="text-text-muted text-center mb-4">Masukkan email admin Anda untuk menerima tautan reset kata sandi.</p>

      <form @submit.prevent="handleForgotPassword">
        <div class="mb-4">
          <label for="email" class="block text-text-muted text-sm font-bold mb-2">Email:</label>
          <input
            type="email"
            id="email"
            v-model="email"
            class="shadow appearance-none border rounded w-full py-2 px-3 text-text-base bg-bg-base leading-tight focus:outline-none focus:shadow-outline"
            placeholder="Masukkan email Anda"
            required
          />
        </div>

        <div class="flex items-center justify-between">
          <button
            type="submit"
            class="btn btn-secondary w-full"
          >
            Kirim Tautan Reset
          </button>
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

<script>
import { ref } from 'vue';
import axios from 'axios'; // Uncomment if you have axios configured globally
import { useToast } from "vue-toastification";

export default {
  name: 'ForgotPassword',
  setup() {
    const email = ref('');
    const toast = useToast();

    const handleForgotPassword = async () => {
      console.log('Forgot password request for:', email.value);
      try {
        const response = await axios.post('/api/forgot-password', {
          email: email.value,
        });
        toast.success(response.data.message || 'Jika akun dengan email tersebut terdaftar, tautan reset kata sandi telah dikirim.');
      } catch (error) {
        console.error('Forgot password error:', error);
        toast.error(error.response?.data?.message || 'Terjadi kesalahan saat mengirim tautan reset.');
      }
    };

    return {
      email,
      handleForgotPassword,
    };
  },
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>
