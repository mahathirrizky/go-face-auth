<template>
  <div class="min-h-screen flex items-center justify-center bg-bg-base py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8 p-10 bg-bg-muted rounded-lg shadow-xl">
      <div class="flex justify-center mb-6">
        <img class="h-20 w-auto" src="/vite.svg" alt="Workflow" />
      </div>
      <div>
        <h2 class="mt-2 text-center text-3xl font-extrabold text-text-base">
          Login SuperAdmin
        </h2>
      </div>
      <form class="mt-8 space-y-6" @submit.prevent="handleLogin">
        <div class="rounded-md shadow-sm -space-y-px">
          <div>
            <label for="email-address" class="sr-only">Email address</label>
            <input id="email-address" name="email" type="email" autocomplete="email" required v-model="email"
              class="shadow-sm appearance-none border border-gray-300 rounded-md relative block w-full px-3 py-2 placeholder-gray-500 text-text-base bg-bg-base focus:outline-none focus:ring-secondary focus:border-secondary sm:text-sm"
              placeholder="Email address">
          </div>
          <div class="mt-4">
            <label for="password" class="sr-only">Password</label>
            <div class="relative">
              <input id="password" name="password" :type="passwordFieldType" autocomplete="current-password" required v-model="password"
                class="shadow-sm appearance-none border border-gray-300 rounded-md relative block w-full pr-10 pl-3 py-2 placeholder-gray-500 text-text-base bg-bg-base focus:outline-none focus:ring-secondary focus:border-secondary sm:text-sm"
                placeholder="Password">
              <span class="absolute inset-y-0 right-0 pr-3 flex items-center cursor-pointer" @click="togglePasswordVisibility">
                <font-awesome-icon :icon="showPassword ? ['far', 'eye-slash'] : ['far', 'eye']" class="text-gray-400 hover:text-gray-600" />
              </span>
            </div>
          </div>
        </div>

        <div class="mt-6"> <!-- Added mt-6 here -->
          <button type="submit"
            class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-secondary hover:bg-opacity-90 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-secondary">
            <span class="absolute left-0 inset-y-0 flex items-center pl-3">
              <!-- Heroicon name: solid/lock-closed -->
              <svg class="h-5 w-5 text-primary group-hover:text-opacity-90" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                <path fill-rule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clip-rule="evenodd" />
              </svg>
            </span>
            Sign in
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue';
import axios from 'axios';
import { useRouter } from 'vue-router';
import { useToast } from "vue-toastification";
import { useAuthStore } from '../../stores/auth';

export default {
  name: 'SuperAdminAuth',
  setup() {
    const email = ref('');
    const password = ref('');
    const router = useRouter();
    const toast = useToast();
    const authStore = useAuthStore();

    const passwordFieldType = ref('password');
    const showPassword = ref(false);

    const togglePasswordVisibility = () => {
      showPassword.value = !showPassword.value;
      passwordFieldType.value = showPassword.value ? 'text' : 'password';
    };

    const handleLogin = async () => {
      try {
        const response = await axios.post('/api/login/superadmin', {
          email: email.value,
          password: password.value,
        });

        if (response.data && response.data.status === 'success') {
          const { token, user } = response.data.data;
          authStore.setAuth(user, token);
          // For superadmin, we might not need company details, or fetch different details
          // If superadmin also has a companyId, you might fetch it here.
          // await authStore.fetchCompanyDetails(); // Uncomment if superadmin needs company details
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
      showPassword,
      togglePasswordVisibility,
      handleLogin,
    };
  },
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>
