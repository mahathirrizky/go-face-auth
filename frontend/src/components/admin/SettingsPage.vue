<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-6">Pengaturan</h2>

    <div class="bg-bg-muted p-6 rounded-lg shadow-md">
      <h3 class="text-xl font-semibold text-text-base mb-4">Pengaturan Umum</h3>
      <div class="mb-4">
        <label for="companyName" class="block text-text-muted text-sm font-bold mb-2">Nama Perusahaan:</label>
        <input
          type="text"
          id="companyName"
          v-model="settings.companyName"
          class="w-full p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary"
        />
      </div>
      <div class="mb-4">
        <label for="companyAddress" class="block text-text-muted text-sm font-bold mb-2">Alamat Perusahaan:</label>
        <input
          type="text"
          id="companyAddress"
          v-model="settings.companyAddress"
          class="w-full p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary"
        />
      </div>
      <button class="btn btn-secondary mt-4">
        Simpan Pengaturan
      </button>
    </div>

    <div class="bg-bg-muted p-6 rounded-lg shadow-md mt-6">
      <h3 class="text-xl font-semibold text-text-base mb-4">Manajemen Akun Admin</h3>
      <div class="mb-4">
        <label class="block text-text-muted text-sm font-bold mb-2">Email Admin:</label>
        <span class="block w-full p-2 rounded-md border border-bg-base bg-bg-base text-text-base">{{ settings.adminEmail }}</span>
      </div>
      <div class="mb-4">
        <label for="oldPassword" class="block text-text-muted text-sm font-bold mb-2">Kata Sandi Lama:</label>
        <div class="relative">
          <input
            :type="showOldPassword ? 'text' : 'password'"
            id="oldPassword"
            v-model="settings.oldPassword"
            class="w-full p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary pr-10"
          />
          <div class="absolute inset-y-0 right-0 pr-3 flex items-center cursor-pointer" @click="toggleOldPasswordVisibility">
            <svg v-if="showOldPassword" class="h-5 w-5 text-text-muted" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.418 0-8-2.915-8-6.5S7.582 6 12 6c1.832 0 3.533.61 4.905 1.635M19 10.5c0 3.585-3.582 6.5-8 6.5a10.05 10.05 0 01-1.875-.175M18 18l-4.243-4.243M12 12l-1.414-1.414M12 12l1.414 1.414M12 12l-1.414 1.414M12 12l1.414 1.414" />
            </svg>
            <svg v-else class="h-5 w-5 text-text-muted" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
            </svg>
          </div>
        </div>
      </div>
      <div class="mb-4">
        <label for="newPassword" class="block text-text-muted text-sm font-bold mb-2">Kata Sandi Baru:</label>
        <div class="relative">
          <input
            :type="showNewPassword ? 'text' : 'password'"
            id="newPassword"
            v-model="settings.newPassword"
            class="w-full p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary pr-10"
          />
          <div class="absolute inset-y-0 right-0 pr-3 flex items-center cursor-pointer" @click="toggleNewPasswordVisibility">
            <svg v-if="showNewPassword" class="h-5 w-5 text-text-muted" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.418 0-8-2.915-8-6.5S7.582 6 12 6c1.832 0 3.533.61 4.905 1.635M19 10.5c0 3.585-3.582 6.5-8 6.5a10.05 10.05 0 01-1.875-.175M18 18l-4.243-4.243M12 12l-1.414-1.414M12 12l1.414 1.414M12 12l-1.414 1.414M12 12l1.414 1.414" />
            </svg>
            <svg v-else class="h-5 w-5 text-text-muted" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
            </svg>
          </div>
        </div>
      </div>
      </div>
      <button @click="changeAdminPassword" class="btn btn-secondary mt-4">
        Ubah Kata Sandi
      </button>
    </div>

    <div class="bg-bg-muted p-6 rounded-lg shadow-md mt-6">
      <h3 class="text-xl font-semibold text-text-base mb-4">Manajemen Shift</h3>
      <p class="text-text-muted mb-4">Kelola jadwal shift kerja karyawan Anda.</p>
      <router-link to="/dashboard/shifts" class="btn btn-secondary">
        Buka Manajemen Shift
      </router-link>
    </div>
  
</template>

<script>
import { ref, onMounted } from 'vue';
import axios from 'axios'; // Import axios
import { useToast } from 'vue-toastification';
import { useAuthStore } from '../../stores/auth';

export default {
  name: 'SettingsPage',
  setup() {
    const authStore = useAuthStore();
    const toast = useToast();

    const settings = ref({
      companyName: authStore.companyName || '',
      companyAddress: authStore.companyAddress || '',
      adminEmail: authStore.adminEmail || '',
      oldPassword: '',
      newPassword: '',
    });

    const showOldPassword = ref(false);
    const showNewPassword = ref(false);

    const toggleOldPasswordVisibility = () => {
      showOldPassword.value = !showOldPassword.value;
    };

    const toggleNewPasswordVisibility = () => {
      showNewPassword.value = !showNewPassword.value;
    };

    const saveSettings = async () => {
      // This function will be implemented when there's a backend endpoint to update company details
      toast.info('Fungsi simpan pengaturan umum belum diimplementasikan.');
      console.log('Saving general settings:', settings.value);
    };

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

    return {
      settings,
      saveSettings,
      changeAdminPassword,
      showOldPassword,
      showNewPassword,
      toggleOldPasswordVisibility,
      toggleNewPasswordVisibility,
    };
  },
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>
