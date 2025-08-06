<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-6">Pengaturan</h2>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <!-- General Settings Card -->
      <router-link to="/dashboard/settings/general" class=" bg-bg-muted p-4 rounded-lg shadow-md hover:shadow-lg transition-shadow duration-200 cursor-pointer flex items-center space-x-3"
        :class="{ 'ring-2 ring-primary': $route.path.includes('/dashboard/settings/general') }">
        <i class="pi pi-cog text-3xl text-blue-500"></i>
        <div>
          <h3 class="text-lg font-semibold text-text-base mb-1">Pengaturan Umum</h3>
        </div>
      </router-link>

      <!-- Shift Management Card -->
      <router-link to="/dashboard/settings/shifts" class=" bg-bg-muted p-4 rounded-lg shadow-md hover:shadow-lg transition-shadow duration-200 cursor-pointer flex items-center space-x-3"
        :class="{ 'ring-2 ring-primary': $route.path.includes('/dashboard/settings/shifts') }">
        <i class="pi pi-clock text-3xl text-green-500"></i>
        <div>
          <h3 class="text-lg font-semibold text-text-base mb-1">Manajemen Shift</h3>
        </div>
      </router-link>

      <!-- Location Management Card -->
      <router-link to="/dashboard/settings/locations" class=" bg-bg-muted p-4 rounded-lg shadow-md hover:shadow-lg transition-shadow duration-200 cursor-pointer flex items-center space-x-3"
        :class="{ 'ring-2 ring-primary': $route.path.includes('/dashboard/settings/locations') }">
        <i class="pi pi-map-marker text-3xl text-red-500"></i>
        <div>
          <h3 class="text-lg font-semibold text-text-base mb-1">Manajemen Lokasi Absensi</h3>
        </div>
      </router-link>

      

      
    </div>

    <!-- Router View for nested settings components -->
    <router-view class="mt-6"></router-view>
  </div>
</template>

<script setup>
import { ref } from 'vue';
</script>

<style scoped>
/* Tailwind handles styling */
</style>