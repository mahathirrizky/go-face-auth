<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-4">Manajemen Langganan</h2>

    <div v-if="loading" class="text-text-muted">Loading subscriptions...</div>
    <div v-else-if="error" class="text-danger">Error: {{ error }}</div>
    <div v-else>
      <div class="bg-bg-muted p-4 rounded-lg shadow-md">
        <table class="min-w-full divide-y divide-bg-base">
          <thead class="bg-bg-base">
            <tr>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-text-muted uppercase tracking-wider">ID Perusahaan</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-text-muted uppercase tracking-wider">Nama Perusahaan</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-text-muted uppercase tracking-wider">Status Langganan</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-text-muted uppercase tracking-wider">Paket Langganan</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-text-muted uppercase tracking-wider">Tanggal Berakhir Percobaan</th>
            </tr>
          </thead>
          <tbody class="bg-bg-muted divide-y divide-bg-base">
            <tr v-for="company in subscriptions" :key="company.id">
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-text-base">{{ company.id }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-text-muted">{{ company.name }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-text-muted">{{ company.subscription_status }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-text-muted">{{ company.subscription_package ? company.subscription_package.name : 'N/A' }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-text-muted">{{ company.trial_end_date ? new Date(company.trial_end_date).toLocaleDateString() : 'N/A' }}</td>
            </tr>
            <tr v-if="subscriptions.length === 0">
              <td colspan="5" class="px-6 py-4 whitespace-nowrap text-sm text-text-muted text-center">Tidak ada langganan ditemukan.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { useToast } from 'vue-toastification';

export default {
  name: 'SuperAdminSubscriptions',
  setup() {
    const subscriptions = ref([]);
    const loading = ref(true);
    const error = ref(null);
    const toast = useToast();

    const fetchSubscriptions = async () => {
      try {
        const response = await axios.get('/api/superadmin/subscriptions');
        if (response.data && response.data.status === 'success') {
          subscriptions.value = response.data.data;
        } else {
          toast.error(response.data.message || 'Failed to fetch subscriptions.');
          error.value = response.data.message || 'Failed to fetch subscriptions.';
        }
      } catch (err) {
        console.error('Error fetching subscriptions:', err);
        toast.error('An error occurred while fetching subscriptions.');
        error.value = err.message;
      } finally {
        loading.value = false;
      }
    };

    onMounted(() => {
      fetchSubscriptions();
    });

    return {
      subscriptions,
      loading,
      error,
    };
  },
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>
