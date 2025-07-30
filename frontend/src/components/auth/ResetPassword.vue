<template>
  <div class="reset-password-page flex flex-col items-center justify-center min-h-screen bg-bg-base p-4">
    <Card class="w-full max-w-md shadow-xl">
      <template #title>
        <h1 class="text-3xl font-bold text-center text-text-base">Reset Kata Sandi Admin</h1>
      </template>
      <template #content>
        <p v-if="!tokenValid" class="text-danger text-center mb-4">Tautan reset kata sandi tidak valid atau sudah kedaluwarsa.</p>
        <p v-else class="text-text-muted text-center mb-4">Masukkan kata sandi baru Anda.</p>

        <Form :resolver="resolver" :initialValues="initialValues" @submit="handleResetPassword" v-if="tokenValid" v-slot="{ errors, handleSubmit }">
          <form @submit="handleSubmit" class="p-fluid mt-4">
            <div class="field mb-4">
              <label for="newPassword">Kata Sandi Baru</label>
              <Password id="newPassword" name="newPassword" v-model="initialValues.newPassword" :feedback="true" :toggleMask="true" :invalid="!!errors.newPassword" fluid />
              <small class="p-error" v-if="errors.newPassword">{{ errors.newPassword }}</small>
            </div>

            <div class="field mb-4">
              <label for="confirmPassword">Konfirmasi Kata Sandi Baru</label>
              <Password id="confirmPassword" name="confirmPassword" v-model="initialValues.confirmPassword" :feedback="false" :toggleMask="true" :invalid="!!errors.confirmPassword" fluid />
              <small class="p-error" v-if="errors.confirmPassword">{{ errors.confirmPassword }}</small>
            </div>

            <Button type="submit" label="Reset Kata Sandi" icon="pi pi-refresh" class="w-full mt-4" />
          </form>
        </Form>

        <div class="text-center mt-4">
          <router-link to="/admin" class="inline-block align-baseline font-bold text-sm text-accent hover:opacity-90">
            Kembali ke Login
          </router-link>
        </div>
      </template>
    </Card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
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
const tokenValid = ref(true);

onMounted(() => {
  token.value = route.query.token || '';

  if (!token.value) {
    tokenValid.value = false;
    router.replace({ name: 'TokenInvalid' });
    return;
  }
});

const passwordSchema = z.object({
  newPassword: z.string()
    .min(8, { message: 'Minimal 8 karakter.' })
    .refine((value) => /[a-z]/.test(value), { message: 'Minimal satu huruf kecil.' })
    .refine((value) => /[A-Z]/.test(value), { message: 'Minimal satu huruf besar.' })
    .refine((value) => /\d/.test(value), { message: 'Minimal satu angka.' })
    .refine((value) => /[^a-zA-Z0-9]/.test(value), { message: 'Minimal satu simbol.' }),
  confirmPassword: z.string(),
}).refine((data) => data.newPassword === data.confirmPassword, {
  message: 'Kata sandi baru dan konfirmasi kata sandi tidak cocok.',
  path: ['confirmPassword'],
});

const resolver = zodResolver(passwordSchema);

const initialValues = ref({
  newPassword: '',
  confirmPassword: '',
});

const handleResetPassword = async ({ valid, values: data }) => {
  if (!valid) {
    toast.add({ severity: 'error', summary: 'Validasi Gagal', detail: 'Silakan periksa kembali input Anda.', life: 3000 });
    return;
  }

  try {
    const response = await axios.post('/api/reset-password', {
      token: token.value,
      new_password: data.newPassword,
    });
    toast.add({ severity: 'success', summary: 'Success', detail: response.data.message || 'Kata sandi berhasil direset.', life: 3000 });
    router.push({ name: 'AuthPage' });
  } catch (error) {
    console.error('Reset password error:', error);
    if (error.response && error.response.status === 400) {
      router.replace({ name: 'TokenInvalid' });
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data?.message || 'Terjadi kesalahan saat mereset kata sandi.', life: 3000 });
    }
    tokenValid.value = false;
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