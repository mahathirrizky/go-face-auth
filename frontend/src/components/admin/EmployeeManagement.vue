<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-6">Manajemen Karyawan</h2>

    <!-- Tab Navigation -->
    <div class="mb-6 border-b border-bg-base">
      <nav class="-mb-px flex space-x-8" aria-label="Tabs">
        <button
          @click="selectedTab = 'all'"
          :class="[
            'whitespace-nowrap py-3 px-1 border-b-2 font-medium text-sm',
            selectedTab === 'all'
              ? 'border-secondary text-secondary'
              : 'border-transparent text-text-muted hover:text-text-base hover:border-gray-300',
          ]"
        >
          Daftar Karyawan
        </button>
        <button
          @click="selectedTab = 'pending'"
          :class="[
            'whitespace-nowrap py-3 px-1 border-b-2 font-medium text-sm',
            selectedTab === 'pending'
              ? 'border-secondary text-secondary'
              : 'border-transparent text-text-muted hover:text-text-base hover:border-gray-300',
          ]"
        >
          Pending
        </button>
      </nav>
    </div>

    <!-- Tab Content: Daftar Karyawan -->
    <div v-if="selectedTab === 'all'">
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
        <BaseButton @click="openBulkImportModal" class="w-full md:w-auto btn-secondary ml-2">
          Import dari Excel
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
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Riwayat</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-bg-base">
            <tr v-for="employee in employees" :key="employee.id">
              <td class="px-6 py-4 whitespace-nowrap text-text-base">{{ employee.name }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ employee.email }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ employee.employee_id_number }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ employee.position }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-center text-sm font-medium">
                <router-link :to="{ name: 'EmployeeAttendanceHistory', params: { employeeId: employee.id } }" custom v-slot="{ navigate }">
                  <BaseButton @click="navigate" role="link" class="btn-info btn-sm">Riwayat Absensi</BaseButton>
                </router-link>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <BaseButton @click="openEditModal(employee)" class="text-accent hover:text-secondary mr-3">
                  <i class="fas fa-edit"></i> Edit
                </BaseButton>
                <BaseButton @click="deleteEmployee(employee.id)" class="text-danger hover:opacity-80">
                  <i class="fas fa-trash-alt"></i> Hapus
                </BaseButton>
              </td>
            </tr>
            <tr v-if="employees.length === 0">
              <td colspan="6" class="px-6 py-4 text-center text-text-muted">Tidak ada data karyawan.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Tab Content: Pending -->
    <div v-if="selectedTab === 'pending'">
      <div class="bg-bg-muted p-4 rounded-lg shadow-md mb-6 flex justify-end">
        <p class="text-text-muted">Karyawan yang belum mengatur kata sandi awal.</p>
      </div>
      <div class="overflow-x-auto bg-bg-muted rounded-lg shadow-md">
        <table class="min-w-full divide-y divide-bg-base">
          <thead class="bg-primary">
            <tr>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Nama</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Email</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-bg-base">
            <tr v-for="employee in pendingEmployees" :key="employee.id">
              <td class="px-6 py-4 whitespace-nowrap text-text-base">{{ employee.name }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ employee.email }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <BaseButton @click="resendPasswordEmail(employee.id)" class="btn-secondary btn-sm">Kirim Ulang Email</BaseButton>
              </td>
            </tr>
            <tr v-if="pendingEmployees.length === 0">
              <td colspan="3" class="px-6 py-4 text-center text-text-muted">Tidak ada karyawan pending.</td>
            </tr>
          </tbody>
        </table>
      </div>
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

    <!-- Confirmation Modal for Deletion -->
    <BaseModal :isOpen="isConfirmModalOpen" @close="cancelDelete" title="Konfirmasi Hapus Karyawan" maxWidth="sm">
      <div class="text-center p-4">
        <p class="text-lg text-text-base mb-6">Apakah Anda yakin ingin menghapus karyawan ini?</p>
        <div class="flex justify-center space-x-4">
          <BaseButton @click="cancelDelete" class="btn-outline-primary">
            Batal
          </BaseButton>
          <BaseButton @click="confirmDelete" class="btn-danger">
            Ya, Hapus
          </BaseButton>
        </div>
      </div>
    </BaseModal>

    <!-- Bulk Import Modal -->
    <BaseModal :isOpen="isBulkImportModalOpen" @close="closeBulkImportModal" title="Import Karyawan dari Excel" maxWidth="lg">
      <div class="p-4">
        <p class="text-text-muted mb-4">
          Gunakan fitur ini untuk menambahkan banyak karyawan sekaligus. Unduh template Excel, isi data karyawan, lalu unggah kembali file tersebut.
        </p>

        <div v-if="hasMultipleShifts" class="bg-blue-100 border-l-4 border-blue-500 text-blue-700 p-4 mb-4" role="alert">
          <p class="font-bold">Penting:</p>
          <p>Untuk memastikan nama shift sesuai, silakan atur shift Anda di halaman
            <router-link to="/dashboard/settings/shifts" class="underline font-bold text-blue-800 hover:text-blue-900">Manajemen Shift</router-link>
            jika Anda memiliki lebih dari satu jenis shift.
          </p>
        </div>

        <div class="mb-4">
          <BaseButton @click="downloadTemplate" class="btn-secondary">
            <i class="fas fa-download"></i> Unduh Template Excel
          </BaseButton>
        </div>

        <div class="mb-4">
          <label for="bulkFile" class="block text-text-muted text-sm font-bold mb-2">Pilih File Excel:</label>
          <input
            type="file"
            id="bulkFile"
            @change="handleBulkFileChange"
            accept=".xlsx, .xls"
            class="block w-full text-sm text-text-base file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-primary file:text-white hover:file:bg-secondary hover:file:text-primary"
          />
        </div>

        <div class="flex justify-end space-x-4">
          <BaseButton type="button" @click="closeBulkImportModal" class="btn-outline-primary">
            Batal
          </BaseButton>
          <BaseButton @click="uploadBulkFile" class="btn-primary">
            <i class="fas fa-upload"></i> Unggah
          </BaseButton>
        </div>

        <div v-if="bulkImportResults" class="mt-6 p-4 rounded-lg" :class="bulkImportResults.failed_count > 0 ? 'bg-red-100 text-red-800' : 'bg-green-100 text-green-800'">
          <h4 class="font-bold mb-2">Hasil Impor:</h4>
          <p>Total Diproses: {{ bulkImportResults.total_processed }}</p>
          <p>Berhasil: {{ bulkImportResults.success_count }}</p>
          <p>Gagal: {{ bulkImportResults.failed_count }}</p>
          <div v-if="bulkImportResults.failed_count > 0" class="mt-4">
            <h5 class="font-semibold">Detail Kegagalan:</h5>
            <ul class="list-disc list-inside">
              <li v-for="(result, index) in bulkImportResults.results" :key="index">
                Baris {{ result.row_number || 'N/A' }}: {{ result.message }}
              </li>
            </ul>
          </div>
        </div>
      </div>
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
const pendingEmployees = ref([]); // New ref for pending employees
const shifts = ref([]);
const isModalOpen = ref(false);
const currentEmployee = ref({});
const searchTerm = ref('');
const editingEmployee = ref(false);
const toast = useToast();
const authStore = useAuthStore();
const selectedTab = ref('all'); // New ref for tab selection

const isConfirmModalOpen = ref(false);
const employeeToDeleteId = ref(null);

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
      console.log('Fetched employees:', response.data.data); // Log fetched data
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

const fetchPendingEmployees = async () => {
  if (!authStore.companyId) {
    toast.error('Company ID not available. Cannot fetch pending employees.');
    return;
  }
  try {
    // Assuming a new backend endpoint for pending employees
    // This endpoint should return employees who have a password reset token but no password set
    const response = await axios.get(`/api/companies/${authStore.companyId}/employees/pending`);
    if (response.data && response.data.status === 'success') {
      pendingEmployees.value = Array.isArray(response.data.data) ? response.data.data : [];
    } else {
      toast.error(response.data?.message || 'Failed to fetch pending employees.');
    }
  } catch (error) {
    console.error('Error fetching pending employees:', error);
    let message = 'Failed to fetch pending employees.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.error(message);
  }
};

const resendPasswordEmail = async (employeeId) => {
  if (!authStore.companyId) {
    toast.error('Company ID not available. Cannot resend email.');
    return;
  }
  if (confirm('Apakah Anda yakin ingin mengirim ulang email pengaturan kata sandi untuk karyawan ini?')) {
    try {
      // Assuming a new backend endpoint to resend password email
      const response = await axios.post(`/api/employees/${employeeId}/resend-password-email`);
      if (response.data && response.data.status === 'success') {
        toast.success('Email pengaturan kata sandi berhasil dikirim ulang!');
      } else {
        toast.error(response.data?.message || 'Gagal mengirim ulang email pengaturan kata sandi.');
      }
    } catch (error) {
      console.error('Error resending password email:', error);
      let message = 'Gagal mengirim ulang email pengaturan kata sandi.';
      if (error.response && error.response.data && error.response.data.message) {
        message = error.response.data.message;
      }
      toast.error(message);
    }
  }
};

onMounted(() => {
  fetchShifts();
  // Initial fetch for the default tab
  if (selectedTab.value === 'all') {
    fetchEmployees();
  } else if (selectedTab.value === 'pending') {
    fetchPendingEmployees();
  }
});

watch(() => authStore.companyId, (newCompanyId) => {
  if (newCompanyId) {
    if (selectedTab.value === 'all') {
      fetchEmployees();
    } else if (selectedTab.value === 'pending') {
      fetchPendingEmployees();
    }
  }
}, { immediate: true });

watch(selectedTab, (newTab) => {
  if (newTab === 'all') {
    fetchEmployees();
  } else if (newTab === 'pending') {
    fetchPendingEmployees();
  }
});

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
    // Refresh the correct tab after saving
    if (selectedTab.value === 'all') {
      fetchEmployees();
    } else if (selectedTab.value === 'pending') {
      fetchPendingEmployees();
    }
  } catch (error) {
    console.error('Error saving employee:', error);
    let message = 'Failed to save employee.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.error(message);
  }
};

const deleteEmployee = (id) => {
  employeeToDeleteId.value = id;
  isConfirmModalOpen.value = true;
};

const confirmDelete = async () => {
  if (!authStore.companyId) {
    toast.error('Company ID not available. Cannot delete employee.');
    return;
  }
  try {
    const response = await axios.delete(`/api/employees/${employeeToDeleteId.value}`);
    toast.success(response.data.message || 'Employee deleted successfully!');
    // Refresh the correct tab after deleting
    if (selectedTab.value === 'all') {
      fetchEmployees();
    } else if (selectedTab.value === 'pending') {
      fetchPendingEmployees();
    }
    cancelDelete(); // Close the confirmation modal
  } catch (error) {
    console.error('Error deleting employee:', error);
    let message = 'Failed to delete employee.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.error(message);
  }
};

const cancelDelete = () => {
  isConfirmModalOpen.value = false;
  employeeToDeleteId.value = null;
};

const isBulkImportModalOpen = ref(false);
const bulkFile = ref(null);
const bulkImportResults = ref(null);

const openBulkImportModal = () => {
  isBulkImportModalOpen.value = true;
  bulkFile.value = null; // Reset file input
  bulkImportResults.value = null; // Clear previous results
};

const closeBulkImportModal = () => {
  isBulkImportModalOpen.value = false;
};

const handleBulkFileChange = (event) => {
  bulkFile.value = event.target.files[0];
};

const downloadTemplate = async () => {
  try {
    const response = await axios.get(`/api/employees/template`, {
      responseType: 'blob', // Important for downloading files
    });

    const url = window.URL.createObjectURL(new Blob([response.data]));
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', 'employee_template.xlsx'); // Or whatever name you want
    document.body.appendChild(link);
    link.click();
    link.remove();
    window.URL.revokeObjectURL(url);
    toast.success('Template Excel berhasil diunduh!');
  } catch (error) {
    console.error('Error downloading template:', error);
    let message = 'Gagal mengunduh template Excel.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.error(message);
  }
};

const uploadBulkFile = async () => {
  if (!bulkFile.value) {
    toast.error('Silakan pilih file Excel untuk diunggah.');
    return;
  }

  const formData = new FormData();
  formData.append('file', bulkFile.value);

  try {
    const response = await axios.post(`/api/employees/bulk`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });

    if (response.data && response.data.status === 'success') {
      bulkImportResults.value = response.data.data;
      toast.success(response.data.message || 'Impor massal selesai.');
      // Refresh employee list after successful import
      if (selectedTab.value === 'all') {
        fetchEmployees();
      } else if (selectedTab.value === 'pending') {
        fetchPendingEmployees();
      }
    } else {
      bulkImportResults.value = response.data.data; // Display errors if any
      toast.error(response.data?.message || 'Impor massal gagal.');
    }
  } catch (error) {
    console.error('Error uploading bulk file:', error);
    let message = 'Terjadi kesalahan saat mengunggah file.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    bulkImportResults.value = { failed_count: 1, results: [{ message: message }] }; // Display generic error
    toast.error(message);
  }
};

const hasMultipleShifts = computed(() => shifts.value.length > 1);

</script>
