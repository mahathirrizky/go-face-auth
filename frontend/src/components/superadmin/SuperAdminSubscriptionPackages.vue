<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-4">Manajemen Paket Langganan</h2>

    <BaseDataTable
      :data="packages"
      :columns="packageColumns"
      :loading="isLoading"
      v-model:filters="filters"
      :globalFilterFields="['package_name', 'features']"
      searchPlaceholder="Cari Paket..."
      editKey="id"
      @cell-edit-complete="onRowEditSave"
    >
      <template #header-actions>
        <!-- Add new package functionality can be re-implemented here if needed -->
      </template>

      <template #column-price_monthly="{ item }">
        {{ formatCurrency(item.price_monthly) }}
      </template>

      <template #editor-price_monthly="{ data, field }">
        <InputNumber v-model="data[field]" mode="currency" currency="IDR" locale="id-ID" class="w-full" autofocus />
      </template>

      <template #editor-package_name="{ data, field }">
        <BaseInput v-model="data[field]" :id="`edit-${field}`" :name="field" autofocus />
      </template>

      <template #column-price_yearly="{ item }">
        {{ formatCurrency(item.price_yearly) }}
      </template>

      <template #editor-price_yearly="{ data, field }">
        <InputNumber v-model="data[field]" mode="currency" currency="IDR" locale="id-ID" class="w-full" autofocus />
      </template>

      <template #editor-max_employees="{ data, field }">
        <InputNumber v-model="data[field]" class="w-full" autofocus />
      </template>

      <template #column-is_active="{ item }">
         <ToggleSwitch v-model="item.is_active" readonly />
      </template>

      <template #editor-is_active="{ data, field }">
        <ToggleSwitch v-model="data[field]" />
      </template>

      <template #editor-features="{ data, field }">
        <Textarea v-model="data[field]" :id="`edit-${field}`" :name="field" rows="3" class="w-full" autofocus />
      </template>

      <template #actions="{ item }">
        <BaseButton @click="deletePackage(item.id)" class="btn-sm btn-danger">
          <i class="pi pi-trash"></i>
        </BaseButton>
      </template>
    </BaseDataTable>

    <ConfirmDialog />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { useConfirm } from 'primevue/useconfirm';
import { FilterMatchMode } from '@primevue/core/api';
import ConfirmDialog from 'primevue/confirmdialog';
import InputText from 'primevue/inputtext';
import InputNumber from 'primevue/inputnumber';
import ToggleSwitch from 'primevue/toggleswitch';
import Textarea from 'primevue/textarea';
import BaseButton from '../ui/BaseButton.vue';
import BaseDataTable from '../ui/BaseDataTable.vue';
import BaseInput from '../ui/BaseInput.vue';

const toast = useToast();
const confirm = useConfirm();
const packages = ref([]);
const isLoading = ref(false);

const packageColumns = ref([
    { field: 'package_name', header: 'Nama Paket', editable: true },
    { field: 'price_monthly', header: 'Harga Bulanan', editable: true },
    { field: 'price_yearly', header: 'Harga Tahunan', editable: true },
    { field: 'max_employees', header: 'Max Karyawan', editable: true },
    { field: 'features', header: 'Fitur', editable: true },
    { field: 'is_active', header: 'Aktif', editable: true }
]);

const filters = ref({
    global: { value: null, matchMode: FilterMatchMode.CONTAINS }
});

const fetchPackages = async () => {
  isLoading.value = true;
  try {
    const response = await axios.get('/api/superadmin/subscription-packages');
    if (response.data.status === 'success') {
      packages.value = response.data.data;
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: response.data.message || 'Gagal mengambil daftar paket.', life: 3000 });
    }
  } catch (error) {
    console.error('Error fetching packages:', error);
    toast.add({ severity: 'error', summary: 'Error', detail: 'Terjadi kesalahan saat mengambil paket.', life: 3000 });
  } finally {
    isLoading.value = false;
  }
};

const onRowEditSave = async (event) => {
  let { data, field, newValue } = event;

  // Create a payload with only the changed field
  const updatePayload = {
    [field]: newValue,
  };

  try {
    await axios.put(`/api/superadmin/subscription-packages/${data.id}`, updatePayload);
    toast.add({ severity: 'success', summary: 'Success', detail: 'Paket berhasil diperbarui!', life: 3000 });
    // Perbarui data di array packages.value secara langsung
    const index = packages.value.findIndex(p => p.id === data.id);
    if (index !== -1) {
      packages.value[index] = { ...packages.value[index], ...updatePayload };
    }
  } catch (error) {
    console.error('Error saving package:', error);
    toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data?.message || 'Gagal menyimpan paket.', life: 3000 });
    // Revert the local change if the backend update fails
    // This requires storing the original value before editing, which is more complex
    // For now, we'll just show an error and rely on a re-fetch if needed
    fetchPackages(); // Re-fetch to ensure data consistency
  }
};

const deletePackage = (id) => {
  confirm.require({
    message: 'Apakah Anda yakin ingin menghapus paket ini?',
    header: 'Konfirmasi Hapus',
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      try {
        await axios.delete(`/api/superadmin/subscription-packages/${id}`);
        toast.add({ severity: 'success', summary: 'Berhasil', detail: 'Paket berhasil dihapus', life: 3000 });
        fetchPackages();
      } catch (error) {
        toast.add({ severity: 'error', summary: 'Gagal', detail: error.response?.data?.message || 'Gagal menghapus paket', life: 3000 });
      }
    },
    reject: () => {
      toast.add({ severity: 'info', summary: 'Dibatalkan', detail: 'Penghapusan paket dibatalkan', life: 3000 });
    }
  });
};

const formatCurrency = (value) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(value);
};

onMounted(() => {
  fetchPackages();
});
</script>