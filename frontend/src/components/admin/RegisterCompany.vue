<template>
  <div class="min-h-screen flex items-center justify-center bg-bg-base">
    <div class="bg-bg-muted p-8 rounded-lg shadow-md w-full max-w-md">
      <h2 class="text-2xl font-bold text-center mb-6 text-text-base">Daftar Perusahaan Baru</h2>
      <form @submit.prevent="registerCompany">
        <BaseInput
          id="companyName"
          label="Nama Perusahaan:"
          v-model="form.company_name"
          required
        />
        <BaseInput
          id="companyAddress"
          label="Alamat Perusahaan:"
          v-model="form.company_address"
        />
        <BaseInput
          id="adminEmail"
          label="Email Admin:"
          v-model="form.admin_email"
          type="email"
          required
        />
        <PasswordInput
          id="adminPassword"
          label="Password Admin:"
          v-model="form.admin_password"
          required
        />
        <BaseInput
          id="subscriptionPackage"
          label="Paket Langganan:"
          v-model="selectedPackageName"
          readonly
          class="cursor-not-allowed"
        />
        <BaseButton :fullWidth="true" class="mt-6">
          Daftar & Mulai Coba Gratis
        </BaseButton>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue';
import { useToast } from 'vue-toastification';
import { useRoute } from 'vue-router';
import axios from 'axios';
import { getBaseDomain } from '../../utils/subdomain';
import BaseInput from '../ui/BaseInput.vue';
import BaseButton from '../ui/BaseButton.vue';
import PasswordInput from '../ui/PasswordInput.vue';

const props = defineProps(['packageId']);
const toast = useToast();
const route = useRoute();

const form = ref({
  company_name: '',
  company_address: '',
  admin_email: '',
  admin_password: '',
  subscription_package_id: parseInt(props.packageId),
  billing_cycle: route.query.billingCycle || 'monthly',
});
const selectedPackageName = ref('');
const subscriptionPackages = ref([]);

const fetchSubscriptionPackages = async () => {
  try {
    const response = await axios.get('/api/subscription-packages');
    subscriptionPackages.value = response.data.data;
  } catch (error) {
    console.error('Error fetching subscription packages:', error);
  }
};

watch([subscriptionPackages, () => form.value.subscription_package_id], ([newPackages, newPackageId]) => {
  if (newPackages.length > 0 && newPackageId) {
    const selectedPackage = newPackages.find(pkg => pkg.id === newPackageId);
    if (selectedPackage) {
      selectedPackageName.value = selectedPackage.package_name;
    } else {
      selectedPackageName.value = 'Paket tidak ditemukan';
    }
  }
}, { immediate: true });

onMounted(() => {
  fetchSubscriptionPackages();
});

const registerCompany = async () => {
  try {
    const response = await axios.post('/api/register-company', form.value);
    toast.success(response.data.message);
    setTimeout(() => {
      const baseDomain = getBaseDomain();
      const adminLoginURL = `${window.location.protocol}//admin.${baseDomain}${window.location.port ? ':' + window.location.port : ''}/`;
      window.location.href = adminLoginURL;
    }, 2000);
  } catch (error) {
    const errorMessage = error.response?.data?.message || error.message;
    toast.error('Registration failed: ' + errorMessage);
  }
};
</script>
