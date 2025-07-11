<template>
  <div class="min-h-screen flex flex-col items-center justify-center bg-bg-base text-center p-4">
    <h1 class="text-9xl font-extrabold text-secondary tracking-wider">404</h1>
    <h2 class="text-4xl font-bold text-text-base mt-4 mb-2">Halaman Tidak Ditemukan</h2>
    <p class="text-lg text-text-muted mb-8 max-w-md">
      Maaf, halaman yang Anda cari tidak ada atau mungkin telah dipindahkan.
    </p>
    <router-link to="/">
      <BaseButton>
        Kembali ke Halaman Utama
      </BaseButton>
    </router-link>
  </div>
</template>

<script setup>
import BaseButton from '../ui/BaseButton.vue';
</script>

<style scoped>
/* Tailwind handles styling */
</style>
