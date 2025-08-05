<template>
  <div class="min-h-screen flex items-center justify-center bg-bg-base py-12 px-4 sm:px-6 lg:px-8">
    <Card class="w-full max-w-md shadow-xl">
      <template #title>
        <h2 class="text-2xl font-bold text-center text-text-base">
          Atur Kata Sandi Awal
        </h2>
        <p class="mt-2 text-center text-sm text-text-muted">
          Silakan atur kata sandi Anda untuk pertama kali.
        </p>
      </template>
      <template #content>
        <Form :resolver="resolver" :initialValues="initialValues" @submit="handlePasswordSetup" v-slot="{ errors, handleSubmit }">
          <form @submit="handleSubmit" class="p-fluid mt-4">
            <div class="field mb-4">
              <label for="password">Kata Sandi Baru</label>
              <Password id="password" name="password" v-model="initialValues.password" :feedback="true" :toggleMask="true" :invalid="!!errors.password" fluid>
                <template #footer>
                  <p class="mt-2">Saran Kata Sandi:</p>
                  <ul class="pl-4 ml-2 mt-0" style="line-height: 1.5">
                    <li :class="{ 'text-green-500': isLengthValid }">Minimal 8 karakter <i v-if="isLengthValid" class="pi pi-check"></i></li>
                    <li :class="{ 'text-green-500': hasLowercase }">Minimal satu huruf kecil (a-z) <i v-if="hasLowercase" class="pi pi-check"></i></li>
                    <li :class="{ 'text-green-500': hasUppercase }">Minimal satu huruf besar (A-Z) <i v-if="hasUppercase" class="pi pi-check"></i></li>
                    <li :class="{ 'text-green-500': hasNumber }">Minimal satu angka (0-9) <i v-if="hasNumber" class="pi pi-check"></i></li>
                  </ul>
                </template>
              </Password>
              <small class="p-error" v-if="errors.password">{{ errors.password }}</small>
            </div>

            <div class="field mb-4">
              <label for="confirm-password">Konfirmasi Kata Sandi</label>
              <Password id="confirm-password" name="confirmPassword" v-model="initialValues.confirmPassword" :feedback="false" :toggleMask="true" :invalid="!!errors.confirmPassword" fluid />
              <small class="p-error" v-if="errors.confirmPassword">{{ errors.confirmPassword }}</small>
            </div>

            <Button type="submit" label="Atur Kata Sandi" class="w-full mt-4" />
          </form>
        </Form>
      </template>
    </Card>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import Button from 'primevue/button';
import Password from 'primevue/password';
import Card from 'primevue/card';
import { Form } from '@primevue/forms';
import { zodResolver } from '@primevue/forms/resolvers/zod';
import { z } from 'zod';

const route = useRoute();
const router = useRouter();
const toast = useToast();

const token = ref('');

onMounted(() => {
  token.value = route.query.token || '';
  if (!token.value) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Token tidak ditemukan. Link tidak valid.', life: 3000 });
    router.push('/');
  }
});

const passwordSchema = z.object({
  password: z.string()
    .min(8, { message: 'Minimal 8 karakter.' })
    .refine((value) => /[a-z]/.test(value), { message: 'Minimal satu huruf kecil.' })
    .refine((value) => /[A-Z]/.test(value), { message: 'Minimal satu huruf besar.' })
    .refine((value) => /\d/.test(value), { message: 'Minimal satu angka.' }),
  confirmPassword: z.string(),
}).refine((data) => data.password === data.confirmPassword, {
  message: 'Kata sandi baru dan konfirmasi kata sandi tidak cocok.',
  path: ['confirmPassword'],
});

const resolver = zodResolver(passwordSchema);

const initialValues = ref({
  password: '',
  confirmPassword: '',
});

const isLengthValid = computed(() => initialValues.value.password.length >= 8);
const hasLowercase = computed(() => /[a-z]/.test(initialValues.value.password));
const hasUppercase = computed(() => /[A-Z]/.test(initialValues.value.password));
const hasNumber = computed(() => /\d/.test(initialValues.value.password));

const handlePasswordSetup = async ({ values: data }) => {

  try {
    const response = await axios.post('/api/initial-password-setup', {
      token: token.value,
      password: data.password,
    });

    if (response.data && response.data.status === 'success') {
      toast.add({ severity: 'success', summary: 'Success', detail: 'Kata sandi berhasil diatur!', life: 3000 });
      router.push('/initial-password-success');
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data.meta.message || 'Gagal mengatur kata sandi.', life: 3000 });
    }
  } catch (error) {
    console.error('Password setup error:', error);
    toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data?.meta?.message || 'Terjadi kesalahan saat mengatur kata sandi.', life: 3000 });
  }
};
</script>

<style scoped>
.field > label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 600;
}
</style>
