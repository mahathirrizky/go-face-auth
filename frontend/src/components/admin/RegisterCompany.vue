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
        <div class="mb-4">
          <label for="adminPassword" class="block text-text-muted text-sm font-bold mb-2">Password Admin:</label>
          <Password
            id="adminPassword"
            v-model="form.admin_password"
            placeholder="Minimal 8 karakter"
            :required="true"
            toggleMask
            :feedback="true"
            class="w-full"
          >
            <template #header>
                <h6>Atur Kata Sandi</h6>
            </template>
            <template #footer>
                <Divider />
                <p class="mt-2">Saran:</p>
                <ul class="pl-2 ml-2 mt-0" style="line-height: 1.5">
                    <li>Minimal satu huruf kecil</li>
                    <li>Minimal satu huruf besar</li>
                    <li>Minimal satu angka</li>
                    <li>Minimal 8 karakter</li>
                </ul>
            </template>
          </Password>
        </div>
        <BaseInput
          id="subscriptionPackage"
          label="Paket Langganan:"
          v-model="selectedPackageName"
          readonly
          class="cursor-not-allowed"
        />
        <BaseButton :fullWidth="true" class="mt-6 btn-primary">
          <i class="pi pi-check"></i> Daftar & Mulai Coba Gratis
        </BaseButton>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue';
import { useToast } from 'primevue/usetoast';
import { useRoute } from 'vue-router';
import axios from 'axios';
import { getBaseDomain } from '../../utils/subdomain';
import BaseInput from '../ui/BaseInput.vue';
import BaseButton from '../ui/BaseButton.vue';
import Password from 'primevue/password';
import Divider from 'primevue/divider';

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
    toast.add({ severity: 'success', summary: 'Success', detail: response.data.message, life: 3000 });
    setTimeout(() => {
      const baseDomain = getBaseDomain();
      const adminLoginURL = `${window.location.protocol}//admin.${baseDomain}${window.location.port ? ':' + window.location.port : ''}/`;
      window.location.href = adminLoginURL;
    }, 2000);
  } catch (error) {
    const errorMessage = error.response?.data?.message || error.message;
    toast.add({ severity: 'error', summary: 'Error', detail: 'Registration failed: ' + errorMessage, life: 3000 });
  }
};
</script>
