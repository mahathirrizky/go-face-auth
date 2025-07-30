<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-4">Manajemen Langganan</h2>

    <DataTable
      :value="subscriptions"
      :loading="loading"
      v-model:filters="filters"
      :globalFilterFields="['name', 'subscription_status']"
      paginator
      :rows="10"
      class="p-datatable-customers"
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink CurrentPageReport RowsPerPageDropdown"
      :rowsPerPageOptions="[10, 25, 50]"
      currentPageReportTemplate="Menampilkan {first} sampai {last} dari {totalRecords} data"
      dataKey="id"
    >
       <template #header>
        <div class="flex justify-end">
          <IconField iconPosition="left">
            <InputIcon class="pi pi-search"></InputIcon>
            <InputText v-model="filters['global'].value" placeholder="Cari Langganan..." fluid />
          </IconField>
        </div>
      </template>
      <template #empty>
        Tidak ada data ditemukan.
      </template>
      <template #loading>
        Memuat data...
      </template>

      <Column field="id" header="ID Perusahaan" :sortable="true"></Column>
      <Column field="name" header="Nama Perusahaan" :sortable="true"></Column>
      <Column field="subscription_status" header="Status Langganan" :sortable="true"></Column>
      <Column field="subscription_package.name" header="Paket Langganan" :sortable="true">
        <template #body="{ data }">
          {{ data.subscription_package ? data.subscription_package.name : 'N/A' }}
        </template>
      </Column>
      <Column field="trial_end_date" header="Tanggal Berakhir Percobaan" :sortable="true">
        <template #body="{ data }">
          {{ data.trial_end_date ? new Date(data.trial_end_date).toLocaleDateString() : 'N/A' }}
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
import InputText from 'primevue/inputtext';
import IconField from 'primevue/iconfield';
import InputIcon from 'primevue/inputicon';

const subscriptions = ref([]);
const loading = ref(true);
const error = ref(null);
const toast = useToast();

const filters = ref({
    global: { value: null, matchMode: FilterMatchMode.CONTAINS }
});

const fetchSubscriptions = async () => {
  try {
    const response = await axios.get('/api/superadmin/subscriptions');
    if (response.data && response.data.status === 'success') {
      subscriptions.value = response.data.data;
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data.message || 'Failed to fetch subscriptions.', life: 3000 });
      error.value = response.data.message || 'Failed to fetch subscriptions.';
    }
  } catch (err) {
    console.error('Error fetching subscriptions:', err);
    toast.add({ severity: 'error', summary: 'Error', detail: 'An error occurred while fetching subscriptions.', life: 3000 });
    error.value = err.message;
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchSubscriptions();
});
</script>

<style scoped>
/* Tailwind handles styling */
</style>
