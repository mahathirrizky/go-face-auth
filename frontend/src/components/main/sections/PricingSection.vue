<template>
  <section id="pricing" class="py-16 md:py-24 bg-bg-muted text-text-base">
    <div class="container mx-auto px-4">
      <h2 class="text-3xl md:text-4xl font-extrabold text-center mb-6">
        Pilih Paket yang Sesuai untuk Bisnis Anda
      </h2>
      
      <!-- Billing Cycle Toggle -->
      <div class="flex justify-center items-center space-x-4 mb-12">
        <span :class="{ 'text-secondary font-bold': billingCycle === 'monthly' }">Bulanan</span>
        <ToggleSwitch v-model="isYearly" />
        <span :class="{ 'text-secondary font-bold': billingCycle === 'yearly' }">Tahunan</span>
        <Tag severity="warning" value="Hemat 2 Bulan!" class="ml-2" />
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 md:gap-8 lg:gap-12">
        <Card
          v-for="pkg in packages"
          :key="pkg.id"
          :class="{ 'border-4 border-secondary shadow-2xl relative': pkg.name === 'Standard' }"
          class="bg-bg-base text-text-base shadow-lg flex flex-col transform hover:scale-105 transition duration-300 ease-in-out"
        >
          <template #header v-if="pkg.name === 'Standard'">
            <div class="absolute -top-6 left-1/2 -translate-x-1/2 bg-secondary text-white text-sm font-extrabold px-4 py-2 rounded-full uppercase">
              Paling Populer
            </div>
          </template>
          <template #title>
            <h3 class="text-2xl font-bold text-center">{{ pkg.package_name }}</h3>
          </template>
          <template #subtitle>
            <p class="text-center text-text-muted">
              {{ pkg.package_name === 'Basic' ? 'Cocok untuk startup & bisnis kecil' : pkg.package_name === 'Standard' ? 'Ideal untuk bisnis berkembang' : 'Solusi khusus untuk perusahaan besar' }}
            </p>
          </template>
          <template #content>
            <div class="text-center mb-8">
              <span class="text-5xl font-extrabold text-secondary">
                Rp {{ new Intl.NumberFormat('id-ID').format(billingCycle === 'monthly' ? pkg.price_monthly : pkg.price_yearly) }}
              </span>
              <span class="text-xl text-text-muted">/{{ billingCycle === 'monthly' ? 'bulan' : 'tahun' }}</span>
            </div>
            <ul class="text-left space-y-3 mb-8 flex-grow">
              <li class="flex items-center">
                <i class="pi pi-check-circle text-green-500 mr-2"></i>
                Hingga {{ pkg.max_employees }} Karyawan
              </li>
              <li class="flex items-center" v-for="(feature, index) in pkg.features.split(',')" :key="index">
                <i class="pi pi-check-circle text-green-500 mr-2"></i>
                {{ feature.trim() }}
              </li>
            </ul>
          </template>
          <template #footer>
            <Button
              @click="selectPackage(pkg.id)"
              class="w-full mt-auto p-button-primary"
              icon="pi pi-play"
              label="Mulai Coba Gratis"
            />
          </template>
        </Card>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import Button from 'primevue/button';
import ToggleSwitch from 'primevue/toggleswitch';
import Card from 'primevue/card';
import Tag from 'primevue/tag';

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