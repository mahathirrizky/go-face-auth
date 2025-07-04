<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100">
    <div class="bg-white p-8 rounded-lg shadow-md w-full max-w-md">
      <h2 class="text-2xl font-bold text-center mb-6">Daftar Perusahaan Baru</h2>
      <form @submit.prevent="registerCompany">
        <div class="mb-4">
          <label for="companyName" class="block text-gray-700 text-sm font-bold mb-2">Nama Perusahaan:</label>
          <input
            type="text"
            id="companyName"
            v-model="form.company_name"
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            required
          />
        </div>
        <div class="mb-4">
          <label for="companyAddress" class="block text-gray-700 text-sm font-bold mb-2">Alamat Perusahaan:</label>
          <input
            type="text"
            id="companyAddress"
            v-model="form.company_address"
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
          />
        </div>
        <div class="mb-4">
          <label for="adminEmail" class="block text-gray-700 text-sm font-bold mb-2">Email Admin:</label>
          <input
            type="email"
            id="adminEmail"
            v-model="form.admin_email"
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            required
          />
        </div>
        <div class="mb-6">
          <label for="adminPassword" class="block text-gray-700 text-sm font-bold mb-2">Password Admin:</label>
          <input
            type="password"
            id="adminPassword"
            v-model="form.admin_password"
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:shadow-outline"
            required
          />
        </div>
        <div class="mb-6">
          <label for="subscriptionPackage" class="block text-gray-700 text-sm font-bold mb-2">Paket Langganan:</label>
          <input
            type="text"
            id="subscriptionPackage"
            v-model="selectedPackageName"
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline bg-gray-100 cursor-not-allowed"
            readonly
          />
        </div>
        <button
          type="submit"
          class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline w-full"
        >
          Daftar & Lanjutkan Pembayaran
        </button>
      </form>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'RegisterCompany',
  props: ['packageId'],
  data() {
    return {
      form: {
        company_name: '',
        company_address: '',
        admin_email: '',
        admin_password: '',
        subscription_package_id: null,
      },
      selectedPackageName: '',
      subscriptionPackages: [], // To store fetched packages
    };
  },
  created() {
    this.form.subscription_package_id = parseInt(this.packageId);
    this.fetchSubscriptionPackages();
  },
  methods: {
    async fetchSubscriptionPackages() {
      try {
        const response = await axios.get('/api/subscription-packages');
        this.subscriptionPackages = response.data.data;
        const selectedPackage = this.subscriptionPackages.find(pkg => pkg.id === this.form.subscription_package_id);
        if (selectedPackage) {
          this.selectedPackageName = selectedPackage.name;
        }
      } catch (error) {
        console.error('Error fetching subscription packages:', error);
      }
    },
    async registerCompany() {
      try {
        const response = await axios.post('/api/register-company', this.form);
        // Redirect to payment page or show success message
        this.$router.push({
          name: 'PaymentPage',
          params: { companyId: response.data.data.company_id },
          query: { packageId: this.packageId } // Pass as query param
        });
      } catch (error) {
        console.error('Registration failed - full error object:', error);
        if (error.response) {
          console.error('Registration failed - error.response:', error.response);
          console.error('Registration failed - error.response.data:', error.response.data);
        }
        const errorMessage = error.response && error.response.data && error.response.data.message ? error.response.data.message : error.message;
        alert('Registration failed: ' + errorMessage);
      }
    },
  },
};
</script>