<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-6">Manajemen Absensi</h2>

    <div class="bg-bg-muted p-4 rounded-lg shadow-md mb-6 flex flex-col md:flex-row justify-between items-center">
      <div class="flex flex-col md:flex-row gap-4 w-full md:w-2/3">
        <input
          type="date"
          class="p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary"
        />
        <select
          class="p-2 rounded-md border border-bg-base bg-bg-base text-text-base focus:outline-none focus:ring-2 focus:ring-secondary"
        >
          <option value="">Pilih Karyawan</option>
          <option value="1">Budi Santoso</option>
          <option value="2">Siti Aminah</option>
        </select>
      </div>
      <button class="btn btn-secondary w-full md:w-auto mt-4 md:mt-0">
        Filter Absensi
      </button>
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
            <td class="px-6 py-4 whitespace-nowrap text-text-base">{{ record.date }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ record.employeeName }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ record.checkIn }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ record.checkOut }}</td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span :class="{
                'px-2 inline-flex text-xs leading-5 font-semibold rounded-full': true,
                'bg-green-100 text-green-800': record.status === 'Hadir',
                'bg-red-100 text-red-800': record.status === 'Terlambat',
                'bg-yellow-100 text-yellow-800': record.status === 'Izin'
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
import { ref } from 'vue';

export default {
  name: 'AttendanceManagement',
  setup() {
    const attendanceRecords = ref([
      { id: 1, date: '2024-07-01', employeeName: 'Budi Santoso', checkIn: '08:00', checkOut: '17:00', status: 'Hadir' },
      { id: 2, date: '2024-07-01', employeeName: 'Siti Aminah', checkIn: '08:30', checkOut: '17:00', status: 'Terlambat' },
      { id: 3, date: '2024-07-01', employeeName: 'Joko Susilo', checkIn: '-', checkOut: '-', status: 'Izin' },
    ]);

    return {
      attendanceRecords,
    };
  },
};
</script>

<style scoped>
/* Tailwind handles styling */
</style>
