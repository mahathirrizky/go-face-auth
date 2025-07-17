<template>
  <DataTable
    :value="data"
    :loading="loading"
    :paginator="true"
    :rows="10"
    :totalRecords="totalRecords"
    :lazy="lazy"
    @page="$emit('page', $event)"
    @filter="$emit('filter', $event)"
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
      :showFilterMenu="col.showFilterMenu !== false"
    >
      <template #body="slotProps">
        <slot :name="`column-${col.field}`" :item="slotProps.data">
          {{ slotProps.data[col.field] }}
        </slot>
      </template>
      <template #filter="{ filterModel }">
        <slot :name="`filter-${col.field}`" :filterModel="filterModel">
          <!-- Fallback to default text input filter if no custom template is provided -->
          <InputText v-if="col.showFilterMenu !== false" type="text" v-model="filterModel.value" class="p-column-filter" :placeholder="`Cari di ${col.header}`" />
        </slot>
      </template>
    </Column>
  </DataTable>
</template>

<script setup>
import { ref, watch } from 'vue';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import InputText from 'primevue/inputtext';
import IconField from 'primevue/iconfield';
import InputIcon from 'primevue/inputicon';

const props = defineProps({
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
  searchPlaceholder: {
    type: String,
    default: 'Cari...',
  },
  lazy: {
    type: Boolean,
    default: false,
  },
  totalRecords: {
    type: Number,
    default: 0,
  },
  // Accept filters from parent
  filters: {
    type: Object,
    required: true
  }
});

defineEmits(['page', 'filter', 'update:filters']);

// Use a local ref that is synced with the prop
const filters = ref(props.filters);
watch(() => props.filters, (newFilters) => {
  filters.value = newFilters;
});

</script>

<style scoped>
/* You can add component-specific styles here if needed */
</style>
