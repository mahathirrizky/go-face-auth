<template>
  <div class="min-h-screen flex items-center justify-center bg-bg-base p-4">
    <div class="bg-bg-muted p-8 rounded-lg shadow-md w-full max-w-md text-center">
      <h1 :class="['text-3xl font-bold mb-4', statusColor]">{{ statusTitle }}</h1>
      <p class="text-text-muted mb-6">{{ statusMessage }}</p>
      
      <div v-if="isLoading" class="flex justify-center items-center mb-4">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-secondary"></div>
      </div>

      <div v-if="!isLoading && invoice" class="bg-bg-base p-4 rounded-lg text-left text-sm mb-6">
        <p><strong class="text-text-base">Order ID:</strong> <span class="text-text-muted">{{ invoice.order_id }}</span></p>
        <p><strong class="text-text-base">Amount:</strong> <span class="text-text-muted">{{ invoice.amount }}</span></p>
        <p><strong class="text-text-base">Status:</strong> <span class="font-semibold" :class="statusColor">{{ invoice.status }}</span></p>
        <p v-if="invoice.paid_at"><strong class="text-text-base">Paid At:</strong> <span class="text-text-muted">{{ new Date(invoice.paid_at).toLocaleString() }}</p>
      </div>

      <BaseButton v-if="!isLoading && (status === 'success' || status === 'failed' || status === 'expired')" @click="performRedirect" class="w-full">
        {{ getRedirectButtonText() }}
      </BaseButton>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import axios from 'axios';
import { useRoute, useRouter } from 'vue-router';
import BaseButton from '../ui/BaseButton.vue';

const route = useRoute();
const router = useRouter();

const status = ref('loading'); // success, failed, expired, pending, loading
const statusTitle = ref('Processing Payment...');
const statusMessage = ref('Please wait while we confirm your payment status.');
const invoice = ref(null);
const isLoading = ref(true);

const statusColor = computed(() => {
  switch (status.value) {
    case 'success': return 'text-green-500';
    case 'failed':
    case 'expired': return 'text-danger';
    case 'pending': return 'text-yellow-500';
    case 'error': return 'text-gray-500';
    default: return 'text-text-base';
  }
});

const performRedirect = () => {
  if (status.value === 'success') {
    const hostname = window.location.hostname;
    const port = window.location.port ? `:${window.location.port}` : '';
    let baseUrl = '';

    if (hostname === 'localhost') {
      baseUrl = hostname + port;
    } else {
      const parts = hostname.split('.');
      baseUrl = parts.slice(-2).join('.');
    }

    const newUrl = `${window.location.protocol}//admin.${baseUrl}`;
    window.location.href = newUrl;

  } else {
    router.push('/');
  }
};

const getRedirectButtonText = () => {
  if (status.value === 'success') {
    return 'Go to Admin Portal';
  } else {
    return 'Return to Home';
  }
};

onMounted(async () => {
  const orderId = route.query.order_id;
  const transactionStatus = route.query.transaction_status;

  if (!orderId) {
    status.value = 'error';
    statusTitle.value = 'Error';
    statusMessage.value = 'Order ID not found in URL.';
    isLoading.value = false;
    return;
  }

  if (transactionStatus === 'settlement') {
    status.value = 'success';
    statusTitle.value = 'Payment Successful!';
    statusMessage.value = 'Your subscription is being activated. Please wait a moment...';
  } else if (transactionStatus === 'pending') {
    status.value = 'pending';
    statusTitle.value = 'Payment Pending';
    statusMessage.value = 'Your payment is still pending. We will update the status shortly.';
  }

  try {
    for (let i = 0; i < 5; i++) {
      const response = await axios.get(`/api/invoices/${orderId}`);
      invoice.value = response.data.data;

      if (invoice.value && invoice.value.status === 'paid') {
        break;
      }
      if (i < 4) await new Promise(resolve => setTimeout(resolve, 2000));
    }

    if (invoice.value && invoice.value.status) {
      switch (invoice.value.status) {
        case 'paid':
          status.value = 'success';
          statusTitle.value = 'Payment Successful!';
          statusMessage.value = 'Your subscription has been activated. You can now access all features.';
          break;
        case 'pending':
          status.value = 'pending';
          statusTitle.value = 'Payment Pending';
          statusMessage.value = 'Your payment is still pending. Please complete the payment process or wait for confirmation.';
          break;
        case 'failed':
        case 'deny':
        case 'cancel':
          status.value = 'failed';
          statusTitle.value = 'Payment Failed';
          statusMessage.value = 'Your payment could not be processed. Please try again.';
          break;
        case 'expire':
          status.value = 'expired';
          statusTitle.value = 'Payment Expired';
          statusMessage.value = 'The payment window has expired. Please initiate a new payment.';
          break;
        default:
          status.value = 'unknown';
          statusTitle.value = 'Unknown Payment Status';
          statusMessage.value = 'We could not determine your payment status. Please contact support.';
      }
    } else {
      status.value = 'error';
      statusTitle.value = 'Error';
      statusMessage.value = 'Could not retrieve final invoice status from backend.';
    }
  } catch (error) {
    console.error('Error fetching invoice status:', error);
    status.value = 'error';
    statusTitle.value = 'Error';
    statusMessage.value = 'There was an error confirming your payment. Please try again or contact support.';
  } finally {
    isLoading.value = false;
  }
});
</script>

<style scoped>
/* Tailwind handles styling */
</style>

