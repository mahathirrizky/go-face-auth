<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-4">Manajemen Paket Langganan</h2>

    <DataTable
      :value="packages"
      :loading="isLoading"
      v-model:filters="filters"
      editMode="cell"
      dataKey="id"
      @cell-edit-complete="onRowEditSave"
      paginator
      :rows="10"
      class="p-datatable-customers"
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink CurrentPageReport RowsPerPageDropdown"
      :rowsPerPageOptions="[10, 25, 50]"
      currentPageReportTemplate="Menampilkan {first} sampai {last} dari {totalRecords} data"
      :globalFilterFields="['package_name', 'features']"
    >
      <template #header>
        <div class="flex flex-wrap justify-between items-center gap-4">
          <IconField iconPosition="left">
            <InputIcon class="pi pi-search"></InputIcon>
            <InputText v-model="filters['global'].value" placeholder="Cari Paket..." fluid />
          </IconField>
        </div>
      </template>
      <template #empty>
        Tidak ada data ditemukan.
      </template>
      <template #loading>
        Memuat data...
      </template>

      <Column field="package_name" header="Nama Paket" :sortable="true" style="min-width: 12rem">
        <template #editor="{ data, field }">
          <InputText v-model="data[field]" autofocus class="w-full" fluid />
        </template>
      </Column>

      <Column field="price_monthly" header="Harga Bulanan" :sortable="true" style="min-width: 12rem">
        <template #body="{ data, field }">
          {{ formatCurrency(data[field]) }}
        </template>
        <template #editor="{ data, field }">
          <InputNumber v-model="data[field]" mode="currency" currency="IDR" locale="id-ID" class="w-full" autofocus fluid />
        </template>
      </Column>

      <Column field="price_yearly" header="Harga Tahunan" :sortable="true" style="min-width: 12rem">
        <template #body="{ data, field }">
          {{ formatCurrency(data[field]) }}
        </template>
        <template #editor="{ data, field }">
          <InputNumber v-model="data[field]" mode="currency" currency="IDR" locale="id-ID" class="w-full" autofocus fluid />
        </template>
      </Column>

      <Column field="max_employees" header="Max Karyawan" :sortable="true" style="min-width: 10rem">
        <template #editor="{ data, field }">
          <InputNumber v-model="data[field]" class="w-full" autofocus fluid />
        </template>
      </Column>

      <Column field="features" header="Fitur" :sortable="true" style="min-width: 15rem">
         <template #editor="{ data, field }">
            <Textarea v-model="data[field]" rows="3" class="w-full" autofocus fluid />
        </template>
      </Column>

      <Column field="is_active" header="Aktif" :sortable="true" style="min-width: 8rem">
        <template #body="{ data, field }">
           <ToggleSwitch v-model="data[field]" readonly />
        </template>
        <template #editor="{ data, field }">
          <ToggleSwitch v-model="data[field]" />
        </template>
      </Column>

      <Column header="Aksi" style="width: 10%; min-width: 8rem" bodyStyle="text-align:center">
        <template #body="{ data }">
          <Button @click="deletePackage(data.id)" class="p-button-sm p-button-danger" icon="pi pi-trash" />
        </template>
      </Column>
    </DataTable>

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
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import IconField from 'primevue/iconfield';
import InputIcon from 'primevue/inputicon';


const toast = useToast();
const confirm = useConfirm();
const packages = ref([]);
const isLoading = ref(false);

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
    // Update data in packages.value directly
    const index = packages.value.findIndex(p => p.id === data.id);
    if (index !== -1) {
      packages.value[index] = { ...packages.value[index], ...updatePayload };
    }
  } catch (error) {
    console.error('Error saving package:', error);
    toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data?.message || 'Gagal menyimpan paket.', life: 3000 });
    // Revert the local change if the backend update fails
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
