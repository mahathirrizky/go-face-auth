<template>
  <div class="admin-login-page flex flex-col items-center justify-center min-h-screen bg-bg-base p-4">
    <div class="bg-bg-muted p-8 rounded-lg shadow-md w-full max-w-md">
      <h1 class="text-3xl font-bold text-center text-text-base mb-6">Login Admin Perusahaan</h1>

      <form @submit.prevent="handleLogin">
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

        <div class="mb-6 relative">
          <label for="password" class="block text-text-muted text-sm font-bold mb-2">Kata Sandi:</label>
          <input
            :type="passwordFieldType"
            id="password"
            v-model="password"
            class="shadow appearance-none border rounded w-full py-2 px-3 text-text-base bg-bg-base mb-3 leading-tight focus:outline-none focus:shadow-outline pr-10"
            placeholder="Masukkan kata sandi Anda"
            required
          />
          <button
            type="button"
            @click="togglePasswordVisibility"
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

        <div class="flex items-center justify-between mb-4">
          <button
            type="submit"
            class="btn btn-secondary w-full"
          >
            Login
          </button>
        </div>
      </form>

      <div class="text-center mt-4">
        <a
          href="#"
          @click.prevent="goToForgotPassword"
          class="inline-block align-baseline font-bold text-sm text-accent hover:opacity-90"
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

<script>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import { useToast } from 'vue-toastification';
import { useAuthStore } from '../../stores/auth';

export default {
  name: 'AdminLandingPage',
  setup() {
    const email = ref('');
    const password = ref('');
    const passwordFieldType = ref('password');
    const router = useRouter();
    const toast = useToast(); // Initialize toast
    const authStore = useAuthStore();

    const togglePasswordVisibility = () => {
      passwordFieldType.value = passwordFieldType.value === 'password' ? 'text' : 'password';
    };

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

    return {
      email,
      password,
      passwordFieldType,
      togglePasswordVisibility,
      handleLogin,
      goToForgotPassword,
      goToPricingSection,
    };
  },
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>