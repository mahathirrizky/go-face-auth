<template>
  <div class="p-6 bg-bg-base min-h-screen">
    <h2 class="text-2xl font-bold text-text-base mb-4">Manajemen Paket Langganan</h2>

    <!-- Add New Package Button -->
    <div class="mb-6">
      <BaseButton @click="openAddModal" class="btn-primary">Tambah Paket Baru</BaseButton>
    </div>

    <!-- Packages List Table -->
    <div class="bg-bg-muted p-6 rounded-lg shadow-md">
      <h3 class="text-xl font-semibold text-text-base mb-4">Daftar Paket</h3>
      <div class="overflow-x-auto">
        <table class="min-w-full bg-bg-muted text-text-base">
          <thead>
            <tr>
              <th class="py-2 px-4 border-b border-bg-base text-left">Nama Paket</th>
              <th class="py-2 px-4 border-b border-bg-base text-left">Harga Bulanan</th>
              <th class="py-2 px-4 border-b border-bg-base text-left">Harga Tahunan</th>
              <th class="py-2 px-4 border-b border-bg-base text-left">Max Karyawan</th>
              <th class="py-2 px-4 border-b border-bg-base text-left">Fitur</th>
              <th class="py-2 px-4 border-b border-bg-base text-left">Aktif</th>
              <th class="py-2 px-4 border-b border-bg-base text-left">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="packages.length === 0">
              <td colspan="7" class="py-4 px-4 text-center text-text-muted">Tidak ada paket langganan.</td>
            </tr>
            <tr v-for="pkg in packages" :key="pkg.id" class="hover:bg-bg-hover">
              <td class="py-2 px-4 border-b border-bg-base">{{ pkg.package_name }}</td>
              <td class="py-2 px-4 border-b border-bg-base">{{ formatCurrency(pkg.price_monthly) }}</td>
              <td class="py-2 px-4 border-b border-bg-base">{{ formatCurrency(pkg.price_yearly) }}</td>
              <td class="py-2 px-4 border-b border-bg-base">{{ pkg.max_employees }}</td>
              <td class="py-2 px-4 border-b border-bg-base">{{ pkg.features }}</td>
              <td class="py-2 px-4 border-b border-bg-base">
                <span :class="pkg.is_active ? 'text-green-500' : 'text-red-500'">
                  {{ pkg.is_active ? 'Ya' : 'Tidak' }}
                </span>
              </td>
              <td class="py-2 px-4 border-b border-bg-base">
                <BaseButton @click="openEditModal(pkg)" class="btn-sm btn-secondary mr-2">Edit</BaseButton>
                <BaseButton @click="deletePackage(pkg.id)" class="btn-sm btn-danger">Hapus</BaseButton>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Add/Edit Package Modal -->
    <BaseModal :isOpen="isModalOpen" @close="closeModal" :title="isEditMode ? 'Edit Paket Langganan' : 'Tambah Paket Langganan Baru'">
      <form @submit.prevent="handleSubmit">
        <BaseInput
          id="packageName"
          label="Nama Paket:"
          v-model="currentPackage.package_name"
          required
        />
        <BaseInput
          id="priceMonthly"
          label="Harga Bulanan:"
          v-model="currentPackage.price_monthly"
          type="number"
          required
        />
        <BaseInput
          id="priceYearly"
          label="Harga Tahunan:"
          v-model="currentPackage.price_yearly"
          type="number"
          required
        />
        <BaseInput
          id="maxEmployees"
          label="Max Karyawan:"
          v-model="currentPackage.max_employees"
          type="number"
          required
        />
        <BaseInput
          id="features"
          label="Fitur (pisahkan dengan koma):"
          v-model="currentPackage.features"
        />
        <div class="mb-4 flex items-center">
          <input type="checkbox" id="isActive" v-model="currentPackage.is_active" class="form-checkbox mr-2">
          <label for="isActive" class="text-text-base text-sm font-bold">Aktif</label>
        </div>
        <div class="flex justify-end">
          <BaseButton type="button" @click="closeModal" class="btn-secondary mr-2">Batal</BaseButton>
          <BaseButton type="submit" class="btn-primary">{{ isEditMode ? 'Update' : 'Tambah' }}</BaseButton>
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

const toast = useToast();
const packages = ref([]);
const isModalOpen = ref(false);
const isEditMode = ref(false);
const currentPackage = ref({
  id: null,
  package_name: '',
  price_monthly: 0,
  price_yearly: 0,
  max_employees: 0,
  features: '',
  is_active: true,
});

const fetchPackages = async () => {
  try {
    const response = await axios.get('/api/superadmin/subscription-packages');
    if (response.data.status === 'success') {
      packages.value = response.data.data;
    } else {
      toast.error(response.data.message || 'Gagal mengambil daftar paket.');
    }
  } catch (error) {
    console.error('Error fetching packages:', error);
    toast.error('Terjadi kesalahan saat mengambil paket.');
  }
};

const openAddModal = () => {
  isEditMode.value = false;
  currentPackage.value = {
    id: null,
    package_name: '',
    price_monthly: 0,
    price_yearly: 0,
    duration_months: 0,
    max_employees: 0,
    features: '',
    is_active: true,
  };
  isModalOpen.value = true;
};

const openEditModal = (pkg) => {
  isEditMode.value = true;
  currentPackage.value = { ...pkg };
  isModalOpen.value = true;
};

const closeModal = () => {
  isModalOpen.value = false;
};

const handleSubmit = async () => {
  try {
    if (isEditMode.value) {
      await axios.put(`/api/superadmin/subscription-packages/${currentPackage.value.id}`, currentPackage.value);
      toast.success('Paket berhasil diperbarui!');
    } else {
      await axios.post('/api/superadmin/subscription-packages', currentPackage.value);
      toast.success('Paket berhasil ditambahkan!');
    }
    closeModal();
    fetchPackages();
  } catch (error) {
    console.error('Error saving package:', error);
    toast.error(error.response?.data?.message || 'Gagal menyimpan paket.');
  }
};

const deletePackage = async (id) => {
  if (confirm('Apakah Anda yakin ingin menghapus paket ini?')) {
    try {
      await axios.delete(`/api/superadmin/subscription-packages/${id}`);
      toast.success('Paket berhasil dihapus!');
      fetchPackages();
    } catch (error) {
      console.error('Error deleting package:', error);
      toast.error(error.response?.data?.message || 'Gagal menghapus paket.');
    }
  }
};

const formatCurrency = (value) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(value);
};

onMounted(() => {
  fetchPackages();
});
</script>
