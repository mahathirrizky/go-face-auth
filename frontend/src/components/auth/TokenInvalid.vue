<template>
  <div class="min-h-screen flex items-center justify-center bg-bg-base py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8 p-10 bg-bg-muted rounded-lg shadow-xl text-center">
      <h2 class="mt-6 text-center text-3xl font-extrabold text-text-base">
        Tautan Tidak Valid atau Kedaluwarsa
      </h2>
      <p class="mt-2 text-center text-sm text-text-muted">
        Tautan reset kata sandi yang Anda gunakan tidak valid, sudah kedaluwarsa, atau sudah digunakan.
        Silakan coba lagi proses reset kata sandi atau hubungi administrator.
      </p>
      <div class="mt-6">
        <router-link to="/admin/forgot-password" class="font-medium text-accent hover:text-accent-dark">
          Minta Tautan Reset Kata Sandi Baru
        </router-link>
      </div>
      <div class="mt-4">
        <router-link to="/admin" class="inline-block align-baseline font-bold text-sm text-accent hover:opacity-90">
          Kembali ke Login Admin
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
// No specific script logic needed for this simple page
</script>

<style scoped>
/* Tailwind handles styling */
</style>
