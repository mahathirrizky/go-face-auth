<template>
  <div class="admin-login-page flex flex-col items-center justify-center min-h-screen bg-bg-base p-4">
    <div class="bg-bg-muted p-8 rounded-lg shadow-md w-full max-w-md">
      <h1 class="text-3xl font-bold text-center text-text-base mb-6">Login Admin Perusahaan</h1>

      <form @submit.prevent="handleLogin">
        <div class="mb-4">
          <label for="email" class="block text-text-muted text-sm font-bold mb-2 sr-only">Email:</label>
          <InputText
            id="email"
            v-model="email"
            type="email"
            placeholder="Masukkan email Anda"
            required
            class="w-full"
          />
        </div>

        <div class="mb-4">
          <label for="password" class="block text-text-muted text-sm font-bold mb-2 sr-only">Kata Sandi:</label>
          <Password
            id="password"
            v-model="password"
            placeholder="Masukkan kata sandi Anda"
            :required="true"
            toggleMask
            :feedback="true"
            class="w-full"
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
          </Password>
        </div>

        <div class="mt-6">
          <BaseButton type="submit" :fullWidth="true">
            <span class="absolute left-0 inset-y-0 flex items-center pl-3">
              <i class="pi pi-sign-in"></i>
            </span>
            Login
          </BaseButton>
        </div>
      </form>

      <div class="text-sm text-center mt-4">
        <a
          href="#"
          @click.prevent="goToForgotPassword"
          class="font-medium text-accent hover:text-accent-dark"
        >
          Lupa Kata Sandi?
        </a>
      </div>

      <div class="text-center mt-4">
        <p class="text-text-muted text-sm">
          Belum punya akun?
          <a
            href="#"
            @click.prevent="goToPricingSection"
            class="font-bold text-accent hover:opacity-90"
          >
            Daftar Sekarang
          </a>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { useAuthStore } from '../../stores/auth';
import BaseButton from '../ui/BaseButton.vue';
import Password from 'primevue/password';
import Divider from 'primevue/divider';
import InputText from 'primevue/inputtext';



const email = ref('');
const password = ref('');
const router = useRouter();
const toast = useToast();
const authStore = useAuthStore();

const handleLogin = async () => {
  try {
    const response = await axios.post('/api/login/admin-company', {
      email: email.value,
      password: password.value,
    });

    if (response.data && response.data.status === 'success') {
      const { token, user } = response.data.data;
      authStore.setAuth(user, token);
      await authStore.fetchCompanyDetails();
      toast.add({ severity: 'success', summary: 'Success', detail: 'Login successful!', life: 3000 });
      router.push('/dashboard');
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data.message || 'Login failed.', life: 3000 });
    }
  } catch (error) {
    console.error('Login error:', error);
    let message = 'Login failed. Please check your credentials.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  }
};

const goToForgotPassword = () => {
  router.push({ name: 'ForgotPassword' });
};

const goToPricingSection = () => {
  const mainFrontendUrl = process.env.VITE_MAIN_FRONTEND_URL || 'http://localhost:5173'; // Fallback for development
  const newUrl = `${mainFrontendUrl}/#pricing`;
  window.location.href = newUrl;
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>
