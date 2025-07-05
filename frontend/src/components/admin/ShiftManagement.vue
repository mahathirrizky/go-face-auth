<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-6">Manajemen Shift</h2>

    <div class="bg-bg-muted p-4 rounded-lg shadow-md mb-6 flex justify-end">
      <button @click="openAddModal" class="btn btn-secondary">
        Tambah Shift
      </button>
    </div>

    <div class="overflow-x-auto bg-bg-muted rounded-lg shadow-md">
      <table class="min-w-full divide-y divide-bg-base">
        <thead class="bg-primary">
          <tr>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Nama Shift</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Waktu Mulai</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Waktu Selesai</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-bg-base">
          <tr v-for="shift in shifts" :key="shift.id">
            <td class="px-6 py-4 whitespace-nowrap text-text-base">{{ shift.name }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ shift.start_time }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-text-muted">{{ shift.end_time }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
              <button @click="openEditModal(shift)" class="text-accent hover:text-secondary mr-3">Edit</button>
              <button @click="deleteShift(shift.id)" class="text-danger hover:opacity-80">Hapus</button>
            </td>
          </tr>
          <tr v-if="shifts.length === 0">
            <td colspan="4" class="px-6 py-4 text-center text-text-muted">Tidak ada data shift.</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Add/Edit Shift Modal -->
    <div v-if="isModalOpen" class="fixed inset-0 bg-black bg-opacity-20 flex items-center justify-center z-50">
      <div class="bg-bg-muted p-8 rounded-lg shadow-lg w-full max-w-md">
        <h3 class="text-2xl font-bold text-text-base mb-6">{{ editingShift ? 'Edit Shift' : 'Tambah Shift' }}</h3>
        <form @submit.prevent="saveShift">
          <div class="mb-4">
            <label for="name" class="block text-text-muted text-sm font-bold mb-2">Nama Shift:</label>
            <input type="text" id="name" v-model="currentShift.name" class="w-full p-2 rounded-md border border-bg-base bg-bg-base text-text-base" required />
          </div>
          <div class="mb-4">
            <label for="start_time" class="block text-text-muted text-sm font-bold mb-2">Waktu Mulai:</label>
            <input type="time" id="start_time" v-model="currentShift.start_time" class="w-full p-2 rounded-md border border-bg-base bg-bg-base text-text-base" required />
          </div>
          <div class="mb-4">
            <label for="end_time" class="block text-text-muted text-sm font-bold mb-2">Waktu Selesai:</label>
            <input type="time" id="end_time" v-model="currentShift.end_time" class="w-full p-2 rounded-md border border-bg-base bg-bg-base text-text-base" required />
          </div>
          <div class="mb-6">
            <label for="grace_period_minutes" class="block text-text-muted text-sm font-bold mb-2">Toleransi Keterlambatan (menit):</label>
            <input type="number" id="grace_period_minutes" v-model="currentShift.grace_period_minutes" class="w-full p-2 rounded-md border border-bg-base bg-bg-base text-text-base" required />
          </div>
          <div class="flex justify-end space-x-4">
            <button type="button" @click="closeModal" class="btn btn-outline-primary">Batal</button>
            <button type="submit" class="btn btn-secondary">Simpan</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { useToast } from 'vue-toastification';

export default {
  name: 'ShiftManagement',
  setup() {
    const shifts = ref([]);
    const isModalOpen = ref(false);
    const currentShift = ref({});
    const editingShift = ref(false);
    const toast = useToast();

    const fetchShifts = async () => {
      try {
        const response = await axios.get('/api/shifts');
        if (response.data && response.data.status === 'success') {
          shifts.value = response.data.data;
        } else {
          toast.error(response.data?.message || 'Gagal mengambil data shift.');
        }
      } catch (error) {
        toast.error('Gagal mengambil data shift.');
      }
    };

    onMounted(fetchShifts);

    const openAddModal = () => {
      currentShift.value = { name: '', start_time: '', end_time: '' };
      editingShift.value = false;
      isModalOpen.value = true;
    };

    const openEditModal = (shift) => {
      currentShift.value = { ...shift };
      editingShift.value = true;
      isModalOpen.value = true;
    };

    const closeModal = () => {
      isModalOpen.value = false;
    };

    const saveShift = async () => {
      try {
        if (editingShift.value) {
          await axios.put(`/api/shifts/${currentShift.value.id}`, currentShift.value);
          toast.success('Shift berhasil diperbarui.');
        } else {
          await axios.post('/api/shifts', currentShift.value);
          toast.success('Shift berhasil ditambahkan.');
        }
        closeModal();
        fetchShifts();
      } catch (error) {
        toast.error('Gagal menyimpan shift.');
      }
    };

    const deleteShift = async (id) => {
      if (confirm('Apakah Anda yakin ingin menghapus shift ini?')) {
        try {
          await axios.delete(`/api/shifts/${id}`);
          toast.success('Shift berhasil dihapus.');
          fetchShifts();
        } catch (error) {
          toast.error('Gagal menghapus shift.');
        }
      }
    };

    return {
      shifts,
      isModalOpen,
      currentShift,
      editingShift,
      openAddModal,
      openEditModal,
      closeModal,
      saveShift,
      deleteShift,
    };
  },
};
</script>
