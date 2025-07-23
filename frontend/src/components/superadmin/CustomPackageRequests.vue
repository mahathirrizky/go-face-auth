<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-6">Permintaan Paket Kustom</h2>

    <BaseDataTable
      :data="customRequests"
      :columns="requestColumns"
      :loading="isLoading"
      :totalRecords="totalRecords"
      :lazy="true"
      v-model:filters="filters"
      @page="onPage"
      @filter="onFilter"
      searchPlaceholder="Cari Permintaan..."
    >
      <template #column-status="{ item }">
        <span :class="{
          'px-2 inline-flex text-xs leading-5 font-semibold rounded-full': true,
          'bg-yellow-100 text-yellow-800': item.status === 'pending',
          'bg-green-100 text-green-800': item.status === 'contacted',
          'bg-blue-100 text-blue-800': item.status === 'resolved',
        }">
          {{ item.status }}
        </span>
      </template>

      <template #column-actions="{ item }">
        <div class="flex space-x-2">
          <BaseButton @click="markAsContacted(item.id)" class="btn-info btn-sm" v-if="item.status === 'pending'">
            <i class="pi pi-phone"></i> Tandai Dihubungi
          </BaseButton>
          <BaseButton @click="markAsResolved(item.id)" class="btn-success btn-sm" v-if="item.status === 'contacted'">
            <i class="pi pi-check"></i> Tandai Selesai
          </BaseButton>
        </div>
      </template>
    </BaseDataTable>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { useWebSocketStore } from '../../stores/websocket';
import { FilterMatchMode } from '@primevue/core/api';

import BaseDataTable from '../ui/BaseDataTable.vue';
import BaseButton from '../ui/BaseButton.vue';

const customRequests = ref([]);
const isLoading = ref(false);
const totalRecords = ref(0);
const lazyParams = ref({});
const toast = useToast();
const webSocketStore = useWebSocketStore();

const filters = ref({
  'global': { value: null, matchMode: FilterMatchMode.CONTAINS },
});

const requestColumns = ref([
  { field: 'company_name', header: 'Perusahaan' },
  { field: 'name', header: 'Nama Kontak' },
  { field: 'email', header: 'Email Kontak' },
  { field: 'phone', header: 'Telepon Kontak' },
  { field: 'message', header: 'Pesan' },
  { field: 'status', header: 'Status' },
  { field: 'created_at', header: 'Tanggal Permintaan' },
  { field: 'actions', header: 'Aksi', sortable: false },
]);

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
