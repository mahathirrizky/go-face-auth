<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-6">Riwayat Absensi Karyawan</h2>

    <div class="bg-bg-muted p-4 rounded-lg shadow-md mb-6 flex flex-col md:flex-row justify-between items-center">
      <div class="flex flex-col md:flex-row items-center space-y-4 md:space-y-0 md:space-x-4 w-full">
        <div class="flex items-center space-x-2">
          <label for="startDate" class="text-text-muted">Dari:</label>
          <input
            type="date"
            id="startDate"
            v-model="startDate"
            class="p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary"
          />
        </div>
        <div class="flex items-center space-x-2">
          <label for="endDate" class="text-text-muted">Sampai:</label>
          <input
            type="date"
            id="endDate"
            v-model="endDate"
            class="p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary"
          />
        </div>
        <button @click="fetchAttendanceHistory" class="btn btn-primary w-full md:w-auto">Filter</button>
      </div>
      <button @click="exportToExcel" class="btn btn-secondary w-full md:w-auto mt-4 md:mt-0">Export ke Excel</button>
    </div>

    <div class="overflow-x-auto bg-bg-muted rounded-lg shadow-md">
      <table class="min-w-full divide-y divide-bg-base">
        <thead class="bg-primary">
          <tr>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Nama Karyawan</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Waktu Check-in</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Waktu Check-out</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Status</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-bg-base">
          <tr v-for="attendance in attendances" :key="attendance.id">
            <td class="px-6 py-4 whitespace-nowrap text-text-base">{{ attendance.Employee ? attendance.Employee.name : 'N/A' }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ formatDateTime(attendance.check_in_time) }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ attendance.check_out_time ? formatDateTime(attendance.check_out_time) : 'Belum Check-out' }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ attendance.status }}</td>
          </tr>
          <tr v-if="attendances.length === 0">
            <td colspan="4" class="px-6 py-4 text-center text-text-muted">Tidak ada riwayat absensi untuk periode ini.</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import axios from 'axios';
import { useToast } from 'vue-toastification';

export default {
  name: 'EmployeeAttendanceHistory',
  setup() {
    const route = useRoute();
    const toast = useToast();
    const employeeId = ref(route.params.employeeId);
    const attendances = ref([]);

    // Calculate start and end of current month
    const today = new Date();
    const firstDayOfMonth = new Date(today.getFullYear(), today.getMonth(), 1);
    const lastDayOfMonth = new Date(today.getFullYear(), today.getMonth() + 1, 0);

    // Format dates to YYYY-MM-DD for input type="date"
    const formatToYYYYMMDD = (date) => {
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      return `${year}-${month}-${day}`;
    };

    const startDate = ref(formatToYYYYMMDD(firstDayOfMonth));
    const endDate = ref(formatToYYYYMMDD(lastDayOfMonth));

    const fetchAttendanceHistory = async () => {
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
          toast.error(response.data?.message || 'Failed to fetch attendance history.');
        }
      } catch (error) {
        console.error('Error fetching attendance history:', error);
        let message = 'Failed to fetch attendance history.';
        if (error.response && error.response.data && error.response.data.message) {
          message = error.response.data.message;
        }
        toast.error(message);
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
          responseType: 'blob', // Important for downloading files
        });

        const blob = new Blob([response.data], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' });
        const link = document.createElement('a');
        link.href = URL.createObjectURL(blob);
        link.download = `riwayat_absensi_karyawan_${employeeId.value}.xlsx`;
        link.click();
        URL.revokeObjectURL(link.href);
        toast.success('File Excel berhasil diunduh!');
      } catch (error) {
        console.error('Error exporting to Excel:', error);
        let message = 'Failed to export attendance to Excel.';
        if (error.response && error.response.data && error.response.data.message) {
          message = error.response.data.message;
        }
        toast.error(message);
      }
    };

    const formatDateTime = (dateTimeString) => {
      if (!dateTimeString) return '';
      const date = new Date(dateTimeString);
      return date.toLocaleString('id-ID', { year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit', second: '2-digit' });
    };

    onMounted(() => {
      fetchAttendanceHistory();
    });

    return {
      employeeId,
      attendances,
      startDate,
      endDate,
      fetchAttendanceHistory,
      exportToExcel,
      formatDateTime,
    };
  },
};
</script>

<style scoped>
/* Add any specific styles for this component here */
</style>
