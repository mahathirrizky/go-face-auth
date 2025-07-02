<template>
  <section id="pricing" class="py-16 md:py-24 bg-blue-700 text-white">
    <div class="container mx-auto px-4">
      <h2 class="text-3xl md:text-4xl font-extrabold text-center mb-12 md:mb-16">
        Pilih Paket yang Sesuai untuk Bisnis Anda
      </h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-8 lg:gap-12">
        <!-- Paket Basic -->
        <div class="bg-white text-gray-800 p-8 rounded-xl shadow-lg flex flex-col transform hover:scale-105 transition duration-300 ease-in-out">
          <h3 class="text-2xl font-bold mb-4 text-center">Basic</h3>
          <p class="text-center text-gray-600 mb-6">Cocok untuk startup & bisnis kecil</p>
          <div class="text-center mb-8">
            <span class="text-5xl font-extrabold text-blue-600">Rp 50.000</span>
            <span class="text-xl text-gray-500">/bulan</span>
          </div>
          <ul class="text-left space-y-3 mb-8 flex-grow">
            <li class="flex items-center"><svg class="w-6 h-6 text-green-500 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg> Hingga 10 Karyawan</li>
            <li class="flex items-center"><svg class="w-6 h-6 text-green-500 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg> Absensi Online & GPS</li>
            <li class="flex items-center"><svg class="w-6 h-6 text-green-500 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg> Laporan Dasar</li>
            <li class="flex items-center text-gray-400"><svg class="w-6 h-6 text-red-500 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg> Pengenalan Wajah</li>
          </ul>
          <button class="mt-auto bg-blue-600 hover:bg-blue-700 text-white font-bold py-3 px-6 rounded-full text-lg transition duration-300 w-full">
            Pilih Paket Basic
          </button>
        </div>

        <!-- Paket Pro (Most Popular) -->
        <div class="bg-white text-gray-800 p-8 rounded-xl shadow-lg flex flex-col transform hover:scale-105 transition duration-300 ease-in-out border-4 border-blue-600 relative">
          <span class="absolute -top-4 left-1/2 -translate-x-1/2 bg-blue-600 text-white text-xs font-bold px-3 py-1 rounded-full uppercase">Paling Populer</span>
          <h3 class="text-2xl font-bold mb-4 text-center">Pro</h3>
          <p class="text-center text-gray-600 mb-6">Ideal untuk bisnis berkembang</p>
          <div class="text-center mb-8">
            <span class="text-5xl font-extrabold text-blue-600">Rp 150.000</span>
            <span class="text-xl text-gray-500">/bulan</span>
          </div>
          <ul class="text-left space-y-3 mb-8 flex-grow">
            <li class="flex items-center"><svg class="w-6 h-6 text-green-500 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg> Hingga 50 Karyawan</li>
            <li class="flex items-center"><svg class="w-6 h-6 text-green-500 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg> Absensi Online & GPS</li>
            <li class="flex items-center"><svg class="w-6 h-6 text-green-500 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg> Pengenalan Wajah</li>
            <li class="flex items-center"><svg class="w-6 h-6 text-green-500 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg> Laporan Lengkap</li>
          </ul>
          <button class="mt-auto bg-blue-600 hover:bg-blue-700 text-white font-bold py-3 px-6 rounded-full text-lg transition duration-300 w-full">
            Pilih Paket Pro
          </button>
        </div>

        <!-- Paket Enterprise -->
        <div class="bg-white text-gray-800 p-8 rounded-xl shadow-lg flex flex-col transform hover:scale-105 transition duration-300 ease-in-out">
          <h3 class="text-2xl font-bold mb-4 text-center">Enterprise</h3>
          <p class="text-center text-gray-600 mb-6">Solusi khusus untuk perusahaan besar</p>
          <div class="text-center mb-8">
            <span class="text-5xl font-extrabold text-blue-600">Custom</span>
            <span class="text-xl text-gray-500">/bulan</span>
          </div>
          <ul class="text-left space-y-3 mb-8 flex-grow">
            <li class="flex items-center"><svg class="w-6 h-6 text-green-500 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg> Karyawan Tidak Terbatas</li>
            <li class="flex items-center"><svg class="w-6 h-6 text-green-500 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg> Semua Fitur Pro</li>
            <li class="flex items-center"><svg class="w-6 h-6 text-green-500 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg> Integrasi API Kustom</li>
            <li class="flex items-center"><svg class="w-6 h-6 text-green-500 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg> Dukungan Prioritas</li>
          </ul>
          <button class="mt-auto bg-blue-600 hover:bg-blue-700 text-white font-bold py-3 px-6 rounded-full text-lg transition duration-300 w-full">
            Hubungi Kami
          </button>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  name: 'PricingSection',
};
</script>
