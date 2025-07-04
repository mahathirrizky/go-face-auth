<template>
  <section id="pricing" class="py-16 md:py-24 bg-blue-700 text-white">
    <div class="container mx-auto px-4">
      <h2 class="text-3xl md:text-4xl font-extrabold text-center mb-12 md:mb-16">
        Pilih Paket yang Sesuai untuk Bisnis Anda
      </h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-8 lg:gap-12">
        <div
          v-for="pkg in packages"
          :key="pkg.id"
          :class="{ 'border-8 border-blue-500 shadow-2xl relative': pkg.name === 'Standard' }"
          class="bg-white text-gray-800 p-8 rounded-xl shadow-lg flex flex-col transform hover:scale-105 transition duration-300 ease-in-out"
        >
          <span
            v-if="pkg.name === 'Standard'"
            class="absolute -top-6 left-1/2 -translate-x-1/2 bg-blue-500 text-white text-sm font-extrabold px-4 py-2 rounded-full uppercase"
          >
            Paling Populer
          </span>
          <h3 class="text-2xl font-bold mb-4 text-center">{{ pkg.name }}</h3>
          <p class="text-center text-gray-600 mb-6">
            {{ pkg.name === 'Basic' ? 'Cocok untuk startup & bisnis kecil' : pkg.name === 'Standard' ? 'Ideal untuk bisnis berkembang' : 'Solusi khusus untuk perusahaan besar' }}
          </p>
          <div class="text-center mb-8">
            <span class="text-5xl font-extrabold text-blue-600">
              {{ pkg.price === 0 ? 'Gratis' : `Rp ${pkg.price}` }}
            </span>
            <span class="text-xl text-gray-500" v-if="pkg.price !== 0">/bulan</span>
          </div>
          <ul class="text-left space-y-3 mb-8 flex-grow">
            <li class="flex items-center">
              <svg class="w-6 h-6 text-green-500 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
              Hingga {{ pkg.max_employees }} Karyawan
            </li>
            <li class="flex items-center" v-for="(feature, index) in pkg.features.split(',')" :key="index">
              <svg class="w-6 h-6 text-green-500 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
              {{ feature.trim() }}
            </li>
          </ul>
          <button
            @click="selectPackage(pkg.id)"
            class="mt-auto bg-blue-600 hover:bg-blue-700 text-white font-bold py-3 px-6 rounded-full text-lg transition duration-300 w-full"
          >
            Pilih Paket {{ pkg.name }}
          </button>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  name: 'PricingSection',
  props: ['packages'],
  methods: {
    selectPackage(packageId) {
      console.log('Selecting package with ID:', packageId);
      this.$router.push({ name: 'RegisterCompany', params: { packageId: packageId } });
    },
  },
};
</script>

<style scoped>
/* No custom scoped styles needed as Tailwind handles styling */
</style>