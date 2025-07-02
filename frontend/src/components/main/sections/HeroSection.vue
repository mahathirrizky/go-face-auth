<template>
  <section class="relative bg-gradient-to-r from-blue-600 to-blue-800 text-white py-20 md:py-32 overflow-hidden min-h-screen flex items-center">
    <div class="container mx-auto px-4 relative z-10 text-center">
      <h1 class="text-4xl md:text-6xl font-extrabold leading-tight mb-6 animate-fade-in-up">
        Solusi Absensi Modern untuk Perusahaan Anda
      </h1>
      <p class="text-lg md:text-xl mb-10 max-w-3xl mx-auto animate-fade-in-up animation-delay-200">
        Kelola absensi karyawan dengan mudah, cepat, dan akurat menggunakan teknologi pengenalan wajah.
      </p>
      <div class="flex flex-col sm:flex-row justify-center gap-4 animate-fade-in-up animation-delay-400">
        <button class="bg-white text-blue-700 hover:bg-gray-100 font-bold py-3 px-8 rounded-full text-lg transition duration-300 shadow-lg">
          Coba Gratis Sekarang
        </button>
        <button @click="scrollToFeatures" class="bg-transparent border-2 border-white text-white hover:bg-white hover:text-blue-700 font-bold py-3 px-8 rounded-full text-lg transition duration-300">
          Pelajari Lebih Lanjut
        </button>
      </div>
    </div>
    <!-- Background shapes/effects for high-tech feel -->
    <div class="absolute inset-0 z-0 opacity-10">
      <svg class="w-full h-full" viewBox="0 0 100 100" preserveAspectRatio="none">
        <defs>
          <radialGradient id="gradient1" cx="50%" cy="50%" r="50%" fx="50%" fy="50%">
            <stop offset="0%" stop-color="rgba(255,255,255,0.2)"></stop>
            <stop offset="100%" stop-color="rgba(255,255,255,0)"></stop>
          </radialGradient>
        </defs>
        <circle cx="20" cy="20" r="15" fill="url(#gradient1)"></circle>
        <circle cx="80" cy="80" r="20" fill="url(#gradient1)"></circle>
        <circle cx="50" cy="10" r="10" fill="url(#gradient1)"></circle>
        <circle cx="10" cy="90" r="12" fill="url(#gradient1)"></circle>
      </svg>
    </div>
  </section>
</template>

<script>
export default {
  name: 'HeroSection',
  methods: {
    scrollToFeatures() {
      document.getElementById('features').scrollIntoView({
        behavior: 'smooth'
      });
    }
  }
};
</script>

<style scoped>
@keyframes fadeInFromBottom {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-fade-in-up {
  animation: fadeInFromBottom 0.8s ease-out forwards;
  opacity: 0; /* Start invisible */
}

.animation-delay-200 {
  animation-delay: 0.2s;
}

.animation-delay-400 {
  animation-delay: 0.4s;
}
</style>
