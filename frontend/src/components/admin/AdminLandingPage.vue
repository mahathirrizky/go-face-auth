<template>
  <div class="admin-landing-page flex flex-col items-center justify-center min-h-[80vh] text-center bg-blue-100 text-blue-800 font-sans">
    <h1 class="text-5xl text-blue-900 mb-4">Selamat Datang di Panel Admin!</h1>
    <p class="text-xl mb-6 max-w-2xl">Ini adalah halaman arahan untuk administrasi.</p>
    <button class="bg-blue-600 hover:bg-blue-700 text-white font-bold py-4 px-8 rounded-lg cursor-pointer transition-colors duration-300">
      Kelola Data
    </button>
  </div>
</template>

<script>
export default {
  name: 'AdminLandingPage',
}
</script>

<style scoped>
/* Tailwind handles styling */
</style>
