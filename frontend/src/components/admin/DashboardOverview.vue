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
        <li v-for="activity in recentActivities" :key="activity.timestamp" class="border-b border-bg-base py-2 text-text-muted">
          {{ activity.description }} - {{ new Date(activity.timestamp).toLocaleString() }}
        </li>
        <li v-if="recentActivities.length === 0" class="py-2 text-text-muted">Tidak ada aktivitas terbaru.</li>
      </ul>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted } from 'vue';
import axios from 'axios';
import { useToast } from 'vue-toastification';
import { useWebSocketStore } from '../../stores/websocket'; // Import WebSocket store

export default {
  name: 'DashboardOverview',
  setup() {
    const summary = ref({
      total_employees: 0,
      present_today: 0,
      absent_today: 0,
      on_leave_today: 0,
    });
    const recentActivities = ref([]); // New ref for recent activities
    const toast = useToast();
    const webSocketStore = useWebSocketStore(); // Initialize WebSocket store

    const fetchDashboardSummary = async () => {
      try {
        const response = await axios.get('/api/dashboard-summary');
        if (response.data && response.data.data) {
          summary.value = response.data.data;
          console.log('DEBUG: DashboardOverview - Received recent_activities:', response.data.data.recent_activities);
          recentActivities.value = response.data.data.recent_activities || [];
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

    // Handler for WebSocket messages
    const handleWebSocketMessage = (data) => {
      if (data) {
        summary.value = data; // Update summary with WebSocket data
        recentActivities.value = data.recent_activities || []; // Update recent activities
      }
    };

    onMounted(() => {
      fetchDashboardSummary(); // Initial fetch via HTTP
      // Register WebSocket message handler
      webSocketStore.onMessage('dashboard_update', handleWebSocketMessage);
    });

    onUnmounted(() => {
      // Unregister WebSocket message handler
      webSocketStore.offMessage('dashboard_update');
    });

    return {
      summary,
      recentActivities, // Return recentActivities
    };
  },
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>
