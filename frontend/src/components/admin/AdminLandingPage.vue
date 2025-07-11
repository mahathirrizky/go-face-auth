<template>
  <div class="admin-login-page flex flex-col items-center justify-center min-h-screen bg-bg-base p-4">
    <div class="bg-bg-muted p-8 rounded-lg shadow-md w-full max-w-md">
      <h1 class="text-3xl font-bold text-center text-text-base mb-6">Login Admin Perusahaan</h1>

      <form @submit.prevent="handleLogin">
        <BaseInput
          id="email"
          label="Email:"
          v-model="email"
          type="email"
          placeholder="Masukkan email Anda"
          required
          :label-sr-only="true"
        />

        <PasswordInput
          id="password"
          label="Kata Sandi:"
          v-model="password"
          placeholder="Masukkan kata sandi Anda"
          required
          :label-sr-only="true"
        />

        <div class="mt-6">
          <BaseButton type="submit" :fullWidth="true">
            <span class="absolute left-0 inset-y-0 flex items-center pl-3">
              <i class="fas fa-sign-in-alt text-primary group-hover:text-opacity-90"></i>
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
import { useToast } from 'vue-toastification';
import { useAuthStore } from '../../stores/auth';
import BaseInput from '../ui/BaseInput.vue';
import PasswordInput from '../ui/PasswordInput.vue';
import BaseButton from '../ui/BaseButton.vue';

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
      toast.success('Login successful!');
      router.push('/dashboard');
    } else {
      toast.error(response.data.message || 'Login failed.');
    }
  } catch (error) {
    console.error('Login error:', error);
    let message = 'Login failed. Please check your credentials.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.error(message);
  }
};

const goToForgotPassword = () => {
  router.push({ name: 'ForgotPassword' });
};

const goToPricingSection = () => {
  let hostname = window.location.hostname;
  let port = window.location.port ? `:${window.location.port}` : '';
  let baseUrl = hostname;

  if (hostname.endsWith('.localhost')) {
    baseUrl = 'localhost';
  } else if (hostname.split('.').length > 2) {
    const parts = hostname.split('.');
    baseUrl = parts.slice(parts.length - 2).join('.');
  }

  const newUrl = `${window.location.protocol}//${baseUrl}${port}/#pricing`;
  window.location.href = newUrl;
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>
