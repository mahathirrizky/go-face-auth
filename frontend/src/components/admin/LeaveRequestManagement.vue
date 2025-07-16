<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-6">Manajemen Pengajuan Cuti & Izin</h2>

    <div class="bg-bg-muted p-4 rounded-lg shadow-md mb-6">
      <h3 class="text-xl font-semibold text-text-base mb-4">Filter Pengajuan</h3>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label for="filterStatus" class="block text-text-muted text-sm font-bold mb-2">Status:</label>
          <Dropdown
            id="filterStatus"
            v-model="filterStatus"
            :options="[
              { label: 'Semua', value: '' },
              { label: 'Pending', value: 'pending' },
              { label: 'Disetujui', value: 'approved' },
              { label: 'Ditolak', value: 'rejected' }
            ]"
            optionLabel="label"
            optionValue="value"
            placeholder="Pilih Status"
            class="w-full"
          />
        </div>
        <BaseInput
          id="filterEmployee"
          label="Nama Karyawan:"
          v-model="filterEmployeeName"
          placeholder="Cari nama karyawan..."
        />
      </div>
      <BaseButton @click="fetchLeaveRequests" class="btn-primary mt-4"><i class="pi pi-filter"></i> Terapkan Filter</BaseButton>
    </div>

    <BaseDataTable
      :data="filteredLeaveRequests"
      :columns="leaveRequestColumns"
      :loading="isLoading"
      :globalFilterFields="['Employee.name', 'Type', 'Reason']"
      searchPlaceholder="Cari Pengajuan..."
    >
      <template #column-status="{ item }">
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
import { ref, onMounted, computed } from 'vue';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { useAuthStore } from '../../stores/auth';
import BaseInput from '../ui/BaseInput.vue';
import BaseButton from '../ui/BaseButton.vue';
import BaseDataTable from '../ui/BaseDataTable.vue';
import Dropdown from 'primevue/dropdown';

const leaveRequests = ref([]);
const filterStatus = ref('');
const filterEmployeeName = ref('');
const toast = useToast();
const authStore = useAuthStore();
const isLoading = ref(false);

const leaveRequestColumns = ref([
    { field: 'Employee.name', header: 'Karyawan' },
    { field: 'Type', header: 'Tipe' },
    { field: 'StartDate', header: 'Tanggal Mulai' },
    { field: 'EndDate', header: 'Tanggal Selesai' },
    { field: 'Reason', header: 'Alasan' },
    { field: 'Status', header: 'Status' },
    { field: 'actions', header: 'Aksi' }
]);

const fetchLeaveRequests = async () => {
  if (!authStore.companyId) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Company ID not available. Cannot fetch leave requests.', life: 3000 });
    return;
  }
  isLoading.value = true;
  try {
    const response = await axios.get(`/api/company-leave-requests`);
    if (response.data && response.data.status === 'success') {
      leaveRequests.value = response.data.data;
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
      fetchLeaveRequests();
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

const filteredLeaveRequests = computed(() => {
  if (!Array.isArray(leaveRequests.value)) {
    return [];
  }
  return leaveRequests.value.filter(request => {
    const matchesStatus = filterStatus.value === '' || request.Status === filterStatus.value;
    const matchesEmployeeName = filterEmployeeName.value === '' ||
                                request.Employee.name.toLowerCase().includes(filterEmployeeName.value.toLowerCase());
    return matchesStatus && matchesEmployeeName;
  });
});

onMounted(() => {
  fetchLeaveRequests();
});
</script>

<style scoped>
/* Tailwind handles styling */
</style>