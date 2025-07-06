<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-6">Pengaturan</h2>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <!-- General Settings Card -->
      <router-link to="/dashboard/settings/general" class="block bg-bg-muted p-6 rounded-lg shadow-md hover:shadow-lg transition-shadow duration-200 cursor-pointer"
        :class="{ 'ring-2 ring-primary': $route.path.includes('/dashboard/settings/general') }">
        <h3 class="text-xl font-semibold text-text-base mb-2">Pengaturan Umum</h3>
        <p class="text-text-muted">Kelola nama dan alamat perusahaan Anda.</p>
      </router-link>

      <!-- Admin Account Settings Card -->
      <router-link to="/dashboard/settings/admin-account" class="block bg-bg-muted p-6 rounded-lg shadow-md hover:shadow-lg transition-shadow duration-200 cursor-pointer"
        :class="{ 'ring-2 ring-primary': $route.path.includes('/dashboard/settings/admin-account') }">
        <h3 class="text-xl font-semibold text-text-base mb-2">Manajemen Akun Admin</h3>
        <p class="text-text-muted">Ubah kata sandi dan kelola detail akun admin.</p>
      </router-link>

      <!-- Shift Management Card -->
      <router-link to="/dashboard/settings/shifts" class="block bg-bg-muted p-6 rounded-lg shadow-md hover:shadow-lg transition-shadow duration-200 cursor-pointer"
        :class="{ 'ring-2 ring-primary': $route.path.includes('/dashboard/settings/shifts') }">
        <h3 class="text-xl font-semibold text-text-base mb-2">Manajemen Shift</h3>
        <p class="text-text-muted">Kelola jadwal shift kerja karyawan Anda.</p>
      </router-link>
    </div>

    <!-- Router View for nested settings components -->
    <router-view class="mt-6"></router-view>
  </div>
</template>

<script>
import { ref } from 'vue';

export default {
  name: 'SettingsPage',
  setup() {
    // This component now acts as a navigation hub, so no specific settings refs are needed here
    return {};
  },
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>