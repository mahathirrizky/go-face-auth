<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-6">Manajemen Karyawan</h2>

    <Tabs v-model:value="activeTab">
      <TabList>
        <Tab :value="0">Daftar Karyawan</Tab>
        <Tab :value="1">Pending</Tab>
      </TabList>

      <TabPanels>
        <TabPanel :value="0">
          <BaseDataTable
            :data="employees"
            :columns="employeeColumns"
            :loading="isLoading"
            :totalRecords="employeesTotalRecords"
            :lazy="true"
            v-model:filters="employeesFilters"
            @page="onPage($event, 'active')"
            @filter="onFilter($event, 'active')"
            searchPlaceholder="Cari Karyawan..."
            editable="true"
            @save="onRowEditSave"
            @cancel="onRowEditCancel"
          >
            <template #header-actions>
              <div class="flex flex-wrap gap-3">
                <BaseButton @click="openAddModal">
                  <i class="pi pi-plus"></i> <span class="hidden sm:inline">Tambah Karyawan</span>
                </BaseButton>
                <BaseButton @click="openBulkImportModal" class="btn-secondary">
                  <i class="pi pi-file-excel"></i> <span class="hidden sm:inline">Import dari Excel</span>
                </BaseButton>
              </div>
            </template>

            <template #column-history="{ item }">
              <router-link :to="{ name: 'EmployeeAttendanceHistory', params: { employeeId: item.id } }" custom v-slot="{ navigate }">
                <BaseButton @click="navigate" role="link" class="btn-info btn-sm"><i class="pi pi-history"></i> <span class="hidden sm:inline">Riwayat Absensi</span></BaseButton>
              </router-link>
            </template>

            <template #editor-name="{ data, field }">
              <BaseInput v-model="data[field]" :id="`edit-${field}`" :name="field" />
            </template>
            <template #editor-email="{ data, field }">
              <BaseInput v-model="data[field]" :id="`edit-${field}`" :name="field" type="email" />
            </template>
            <template #editor-employee_id_number="{ data, field }">
              <BaseInput v-model="data[field]" :id="`edit-${field}`" :name="field" />
            </template>
            <template #editor-position="{ data, field }">
              <BaseInput v-model="data[field]" :id="`edit-${field}`" :name="field" />
            </template>
            <template #editor-shift_id="{ data, field }">
              <Dropdown
                v-model="data[field]"
                :options="shifts"
                optionLabel="name"
                optionValue="id"
                placeholder="Pilih Shift"
                class="w-full"
              />
            </template>

            <template #actions="{ item }">
              <div class="flex flex-wrap gap-2">
                <BaseButton @click="deleteEmployee(item.id)" class="btn-danger btn-sm">
                  <i class="pi pi-trash"></i> Hapus
                </BaseButton>
              </div>
            </template>
          </BaseDataTable>
        </TabPanel>
        <TabPanel :value="1">
            <BaseDataTable
                :data="pendingEmployees"
                :columns="pendingEmployeeColumns"
                :loading="isLoading"
                :totalRecords="pendingTotalRecords"
                :lazy="true"
                v-model:filters="pendingFilters"
                @page="onPage($event, 'pending')"
                @filter="onFilter($event, 'pending')"
                searchPlaceholder="Cari Karyawan..."
            >
                <template #header-actions>
                  <p class="text-text-muted">Karyawan yang belum mengatur kata sandi awal.</p>
                </template>
                <template #column-actions="{ item }">
                    <div class="flex flex-wrap gap-2">
                        <BaseButton @click="resendPasswordEmail(item.id)" class="btn-secondary btn-sm"><i class="pi pi-envelope"></i> <span class="hidden md:inline">Kirim Ulang Email</span></BaseButton>
                    </div>
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
          <BaseButton @click="downloadTemplate" class="btn-secondary">  <i class="pi pi-download"></i> Unduh Template Excel</BaseButton>  
        <div v-if="!hasMultipleShifts" class="bg-blue-100 border-l-4 border-blue-500 text-blue-700 p-4 mb-4" role="alert">
          <p class="font-bold">Penting:</p>
          <p>Perusahaan Anda saat ini hanya memiliki satu shift atau belum mengkonfigurasi shift. Pastikan Anda telah mengatur shift yang sesuai di halaman
            <router-link to="/dashboard/settings/shifts" class="underline font-bold text-blue-800 hover:text-blue-900">Manajemen Shift</router-link>
            untuk memastikan data absensi karyawan tercatat dengan benar.
          </p>
        </div>

       

         <div class="mb-4"> <label for="bulkFile" class="block text-text-muted text-sm font-bold  mb-2">Pilih File Excel:</label>   <FileUpload  name="bulkFile"  @uploader="uploadBulkFile"   :customUpload="true"  :multiple="false"  accept=".xlsx, .xls"  :maxFileSize="1000000" chooseLabel="Pilih File" uploadLabel="Unggah"   cancelLabel="Batal"  class="w-full"   > <template #empty> <p class="text-center text-text-muted">Seret dan lepas file di sini untuk mengunggah.</p>    </template>   </FileUpload> </div>  <div class="flex justify-end space-x-4">     <BaseButton type="button" @click="closeBulkImportModal" class="btn-outline-primary"> Batal  </BaseButton> </div> <div v-if="bulkImportResults" class="mt-6 p-4 rounded-lg" :class="bulkImportResults.failed_count > 0 ? 'bg-red-100 text-red-800' : 'bg-green-100 text-green-800'"> <h4 class="font-bold mb-2">Hasil Impor:</h4> <p>Total Diproses: {{ bulkImportResults.total_processed }}</p> <p>Berhasil: {{ bulkImportResults.success_count }}</p> <p>Gagal: {{ bulkImportResults.failed_count }}</p>  <div v-if="bulkImportResults.failed_count > 0" class="mt-4"> <h5 class="font-semibold">Detail Kegagalan:</h5> <ul class="list-disc list-inside">  <li v-for="(result, index) in bulkImportResults.results" :key="index">  {{ result.row_number || 'N/A' }}: {{ result.message }} </li>   </ul> </div>   </div>   </div>  </BaseModal>   </div>  
  <ConfirmDialog></ConfirmDialog>
</template>

<script setup>
import { ref, onMounted, watch,computed } from 'vue';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { useConfirm } from 'primevue/useconfirm';
import { useAuthStore } from '../../stores/auth';
import { FilterMatchMode } from '@primevue/core/api';
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
import ConfirmDialog from 'primevue/confirmdialog';

const toast = useToast();
const confirm = useConfirm();
const authStore = useAuthStore();

const employees = ref([]);
const pendingEmployees = ref([]);
const shifts = ref([]);

const isLoading = ref(false);
const activeTab = ref(0);

// State for active employees table
const employeesTotalRecords = ref(0);
const employeesLazyParams = ref({});
const employeesFilters = ref({ 'global': { value: null, matchMode: FilterMatchMode.CONTAINS } });

// State for pending employees table
const pendingTotalRecords = ref(0);
const pendingLazyParams = ref({});
const pendingFilters = ref({ 'global': { value: null, matchMode: FilterMatchMode.CONTAINS } });

const employeeColumns = ref([
    { field: 'name', header: 'Nama' },
    { field: 'email', header: 'Email' },
    { field: 'employee_id_number', header: 'Nomor ID' },
    { field: 'position', header: 'Jabatan' },
    { field: 'history', header: 'Riwayat', sortable: false, editable: false }
]);

const pendingEmployeeColumns = ref([
    { field: 'name', header: 'Nama' },
    { field: 'email', header: 'Email' },
    { field: 'actions', header: 'Aksi', sortable: false }
]);

const fetchEmployees = async () => {
  if (!authStore.companyId) return;
  isLoading.value = true;
  try {
    const params = {
      page: employeesLazyParams.value.page + 1,
      limit: employeesLazyParams.value.rows,
      search: employeesFilters.value.global.value || ''
    };
    const response = await axios.get(`/api/companies/${authStore.companyId}/employees`, { params });
    if (response.data && response.data.status === 'success') {
      employees.value = response.data.data.items;
      employeesTotalRecords.value = response.data.data.total_records;
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal mengambil data karyawan.', life: 3000 });
    }
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Terjadi kesalahan saat mengambil data karyawan.', life: 3000 });
  } finally {
    isLoading.value = false;
  }
};

const fetchPendingEmployees = async () => {
  if (!authStore.companyId) return;
  isLoading.value = true;
  try {
    const params = {
      page: pendingLazyParams.value.page + 1,
      limit: pendingLazyParams.value.rows,
      search: pendingFilters.value.global.value || ''
    };
    const response = await axios.get(`/api/companies/${authStore.companyId}/employees/pending`, { params });
    if (response.data && response.data.status === 'success') {
      pendingEmployees.value = response.data.data.items;
      pendingTotalRecords.value = response.data.data.total_records;
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal mengambil data karyawan pending.', life: 3000 });
    }
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Terjadi kesalahan saat mengambil data karyawan pending.', life: 3000 });
  } finally {
    isLoading.value = false;
  }
};

const onPage = (event, type) => {
  if (type === 'active') {
    employeesLazyParams.value = event;
    fetchEmployees();
  } else if (type === 'pending') {
    pendingLazyParams.value = event;
    fetchPendingEmployees();
  }
};

const onFilter = (event, type) => {
  if (type === 'active') {
    // PrimeVue updates filters via v-model, just need to trigger fetch
    fetchEmployees();
  } else if (type === 'pending') {
    fetchPendingEmployees();
  }
};

watch(activeTab, (newTab) => {
  if (newTab === 0) {
    fetchEmployees();
  } else if (newTab === 1) {
    fetchPendingEmployees();
  }
});

const onRowEditSave = async (event) => {
  let { newData, index } = event;
  try {
    const response = await axios.put(`/api/employees/${newData.id}`, newData);
    if (response.data && response.data.status === 'success') {
      employees.value[index] = newData;
      toast.add({ severity: 'success', summary: 'Success', detail: 'Karyawan berhasil diperbarui!', life: 3000 });
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data?.message || 'Gagal memperbarui karyawan.', life: 3000 });
    }
  } catch (error) {
    let message = 'Terjadi kesalahan saat memperbarui karyawan.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  }
};

const onRowEditCancel = (event) => {
  toast.add({ severity: 'info', summary: 'Dibatalkan', detail: 'Perubahan dibatalkan', life: 3000 });
};

const fetchShifts = async () => {
  try {
    const response = await axios.get(`/api/shifts`);
    if (response.data && response.data.status === 'success') {
      shifts.value = response.data.data;
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data?.message || 'Failed to fetch shifts.', life: 3000 });
    }
  } catch (error) {
    console.error('Error fetching shifts:', error);
    let message = 'Failed to fetch shifts.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  }
};





onMounted(() => {
  fetchShifts();
  employeesLazyParams.value = { first: 0, rows: 10, page: 0 };
  pendingLazyParams.value = { first: 0, rows: 10, page: 0 };
  fetchEmployees(); // Fetch initial data for the first tab
});

const isModalOpen = ref(false);
const currentEmployee = ref({});
const editingEmployee = ref(false);

const openAddModal = () => {
  currentEmployee.value = { name: '', email: '', position: '', employee_id_number: '', shift_id: null };
  editingEmployee.value = false;
  isModalOpen.value = true;
};

const closeModal = () => {
  isModalOpen.value = false;
  currentEmployee.value = {};
};

const deleteEmployee = (id) => {
  confirm.require({
    message: 'Apakah Anda yakin ingin menghapus karyawan ini?',
    header: 'Konfirmasi Hapus Karyawan',
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      try {
        await axios.delete(`/api/employees/${id}`);
        toast.add({ severity: 'success', summary: 'Success', detail: 'Employee deleted successfully!', life: 3000 });
        fetchEmployees(); // Always refresh the active list
      } catch (error) {
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

const resendPasswordEmail = async (employeeId) => {
  try {
    const response = await axios.post(`/api/employees/${employeeId}/resend-password-email`);
    if (response.data && response.data.status === 'success') {
      toast.add({ severity: 'success', summary: 'Success', detail: 'Email pengaturan kata sandi berhasil dikirim ulang!', life: 3000 });
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data?.message || 'Gagal mengirim ulang email.', life: 3000 });
    }
  } catch (error) {
    let message = 'Gagal mengirim ulang email.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  }
  fetchPendingEmployees(); // Refresh pending list
};

// Bulk Import Modal related state and functions
const isBulkImportModalOpen = ref(false);
const bulkImportResults = ref(null);

const openBulkImportModal = () => {
  isBulkImportModalOpen.value = true;
  bulkImportResults.value = null;
};

const closeBulkImportModal = () => {
  isBulkImportModalOpen.value = false;
};

const uploadBulkFile = async (event) => {
  const file = event.files[0];
  if (!file) return;
  const formData = new FormData();
  formData.append('file', file);
  try {
    const response = await axios.post(`/api/employees/bulk`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
    if (response.data && response.data.status === 'success') {
      bulkImportResults.value = response.data.data;
      toast.add({ severity: 'success', summary: 'Success', detail: 'Impor massal selesai.', life: 3000 });
      fetchEmployees();
    } else {
      bulkImportResults.value = response.data.data;
      toast.add({ severity: 'error', summary: 'Error', detail: response.data?.message || 'Impor massal gagal.', life: 3000 });
    }
  } catch (error) {
    let message = 'Terjadi kesalahan saat mengunggah file.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  }
};

const downloadTemplate = async () => {
  try {
    const response = await axios.get(`/api/employees/template`, { responseType: 'blob' });
    const url = window.URL.createObjectURL(new Blob([response.data]));
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', 'employee_template.xlsx');
    document.body.appendChild(link);
    link.click();
    link.remove();
    window.URL.revokeObjectURL(url);
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal mengunduh template.', life: 3000 });
  }
};

const hasMultipleShifts = computed(() => shifts.value.length > 1);


</script>