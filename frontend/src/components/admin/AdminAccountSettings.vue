<template>
  <div class="bg-bg-muted p-6 rounded-lg shadow-md">
    <h3 class="text-xl font-semibold text-text-base mb-4">Manajemen Akun Admin</h3>
    <div class="mb-4">
      <label class="block text-text-muted text-sm font-bold mb-2">Email Admin:</label>
      <span class="block w-full p-2 rounded-md border border-bg-base bg-bg-base text-text-base">{{ initialValues.adminEmail }}</span>
    </div>


      <BaseForm :resolver="resolver" :initialValues="initialValues" @submit="changeAdminPassword" v-slot="{ $form }">
        <BaseInput
        id="oldPassword"
        name="oldPassword"
        label="Kata Sandi Lama:"
        type="password"
        :feedback="false"
        :toggleMask="true"
        :invalid="$form.oldPassword?.invalid"
        :errors="$form.oldPassword?.errors"
        :fluid="true"
        />
        <BaseInput
        id="newPassword"
        name="newPassword"
        label="Kata Sandi Baru:"
        type="password"
        :feedback="false"
        :toggleMask="true"
        :invalid="$form.newPassword?.invalid"
        :fluid="true"
        :errors="$form.newPassword?.errors"
        />
        <BaseInput
        id="confirmNewPassword"
        name="confirmNewPassword"
        label="Konfirmasi Kata Sandi Baru:"
        type="password"
        :feedback="false"
        :toggleMask="true"
        :invalid="$form.confirmNewPassword?.invalid"
        :fluid="true"
        :errors="$form.confirmNewPassword?.errors"
        />
        <BaseButton type="submit" class="mt-4">
          <i class="pi pi-key"></i> Ubah Kata Sandi
        </BaseButton>
      </BaseForm>
    </div>

</template>

<script setup>
import { ref } from 'vue';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { useAuthStore } from '../../stores/auth';
import BaseButton from '../ui/BaseButton.vue';
import BaseInput from '../ui/BaseInput.vue';
import BaseForm from '../ui/BaseForm.vue'; // Import BaseForm
import { zodResolver } from '@primevue/forms/resolvers/zod';
import { z } from 'zod';

const authStore = useAuthStore();
const toast = useToast();

const initialValues = ref({
  adminEmail: authStore.adminEmail || '',
  oldPassword: '',
  newPassword: '',
  confirmNewPassword: '',
});

const passwordSchema = z.object({
  oldPassword: z.string().min(1, { message: 'Kata sandi lama wajib diisi.' }),
  newPassword: z.string()
    .min(8, { message: 'Minimal 8 karakter.' })
    .refine((value) => /[a-z]/.test(value), {
      message: 'Minimal satu huruf kecil.'
    })
    .refine((value) => /[A-Z]/.test(value), {
      message: 'Minimal satu huruf besar.'
    })
    .refine((value) => /\d/.test(value), {
      message: 'Minimal satu angka.'
    }),
  confirmNewPassword: z.string(),
}).refine((data) => data.newPassword === data.confirmNewPassword, {
  message: 'Kata sandi baru dan konfirmasi kata sandi tidak cocok.',
  path: ['confirmNewPassword'],
});

const resolver = zodResolver(passwordSchema);

const changeAdminPassword = async (event) => {
  const { valid, data } = event;

  if (!valid) {
    toast.add({ severity: 'error', summary: 'Validasi Gagal', detail: 'Silakan periksa kembali input Anda.', life: 3000 });
    return;
  }

  try {
    const response = await axios.put('/api/admin/change-password', {
      old_password: data.oldPassword,
      new_password: data.newPassword,
    });
    if (response.data && response.data.status === 'success') {
      toast.add({ severity: 'success', summary: 'Success', detail: response.data.message || 'Kata sandi berhasil diubah.', life: 3000 });
      initialValues.value.oldPassword = '';
      initialValues.value.newPassword = '';
      initialValues.value.confirmNewPassword = '';
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data?.message || 'Gagal mengubah kata sandi.', life: 3000 });
    }
  } catch (error) {
    console.error('Error changing admin password:', error);
    let message = 'Gagal mengubah kata sandi.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  }
};
</script>
