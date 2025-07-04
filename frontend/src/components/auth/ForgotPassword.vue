<template>
  <div class="forgot-password-page flex flex-col items-center justify-center min-h-screen bg-gray-100 p-4">
    <div class="bg-white p-8 rounded-lg shadow-md w-full max-w-md">
      <h1 class="text-3xl font-bold text-center text-gray-800 mb-6">Lupa Kata Sandi Admin</h1>

      <p class="text-gray-600 text-center mb-4">Masukkan email admin Anda untuk menerima tautan reset kata sandi.</p>

      <form @submit.prevent="handleForgotPassword">
        <div class="mb-4">
          <label for="email" class="block text-gray-700 text-sm font-bold mb-2">Email:</label>
          <input
            type="email"
            id="email"
            v-model="email"
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            placeholder="Masukkan email Anda"
            required
          />
        </div>

        <div class="flex items-center justify-between">
          <button
            type="submit"
            class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline w-full"
          >
            Kirim Tautan Reset
          </button>
        </div>
      </form>

      <div class="text-center mt-4">
        <router-link to="/admin" class="inline-block align-baseline font-bold text-sm text-blue-500 hover:text-blue-800">
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
