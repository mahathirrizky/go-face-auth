<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-6">Manajemen Absensi</h2>

    <Tabs v-model:value="selectedTab">
      <TabList>
        <Tab :value="0">Semua Absensi</Tab>
        <Tab :value="1">Karyawan Tidak Absen</Tab>
        <Tab :value="2">Lembur</Tab>
      </TabList>

      <TabPanels>
        <TabPanel :value="0">
          <BaseDataTable
            :data="attendanceRecords"
            :columns="attendanceColumns"
            :loading="isLoading"
            :totalRecords="attendancesTotalRecords"
            :lazy="true"
            v-model:filters="attendancesFilters"
            @page="onPage($event, 'attendances')"
            @filter="onFilter($event, 'attendances')"
            searchPlaceholder="Cari Absensi..."
          >
            <template #header-actions>
              <div class="flex flex-wrap items-center gap-2">
                <div class="flex items-center">
                  <label for="startDate" class="text-text-muted mr-2">Dari:</label>
                  <BaseInput
                    type="date"
                    id="startDate"
                    v-model="startDate"
                    class="p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary"
                    :label-sr-only="true"
                  />
                </div>
                <div class="flex items-center">
                  <label for="endDate" class="text-text-muted mr-2">Sampai:</label>
                  <BaseInput
                    type="date"
                    id="endDate"
                    v-model="endDate"
                    class="p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary"
                    :label-sr-only="true"
                  />
                </div>
                <BaseButton @click="fetchAttendances" class="btn-primary"><i class="pi pi-filter"></i> Filter</BaseButton>
                <BaseButton @click="exportAllToExcel" class="btn-secondary whitespace-nowrap"><i class="pi pi-file-excel"></i> Export to Excel</BaseButton>
              </div>
            </template>
            <template #column-date="{ item }">
              {{ new Date(item.check_in_time).toLocaleDateString() }}
            </template>

            <template #column-check_in_time="{ item }">
              {{ new Date(item.check_in_time).toLocaleTimeString() }}
            </template>

            <template #column-check_out_time="{ item }">
              {{ item.check_out_time ? new Date(item.check_out_time).toLocaleTimeString() : '-' }}
            </template>

            <template #column-status="{ item }">
              <span :class="{
                'px-2 inline-flex text-xs leading-5 font-semibold rounded-full': true,
                'bg-green-100 text-green-600': item.status === 'on_time',
                'bg-yellow-100 text-yellow-600': item.status === 'late',
                'bg-blue-100 text-blue-600': item.status === 'overtime_in' || item.status === 'overtime_out',
              }">
                {{ item.status === 'on_time' ? 'Tepat Waktu' : item.status === 'late' ? 'Terlambat' : item.status === 'overtime_in' ? 'Lembur Masuk' : item.status === 'overtime_out' ? 'Lembur Keluar' : item.status }}
              </span>
            </template>

            <template #column-employee.name="{ item }">
              {{ item.employee ? item.employee.name : 'N/A' }}
            </template>
          </BaseDataTable>
        </TabPanel>

        <TabPanel :value="1">
          <BaseDataTable
            :data="unaccountedEmployees"
            :columns="unaccountedEmployeeColumns"
            :loading="isLoading"
            :totalRecords="unaccountedTotalRecords"
            :lazy="true"
            v-model:filters="unaccountedFilters"
            @page="onPage($event, 'unaccounted')"
            @filter="onFilter($event, 'unaccounted')"
            searchPlaceholder="Cari Karyawan..."
          >
            <template #header-actions>
              <div class="flex flex-wrap items-center gap-2">
                <div class="flex items-center">
                  <label for="unaccountedStartDate" class="text-text-muted mr-2">Dari:</label>
                  <BaseInput
                    type="date"
                    id="unaccountedStartDate"
                    v-model="unaccountedStartDate"
                    class="p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary"
                    :label-sr-only="true"
                  />
                </div>
                <div class="flex items-center">
                  <label for="unaccountedEndDate" class="text-text-muted mr-2">Sampai:</label>
                  <BaseInput
                    type="date"
                    id="unaccountedEndDate"
                    v-model="unaccountedEndDate"
                    class="p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary"
                    :label-sr-only="true"
                  />
                </div>
                <BaseButton @click="fetchUnaccountedEmployees" class="btn-primary"><i class="pi pi-search"></i> Cari</BaseButton>
              </div>
            </template>
          </BaseDataTable>
        </TabPanel>

        <TabPanel :value="2">
          <BaseDataTable
            :data="overtimeRecords"
            :columns="overtimeColumns"
            :loading="isLoading"
            :totalRecords="overtimeTotalRecords"
            :lazy="true"
            v-model:filters="overtimeFilters"
            @page="onPage($event, 'overtime')"
            @filter="onFilter($event, 'overtime')"
            searchPlaceholder="Cari Lembur..."
          >
            <template #header-actions>
              <div class="flex flex-wrap items-center gap-2">
                <div class="flex items-center">
                  <label for="overtimeStartDate" class="text-text-muted mr-2">Dari:</label>
                  <BaseInput
                    type="date"
                    id="overtimeStartDate"
                    v-model="overtimeStartDate"
                    class="p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary"
                    :label-sr-only="true"
                  />
                </div>
                <div class="flex items-center">
                  <label for="overtimeEndDate" class="text-text-muted mr-2">Sampai:</label>
                  <BaseInput
                    type="date"
                    id="overtimeEndDate"
                    v-model="overtimeEndDate"
                    class="p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary"
                    :label-sr-only="true"
                  />
                </div>
                <BaseButton @click="fetchOvertimeAttendances" class="btn-primary"><i class="pi pi-search"></i> Cari</BaseButton>
              </div>
            </template>
            <template #column-check_in_time="{ item }">
              {{ new Date(item.check_in_time).toLocaleString() }}
            </template>

            <template #column-check_out_time="{ item }">
              {{ item.check_out_time ? new Date(item.check_out_time).toLocaleString() : '-' }}
            </template>

            <template #column-employee.name="{ item }">
              {{ item.employee ? item.employee.name : 'N/A' }}
            </template>
          </BaseDataTable>
        </TabPanel>
      </TabPanels>
    </Tabs>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { useAuthStore } from '../../stores/auth';
import BaseInput from '../ui/BaseInput.vue';
import BaseButton from '../ui/BaseButton.vue';
import BaseDataTable from '../ui/BaseDataTable.vue';
import Tabs from 'primevue/tabs';
import Tab from 'primevue/tab';
import TabList from 'primevue/tablist';
import TabPanels from 'primevue/tabpanels';
import TabPanel from 'primevue/tabpanel';
import { FilterMatchMode } from '@primevue/core/api';

const attendanceRecords = ref([]);
const unaccountedEmployees = ref([]);
const overtimeRecords = ref([]);
const selectedTab = ref(0); // Changed to numerical index
const toast = useToast();
const authStore = useAuthStore();
const isLoading = ref(false);

// State for All Attendances table
const attendancesTotalRecords = ref(0);
const attendancesLazyParams = ref({});
const attendancesFilters = ref({ 'global': { value: null, matchMode: FilterMatchMode.CONTAINS } });

// State for Unaccounted Employees table
const unaccountedTotalRecords = ref(0);
const unaccountedLazyParams = ref({});
const unaccountedFilters = ref({ 'global': { value: null, matchMode: FilterMatchMode.CONTAINS } });

// State for Overtime Attendances table
const overtimeTotalRecords = ref(0);
const overtimeLazyParams = ref({});
const overtimeFilters = ref({ 'global': { value: null, matchMode: FilterMatchMode.CONTAINS } });

const attendanceColumns = ref([
    { field: 'date', header: 'Tanggal' },
    { field: 'employee.name', header: 'Nama Karyawan' },
    { field: 'check_in_time', header: 'Waktu Masuk' },
    { field: 'check_out_time', header: 'Waktu Keluar' },
    { field: 'status', header: 'Status' }
]);

const unaccountedEmployeeColumns = ref([
    { field: 'name', header: 'Nama Karyawan' },
    { field: 'email', header: 'Email' },
    { field: 'position', header: 'Posisi' }
]);

const overtimeColumns = ref([
    { field: 'employee.name', header: 'Nama Karyawan' },
    { field: 'check_in_time', header: 'Waktu Masuk Lembur' },
    { field: 'check_out_time', header: 'Waktu Keluar Lembur' },
    { field: 'overtime_minutes', header: 'Durasi Lembur (Menit)' }
]);

const formatToYYYYMMDD = (date) => {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
};

const today = new Date();
const firstDayOfMonth = new Date(today.getFullYear(), today.getMonth(), 1);

const startDate = ref(formatToYYYYMMDD(firstDayOfMonth));
const endDate = ref(formatToYYYYMMDD(today));
const unaccountedStartDate = ref(formatToYYYYMMDD(firstDayOfMonth));
const unaccountedEndDate = ref(formatToYYYYMMDD(today));
const overtimeStartDate = ref(formatToYYYYMMDD(firstDayOfMonth));
const overtimeEndDate = ref(formatToYYYYMMDD(today));

const fetchAttendances = async () => {
  if (!authStore.companyId) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Company ID not available. Cannot fetch attendances.', life: 3000 });
    return;
  }
  isLoading.value = true;
  try {
    const params = {
      page: attendancesLazyParams.value.page + 1,
      limit: attendancesLazyParams.value.rows,
      search: attendancesFilters.value.global.value || '',
      startDate: startDate.value,
      endDate: endDate.value,
    };

    const response = await axios.get(`/api/attendances`, { params });
    if (response.data && response.data.status === 'success') {
      attendanceRecords.value = response.data.data.items;
      attendancesTotalRecords.value = response.data.data.total_records;
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data?.message || 'Failed to fetch attendances.', life: 3000 });
    }
  } catch (error) {
    console.error('Error fetching attendances:', error);
    let message = 'Failed to fetch attendances.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  } finally {
    isLoading.value = false;
  }
};

const fetchUnaccountedEmployees = async () => {
  if (!authStore.companyId) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Company ID not available.', life: 3000 });
    return;
  }
  isLoading.value = true;
  try {
    const params = {
      page: unaccountedLazyParams.value.page + 1,
      limit: unaccountedLazyParams.value.rows,
      search: unaccountedFilters.value.global.value || '',
      startDate: unaccountedStartDate.value,
      endDate: unaccountedEndDate.value,
    };

    const response = await axios.get(`/api/attendances/unaccounted`, { params });
    if (response.data && response.data.status === 'success') {
      unaccountedEmployees.value = response.data.data.items;
      unaccountedTotalRecords.value = response.data.data.total_records;
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data?.message || 'Failed to fetch unaccounted employees.', life: 3000 });
    }
  } catch (error) {
    console.error('Error fetching unaccounted employees:', error);
    let message = 'Failed to fetch unaccounted employees.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  } finally {
    isLoading.value = false;
  }
};

const fetchOvertimeAttendances = async () => {
  if (!authStore.companyId) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Company ID not available.', life: 3000 });
    return;
  }
  isLoading.value = true;
  try {
    const params = {
      page: overtimeLazyParams.value.page + 1,
      limit: overtimeLazyParams.value.rows,
      search: overtimeFilters.value.global.value || '',
      startDate: overtimeStartDate.value,
      endDate: overtimeEndDate.value,
    };

    const response = await axios.get(`/api/attendances/overtime`, { params });
    if (response.data && response.data.status === 'success') {
      overtimeRecords.value = response.data.data.items;
      overtimeTotalRecords.value = response.data.data.total_records;
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data?.message || 'Failed to fetch overtime attendances.', life: 3000 });
    }
  } catch (error) {
    console.error('Error fetching overtime attendances:', error);
    let message = 'Failed to fetch overtime attendances.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  } finally {
    isLoading.value = false;
  }
};

const onPage = (event, type) => {
  if (type === 'attendances') {
    attendancesLazyParams.value = event;
    fetchAttendances();
  } else if (type === 'unaccounted') {
    unaccountedLazyParams.value = event;
    fetchUnaccountedEmployees();
  } else if (type === 'overtime') {
    overtimeLazyParams.value = event;
    fetchOvertimeAttendances();
  }
};

const onFilter = (event, type) => {
  // PrimeVue updates filters via v-model, just need to trigger fetch
  if (type === 'attendances') {
    fetchAttendances();
  } else if (type === 'unaccounted') {
    fetchUnaccountedEmployees();
  } else if (type === 'overtime') {
    fetchOvertimeAttendances();
  }
};

const exportAllToExcel = async () => {
  if (!authStore.companyId) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Company ID not available. Cannot export.', life: 3000 });
    return;
  }
  try {
    const params = {};
    if (startDate.value) {
      params.startDate = startDate.value;
    }
    if (endDate.value) {
      params.endDate = endDate.value;
    }

    const response = await axios.get(`/api/attendances/export`, {
      params,
      responseType: 'blob',
    });

    const blob = new Blob([response.data], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = `all_company_attendance.xlsx`;
    link.click();
    URL.revokeObjectURL(link.href);
    toast.add({ severity: 'success', summary: 'Success', detail: 'File Excel semua absensi berhasil diunduh!', life: 3000 });
  } catch (error) {
    console.error('Error exporting all attendances to Excel:', error);
    let message = 'Failed to export all attendances to Excel.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  }
};

onMounted(() => {
  // Initialize lazy params for all tables
  attendancesLazyParams.value = { first: 0, rows: 10, page: 0 };
  unaccountedLazyParams.value = { first: 0, rows: 10, page: 0 };
  overtimeLazyParams.value = { first: 0, rows: 10, page: 0 };

  // Fetch initial data for the default tab
  fetchAttendances();
});

watch(() => selectedTab.value, (newTab) => {
  if (newTab === 0) {
    fetchAttendances();
  } else if (newTab === 1) {
    fetchUnaccountedEmployees();
  } else if (newTab === 2) {
    fetchOvertimeAttendances();
  }
});
</script>

<style scoped>
/* Tailwind handles styling */
</style>
