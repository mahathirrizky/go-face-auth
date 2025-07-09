<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-4">Manajemen Perusahaan</h2>

    <div v-if="loading" class="text-text-muted">Loading companies...</div>
    <div v-else-if="error" class="text-danger">Error: {{ error }}</div>
    <div v-else>
      <div class="bg-bg-muted p-4 rounded-lg shadow-md">
        <table class="min-w-full divide-y divide-bg-base">
          <thead class="bg-bg-base">
            <tr>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-text-muted uppercase tracking-wider">ID</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-text-muted uppercase tracking-wider">Nama Perusahaan</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-text-muted uppercase tracking-wider">Alamat</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-text-muted uppercase tracking-wider">Status Langganan</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-text-muted uppercase tracking-wider">Paket Langganan</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-text-muted uppercase tracking-wider">Tanggal Dibuat</th>
            </tr>
          </thead>
          <tbody class="bg-bg-muted divide-y divide-bg-base">
            <tr v-for="company in companies" :key="company.id">
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-text-base">{{ company.id }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-text-muted">{{ company.name }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-text-muted">{{ company.address }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-text-muted">{{ company.subscription_status }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-text-muted">{{ company.subscription_package ? company.subscription_package.name : 'N/A' }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-text-muted">{{ new Date(company.created_at).toLocaleDateString() }}</td>
            </tr>
            <tr v-if="companies.length === 0">
              <td colspan="6" class="px-6 py-4 whitespace-nowrap text-sm text-text-muted text-center">Tidak ada perusahaan ditemukan.</td>
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
  name: 'SuperAdminCompanies',
  setup() {
    const companies = ref([]);
    const loading = ref(true);
    const error = ref(null);
    const toast = useToast();

    const fetchCompanies = async () => {
      try {
        const response = await axios.get('/api/superadmin/companies');
        if (response.data && response.data.status === 'success') {
          companies.value = response.data.data;
        } else {
          toast.error(response.data.message || 'Failed to fetch companies.');
          error.value = response.data.message || 'Failed to fetch companies.';
        }
      } catch (err) {
        console.error('Error fetching companies:', err);
        toast.error('An error occurred while fetching companies.');
        error.value = err.message;
      } finally {
        loading.value = false;
      }
    };

    onMounted(() => {
      fetchCompanies();
    });

    return {
      companies,
      loading,
      error,
    };
  },
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>
