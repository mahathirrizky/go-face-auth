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
      :editable="true"
      editKey="id"
      @save="onRowEditSave"
    >
      <template #header-actions>
        <BaseButton @click="openAddModal" class="btn-primary"><i class="pi pi-plus"></i> Tambah Paket Baru</BaseButton>
      </template>

      <template #column-price_monthly="{ item }">
        {{ formatCurrency(item.price_monthly) }}
      </template>

      <template #editor-price_monthly="{ data, field }">
        <InputNumber v-model="data[field]" mode="currency" currency="IDR" locale="id-ID" class="w-full" />
      </template>

      <template #column-price_yearly="{ item }">
        {{ formatCurrency(item.price_yearly) }}
      </template>

      <template #editor-price_yearly="{ data, field }">
        <InputNumber v-model="data[field]" mode="currency" currency="IDR" locale="id-ID" class="w-full" />
      </template>

      <template #editor-max_employees="{ data, field }">
        <InputNumber v-model="data[field]" class="w-full" />
      </template>

      <template #column-is_active="{ item }">
         <ToggleSwitch v-model="item.is_active" readonly />
      </template>

      <template #editor-is_active="{ data, field }">
        <ToggleSwitch v-model="data[field]" />
      </template>

      <template #actions="{ item, editor }">
        <div class="flex items-center justify-center space-x-2">
          <button v-if="!editor" v-row-editor-init class="p-link p-0 m-0">
            <BaseButton class="btn-sm btn-accent">
              <i class="pi pi-pencil"></i>
            </BaseButton>
          </button>
          <button v-else v-row-editor-save class="p-link p-0 m-0">
            <BaseButton class="btn-sm btn-success">
              <i class="pi pi-check"></i>
            </BaseButton>
          </button>
          <button v-if="!editor" @click="deletePackage(item.id)" class="p-link p-0 m-0">
            <BaseButton class="btn-sm btn-danger">
              <i class="pi pi-trash"></i>
            </BaseButton>
          </button>
          <button v-else v-row-editor-cancel class="p-link p-0 m-0">
            <BaseButton class="btn-sm btn-danger">
              <i class="pi pi-times"></i>
            </BaseButton>
          </button>
        </div>
      </template>
    </BaseDataTable>

    <!-- Add/Edit Package Modal -->
    <BaseModal :isOpen="isModalOpen" @close="closeModal" :title="isEditMode ? 'Edit Paket Langganan' : 'Tambah Paket Langganan Baru'">
      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div>
          <label for="packageName" class="block text-sm font-medium text-text-base mb-1">Nama Paket:</label>
          <InputText
            id="packageName"
            v-model="currentPackage.package_name"
            required
            class="w-full"
          />
        </div>
        <div>
          <label for="priceMonthly" class="block text-sm font-medium text-text-base mb-1">Harga Bulanan:</label>
          <InputNumber
            id="priceMonthly"
            v-model="currentPackage.price_monthly"
            mode="currency"
            currency="IDR"
            locale="id-ID"
            required
            class="w-full"
          />
        </div>
        <div>
          <label for="priceYearly" class="block text-sm font-medium text-text-base mb-1">Harga Tahunan:</label>
          <InputNumber
            id="priceYearly"
            v-model="currentPackage.price_yearly"
            mode="currency"
            currency="IDR"
            locale="id-ID"
            required
            class="w-full"
          />
        </div>
        <div>
          <label for="maxEmployees" class="block text-sm font-medium text-text-base mb-1">Max Karyawan:</label>
          <InputNumber
            id="maxEmployees"
            v-model="currentPackage.max_employees"
            required
            class="w-full"
          />
        </div>
        <div>
          <label for="features" class="block text-sm font-medium text-text-base mb-1">Fitur (pisahkan dengan koma):</label>
          <InputText
            id="features"
            v-model="currentPackage.features"
            class="w-full"
          />
        </div>
        <div class="flex items-center">
          <ToggleSwitch id="isActive" v-model="currentPackage.is_active" />
          <label for="isActive" class="text-text-base text-sm font-bold ml-2">Aktif</label>
        </div>
        <div class="flex justify-end pt-4">
          <BaseButton type="button" @click="closeModal" class="btn-secondary mr-2"><i class="pi pi-times"></i> Batal</BaseButton>
          <BaseButton type="submit" class="btn-primary"><i class="pi pi-save"></i> {{ isEditMode ? 'Update' : 'Tambah' }}</BaseButton>
        </div>
      </form>
    </BaseModal>

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
import BaseButton from '../ui/BaseButton.vue';
import BaseDataTable from '../ui/BaseDataTable.vue';

const toast = useToast();
const confirm = useConfirm();
const packages = ref([]);
const isLoading = ref(false);

const filters = ref({
    global: { value: null, matchMode: FilterMatchMode.CONTAINS }
});

const packageColumns = ref([
    { field: 'package_name', header: 'Nama Paket' },
    { field: 'price_monthly', header: 'Harga Bulanan' },
    { field: 'price_yearly', header: 'Harga Tahunan' },
    { field: 'max_employees', header: 'Max Karyawan' },
    { field: 'features', header: 'Fitur' },
    { field: 'is_active', header: 'Aktif' }
]);

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
  let { newData } = event;
  try {
    await axios.put(`/api/superadmin/subscription-packages/${newData.id}`, newData);
    toast.add({ severity: 'success', summary: 'Success', detail: 'Paket berhasil diperbarui!', life: 3000 });
    // Optionally re-fetch data to ensure consistency
    fetchPackages();
  } catch (error) {
    console.error('Error saving package:', error);
    toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data?.message || 'Gagal menyimpan paket.', life: 3000 });
    // You might want to revert the changes in the UI on failure
    // This can be done by finding the original data and replacing it back
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