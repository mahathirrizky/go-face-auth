<template>
  <div>
    <div class="bg-bg-muted p-6 rounded-lg shadow-md mt-6">
      <h3 class="text-xl font-semibold text-text-base mb-4">Pengaturan Umum Perusahaan</h3>
      <div class="p-fluid">
        <div class="field mb-4">
          <FloatLabel variant="on">
            <label for="companyName" class="block text-text-muted text-sm font-bold mb-2">Nama Perusahaan:</label>
            <InputText id="companyName" v-model="settings.companyName" fluid />
          </FloatLabel>
        </div>
        <div class="field mb-4">
          <FloatLabel variant="on">
            <label for="companyAddress" class="block text-text-muted text-sm font-bold mb-2">Alamat Perusahaan:</label>
            <InputText id="companyAddress" v-model="settings.companyAddress" fluid />
          </FloatLabel>
        </div>
        <div class="field mb-4">
          <FloatLabel variant="on">
            <Select id="timezone" v-model="settings.timezone" :options="timezones" optionLabel="label" optionValue="value"  fluid />
            <label for="timezone" class="block text-text-muted text-sm font-bold mb-2">Zona Waktu:</label>
          </FloatLabel>
        </div>
        <Button @click="saveSettings" class="mt-4" :loading="isSaving" :label="isSaving ? 'Menyimpan...' : 'Simpan Pengaturan'" icon="pi pi-save" />
      </div>
    </div>

    <div class="bg-bg-muted p-6 rounded-lg shadow-md mt-6">
      <h3 class="text-xl font-semibold text-text-base mb-4">Manajemen Akun Admin</h3>
      <div class="mb-4">
        <FloatLabel variant="on">
          <InputText v-model="value" disabled :placeholder="adminAccountForm.adminEmail" fluid />
          <label class="block text-text-muted text-sm font-bold mb-2">Email Admin:</label>
        </FloatLabel>
      </div>

      <Form :resolver="passwordResolver" :initialValues="adminAccountForm" @submit="changeAdminPassword" v-slot="{ errors = {}, handleSubmit }">
        <div class="p-fluid">
          <div class="field mb-4">
            <FloatLabel variant="on">   
              <Password id="oldPassword" name="oldPassword" v-model="adminAccountForm.oldPassword" :feedback="false" :toggleMask="true" :invalid="!!errors?.oldPassword" fluid />
              <label for="oldPassword" class="block text-text-muted text-sm font-bold mb-2">Kata Sandi Lama:</label>
              <small class="p-error" v-if="errors?.oldPassword">{{ errors?.oldPassword }}</small>
            </FloatLabel>
          </div>
          <div class="field mb-4">
            <FloatLabel variant="on">
              <Password id="newPassword" name="newPassword" v-model="adminAccountForm.newPassword" :feedback="true" :toggleMask="true" :invalid="!!errors?.newPassword" fluid>
                <template #footer>
                  <p class="mt-2">Saran Kata Sandi:</p>
                  <ul class="pl-4 ml-2 mt-0" style="line-height: 1.5">
                    <li :class="{ 'text-green-500': isLengthValid }">Minimal 8 karakter <i v-if="isLengthValid" class="pi pi-check"></i></li>
                    <li :class="{ 'text-green-500': hasLowercase }">Minimal satu huruf kecil (a-z) <i v-if="hasLowercase" class="pi pi-check"></i></li>
                    <li :class="{ 'text-green-500': hasUppercase }">Minimal satu huruf besar (A-Z) <i v-if="hasUppercase" class="pi pi-check"></i></li>
                    <li :class="{ 'text-green-500': hasNumber }">Minimal satu angka (0-9) <i v-if="hasNumber" class="pi pi-check"></i></li>
                  </ul>
                  <small class="p-error" v-if="errors?.newPassword">{{ errors?.newPassword }}</small>
                </template>
              </Password>
              <label for="newPassword"  class="block text-text-muted text-sm font-bold mb-2">Kata Sandi Baru:</label>
            </FloatLabel>  
          </div>
          <div class="field mb-4">
            <FloatLabel variant="on">

              <Password id="confirmNewPassword" name="confirmNewPassword" v-model="adminAccountForm.confirmNewPassword" :feedback="false" :toggleMask="true" :invalid="!!errors?.confirmNewPassword" fluid />
              <label for="confirmNewPassword" class="block text-text-muted text-sm font-bold mb-2">Konfirmasi Kata Sandi Baru:</label>
              <small class="p-error" v-if="errors?.confirmNewPassword">{{ errors?.confirmNewPassword }}</small>
            </FloatLabel>
          </div>
          <Button type="submit" class="mt-4" :loading="isChangingPassword" :label="isChangingPassword ? 'Mengubah...' : 'Ubah Kata Sandi'" icon="pi pi-key" />
        </div>
      </Form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useToast } from 'primevue/usetoast';
import { useAuthStore } from '../../stores/auth';
import axios from 'axios';
import InputText from 'primevue/inputtext';
import Select from 'primevue/select';
import Button from 'primevue/button';
import Password from 'primevue/password';
import { Form } from '@primevue/forms';
import { zodResolver } from '@primevue/forms/resolvers/zod';
import { z } from 'zod';
import FloatLabel from 'primevue/floatlabel';

const authStore = useAuthStore();
const toast = useToast();
const isSaving = ref(false);
const isChangingPassword = ref(false);
const value = ref(''); 

const isLengthValid = computed(() => adminAccountForm.value.newPassword.length >= 8);
const hasLowercase = computed(() => /[a-z]/.test(adminAccountForm.value.newPassword));
const hasUppercase = computed(() => /[A-Z]/.test(adminAccountForm.value.newPassword));
const hasNumber = computed(() => /\d/.test(adminAccountForm.value.newPassword)); 

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
    .refine((value) => /[a-z]/.test(value), { message: 'Minimal satu huruf kecil.' })
    .refine((value) => /[A-Z]/.test(value), { message: 'Minimal satu huruf besar.' })
    .refine((value) => /\d/.test(value), { message: 'Minimal satu angka.' }),
  confirmNewPassword: z.string(),
}).refine((data) => data.newPassword === data.confirmNewPassword, {
  message: 'Kata sandi baru dan konfirmasi kata sandi tidak cocok.',
  path: ['confirmNewPassword'],
});

const passwordResolver = zodResolver(passwordSchema);

const changeAdminPassword = async ({ values }) => {
  isChangingPassword.value = true;
  try {
    const response = await axios.put('/api/admin/change-password', {
      old_password: values.oldPassword,
      new_password: values.newPassword,
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
    let message = error.response?.data?.message || 'Gagal mengubah kata sandi.';
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  } finally {
    isChangingPassword.value = false;
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