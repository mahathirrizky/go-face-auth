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
          Karyawan Tidak Absen
        </button>
        <button
          @click="selectedTab = 'overtime'"
          :class="[
            'whitespace-nowrap py-3 px-1 border-b-2 font-medium text-sm',
            selectedTab === 'overtime'
              ? 'border-secondary text-secondary'
              : 'border-transparent text-text-muted hover:text-text-base hover:border-gray-300',
          ]"
        >
          Lembur
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
        <button @click="exportAllToExcel" class="btn btn-secondary w-full md:w-auto mt-4 md:mt-0">Export to Excel</button>
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

    <!-- Tab Content: Karyawan Tidak Absen -->
    <div v-if="selectedTab === 'unaccounted'">
      <div class="bg-bg-muted p-4 rounded-lg shadow-md mb-6 flex flex-col md:flex-row justify-between items-center">
        <div class="flex items-center space-x-2">
          <label for="unaccountedStartDate" class="text-text-muted">Dari:</label>
          <input
            type="date"
            id="unaccountedStartDate"
            v-model="unaccountedStartDate"
            class="p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary"
          />
        </div>
        <div class="flex items-center space-x-2">
          <label for="unaccountedEndDate" class="text-text-muted">Sampai:</label>
          <input
            type="date"
            id="unaccountedEndDate"
            v-model="unaccountedEndDate"
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
            <tr v-if="!Array.isArray(unaccountedEmployees) || unaccountedEmployees.length === 0">
              <td colspan="3" class="px-6 py-4 text-center text-text-muted">Tidak ada karyawan tidak absen untuk rentang tanggal ini.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Tab Content: Lembur -->
    <div v-if="selectedTab === 'overtime'">
      <div class="bg-bg-muted p-4 rounded-lg shadow-md mb-6 flex flex-col md:flex-row justify-between items-center">
        <div class="flex items-center space-x-2">
          <label for="overtimeStartDate" class="text-text-muted">Dari:</label>
          <input
            type="date"
            id="overtimeStartDate"
            v-model="overtimeStartDate"
            class="p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary"
          />
        </div>
        <div class="flex items-center space-x-2">
          <label for="overtimeEndDate" class="text-text-muted">Sampai:</label>
          <input
            type="date"
            id="overtimeEndDate"
            v-model="overtimeEndDate"
            class="p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary"
          />
        </div>
        <button @click="fetchOvertimeAttendances" class="btn btn-primary w-full md:w-auto mt-4 md:mt-0">Cari</button>
      </div>

      <div class="overflow-x-auto bg-bg-muted rounded-lg shadow-md">
        <table class="min-w-full divide-y divide-bg-base">
          <thead class="bg-primary">
            <tr>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Nama Karyawan</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Waktu Masuk Lembur</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Waktu Keluar Lembur</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Durasi Lembur (Menit)</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-bg-base">
            <tr v-for="record in overtimeRecords" :key="record.id">
              <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ record.Employee.name }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ new Date(record.check_in_time).toLocaleString() }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ record.check_out_time ? new Date(record.check_out_time).toLocaleString() : '-' }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ record.overtime_minutes || 0 }}</td>
            </tr>
            <tr v-if="!Array.isArray(overtimeRecords) || overtimeRecords.length === 0">
              <td colspan="4" class="px-6 py-4 text-center text-text-muted">Tidak ada data lembur untuk rentang tanggal ini.</td>
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
    const overtimeRecords = ref([]); // New ref for overtime records
    const selectedTab = ref('all'); // Default to 'all' tab
    const toast = useToast();
    const authStore = useAuthStore();

    // Helper to format date to YYYY-MM-DD
    const formatToYYYYMMDD = (date) => {
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      return `${year}-${month}-${day}`;
    };

    // Calculate start of current month and today's date
    const today = new Date();
    const firstDayOfMonth = new Date(today.getFullYear(), today.getMonth(), 1);

    const startDate = ref(formatToYYYYMMDD(firstDayOfMonth));
    const endDate = ref(formatToYYYYMMDD(today));
    const unaccountedStartDate = ref(formatToYYYYMMDD(firstDayOfMonth));
    const unaccountedEndDate = ref(formatToYYYYMMDD(today));
    const overtimeStartDate = ref(formatToYYYYMMDD(firstDayOfMonth)); // New date ref for overtime
    const overtimeEndDate = ref(formatToYYYYMMDD(today)); // New date ref for overtime

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
          params: { 
            startDate: unaccountedStartDate.value,
            endDate: unaccountedEndDate.value
          },
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

    const fetchOvertimeAttendances = async () => {
      if (!authStore.companyId) {
        toast.error('Company ID not available. Cannot fetch overtime attendances.');
        return;
      }
      try {
        const response = await axios.get(`/api/overtime-attendances`, {
          params: {
            startDate: overtimeStartDate.value,
            endDate: overtimeEndDate.value,
          },
        });
        if (response.data && response.data.status === 'success') {
          overtimeRecords.value = response.data.data;
        } else {
          toast.error(response.data?.message || 'Failed to fetch overtime attendances.');
        }
      } catch (error) {
        console.error('Error fetching overtime attendances:', error);
        let message = 'Failed to fetch overtime attendances.';
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
      fetchOvertimeAttendances(); // Fetch initial data for the overtime tab
    });

    return {
      attendanceRecords,
      unaccountedEmployees,
      overtimeRecords,
      selectedTab,
      startDate,
      endDate,
      unaccountedStartDate,
      unaccountedEndDate,
      overtimeStartDate,
      overtimeEndDate,
      fetchAttendances,
      fetchUnaccountedEmployees,
      fetchOvertimeAttendances,
      exportAllToExcel,
    };
  },
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>

