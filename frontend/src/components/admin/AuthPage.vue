<template>
  <div class="min-h-screen flex items-center justify-center bg-bg-base py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8 p-10 bg-bg-muted rounded-lg shadow-md">
      <div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-text-base">
          Masuk ke Akun Admin Anda
        </h2>
      </div>
      <form class="mt-8 space-y-6" @submit.prevent="handleLogin">
        <div class="rounded-md shadow-sm -space-y-px">
          <div>
            <label for="email-address" class="sr-only">Alamat Email</label>
            <input id="email-address" name="email" type="email" autocomplete="email" required v-model="email"
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-text-base bg-bg-base rounded-t-md focus:outline-none focus:ring-secondary focus:border-secondary focus:z-10 sm:text-sm"
              placeholder="Alamat Email">
          </div>
          <div class="relative">
            <label for="password" class="sr-only">Kata Sandi</label>
            <input id="password" name="password" :type="passwordFieldType" autocomplete="current-password" required v-model="password"
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-text-base bg-bg-base rounded-b-md focus:outline-none focus:ring-secondary focus:border-secondary focus:z-10 sm:text-sm pr-10"
              placeholder="Kata Sandi">
            <button
              type="button"
              @click="togglePasswordVisibility"
              class="absolute inset-y-0 right-0 pr-3 flex items-center text-sm leading-5"
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
        </div>

        <div class="flex items-center justify-between">
          <div class="flex items-center">
            <input id="remember-me" name="remember-me" type="checkbox"
              class="h-4 w-4 text-secondary focus:ring-secondary border-gray-300 rounded">
            <label for="remember-me" class="ml-2 block text-sm text-text-muted">
              Ingat saya
            </label>
          </div>

          <div class="text-sm">
            <a href="#" class="font-medium text-accent hover:opacity-90">
              Lupa kata sandi?
            </a>
          </div>
        </div>

        <div>
          <button type="submit"
            class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-secondary hover:bg-opacity-90 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-secondary">
            <span class="absolute left-0 inset-y-0 flex items-center pl-3">
              <!-- Heroicon name: solid/lock-closed -->
              <svg class="h-5 w-5 text-primary group-hover:text-opacity-90" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                <path fill-rule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clip-rule="evenodd" />
              </svg>
            </span>
            Masuk
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue';
import axios from 'axios';
import { useAuthStore } from '../../stores/auth';
import { useRouter } from 'vue-router';
import { useToast } from "vue-toastification";

export default {
  name: 'AuthPage',
  setup() {
    const email = ref('');
    const password = ref('');
    const passwordFieldType = ref('password');
    const authStore = useAuthStore();
    const router = useRouter();
    const toast = useToast();

    const togglePasswordVisibility = () => {
      passwordFieldType.value = passwordFieldType.value === 'password' ? 'text' : 'password';
    };

    const handleLogin = async () => {
      try {
        const response = await axios.post('/api/admin/login', {
          email: email.value,
          password: password.value,
        });

        console.log('Login API Response:', response);
        if (response.data && response.data.status === 'success') {
          const { token, user } = response.data.data;
          authStore.setAuth(user, token);
          toast.success('Login successful!');
          router.push('/dashboard');
        } else {
          toast.error(response.data.meta.message || 'Login failed.');
        }
      } catch (error) {
        console.error('Login error:', error);
        toast.error(error.response?.data?.meta?.message || 'An error occurred during login.');
      }
    };

    return {
      email,
      password,
      passwordFieldType,
      togglePasswordVisibility,
      handleLogin,
    };
  },
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>