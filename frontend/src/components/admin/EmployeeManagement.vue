<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <Toast />
    <ConfirmDialog />
    <h2 class="text-2xl font-bold text-text-base mb-6">Manajemen Karyawan</h2>

    <Tabs v-model:value="activeTab">
      <TabList>
        <Tab value="active">Daftar Karyawan</Tab>
        <Tab value="pending">Pending</Tab>
      </TabList>
      <TabPanels>
        <TabPanel value="active">
          <DataTable
            :value="employees"
            :loading="isLoading"
            :totalRecords="employeesTotalRecords"
            :lazy="true"
            v-model:filters="employeesFilters"
            @page="onPage($event, 'active')"
            @filter="onFilter($event, 'active')"
            paginator :rows="10" :rowsPerPageOptions="[10, 25, 50]"
            editMode="cell"
            @cell-edit-complete="onCellEditComplete"
            dataKey="id"
            :globalFilterFields="['name', 'email', 'employee_id_number', 'position']"
          >
            <template #header>
              <div class="flex flex-wrap justify-between items-center gap-4 spa">
                  <IconField iconPosition="left">
                      <InputIcon class="pi pi-search"></InputIcon>
                      <InputText v-model="employeesFilters['global'].value" placeholder="Cari Karyawan..." @keydown.enter="onFilter($event, 'active')"/>
                  </IconField>
                  <div class="flex flex-wrap gap-2 space-x-2">
                      <Button @click="openAddModal" icon="pi pi-plus" label="Tambah Karyawan" />
                      <Button @click="openBulkImportModal" icon="pi pi-file-excel" label="Import dari Excel" class="p-button-secondary" />
                  </div>
              </div>
            </template>
            
            <Column field="name" header="Nama" :sortable="true" style="min-width: 12rem">
              <template #editor="{ data, field }">
                  <InputText v-model="data[field]" autofocus />
              </template>
            </Column>
            <Column field="email" header="Email" :sortable="true" style="min-width: 12rem">
               <template #editor="{ data, field }">
                  <InputText v-model="data[field]" autofocus />
              </template>
            </Column>
            <Column field="employee_id_number" header="Nomor ID" :sortable="true" style="min-width: 10rem">
               <template #editor="{ data, field }">
                  <InputText v-model="data[field]" autofocus />
              </template>
            </Column>
            <Column field="position" header="Jabatan" :sortable="true" style="min-width: 10rem">
               <template #editor="{ data, field }">
                  <InputText v-model="data[field]" autofocus />
              </template>
            </Column>
            <Column field="shift_id" header="Shift" :sortable="true" style="min-width: 10rem">
              <template #body="{ data }">{{ data.Shift ? data.Shift.Name : 'N/A' }}</template>
              <template #editor="{ data, field }">
                  <Select v-model="data[field]" :options="shifts" optionLabel="name" optionValue="id" placeholder="Pilih Shift" class="w-full" />
              </template>
            </Column>
            <Column header="Riwayat" style="min-width: 10rem">
              <template #body="{ data }">
                <router-link :to="{ name: 'EmployeeAttendanceHistory', params: { employeeId: data.id } }" custom v-slot="{ navigate }">
                  <Button @click="navigate" role="link" icon="pi pi-history" class="p-button-info p-button-sm" label="Riwayat" />
                </router-link>
              </template>
            </Column>
            <Column header="Aksi" style="min-width: 8rem">
              <template #body="{ data }">
                <Button @click="deleteEmployee(data.id)" icon="pi pi-trash" class="p-button-danger p-button-sm" label="Hapus" />
              </template>
            </Column>
          </DataTable>
        </TabPanel>
        <TabPanel value="pending">
          <DataTable
            :value="pendingEmployees"
            :loading="isLoading"
            :totalRecords="pendingTotalRecords"
            :lazy="true"
            v-model:filters="pendingFilters"
            @page="onPage($event, 'pending')"
            @filter="onFilter($event, 'pending')"
            paginator :rows="10" :rowsPerPageOptions="[10, 25, 50]"
            dataKey="id"
            :globalFilterFields="['name', 'email']"
          >
            <template #header>
               <div class="flex flex-wrap justify-between items-center gap-4">
                  <p class="text-text-muted m-0">Karyawan yang belum mengatur kata sandi awal.</p>
                  <IconField iconPosition="left">
                      <InputIcon class="pi pi-search"></InputIcon>
                      <InputText v-model="pendingFilters['global'].value" placeholder="Cari Karyawan..." @keydown.enter="onFilter($event, 'pending')"/>
                  </IconField>
              </div>
            </template>
            <Column field="name" header="Nama" :sortable="true"></Column>
            <Column field="email" header="Email" :sortable="true"></Column>
            <Column header="Aksi">
              <template #body="{ data }">
                <Button @click="resendPasswordEmail(data.id)" icon="pi pi-envelope" class="p-button-secondary p-button-sm" label="Kirim Ulang Email" :loading="isResending" />
              </template>
            </Column>
          </DataTable>
        </TabPanel>
      </TabPanels>
    </Tabs>

    <Dialog v-model:visible="isModalOpen" :header="editingEmployee ? 'Edit Karyawan' : 'Tambah Karyawan'" :modal="true" class="w-full max-w-md">
      <form @submit.prevent="saveEmployee" class="p-fluid mt-4">
        <div class="field mb-4">
          <FloatLabel variant="on">

            <label for="name">Nama :</label>
          <InputText id="name" v-model="currentEmployee.name" required fluid/>
          </FloatLabel>
        </div>
        <div class="field mb-4">
           <FloatLabel variant="on">

             <label for="email">Email :</label>
             <InputText id="email" v-model="currentEmployee.email" type="email" required  fluid/>
            </FloatLabel>
            </div>
        <div class="field mb-4">
           <FloatLabel variant="on">

             <label for="employeeIdNumber">Nomor ID Karyawan :</label>
             <InputText id="employeeIdNumber" v-model="currentEmployee.employee_id_number" required fluid/>
            </FloatLabel>
        </div>
        <div class="field mb-4">
           <FloatLabel variant="on">

             <label for="position">Jabatan :</label>
             <InputText id="position"  v-model="currentEmployee.position" required fluid/>
            </FloatLabel>
        </div>
        <div class="field mb-4">
           <FloatLabel variant="on">

             <Select id="shift" v-model="currentEmployee.shift_id" :options="shifts" optionLabel="name" optionValue="id" fluid />
             <label for="shift">Pilih Shift :</label>
            </FloatLabel>
        </div>
        <div class="flex justify-end space-x-2 mt-6">
          <Button type="button" @click="closeModal" label="Batal" class="p-button-text"/>
          <Button type="submit" :loading="isSaving" :label="isSaving ? 'Menyimpan...' : 'Simpan'" />
        </div>
      </form>
    </Dialog>

    <Dialog v-model:visible="isBulkImportModalOpen" header="Import Karyawan dari Excel" :modal="true" class="w-full max-w-lg">
      <div class="p-4">
        <p class="text-text-muted mb-4">Gunakan fitur ini untuk menambahkan banyak karyawan sekaligus. Unduh template Excel, isi data karyawan, lalu unggah kembali file tersebut.</p>
        <Button @click="downloadTemplate" icon="pi pi-download" label="Unduh Template Excel" class="p-button-secondary mb-4" :loading="isDownloading"/>
        <Message v-if="!hasMultipleShifts" severity="info">Perusahaan Anda saat ini hanya memiliki satu shift atau belum mengkonfigurasi shift. Pastikan Anda telah mengatur shift yang sesuai di halaman <router-link to="/dashboard/settings/shifts" class="underline font-bold">Manajemen Shift</router-link> untuk memastikan data absensi karyawan tercatat dengan benar.</Message>
        
        <FileUpload name="bulkFile" @uploader="uploadBulkFile" :customUpload="true" :multiple="false" accept=".xlsx, .xls" :maxFileSize="1000000" :disabled="isUploading">
            <template #header="{ choose, upload, clear, files }">
                <div class="flex justify-between items-center flex-1 gap-2">
                    <div class="flex gap-2">
                        <Button @click="choose()" icon="pi pi-file" class="p-button-outlined" label="Pilih File"></Button>
                        <Button @click="upload()" icon="pi pi-cloud-upload" label="Unggah" :disabled="!files || files.length === 0"></Button>
                        <Button @click="clear()" icon="pi pi-times" class="p-button-outlined p-button-danger" label="Hapus" :disabled="!files || files.length === 0"></Button>
                    </div>
                </div>
            </template>
            <template #empty>
                <p class="text-center text-text-muted">Seret dan lepas file di sini untuk mengunggah.</p>
            </template>
        </FileUpload>

        <div v-if="bulkImportResults" class="mt-6 p-4 rounded-lg" :class="bulkImportResults.failed_count > 0 ? 'bg-red-100 text-red-800' : 'bg-green-100 text-green-800'">
          <h4 class="font-bold mb-2">Hasil Impor:</h4>
          <p>Total Diproses: {{ bulkImportResults.total_processed }}</p>
          <p>Berhasil: {{ bulkImportResults.success_count }}</p>
          <p>Gagal: {{ bulkImportResults.failed_count }}</p>
          <div v-if="bulkImportResults.failed_count > 0" class="mt-4">
            <h5 class="font-semibold">Detail Kegagalan:</h5>
            <ul class="list-disc list-inside">
              <li v-for="(result, index) in bulkImportResults.results" :key="index">{{ result.row_number || 'N/A' }}: {{ result.message }}</li>
            </ul>
          </div>
        </div>
      </div>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, computed } from 'vue';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { useConfirm } from 'primevue/useconfirm';
import { useAuthStore } from '../../stores/auth';
import { FilterMatchMode } from '@primevue/core/api';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import Select from 'primevue/select';
import FileUpload from 'primevue/fileupload';
import Tabs from 'primevue/tabs';
import TabList from 'primevue/tablist';
import Tab from 'primevue/tab';
import TabPanels from 'primevue/tabpanels';
import TabPanel from 'primevue/tabpanel';
import IconField from 'primevue/iconfield';
import InputIcon from 'primevue/inputicon';
import Message from 'primevue/message';
import Toast from 'primevue/toast';
import ConfirmDialog from 'primevue/confirmdialog';
import FloatLabel from 'primevue/floatlabel';
const toast = useToast();
const confirm = useConfirm();
const authStore = useAuthStore();

const employees = ref([]);
const pendingEmployees = ref([]);
const shifts = ref([]);

const isLoading = ref(false);
const isSaving = ref(false);
const isResending = ref(false);
const isUploading = ref(false);
const isDownloading = ref(false);
const activeTab = ref('active');

const employeesTotalRecords = ref(0);
const employeesLazyParams = ref({});
const employeesFilters = ref({ 'global': { value: null, matchMode: FilterMatchMode.CONTAINS } });

const pendingTotalRecords = ref(0);
const pendingLazyParams = ref({});
const pendingFilters = ref({ 'global': { value: null, matchMode: FilterMatchMode.CONTAINS } });

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
      toast.add({ severity: 'error', summary: 'Error', detail: response.data?.message || 'Gagal mengambil data karyawan pending.', life: 3000 });
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
        employeesLazyParams.value.page = 0;
        fetchEmployees();
    } else if (type === 'pending') {
        pendingLazyParams.value.page = 0;
        fetchPendingEmployees();
    }
};

watch(activeTab, (newTab) => {
  if (newTab === 'active') {
    fetchEmployees();
  } else if (newTab === 'pending') {
    fetchPendingEmployees();
  }
});

const onCellEditComplete = async (event) => {
  let { data, newValue, field } = event;
  data[field] = newValue;

  try {
    const updatePayload = { [field]: newValue };
    const response = await axios.put(`/api/employees/${data.id}`, updatePayload);
    if (response.data && response.data.status === 'success') {
      toast.add({ severity: 'success', summary: 'Success', detail: 'Karyawan berhasil diperbarui!', life: 3000 });
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data?.message || 'Gagal memperbarui karyawan.', life: 3000 });
      fetchEmployees();
    }
  } catch (error) {
    let message = error.response?.data?.message || 'Terjadi kesalahan saat memperbarui karyawan.';
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
    fetchEmployees();
  }
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
    let message = error.response?.data?.message || 'Failed to fetch shifts.';
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  }
};

onMounted(() => {
  fetchShifts();
  employeesLazyParams.value = { first: 0, rows: 10, page: 0 };
  pendingLazyParams.value = { first: 0, rows: 10, page: 0 };
  fetchEmployees();
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

const saveEmployee = async () => {
  if (!authStore.companyId) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Company ID not available.', life: 3000 });
    return;
  }
  isSaving.value = true;
  try {
    const response = await axios.post(`/api/employees`, currentEmployee.value);
    toast.add({ severity: 'success', summary: 'Success', detail: response.data.message || 'Employee created successfully.', life: 3000 });
    closeModal();
    fetchEmployees();
  } catch (error) {
    let message = error.response?.data?.message || 'Failed to save employee.';
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  } finally {
    isSaving.value = false;
  }
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
        fetchEmployees();
      } catch (error) {
        let message = error.response?.data?.message || 'Failed to delete employee.';
        toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
      }
    }
  });
};

const resendPasswordEmail = async (employeeId) => {
  isResending.value = true;
  try {
    const response = await axios.post(`/api/employees/${employeeId}/resend-password-email`);
    if (response.data && response.data.status === 'success') {
      toast.add({ severity: 'success', summary: 'Success', detail: 'Email pengaturan kata sandi berhasil dikirim ulang!', life: 3000 });
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data?.message || 'Gagal mengirim ulang email.', life: 3000 });
    }
  } catch (error) {
    let message = error.response?.data?.message || 'Gagal mengirim ulang email.';
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  } finally {
    isResending.value = false;
  }
  fetchPendingEmployees();
};

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
  isUploading.value = true;
  try {
    const response = await axios.post(`/api/employees/bulk`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
    if (response.data && response.data.status === 'success') {
      bulkImportResults.value = response.data.data;
      toast.add({ severity: 'success', summary: 'Success', detail: 'Impor massal selesai.', life: 3000 });
      fetchEmployees();
      event.files.length = 0;
    } else {
      bulkImportResults.value = response.data.data;
      toast.add({ severity: 'error', summary: 'Error', detail: response.data?.message || 'Impor massal gagal.', life: 3000 });
    }
  } catch (error) {
    let message = error.response?.data?.message || 'Terjadi kesalahan saat mengunggah file.';
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  } finally {
    isUploading.value = false;
  }
};

const downloadTemplate = async () => {
  isDownloading.value = true;
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
  } finally {
    isDownloading.value = false;
  }
};

const hasMultipleShifts = computed(() => shifts.value.length > 1);
</script>

<style scoped>
.field > label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 600;
}
</style>