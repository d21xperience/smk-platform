<!-- apps/siakad-tu/src/components/student/StudentStatusBadge.vue -->

<template>
  <q-badge :color="statusColor" :label="statusLabel" />
</template>

<script setup>
import { computed } from 'vue';
import { formatStudentStatus, formatStudentStatusColor } from '@/composables/student/helpers/studentFormatters';

/**
 * StudentStatusBadge — Badge untuk menampilkan status siswa.
 * 
 * Props:
 * - status: string (ACTIVE, GRADUATED, TRANSFERRED, dll)
 * 
 * Menggunakan formatters dari composable untuk konsistensi.
 */
const props = defineProps({
  status: {
    type: String,
    required: true
  }
});

const statusLabel = computed(() => formatStudentStatus(props.status));
const statusColor = computed(() => formatStudentStatusColor(props.status));
</script>