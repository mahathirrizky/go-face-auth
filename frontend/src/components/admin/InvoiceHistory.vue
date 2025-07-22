<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-6">Riwayat Tagihan & Pembayaran</h2>

    <BaseDataTable
      :data="invoices"
      :columns="invoiceColumns"
      :loading="isLoading"
      :totalRecords="invoices.length"
      :lazy="false"
      searchPlaceholder="Cari Tagihan..."
    >
      <template #column-issued_at="{ item }">
        {{ new Date(item.issued_at).toLocaleDateString() }}
      </template>

      <template #column-due_date="{ item }">
        {{ new Date(item.due_date).toLocaleDateString() }}
      </template>

      <template #column-paid_at="{ item }">
        {{ item.paid_at ? new Date(item.paid_at).toLocaleDateString() : '-' }}
      </template>

      <template #column-amount="{ item }">
        {{ new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(item.amount) }}
      </template>

      <template #column-status="{ item }">
        <span :class="{
          'px-2 inline-flex text-xs leading-5 font-semibold rounded-full': true,
          'bg-green-100 text-green-600': item.status === 'paid',
          'bg-yellow-100 text-yellow-600': item.status === 'pending',
          'bg-red-100 text-red-600': ['failed', 'expired', 'cancelled'].includes(item.status),
        }">
          {{ item.status }}
        </span>
      </template>

      <template #column-actions="{ item }">
        <BaseButton
          v-if="item.status === 'paid'"
          @click="downloadInvoice(item.order_id)"
          class="btn-secondary btn-sm"
        >
          <i class="pi pi-download"></i> <span class="hidden sm:inline">Unduh PDF</span>
        </BaseButton>
      </template>
    </BaseDataTable>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import BaseDataTable from '../ui/BaseDataTable.vue';
import BaseButton from '../ui/BaseButton.vue';

const invoices = ref([]);
const isLoading = ref(false);
const toast = useToast();

const invoiceColumns = ref([
    { field: 'order_id', header: 'Order ID' },
    { field: 'subscription_package.package_name', header: 'Paket Langganan' },
    { field: 'billing_cycle', header: 'Siklus' },
    { field: 'amount', header: 'Jumlah' },
    { field: 'status', header: 'Status' },
    { field: 'issued_at', header: 'Tanggal Dibuat' },
    { field: 'paid_at', header: 'Tanggal Dibayar' },
    { field: 'actions', header: 'Aksi' },
]);

const fetchInvoices = async () => {
  isLoading.value = true;
  try {
    const response = await axios.get('/api/invoices');
    if (response.data && response.data.status === 'success') {
      invoices.value = response.data.data;
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data?.message || 'Gagal mengambil riwayat tagihan.', life: 3000 });
    }
  } catch (error) {
    console.error('Error fetching invoices:', error);
    let message = 'Terjadi kesalahan saat mengambil riwayat tagihan.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  } finally {
    isLoading.value = false;
  }
};

const downloadInvoice = async (orderId) => {
  try {
    const response = await axios.get(`/api/invoices/${orderId}/download`, {
      responseType: 'blob',
    });

    const blob = new Blob([response.data], { type: 'application/pdf' });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = `Invoice-${orderId}.pdf`;
    link.click();
    URL.revokeObjectURL(link.href);
    toast.add({ severity: 'success', summary: 'Success', detail: 'Invoice PDF berhasil diunduh!', life: 3000 });
  } catch (error) {
    console.error('Error downloading invoice:', error);
    let message = 'Gagal mengunduh invoice.';
    if (error.response && error.response.data) {
        // Since the response is a blob, we need to read it as text to get the error message
        const reader = new FileReader();
        reader.onload = () => {
            try {
                const errorData = JSON.parse(reader.result);
                message = errorData.message || message;
            } catch (e) {
                // The error response might not be JSON
            }
            toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
        };
        reader.readAsText(error.response.data);
    } else {
        toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
    }
  }
};

onMounted(() => {
  fetchInvoices();
});
</script>
