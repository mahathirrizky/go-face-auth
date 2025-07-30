<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <Toast />
    <ConfirmDialog />
    <h2 class="text-2xl font-bold text-text-base mb-6">Manajemen Shift</h2>

    <DataTable
      :value="shifts"
      :loading="isLoading"
      :globalFilterFields="['name']"
      paginator
      :rows="10"
      class="p-datatable-customers"
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink CurrentPageReport RowsPerPageDropdown"
      :rowsPerPageOptions="[10, 25, 50]"
      currentPageReportTemplate="Menampilkan {first} sampai {last} dari {totalRecords} data"
      dataKey="id"
      v-model:filters="filters"
    >
      <template #header>
        <div class="flex flex-wrap justify-between items-center gap-4">
          <IconField iconPosition="left">
            <InputIcon class="pi pi-search"></InputIcon>
            <InputText v-model="filters['global'].value" placeholder="Cari Shift..." fluid />
          </IconField>
          <Button @click="openAddModal" icon="pi pi-plus" label="Tambah Shift" />
        </div>
      </template>
      <template #empty>
        Tidak ada data ditemukan.
      </template>
      <template #loading>
        Memuat data...
      </template>

      <Column field="name" header="Nama Shift" :sortable="true"></Column>
      <Column field="start_time" header="Waktu Mulai" :sortable="true"></Column>
      <Column field="end_time" header="Waktu Selesai" :sortable="true"></Column>
      <Column header="Aksi" style="min-width: 12rem">
        <template #body="{ data }">
          <Button @click="openEditModal(data)" icon="pi pi-pencil" class="p-button-rounded p-button-success mr-2" />
          <Button @click="deleteShift(data.id)" icon="pi pi-trash" class="p-button-rounded p-button-danger" />
        </template>
      </Column>
    </DataTable>

    <Dialog v-model:visible="isModalOpen" :header="editingShift ? 'Edit Shift' : 'Tambah Shift'" :modal="true" class="w-full max-w-md">
      <form @submit.prevent="saveShift" class="p-fluid">
        <div class="field mb-4">
          <label for="name">Nama Shift</label>
          <InputText id="name" v-model="currentShift.name" required :class="{ 'p-invalid': !currentShift.name }" fluid/>
        </div>
        <div class="field mb-4">
          <label for="start_time">Waktu Mulai</label>
          <InputMask id="start_time" v-model="currentShift.start_time" mask="99:99" required fluid />
        </div>
        <div class="field mb-4">
          <label for="end_time">Waktu Selesai</label>
          <InputMask id="end_time" v-model="currentShift.end_time" mask="99:99" required fluid />
        </div>
        <div class="field mb-4">
          <label for="grace_period_minutes">Toleransi Keterlambatan (menit)</label>
          <InputNumber id="grace_period_minutes" v-model="currentShift.grace_period_minutes" required fluid />
        </div>
        <div class="flex justify-end space-x-2 mt-6">
          <Button type="button" @click="closeModal" label="Batal" class="p-button-text"/>
          <Button type="submit" :loading="isSaving" :label="isSaving ? 'Menyimpan...' : 'Simpan'" />
        </div>
      </form>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { useConfirm } from 'primevue/useconfirm';
import { FilterMatchMode } from '@primevue/core/api';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import InputNumber from 'primevue/inputnumber';
import InputMask from 'primevue/inputmask';
import IconField from 'primevue/iconfield';
import InputIcon from 'primevue/inputicon';
import Toast from 'primevue/toast';
import ConfirmDialog from 'primevue/confirmdialog';

const shifts = ref([]);
const isModalOpen = ref(false);
const currentShift = ref({});
const editingShift = ref(false);
const toast = useToast();
const isLoading = ref(false);
const isSaving = ref(false);
const confirm = useConfirm();

const filters = ref({
    global: { value: null, matchMode: FilterMatchMode.CONTAINS }
});

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
  currentShift.value = { name: '', start_time: '', end_time: '', grace_period_minutes: 0 };
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
  isSaving.value = true;
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
    const errorMessage = error.response?.data?.error || 'Gagal menyimpan shift.';
    toast.add({ severity: 'error', summary: 'Error', detail: errorMessage, life: 3000 });
  } finally {
    isSaving.value = false;
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

<style scoped>
.field > label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 600;
}

.p-button-rounded {
    width: 2.5rem;
    height: 2.5rem;
}
</style>