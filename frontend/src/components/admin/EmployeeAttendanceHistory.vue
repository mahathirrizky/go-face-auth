<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-6">Riwayat Absensi Karyawan</h2>

    <div class="bg-bg-muted p-4 rounded-lg shadow-md mb-6 flex flex-col md:flex-row justify-between items-center">
      <div class="flex flex-col md:flex-row items-center space-y-4 md:space-y-0 md:space-x-4 w-full">
        <div class="flex items-center space-x-2">
          <label for="startDate" class="text-text-muted">Dari:</label>
          <BaseInput
            type="date"
            id="startDate"
            v-model="startDate"
            class="p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary"
            :label-sr-only="true"
          />
        </div>
        <div class="flex items-center space-x-2">
          <label for="endDate" class="text-text-muted">Sampai:</label>
          <BaseInput
            type="date"
            id="endDate"
            v-model="endDate"
            class="p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary"
            :label-sr-only="true"
          />
        </div>
        <BaseButton @click="fetchAttendanceHistory" class="btn-primary w-full md:w-auto"><i class="fas fa-filter"></i> Filter</BaseButton>
      </div>
      <BaseButton @click="exportToExcel" class="btn-secondary w-full md:w-auto mt-4 md:mt-0 whitespace-nowrap"><i class="fas fa-file-excel"></i> Export ke Excel</BaseButton>
    </div>

    <BaseDataTable
      :data="attendances"
      :columns="attendanceColumns"
      :loading="isLoading"
      :globalFilterFields="['status']"
      searchPlaceholder="Cari Riwayat..."
    >
      <template #column-date="{ item }">
        {{ formatDateTime(item.check_in_time, 'date') }}
      </template>

      <template #column-check_in_time="{ item }">
        {{ formatDateTime(item.check_in_time, 'time') }}
      </template>

      <template #column-check_out_time="{ item }">
        {{ item.check_out_time ? formatDateTime(item.check_out_time, 'time') : 'Belum Check-out' }}
      </template>
    </BaseDataTable>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import BaseInput from '../ui/BaseInput.vue';
import BaseButton from '../ui/BaseButton.vue';
import BaseDataTable from '../ui/BaseDataTable.vue';

const route = useRoute();
const toast = useToast();
const employeeId = ref(route.params.employeeId);
const attendances = ref([]);
const isLoading = ref(false);

const attendanceColumns = ref([
    { field: 'date', header: 'Tanggal' },
    { field: 'check_in_time', header: 'Waktu Check-in' },
    { field: 'check_out_time', header: 'Waktu Check-out' },
    { field: 'status', header: 'Status' }
]);

const today = new Date();
const firstDayOfMonth = new Date(today.getFullYear(), today.getMonth(), 1);
const lastDayOfMonth = new Date(today.getFullYear(), today.getMonth() + 1, 0);

const formatToYYYYMMDD = (date) => {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
};

const startDate = ref(formatToYYYYMMDD(firstDayOfMonth));
const endDate = ref(formatToYYYYMMDD(lastDayOfMonth));

const fetchAttendanceHistory = async () => {
  isLoading.value = true;
  try {
    let url = `/api/employees/${employeeId.value}/attendances`;
    const params = {};
    if (startDate.value) {
      params.startDate = startDate.value;
    }
    if (endDate.value) {
      params.endDate = endDate.value;
    }

    const response = await axios.get(url, { params });
    if (response.data && response.data.status === 'success') {
      attendances.value = response.data.data;
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data?.message || 'Failed to fetch attendance history.', life: 3000 });
    }
  } catch (error) {
    console.error('Error fetching attendance history:', error);
    let message = 'Failed to fetch attendance history.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  } finally {
    isLoading.value = false;
  }
};

const exportToExcel = async () => {
  try {
    let url = `/api/employees/${employeeId.value}/attendances/export`;
    const params = {};
    if (startDate.value) {
      params.startDate = startDate.value;
    }
    if (endDate.value) {
      params.endDate = endDate.value;
    }

    const response = await axios.get(url, {
      params,
      responseType: 'blob',
    });

    const blob = new Blob([response.data], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = `riwayat_absensi_karyawan_${employeeId.value}.xlsx`;
    link.click();
    URL.revokeObjectURL(link.href);
    toast.add({ severity: 'success', summary: 'Success', detail: 'File Excel berhasil diunduh!', life: 3000 });
  } catch (error) {
    console.error('Error exporting to Excel:', error);
    let message = 'Failed to export attendance to Excel.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  }
};

const formatDateTime = (dateTimeString, type = 'datetime') => {
  if (!dateTimeString) return '';
  const date = new Date(dateTimeString);
  const options = {
    date: { year: 'numeric', month: 'long', day: 'numeric' },
    time: { hour: '2-digit', minute: '2-digit', second: '2-digit' },
    datetime: { year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit', second: '2-digit' }
  };
  return date.toLocaleString('id-ID', options[type]);
};

onMounted(() => {
  fetchAttendanceHistory();
});
</script>

<style scoped>
/* Add any specific styles for this component here */
</style>