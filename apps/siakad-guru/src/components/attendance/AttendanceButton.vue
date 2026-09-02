<template>
  <div class="q-gutter-xs">
    <q-btn
      v-for="option in options"
      :key="option.value"
      :label="option.label"
      :color="isActive(option.value) ? option.color : 'grey-4'"
      :text-color="isActive(option.value) ? 'white' : 'dark'"
      size="sm"
      dense
      unelevated
      @click="selectStatus(option.value)"
    />
  </div>
</template>

<script setup>
import { defineProps, defineEmits } from 'vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: null, // null, 'PRESENT', 'SICK', 'PERMIT', 'ALPHA', 'LATE'
  },
})

const emit = defineEmits(['update:modelValue', 'change'])

const options = [
  { value: 'PRESENT', label: 'Hadir', color: 'positive' },
  { value: 'SICK', label: 'Sakit', color: 'warning' },
  { value: 'PERMIT', label: 'Izin', color: 'info' },
  { value: 'ALPHA', label: 'Alpha', color: 'negative' },
  { value: 'LATE', label: 'Terlambat', color: 'orange' },
]

const isActive = (value) => props.modelValue === value

const selectStatus = (value) => {
  const newValue = props.modelValue === value ? null : value // Toggle off jika klik sama
  emit('update:modelValue', newValue)
  emit('change', newValue)
}
</script>
