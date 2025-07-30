<template>
  <div class="p-8 bg-bg-base">
    <Toast />
    <div v-if="isSubscribed" class="subscribed-view">
      <h1 class="text-3xl font-bold mb-2 text-text-base">Status Langganan Anda</h1>
      <p class="text-text-muted mb-8">Berikut adalah detail langganan aktif Anda.</p>

      <Card class="bg-bg-muted text-text-base shadow-lg mb-8">
        <template #title>
          <h3 class="text-2xl font-bold">{{ subscriptionDetails.package_name }}</h3>
        </template>
        <template #content>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <p><span class="font-semibold">Status:</span> <Tag :value="subscriptionDetails.subscription_status === 'trial' ? 'Percobaan' : 'Aktif'" :severity="subscriptionDetails.subscription_status === 'trial' ? 'warning' : 'success'" /></p>
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
            <Button @click="showPackageSelection" class="mt-4 p-button-outlined p-button-primary" label="Ubah Paket" />
          </div>
        </template>
      </Card>

      <Fieldset v-if="subscriptionDetails.subscription_status === 'trial'" legend="Masa Percobaan Gratis" class="bg-yellow-100 border-yellow-500 text-yellow-700">
          <div class="flex items-center">
            <i class="pi pi-info-circle text-yellow-500 text-2xl mr-4"></i>
            <div>
              <p class="font-bold">Sisa waktu Anda: {{ trialDaysRemaining }} hari.</p>
              <Button @click="switchToSubscription" class="mt-4 p-button-primary" label="Berlangganan Sekarang" />
            </div>
          </div>
      </Fieldset>
    </div>

    <div v-else class="not-subscribed-view">
      <h1 class="text-3xl font-bold mb-2 text-text-base">Pilih Paket Langganan</h1>
      <p class="text-text-muted mb-8">Masa percobaan Anda akan segera berakhir. Pilih paket untuk melanjutkan layanan.</p>

      <div class="flex justify-center items-center space-x-4 mb-12">
        <span :class="{ 'text-secondary font-bold': billingCycle === 'monthly' }">Bulanan</span>
        <ToggleSwitch v-model="isYearly" />
        <span :class="{ 'text-secondary font-bold': billingCycle === 'yearly' }">Tahunan</span>
        <Tag severity="warning" value="Hemat 2 Bulan!"></Tag>
      </div>

      <div v-if="isLoadingPackages" class="flex items-center justify-center py-4">
        <ProgressSpinner />
        <span class="ml-2">Memuat paket langganan...</span>
      </div>
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 md:gap-8">
        <Card
          v-for="pkg in packages"
          :key="pkg.id"
          class="bg-bg-muted text-text-base shadow-lg flex flex-col transform hover:scale-105 transition duration-300 ease-in-out"
        >
          <template #title><h3 class="text-2xl font-bold text-center">{{ pkg.package_name }}</h3></template>
          <template #content>
            <div class="flex flex-col flex-grow">
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
              <Button
                @click="selectPackage(pkg.id)"
                class="w-full mt-auto p-button-primary"
                icon="pi pi-shopping-cart"
                label="Pilih & Bayar"
              />
            </div>
          </template>
        </Card>
      </div>

      <div class="text-center mt-8">
        <p class="text-text-muted mb-4">Tidak menemukan paket yang sesuai? Kami dapat membuat penawaran khusus untuk Anda.</p>
        <Button @click="contactAdminForCustomPackage" class="p-button-primary" icon="pi pi-envelope" label="Minta Penawaran Kustom" />
      </div>
    </div>

    <Dialog v-model:visible="showSummaryModal" header="Ringkasan Pembayaran" :modal="true" class="w-full max-w-md">
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
          <div class="flex justify-between font-bold text-xl mt-4 pt-4 border-t border-gray-700">
            <span>Total Pembayaran:</span>
            <span class="text-secondary">{{ new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(selectedPackageDetails.finalPrice) }}</span>
          </div>
        </div>
        <p class="text-sm text-text-muted mb-6 text-center">Anda akan diarahkan ke halaman pembayaran Midtrans setelah melanjutkan.</p>
      </div>
      <template #footer>
        <Button @click="showSummaryModal = false" label="Batal" icon="pi pi-times" class="p-button-text"/>
        <Button @click="proceedToPayment" label="Lanjutkan ke Pembayaran" icon="pi pi-arrow-right" class="p-button-primary"/>
      </template>
    </Dialog>

    <Dialog v-model:visible="showContactModal" header="Minta Penawaran Kustom" :modal="true" class="w-full max-w-md">
      <form @submit.prevent="submitCustomPackageRequest" class="p-fluid">
        <div class="field mb-4">
          <label for="contactPhone" class="block text-text-muted text-sm font-bold mb-2">Nomor Telepon:</label>
          <InputText id="contactPhone" v-model="customPackageRequest.phone" type="tel" placeholder="Contoh: 08123456789" fluid />
        </div>
        <div class="field mb-4">
          <label for="message" class="block text-text-muted text-sm font-bold mb-2">Pesan/Kebutuhan (Opsional):</label>
          <Textarea id="message" v-model="customPackageRequest.message" rows="5" placeholder="Jelaskan kebutuhan spesifik Anda..." :autoResize="true" fluid/>
        </div>
        <div class="flex justify-end space-x-2 mt-6">
          <Button type="button" @click="showContactModal = false" label="Batal" class="p-button-outlined" />
          <Button type="submit" :loading="isSubmittingCustomRequest" label="Kirim Permintaan" class="p-button-primary" />
        </div>
      </form>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../../stores/auth';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import Button from 'primevue/button';
import Dialog from 'primevue/dialog';
import ToggleSwitch from 'primevue/toggleswitch';
import Textarea from 'primevue/textarea';
import InputText from 'primevue/inputtext';
import Card from 'primevue/card';
import Tag from 'primevue/tag';
import Fieldset from 'primevue/fieldset';
import ProgressSpinner from 'primevue/progressspinner';
import Toast from 'primevue/toast';

const packages = ref([]);
const router = useRouter();
const authStore = useAuthStore();
const toast = useToast();
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
    package_name: authStore.companyName,
    subscription_status: authStore.subscriptionStatus,
    trial_end_date: authStore.trialEndDate,
    billing_cycle: authStore.billingCycle,
    subscription_package_id: authStore.subscriptionPackageId,
    max_employees: 0,
    features: '',
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
    toast.add({ severity: 'error', summary: 'Error', detail: 'Terjadi kesalahan saat memuat detail langganan Anda.', life: 3000 });
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
    toast.add({ severity: 'error', summary: 'Gagal Memuat', detail: 'Tidak dapat memuat paket langganan.', life: 3000 });
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
    toast.add({ severity: 'error', summary: 'Error', detail: 'Paket yang dipilih tidak ditemukan.', life: 3000 });
  }
};

const proceedToPayment = () => {
    const companyId = authStore.companyId;
    if (!companyId) {
      toast.add({ severity: 'error', summary: 'Error', detail: 'ID Perusahaan tidak ditemukan.', life: 3000 });
      return;
    }
    if (!selectedPackageDetails.value) {
      toast.add({ severity: 'error', summary: 'Error', detail: 'Tidak ada paket yang dipilih untuk pembayaran.', life: 3000 });
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
  customPackageRequest.value = { phone: '', message: '' };
};

const submitCustomPackageRequest = async () => {
  isSubmittingCustomRequest.value = true;
  try {
    const response = await axios.post('/api/custom-package-requests', customPackageRequest.value);
    if (response.data && response.data.status === 'success') {
      toast.add({ severity: 'success', summary: 'Berhasil', detail: 'Permintaan Anda telah terkirim. Admin kami akan segera menghubungi Anda!', life: 4000 });
      showContactModal.value = false;
    } else {
      toast.add({ severity: 'error', summary: 'Gagal', detail: response.data?.message || 'Gagal mengirim permintaan. Silakan coba lagi.', life: 3000 });
    }
  } catch (error) {
    console.error('Error submitting custom package request:', error);
    toast.add({ severity: 'error', summary: 'Error', detail: 'Terjadi kesalahan saat mengirim permintaan.', life: 3000 });
  } finally {
    isSubmittingCustomRequest.value = false;
  }
};

onMounted(async () => {
  isSubscribed.value = authStore.subscriptionStatus === 'active' || authStore.subscriptionStatus === 'trial';
  await fetchSubscriptionPackages();
});
</script>