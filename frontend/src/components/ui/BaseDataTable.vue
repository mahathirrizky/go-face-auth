<template>
  <DataTable
    :value="data"
    :loading="loading"
    :paginator="true"
    :rows="10"
    :globalFilterFields="globalFilterFields"
    v-model:filters="filters"
    class="p-datatable-customers"
    paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink CurrentPageReport RowsPerPageDropdown"
    :rowsPerPageOptions="[10, 25, 50]"
    currentPageReportTemplate="Menampilkan {first} sampai {last} dari {totalRecords} data"
  >
    <template #header>
      <div class="flex justify-between items-center">
        <IconField iconPosition="left">
          <InputIcon class="pi pi-search"></InputIcon>
          <InputText v-model="filters.global.value" :placeholder="searchPlaceholder" />
        </IconField>
        <div class="flex space-x-2">
          <slot name="header-actions"></slot>
        </div>
      </div>
    </template>
    <template #empty>
      Tidak ada data ditemukan.
    </template>
    <template #loading>
      Memuat data...
    </template>

    <Column
      v-for="col in columns"
      :key="col.field"
      :field="col.field"
      :header="col.header"
      :sortable="col.sortable !== false"
    >
      <template #body="slotProps">
        <slot :name="`column-${col.field}`" :item="slotProps.data">
          {{ slotProps.data[col.field] }}
        </slot>
      </template>
    </Column>
  </DataTable>
</template>

<script setup>
import { ref, defineProps } from 'vue';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import InputText from 'primevue/inputtext';
import IconField from 'primevue/iconfield';
import InputIcon from 'primevue/inputicon';
import { FilterMatchMode } from '@primevue/core/api'; // Corrected import

defineProps({
  data: {
    type: Array,
    required: true,
  },
  columns: {
    type: Array,
    required: true,
  },
  loading: {
    type: Boolean,
    default: false,
  },
  globalFilterFields: {
    type: Array,
    default: () => [],
  },
  searchPlaceholder: {
    type: String,
    default: 'Cari...',
  },
});

const filters = ref({
  global: { value: null, matchMode: FilterMatchMode.CONTAINS },
});
</script>

<style scoped>
/* You can add component-specific styles here if needed */
</style>
