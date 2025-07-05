<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-4">Ringkasan Dashboard</h2>
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <!-- Total Karyawan Card -->
      <div class="bg-bg-muted p-6 rounded-lg shadow-md">
        <h3 class="text-lg font-semibold text-text-base mb-2">Total Karyawan</h3>
        <p class="text-3xl font-bold text-secondary">{{ summary.total_employees }}</p>
      </div>
      <!-- Absensi Hari Ini Card -->
      <div class="bg-bg-muted p-6 rounded-lg shadow-md">
        <h3 class="text-lg font-semibold text-text-base mb-2">Absensi Hari Ini</h3>
        <p class="text-3xl font-bold text-accent">{{ summary.present_today }} Hadir</p>
      </div>
      <!-- Izin/Cuti Card -->
      <div class="bg-bg-muted p-6 rounded-lg shadow-md">
        <h3 class="text-lg font-semibold text-text-base mb-2">Izin/Cuti</h3>
        <p class="text-3xl font-bold text-danger">{{ summary.on_leave_today }}</p>
      </div>
    </div>

    <div class="mt-8 bg-bg-muted p-6 rounded-lg shadow-md">
      <h3 class="text-xl font-semibold text-text-base mb-4">Aktivitas Terbaru</h3>
      <ul>
        <li class="border-b border-bg-base py-2 text-text-muted">Karyawan John Doe absen masuk pukul 08:00</li>
        <li class="border-b border-bg-base py-2 text-text-muted">Karyawan Jane Smith absen keluar pukul 17:00</li>
        <li class="py-2 text-text-muted">Laporan absensi bulan Juni telah dibuat</li>
      </ul>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { useToast } from 'vue-toastification';

export default {
  name: 'DashboardOverview',
  setup() {
    const summary = ref({
      total_employees: 0,
      present_today: 0,
      absent_today: 0,
      on_leave_today: 0,
    });
    const toast = useToast();

    const fetchDashboardSummary = async () => {
      try {
        const response = await axios.get('/api/dashboard-summary');
        if (response.data && response.data.data) {
          summary.value = response.data.data;
        } else {
          toast.error('Failed to fetch dashboard summary.');
        }
      } catch (error) {
        console.error('Error fetching dashboard summary:', error);
        let message = 'Failed to load dashboard summary.';
        if (error.response && error.response.data && error.response.data.message) {
          message = error.response.data.message;
        }
        toast.error(message);
      }
    };

    onMounted(() => {
      fetchDashboardSummary();
    });

    return {
      summary,
    };
  },
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>
