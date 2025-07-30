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
            <InputText
              id="email-address"
              v-model="email"
              type="email"
              autocomplete="email"
              required
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
              placeholder="Email address"
              fluid
            />
          </div>
          <div>
            <label for="password" class="sr-only">Password</label>
            <Password
              id="password"
              v-model="password"
              autocomplete="current-password"
              required
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
              placeholder="Password"
              :toggleMask="true"
              :feedback="false"
              fluid
            />
          </div>
        </div>

        <div class="mt-6">
          <Button type="submit" class="w-full">
            <span class="absolute left-0 inset-y-0 flex items-center pl-3">
              <i class="pi pi-sign-in"></i>
            </span>
            Sign in
          </Button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import axios from 'axios';
import { useRouter } from 'vue-router';
import { useToast } from 'primevue/usetoast';
import { useAuthStore } from '../../stores/auth';
import InputText from 'primevue/inputtext';
import Password from 'primevue/password';
import Button from 'primevue/button';

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
      toast.add({ severity: 'success', summary: 'Success', detail: 'Login successful!', life: 3000 });
      router.push('/dashboard');
    } else {
      const message = response.data?.message || 'Login failed.';
      toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
    }
  } catch (error) {
    console.error('Login error:', error);
    const message = error.response?.data?.message || 'An error occurred during login.';
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  }
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>