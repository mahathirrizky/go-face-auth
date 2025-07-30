<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-6">Permintaan Paket Kustom</h2>

    <DataTable
      :value="customRequests"
      :loading="isLoading"
      :totalRecords="totalRecords"
      :lazy="true"
      v-model:filters="filters"
      @page="onPage"
      paginator
      :rows="10"
      :rowsPerPageOptions="[10, 25, 50]"
      currentPageReportTemplate="Menampilkan {first} sampai {last} dari {totalRecords} data"
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink CurrentPageReport RowsPerPageDropdown"
      dataKey="id"
      :globalFilterFields="['company_name', 'name', 'email', 'phone', 'message']"
    >
      <template #header>
        <div class="flex justify-end">
          <IconField iconPosition="left">
            <InputIcon class="pi pi-search"></InputIcon>
            <InputText v-model="filters['global'].value" placeholder="Cari Permintaan..." @keydown.enter="onFilter" fluid />
          </IconField>
        </div>
      </template>

      <template #empty>
        Tidak ada data ditemukan.
      </template>
      <template #loading>
        Memuat data...
      </template>

      <Column field="company_name" header="Perusahaan" :sortable="true"></Column>
      <Column field="name" header="Nama Kontak" :sortable="true"></Column>
      <Column field="email" header="Email Kontak" :sortable="true"></Column>
      <Column field="phone" header="Telepon Kontak" :sortable="true"></Column>
      <Column field="message" header="Pesan" :sortable="true"></Column>
      <Column field="status" header="Status" :sortable="true">
        <template #body="{ data }">
          <Tag :value="data.status" :severity="getStatusSeverity(data.status)" />
        </template>
      </Column>
      <Column field="created_at" header="Tanggal Permintaan" :sortable="true">
        <template #body="{ data }">
          {{ new Date(data.created_at).toLocaleString('id-ID') }}
        </template>
      </Column>
      <Column header="Aksi">
        <template #body="{ data }">
          <div class="flex space-x-2">
            <Button @click="markAsContacted(data.id)" class="p-button-info p-button-sm" v-if="data.status === 'pending'" icon="pi pi-phone" label="Hubungi" />
            <Button @click="markAsResolved(data.id)" class="p-button-success p-button-sm" v-if="data.status === 'contacted'" icon="pi pi-check" label="Selesai" />
          </div>
        </template>
      </Column>
    </DataTable>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { useWebSocketStore } from '../../stores/websocket';
import { FilterMatchMode } from '@primevue/core/api';

import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import IconField from 'primevue/iconfield';
import InputIcon from 'primevue/inputicon';
import Tag from 'primevue/tag';

const customRequests = ref([]);
const isLoading = ref(false);
const totalRecords = ref(0);
const lazyParams = ref({});
const toast = useToast();
const webSocketStore = useWebSocketStore();

const filters = ref({
  'global': { value: null, matchMode: FilterMatchMode.CONTAINS },
});

const fetchCustomRequests = async () => {
  isLoading.value = true;
  try {
    const params = {
      page: lazyParams.value.page + 1,
      limit: lazyParams.value.rows,
      search: filters.value.global.value || '',
    };
    const response = await axios.get('/api/superadmin/custom-package-requests', { params });
    if (response.data && response.data.status === 'success') {
      customRequests.value = response.data.data.items;
      totalRecords.value = response.data.data.total_records;
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data?.message || 'Gagal mengambil permintaan kustom.', life: 3000 });
    }
  } catch (error) {
    console.error('Error fetching custom requests:', error);
    toast.add({ severity: 'error', summary: 'Error', detail: 'Terjadi kesalahan saat mengambil permintaan kustom.', life: 3000 });
  } finally {
    isLoading.value = false;
  }
};

const onPage = (event) => {
  lazyParams.value = event;
  fetchCustomRequests();
};

const onFilter = () => {
  lazyParams.value.page = 0; // Reset to first page on filter
  fetchCustomRequests();
};

const markAsContacted = async (id) => {
  try {
    const response = await axios.put(`/api/superadmin/custom-package-requests/${id}/contacted`);
    if (response.data && response.data.status === 'success') {
      toast.add({ severity: 'success', summary: 'Success', detail: 'Permintaan ditandai sebagai telah dihubungi.', life: 3000 });
      fetchCustomRequests(); // Refresh data
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data?.message || 'Gagal menandai permintaan.', life: 3000 });
    }
  } catch (error) {
    console.error('Error marking as contacted:', error);
    toast.add({ severity: 'error', summary: 'Error', detail: 'Terjadi kesalahan saat menandai permintaan.', life: 3000 });
  }
};

const markAsResolved = async (id) => {
  try {
    const response = await axios.put(`/api/superadmin/custom-package-requests/${id}/resolved`);
    if (response.data && response.data.status === 'success') {
      toast.add({ severity: 'success', summary: 'Success', detail: 'Permintaan ditandai sebagai selesai.', life: 3000 });
      fetchCustomRequests(); // Refresh data
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data?.message || 'Gagal menandai permintaan.', life: 3000 });
    }
  } catch (error) {
    console.error('Error marking as resolved:', error);
    toast.add({ severity: 'error', summary: 'Error', detail: 'Terjadi kesalahan saat menandai permintaan.', life: 3000 });
  }
};

const getStatusSeverity = (status) => {
  switch (status) {
    case 'pending':
      return 'warning';
    case 'contacted':
      return 'info';
    case 'resolved':
      return 'success';
    default:
      return null;
  }
};

// WebSocket handling
const handleWebSocketMessage = (data) => {
  if (data.type === 'new_custom_package_request') {
    toast.add({ severity: 'info', summary: 'Permintaan Baru!', detail: `Permintaan paket kustom baru dari ${data.company_name}.`, life: 5000 });
    fetchCustomRequests(); // Refresh the list to show the new request
  }
};

onMounted(() => {
  lazyParams.value = {
    first: 0,
    rows: 10,
    page: 0,
  };
  fetchCustomRequests();
  webSocketStore.onMessage('superadmin_notification', handleWebSocketMessage);
});

onUnmounted(() => {
  webSocketStore.offMessage('superadmin_notification');
});
</script>