<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <Toast />
    <h2 class="text-2xl font-bold text-text-base mb-6">Riwayat Tagihan & Pembayaran</h2>

    <DataTable
      :value="invoices"
      :loading="isLoading"
      paginator
      :rows="10"
      :rowsPerPageOptions="[10, 25, 50]"
      currentPageReportTemplate="Menampilkan {first} sampai {last} dari {totalRecords} data"
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink CurrentPageReport RowsPerPageDropdown"
      dataKey="order_id"
      v-model:filters="filters"
      :globalFilterFields="['order_id', 'subscription_package.package_name', 'status']"
    >
      <template #header>
        <div class="flex justify-end">
          <IconField iconPosition="left">
            <InputIcon class="pi pi-search"></InputIcon>
            <InputText v-model="filters['global'].value" placeholder="Cari Tagihan..." fluid />
          </IconField>
        </div>
      </template>

      <template #empty>
        Tidak ada data ditemukan.
      </template>
      <template #loading>
        Memuat data...
      </template>

      <Column field="order_id" header="Order ID" :sortable="true"></Column>
      <Column field="subscription_package.package_name" header="Paket Langganan" :sortable="true"></Column>
      <Column field="billing_cycle" header="Siklus" :sortable="true"></Column>
      <Column field="amount" header="Jumlah" :sortable="true">
        <template #body="{ data }">
          {{ new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(data.amount) }}
        </template>
      </Column>
      <Column field="status" header="Status" :sortable="true">
        <template #body="{ data }">
          <Tag :value="data.status" :severity="getStatusSeverity(data.status)" />
        </template>
      </Column>
      <Column field="issued_at" header="Tanggal Dibuat" :sortable="true">
        <template #body="{ data }">
          {{ new Date(data.issued_at).toLocaleDateString('id-ID') }}
        </template>
      </Column>
      <Column field="paid_at" header="Tanggal Dibayar" :sortable="true">
        <template #body="{ data }">
          {{ data.paid_at ? new Date(data.paid_at).toLocaleDateString('id-ID') : '-' }}
        </template>
      </Column>
      <Column header="Aksi">
        <template #body="{ data }">
          <Button
            v-if="data.status === 'paid'"
            @click="downloadInvoice(data.order_id)"
            icon="pi pi-download"
            class="p-button-secondary p-button-sm mr-2"
            :loading="isDownloading"
            v-tooltip.top="'Unduh PDF'"
          />
          <Button
            v-if="data.status === 'pending' && data.payment_url"
            :as="'a'"
            :href="data.payment_url"
            target="_blank"
            rel="noopener noreferrer"
            icon="pi pi-money-bill"
            class="p-button-primary p-button-sm"
            v-tooltip.top="'Bayar Sekarang'"
          />
        </template>
      </Column>
    </DataTable>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { FilterMatchMode } from '@primevue/core/api';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import IconField from 'primevue/iconfield';
import InputIcon from 'primevue/inputicon';
import Tag from 'primevue/tag';
import Toast from 'primevue/toast';

const invoices = ref([]);
const isLoading = ref(false);
const isDownloading = ref(false);
const toast = useToast();

const filters = ref({
    global: { value: null, matchMode: FilterMatchMode.CONTAINS }
});

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
  isDownloading.value = true;
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
        const reader = new FileReader();
        reader.onload = () => {
            try {
                const errorData = JSON.parse(reader.result);
                message = errorData.message || message;
            } catch (e) {
                // Error response might not be JSON
            }
            toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
        };
        reader.readAsText(error.response.data);
    } else {
        toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
    }
  } finally {
    isDownloading.value = false;
  }
};

const getStatusSeverity = (status) => {
  switch (status) {
    case 'paid': return 'success';
    case 'pending': return 'warning';
    case 'failed':
    case 'expired':
    case 'cancelled': return 'danger';
    default: return 'info';
  }
};

onMounted(() => {
  fetchInvoices();
});
</script>