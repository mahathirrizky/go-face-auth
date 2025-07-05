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
        <label for="adminEmail" class="block text-text-muted text-sm font-bold mb-2">Email Admin:</label>
        <input
          type="email"
          id="adminEmail"
          v-model="settings.adminEmail"
          class="w-full p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary"
        />
      </div>
      <div class="mb-4">
        <label for="newPassword" class="block text-text-muted text-sm font-bold mb-2">Kata Sandi Baru:</label>
        <input
          type="password"
          id="newPassword"
          v-model="settings.newPassword"
          class="w-full p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary"
        />
      </div>
      <button class="btn btn-secondary mt-4">
        Ubah Kata Sandi
      </button>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue';

export default {
  name: 'SettingsPage',
  setup() {
    const settings = ref({
      companyName: 'PT. HadirBos Indonesia',
      companyAddress: 'Jl. Contoh No. 123, Jakarta',
      adminEmail: 'admin@hadirbos.com',
      newPassword: '',
    });

    return {
      settings,
    };
  },
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>
