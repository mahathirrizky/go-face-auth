<template>
  <div class="payment-status-container">
    <h1 class="title">{{ statusTitle }}</h1>
    <p class="message">{{ statusMessage }}</p>
    <div v-if="isLoading" class="loading-spinner"></div>
    <div v-if="!isLoading && invoice" class="invoice-details">
      <p><strong>Order ID:</strong> {{ invoice.order_id }}</p>
      <p><strong>Amount:</strong> {{ invoice.amount }}</p>
      <p><strong>Status:</strong> {{ invoice.status }}</p>
      <p v-if="invoice.paid_at"><strong>Paid At:</strong> {{ new Date(invoice.paid_at).toLocaleString() }}</p>
    </div>
    <button v-if="!isLoading && (status === 'success' || status === 'failed' || status === 'expired')" @click="performRedirect" class="btn-action">
      {{ getRedirectButtonText() }}
    </button>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'PaymentStatusPage',
  data() {
    return {
      status: 'loading', // success, failed, expired, pending, loading
      statusTitle: 'Processing Payment...',
      statusMessage: 'Please wait while we confirm your payment status.',
      invoice: null,
      isLoading: true,
    };
  },
  async created() {
    const orderId = this.$route.query.order_id;
    const transactionStatus = this.$route.query.transaction_status;

    if (!orderId) {
      this.status = 'error';
      this.statusTitle = 'Error';
      this.statusMessage = 'Order ID not found in URL.';
      this.isLoading = false;
      return;
    }

    // Provide immediate feedback based on URL parameter for better UX
    if (transactionStatus === 'settlement') {
      this.status = 'success';
      this.statusTitle = 'Payment Successful!';
      this.statusMessage = 'Your subscription is being activated. Please wait a moment...';
    } else if (transactionStatus === 'pending') {
      this.status = 'pending';
      this.statusTitle = 'Payment Pending';
      this.statusMessage = 'Your payment is still pending. We will update the status shortly.';
    }

    // Verify the final status with the backend
    try {
      // Poll the backend a few times to give the webhook time to process
      for (let i = 0; i < 5; i++) { // Poll up to 5 times
        const response = await axios.get(`/api/invoices/${orderId}`);
        this.invoice = response.data.data; // Assuming API returns { data: invoiceObject }

        if (this.invoice && this.invoice.status === 'paid') {
          break; // Exit loop if status is paid
        }
        // Wait 2 seconds before retrying
        if (i < 4) await new Promise(resolve => setTimeout(resolve, 2000));
      }

      if (this.invoice && this.invoice.status) {
        switch (this.invoice.status) {
          case 'paid':
            this.status = 'success';
            this.statusTitle = 'Payment Successful!';
            this.statusMessage = 'Your subscription has been activated. You can now access all features.';
            break;
          case 'pending':
            this.status = 'pending';
            this.statusTitle = 'Payment Pending';
            this.statusMessage = 'Your payment is still pending. Please complete the payment process or wait for confirmation.';
            break;
          case 'failed':
          case 'deny':
          case 'cancel':
            this.status = 'failed';
            this.statusTitle = 'Payment Failed';
            this.statusMessage = 'Your payment could not be processed. Please try again.';
            break;
          case 'expire':
            this.status = 'expired';
            this.statusTitle = 'Payment Expired';
            this.statusMessage = 'The payment window has expired. Please initiate a new payment.';
            break;
          default:
            this.status = 'unknown';
            this.statusTitle = 'Unknown Payment Status';
            this.statusMessage = 'We could not determine your payment status. Please contact support.';
        }
      } else {
        this.status = 'error';
        this.statusTitle = 'Error';
        this.statusMessage = 'Could not retrieve final invoice status from backend.';
      }
    } catch (error) {
      console.error('Error fetching invoice status:', error);
      this.status = 'error';
      this.statusTitle = 'Error';
      this.statusMessage = 'There was an error confirming your payment. Please try again or contact support.';
    } finally {
      this.isLoading = false;
    }
  },
  methods: {
    performRedirect() {
      if (this.status === 'success') {
        const hostname = window.location.hostname;
        const port = window.location.port ? `:${window.location.port}` : '';
        let baseUrl = '';

        if (hostname === 'localhost') {
          // For local development, redirect to admin.localhost which needs to be in /etc/hosts
          baseUrl = hostname + port;
        } else {
          // For production, derive the base domain (e.g., from www.example.com to example.com)
          const parts = hostname.split('.');
          baseUrl = parts.slice(-2).join('.');
        }

        const newUrl = `${window.location.protocol}//admin.${baseUrl}`;
        window.location.href = newUrl;

      } else {
        this.$router.push('/'); // Redirect to home for failed/expired cases
      }
    },
    getRedirectButtonText() {
      if (this.status === 'success') {
        return 'Go to Admin Portal';
      } else {
        return 'Return to Home';
      }
    }
  }
};
</script>

<style scoped>
.payment-status-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 80vh;
  text-align: center;
  padding: 20px;
  background-color: #f0f2f5;
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  max-width: 600px;
  margin: 50px auto;
}

.title {
  font-size: 2.5em;
  color: #333;
  margin-bottom: 15px;
}

.message {
  font-size: 1.2em;
  color: #555;
  margin-bottom: 30px;
}

.loading-spinner {
  border: 4px solid rgba(0, 0, 0, 0.1);
  border-left-color: #007bff;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  animation: spin 1s linear infinite;
  margin-bottom: 20px;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.invoice-details {
  background-color: #e9ecef;
  padding: 15px;
  border-radius: 5px;
  margin-top: 20px;
  text-align: left;
  width: 100%;
}

.invoice-details p {
  margin: 5px 0;
  color: #444;
}

.btn-action {
  background-color: #007bff;
  color: white;
  padding: 10px 20px;
  border-radius: 5px;
  text-decoration: none;
  font-size: 1.1em;
  margin-top: 30px;
  transition: background-color 0.3s ease;
}

.btn-action:hover {
  background-color: #0056b3;
}

/* Status specific colors */
.payment-status-container .title[data-status="success"] { color: #28a745; }
.payment-status-container .title[data-status="failed"],
.payment-status-container .title[data-status="expired"] { color: #dc3545; }
.payment-status-container .title[data-status="pending"] { color: #ffc107; }
.payment-status-container .title[data-status="error"] { color: #6c757d; }
</style>
