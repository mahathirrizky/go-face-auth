<template>
  <div class="bg-bg-muted p-6 rounded-lg shadow-md">
    <h3 class="text-xl font-semibold text-text-base mb-4">Manajemen Akun Admin</h3>
    <div class="mb-4">
      <label class="block text-text-muted text-sm font-bold mb-2">Email Admin:</label>
      <span class="block w-full p-2 rounded-md border border-bg-base bg-bg-base text-text-base">{{ settings.adminEmail }}</span>
    </div>
    <PasswordInput
      id="oldPassword"
      label="Kata Sandi Lama:"
      v-model="settings.oldPassword"
    />
    <PasswordInput
      id="newPassword"
      label="Kata Sandi Baru:"
      v-model="settings.newPassword"
    />
    <BaseButton @click="changeAdminPassword" class="mt-4">
      <i class="fas fa-key"></i> Ubah Kata Sandi
    </BaseButton>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import axios from 'axios';
import { useToast } from 'vue-toastification';
import { useAuthStore } from '../../stores/auth';
import PasswordInput from '../ui/PasswordInput.vue';
import BaseButton from '../ui/BaseButton.vue';

const authStore = useAuthStore();
const toast = useToast();

const settings = ref({
  adminEmail: authStore.adminEmail || '',
  oldPassword: '',
  newPassword: '',
});

const changeAdminPassword = async () => {
  if (!settings.value.oldPassword || !settings.value.newPassword) {
    toast.error('Kata sandi lama dan kata sandi baru harus diisi.');
    return;
  }
  if (settings.value.newPassword.length < 6) {
    toast.error('Kata sandi baru minimal 6 karakter.');
    return;
  }

  try {
    const response = await axios.put('/api/admin/change-password', {
      old_password: settings.value.oldPassword,
      new_password: settings.value.newPassword,
    });
    if (response.data && response.data.status === 'success') {
      toast.success(response.data.message || 'Kata sandi berhasil diubah.');
      settings.value.oldPassword = '';
      settings.value.newPassword = '';
    } else {
      toast.error(response.data?.message || 'Gagal mengubah kata sandi.');
    }
  } catch (error) {
    console.error('Error changing admin password:', error);
    let message = 'Gagal mengubah kata sandi.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.error(message);
  }
};
</script>
