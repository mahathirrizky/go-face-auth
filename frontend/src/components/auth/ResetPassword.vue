<template>
  <div class="reset-password-page flex flex-col items-center justify-center min-h-screen bg-gray-100 p-4">
    <div class="bg-white p-8 rounded-lg shadow-md w-full max-w-md">
      <h1 class="text-3xl font-bold text-center text-gray-800 mb-6">Reset Kata Sandi Admin</h1>

      <p v-if="!tokenValid" class="text-red-500 text-center mb-4">Tautan reset kata sandi tidak valid atau sudah kedaluwarsa.</p>
      <p v-else class="text-gray-600 text-center mb-4">Masukkan kata sandi baru Anda.</p>

      <form @submit.prevent="handleResetPassword" v-if="tokenValid">
        <div class="mb-4">
          <label for="newPassword" class="block text-gray-700 text-sm font-bold mb-2">Kata Sandi Baru:</label>
          <input
            type="password"
            id="newPassword"
            v-model="newPassword"
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            placeholder="Minimal 6 karakter"
            required
          />
        </div>

        <div class="mb-6">
          <label for="confirmPassword" class="block text-gray-700 text-sm font-bold mb-2">Konfirmasi Kata Sandi Baru:</label>
          <input
            type="password"
            id="confirmPassword"
            v-model="confirmPassword"
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:shadow-outline"
            placeholder="Konfirmasi kata sandi Anda"
            required
          />
        </div>

        <div class="flex items-center justify-between">
          <button
            type="submit"
            class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline w-full"
          >
            Reset Kata Sandi
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
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
// import axios from 'axios'; // Uncomment if you have axios configured globally

export default {
  name: 'ResetPassword',
  setup() {
    const route = useRoute();
    const router = useRouter();
    const newPassword = ref('');
    const confirmPassword = ref('');
    const token = ref('');
    const tokenValid = ref(true); // Assume valid until checked

    onMounted(() => {
      token.value = route.query.token || '';

      if (!token.value) {
        tokenValid.value = false;
      }
      // In a real app, you might want to make an API call here to validate the token immediately
      // For now, we'll rely on the backend validation during the reset attempt.
    });

    const handleResetPassword = async () => {
      if (newPassword.value !== confirmPassword.value) {
        alert('Kata sandi baru dan konfirmasi kata sandi tidak cocok.');
        return;
      }

      if (newPassword.value.length < 6) {
        alert('Kata sandi minimal harus 6 karakter.');
        return;
      }

      console.log('Reset password attempt for:', token.value, newPassword.value);
      // Implement API call here
      // try {
      //   const response = await axios.post('/api/reset-password', {
      //     token: token.value,
      //     new_password: newPassword.value,
      //   });
      //   alert(response.data.message || 'Kata sandi berhasil direset.');
      //   router.push({ name: 'AuthPage' }); // Redirect to login after successful reset
      // } catch (error) {
      //   console.error('Reset password error:', error);
      //   alert(error.response?.data?.message || 'Terjadi kesalahan saat mereset kata sandi.');
      //   tokenValid.value = false; // Mark token as invalid on backend error
      // }
    };

    return {
      newPassword,
      confirmPassword,
      tokenValid,
      handleResetPassword,
    };
  },
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>
