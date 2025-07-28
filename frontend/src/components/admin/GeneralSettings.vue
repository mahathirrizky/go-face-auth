<template>
  <div class="bg-bg-muted p-6 rounded-lg shadow-md mt-6">
    <h3 class="text-xl font-semibold text-text-base mb-4">Pengaturan Umum Perusahaan</h3>
    <div class="mb-4">
      <label for="companyName" class="block text-text-muted text-sm font-bold mb-2">Nama Perusahaan:</label>
      <InputText
        id="companyName"
        v-model="settings.companyName"
        class="w-full"
      />
    </div>
    <div class="mb-4">
      <label for="companyAddress" class="block text-text-muted text-sm font-bold mb-2">Alamat Perusahaan:</label>
      <InputText
        id="companyAddress"
        v-model="settings.companyAddress"
        class="w-full"
      />
    </div>
    <FloatLabel class="mb-4" variant="on">
      
      <Select
        id="timezone"
        v-model="settings.timezone"
        :options="timezones"
        optionLabel="label"
        optionValue="value"
        placeholder="Pilih Zona Waktu"
        class="w-full"
      />
      <label for="timezone" class="block text-text-muted text-sm font-bold mb-2">Zona Waktu:</label>
    </FloatLabel>
    <BaseButton @click="saveSettings" class="mt-4" :disabled="isSaving">
      <i v-if="!isSaving" class="pi pi-save"></i>
      <i v-else class="pi pi-spin pi-spinner"></i>
      {{ isSaving ? 'Menyimpan...' : 'Simpan Pengaturan' }}
    </BaseButton>
  </div>

  <div class="bg-bg-muted p-6 rounded-lg shadow-md mt-6">
    <h3 class="text-xl font-semibold text-text-base mb-4">Manajemen Akun Admin</h3>
    <div class="mb-4">
      <label class="block text-text-muted text-sm font-bold mb-2">Email Admin:</label>
      <span class="block w-full p-2 rounded-md border border-bg-base bg-bg-base text-text-base">{{ adminAccountForm.adminEmail }}</span>
    </div>


      <BaseForm :resolver="passwordResolver" :initialValues="adminAccountForm" @submit="changeAdminPassword" v-slot="{ $form }">
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
        <BaseButton type="submit" class="mt-4" :disabled="isChangingPassword">
          <i v-if="!isChangingPassword" class="pi pi-key"></i>
          <i v-else class="pi pi-spin pi-spinner"></i>
          {{ isChangingPassword ? 'Mengubah...' : 'Ubah Kata Sandi' }}
        </BaseButton>
      </BaseForm>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useToast } from 'primevue/usetoast';
import { useAuthStore } from '../../stores/auth';
import axios from 'axios';
import InputText from 'primevue/inputtext';
import Select from 'primevue/select';
import BaseButton from '../ui/BaseButton.vue';
import FloatLabel from 'primevue/floatlabel';
import BaseInput from '../ui/BaseInput.vue';
import BaseForm from '../ui/BaseForm.vue';
import { zodResolver } from '@primevue/forms/resolvers/zod';
import { z } from 'zod';

const authStore = useAuthStore();
const toast = useToast();
const isSaving = ref(false);
const isChangingPassword = ref(false);

const settings = ref({
  companyName: authStore.companyName || '',
  companyAddress: authStore.companyAddress || '',
  timezone: authStore.companyTimezone || 'Asia/Jakarta',
});

const timezones = ref([
  { value: 'Asia/Jakarta', label: 'Asia/Jakarta (WIB)' },
  { value: 'Asia/Makassar', label: 'Asia/Makassar (WITA)' },
  { value: 'Asia/Jayapura', label: 'Asia/Jayapura (WIT)' },
  { value: 'UTC', label: 'UTC' },
]);

const saveSettings = async () => {
  isSaving.value = true;
  try {
    const payload = {
      name: settings.value.companyName,
      address: settings.value.companyAddress,
      timezone: settings.value.timezone,
    };
    const response = await axios.put('/api/company-details', payload, {
      headers: { Authorization: `Bearer ${authStore.token}` },
    });
    if (response.data && response.data.status === 'success') {
      toast.add({ severity: 'success', summary: 'Success', detail: 'Pengaturan umum berhasil disimpan!', life: 3000 });
      authStore.companyName = settings.value.companyName;
      authStore.companyAddress = settings.value.companyAddress;
      authStore.companyTimezone = settings.value.timezone;
      authStore.hasConfiguredTimezone = true;
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data.message || 'Gagal menyimpan pengaturan umum.', life: 3000 });
    }
  } catch (error) {
    console.error('Error saving general settings:', error);
    toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data?.message || 'Terjadi kesalahan saat menyimpan pengaturan.', life: 3000 });
  } finally {
    isSaving.value = false;
  }
};

// Admin Account Management Logic
const adminAccountForm = ref({
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

const passwordResolver = zodResolver(passwordSchema);

const changeAdminPassword = async (event) => {
  const { valid, data } = event;

  if (!valid) {
    toast.add({ severity: 'error', summary: 'Validasi Gagal', detail: 'Silakan periksa kembali input Anda.', life: 3000 });
    return;
  }

  isChangingPassword.value = true; // Start loading
  try {
    const response = await axios.put('/api/admin/change-password', {
      old_password: data.oldPassword,
      new_password: data.newPassword,
    });
    if (response.data && response.data.status === 'success') {
      toast.add({ severity: 'success', summary: 'Success', detail: response.data.message || 'Kata sandi berhasil diubah.', life: 3000 });
      adminAccountForm.value.oldPassword = '';
      adminAccountForm.value.newPassword = '';
      adminAccountForm.value.confirmNewPassword = '';
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
  } finally {
    isChangingPassword.value = false; // End loading
  }
};
</script>