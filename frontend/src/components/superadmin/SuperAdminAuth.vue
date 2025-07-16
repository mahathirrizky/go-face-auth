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
        <div class="mb-4">
          <label for="password" class="block text-text-muted text-sm font-bold mb-2 sr-only">Password</label>
          <Password
            id="password"
            v-model="password"
            placeholder="Password"
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
          <BaseButton :fullWidth="true">
            <span class="absolute left-0 inset-y-0 flex items-center pl-3">
              <i class="pi pi-sign-in"></i>
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
import { useToast } from 'primevue/usetoast';
import { useAuthStore } from '../../stores/auth';
import BaseInput from '../ui/BaseInput.vue';
import BaseButton from '../ui/BaseButton.vue';
import Password from 'primevue/password';
import Divider from 'primevue/divider';

const email = ref('');
const password = ref('');
const router = useRouter();
const toast = useToast();
const authStore = useAuthStore();

const handleLogin = async () => {
  try {
    const response = await axios.post('/login/superadmin', {
      email: email.value,
      password: password.value,
    });

    if (response.data && response.data.status === 'success') {
      const { token, user } = response.data.data;
      authStore.setAuth(user, token);
      toast.add({ severity: 'success', summary: 'Success', detail: 'Login successful!', life: 3000 });
      router.push('/dashboard');
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data.meta.message || 'Login failed.', life: 3000 });
    }
  } catch (error) {
    console.error('Login error:', error);
    toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data?.meta?.message || 'An error occurred during login.', life: 3000 });
  }
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>
