<template>
  <section id="pricing" class="py-20 bg-white">
    <div class="container mx-auto px-4">
      <h2 class="text-3xl font-bold text-center text-gray-800 mb-12">Pilihan Harga</h2>
      <div class="flex flex-wrap justify-center -mx-4">
        <!-- Pricing Card 1 -->
        <div class="w-full md:w-1/3 px-4 mb-8">
          <div class="border rounded-lg p-8 text-center shadow-lg hover:shadow-2xl transition-shadow duration-300 h-full flex flex-col">
            <h3 class="text-2xl font-bold mb-2">Basic</h3>
            <p class="text-4xl font-bold mb-4">$10<span class="text-lg font-normal">/bln</span></p>
            <ul class="text-gray-600 mb-6 text-left">
              <li class="py-1">✓ 100 Pengguna</li>
              <li class="py-1">✓ Otentikasi Dasar</li>
              <li class="py-1">✓ Dukungan Email</li>
            </ul>
            <BaseButton class="mt-auto btn-primary"><i class="fas fa-check"></i> Pilih Paket</BaseButton>
          </div>
        </div>
        <!-- Pricing Card 2 -->
        <div class="w-full md:w-1/3 px-4 mb-8">
          <div class="border-2 border-blue-600 rounded-lg p-8 text-center shadow-lg hover:shadow-2xl transition-shadow duration-300 relative h-full flex flex-col">
             <span class="bg-blue-600 text-white text-xs font-bold uppercase px-3 py-1 rounded-full absolute top-0 right-4 -mt-3">Populer</span>
            <h3 class="text-2xl font-bold mb-2">Pro</h3>
            <p class="text-4xl font-bold mb-4">$25<span class="text-lg font-normal">/bln</span></p>
            <ul class="text-gray-600 mb-6 text-left">
              <li class="py-1">✓ 1000 Pengguna</li>
              <li class="py-1">✓ Otentikasi Lanjutan</li>
              <li class="py-1">✓ Dukungan Prioritas</li>
            </ul>
            <BaseButton class="mt-auto btn-primary"><i class="fas fa-check"></i> Pilih Paket</BaseButton>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import BaseButton from '../../ui/BaseButton.vue';
</script>
