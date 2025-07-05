<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-6">Manajemen Absensi</h2>

    <div class="overflow-x-auto bg-bg-muted rounded-lg shadow-md">
      <table class="min-w-full divide-y divide-bg-base">
        <thead class="bg-primary">
          <tr>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Tanggal</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Nama Karyawan</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Waktu Masuk</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Waktu Keluar</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Status</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-bg-base">
          <tr v-for="record in attendanceRecords" :key="record.id">
            <td class="px-6 py-4 whitespace-nowrap text-text-base">{{ new Date(record.check_in_time).toLocaleDateString() }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ record.Employee.name }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ new Date(record.check_in_time).toLocaleTimeString() }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ record.check_out_time ? new Date(record.check_out_time).toLocaleTimeString() : '-' }}</td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span :class="{
                'px-2 inline-flex text-xs leading-5 font-semibold rounded-full': true,
                'bg-green-100 text-green-800': record.status === 'present',
                'bg-red-100 text-red-800': record.status === 'absent',
                'bg-yellow-100 text-yellow-800': record.status === 'leave'
              }">
                {{ record.status }}
              </span>
            </td>
          </tr>
          <tr v-if="attendanceRecords.length === 0">
            <td colspan="5" class="px-6 py-4 text-center text-text-muted">Tidak ada data absensi.</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { useToast } from 'vue-toastification';
import { useAuthStore } from '../../stores/auth';

export default {
  name: 'AttendanceManagement',
  setup() {
    const attendanceRecords = ref([]);
    const toast = useToast();
    const authStore = useAuthStore();

    const fetchAttendances = async () => {
      if (!authStore.companyId) {
        toast.error('Company ID not available. Cannot fetch attendances.');
        return;
      }
      try {
        const response = await axios.get(`/api/attendances`);
        if (response.data && response.data.status === 'success') {
          attendanceRecords.value = response.data.data;
        } else {
          toast.error(response.data?.message || 'Failed to fetch attendances.');
        }
      } catch (error) {
        console.error('Error fetching attendances:', error);
        let message = 'Failed to fetch attendances.';
        if (error.response && error.response.data && error.response.data.message) {
          message = error.response.data.message;
        }
        toast.error(message);
      }
    };

    onMounted(() => {
      fetchAttendances();
    });

    return {
      attendanceRecords,
    };
  },
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>
