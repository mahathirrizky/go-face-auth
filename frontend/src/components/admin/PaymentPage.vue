<template>
  <div class="min-h-screen flex items-center justify-center bg-bg-base">
    <div class="bg-bg-muted p-8 rounded-lg shadow-md w-full max-w-md text-center">
      <h2 class="text-2xl font-bold mb-4 text-text-base">Halaman Pembayaran</h2>
      <p class="text-text-muted mb-4">Anda akan diarahkan ke halaman pembayaran Midtrans.</p>
      <p class="text-text-muted mb-4">Company ID: {{ companyId }}</p>
      <p class="text-text-muted mb-4">Package ID: {{ packageId.value }}</p>

      <div v-if="loading" class="mt-4 text-accent">
        <p>Memproses pembayaran Anda...</p>
        <p>Mohon tunggu sebentar, Anda akan segera diarahkan ke halaman Midtrans.</p>
      </div>
      <div v-else-if="error" class="mt-4 text-danger">{{ error }}</div>
      <div v-else class="mt-4 text-green-500">
        <p>Berhasil membuat transaksi. Mengarahkan ke Midtrans...</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import axios from 'axios';
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';

const props = defineProps(['companyId']);
const route = useRoute();
const loading = ref(true);
const error = ref(null);

const companyId = props.companyId;
const packageId = ref(route.query.packageId);
const billingCycle = ref(route.query.billingCycle);

console.log('PaymentPage setup: companyId:', companyId);
console.log('PaymentPage setup: packageId:', packageId.value);
console.log('PaymentPage setup: billingCycle:', billingCycle.value);

const createMidtransTransaction = async () => {
  loading.value = true;
  error.value = null;
  console.log('Attempting to create Midtrans transaction...');
  try {
    const payload = {
      company_id: parseInt(companyId),
      subscription_package_id: parseInt(packageId.value),
      billing_cycle: billingCycle.value,
    };
    console.log('Payload being sent:', payload);

    const response = await axios.post(
      '/api/midtrans/create-transaction',
      payload
    );
    console.log('Axios post call completed successfully.');

    console.log('Midtrans transaction response:', response);
    console.log('Midtrans transaction data:', response.data);

    if (response.data && response.data.data && response.data.data.redirect_url) {
      console.log('Redirecting to Midtrans URL:', response.data.data.redirect_url);
      window.location.href = response.data.data.redirect_url;
    } else {
      error.value = 'Redirect URL not found in Midtrans response.';
      console.error(error.value, response.data);
    }
  } catch (err) {
    console.error('An error occurred during Midtrans transaction creation.');
    console.error('Full error object:', err);
    if (err.response) {
      console.error('Error response data:', err.response.data);
      console.error('Error response status:', err.response.status);
      console.error('Error response headers:', err.response.headers);
      error.value = err.response.data.message || `Error: ${err.response.status}`;
    } else if (err.request) {
      console.error('Error request:', err.request);
      error.value = 'Network Error: No response received from server.';
    } else {
      console.error('Error message:', err.message);
      error.value = err.message;
    }
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  const companyIdNum = parseInt(companyId);
  const packageIdNum = packageId.value ? parseInt(packageId.value) : NaN;

  if (isNaN(companyIdNum) || !companyId) {
    error.value = 'Invalid or missing Company ID.';
    loading.value = false;
    console.error('PaymentPage Error: Invalid or missing Company ID.', companyId);
    return;
  }
  if (isNaN(packageIdNum) || !packageId.value) {
    error.value = 'Invalid or missing Package ID.';
    loading.value = false;
    console.error('PaymentPage Error: Invalid or missing Package ID.', packageId.value);
    return;
  }
  if (!billingCycle.value) {
    error.value = 'Missing Billing Cycle.';
    loading.value = false;
    console.error('PaymentPage Error: Missing Billing Cycle.');
    return;
  }

  createMidtransTransaction();
});
</script>