<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-6">Manajemen Shift</h2>

    <div class="bg-bg-muted p-4 rounded-lg shadow-md mb-6 flex justify-end">
      <BaseButton @click="openAddModal">
        Tambah Shift
      </BaseButton>
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
    <BaseModal :isOpen="isModalOpen" @close="closeModal" :title="editingShift ? 'Edit Shift' : 'Tambah Shift'">
      <form @submit.prevent="saveShift">
        <BaseInput
          id="name"
          label="Nama Shift:"
          v-model="currentShift.name"
          required
        />
        <BaseInput
          id="start_time"
          label="Waktu Mulai:"
          v-model="currentShift.start_time"
          type="time"
          required
        />
        <BaseInput
          id="end_time"
          label="Waktu Selesai:"
          v-model="currentShift.end_time"
          type="time"
          required
        />
        <BaseInput
          id="grace_period_minutes"
          label="Toleransi Keterlambatan (menit):"
          v-model="currentShift.grace_period_minutes"
          type="number"
          required
        />
        <div class="flex justify-end space-x-4 mt-6">
          <BaseButton type="button" @click="closeModal" class="btn-outline-primary">Batal</BaseButton>
          <BaseButton type="submit">Simpan</BaseButton>
        </div>
      </form>
    </BaseModal>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { useToast } from 'vue-toastification';
import BaseModal from '../ui/BaseModal.vue';
import BaseInput from '../ui/BaseInput.vue';
import BaseButton from '../ui/BaseButton.vue';

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
</script>
