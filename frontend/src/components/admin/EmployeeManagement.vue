<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-6">Manajemen Karyawan</h2>

    <Tabs value="0">
      <TabList>
        <Tab value="0">Daftar Karyawan</Tab>
        <Tab value="1">Pending</Tab>
      </TabList>

      <TabPanels>
        <TabPanel value="0">
          <BaseDataTable
            :data="employees"
            :columns="employeeColumns"
            :loading="isLoading"
            :globalFilterFields="['name', 'email', 'employee_id_number', 'position']"
            searchPlaceholder="Cari Karyawan..."
          >
            <template #header-actions>
              <BaseButton @click="openAddModal">
                Tambah Karyawan
              </BaseButton>
              <BaseButton @click="openBulkImportModal" class="btn-secondary">
                Import dari Excel
              </BaseButton>
            </template>

            <template #column-history="{ item }">
              <router-link :to="{ name: 'EmployeeAttendanceHistory', params: { employeeId: item.id } }" custom v-slot="{ navigate }">
                <BaseButton @click="navigate" role="link" class="btn-info btn-sm">Riwayat Absensi</BaseButton>
              </router-link>
            </template>

            <template #column-actions="{ item }">
              <BaseButton @click="openEditModal(item)" class="text-accent hover:opacity-80 mr-3">
                <i class="pi pi-pencil"></i> Edit
              </BaseButton>
              <BaseButton @click="deleteEmployee(item.id)" class="text-danger hover:opacity-80">
                <i class="pi pi-trash"></i> Hapus
              </BaseButton>
            </template>
          </BaseDataTable>
        </TabPanel>
        <TabPanel value="1">
            <BaseDataTable
                :data="pendingEmployees"
                :columns="pendingEmployeeColumns"
                :loading="isLoading"
                searchPlaceholder="Cari Karyawan..."
            >
                <template #header>
                  <div class="flex justify-end">
                    <p class="text-text-muted">Karyawan yang belum mengatur kata sandi awal.</p>
                  </div>
                </template>
                <template #column-actions="{ item }">
                    <BaseButton @click="resendPasswordEmail(item.id)" class="btn-secondary btn-sm">Kirim Ulang Email</BaseButton>
                </template>
            </BaseDataTable>
        </TabPanel>
      </TabPanels>
    </Tabs>

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
          <Dropdown
            id="shift"
            v-model="currentEmployee.shift_id"
            :options="shifts"
            optionLabel="name"
            optionValue="id"
            placeholder="Pilih Shift"
            class="w-full"
          />
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
            <i class="pi pi-download"></i> Unduh Template Excel
          </BaseButton>
        </div>

        <div class="mb-4">
          <label for="bulkFile" class="block text-text-muted text-sm font-bold mb-2">Pilih File Excel:</label>
          <FileUpload
            name="bulkFile"
            @uploader="uploadBulkFile"
            :customUpload="true"
            :multiple="false"
            accept=".xlsx, .xls"
            :maxFileSize="1000000"
            chooseLabel="Pilih File"
            uploadLabel="Unggah"
            cancelLabel="Batal"
            class="w-full"
          >
            <template #empty>
              <p class="text-center text-text-muted">Seret dan lepas file di sini untuk mengunggah.</p>
            </template>
          </FileUpload>
        </div>

        <div class="flex justify-end space-x-4">
          <BaseButton type="button" @click="closeBulkImportModal" class="btn-outline-primary">
            Batal
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
import { useToast } from 'primevue/usetoast';
import { useConfirm } from 'primevue/useconfirm';

import { useAuthStore } from '../../stores/auth';
import { RouterLink } from 'vue-router';
import BaseInput from '../ui/BaseInput.vue';
import BaseButton from '../ui/BaseButton.vue';
import BaseModal from '../ui/BaseModal.vue';
import BaseDataTable from '../ui/BaseDataTable.vue';
import Tabs from 'primevue/tabs';
import Tab from 'primevue/tab';
import TabList from 'primevue/tablist';
import TabPanels from 'primevue/tabpanels';
import TabPanel from 'primevue/tabpanel';
import Dropdown from 'primevue/dropdown';
import FileUpload from 'primevue/fileupload';

const employees = ref([]);
const pendingEmployees = ref([]); // New ref for pending employees
const shifts = ref([]);
const isModalOpen = ref(false);
const currentEmployee = ref({});
const editingEmployee = ref(false);
const toast = useToast();
const authStore = useAuthStore();
const selectedTab = ref(0); // New ref for tab selection, 0 for 'all', 1 for 'pending'
const isLoading = ref(false);
const confirm = useConfirm();

const employeeColumns = ref([
    { field: 'name', header: 'Nama' },
    { field: 'email', header: 'Email' },
    { field: 'employee_id_number', header: 'Nomor ID' },
    { field: 'position', header: 'Jabatan' },
    { field: 'history', header: 'Riwayat' },
    { field: 'actions', header: 'Aksi' }
]);

const pendingEmployeeColumns = ref([
    { field: 'name', header: 'Nama' },
    { field: 'email', header: 'Email' },
    { field: 'actions', header: 'Aksi' }
]);

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
    toast.add({ severity: 'error', summary: 'Error', detail: 'Company ID not available. Cannot fetch employees.', life: 3000 });
    return;
  }
  isLoading.value = true;
  try {
    const url = `/api/companies/${authStore.companyId}/employees`;
    
    const response = await axios.get(url);
    if (response.data && response.data.status === 'success') {
      console.log('Fetched employees:', response.data.data); // Log fetched data
      employees.value = Array.isArray(response.data.data) ? response.data.data : [];
      
      if (response.data.data !== undefined && response.data.data !== null && !Array.isArray(response.data.data)) {
        toast.add({ severity: 'warning', summary: 'Warning', detail: 'Received non-array data for employees, treating as empty list.', life: 3000 });
      }
    } else {
      console.log('Unexpected API response for employees:', response);
      employees.value = [];
      toast.add({ severity: 'error', summary: 'Error', detail: response.data?.message || 'Failed to fetch employees due to an unexpected response format.', life: 3000 });
    }
  } catch (error) {
    console.error('Error fetching employees:', error);
    let message = 'Failed to fetch employees.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.error(message);
  } finally {
    isLoading.value = false;
  }
};

const fetchPendingEmployees = async () => {
  if (!authStore.companyId) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Company ID not available. Cannot fetch pending employees.', life: 3000 });
    return;
  }
  isLoading.value = true;
  try {
    // Assuming a new backend endpoint for pending employees
    // This endpoint should return employees who have a password reset token but no password set
    const response = await axios.get(`/api/companies/${authStore.companyId}/employees/pending`);
    if (response.data && response.data.status === 'success') {
      pendingEmployees.value = Array.isArray(response.data.data) ? response.data.data : [];
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data?.message || 'Failed to fetch pending employees.', life: 3000 });
    }
  } catch (error) {
    console.error('Error fetching pending employees:', error);
    let message = 'Failed to fetch pending employees.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  } finally {
    isLoading.value = false;
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
        toast.add({ severity: 'success', summary: 'Success', detail: 'Email pengaturan kata sandi berhasil dikirim ulang!', life: 3000 });
      } else {
        toast.add({ severity: 'error', summary: 'Error', detail: response.data?.message || 'Gagal mengirim ulang email pengaturan kata sandi.', life: 3000 });
      }
    } catch (error) {
    console.error('Error resending password email:', error);
    let message = 'Gagal mengirim ulang email pengaturan kata sandi.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  }
};
};

onMounted(() => {
  fetchShifts();
  // Initial fetch for the default tab
  if (selectedTab.value === 0) { // 'all' tab
    fetchEmployees();
  } else if (selectedTab.value === 1) { // 'pending' tab
    fetchPendingEmployees();
  }
});

watch(() => authStore.companyId, (newCompanyId) => {
  if (newCompanyId) {
    if (selectedTab.value === 0) {
      fetchEmployees();
    } else if (selectedTab.value === 1) {
      fetchPendingEmployees();
    }
  }
}, { immediate: true });

watch(selectedTab, (newTab) => {
  if (newTab === 0) {
    fetchEmployees();
  } else if (newTab === 1) {
    fetchPendingEmployees();
  }
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
      toast.add({ severity: 'success', summary: 'Success', detail: response.data.message || 'Employee updated successfully!', life: 3000 });
    } else {
      console.log('Creating employee:', currentEmployee.value);
      const response = await axios.post(`/api/employees`, {
        name: currentEmployee.value.name,
        email: currentEmployee.value.email,
        position: currentEmployee.value.position,
        employee_id_number: currentEmployee.value.employee_id_number,
        shift_id: currentEmployee.value.shift_id,
      });
      toast.add({ severity: 'success', summary: 'Success', detail: response.data.message || 'Employee created successfully. An email with initial password setup link has been sent.', life: 3000 });
    }
    closeModal();
    // Refresh the correct tab after saving
    if (selectedTab.value === 0) {
      fetchEmployees();
    } else if (selectedTab.value === 1) {
      fetchPendingEmployees();
    }
  } catch (error) {
    console.error('Error saving employee:', error);
    let message = 'Failed to save employee.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  }
};

const deleteEmployee = (id) => {
  confirm.require({
    message: 'Apakah Anda yakin ingin menghapus karyawan ini?',
    header: 'Konfirmasi Hapus Karyawan',
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      if (!authStore.companyId) {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Company ID not available. Cannot delete employee.', life: 3000 });
        return;
      }
      try {
        const response = await axios.delete(`/api/employees/${id}`);
        toast.add({ severity: 'success', summary: 'Success', detail: 'Employee deleted successfully!', life: 3000 });
        if (selectedTab.value === 0) {
          fetchEmployees();
        } else if (selectedTab.value === 1) {
          fetchPendingEmployees();
        }
      } catch (error) {
        console.error('Error deleting employee:', error);
        let message = 'Failed to delete employee.';
        if (error.response && error.response.data && error.response.data.message) {
          message = error.response.data.message;
        }
        toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
      }
    },
    reject: () => {
      toast.add({ severity: 'info', summary: 'Dibatalkan', detail: 'Penghapusan karyawan dibatalkan', life: 3000 });
    }
  });
};

const isBulkImportModalOpen = ref(false);
const bulkImportResults = ref(null);

const openBulkImportModal = () => {
  isBulkImportModalOpen.value = true;
  bulkImportResults.value = null; // Clear previous results
};

const closeBulkImportModal = () => {
  isBulkImportModalOpen.value = false;
};

const uploadBulkFile = async (event) => {
  const file = event.files[0];
  if (!file) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Silakan pilih file Excel untuk diunggah.', life: 3000 });
    return;
  }

  const formData = new FormData();
  formData.append('file', file);

  try {
    const response = await axios.post(`/api/employees/bulk`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });

    if (response.data && response.data.status === 'success') {
      bulkImportResults.value = response.data.data;
      toast.add({ severity: 'success', summary: 'Success', detail: response.data.message || 'Impor massal selesai.', life: 3000 });
      // Refresh employee list after successful import
      if (selectedTab.value === 0) {
        fetchEmployees();
      }
    } else {
      bulkImportResults.value = response.data.data; // Display errors if any
      toast.add({ severity: 'error', summary: 'Error', detail: response.data?.message || 'Impor massal gagal.', life: 3000 });
    }
  } catch (error) {
    console.error('Error uploading bulk file:', error);
    let message = 'Terjadi kesalahan saat mengunggah file.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    bulkImportResults.value = { failed_count: 1, results: [{ message: message }] }; // Display generic error
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  }
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
    toast.add({ severity: 'success', summary: 'Success', detail: 'Template Excel berhasil diunduh!', life: 3000 });
  } catch (error) {
    console.error('Error downloading template:', error);
    let message = 'Gagal mengunduh template Excel.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  }
};

const hasMultipleShifts = computed(() => shifts.value.length > 1);

</script>
