<template>
  <section id="pricing" class="py-16 md:py-24 bg-bg-muted text-text-base">
    <div class="container mx-auto px-4">
      <h2 class="text-3xl md:text-4xl font-extrabold text-center mb-6">
        Pilih Paket yang Sesuai untuk Bisnis Anda
      </h2>
      
      <!-- Billing Cycle Toggle -->
      <div class="flex justify-center items-center space-x-4 mb-12">
        <span :class="{ 'text-secondary font-bold': billingCycle === 'monthly' }">Bulanan</span>
        <label class="relative inline-flex items-center cursor-pointer">
          <input type="checkbox" v-model="isYearly" class="sr-only peer">
          <div class="w-14 h-7 bg-gray-200 peer-focus:outline-none rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-0.5 after:left-[4px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-6 after:w-6 after:transition-all dark:border-gray-600 peer-checked:bg-secondary"></div>
        </label>
        <span :class="{ 'text-secondary font-bold': billingCycle === 'yearly' }">Tahunan</span>
        <span class="bg-yellow-200 text-yellow-800 text-xs font-semibold ml-2 px-2.5 py-0.5 rounded-full">Hemat 2 Bulan!</span>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 md:gap-8 lg:gap-12">
        <div
          v-for="pkg in packages"
          :key="pkg.id"
          :class="{ 'border-4 border-secondary shadow-2xl relative': pkg.name === 'Standard' }"
          class="bg-bg-base text-text-base p-4 md:p-8 rounded-xl shadow-lg flex flex-col transform hover:scale-105 transition duration-300 ease-in-out"
        >
          <span
            v-if="pkg.name === 'Standard'"
            class="absolute -top-6 left-1/2 -translate-x-1/2 bg-secondary text-white text-sm font-extrabold px-4 py-2 rounded-full uppercase"
          >
            Paling Populer
          </span>
          <h3 class="text-2xl font-bold mb-4 text-center">{{ pkg.package_name }}</h3>
          <p class="text-center text-text-muted mb-6">
            {{ pkg.package_name === 'Basic' ? 'Cocok untuk startup & bisnis kecil' : pkg.package_name === 'Standard' ? 'Ideal untuk bisnis berkembang' : 'Solusi khusus untuk perusahaan besar' }}
          </p>
          <div class="text-center mb-8">
            <span class="text-5xl font-extrabold text-secondary">
              Rp {{ billingCycle === 'monthly' ? pkg.price_monthly : pkg.price_yearly }}
            </span>
            <span class="text-xl text-text-muted">/{{ billingCycle === 'monthly' ? 'bulan' : 'tahun' }}</span>
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
          <BaseButton
            @click="selectPackage(pkg.id)"
            class="w-full mt-auto"
          >
            Mulai Coba Gratis
          </BaseButton>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import BaseButton from '../../ui/BaseButton.vue';

const props = defineProps(['packages']);
const isYearly = ref(false);
const router = useRouter();

const billingCycle = computed(() => {
  return isYearly.value ? 'yearly' : 'monthly';
});

const selectPackage = (packageId) => {
  console.log('Selecting package with ID:', packageId, 'for billing cycle:', billingCycle.value);
  router.push({
    name: 'RegisterCompany',
    params: { packageId: packageId },
    query: { billingCycle: billingCycle.value }
  });
};
</script>

<style scoped>
/* No custom scoped styles needed as Tailwind handles styling */
</style>
