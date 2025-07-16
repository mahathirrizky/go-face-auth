<template>
  <div class="reset-password-page flex flex-col items-center justify-center min-h-screen bg-bg-base p-4">
    <div class="bg-bg-muted p-8 rounded-lg shadow-md w-full max-w-md">
      <h1 class="text-3xl font-bold text-center text-text-base mb-6">Reset Kata Sandi Admin</h1>

      <p v-if="!tokenValid" class="text-danger text-center mb-4">Tautan reset kata sandi tidak valid atau sudah kedaluwarsa.</p>
      <p v-else class="text-text-muted text-center mb-4">Masukkan kata sandi baru Anda.</p>

      <form @submit.prevent="handleResetPassword" v-if="tokenValid">
        <BaseInput
          id="newPassword"
          label="Kata Sandi Baru:"
          v-model="newPassword"
          type="password"
          placeholder="Minimal 6 karakter"
          :required="true"
          :toggleMask="true"
          :feedback="true"
        >
          <template #header>
              <h6>Atur Kata Sandi</h6>
          </template>
          <template #footer>
              <Divider />
              <p class="mt-2">Saran:</p>
              <ul class="pl-2 ml-2 mt-0" style="line-height: 1.5">
                  <li>Minimal satu huruf kecil</li>
                  <li>Minimal satu huruf besar</li>
                  <li>Minimal satu angka</li>
                  <li>Minimal 8 karakter</li>
              </ul>
          </template>
        </BaseInput>

        <BaseInput
          id="confirmPassword"
          label="Konfirmasi Kata Sandi Baru:"
          v-model="confirmPassword"
          type="password"
          placeholder="Konfirmasi kata sandi Anda"
          :required="true"
          :toggleMask="true"
          :feedback="false"
        />

        <div class="flex items-center justify-between mt-6">
          <BaseButton :fullWidth="true">
            <i class="pi pi-refresh"></i> Reset Kata Sandi
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
import BaseButton from '../ui/BaseButton.vue';
import BaseInput from '../ui/BaseInput.vue';
import Divider from 'primevue/divider';
import axios from 'axios';

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
