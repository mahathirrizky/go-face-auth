<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-4">Manajemen Perusahaan</h2>

    <DataTable
      :value="companies"
      :loading="loading"
      v-model:filters="filters"
      :globalFilterFields="['name', 'address', 'subscription_status']"
      paginator
      :rows="10"
      class="p-datatable-customers"
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink CurrentPageReport RowsPerPageDropdown"
      :rowsPerPageOptions="[10, 25, 50]"
      currentPageReportTemplate="Menampilkan {first} sampai {last} dari {totalRecords} data"
      dataKey="id"
    >
      <template #header>
        <div class="flex flex-wrap justify-between items-center gap-4">
          <IconField iconPosition="left">
            <InputIcon class="pi pi-search"></InputIcon>
            <InputText v-model="filters['global'].value" placeholder="Cari Perusahaan..." fluid />
          </IconField>
        </div>
      </template>
      <template #empty>
        Tidak ada data ditemukan.
      </template>
      <template #loading>
        Memuat data...
      </template>

      <Column field="id" header="ID" :sortable="true"></Column>
      <Column field="name" header="Nama Perusahaan" :sortable="true"></Column>
      <Column field="address" header="Alamat" :sortable="true"></Column>
      <Column field="subscription_status" header="Status Langganan" :sortable="true"></Column>
      <Column field="subscription_package.package_name" header="Paket Langganan" :sortable="true">
         <template #body="{ data }">
            {{ data.subscription_package ? data.subscription_package.package_name : 'N/A' }}
        </template>
      </Column>
      <Column field="billing_cycle" header="Siklus Penagihan" :sortable="true"></Column>
      <Column field="created_at" header="Tanggal Dibuat" :sortable="true">
        <template #body="{ data }">
          {{ new Date(data.created_at).toLocaleDateString() }}
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

const companies = ref([]);
const loading = ref(true);
const error = ref(null);
const toast = useToast();

const filters = ref({
    global: { value: null, matchMode: FilterMatchMode.CONTAINS }
});

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
