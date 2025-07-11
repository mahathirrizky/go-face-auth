<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-6">Manajemen Karyawan</h2>

    <div class="bg-bg-muted p-4 rounded-lg shadow-md mb-6 flex flex-col md:flex-row justify-between items-center">
      <BaseInput
        id="searchTerm"
        label="Cari karyawan..."
        v-model="searchTerm"
        placeholder="Cari karyawan..."
        :label-sr-only="true"
        class="w-full md:w-1/3 mb-4 md:mb-0"
      />
      <BaseButton @click="openAddModal" class="w-full md:w-auto">
        Tambah Karyawan
      </BaseButton>
    </div>

    <div class="overflow-x-auto bg-bg-muted rounded-lg shadow-md">
      <table class="min-w-full divide-y divide-bg-base">
        <thead class="bg-primary">
          <tr>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Nama</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Email</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Nomor ID</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Jabatan</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-bg-base">
          <tr v-for="employee in employees" :key="employee.id">
            <td class="px-6 py-4 whitespace-nowrap text-text-base">{{ employee.name }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ employee.email }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ employee.employee_id_number }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ employee.position }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
              <BaseButton @click="openEditModal(employee)" class="text-accent hover:text-secondary mr-3">Edit</BaseButton>
              <router-link :to="{ name: 'EmployeeAttendanceHistory', params: { employeeId: employee.id } }" class="text-info hover:text-info-dark mr-3">Riwayat Absensi</router-link>
              <BaseButton @click="deleteEmployee(employee.id)" class="text-danger hover:opacity-80">Hapus</BaseButton>
            </td>
          </tr>
          <tr v-if="employees.length === 0">
            <td colspan="4" class="px-6 py-4 text-center text-text-muted">Tidak ada data karyawan.</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Add/Edit Employee Modal -->
    <BaseModal :isOpen="isModalOpen" @close="closeModal" :title="editingEmployee ? 'Edit Karyawan' : 'Tambah Karyawan'" maxWidth="md">
      <form @submit.prevent="saveEmployee">
        <BaseInput
          id="name"
          label="Nama:"
          v-model="currentEmployee.name"
          required
        />
        <BaseInput
          id="email"
          label="Email:"
          v-model="currentEmployee.email"
          type="email"
          required
        />
        <BaseInput
          id="employeeIdNumber"
          label="Nomor ID Karyawan:"
          v-model="currentEmployee.employee_id_number"
          required
        />
        <BaseInput
          id="position"
          label="Jabatan:"
          v-model="currentEmployee.position"
          required
        />
        <div class="mb-6">
          <label for="shift" class="block text-text-muted text-sm font-bold mb-2">Shift:</label>
          <select
            id="shift"
            v-model="currentEmployee.shift_id"
            class="w-full p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary"
          >
            <option :value="null">Tidak Ada Shift</option>
            <option v-for="shift in shifts" :key="shift.id" :value="shift.id">
              {{ shift.name }} ({{ shift.start_time }} - {{ shift.end_time }})
            </option>
          </select>
        </div>
        <div class="flex justify-end space-x-4">
          <BaseButton type="button" @click="closeModal" class="btn-outline-primary">
            Batal
          </BaseButton>
          <BaseButton type="submit">
            Simpan
          </BaseButton>
        </div>
      </form>
    </BaseModal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import axios from 'axios';
import { useToast } from 'vue-toastification';
import { useAuthStore } from '../../stores/auth';
import { RouterLink } from 'vue-router';
import BaseInput from '../ui/BaseInput.vue';
import BaseButton from '../ui/BaseButton.vue';
import BaseModal from '../ui/BaseModal.vue';

const employees = ref([]);
const shifts = ref([]);
const isModalOpen = ref(false);
const currentEmployee = ref({});
const searchTerm = ref('');
const editingEmployee = ref(false);
const toast = useToast();
const authStore = useAuthStore();

const fetchShifts = async () => {
  try {
    const response = await axios.get(`/api/shifts`);
    if (response.data && response.data.status === 'success') {
      shifts.value = response.data.data;
    } else {
      toast.error(response.data?.message || 'Failed to fetch shifts.');
    }
  } catch (error) {
    console.error('Error fetching shifts:', error);
    let message = 'Failed to fetch shifts.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.error(message);
  }
};

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
      employees.value = Array.isArray(response.data.data) ? response.data.data : [];
      
      if (response.data.data !== undefined && response.data.data !== null && !Array.isArray(response.data.data)) {
        toast.warning('Received non-array data for employees, treating as empty list.');
      }
    } else {
      console.log('Unexpected API response for employees:', response);
      employees.value = [];
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

onMounted(() => {
  fetchShifts();
});

watch(() => authStore.companyId, (newCompanyId) => {
  if (newCompanyId) {
    fetchEmployees();
  }
}, { immediate: true });

let searchTimeout = null;
watch(searchTerm, (newSearchTerm) => {
  clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    fetchEmployees();
  }, 300);
});

const openAddModal = () => {
  currentEmployee.value = { name: '', email: '', position: '', employee_id_number: '', shift_id: null };
  editingEmployee.value = false;
  isModalOpen.value = true;
};

const openEditModal = (employee) => {
  currentEmployee.value = { ...employee, shift_id: employee.shift_id || null };
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
        employee_id_number: currentEmployee.value.employee_id_number,
        shift_id: currentEmployee.value.shift_id,
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
</script>

</script>