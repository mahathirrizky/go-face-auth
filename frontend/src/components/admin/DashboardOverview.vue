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
<script>
import { ref, onMounted, onUnmounted } from 'vue';
import axios from 'axios';
import { useToast } from 'vue-toastification';
import { useAuthStore } from '../../stores/auth'; // Import auth store

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
    const authStore = useAuthStore(); // Initialize auth store
    let ws = null; // WebSocket instance

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

    const connectWebSocket = () => {
      if (!authStore.token) {
        console.warn('No auth token found, cannot establish WebSocket connection.');
        return;
      }

      // Determine WebSocket URL based on current location
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const host = window.location.host.split(':')[0]; // Get hostname without port
      const port = window.location.port ? ':' + window.location.port : '';
      const wsUrl = `${protocol}//${host}${port}/ws/dashboard`;

      ws = new WebSocket(wsUrl);

      ws.onopen = () => {
        console.log('WebSocket connected.');
        // Send auth token immediately after connection is open
        ws.send(JSON.stringify({ token: authStore.token }));
      };

      ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        if (data) {
          summary.value = data; // Update summary with WebSocket data
        }
      };

      ws.onclose = (event) => {
        console.log('WebSocket disconnected:', event.code, event.reason);
        // Attempt to reconnect after a delay if connection was not closed cleanly
        if (event.code !== 1000) { // 1000 is normal closure
          setTimeout(() => {
            console.log('Attempting to reconnect WebSocket...');
            connectWebSocket();
          }, 3000); // Reconnect after 3 seconds
        }
      };

      ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        toast.error('WebSocket connection error. Dashboard updates may be delayed.');
      };
    };

    onMounted(() => {
      fetchDashboardSummary(); // Initial fetch via HTTP
      connectWebSocket(); // Establish WebSocket connection
    });

    onUnmounted(() => {
      if (ws) {
        ws.close(); // Close WebSocket connection when component is unmounted
      }
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
