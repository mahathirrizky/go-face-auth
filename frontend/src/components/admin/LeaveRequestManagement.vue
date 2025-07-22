<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-6">Manajemen Pengajuan Cuti & Izin</h2>

    <BaseDataTable
      :data="leaveRequests"
      :columns="leaveRequestColumns"
      :loading="isLoading"
      :totalRecords="totalRecords"
      :lazy="true"
      v-model:filters="filters"
      @page="onPage"
      @filter="onFilter"
      searchPlaceholder="Cari Nama Karyawan..."
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
          <BaseButton @click="fetchLeaveRequests" class="btn-primary"><i class="pi pi-filter"></i> Filter</BaseButton>
          <BaseButton @click="exportLeaveRequestsToExcel" class="btn-secondary whitespace-nowrap"><i class="pi pi-file-excel"></i> Export to Excel</BaseButton>
        </div>
      </template>

      <template #filter-Status="{ filterModel }">
        <Select
          v-model="filterModel.value"
          :options="statusOptions"
          optionLabel="label"
          optionValue="value"
          placeholder="Pilih Status"
          class="p-column-filter w-full"
          :showClear="true"
        />
      </template>

      <template #column-Status="{ item }">
        <span :class="{
          'px-2 inline-flex text-xs leading-5 font-semibold rounded-full': true,
          'bg-yellow-100 text-yellow-800': item.Status === 'pending',
          'bg-green-100 text-green-800': item.Status === 'approved',
          'bg-red-100 text-red-800': item.Status === 'rejected',
        }">
          {{ item.Status }}
        </span>
      </template>

      <template #column-actions="{ item }">
        <div v-if="item.Status === 'pending'" class="flex space-x-2">
          <BaseButton @click="reviewLeaveRequest(item.ID, 'approved')" class="btn-success btn-sm"><i class="pi pi-check"></i> Setujui</BaseButton>
          <BaseButton @click="reviewLeaveRequest(item.ID, 'rejected')" class="btn-danger btn-sm"><i class="pi pi-times"></i> Tolak</BaseButton>
        </div>
        <span v-else class="text-text-muted">Sudah Ditinjau</span>
      </template>
    </BaseDataTable>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { useAuthStore } from '../../stores/auth';
import BaseButton from '../ui/BaseButton.vue';
import BaseDataTable from '../ui/BaseDataTable.vue';
import Select from 'primevue/select';
import BaseInput from '../ui/BaseInput.vue'; // Import BaseInput
import { FilterMatchMode } from '@primevue/core/api';

const leaveRequests = ref([]);
const toast = useToast();
const authStore = useAuthStore();
const isLoading = ref(false);
const totalRecords = ref(0);
const lazyParams = ref({});

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

const filters = ref({
    'global': { value: null, matchMode: FilterMatchMode.CONTAINS },
    'Status': { value: null, matchMode: FilterMatchMode.EQUALS },
});

const statusOptions = ref([
  { label: 'Pending', value: 'pending' },
  { label: 'Disetujui', value: 'approved' },
  { label: 'Ditolak', value: 'rejected' }
]);

const leaveRequestColumns = ref([
    { field: 'Employee.name', header: 'Karyawan', showFilterMenu: false },
    { field: 'Type', header: 'Tipe', showFilterMenu: false },
    { field: 'StartDate', header: 'Tanggal Mulai', showFilterMenu: false },
    { field: 'EndDate', header: 'Tanggal Selesai', showFilterMenu: false },
    { field: 'Reason', header: 'Alasan', showFilterMenu: false },
    { field: 'Status', header: 'Status', filter: true, showFilterMenu: true },
    { field: 'actions', header: 'Aksi', showFilterMenu: false, sortable: false }
]);

const fetchLeaveRequests = async () => {
  if (!authStore.companyId) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Company ID not available.', life: 3000 });
    return;
  }
  isLoading.value = true;
  try {
    const params = {
      page: lazyParams.value.page + 1,
      limit: lazyParams.value.rows,
      status: filters.value.Status.value || '',
      search: filters.value.global.value || '',
      startDate: startDate.value,
      endDate: endDate.value,
    };

    const response = await axios.get('/api/company-leave-requests', { params });
    
    if (response.data && response.data.status === 'success') {
      leaveRequests.value = response.data.data.items;
      totalRecords.value = response.data.data.total_records;
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data?.message || 'Failed to fetch leave requests.', life: 3000 });
    }
  } catch (error) {
    console.error('Error fetching leave requests:', error);
    let message = 'Failed to fetch leave requests.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  } finally {
    isLoading.value = false;
  }
};

const reviewLeaveRequest = async (id, status) => {
  try {
    const response = await axios.put(`/api/leave-requests/${id}/review`, { status });
    if (response.data && response.data.status === 'success') {
      toast.add({ severity: 'success', summary: 'Success', detail: `Pengajuan ${status === 'approved' ? 'disetujui' : 'ditolak'}.`, life: 3000 });
      fetchLeaveRequests(); // Refresh data
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data?.message || 'Gagal meninjau pengajuan.', life: 3000 });
    }
  } catch (error) {
    console.error('Error reviewing leave request:', error);
    let message = 'Gagal meninjau pengajuan.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  }
};

const onPage = (event) => {
    lazyParams.value = event;
    fetchLeaveRequests();
};

const onFilter = () => {
    // The v-model:filters binding handles the state update.
    // We just need to trigger a refetch.
    fetchLeaveRequests();
};

const exportLeaveRequestsToExcel = async () => {
  if (!authStore.companyId) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Company ID not available. Cannot export.', life: 3000 });
    return;
  }
  try {
    const params = {
      status: filters.value.Status.value || '',
      search: filters.value.global.value || ''
    };

    const response = await axios.get(`/api/company-leave-requests/export`, {
      params,
      responseType: 'blob',
    });

    const blob = new Blob([response.data], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = `company_leave_requests.xlsx`;
    link.click();
    URL.revokeObjectURL(link.href);
    toast.add({ severity: 'success', summary: 'Success', detail: 'File Excel pengajuan cuti berhasil diunduh!', life: 3000 });
  } catch (error) {
    console.error('Error exporting leave requests to Excel:', error);
    let message = 'Failed to export leave requests to Excel.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  }
};

onMounted(() => {
  lazyParams.value = {
    first: 0,
    rows: 10,
    page: 0,
  };
  fetchLeaveRequests();
});
</script>

<style scoped>
/* Tailwind handles styling */
</style>
