<template>
  <section id="testimonials" class="py-16 md:py-24 bg-gray-50">
    <div class="container mx-auto px-4">
      <h2 class="text-3xl md:text-4xl font-extrabold text-center text-gray-800 mb-12 md:mb-16">
        Apa Kata Klien Kami?
      </h2>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        <div class="bg-white p-8 rounded-xl shadow-lg flex flex-col items-center text-center">
          <img src="https://via.placeholder.com/100" alt="Client Avatar" class="w-24 h-24 rounded-full mb-6 border-4 border-blue-500 object-cover">
          <p class="text-lg text-gray-700 mb-4 italic">"HadirBos.com sangat membantu dalam mengelola absensi karyawan kami. Fitur pengenalan wajahnya sangat akurat dan efisien!"</p>
          <p class="font-bold text-gray-900">Budi Santoso</p>
          <p class="text-sm text-gray-500">HR Manager, PT Maju Jaya</p>
        </div>

        <div class="bg-white p-8 rounded-xl shadow-lg flex flex-col items-center text-center">
          <img src="https://via.placeholder.com/100" alt="Client Avatar" class="w-24 h-24 rounded-full mb-6 border-4 border-blue-500 object-cover">
          <p class="text-lg text-gray-700 mb-4 italic">"Laporan real-time yang disediakan sangat memudahkan kami dalam memantau kehadiran karyawan. Sangat direkomendasikan!"</p>
          <p class="font-bold text-gray-900">Siti Aminah</p>
          <p class="text-sm text-gray-500">Owner, Kedai Kopi Nikmat</p>
        </div>

        <div class="bg-white p-8 rounded-xl shadow-lg flex flex-col items-center text-center">
          <img src="https://via.placeholder.com/100" alt="Client Avatar" class="w-24 h-24 rounded-full mb-6 border-4 border-blue-500 object-cover">
          <p class="text-lg text-gray-700 mb-4 italic">"Integrasi yang mudah dan dukungan pelanggan yang responsif membuat pengalaman kami menggunakan HadirBos.com sangat menyenangkan."</p>
          <p class="font-bold text-gray-900">Joko Susilo</p>
          <p class="text-sm text-gray-500">CEO, CV Inovasi Digital</p>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  name: 'TestimonialsSection',
};
</script>
