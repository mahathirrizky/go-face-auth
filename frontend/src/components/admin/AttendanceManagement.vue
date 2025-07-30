<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <Toast />
    <h2 class="text-2xl font-bold text-text-base mb-6">Manajemen Absensi</h2>

    <Tabs v-model:value="selectedTab">
      <TabList>
        <Tab value="all">Semua Absensi</Tab>
        <Tab value="unaccounted">Karyawan Tidak Absen</Tab>
        <Tab value="overtime">Lembur</Tab>
      </TabList>
      <TabPanels>
        <TabPanel value="all">
          <DataTable
            :value="attendanceRecords"
            :loading="isLoading"
            :totalRecords="attendancesTotalRecords"
            :lazy="true"
            v-model:filters="attendancesFilters"
            @page="onPage($event, 'attendances')"
            @filter="onFilter($event, 'attendances')"
            paginator :rows="10" :rowsPerPageOptions="[10, 25, 50]"
            dataKey="id"
            :globalFilterFields="['employee.name', 'status']"
          >
            <template #header>
              <div class="flex flex-wrap items-center justify-between gap-4">
                  <IconField iconPosition="left">
                      <InputIcon class="pi pi-search"></InputIcon>
                      <InputText v-model="attendancesFilters['global'].value" placeholder="Cari Absensi..." @keydown.enter="onFilter($event, 'attendances')"/>
                  </IconField>
                  <div class="flex flex-wrap items-center gap-2">
                      <DatePicker v-model="startDate" dateFormat="yy-mm-dd" placeholder="Dari Tanggal"/>
                      <DatePicker v-model="endDate" dateFormat="yy-mm-dd" placeholder="Sampai Tanggal"/>
                      <Button @click="fetchAttendances" icon="pi pi-filter" label="Filter" />
                      <Button @click="exportAllToExcel" icon="pi pi-file-excel" label="Export" class="p-button-secondary" />
                      <Button @click="openCorrectionModal()" icon="pi pi-plus" label="Koreksi" class="p-button-help" />
                  </div>
              </div>
            </template>
            <Column field="date" header="Tanggal" :sortable="true">
              <template #body="{ data }">{{ new Date(data.check_in_time).toLocaleDateString('id-ID') }}</template>
            </Column>
            <Column field="employee.name" header="Nama Karyawan" :sortable="true"></Column>
            <Column field="check_in_time" header="Waktu Masuk" :sortable="true">
              <template #body="{ data }">{{ new Date(data.check_in_time).toLocaleTimeString('id-ID') }}</template>
            </Column>
            <Column field="check_out_time" header="Waktu Keluar" :sortable="true">
              <template #body="{ data }">{{ data.check_out_time ? new Date(data.check_out_time).toLocaleTimeString('id-ID') : '-' }}</template>
            </Column>
            <Column field="status" header="Status" :sortable="true">
              <template #body="{ data }">
                <Tag :value="data.is_correction ? 'Dikoreksi' : data.status" :severity="getStatusSeverity(data.status, data.is_correction)" />
              </template>
            </Column>
            <Column header="Aksi">
              <template #body="{ data }">
                <Button v-if="!data.check_out_time" @click="openCorrectionModal(data)" icon="pi pi-pencil" class="p-button-info p-button-sm" v-tooltip.top="'Tambah Waktu Pulang'" />
              </template>
            </Column>
          </DataTable>
        </TabPanel>
        <TabPanel value="unaccounted">
          <DataTable
            :value="unaccountedEmployees"
            :loading="isLoading"
            :totalRecords="unaccountedEmployeesTotalRecords"
            :lazy="true"
            v-model:filters="unaccountedEmployeesFilters"
            @page="onPage($event, 'unaccounted')"
            @filter="onFilter($event, 'unaccounted')"
            paginator :rows="10" :rowsPerPageOptions="[10, 25, 50]"
            dataKey="id"
            :globalFilterFields="['name', 'email']"
          >
            <template #header>
              <div class="flex justify-between items-center">
                  <Button @click="fetchUnaccountedEmployees" icon="pi pi-refresh" label="Refresh" />
                  <IconField iconPosition="left">
                      <InputIcon class="pi pi-search"></InputIcon>
                      <InputText v-model="unaccountedEmployeesFilters['global'].value" placeholder="Cari Karyawan..." @keydown.enter="onFilter($event, 'unaccounted')"/>
                  </IconField>
              </div>
            </template>
            <Column field="name" header="Nama Karyawan" :sortable="true"></Column>
            <Column field="email" header="Email" :sortable="true"></Column>
            <Column header="Aksi">
              <template #body="{ data }">
                <Button @click="openCorrectionModal(data)" icon="pi pi-plus" class="p-button-info p-button-sm" v-tooltip.top="'Tambah Absensi'" />
              </template>
            </Column>
          </DataTable>
        </TabPanel>
        <TabPanel value="overtime">
          <DataTable
            :value="overtimeRecords"
            :loading="isLoading"
            :totalRecords="overtimeTotalRecords"
            :lazy="true"
            v-model:filters="overtimeFilters"
            @page="onPage($event, 'overtime')"
            @filter="onFilter($event, 'overtime')"
            paginator :rows="10" :rowsPerPageOptions="[10, 25, 50]"
            dataKey="id"
            :globalFilterFields="['employee.name']"
          >
            <template #header>
              <div class="flex justify-between items-center">
                  <Button @click="fetchOvertimeRecords" icon="pi pi-refresh" label="Refresh" />
                  <IconField iconPosition="left">
                      <InputIcon class="pi pi-search"></InputIcon>
                      <InputText v-model="overtimeFilters['global'].value" placeholder="Cari Lembur..." @keydown.enter="onFilter($event, 'overtime')"/>
                  </IconField>
              </div>
            </template>
            <Column field="employee.name" header="Nama Karyawan" :sortable="true"></Column>
            <Column field="date" header="Tanggal" :sortable="true"></Column>
            <Column field="check_in_time" header="Waktu Masuk" :sortable="true"></Column>
            <Column field="check_out_time" header="Waktu Keluar" :sortable="true"></Column>
            <Column field="overtime_minutes" header="Durasi Lembur" :sortable="true">
              <template #body="{ data }">{{ data.overtime_minutes ? `${data.overtime_minutes} menit` : '-' }}</template>
            </Column>
          </DataTable>
        </TabPanel>
      </TabPanels>
    </Tabs>

    <Dialog v-model:visible="isCorrectionModalOpen" header="Koreksi Absensi" :modal="true" class="w-full max-w-md">
      <form @submit.prevent="submitCorrection" class="p-fluid mt-4">
        <div class="field mb-4">
          <label for="employee">Karyawan</label>
          <Select v-model="correctionForm.employee_id" :options="allEmployees" optionLabel="name" optionValue="id" placeholder="Pilih Karyawan" :disabled="!!correctionForm.id" filter />
        </div>
        <div class="field mb-4">
          <label for="correction_time">Tanggal & Waktu Koreksi</label>
          <DatePicker v-model="correctionForm.correction_time" showTime hourFormat="24" />
        </div>
        <div class="field mb-4">
          <label>Tipe Koreksi</label>
          <div class="flex items-center gap-4">
            <div class="flex items-center">
              <RadioButton v-model="correctionForm.correction_type" inputId="type_check_in" name="correction_type" value="check_in" :disabled="!!correctionForm.id" />
              <label for="type_check_in" class="ml-2">Check-in</label>
            </div>
            <div class="flex items-center">
              <RadioButton v-model="correctionForm.correction_type" inputId="type_check_out" name="correction_type" value="check_out" :disabled="!correctionForm.id && correctionForm.correction_type !== 'check_out'" />
              <label for="type_check_out" class="ml-2">Check-out</label>
            </div>
          </div>
        </div>
        <div class="field mb-4">
          <label for="notes">Alasan Koreksi</label>
          <Textarea v-model="correctionForm.notes" rows="3" placeholder="Contoh: Karyawan lupa absen pulang." />
        </div>
        <div class="flex justify-end space-x-2 mt-6">
          <Button type="button" @click="closeCorrectionModal" label="Batal" class="p-button-text"/>
          <Button type="submit" :loading="isSubmitting" label="Simpan" />
        </div>
      </form>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { useAuthStore } from '../../stores/auth';
import { FilterMatchMode } from '@primevue/core/api';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import Select from 'primevue/select';
import DatePicker from 'primevue/datepicker';
import RadioButton from 'primevue/radiobutton';
import Textarea from 'primevue/textarea';
import Tabs from 'primevue/tabs';
import TabList from 'primevue/tablist';
import Tab from 'primevue/tab';
import TabPanels from 'primevue/tabpanels';
import TabPanel from 'primevue/tabpanel';
import IconField from 'primevue/iconfield';
import InputIcon from 'primevue/inputicon';
import Tag from 'primevue/tag';
import Toast from 'primevue/toast';

const attendanceRecords = ref([]);
const unaccountedEmployees = ref([]);
const overtimeRecords = ref([]);
const selectedTab = ref('all');
const toast = useToast();
const authStore = useAuthStore();
const isLoading = ref(false);
const isSubmitting = ref(false);

const startDate = ref(null);
const endDate = ref(null);

const isCorrectionModalOpen = ref(false);
const allEmployees = ref([]);
const correctionForm = ref({
  id: null,
  employee_id: null,
  correction_time: null,
  correction_type: 'check_in',
  notes: ''
});

const attendancesTotalRecords = ref(0);
const attendancesLazyParams = ref({});
const attendancesFilters = ref({ 'global': { value: null, matchMode: FilterMatchMode.CONTAINS } });

const unaccountedEmployeesTotalRecords = ref(0);
const unaccountedEmployeesLazyParams = ref({});
const unaccountedEmployeesFilters = ref({ 'global': { value: null, matchMode: FilterMatchMode.CONTAINS } });

const overtimeTotalRecords = ref(0);
const overtimeLazyParams = ref({});
const overtimeFilters = ref({ 'global': { value: null, matchMode: FilterMatchMode.CONTAINS } });

const fetchAttendances = async () => {
  isLoading.value = true;
  try {
    const params = {
      page: attendancesLazyParams.value.page ? attendancesLazyParams.value.page + 1 : 1,
      limit: attendancesLazyParams.value.rows || 10,
      sortField: attendancesLazyParams.value.sortField,
      sortOrder: attendancesLazyParams.value.sortOrder,
      search: attendancesFilters.value.global.value,
      startDate: startDate.value ? startDate.value.toISOString().split('T')[0] : null,
      endDate: endDate.value ? endDate.value.toISOString().split('T')[0] : null,
    };
    const response = await axios.get('/api/attendances', { params });
    attendanceRecords.value = response.data.data.items || [];
    attendancesTotalRecords.value = response.data.data.total_records || 0;
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memuat data absensi.', life: 3000 });
  } finally {
    isLoading.value = false;
  }
};

const fetchUnaccountedEmployees = async () => {
  isLoading.value = true;
  try {
    const params = {
      page: unaccountedEmployeesLazyParams.value.page ? unaccountedEmployeesLazyParams.value.page + 1 : 1,
      limit: unaccountedEmployeesLazyParams.value.rows || 10,
    };
    const response = await axios.get('/api/attendances/unaccounted', { params });
    unaccountedEmployees.value = response.data.data.items || [];
    unaccountedEmployeesTotalRecords.value = response.data.data.total_records || 0;
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memuat data karyawan tidak absen.', life: 3000 });
  } finally {
    isLoading.value = false;
  }
};

const fetchOvertimeRecords = async () => {
  isLoading.value = true;
  try {
    const params = {
      page: overtimeLazyParams.value.page ? overtimeLazyParams.value.page + 1 : 1,
      limit: overtimeLazyParams.value.rows || 10,
      sortField: overtimeLazyParams.value.sortField,
      sortOrder: overtimeLazyParams.value.sortOrder,
      search: overtimeFilters.value.global.value,
    };
    const response = await axios.get('/api/attendances/overtime', { params });
    overtimeRecords.value = response.data.data.items || [];
    overtimeTotalRecords.value = response.data.data.total_records || 0;
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memuat data lembur.', life: 3000 });
  } finally {
    isLoading.value = false;
  }
};

const exportAllToExcel = () => {
  toast.add({ severity: 'info', summary: 'Info', detail: 'Fitur export sedang dalam pengembangan.', life: 3000 });
};

const onPage = (event, tableType) => {
  if (tableType === 'attendances') {
    attendancesLazyParams.value = event;
    fetchAttendances();
  } else if (tableType === 'unaccounted') {
    unaccountedEmployeesLazyParams.value = event;
    fetchUnaccountedEmployees();
  } else if (tableType === 'overtime') {
    overtimeLazyParams.value = event;
    fetchOvertimeRecords();
  }
};

const onFilter = (event, tableType) => {
    const lazyParams = tableType === 'attendances' ? attendancesLazyParams : (tableType === 'unaccounted' ? unaccountedEmployeesLazyParams : overtimeLazyParams);
    lazyParams.value.page = 0;
    if (tableType === 'attendances') fetchAttendances();
    else if (tableType === 'unaccounted') fetchUnaccountedEmployees();
    else if (tableType === 'overtime') fetchOvertimeRecords();
};

const fetchAllEmployees = async () => {
  if (!authStore.companyId) return;
  try {
    const response = await axios.get(`/api/companies/${authStore.companyId}/employees?limit=1000`); // Fetch all for Select
    allEmployees.value = response.data.data.items;
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memuat daftar karyawan.', life: 3000 });
  }
};

const openCorrectionModal = (item = null) => {
  if (item && item.check_in_time) { // Editing existing record for check-out
    correctionForm.value = {
      id: item.id,
      employee_id: item.employee_id,
      correction_time: new Date(),
      correction_type: 'check_out',
      notes: ''
    };
  } else { // Creating new record or correcting unaccounted
    correctionForm.value = {
      id: null,
      employee_id: item ? item.id : null,
      correction_time: new Date(),
      correction_type: 'check_in',
      notes: ''
    };
  }
  isCorrectionModalOpen.value = true;
};

const closeCorrectionModal = () => {
  isCorrectionModalOpen.value = false;
};

const submitCorrection = async () => {
  isSubmitting.value = true;
  try {
    const payload = {
      employee_id: correctionForm.value.employee_id,
      correction_time: correctionForm.value.correction_time.toISOString(),
      correction_type: correctionForm.value.correction_type,
      notes: correctionForm.value.notes,
      attendance_id: correctionForm.value.id
    };
    
    await axios.post('/api/attendances/correction', payload);
    
    toast.add({ severity: 'success', summary: 'Sukses', detail: 'Absensi berhasil dikoreksi.', life: 3000 });
    closeCorrectionModal();
    fetchAttendances();
    fetchUnaccountedEmployees();
  } catch (error) {
    const message = error.response?.data?.message || 'Gagal menyimpan koreksi.';
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  } finally {
    isSubmitting.value = false;
  }
};

const getStatusSeverity = (status, isCorrection) => {
    if (isCorrection) return 'info';
    switch (status) {
        case 'Tepat Waktu': return 'success';
        case 'Terlambat': return 'warning';
        case 'Alpha': return 'danger';
        case 'Cuti': return 'help';
        case 'Sakit': return 'contrast';
        default: return 'secondary';
    }
};

watch(selectedTab, (newTab) => {
  if (newTab === 'all') fetchAttendances();
  else if (newTab === 'unaccounted') fetchUnaccountedEmployees();
  else if (newTab === 'overtime') fetchOvertimeRecords();
});

onMounted(() => {
  fetchAllEmployees();
  fetchAttendances();
});
</script>