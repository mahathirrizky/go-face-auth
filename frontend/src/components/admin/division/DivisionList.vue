<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <Toast />
    <ConfirmDialog />
    <h2 class="text-2xl font-bold text-text-base mb-6">Manajemen Divisi</h2>

    <DataTable
      :value="divisions"
      :loading="isLoading"
      v-model:filters="filters"
      :globalFilterFields="['name']"
      paginator :rows="10" :rowsPerPageOptions="[10, 25, 50]"
      editMode="cell"
      dataKey="id"
      @cell-edit-complete="onCellEditComplete"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <Button @click="openAddModal" icon="pi pi-plus" label="Tambah Divisi" />
          <IconField iconPosition="left">
            <InputIcon class="pi pi-search"></InputIcon>
            <InputText v-model="filters['global'].value" placeholder="Cari Divisi..." />
          </IconField>
        </div>
      </template>

      <Column field="name" header="Nama Divisi" :sortable="true" style="min-width: 12rem">
        <template #editor="{ data, field }">
          <InputText v-model="data[field]" autofocus />
        </template>
      </Column>
      <Column field="description" header="Deskripsi" :sortable="true" style="min-width: 15rem">
        <template #editor="{ data, field }">
          <Textarea v-model="data[field]" rows="3" class="w-full" />
        </template>
      </Column>
      <Column field="shifts" header="Shift" style="min-width: 15rem">
        <template #body="{ data }">{{ data.Shifts ? data.Shifts.map(s => s.name).join(', ') : '-' }}</template>
        <template #editor="{ data, field }">
          <MultiSelect v-model="data.shift_ids" :options="shifts" optionLabel="name" optionValue="id" placeholder="Pilih Shift" display="chip" class="w-full" />
        </template>
      </Column>
      <Column field="locations" header="Lokasi Absen" style="min-width: 15rem">
        <template #body="{ data }">{{ data.Locations ? data.Locations.map(l => l.Name).join(', ') : '-' }}</template>
        <template #editor="{ data, field }">
          <MultiSelect v-model="data.location_ids" :options="locations" optionLabel="name" optionValue="id" placeholder="Pilih Lokasi" display="chip" class="w-full" />
        </template>
      </Column>
      <Column header="Aksi" style="width: 10rem; text-align: center">
        <template #body="{ data }">
          <Button @click="openEditModal(data)" icon="pi pi-pencil" class="p-button-rounded p-button-success mr-2" />
          <Button @click="deleteDivision(data.id)" icon="pi pi-trash" class="p-button-rounded p-button-danger" />
        </template>
      </Column>
    </DataTable>

    <Dialog v-model:visible="isModalOpen" :header="modalTitle" :modal="true" class="w-full max-w-lg">
      <form @submit.prevent="saveDivision" class="space-y-3 ">
      
          <FloatLabel variant="on" class="mt-3">
            <InputText id="name" v-model="currentDivision.name" required fluid/>
            <label for="name">Nama Divisi</label>
          </FloatLabel>
       
 
          <FloatLabel variant="on">
            <Textarea id="description" v-model="currentDivision.description" rows="3" fluid />
            <label for="description">Deskripsi</label>
          </FloatLabel>
       
     
          <FloatLabel variant="on">
            <MultiSelect id="shifts" v-model="currentDivision.shift_ids" :options="shifts" optionLabel="name" optionValue="id" display="chip" :loading="isFetchingShifts" class="w-full" />
            <label for="shifts">Pilih Shift</label>
          </FloatLabel>
     
        <div class="field mb-4">
            <label for="locations" class="block mb-2">Pilih Lokasi Absen</label>
            <div v-if="isFetchingLocations" class="flex items-center">
                <ProgressSpinner style="width: 2rem; height: 2rem" class="mr-2" />
                <span>Memuat lokasi...</span>
            </div>
            <div v-else-if="locations.length === 0">
                <Message severity="warn" :closable="false">
                    Tidak ada lokasi tersedia. 
                    <router-link :to="{ name: 'LocationManagement' }" class="underline font-bold text-yellow-700 hover:text-yellow-600">Tambah Lokasi Baru</router-link>.
                </Message>
            </div>
            <FloatLabel v-else variant="on">
                <MultiSelect id="locations" v-model="currentDivision.location_ids" :options="locations" optionLabel="name" optionValue="id" placeholder="Pilih Lokasi" display="chip" class="w-full" />
                <label for="locations">Pilih Lokasi</label>
            </FloatLabel>
        </div>
        <div class="flex justify-end space-x-2 mt-8">
          <Button @click="closeModal" type="button" label="Batal" class="p-button-text"/>
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
import { useAuthStore } from '../../../stores/auth';
import { FilterMatchMode } from '@primevue/core/api';
import { useRouter } from 'vue-router';

import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import Textarea from 'primevue/textarea';
import MultiSelect from 'primevue/multiselect';
import IconField from 'primevue/iconfield';
import InputIcon from 'primevue/inputicon';
import Message from 'primevue/message';
import ProgressSpinner from 'primevue/progressspinner';
import Toast from 'primevue/toast';
import ConfirmDialog from 'primevue/confirmdialog';
import FloatLabel from 'primevue/floatlabel';

const toast = useToast();
const confirm = useConfirm();
const authStore = useAuthStore();
const router = useRouter();

const divisions = ref([]);
const isLoading = ref(false);
const isModalOpen = ref(false);
const modalTitle = ref('');
const isEditMode = ref(false);
const currentDivision = ref({});
const isSaving = ref(false);

const filters = ref({ global: { value: null, matchMode: FilterMatchMode.CONTAINS } });

const shifts = ref([]);
const locations = ref([]);
const isFetchingShifts = ref(false);
const isFetchingLocations = ref(false);

const fetchDivisions = async () => {
  isLoading.value = true;
  try {
    const response = await axios.get('/api/admin/divisions');
    divisions.value = response.data.data.map(div => ({
      ...div,
      id: div.ID,
      name: div.Name,
      description: div.Description,
      shift_ids: div.Shifts ? div.Shifts.map(s => s.id) : [],
      location_ids: div.Locations ? div.Locations.map(l => l.ID) : []
    }));
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memuat divisi.', life: 3000 });
  } finally {
    isLoading.value = false;
  }
};

const fetchShifts = async () => {
  isFetchingShifts.value = true;
  try {
    const response = await axios.get('/api/shifts');
    shifts.value = response.data.data.map(s => ({ id: s.id, name: s.name }));
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memuat shift.', life: 3000 });
  } finally {
    isFetchingShifts.value = false;
  }
};

const fetchLocations = async () => {
  isFetchingLocations.value = true;
  try {
    const response = await axios.get('/api/company/locations');
    locations.value = (response.data.data || []).map(loc => ({ id: loc.ID, name: loc.Name }));
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memuat lokasi.', life: 3000 });
    locations.value = [];
  } finally {
    isFetchingLocations.value = false;
  }
};

const openAddModal = () => {
  isEditMode.value = false;
  modalTitle.value = 'Tambah Divisi Baru';
  currentDivision.value = { name: '', description: '', shift_ids: [], location_ids: [] };
  isModalOpen.value = true;
};

const openEditModal = (division) => {
  isEditMode.value = true;
  modalTitle.value = 'Edit Divisi';
  currentDivision.value = { ...division };
  isModalOpen.value = true;
};

const closeModal = () => {
  isModalOpen.value = false;
};

const saveDivision = async () => {
  if (!currentDivision.value.name.trim()) {
    toast.add({ severity: 'warn', summary: 'Peringatan', detail: 'Nama divisi tidak boleh kosong.', life: 3000 });
    return;
  }
  isSaving.value = true;
  try {
    const payload = {
      name: currentDivision.value.name,
      description: currentDivision.value.description,
      shift_ids: currentDivision.value.shift_ids,
      location_ids: currentDivision.value.location_ids,
    };
    const url = isEditMode.value ? `/api/admin/divisions/${currentDivision.value.id}` : '/api/admin/divisions';
    const method = isEditMode.value ? 'put' : 'post';
    await axios[method](url, payload);
    toast.add({ severity: 'success', summary: 'Berhasil', detail: `Divisi berhasil ${isEditMode.value ? 'diperbarui' : 'ditambahkan'}.`, life: 3000 });
    closeModal();
    fetchDivisions();
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.response?.data?.message || 'Gagal menyimpan divisi.', life: 3000 });
  } finally {
    isSaving.value = false;
  }
};

const onCellEditComplete = async (event) => {
  let { data, field, value } = event;
  let payload = { [field]: value };

  if (field === 'shifts') payload = { shift_ids: data.shift_ids };
  if (field === 'locations') payload = { location_ids: data.location_ids };

  try {
    await axios.put(`/api/admin/divisions/${data.id}`, payload);
    toast.add({ severity: 'success', summary: 'Berhasil', detail: 'Divisi berhasil diperbarui.', life: 3000 });
    fetchDivisions();
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memperbarui divisi.', life: 3000 });
    fetchDivisions();
  }
};

const deleteDivision = (id) => {
  confirm.require({
    message: 'Apakah Anda yakin ingin menghapus divisi ini?',
    header: 'Konfirmasi Hapus',
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      try {
        await axios.delete(`/api/admin/divisions/${id}`);
        toast.add({ severity: 'success', summary: 'Berhasil', detail: 'Divisi berhasil dihapus.', life: 3000 });
        fetchDivisions();
      } catch (err) {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal menghapus divisi.', life: 3000 });
      }
    }
  });
};

onMounted(() => {
  fetchDivisions();
  fetchShifts();
  fetchLocations();
});
</script>
