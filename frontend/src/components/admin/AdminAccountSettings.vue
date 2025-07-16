<template>
  <div class="bg-bg-muted p-6 rounded-lg shadow-md">
    <h3 class="text-xl font-semibold text-text-base mb-4">Manajemen Akun Admin</h3>
    <div class="mb-4">
      <label class="block text-text-muted text-sm font-bold mb-2">Email Admin:</label>
      <span class="block w-full p-2 rounded-md border border-bg-base bg-bg-base text-text-base">{{ settings.adminEmail }}</span>
    </div>
    <BaseInput
      id="oldPassword"
      label="Kata Sandi Lama:"
      v-model="settings.oldPassword"
      type="password"
      :feedback="false"
      :toggleMask="true"
    />
    <BaseInput
      id="newPassword"
      label="Kata Sandi Baru:"
      v-model="settings.newPassword"
      type="password"
      :feedback="true"
      :toggleMask="true"
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
    </BaseInput>
    <BaseButton @click="changeAdminPassword" class="mt-4">
      <i class="pi pi-key"></i> Ubah Kata Sandi
    </BaseButton>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { useAuthStore } from '../../stores/auth';
import BaseButton from '../ui/BaseButton.vue';
import Divider from 'primevue/divider';
import BaseInput from '../ui/BaseInput.vue';

const authStore = useAuthStore();
const toast = useToast();

const settings = ref({
  adminEmail: authStore.adminEmail || '',
  oldPassword: '',
  newPassword: '',
});

const changeAdminPassword = async () => {
  if (!settings.value.oldPassword || !settings.value.newPassword) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Kata sandi lama dan kata sandi baru harus diisi.', life: 3000 });
    return;
  }
  if (settings.value.newPassword.length < 6) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Kata sandi baru minimal 6 karakter.', life: 3000 });
    return;
  }

  try {
    const response = await axios.put('/api/admin/change-password', {
      old_password: settings.value.oldPassword,
      new_password: settings.value.newPassword,
    });
    if (response.data && response.data.status === 'success') {
      toast.add({ severity: 'success', summary: 'Success', detail: response.data.message || 'Kata sandi berhasil diubah.', life: 3000 });
      settings.value.oldPassword = '';
      settings.value.newPassword = '';
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
