<template>
  <div class="reset-password-page flex flex-col items-center justify-center min-h-screen bg-bg-base p-4">
    <div class="bg-bg-muted p-8 rounded-lg shadow-md w-full max-w-md">
      <h1 class="text-3xl font-bold text-center text-text-base mb-6">Reset Kata Sandi Admin</h1>

      <p v-if="!tokenValid" class="text-danger text-center mb-4">Tautan reset kata sandi tidak valid atau sudah kedaluwarsa.</p>
      <p v-else class="text-text-muted text-center mb-4">Masukkan kata sandi baru Anda.</p>

      <form @submit.prevent="handleResetPassword" v-if="tokenValid">
        <PasswordInput
          id="newPassword"
          label="Kata Sandi Baru:"
          v-model="newPassword"
          placeholder="Minimal 6 karakter"
          required
        />

        <PasswordInput
          id="confirmPassword"
          label="Konfirmasi Kata Sandi Baru:"
          v-model="confirmPassword"
          placeholder="Konfirmasi kata sandi Anda"
          required
        />

        <div class="flex items-center justify-between mt-6">
          <BaseButton :fullWidth="true">
            Reset Kata Sandi
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
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import PasswordInput from '../ui/PasswordInput.vue';
import BaseButton from '../ui/BaseButton.vue';
// import axios from 'axios'; // Uncomment if you have axios configured globally

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
  try {
    const response = await axios.post('/api/reset-password', {
      token: token.value,
      new_password: newPassword.value,
    });
    alert(response.data.message || 'Kata sandi berhasil direset.');
    router.push({ name: 'AuthPage' }); // Redirect to login after successful reset
  } catch (error) {
    console.error('Reset password error:', error);
    alert(error.response?.data?.message || 'Terjadi kesalahan saat mereset kata sandi.');
    tokenValid.value = false; // Mark token as invalid on backend error
  }
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>
