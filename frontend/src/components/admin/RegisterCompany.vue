<template>
  <div class="min-h-screen flex items-center justify-center bg-bg-base">
    <div class="bg-bg-muted p-8 rounded-lg shadow-md w-full max-w-md">
      <h2 class="text-2xl font-bold text-center mb-6 text-text-base">Daftar Perusahaan Baru</h2>
      <form @submit.prevent="registerCompany">
        <div class="mb-4">
          <label for="companyName" class="block text-text-muted text-sm font-bold mb-2">Nama Perusahaan:</label>
          <input
            type="text"
            id="companyName"
            v-model="form.company_name"
            class="shadow appearance-none border rounded w-full py-2 px-3 text-text-base bg-bg-base leading-tight focus:outline-none focus:shadow-outline"
            required
          />
        </div>
        <div class="mb-4">
          <label for="companyAddress" class="block text-text-muted text-sm font-bold mb-2">Alamat Perusahaan:</label>
          <input
            type="text"
            id="companyAddress"
            v-model="form.company_address"
            class="shadow appearance-none border rounded w-full py-2 px-3 text-text-base bg-bg-base leading-tight focus:outline-none focus:shadow-outline"
          />
        </div>
        <div class="mb-4">
          <label for="adminEmail" class="block text-text-muted text-sm font-bold mb-2">Email Admin:</label>
          <input
            type="email"
            id="adminEmail"
            v-model="form.admin_email"
            class="shadow appearance-none border rounded w-full py-2 px-3 text-text-base bg-bg-base leading-tight focus:outline-none focus:shadow-outline"
            required
          />
        </div>
        <div class="mb-6">
          <label for="adminPassword" class="block text-text-muted text-sm font-bold mb-2">Password Admin:</label>
          <div class="relative">
            <input
              :type="passwordFieldType"
              id="adminPassword"
              v-model="form.admin_password"
              class="shadow appearance-none border rounded w-full py-2 px-3 text-text-base bg-bg-base leading-tight focus:outline-none focus:shadow-outline pr-10"
              required
            />
            <button
              type="button"
              @click="togglePasswordVisibility"
              class="absolute inset-y-0 right-0 pr-3 flex items-center focus:outline-none"
            >
              <svg v-if="passwordFieldType === 'password'" class="h-6 w-6 text-text-muted" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
              </svg>
              <svg v-else class="h-6 w-6 text-text-muted" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.542-7 .985-3.14 3.29-5.578 6.16-7.037m6.715 6.715A3 3 0 0112 15a3 3 0 01-3-3m-6.715 6.715L3 21m9-9l9 9" />
              </svg>
            </button>
          </div>
        </div>
        <div class="mb-6">
          <label for="subscriptionPackage" class="block text-text-muted text-sm font-bold mb-2">Paket Langganan:</label>
          <input
            type="text"
            id="subscriptionPackage"
            v-model="selectedPackageName"
            class="shadow appearance-none border rounded w-full py-2 px-3 text-text-base bg-bg-base leading-tight focus:outline-none focus:shadow-outline cursor-not-allowed"
            readonly
          />
        </div>
        <button
          type="submit"
          class="btn btn-secondary w-full"
        >
          Daftar & Mulai Coba Gratis
        </button>
      </form>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import { ref, watch, onMounted } from 'vue'; // Import watch and onMounted
import { useToast } from 'vue-toastification';
import { useRoute } from 'vue-router'; // Import useRoute to get query params

export default {
  name: 'RegisterCompany',
  props: ['packageId'], // packageId is a prop
  setup(props) {
    const toast = useToast();
    const route = useRoute(); // Get route instance to access query params

    const form = ref({
      company_name: '',
      company_address: '',
      admin_email: '',
      admin_password: '',
      subscription_package_id: parseInt(props.packageId),
      billing_cycle: route.query.billingCycle || 'monthly', // Get billingCycle from query params
    });
    const selectedPackageName = ref('');
    const subscriptionPackages = ref([]);
    const passwordFieldType = ref('password');

    console.log('RegisterCompany props.packageId:', props.packageId);
    console.log('RegisterCompany route.query.billingCycle:', route.query.billingCycle);

    const togglePasswordVisibility = () => {
      passwordFieldType.value = passwordFieldType.value === 'password' ? 'text' : 'password';
    };

    const fetchSubscriptionPackages = async () => {
      try {
        const response = await axios.get('/api/subscription-packages');
        subscriptionPackages.value = response.data.data;
        console.log('Fetched subscription packages:', subscriptionPackages.value);
      } catch (error) {
        console.error('Error fetching subscription packages:', error);
      }
    };

    // Watch for changes in subscriptionPackages or subscription_package_id
    watch([subscriptionPackages, () => form.value.subscription_package_id], ([newPackages, newPackageId]) => {
      if (newPackages.length > 0 && newPackageId) {
        const selectedPackage = newPackages.find(pkg => pkg.id === newPackageId);
        if (selectedPackage) {
          selectedPackageName.value = selectedPackage.package_name; // Use package_name as per backend
          console.log('Selected package name updated to:', selectedPackageName.value);
        } else {
          selectedPackageName.value = 'Paket tidak ditemukan';
          console.warn('Selected package not found for ID:', newPackageId);
        }
      }
    }, { immediate: true }); // Run immediately on component creation

    onMounted(() => {
      fetchSubscriptionPackages(); // Call on component mounted
    });

    const registerCompany = async () => {
      try {
        const response = await axios.post('/api/register-company', form.value);
        toast.success(response.data.message);
        // Redirect logic
        setTimeout(() => {
          const currentHost = window.location.hostname;
          const parts = currentHost.split('.');
          let adminDomain = '';

          if (parts.length >= 2) {
            if (parts[parts.length - 1] === 'localhost') {
              adminDomain = `admin.localhost`;
            } else {
              const domain = parts[parts.length - 2] + '.' + parts[parts.length - 1];
              adminDomain = `admin.${domain}`;
            }
          } else {
            adminDomain = `admin.${currentHost}`;
          }

          const adminLoginURL = `${window.location.protocol}//${adminDomain}${window.location.port ? ':' + window.location.port : ''}/`;
          window.location.href = adminLoginURL;
        }, 2000);
      } catch (error) {
        console.error('Registration failed - full error object:', error);
        if (error.response) {
          console.error('Registration failed - error.response:', error.response);
          console.error('Registration failed - error.response.data:', error.response.data);
        }
        const errorMessage = error.response && error.response.data && error.response.data.message ? error.response.data.message : error.message;
        toast.error('Registration failed: ' + errorMessage);
      }
    };

    return {
      form,
      selectedPackageName,
      passwordFieldType,
      togglePasswordVisibility,
      registerCompany,
    };
  },
};
</script>