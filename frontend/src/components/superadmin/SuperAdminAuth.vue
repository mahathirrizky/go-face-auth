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
        <BaseInput
          id="email-address"
          label="Email address"
          v-model="email"
          type="email"
          placeholder="Email address"
          required
          :label-sr-only="true"
        />
        <PasswordInput
          id="password"
          label="Password"
          v-model="password"
          placeholder="Password"
          required
          :label-sr-only="true"
        />

        <div class="mt-6">
          <BaseButton :fullWidth="true">
            <span class="absolute left-0 inset-y-0 flex items-center pl-3">
              <i class="fas fa-sign-in-alt text-primary group-hover:text-opacity-90"></i>
            </span>
            Sign in
          </BaseButton>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import axios from 'axios';
import { useRouter } from 'vue-router';
import { useToast } from "vue-toastification";
import { useAuthStore } from '../../stores/auth';
import BaseInput from '../ui/BaseInput.vue';
import BaseButton from '../ui/BaseButton.vue';
import PasswordInput from '../ui/PasswordInput.vue';

const email = ref('');
const password = ref('');
const router = useRouter();
const toast = useToast();
const authStore = useAuthStore();

const handleLogin = async () => {
  try {
    const response = await axios.post('/api/login/superadmin', {
      email: email.value,
      password: password.value,
    });

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
</script>

<style scoped>
/* Tailwind handles styling */
</style>
