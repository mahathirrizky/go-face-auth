<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-6">Manajemen Karyawan</h2>

    <div class="bg-bg-muted p-4 rounded-lg shadow-md mb-6 flex flex-col md:flex-row justify-between items-center">
      <input
        type="text"
        placeholder="Cari karyawan..."
        v-model="searchTerm"
        class="w-full md:w-1/3 p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary mb-4 md:mb-0"
      />
      <button @click="openAddModal" class="btn btn-secondary w-full md:w-auto">
        Tambah Karyawan
      </button>
    </div>

    <div class="overflow-x-auto bg-bg-muted rounded-lg shadow-md">
      <table class="min-w-full divide-y divide-bg-base">
        <thead class="bg-primary">
          <tr>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Nama</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Email</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Jabatan</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-bg-base">
          <tr v-for="employee in employees" :key="employee.id">
            <td class="px-6 py-4 whitespace-nowrap text-text-base">{{ employee.name }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ employee.email }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ employee.position }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
              <button @click="openEditModal(employee)" class="text-accent hover:text-secondary mr-3">Edit</button>
              <button @click="deleteEmployee(employee.id)" class="text-danger hover:opacity-80">Hapus</button>
            </td>
          </tr>
          <tr v-if="employees.length === 0">
            <td colspan="4" class="px-6 py-4 text-center text-text-muted">Tidak ada data karyawan.</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Add/Edit Employee Modal -->
    <div v-if="isModalOpen" class="fixed inset-0 bg-black bg-opacity-20 flex items-center justify-center z-50">
      <div class="bg-bg-muted p-8 rounded-lg shadow-lg w-full max-w-md">
        <h3 class="text-2xl font-bold text-text-base mb-6">{{ editingEmployee ? 'Edit Karyawan' : 'Tambah Karyawan' }}</h3>
        <form @submit.prevent="saveEmployee">
          <div class="mb-4">
            <label for="name" class="block text-text-muted text-sm font-bold mb-2">Nama:</label>
            <input
              type="text"
              id="name"
              v-model="currentEmployee.name"
              class="w-full p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary"
              required
            />
          </div>
          <div class="mb-4">
            <label for="email" class="block text-text-muted text-sm font-bold mb-2">Email:</label>
            <input
              type="email"
              id="email"
              v-model="currentEmployee.email"
              class="w-full p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary"
              required
            />
          </div>
          <div class="mb-6">
            <label for="position" class="block text-text-muted text-sm font-bold mb-2">Jabatan:</label>
            <input
              type="text"
              id="position"
              v-model="currentEmployee.position"
              class="w-full p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary"
              required
            />
          </div>
          <div class="flex justify-end space-x-4">
            <button type="button" @click="closeModal" class="btn btn-outline-primary">
              Batal
            </button>
            <button type="submit" class="btn btn-secondary">
              Simpan
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue';
import axios from 'axios';
import { useToast } from 'vue-toastification';
import { useAuthStore } from '../../stores/auth';

export default {
  name: 'EmployeeManagement',
  setup() {
    const employees = ref([]);
    const isModalOpen = ref(false);
    const currentEmployee = ref({});
    const searchTerm = ref('');
    const editingEmployee = ref(false);
    const toast = useToast();
    const authStore = useAuthStore();

    // Function to fetch employees from backend
    const fetchEmployees = async () => {
      if (!authStore.companyId) {
        toast.error('Company ID not available. Cannot fetch employees.');
        return;
      }
      try {
        const url = searchTerm.value
          ? `/api/companies/${authStore.companyId}/employees/search?name=${searchTerm.value}`
          : `/api/companies/${authStore.companyId}/employees`;
        
        const response = await axios.get(url);
        if (response.data && response.data.status === 'success') {
          // Ensure employees.value is always an array
          employees.value = Array.isArray(response.data.data) ? response.data.data : [];
          
          // Only show warning if response.data.data was not an array but was present
          if (response.data.data !== undefined && response.data.data !== null && !Array.isArray(response.data.data)) {
            toast.warning('Received non-array data for employees, treating as empty list.');
          }
        } else {
          console.log('Unexpected API response for employees:', response); // Log the full response
          employees.value = []; // Clear employees on error or non-success status
          toast.error(response.data?.message || 'Failed to fetch employees due to an unexpected response format.');
        }
      } catch (error) {
        console.error('Error fetching employees:', error);
        let message = 'Failed to fetch employees.';
        if (error.response && error.response.data && error.response.data.message) {
          message = error.response.data.message;
        }
        toast.error(message);
      }
    };

    // Initial fetch when component is mounted
    onMounted(() => {
      fetchEmployees();
    });

    // Debounce for search input
    let searchTimeout = null;
    watch(searchTerm, (newSearchTerm) => {
      clearTimeout(searchTimeout);
      searchTimeout = setTimeout(() => {
        fetchEmployees();
      }, 300);
    });

    const openAddModal = () => {
      currentEmployee.value = { name: '', email: '', position: '' };
      editingEmployee.value = false;
      isModalOpen.value = true;
    };

    const openEditModal = (employee) => {
      currentEmployee.value = { ...employee };
      editingEmployee.value = true;
      isModalOpen.value = true;
    };

    const closeModal = () => {
      isModalOpen.value = false;
      currentEmployee.value = {};
    };

    const saveEmployee = async () => {
      if (!authStore.companyId) {
        toast.error('Company ID not available. Cannot save employee.');
        return;
      }
      try {
        if (currentEmployee.value.id) {
          console.log('Updating employee:', currentEmployee.value);
          const response = await axios.put(`/api/employees/${currentEmployee.value.id}`, currentEmployee.value);
          toast.success(response.data.message || 'Employee updated successfully!');
        } else {
          console.log('Creating employee:', currentEmployee.value);
          const response = await axios.post(`/api/employees`, {
            name: currentEmployee.value.name,
            email: currentEmployee.value.email,
            position: currentEmployee.value.position,
          });
          toast.success(response.data.message || 'Employee created successfully. An email with initial password setup link has been sent.');
        }
        closeModal();
        fetchEmployees();
      } catch (error) {
        console.error('Error saving employee:', error);
        let message = 'Failed to save employee.';
        if (error.response && error.response.data && error.response.data.message) {
          message = error.response.data.message;
        }
        toast.error(message);
      }
    };

    const deleteEmployee = async (id) => {
      if (confirm('Apakah Anda yakin ingin menghapus karyawan ini?')) {
        try {
          const response = await axios.delete(`/api/employees/${id}`);
          toast.success(response.data.message || 'Employee deleted successfully!');
          fetchEmployees();
        } catch (error) {
          console.error('Error deleting employee:', error);
          let message = 'Failed to delete employee.';
          if (error.response && error.response.data && error.response.data.message) {
            message = error.response.data.message;
          }
          toast.error(message);
        }
      }
    };

    return {
      employees,
      isModalOpen,
      currentEmployee,
      searchTerm,
      editingEmployee,
      openAddModal,
      openEditModal,
      closeModal,
      saveEmployee,
      deleteEmployee,
    };
  },
};
</script>