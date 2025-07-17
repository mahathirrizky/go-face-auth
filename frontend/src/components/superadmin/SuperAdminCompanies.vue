<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-4">Manajemen Perusahaan</h2>

    <BaseDataTable
      :data="companies"
      :columns="companyColumns"
      :loading="loading"
      v-model:filters="filters"
      :globalFilterFields="['name', 'address', 'subscription_status']"
      searchPlaceholder="Cari Perusahaan..."
    >
      <template #column-subscription_package="{ item }">
        {{ item.subscription_package ? item.subscription_package.name : 'N/A' }}
      </template>

      <template #column-created_at="{ item }">
        {{ new Date(item.created_at).toLocaleDateString() }}
      </template>
    </BaseDataTable>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { FilterMatchMode } from 'primevue/api';
import BaseDataTable from '../ui/BaseDataTable.vue';

const companies = ref([]);
const loading = ref(true);
const error = ref(null);
const toast = useToast();

const filters = ref({
    global: { value: null, matchMode: FilterMatchMode.CONTAINS }
});

const companyColumns = ref([
    { field: 'id', header: 'ID' },
    { field: 'name', header: 'Nama Perusahaan' },
    { field: 'address', header: 'Alamat' },
    { field: 'subscription_status', header: 'Status Langganan' },
    { field: 'subscription_package', header: 'Paket Langganan' },
    { field: 'created_at', header: 'Tanggal Dibuat' }
]);

const fetchCompanies = async () => {
  try {
    const response = await axios.get('/api/superadmin/companies');
    if (response.data && response.data.status === 'success') {
      companies.value = response.data.data;
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data.message || 'Failed to fetch companies.', life: 3000 });
      error.value = response.data.message || 'Failed to fetch companies.';
    }
  } catch (err) {
    console.error('Error fetching companies:', err);
    toast.add({ severity: 'error', summary: 'Error', detail: 'An error occurred while fetching companies.', life: 3000 });
    error.value = err.message;
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchCompanies();
});
</script>

<style scoped>
/* Tailwind handles styling */
</style>