<template>
  <div class="min-h-screen flex items-center justify-center bg-bg-base py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8 p-10 bg-bg-muted rounded-lg shadow-xl text-center">
      <div class="flex justify-center mb-6">
        <img class="h-20 w-auto" src="/vite.svg" alt="Logo" />
      </div>
      <h2 class="mt-2 text-center text-3xl font-extrabold text-text-base">
        Selamat Bergabung!
      </h2>
      <p class="mt-2 text-center text-lg text-text-muted">
        Silakan akses akun Anda melalui Mobile App.
      </p>
      <div class="mt-6">
        <router-link to="/login" class="text-secondary hover:underline">
          Kembali ke Halaman Login
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
// No script needed for this simple success page
</script>

<style scoped>
/* Tailwind handles styling */
</style>
