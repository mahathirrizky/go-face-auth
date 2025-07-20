<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-4">Ringkasan Dashboard SuperAdmin</h2>

    <div class="flex flex-col lg:flex-row lg:space-x-6 mt-8">
      <!-- Chart on the left -->
      <div class="bg-bg-muted p-6 rounded-lg shadow-md mb-6 lg:mb-0 lg:w-1/2 max-h-96">
        <h3 class="text-xl font-semibold text-text-base mb-4">Pendapatan Bulanan</h3>
        <Chart type="line" :data="chartData" :options="chartOptions" />
      </div>

      <!-- Cards on the right -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-2 gap-6 lg:w-1/2">
        <!-- Total Companies Card -->
        <div class="bg-bg-muted p-6 rounded-lg shadow-md">
          <h3 class="text-lg font-semibold text-text-base mb-2">Total Perusahaan</h3>
          <p class="text-3xl font-bold text-secondary">{{ summary.total_companies }}</p>
        </div>
        <!-- Active Subscriptions Card -->
        <div class="bg-bg-muted p-6 rounded-lg shadow-md">
          <h3 class="text-lg font-semibold text-text-base mb-2">Langganan Aktif</h3>
          <p class="text-3xl font-bold text-accent">{{ summary.active_subscriptions }}</p>
        </div>
        <!-- Expired Subscriptions Card -->
        <div class="bg-bg-muted p-6 rounded-lg shadow-md">
          <h3 class="text-lg font-semibold text-text-base mb-2">Langganan Kedaluwarsa</h3>
          <p class="text-3xl font-bold text-danger">{{ summary.expired_subscriptions }}</p>
        </div>
        <!-- Trial Subscriptions Card -->
        <div class="bg-bg-muted p-6 rounded-lg shadow-md">
          <h3 class="text-lg font-semibold text-text-base mb-2">Langganan Percobaan</h3>
          <p class="text-3xl font-bold text-blue-400">{{ summary.trial_subscriptions }}</p>
        </div>
      </div>
    </div>

    <div class="mt-8 bg-bg-muted p-6 rounded-lg shadow-md">
      <h3 class="text-xl font-semibold text-text-base mb-4">Aktivitas Terbaru</h3>
      <ul>
        <li v-for="activity in recentActivities" :key="activity.id" class="border-b border-bg-base py-2 text-gray-300">
          {{ activity.description }} - {{ new Date(activity.timestamp).toLocaleString() }}
        </li>
        <li v-if="recentActivities.length === 0" class="py-2 text-gray-300">Tidak ada aktivitas terbaru.</li>
      </ul>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useWebSocketStore } from '../../stores/websocket'; // Import WebSocket store
import { useToast } from 'primevue/usetoast';
import Chart from 'primevue/chart';

export default {
  name: 'SuperAdminDashboardOverview',
  components: {
    Chart,
  },
  setup() {
    const summary = ref({
      total_companies: 0,
      active_subscriptions: 0,
      expired_subscriptions: 0,
      trial_subscriptions: 0,
    });
    const recentActivities = ref([]);
    const revenueData = ref([]);
    const toast = useToast();
    const webSocketStore = useWebSocketStore(); // Initialize WebSocket store

    const chartData = computed(() => {
      const allMonths = [
        'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
        'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
      ];

      const currentYear = new Date().getFullYear();
      const labels = [];
      const data = [];

      // Map revenueData to a more accessible format
      const revenueMap = new Map();
      revenueData.value.forEach(item => {
        // Ensure we only consider data for the current year
        if (parseInt(item.year) === currentYear) { // Ensure comparison is numeric
          revenueMap.set(item.month, item.total_revenue);
        }
      });

      // Populate labels and data for all 12 months of the current year
      for (let i = 0; i < 12; i++) {
        const monthName = allMonths[i];
        // Format month to YYYY-MM for matching with backend data
        const monthNumber = (i + 1).toString().padStart(2, '0');
        const monthYearKey = `${currentYear}-${monthNumber}`;

        labels.push(`${monthName} ${currentYear}`);
        data.push(revenueMap.get(monthYearKey) || 0);
      }

      return {
        labels: labels,
        datasets: [
          {
            label: 'Pendapatan',
            backgroundColor: '#42A5F5',
            borderColor: '#42A5F5',
            data: data,
            tension: 0.4,
          },
        ],
      };
    });

    const chartOptions = ref({
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          labels: {
            color: '#a0aec0', // text-text-muted
          },
        },
        tooltip: {
          callbacks: {
            label: function(context) {
              let label = context.dataset.label || '';
              if (label) {
                label += ': ';
              }
              if (context.parsed.y !== null) {
                label += new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(context.parsed.y);
              }
              return label;
            }
          }
        }
      },
      scales: {
        x: {
          ticks: {
            color: '#a0aec0',
          },
          grid: {
            color: 'rgba(160, 174, 192, 0.1)',
          },
        },
        y: {
          min: 0, // Ensure Y-axis starts from 0
          ticks: {
            color: '#a0aec0',
            callback: function(value) {
              return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', currencyDisplay: 'symbol' }).format(value);
            }
          },
          grid: {
            color: 'rgba(160, 174, 192, 0.1)',
          },
        },
      },
    });

    // Handler for WebSocket messages
    const handleWebSocketMessage = (data) => {
      console.log("WebSocket message received in handler:", data); // Added log
      if (data) {
        summary.value = data;
        recentActivities.value = data.recent_activities || [];
        revenueData.value = data.monthly_revenue || []; // Update revenue data from WebSocket
        console.log("Summary after update:", summary.value); // Added log
        console.log("Recent Activities after update:", recentActivities.value); // Added log
        console.log("Revenue Data after update:", revenueData.value); // Added log
      }
    };

    onMounted(() => {
      // If data already exists in the store (e.g., from a previous visit), use it
      if (webSocketStore.superAdminDashboardData) {
        handleWebSocketMessage(webSocketStore.superAdminDashboardData);
      }
      // Register WebSocket message handler
      webSocketStore.onMessage('superadmin_dashboard_update', handleWebSocketMessage);
    });

    onUnmounted(() => {
      // Unregister WebSocket message handler
      webSocketStore.offMessage('superadmin_dashboard_update');
    });

    return {
      summary,
      recentActivities,
      chartData,
      chartOptions,
    };
  },
};

const formatCurrency = (value) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(value);
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>
