<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-4">Grafik Pendapatan</h2>

    <div class="bg-bg-muted p-6 rounded-lg shadow-md mb-6">
      <h3 class="text-xl font-semibold text-text-base mb-4">Pilih Rentang Tanggal</h3>
      <div class="flex flex-col md:flex-row items-center space-y-4 md:space-y-0 md:space-x-4">
        <div class="flex flex-col">
          <label for="startDate" class="text-text-muted text-sm mb-1">Tanggal Mulai:</label>
          <DatePicker
            v-model="startDate"
            dateFormat="yy-mm-dd"
            :maxDate="new Date()"
            placeholder="Pilih Tanggal Mulai"
            class="w-full md:w-auto"
          />
        </div>
        <div class="flex flex-col">
          <label for="endDate" class="text-text-muted text-sm mb-1">Tanggal Akhir:</label>
          <DatePicker
            v-model="endDate"
            dateFormat="yy-mm-dd"
            :maxDate="new Date()"
            :minDate="startDate"
            placeholder="Pilih Tanggal Akhir"
            class="w-full md:w-auto"
          />
        </div>
        
      </div>
    </div>

    <div class="bg-bg-muted p-6 rounded-lg shadow-md">
      <h3 class="text-xl font-semibold text-text-base mb-4">Pendapatan Bulanan</h3>
      <Chart type="line" :data="chartData" :options="chartOptions" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import Chart from 'primevue/chart';
import DatePicker from 'primevue/datepicker';

const toast = useToast();
const startDate = ref(null);
const endDate = ref(null);
const revenueData = ref([]);

const fetchRevenueData = async () => {
  try {
    const response = await axios.get('/api/superadmin/revenue-summary', {
      params: {
        start_date: startDate.value ? startDate.value.toISOString().slice(0, 10) : null,
        end_date: endDate.value ? endDate.value.toISOString().slice(0, 10) : null,
      },
    });
    if (response.data && response.data.status === 'success') {
      revenueData.value = response.data.data;
      console.log("Fetched revenueData:", revenueData.value); // Added log
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data.message || 'Failed to fetch revenue data.', life: 3000 });
    }
  } catch (error) {
    console.error('Error fetching revenue data:', error);
    toast.add({ severity: 'error', summary: 'Error', detail: 'An error occurred while fetching revenue data.', life: 3000 });
  }
};

const chartData = computed(() => {
  const allMonths = [
    'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
    'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
  ];

  const labels = [];
  const data = [];

  const revenueMap = new Map();
  if (revenueData.value) {
    revenueData.value.forEach(item => {
      revenueMap.set(item.month, item.total_revenue);
    });
  }
  console.log("Revenue Map:", revenueMap); // Added log

  let start = startDate.value ? new Date(startDate.value) : new Date();
  let end = endDate.value ? new Date(endDate.value) : new Date();

  let currentDate = new Date(start.getFullYear(), start.getMonth(), 1);
  while (currentDate <= end) {
    const monthName = allMonths[currentDate.getMonth()];
    const year = currentDate.getFullYear();
    const monthNumber = (currentDate.getMonth() + 1).toString().padStart(2, '0');
    const monthKey = `${year}-${monthNumber}`;

    console.log(`Processing month: ${monthKey}, Value from map: ${revenueMap.get(monthKey)}`); // Added log

    labels.push(`${monthName} ${year}`);
    data.push(revenueMap.get(monthKey) || 0);

    currentDate.setMonth(currentDate.getMonth() + 1);
  }

  console.log("Final Chart Labels:", labels); // Added log
  console.log("Final Chart Data:", data); // Added log

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
        color: '#a0aec0',
      },
    },
    tooltip: {
      callbacks: {
        label: (context) => {
          let label = context.dataset.label || '';
          if (label) {
            label += ': ';
          }
          if (context.parsed.y !== null) {
            label += new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', currencyDisplay: 'symbol' }).format(context.parsed.y);
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
      min: 0,
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

onMounted(() => {
  const today = new Date();
  const firstDayOfYear = new Date(today.getFullYear(), 0, 1);
  startDate.value = firstDayOfYear;
  endDate.value = today;
  fetchRevenueData();
});

watch([startDate, endDate], ([newStartDate, newEndDate], [oldStartDate, oldEndDate]) => {
  if ((newStartDate !== oldStartDate || newEndDate !== oldEndDate) && newStartDate && newEndDate) {
    fetchRevenueData();
  }
});
</script>