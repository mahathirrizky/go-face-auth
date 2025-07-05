<template>
  <div class="reset-password-page flex flex-col items-center justify-center min-h-screen bg-bg-base p-4">
    <div class="bg-bg-muted p-8 rounded-lg shadow-md w-full max-w-md">
      <h1 class="text-3xl font-bold text-center text-text-base mb-6">Reset Kata Sandi Admin</h1>

      <p v-if="!tokenValid" class="text-danger text-center mb-4">Tautan reset kata sandi tidak valid atau sudah kedaluwarsa.</p>
      <p v-else class="text-text-muted text-center mb-4">Masukkan kata sandi baru Anda.</p>

      <form @submit.prevent="handleResetPassword" v-if="tokenValid">
        <div class="mb-4 relative">
          <label for="newPassword" class="block text-text-muted text-sm font-bold mb-2">Kata Sandi Baru:</label>
          <input
            :type="passwordFieldType"
            id="newPassword"
            v-model="newPassword"
            class="shadow appearance-none border rounded w-full py-2 px-3 text-text-base bg-bg-base leading-tight focus:outline-none focus:shadow-outline pr-10"
            placeholder="Minimal 6 karakter"
            required
          />
          <button
            type="button"
            @click="toggleNewPasswordVisibility"
            class="absolute inset-y-0 right-0 pr-3 flex items-center text-sm leading-5 mt-6"
          >
            <svg v-if="passwordFieldType === 'password'" class="h-5 w-5 text-text-muted" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
            </svg>
            <svg v-else class="h-5 w-5 text-text-muted" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.542-7 .985-3.14 3.29-5.578 6.16-7.037m6.715 6.715A3 3 0 0112 15a3 3 0 01-3-3m-6.715 6.715L3 21m9-9l9 9" />
            </svg>
          </button>
        </div>

        <div class="mb-6 relative">
          <label for="confirmPassword" class="block text-text-muted text-sm font-bold mb-2">Konfirmasi Kata Sandi Baru:</label>
          <input
            :type="confirmPasswordFieldType"
            id="confirmPassword"
            v-model="confirmPassword"
            class="shadow appearance-none border rounded w-full py-2 px-3 text-text-base bg-bg-base mb-3 leading-tight focus:outline-none focus:shadow-outline pr-10"
            placeholder="Konfirmasi kata sandi Anda"
            required
          />
          
          <button
            type="button"
            @click="toggleConfirmPasswordVisibility"
            class="absolute inset-y-0 right-0 pr-3 flex items-center text-sm leading-5 mt-6"
          >
            <svg v-if="confirmPasswordFieldType === 'password'" class="h-5 w-5 text-text-muted" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
            </svg>
            <svg v-else class="h-5 w-5 text-text-muted" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.542-7 .985-3.14 3.29-5.578 6.16-7.037m6.715 6.715A3 3 0 0112 15a3 3 0 01-3-3m-6.715 6.715L3 21m9-9l9 9" />
            </svg>
          </button>
        </div>

        <div class="flex items-center justify-between">
          <button
            type="submit"
            class="btn btn-secondary w-full"
          >
            Reset Kata Sandi
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
    const passwordFieldType = ref('password');
    const confirmPasswordFieldType = ref('password');

    const toggleNewPasswordVisibility = () => {
      passwordFieldType.value = passwordFieldType.value === 'password' ? 'text' : 'password';
    };

    const toggleConfirmPasswordVisibility = () => {
      confirmPasswordFieldType.value = confirmPasswordFieldType.value === 'password' ? 'text' : 'password';
    };

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
      passwordFieldType,
      confirmPasswordFieldType,
      toggleNewPasswordVisibility,
      toggleConfirmPasswordVisibility,
      handleResetPassword,
    };
  },
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>
