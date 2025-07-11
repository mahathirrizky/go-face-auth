<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-6">Manajemen Absensi</h2>

    <!-- Tab Navigation -->
    <div class="mb-6 border-b border-bg-base">
      <nav class="-mb-px flex space-x-8" aria-label="Tabs">
        <button
          @click="selectedTab = 'all'"
          :class="[
            'whitespace-nowrap py-3 px-1 border-b-2 font-medium text-sm',
            selectedTab === 'all'
              ? 'border-secondary text-secondary'
              : 'border-transparent text-text-muted hover:text-text-base hover:border-gray-300',
          ]"
        >
          Semua Absensi
        </button>
        <button
          @click="selectedTab = 'unaccounted'"
          :class="[
            'whitespace-nowrap py-3 px-1 border-b-2 font-medium text-sm',
            selectedTab === 'unaccounted'
              ? 'border-secondary text-secondary'
              : 'border-transparent text-text-muted hover:text-text-base hover:border-gray-300',
          ]"
        >
          Karyawan Tidak Terdata
        </button>
      </nav>
    </div>

    <!-- Tab Content: Semua Absensi -->
    <div v-if="selectedTab === 'all'">
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
          <button @click="fetchAttendances" class="btn btn-primary w-full md:w-auto">Filter</button>
        </div>
        <button @click="exportAllToExcel" class="btn btn-secondary w-full md:w-auto mt-4 md:mt-0">Export Semua ke Excel</button>
      </div>

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
                  'bg-green-100 text-green-800': record.status === 'on_time',
                  'bg-yellow-100 text-yellow-800': record.status === 'late',
                  'bg-blue-100 text-blue-800': record.status === 'overtime_in' || record.status === 'overtime_out',
                }">
                  {{ record.status === 'on_time' ? 'Tepat Waktu' : record.status === 'late' ? 'Terlambat' : record.status === 'overtime_in' ? 'Lembur Masuk' : record.status === 'overtime_out' ? 'Lembur Keluar' : record.status }}
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

    <!-- Tab Content: Karyawan Tidak Terdata -->
    <div v-if="selectedTab === 'unaccounted'">
      <div class="bg-bg-muted p-4 rounded-lg shadow-md mb-6 flex flex-col md:flex-row justify-between items-center">
        <div class="flex items-center space-x-2">
          <label for="unaccountedDate" class="text-text-muted">Tanggal:</label>
          <input
            type="date"
            id="unaccountedDate"
            v-model="unaccountedDate"
            class="p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary"
          />
        </div>
        <button @click="fetchUnaccountedEmployees" class="btn btn-primary w-full md:w-auto mt-4 md:mt-0">Cari</button>
      </div>

      <div class="overflow-x-auto bg-bg-muted rounded-lg shadow-md">
        <table class="min-w-full divide-y divide-bg-base">
          <thead class="bg-primary">
            <tr>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Nama Karyawan</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Email</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Posisi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-bg-base">
            <tr v-for="employee in unaccountedEmployees" :key="employee.id">
              <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ employee.name }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ employee.email }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ employee.position }}</td>
            </tr>
            <tr v-if="unaccountedEmployees.length === 0">
              <td colspan="3" class="px-6 py-4 text-center text-text-muted">Tidak ada karyawan tidak terdata untuk tanggal ini.</td>
            </tr>
          </tbody>
        </table>
      </div>
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
    const unaccountedEmployees = ref([]);
    const selectedTab = ref('all'); // Default to 'all' tab
    const toast = useToast();
    const authStore = useAuthStore();

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
    const unaccountedDate = ref(formatToYYYYMMDD(today)); // Default to today for unaccounted employees

    const fetchAttendances = async () => {
      if (!authStore.companyId) {
        toast.error('Company ID not available. Cannot fetch attendances.');
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

        const response = await axios.get(`/api/attendances`, { params });
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

    const fetchUnaccountedEmployees = async () => {
      if (!authStore.companyId) {
        toast.error('Company ID not available. Cannot fetch unaccounted employees.');
        return;
      }
      try {
        const response = await axios.get(`/api/unaccounted-employees`, {
          params: { date: unaccountedDate.value },
        });
        if (response.data && response.data.status === 'success') {
          unaccountedEmployees.value = response.data.data;
        } else {
          toast.error(response.data?.message || 'Failed to fetch unaccounted employees.');
        }
      } catch (error) {
        console.error('Error fetching unaccounted employees:', error);
        let message = 'Failed to fetch unaccounted employees.';
        if (error.response && error.response.data && error.response.data.message) {
          message = error.response.data.message;
        }
        toast.error(message);
      }
    };

    const exportAllToExcel = async () => {
      try {
        let url = `/api/attendances/export`;
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
        link.download = `all_company_attendance.xlsx`;
        link.click();
        URL.revokeObjectURL(link.href);
        toast.success('File Excel semua absensi berhasil diunduh!');
      } catch (error) {
        console.error('Error exporting all attendances to Excel:', error);
        let message = 'Failed to export all attendances to Excel.';
        if (error.response && error.response.data && error.response.data.message) {
          message = error.response.data.message;
        }
        toast.error(message);
      }
    };

    onMounted(() => {
      fetchAttendances();
      fetchUnaccountedEmployees(); // Fetch initial data for the unaccounted tab
    });

    return {
      attendanceRecords,
      unaccountedEmployees,
      selectedTab,
      startDate,
      endDate,
      unaccountedDate,
      fetchAttendances,
      fetchUnaccountedEmployees,
      exportAllToExcel,
    };
  },
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>
