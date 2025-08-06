<template>
  <div class="admin-login-page flex flex-col items-center justify-center min-h-screen bg-bg-base p-4">
    <Card class="w-full max-w-md shadow-xl">
        <template #title>
            <h1 class="text-3xl font-bold text-center text-text-base">Login Admin Perusahaan</h1>
        </template>
        <template #content>
            <form class="mt-4" @submit.prevent="handleLogin">
                <div class="p-fluid">
                    <div class="field mb-4">
                      <FloatLabel variant="on">
                        
                        <InputText id="email" v-model="email" type="email" required fluid/>
                        
                        <label for="email">Email :</label>
                      </FloatLabel>
                    </div>
                    <div class="field mb-4">
                      <FloatLabel variant="on">

                        <Password id="password" v-model="password" :toggleMask="true" :feedback="false" required fluid/>
                        <label for="password">Password :</label>
                      </FloatLabel>
                    </div>
                    <Button type="submit" :loading="loading" :label="loading ? 'Logging in...' : 'Login'" class="w-full" />
                </div>
            </form>

            <div class="text-sm text-center mt-4">
                <a href="#" @click.prevent="goToForgotPassword" class="font-medium text-accent hover:text-accent-dark">
                Lupa Kata Sandi?
                </a>
            </div>

            <div class="text-center mt-4">
                <p class="text-text-muted text-sm">
                Belum punya akun?
                <a href="#" @click.prevent="goToPricingSection" class="font-bold text-accent hover:opacity-90">
                    Daftar Sekarang
                </a>
                </p>
            </div>
        </template>
    </Card>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { useAuthStore } from '../../stores/auth';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import Password from 'primevue/password';
import Card from 'primevue/card';
import FloatLabel  from 'primevue/floatlabel';

const email = ref('');
const password = ref('');
const router = useRouter();
const toast = useToast();
const authStore = useAuthStore();
const loading = ref(false);

const handleLogin = async () => {
  loading.value = true;
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
    let message = error.response?.data?.message || 'Login failed. Please check your credentials.';
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  } finally {
    loading.value = false;
  }
};

const goToForgotPassword = () => {
  router.push({ name: 'ForgotPassword' });
};

const goToPricingSection = () => {
  const mainFrontendUrl = import.meta.env.VITE_MAIN_FRONTEND_URL || 'http://localhost:5173';
  const newUrl = `${mainFrontendUrl}/#pricing`;
  window.location.href = newUrl;
};
</script>

<style scoped>
.field > label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 600;
}
</style>