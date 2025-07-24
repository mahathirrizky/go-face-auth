<template>
  <div class="p-8 bg-bg-base">
    <h1 class="text-3xl font-bold mb-2 text-text-base">Pilih Paket Langganan</h1>
    <p class="text-text-muted mb-8">Masa percobaan Anda akan segera berakhir. Pilih paket untuk melanjutkan layanan.</p>

    <!-- Billing Cycle Toggle -->
    <div class="flex justify-center items-center space-x-4 mb-12">
      <span :class="{ 'text-secondary font-bold': billingCycle === 'monthly' }">Bulanan</span>
      <ToggleSwitch v-model="isYearly" />
      <span :class="{ 'text-secondary font-bold': billingCycle === 'yearly' }">Tahunan</span>
      <span class="bg-yellow-200 text-yellow-800 text-xs font-semibold ml-2 px-2.5 py-0.5 rounded-full">Hemat 2 Bulan!</span>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
      <div
        v-for="pkg in packages"
        :key="pkg.id"
        class="bg-bg-muted text-text-base p-8 rounded-xl shadow-lg flex flex-col transform hover:scale-105 transition duration-300 ease-in-out"
      >
        <h3 class="text-2xl font-bold mb-4 text-center">{{ pkg.package_name }}</h3>
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
          class="w-full mt-auto btn-primary"
        >
          <i class="pi pi-shopping-cart"></i> Pilih Paket & Bayar
        </BaseButton>
      </div>
    </div>

    <!-- Custom Package Section -->
    <div class="mt-12 text-center">
      <h2 class="text-2xl font-bold text-text-base mb-4">Butuh Solusi yang Disesuaikan?</h2>
      <p class="text-text-muted mb-6">Jika paket yang tersedia tidak sesuai dengan kebutuhan unik perusahaan Anda, kami dapat membuat penawaran khusus.</p>
      <BaseButton @click="contactAdminForCustomPackage" class="btn-primary">
        <i class="pi pi-envelope"></i> Hubungi Admin untuk Paket Kustom
      </BaseButton>
    </div>

    <!-- Payment Summary Modal -->
    <BaseModal :isOpen="showSummaryModal" @close="showSummaryModal = false" title="Ringkasan Pembayaran" maxWidth="md">
      <div v-if="selectedPackageDetails">
        <div class="mb-4 border-b border-gray-700 pb-4">
          <p class="text-lg font-semibold">Paket yang Dipilih:</p>
          <p class="text-xl font-bold text-secondary">{{ selectedPackageDetails.package_name }}</p>
          <p class="text-sm text-text-muted">Siklus Penagihan: {{ selectedPackageDetails.displayBillingCycle }}</p>
        </div>
        <div class="mb-6">
          <p class="text-lg font-semibold">Detail Pembayaran:</p>
          <div class="flex justify-between mb-2">
            <span>Harga Paket:</span>
            <span>{{ new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(selectedPackageDetails.finalPrice) }}</span>
          </div>
          <!-- Add more details like tax, discount if applicable -->
          <div class="flex justify-between font-bold text-xl mt-4 pt-4 border-t border-gray-700">
            <span>Total Pembayaran:</span>
            <span class="text-secondary">{{ new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(selectedPackageDetails.finalPrice) }}</span>
          </div>
        </div>
        <p class="text-sm text-text-muted mb-6 text-center">Anda akan diarahkan ke halaman pembayaran Midtrans setelah melanjutkan.</p>
      </div>
      <template #footer>
        <BaseButton @click="showSummaryModal = false" class="btn-outline mr-2">
          <i class="pi pi-times"></i> Batal
        </BaseButton>
        <BaseButton @click="proceedToPayment" class="btn-primary">
          <i class="pi pi-arrow-right"></i> Lanjutkan ke Pembayaran
        </BaseButton>
      </template>
    </BaseModal>

    <!-- Contact Modal -->
    <BaseModal :isOpen="showContactModal" @close="showContactModal = false" title="Minta Penawaran Kustom" maxWidth="md">
      <form @submit.prevent="submitCustomPackageRequest">
        <div class="mb-4">
          <label for="contactPhone" class="block text-text-muted text-sm font-bold mb-2">Nomor Telepon:</label>
          <BaseInput
            id="contactPhone"
            v-model="customPackageRequest.phone"
            type="tel"
          />
        </div>
        <div class="mb-4">
          <label for="message" class="block text-text-muted text-sm font-bold mb-2">Pesan/Kebutuhan (Opsional):</label>
          <Textarea
            id="message"
            v-model="customPackageRequest.message"
            rows="5"
            class="w-full"
            placeholder="Jelaskan kebutuhan spesifik Anda..."
          />
        </div>
        <div class="flex justify-end space-x-4 mt-6">
          <BaseButton type="button" @click="showContactModal = false" class="btn-outline-primary">
            Batal
          </BaseButton>
          <BaseButton type="submit" :loading="isSubmittingCustomRequest">
            Kirim Permintaan
          </BaseButton>
        </div>
      </form>
    </BaseModal>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../../stores/auth';
import axios from 'axios';
import BaseButton from '../ui/BaseButton.vue';
import BaseModal from '../ui/BaseModal.vue';
import ToggleSwitch from 'primevue/toggleswitch';
import Textarea from 'primevue/textarea';
import BaseInput from '../ui/BaseInput.vue';
const packages = ref([]);
const router = useRouter();
const authStore = useAuthStore();
const isYearly = ref(false);
const showSummaryModal = ref(false);
const selectedPackageDetails = ref(null);
const showContactModal = ref(false);
const isSubmittingCustomRequest = ref(false); // New ref for loading state

const customPackageRequest = ref({
  phone: '',
  message: '',
});

const billingCycle = computed(() => {
  return isYearly.value ? 'yearly' : 'monthly';
});

const fetchSubscriptionPackages = async () => {
  try {
    const response = await axios.get('/api/subscription-packages');
    packages.value = response.data.data;
  } catch (error) {
    console.error('Error fetching subscription packages:', error);
  }
};

const selectPackage = (packageId) => {
  const selectedPkg = packages.value.find(pkg => pkg.id === packageId);
  if (selectedPkg) {
    const price = billingCycle.value === 'monthly' ? selectedPkg.price_monthly : selectedPkg.price_yearly;
    selectedPackageDetails.value = {
      ...selectedPkg,
      finalPrice: price,
      displayBillingCycle: billingCycle.value === 'monthly' ? 'Bulanan' : 'Tahunan',
    };
    showSummaryModal.value = true;
  } else {
    console.error('Selected package not found.');
  }
};

const proceedToPayment = () => {
    const companyId = authStore.companyId;
    console.log('SubscriptionPage: companyId before navigation:', companyId);
    if (!companyId) {
      console.error('Company ID not found in store.');
      return;
    }
    if (!selectedPackageDetails.value) {
      console.error('No package selected for payment.');
      return;
    }

    router.push({
      name: 'PaymentPage',
      params: { companyId: String(companyId) },
      query: {
        packageId: selectedPackageDetails.value.id,
        billingCycle: billingCycle.value
      }
    });
    showSummaryModal.value = false;
  };

const contactAdminForCustomPackage = () => {
  showContactModal.value = true;
  // Reset form fields when opening the modal
  customPackageRequest.value = {
    message: '',
  };
};

const submitCustomPackageRequest = async () => {
  isSubmittingCustomRequest.value = true;
  try {
    // Replace with your actual backend endpoint for custom package requests
    const response = await axios.post('/api/custom-package-requests', customPackageRequest.value);
    if (response.data && response.data.status === 'success') {
      alert('Permintaan Anda telah terkirim. Admin kami akan segera menghubungi Anda!');
      showContactModal.value = false;
    } else {
      alert('Gagal mengirim permintaan. Silakan coba lagi.');
    }
  } catch (error) {
    console.error('Error submitting custom package request:', error);
    alert('Terjadi kesalahan saat mengirim permintaan. Silakan coba lagi nanti.');
  } finally {
    isSubmittingCustomRequest.value = false;
  }
};

onMounted(() => {
  fetchSubscriptionPackages();
});
</script>

