<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100">
    <div class="bg-white p-8 rounded-lg shadow-md w-full max-w-md text-center">
      <h2 class="text-2xl font-bold mb-4">Halaman Pembayaran</h2>
      <p class="text-gray-700 mb-4">Anda akan diarahkan ke halaman pembayaran Midtrans.</p>
      <p class="text-gray-700 mb-4">Company ID: {{ companyId }}</p>
      <p class="text-gray-700 mb-4">Package ID: {{ packageIdFromQuery }}</p>

      <div v-if="loading" class="mt-4 text-blue-500">
        <p>Memproses pembayaran Anda...</p>
        <p>Mohon tunggu sebentar, Anda akan segera diarahkan ke halaman Midtrans.</p>
      </div>
      <div v-else-if="error" class="mt-4 text-red-500">{{ error }}</div>
      <div v-else class="mt-4 text-green-500">
        <p>Berhasil membuat transaksi. Mengarahkan ke Midtrans...</p>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'PaymentPage',
  props: ['companyId'], // Only companyId is a prop from params
  data() {
    return {
      loading: true,
      error: null,
    };
  },
  computed: {
    packageIdFromQuery() {
      return this.$route.query.packageId;
    }
  },
  async created() {
    console.log('PaymentPage created. companyId (prop):', this.companyId, 'packageId (query):', this.packageIdFromQuery);
    if (this.companyId && this.packageIdFromQuery) {
      await this.createMidtransTransaction();
    } else {
      this.error = 'Missing company ID or package ID for payment.';
      this.loading = false;
    }
  },
  methods: {
    async createMidtransTransaction() {
      this.loading = true;
      this.error = null;
      console.log('Attempting to create Midtrans transaction with companyId:', this.companyId, 'and packageId:', this.packageIdFromQuery);
      try {
        const response = await axios.post(
          '/api/midtrans/create-transaction',
          {
            company_id: parseInt(this.companyId),
            subscription_package_id: parseInt(this.packageIdFromQuery),
          }
        );
        console.log('Midtrans transaction created:', response.data);
        console.log('Redirecting to Midtrans URL:', response.data.data.redirect_url);
        window.location.href = response.data.data.redirect_url;
      } catch (error) {
        console.error('Error creating Midtrans transaction:', error.response ? error.response.data : error.message);
        this.error = error.response ? error.response.data.message : error.message;
      } finally {
        this.loading = false;
      }
    },
  },
};
</script>
