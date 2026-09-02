<!-- apps/siakad-tu/src/components/student/NisnAvailabilityIndicator.vue -->

<template>
  <div class="nisn-indicator">
    <q-spinner v-if="status === 'checking'" color="primary" size="sm" />
    <q-icon v-else-if="status === 'available'" name="check_circle" color="positive" size="sm" />
    <q-icon v-else-if="status === 'taken'" name="cancel" color="negative" size="sm" />
    <span v-if="status === 'available'" class="text-positive text-caption q-ml-xs">NISN tersedia</span>
    <span v-else-if="status === 'taken'" class="text-negative text-caption q-ml-xs">NISN sudah terdaftar</span>
  </div>
</template>

<script setup>
/**
 * NisnAvailabilityIndicator — Indicator untuk ketersediaan NISN.
 *
 * Props:
 * - status: 'idle' | 'checking' | 'available' | 'taken'
 *
 * TIDAK mengandung logika. Hanya menampilkan status dari parent.
 */
defineProps({
  status: {
    type: String,
    required: true,
    validator: (value) => ['idle', 'checking', 'available', 'taken'].includes(value)
  }
});
</script>

<style scoped>
.nisn-indicator {
  display: flex;
  align-items: center;
  min-height: 24px;
}
</style>
