<template>
  <div class="min-h-screen flex items-center justify-center bg-bg-base py-12 px-4 sm:px-6 lg:px-8">
    <Card class="w-full max-w-md shadow-xl">
      <template #content>
        <div v-if="isLoading" class="text-center text-text-muted">
          <ProgressSpinner />
          <p class="mt-2">Memuat penawaran kustom...</p>
        </div>
        <div v-else-if="error">
          <Message severity="error" :closable="false">{{ error }}</Message>
          <p v-if="error.includes('Unauthorized') || error.includes('Forbidden')" class="mt-4 text-center">
            Silakan <router-link to="/login/admin" class="text-blue-500 hover:underline">login</router-link> sebagai admin perusahaan untuk melihat penawaran ini.
          </p>
        </div>
        <div v-else-if="offer">
          <h2 class="text-center text-2xl font-bold text-text-base mb-2">
            Penawaran Paket Kustom
          </h2>
          <p class="text-center text-lg font-semibold text-secondary mb-6">{{ offer.company_name }}</p>

          <div class="space-y-4">
            <div class="bg-bg-base p-4 rounded-md shadow-inner">
              <div class="flex justify-between mb-2">
                <span class="font-medium text-text-muted">Nama Paket:</span>
                <span class="font-semibold text-text-base">{{ offer.package_name }}</span>
              </div>
              <div class="flex justify-between mb-2">
                <span class="font-medium text-text-muted">Jml. Karyawan:</span>
                <span class="font-semibold text-text-base">{{ offer.max_employees }}</span>
              </div>
              <div class="flex justify-between mb-2">
                <span class="font-medium text-text-muted">Jml. Lokasi:</span>
                <span class="font-semibold text-text-base">{{ offer.max_locations }}</span>
              </div>
              <div class="flex justify-between mb-2">
                <span class="font-medium text-text-muted">Jml. Shift:</span>
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

            <Button
              @click="proceedToPayment"
              :disabled="offer.status !== 'pending'"
              class="w-full py-3"
              icon="pi pi-shopping-cart"
              :label="offer.status === 'pending' ? 'Bayar Sekarang' : 'Penawaran Tidak Tersedia'"
            />
            <Message v-if="offer.status !== 'pending'" severity="warn" :closable="false">
              Penawaran ini {{ offer.status === 'used' ? 'sudah digunakan' : offer.status === 'expired' ? 'sudah kadaluarsa' : 'tidak tersedia' }}.
            </Message>
          </div>
        </div>
        <div v-else>
            <Message severity="warn" :closable="false">Penawaran tidak ditemukan atau tidak valid.</Message>
        </div>
      </template>
    </Card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import axios from 'axios';
import Button from 'primevue/button';
import Card from 'primevue/card';
import ProgressSpinner from 'primevue/progressspinner';
import Message from 'primevue/message';
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

    if (!authStore.isAuthenticated || !authStore.companyId) {
      error.value = 'Unauthorized: Anda perlu login sebagai admin perusahaan.';
      isLoading.value = false;
      return;
    }

    const response = await axios.get(`/api/offer/${token}`);
    if (response.data && response.data.status === 'success') {
      if (response.data.data.company_id !== authStore.companyId) {
        error.value = 'Forbidden: Penawaran ini bukan untuk perusahaan Anda.';
        offer.value = null;
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