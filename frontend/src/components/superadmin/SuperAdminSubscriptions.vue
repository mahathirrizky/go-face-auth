<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-4">Manajemen Langganan</h2>

    <BaseDataTable
      :data="subscriptions"
      :columns="subscriptionColumns"
      :loading="loading"
      :globalFilterFields="['name', 'subscription_status']"
      searchPlaceholder="Cari Langganan..."
    >
      <template #column-subscription_package="{ item }">
        {{ item.subscription_package ? item.subscription_package.name : 'N/A' }}
      </template>

      <template #column-trial_end_date="{ item }">
        {{ item.trial_end_date ? new Date(item.trial_end_date).toLocaleDateString() : 'N/A' }}
      </template>
    </BaseDataTable>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import BaseDataTable from '../ui/BaseDataTable.vue';

const subscriptions = ref([]);
const loading = ref(true);
const error = ref(null);
const toast = useToast();

const subscriptionColumns = ref([
    { field: 'id', header: 'ID Perusahaan' },
    { field: 'name', header: 'Nama Perusahaan' },
    { field: 'subscription_status', header: 'Status Langganan' },
    { field: 'subscription_package', header: 'Paket Langganan' },
    { field: 'trial_end_date', header: 'Tanggal Berakhir Percobaan' }
]);

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