<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-6">Manajemen Shift</h2>

    <BaseDataTable
      :data="shifts"
      :columns="shiftColumns"
      :loading="isLoading"
      :globalFilterFields="['name']"
      searchPlaceholder="Cari Shift..."
    >
      <template #header-actions>
        <BaseButton @click="openAddModal">
          <i class="pi pi-plus"></i> Tambah Shift
        </BaseButton>
      </template>

      <template #column-actions="{ item }">
        <BaseButton @click="openEditModal(item)" class="text-accent hover:opacity-80 mr-3"><i class="pi pi-pencil"></i> Edit</BaseButton>
        <BaseButton @click="deleteShift(item.id)" class="text-danger hover:opacity-80"><i class="pi pi-trash"></i> Hapus</BaseButton>
      </template>
    </BaseDataTable>

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
          <BaseButton type="button" @click="closeModal" class="btn-outline-primary"><i class="pi pi-times"></i> Batal</BaseButton>
          <BaseButton type="submit"><i class="pi pi-save"></i> Simpan</BaseButton>
        </div>
      </form>
    </BaseModal>

    
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { useConfirm } from 'primevue/useconfirm';

import BaseModal from '../ui/BaseModal.vue';
import BaseInput from '../ui/BaseInput.vue';
import BaseButton from '../ui/BaseButton.vue';
import BaseDataTable from '../ui/BaseDataTable.vue';

const shifts = ref([]);
const isModalOpen = ref(false);
const currentShift = ref({});
const editingShift = ref(false);
const toast = useToast();
const isLoading = ref(false);
const confirm = useConfirm();

const shiftColumns = ref([
    { field: 'name', header: 'Nama Shift' },
    { field: 'start_time', header: 'Waktu Mulai' },
    { field: 'end_time', header: 'Waktu Selesai' },
    { field: 'actions', header: 'Aksi' }
]);

const fetchShifts = async () => {
  isLoading.value = true;
  try {
    const response = await axios.get('/api/shifts');
    if (response.data && response.data.status === 'success') {
      shifts.value = response.data.data;
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data?.message || 'Gagal mengambil data shift.', life: 3000 });
    }
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal mengambil data shift.', life: 3000 });
  } finally {
    isLoading.value = false;
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
      toast.add({ severity: 'success', summary: 'Success', detail: 'Shift berhasil diperbarui.', life: 3000 });
    } else {
      await axios.post('/api/shifts', currentShift.value);
      toast.add({ severity: 'success', summary: 'Success', detail: 'Shift berhasil ditambahkan.', life: 3000 });
    }
    closeModal();
    fetchShifts();
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal menyimpan shift.', life: 3000 });
  }
};

const deleteShift = (id) => {
  confirm.require({
    message: 'Apakah Anda yakin ingin menghapus shift ini? Tindakan ini tidak dapat dibatalkan.',
    header: 'Konfirmasi Hapus Shift',
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      try {
        await axios.delete(`/api/shifts/${id}`);
        toast.add({ severity: 'success', summary: 'Success', detail: 'Shift berhasil dihapus.', life: 3000 });
        fetchShifts();
      } catch (error) {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal menghapus shift.', life: 3000 });
      }
    },
    reject: () => {
      toast.add({ severity: 'info', summary: 'Dibatalkan', detail: 'Penghapusan shift dibatalkan', life: 3000 });
    }
  });
};
</script>
