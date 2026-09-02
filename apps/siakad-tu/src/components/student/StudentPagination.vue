<template>
  <div class="row justify-between items-center q-pa-md">
    <!-- Info Data -->
    <div class="text-caption text-grey-7">
      Menampilkan {{ startItem }} - {{ endItem }} dari {{ total }} data
    </div>

    <!-- Pagination Controls -->
    <div class="row items-center q-gutter-sm">
      <q-btn flat dense round icon="chevron_left" :disable="page <= 1" @click="$emit('prev')" />

      <q-pagination v-model="currentPage" :max="totalPages" :max-pages="5" direction-links flat color="primary"
        active-color="primary" @update:model-value="$emit('goToPage', $event)" />

      <q-btn flat dense round icon="chevron_right" :disable="page >= totalPages" @click="$emit('next')" />
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  page: { type: Number, required: true },
  totalPages: { type: Number, required: true },
  total: { type: Number, required: true },
  limit: { type: Number, required: true },
  hasMore: { type: Boolean, default: false }
});

defineEmits(['next', 'prev', 'goToPage']);

// Computed untuk mengikat nilai ke q-pagination tanpa memicu update langsung
const currentPage = computed({
  get: () => props.page,
  set: () => { } // Update ditangani oleh event @update:model-value
});

const startItem = computed(() => props.total === 0 ? 0 : (props.page - 1) * props.limit + 1);
const endItem = computed(() => Math.min(props.page * props.limit, props.total));
</script>
