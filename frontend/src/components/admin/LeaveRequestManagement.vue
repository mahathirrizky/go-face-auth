<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-6">Manajemen Pengajuan Cuti & Izin</h2>

    <div class="bg-bg-muted p-4 rounded-lg shadow-md mb-6">
      <h3 class="text-xl font-semibold text-text-base mb-4">Filter Pengajuan</h3>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label for="filterStatus" class="block text-text-muted text-sm font-bold mb-2">Status:</label>
          <select
            id="filterStatus"
            v-model="filterStatus"
            class="shadow appearance-none border rounded w-full py-2 px-3 text-text-base leading-tight focus:outline-none focus:shadow-outline bg-bg-base border-bg-base"
          >
            <option value="">Semua</option>
            <option value="pending">Pending</option>
            <option value="approved">Disetujui</option>
            <option value="rejected">Ditolak</option>
          </select>
        </div>
        <BaseInput
          id="filterEmployee"
          label="Nama Karyawan:"
          v-model="filterEmployeeName"
          placeholder="Cari nama karyawan..."
        />
      </div>
      <BaseButton @click="fetchLeaveRequests" class="btn-primary mt-4"><i class="fas fa-filter"></i> Terapkan Filter</BaseButton>
    </div>

    <div class="overflow-x-auto bg-bg-muted rounded-lg shadow-md">
      <table class="min-w-full divide-y divide-bg-base">
        <thead class="bg-primary">
          <tr>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Karyawan</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Tipe</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Tanggal Mulai</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Tanggal Selesai</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Alasan</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Status</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-bg-base">
          <tr v-for="request in filteredLeaveRequests" :key="request.ID">
            <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ request.Employee.name }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ request.Type }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ new Date(request.StartDate).toLocaleDateString() }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ new Date(request.EndDate).toLocaleDateString() }}</td>
            <td class="px-6 py-4 text-text-muted">{{ request.Reason }}</td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span :class="{
                'px-2 inline-flex text-xs leading-5 font-semibold rounded-full': true,
                'bg-yellow-100 text-yellow-800': request.Status === 'pending',
                'bg-green-100 text-green-800': request.Status === 'approved',
                'bg-red-100 text-red-800': request.Status === 'rejected',
              }">
                {{ request.Status }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
              <div v-if="request.Status === 'pending'" class="flex space-x-2">
                <BaseButton @click="reviewLeaveRequest(request.ID, 'approved')" class="btn-success btn-sm"><i class="fas fa-check"></i> Setujui</BaseButton>
                <BaseButton @click="reviewLeaveRequest(request.ID, 'rejected')" class="btn-danger btn-sm"><i class="fas fa-times"></i> Tolak</BaseButton>
              </div>
              <span v-else class="text-text-muted">Sudah Ditinjau</span>
            </td>
          </tr>
          <tr v-if="filteredLeaveRequests.length === 0">
            <td colspan="7" class="px-6 py-4 text-center text-text-muted">Tidak ada pengajuan cuti/izin.</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import axios from 'axios';
import { useToast } from 'vue-toastification';
import { useAuthStore } from '../../stores/auth';
import BaseInput from '../ui/BaseInput.vue';
import BaseButton from '../ui/BaseButton.vue';

const leaveRequests = ref([]);
const filterStatus = ref('');
const filterEmployeeName = ref('');
const toast = useToast();
const authStore = useAuthStore();

const fetchLeaveRequests = async () => {
  if (!authStore.companyId) {
    toast.error('Company ID not available. Cannot fetch leave requests.');
    return;
  }
  try {
    const response = await axios.get(`/api/company-leave-requests`);
    if (response.data && response.data.status === 'success') {
      leaveRequests.value = response.data.data;
    } else {
      toast.error(response.data?.message || 'Failed to fetch leave requests.');
    }
  } catch (error) {
    console.error('Error fetching leave requests:', error);
    let message = 'Failed to fetch leave requests.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.error(message);
  }
};

const reviewLeaveRequest = async (id, status) => {
  try {
    const response = await axios.put(`/api/leave-requests/${id}/review`, { status });
    if (response.data && response.data.status === 'success') {
      toast.success(`Pengajuan ${status === 'approved' ? 'disetujui' : 'ditolak'}.`);
      fetchLeaveRequests();
    } else {
      toast.error(response.data?.message || 'Gagal meninjau pengajuan.');
    }
  } catch (error) {
    console.error('Error reviewing leave request:', error);
    let message = 'Gagal meninjau pengajuan.';
    if (error.response && error.response.data && error.response.data.message) {
      message = error.response.data.message;
    }
    toast.error(message);
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
