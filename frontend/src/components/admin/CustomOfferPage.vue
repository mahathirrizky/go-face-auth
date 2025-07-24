<template>
  <div class="min-h-screen flex items-center justify-center bg-bg-base py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8 p-10 bg-bg-muted rounded-lg shadow-xl">
      <div v-if="isLoading" class="text-center text-text-muted">
        Memuat penawaran kustom...
      </div>
      <div v-else-if="error" class="text-center text-red-500">
        {{ error }}
        <p v-if="error.includes('Unauthorized') || error.includes('Forbidden')" class="mt-4">
          Silakan <router-link to="/login/admin" class="text-blue-500 hover:underline">login</router-link> sebagai admin perusahaan untuk melihat penawaran ini.
        </p>
      </div>
      <div v-else-if="offer">
        <h2 class="mt-6 text-center text-3xl font-extrabold text-text-base">
          Penawaran Paket Kustom untuk {{ offer.company_name }}
        </h2>
        <p class="mt-2 text-center text-sm text-text-muted">
          Tinjau detail penawaran khusus Anda di bawah ini.
        </p>

        <div class="mt-8 space-y-6">
          <div class="bg-bg-base p-6 rounded-md shadow-inner">
            <div class="flex justify-between mb-2">
              <span class="font-medium text-text-muted">Nama Paket:</span>
              <span class="font-semibold text-text-base">{{ offer.package_name }}</span>
            </div>
            <div class="flex justify-between mb-2">
              <span class="font-medium text-text-muted">Jumlah Karyawan:</span>
              <span class="font-semibold text-text-base">{{ offer.max_employees }}</span>
            </div>
            <div class="flex justify-between mb-2">
              <span class="font-medium text-text-muted">Jumlah Lokasi:</span>
              <span class="font-semibold text-text-base">{{ offer.max_locations }}</span>
            </div>
            <div class="flex justify-between mb-2">
              <span class="font-medium text-text-muted">Jumlah Shift:</span>
              <span class="font-semibold text-text-base">{{ offer.max_shifts }}</span>
            </div>
            <div class="mb-2">
              <span class="font-medium text-text-muted">Fitur:</span>
              <ul class="list-disc list-inside text-text-base ml-4">
                <li v-for="(feature, index) in offer.features.split(',')" :key="index">{{ feature.trim() }}</li>
              </ul>
            </div>
            <div class="flex justify-between mb-2">
              <span class="font-medium text-text-muted">Siklus Penagihan:</span>
              <span class="font-semibold text-text-base">{{ offer.billing_cycle === 'monthly' ? 'Bulanan' : 'Tahunan' }}</span>
            </div>
            <div class="flex justify-between border-t border-gray-700 pt-4 mt-4">
              <span class="text-xl font-bold text-text-base">Harga Akhir:</span>
              <span class="text-2xl font-extrabold text-secondary">{{ formatCurrency(offer.final_price) }}</span>
            </div>
          </div>

          <BaseButton
            @click="proceedToPayment"
            :disabled="offer.status !== 'pending'"
            class="w-full btn-primary py-2 px-4"
          >
            <i class="pi pi-shopping-cart"></i>
            {{ offer.status === 'pending' ? 'Bayar Sekarang' : 'Penawaran Tidak Tersedia' }}
          </BaseButton>
          <p v-if="offer.status !== 'pending'" class="text-center text-sm text-red-500 mt-2">
            Penawaran ini {{ offer.status === 'used' ? 'sudah digunakan' : offer.status === 'expired' ? 'sudah kadaluarsa' : 'tidak tersedia' }}.
          </p>
        </div>
      </div>
      <div v-else class="text-center text-text-muted">
        Penawaran tidak ditemukan atau tidak valid.
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import axios from 'axios';
import BaseButton from '../ui/BaseButton.vue';
import { useAuthStore } from '../../stores/auth';

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();
const offer = ref(null);
const isLoading = ref(true);
const error = ref(null);

const fetchCustomOffer = async () => {
  try {
    const token = route.params.token;
    if (!token) {
      error.value = 'Token penawaran tidak ditemukan.';
      return;
    }

    // Check if user is authenticated and companyID matches
    if (!authStore.isAuthenticated || !authStore.companyId) {
      error.value = 'Unauthorized: Anda perlu login sebagai admin perusahaan.';
      isLoading.value = false;
      return;
    }

    const response = await axios.get(`/api/offer/${token}`);
    if (response.data && response.data.status === 'success') {
      // Additional client-side check to ensure the offer belongs to the logged-in company
      if (response.data.data.company_id !== authStore.companyId) {
        error.value = 'Forbidden: Penawaran ini bukan untuk perusahaan Anda.';
        offer.value = null; // Clear offer data
      } else {
        offer.value = response.data.data;
      }
    } else {
      error.value = response.data?.message || 'Gagal mengambil detail penawaran.';
    }
  } catch (err) {
    console.error('Error fetching custom offer:', err);
    if (err.response && (err.response.status === 401 || err.response.status === 403)) {
      error.value = 'Unauthorized: Anda perlu login sebagai admin perusahaan untuk melihat penawaran ini.';
    } else {
      error.value = err.response?.data?.message || 'Terjadi kesalahan saat mengambil penawaran.';
    }
  } finally {
    isLoading.value = false;
  }
};

const formatCurrency = (value) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(value);
};

const proceedToPayment = () => {
  if (offer.value && offer.value.status === 'pending') {
    router.push({
      name: 'PaymentPage',
      query: {
        customOfferToken: offer.value.token,
        companyId: offer.value.company_id
      }
    });
  } else {
    alert('Penawaran ini tidak dapat diproses.');
  }
};

onMounted(() => {
  fetchCustomOffer();
});
</script>
