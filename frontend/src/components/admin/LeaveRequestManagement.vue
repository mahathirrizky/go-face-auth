<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <Toast />
    <h2 class="text-2xl font-bold text-text-base mb-6">Manajemen Pengajuan Cuti & Izin</h2>

    <DataTable
      :value="leaveRequests"
      :loading="isLoading"
      :totalRecords="totalRecords"
      :lazy="true"
      v-model:filters="filters"
      @page="onPage"
      paginator
      :rows="10"
      :rowsPerPageOptions="[10, 25, 50]"
      currentPageReportTemplate="Menampilkan {first} sampai {last} dari {totalRecords} data"
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink CurrentPageReport RowsPerPageSelect"
      dataKey="ID"
      :globalFilterFields="['Employee.name']"
      @filter="onFilter"
    >
      <template #header>
        <div class="flex flex-wrap items-center justify-between gap-4">
            <IconField iconPosition="left">
                <InputIcon class="pi pi-search"></InputIcon>
                <InputText v-model="filters['global'].value" placeholder="Cari Nama Karyawan..." @keydown.enter="onFilter" fluid/>
            </IconField>
            <div class="flex flex-wrap items-center gap-2">
                <DatePicker v-model="startDate" dateFormat="yy-mm-dd" placeholder="Dari Tanggal" fluid/>
                <DatePicker v-model="endDate" dateFormat="yy-mm-dd" placeholder="Sampai Tanggal" fluid/>
                <Button @click="fetchLeaveRequests" icon="pi pi-filter" label="Filter" />
                <Button @click="exportLeaveRequestsToExcel" icon="pi pi-file-excel" label="Export" :loading="isExporting" class="p-button-secondary"/>
            </div>
        </div>
      </template>

      <template #empty>
        Tidak ada data ditemukan.
      </template>
      <template #loading>
        Memuat data...
      </template>

      <Column field="Employee.name" header="Karyawan" :sortable="true"></Column>
      <Column field="Type" header="Tipe" :sortable="true"></Column>
      <Column field="StartDate" header="Tanggal Mulai" :sortable="true">
        <template #body="{ data }">{{ new Date(data.StartDate).toLocaleDateString('id-ID') }}</template>
      </Column>
      <Column field="EndDate" header="Tanggal Selesai" :sortable="true">
        <template #body="{ data }">{{ new Date(data.EndDate).toLocaleDateString('id-ID') }}</template>
      </Column>
      <Column field="Reason" header="Alasan"></Column>
      <Column field="Status" header="Status" :sortable="true" :showFilterMenu="true">
        <template #body="{ data }">
          <Tag :value="data.Status" :severity="getStatusSeverity(data.Status)" />
        </template>
        <template #filter="{ filterModel, filterCallback }">
            <Select v-model="filterModel.value" @change="filterCallback()" :options="statusOptions" placeholder="Pilih Status" class="p-column-filter" :showClear="true" fluid>
            </Select>
        </template>
      </Column>
      <Column field="CancelledBy" header="Dibatalkan Oleh">
        <template #body="{ data }">
          <span v-if="data.Status === 'cancelled'">{{ data.CancelledByActorType === 'employee' ? 'Karyawan' : 'Admin' }}</span>
          <span v-else>-</span>
        </template>
      </Column>
      <Column header="Aksi" style="min-width: 20rem">
        <template #body="{ data }">
          <div v-if="data.Status === 'pending'" class="flex space-x-2">
            <Button @click="reviewLeaveRequest(data.ID, 'approved')" class="p-button-success p-button-sm" :loading="isReviewing" icon="pi pi-check" label="Setujui" />
            <Button @click="reviewLeaveRequest(data.ID, 'rejected')" class="p-button-danger p-button-sm" :loading="isReviewing" icon="pi pi-times" label="Tolak" />
            <Button @click="reviewLeaveRequest(data.ID, 'cancelled')" class="p-button-warning p-button-sm" :loading="isReviewing" icon="pi pi-ban" label="Batalkan" />
          </div>
          <div v-else-if="data.Status === 'approved'" class="flex space-x-2">
            <Button @click="adminCancelApprovedLeave(data.ID)" class="p-button-warning p-button-sm" :loading="isCancelling" icon="pi pi-ban" label="Batalkan Cuti" />
          </div>
          <span v-else class="text-text-muted">Sudah Ditinjau</span>
        </template>
      </Column>
    </DataTable>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { useAuthStore } from '../../stores/auth';
import { FilterMatchMode } from '@primevue/core/api';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import IconField from 'primevue/iconfield';
import InputIcon from 'primevue/inputicon';
import Select from 'primevue/select';
import Tag from 'primevue/tag';
import DatePicker from 'primevue/datepicker';
import Toast from 'primevue/toast';

const leaveRequests = ref([]);
const toast = useToast();
const authStore = useAuthStore();
const isLoading = ref(false);
const isReviewing = ref(false);
const isCancelling = ref(false);
const isExporting = ref(false);
const totalRecords = ref(0);
const lazyParams = ref({});

const formatToYYYYMMDD = (date) => {
  if (!date) return null;
  const d = new Date(date);
  const year = d.getFullYear();
  const month = String(d.getMonth() + 1).padStart(2, '0');
  const day = String(d.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
};

const today = new Date();
const firstDayOfMonth = new Date(today.getFullYear(), today.getMonth(), 1);

const startDate = ref(firstDayOfMonth);
const endDate = ref(today);

const filters = ref({
    'global': { value: null, matchMode: FilterMatchMode.CONTAINS },
    'Status': { value: null, matchMode: FilterMatchMode.EQUALS },
});

const statusOptions = ref(['pending', 'approved', 'rejected', 'cancelled']);

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
      startDate: formatToYYYYMMDD(startDate.value),
      endDate: formatToYYYYMMDD(endDate.value),
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
  isReviewing.value = true;
  try {
    const response = await axios.put(`/api/leave-requests/${id}/review`, { status });
    if (response.data && response.data.status === 'success') {
      toast.add({ severity: 'success', summary: 'Success', detail: `Pengajuan ${status === 'approved' ? 'disetujui' : (status === 'rejected' ? 'ditolak' : 'dibatalkan')}.`, life: 3000 });
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
  } finally {
    isReviewing.value = false;
  }
};

const adminCancelApprovedLeave = async (id) => {
  isCancelling.value = true;
  try {
    const response = await axios.put(`/api/leave-requests/${id}/admin-cancel`);
    if (response.data && response.data.status === 'success') {
      toast.add({ severity: 'success', summary: 'Success', detail: 'Cuti yang disetujui berhasil dibatalkan.', life: 3000 });
      fetchLeaveRequests(); // Refresh data
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data?.message || 'Gagal membatalkan cuti yang disetujui.', life: 3000 });
    }
  } catch (error) {
    console.error('Error cancelling approved leave request:', error);
    let message = 'Gagal membatalkan cuti yang disetujui.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.add({ severity: 'error', summary: 'Error', detail: message, life: 3000 });
  } finally {
    isCancelling.value = false;
  }
};

const onPage = (event) => {
    lazyParams.value = event;
    fetchLeaveRequests();
};

const onFilter = () => {
    lazyParams.value.page = 0;
    fetchLeaveRequests();
};

const exportLeaveRequestsToExcel = async () => {
  if (!authStore.companyId) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Company ID not available. Cannot export.', life: 3000 });
    return;
  }
  isExporting.value = true;
  try {
    const params = {
      status: filters.value.Status.value || '',
      search: filters.value.global.value || '',
      startDate: formatToYYYYMMDD(startDate.value),
      endDate: formatToYYYYMMDD(endDate.value),
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
  } finally {
    isExporting.value = false;
  }
};

const getStatusSeverity = (status) => {
  switch (status) {
    case 'pending': return 'warning';
    case 'approved': return 'success';
    case 'rejected': return 'danger';
    case 'cancelled': return 'info';
    default: return null;
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