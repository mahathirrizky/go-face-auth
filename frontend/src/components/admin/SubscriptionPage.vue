<template>
  <div class="p-8 bg-bg-base">
    <div v-if="isSubscribed" class="subscribed-view">
      <h1 class="text-3xl font-bold mb-2 text-text-base">Status Langganan Anda</h1>
      <p class="text-text-muted mb-8">Berikut adalah detail langganan aktif Anda.</p>

      <div class="bg-bg-muted text-text-base p-8 rounded-xl shadow-lg mb-8">
        <h3 class="text-2xl font-bold mb-4">{{ subscriptionDetails.package_name }}</h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <p><span class="font-semibold">Status:</span> <span class="text-green-500 font-bold">{{ subscriptionDetails.subscription_status === 'trial' ? 'Percobaan' : 'Aktif' }}</span></p>
            <p v-if="subscriptionDetails.trial_end_date"><span class="font-semibold">Masa Percobaan Berakhir:</span> {{ new Date(subscriptionDetails.trial_end_date).toLocaleDateString('id-ID') }}</p>
            <p v-if="subscriptionDetails.subscription_end_date"><span class="font-semibold">Berakhir Pada:</span> {{ new Date(subscriptionDetails.subscription_end_date).toLocaleDateString('id-ID') }}</p>
            <p><span class="font-semibold">Siklus Penagihan:</span> {{ subscriptionDetails.billing_cycle === 'monthly' ? 'Bulanan' : 'Tahunan' }}</p>
          </div>
          <div>
            <p><span class="font-semibold">Maksimal Karyawan:</span> {{ subscriptionDetails.max_employees }}</p>
            <p class="font-semibold mb-2">Fitur:</p>
            <ul v-if="subscriptionDetails && subscriptionDetails.features" class="list-disc list-inside">
              <li v-for="(feature, index) in subscriptionDetails.features.split(',')" :key="index">{{ feature.trim() }}</li>
            </ul>
          </div>
        </div>
        <div class="mt-6 text-center">
          <p class="text-text-muted">Terima kasih telah berlangganan layanan kami!</p>
          <BaseButton @click="showPackageSelection" class="mt-4 btn-outline-primary">
            Ubah Paket
          </BaseButton>
        </div>
      </div>

      <div v-if="subscriptionDetails.subscription_status === 'trial'" class="bg-yellow-100 border-l-4 border-yellow-500 text-yellow-700 p-4 mt-8 rounded-lg shadow-md">
        <div class="flex">
          <div class="py-1"><svg class="fill-current h-6 w-6 text-yellow-500 mr-4" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20"><path d="M2.93 17.07A10 10 0 1 1 17.07 2.93 10 10 0 0 1 2.93 17.07zM9 5v6h2V5H9zm0 8h2v2H9v-2z"/></svg></div>
          <div>
            <p class="font-bold">Anda dalam masa coba gratis.</p>
            <p class="text-sm">Sisa waktu Anda: {{ trialDaysRemaining }} hari.</p>
            <BaseButton @click="switchToSubscription" class="mt-4 btn-primary">Berlangganan Sekarang</BaseButton>
          </div>
        </div>
      </div>


      

    </div>

    <div v-else class="not-subscribed-view">
      <h1 class="text-3xl font-bold mb-2 text-text-base">Pilih Paket Langganan</h1>
      <p class="text-text-muted mb-8">Masa percobaan Anda akan segera berakhir. Pilih paket untuk melanjutkan layanan.</p>

      <!-- Billing Cycle Toggle -->
      <div class="flex justify-center items-center space-x-4 mb-12">
        <span :class="{ 'text-secondary font-bold': billingCycle === 'monthly' }">Bulanan</span>
        <ToggleSwitch v-model="isYearly" />
        <span :class="{ 'text-secondary font-bold': billingCycle === 'yearly' }">Tahunan</span>
        <span class="bg-yellow-200 text-yellow-800 text-xs font-semibold ml-2 px-2.5 py-0.5 rounded-full">Hemat 2 Bulan!</span>
      </div>

      <div v-if="isLoadingPackages" class="flex items-center justify-center py-4">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
        <span class="ml-2">Memuat paket langganan...</span>
      </div>
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 md:gap-8">
        <div
          v-for="pkg in packages"
          :key="pkg.id"
          class="bg-bg-muted text-text-base p-4 md:p-8 rounded-xl shadow-lg flex flex-col transform hover:scale-105 transition duration-300 ease-in-out"
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
            <i class="pi pi-shopping-cart"></i> Pilih & Bayar
          </BaseButton>
        </div>
      </div>

      <div class="text-center mt-8">
        <p class="text-text-muted mb-4">Tidak menemukan paket yang sesuai? Kami dapat membuat penawaran khusus untuk Anda.</p>
        <BaseButton @click="contactAdminForCustomPackage" class="btn-primary">
          <i class="pi pi-envelope"></i> Minta Penawaran Kustom
        </BaseButton>
      </div>

      
      
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
const isSubmittingCustomRequest = ref(false);
const isLoadingPackages = ref(false);
const customPackageRequest = ref({
  phone: '',
  message: '',
});

const isSubscribed = ref(false);
const subscriptionDetails = computed(() => {
  const details = {
    package_name: authStore.companyName, // Default, will be overridden if package found
    subscription_status: authStore.subscriptionStatus,
    trial_end_date: authStore.trialEndDate,
    billing_cycle: authStore.billingCycle,
    subscription_package_id: authStore.subscriptionPackageId,
    max_employees: 0, // Default
    features: '', // Default
  };

  if (authStore.subscriptionPackageId && packages.value.length > 0) {
    const subscribedPackage = packages.value.find(
      (pkg) => pkg.id === authStore.subscriptionPackageId
    );
    if (subscribedPackage) {
      details.package_name = subscribedPackage.package_name;
      details.max_employees = subscribedPackage.max_employees;
      details.features = subscribedPackage.features;
    }
  }
  return details;
});

const billingCycle = computed(() => {
  return isYearly.value ? 'yearly' : 'monthly';
});

const trialDaysRemaining = computed(() => {
  if (subscriptionDetails.value && subscriptionDetails.value.trial_end_date) {
    const endDate = new Date(subscriptionDetails.value.trial_end_date);
    const now = new Date();
    const diffTime = Math.abs(endDate - now);
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    return diffDays;
  }
  return 0;
});

const switchToSubscription = () => {
  if (authStore.subscriptionPackageId) {
    isYearly.value = authStore.billingCycle === 'yearly';
    selectPackage(authStore.subscriptionPackageId);
  } else {
    console.error('Could not show payment summary: Subscription package ID not found in auth store.');
    alert('Terjadi kesalahan saat memuat detail langganan Anda. Silakan coba lagi.');
  }
};

const showPackageSelection = () => {
  isSubscribed.value = false;
};

const fetchSubscriptionPackages = async () => {
  isLoadingPackages.value = true;
  try {
    const response = await axios.get('/api/subscription-packages');
    packages.value = response.data.data;
  } catch (error) {
    console.error('Error fetching subscription packages:', error);
  } finally {
    isLoadingPackages.value = false;
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

onMounted(async () => {
  // Initialize isSubscribed based on authStore status
  isSubscribed.value = authStore.subscriptionStatus === 'active' || authStore.subscriptionStatus === 'trial';
  
  // Fetch packages
  await fetchSubscriptionPackages();
});
</script>

