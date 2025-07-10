<template>
  <div class="admin-login-page flex flex-col items-center justify-center min-h-screen bg-bg-base p-4">
    <div class="bg-bg-muted p-8 rounded-lg shadow-md w-full max-w-md">
      <h1 class="text-3xl font-bold text-center text-text-base mb-6">Login Admin Perusahaan</h1>

      <form @submit.prevent="handleLogin">
        <div class="mb-4">
          <label for="email" class="sr-only">Email:</label>
          <input
            type="email"
            id="email"
            v-model="email"
            class="shadow-sm appearance-none border border-gray-300 rounded-md relative block w-full px-3 py-2 placeholder-gray-500 text-text-base bg-bg-base focus:outline-none focus:ring-secondary focus:border-secondary sm:text-sm"
            placeholder="Masukkan email Anda"
            required
          />
        </div>

        <div class="mb-6 relative">
          <label for="password" class="sr-only">Kata Sandi:</label>
          <input
            :type="passwordFieldType"
            id="password"
            v-model="password"
            class="shadow-sm appearance-none border border-gray-300 rounded-md relative block w-full pr-10 pl-3 py-2 placeholder-gray-500 text-text-base bg-bg-base focus:outline-none focus:ring-secondary focus:border-secondary sm:text-sm"
            placeholder="Masukkan kata sandi Anda"
            required
          />
          <span
            @click="togglePasswordVisibility"
            class="absolute inset-y-0 right-0 pr-3 flex items-center cursor-pointer"
          >
            <font-awesome-icon :icon="showPassword ? ['far', 'eye-slash'] : ['far', 'eye']" class="h-5 w-5 text-gray-400 hover:text-gray-600" />
          </span>
        </div>

        <div class="mt-6">
          <button
            type="submit"
            class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-secondary hover:bg-opacity-90 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-secondary"
          >
            <span class="absolute left-0 inset-y-0 flex items-center pl-3">
              <!-- Heroicon name: solid/lock-closed -->
              <svg class="h-5 w-5 text-primary group-hover:text-opacity-90" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                <path fill-rule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clip-rule="evenodd" />
              </svg>
            </span>
            Login
          </button>
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
    const showPassword = ref(false);
    const router = useRouter();
    const toast = useToast(); // Initialize toast
    const authStore = useAuthStore();

    const togglePasswordVisibility = () => {
      showPassword.value = !showPassword.value;
      passwordFieldType.value = showPassword.value ? 'text' : 'password';
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
      showPassword,
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