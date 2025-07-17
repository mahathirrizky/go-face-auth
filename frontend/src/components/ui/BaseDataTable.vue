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
    v-model:filters="localFilters"
    class="p-datatable-customers"
    paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink CurrentPageReport RowsPerPageDropdown"
    :rowsPerPageOptions="[10, 25, 50]"
    currentPageReportTemplate="Menampilkan {first} sampai {last} dari {totalRecords} data"
    
    :editMode="editable ? 'row' : null"
    :dataKey="editKey"
    v-model:editingRows="editingRows"
    @row-edit-save="onRowEditSave"
    @row-edit-cancel="onRowEditCancel"
  >
    <template #header>
      <div class="flex justify-between items-center">
        <IconField iconPosition="left">
          <InputIcon class="pi pi-search"></InputIcon>
          <InputText v-model="localFilters.global.value" :placeholder="searchPlaceholder" />
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
      <template #editor="{ data, field }">
        <slot :name="`editor-${col.field}`" :data="data" :field="field">
            <!-- Fallback to a simple text input if no specific editor is provided -->
            <InputText v-model="data[field]" class="w-full" />
        </slot>
      </template>
    </Column>

    <Column v-if="editable" :rowEditor="true" style="width: 10%; min-width: 8rem" bodyStyle="text-align:center"></Column>

  </DataTable>
</template>

<script setup>
import { ref, watch, computed } from 'vue';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import InputText from 'primevue/inputtext';
import IconField from 'primevue/iconfield';
import InputIcon from 'primevue/inputicon';
import { FilterMatchMode } from '@primevue/core/api';

const props = defineProps({
  data: { type: Array, required: true },
  columns: { type: Array, required: true },
  loading: { type: Boolean, default: false },
  searchPlaceholder: { type: String, default: 'Cari...' },
  lazy: { type: Boolean, default: false },
  totalRecords: { type: Number, default: 0 },
  filters: {
    type: Object,
    default: () => ({ global: { value: null, matchMode: FilterMatchMode.CONTAINS } })
  },
  // Props for Row Editing
  editable: { type: Boolean, default: false },
  editKey: { type: String, default: 'id' }
});

const emit = defineEmits(['page', 'filter', 'update:filters', 'save', 'cancel']);

const localFilters = computed({
  get: () => props.filters,
  set: (value) => emit('update:filters', value)
});

const editingRows = ref([]);

const onRowEditSave = (event) => {
  emit('save', event);
};

const onRowEditCancel = (event) => {
  emit('cancel', event);
};
</script>

<style scoped>
/* You can add component-specific styles here if needed */
</style>
