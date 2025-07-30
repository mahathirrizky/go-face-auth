<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <Toast />
    <h2 class="text-2xl font-bold text-text-base mb-6">Riwayat Absensi Karyawan</h2>

    <div class="bg-bg-muted p-4 rounded-lg shadow-md mb-6">
        <div class="flex flex-col md:flex-row justify-between items-center gap-4">
            <div class="flex flex-col md:flex-row items-center gap-4 w-full">
                <div class="p-fluid flex-1 w-full md:w-auto">
                    <DatePicker v-model="startDate" dateFormat="yy-mm-dd" placeholder="Dari Tanggal" class="w-full" fluid />
                </div>
                <div class="p-fluid flex-1 w-full md:w-auto">
                    <DatePicker v-model="endDate" dateFormat="yy-mm-dd" placeholder="Sampai Tanggal" class="w-full" fluid />
                </div>
                <Button @click="fetchAttendanceHistory" icon="pi pi-filter" label="Filter" class="w-full md:w-auto" />
            </div>
            <Button @click="exportToExcel" icon="pi pi-file-excel" label="Export ke Excel" class="p-button-secondary w-full md:w-auto whitespace-nowrap" />
        </div>
    </div>

    <DataTable
      :value="attendances"
      :loading="isLoading"
      paginator :rows="10" :rowsPerPageOptions="[10, 25, 50]"
      currentPageReportTemplate="Menampilkan {first} sampai {last} dari {totalRecords} data"
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink CurrentPageReport RowsPerPageDropdown"
      dataKey="id"
      :globalFilterFields="['status']"
      v-model:filters="filters"
    >
      <template #header>
        <div class="flex justify-end">
            <IconField iconPosition="left">
                <InputIcon class="pi pi-search"></InputIcon>
                <InputText v-model="filters['global'].value" placeholder="Cari Riwayat..." fluid />
            </IconField>
        </div>
      </template>

      <template #empty>Tidak ada data riwayat absensi.</template>
      <template #loading>Memuat data...</template>

      <Column field="date" header="Tanggal" :sortable="true">
        <template #body="{ item }">{{ formatDateTime(item.check_in_time, 'date') }}</template>
      </Column>
      <Column field="check_in_time" header="Waktu Check-in" :sortable="true">
        <template #body="{ item }">{{ formatDateTime(item.check_in_time, 'time') }}</template>
      </Column>
      <Column field="check_out_time" header="Waktu Check-out" :sortable="true">
        <template #body="{ item }">{{ item.check_out_time ? formatDateTime(item.check_out_time, 'time') : 'Belum Check-out' }}</template>
      </Column>
      <Column field="status" header="Status" :sortable="true">
        <template #body="{ data }">
          <Tag :value="data.status" :severity="getStatusSeverity(data.status)" />
        </template>
      </Column>
    </DataTable>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { FilterMatchMode } from '@primevue/core/api';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import DatePicker from 'primevue/datepicker';
import IconField from 'primevue/iconfield';
import InputIcon from 'primevue/inputicon';
import Tag from 'primevue/tag';
import Toast from 'primevue/toast';

const route = useRoute();
const toast = useToast();
const employeeId = ref(route.params.employeeId);
const attendances = ref([]);
const isLoading = ref(false);

const filters = ref({
    global: { value: null, matchMode: FilterMatchMode.CONTAINS }
});

const today = new Date();
const firstDayOfMonth = new Date(today.getFullYear(), today.getMonth(), 1);
const lastDayOfMonth = new Date(today.getFullYear(), today.getMonth() + 1, 0);

const startDate = ref(firstDayOfMonth);
const endDate = ref(lastDayOfMonth);

const formatToYYYYMMDD = (date) => {
  if (!date) return null;
  const d = new Date(date);
  const year = d.getFullYear();
  const month = String(d.getMonth() + 1).padStart(2, '0');
  const day = String(d.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
};

const fetchAttendanceHistory = async () => {
  isLoading.value = true;
  try {
    const params = {
      startDate: formatToYYYYMMDD(startDate.value),
      endDate: formatToYYYYMMDD(endDate.value),
    };

    const response = await axios.get(`/api/employees/${employeeId.value}/attendances`, { params });
    if (response.data && response.data.status === 'success') {
      attendances.value = response.data.data;
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data?.message || 'Failed to fetch attendance history.', life: 3000 });
    }
  } catch (error) {
    console.error('Error fetching attendance history:', error);
    let message = error.response?.data?.message || 'Failed to fetch attendance history.';
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  } finally {
    isLoading.value = false;
  }
};

const exportToExcel = async () => {
  try {
    const params = {
      startDate: formatToYYYYMMDD(startDate.value),
      endDate: formatToYYYYMMDD(endDate.value),
    };

    const response = await axios.get(`/api/employees/${employeeId.value}/attendances/export`, {
      params,
      responseType: 'blob',
    });

    const blob = new Blob([response.data], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = `riwayat_absensi_${employeeId.value}.xlsx`;
    link.click();
    URL.revokeObjectURL(link.href);
    toast.add({ severity: 'success', summary: 'Success', detail: 'File Excel berhasil diunduh!', life: 3000 });
  } catch (error) {
    console.error('Error exporting to Excel:', error);
    let message = error.response?.data?.message || 'Failed to export to Excel.';
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  }
};

const formatDateTime = (dateTimeString, type = 'datetime') => {
  if (!dateTimeString) return '';
  const date = new Date(dateTimeString);
  const options = {
    date: { year: 'numeric', month: 'long', day: 'numeric' },
    time: { hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false },
    datetime: { year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false }
  };
  return date.toLocaleString('id-ID', options[type]);
};

const getStatusSeverity = (status) => {
    switch (status) {
        case 'Tepat Waktu': return 'success';
        case 'Terlambat': return 'warning';
        case 'Alpha': return 'danger';
        default: return 'info';
    }
};

onMounted(() => {
  fetchAttendanceHistory();
});
</script>
