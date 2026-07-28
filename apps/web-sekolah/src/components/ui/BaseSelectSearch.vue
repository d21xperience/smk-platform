<!-- <template>
  <q-select v-model="modelValueComputed" :options="filteredOptions" :option-label="optionLabel"
    :option-value="optionValue" :label="label" :disable="disable" :rules="rules" emit-value map-options outlined dense
    use-input fill-input hide-selected input-debounce="300" @filter="filterFn" />
</template>

<script setup>
import { ref, computed } from 'vue'

// 1. Definisikan Props untuk menerima data dari luar
const props = defineProps({
  modelValue: { type: [String, Number, Object], default: null },
  options: { type: Array, required: true, default: () => [] },
  optionLabel: { type: String, default: 'label' },
  optionValue: { type: String, default: 'value' },
  label: { type: String, default: 'Pilih Data' },
  disable: { type: Boolean, default: false },
  rules: { type: Array, default: () => [] }
})

// 2. Definisikan Emits untuk update v-model ke parent
const emit = defineEmits(['update:modelValue'])

// 3. Computed properti agar v-model sinkron antara child dan parent
const modelValueComputed = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

// 4. State lokal khusus untuk menampung hasil filter pencarian
const filteredOptions = ref([])

// 5. Fungsi filter generic (bisa mencari key apa saja sesuai props optionLabel)
const filterFn = (val, update) => {
  if (val === '') {
    update(() => {
      filteredOptions.value = props.options
    })
    return
  }

  update(() => {
    const needle = val.toLowerCase()
    filteredOptions.value = props.options.filter(item => {
      // Mengambil teks label secara dinamis berdasarkan properti optionLabel
      const labelText = item[props.optionLabel] ? String(item[props.optionLabel]) : ''
      return labelText.toLowerCase().indexOf(needle) > -1
    })
  })
}
</script> -->
<template>
  <q-select ref="selectRef" v-model="modelValueComputed" :options="filteredOptions" :option-label="optionLabel"
    :option-value="optionValue" :label="label" :disable="disable" :rules="rules" emit-value map-options outlined dense
    use-input fill-input hide-selected input-debounce="300" @filter="filterFn" @input-value="setInputValue" />
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const props = defineProps({
  modelValue: { type: [String, Number, Object], default: null },
  options: { type: Array, required: true, default: () => [] },
  optionLabel: { type: String, default: 'label' },
  optionValue: { type: String, default: 'value' },
  label: { type: String, default: 'Pilih Data' },
  disable: { type: Boolean, default: false },
  rules: { type: Array, default: () => [] }
})

const emit = defineEmits(['update:modelValue'])

const selectRef = ref(null)
const filteredOptions = ref([])

// Computed properti untuk sinkronisasi v-model
const modelValueComputed = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

// Fungsi bantu untuk memperbarui teks yang terlihat di dalam q-select
const setInputValue = (val) => {
  if (selectRef.value) {
    selectRef.value.updateInputValue(val)
  }
}

// 🚀 SOLUSI UTAMA: Perhatikan perubahan v-model saat dialog dibuka dalam mode edit
watch(
  () => props.modelValue,
  (newVal) => {
    if (newVal && props.options.length > 0) {
      // Cari objek siswa di dalam list berdasarkan ID yang dikirim
      const selectedItem = props.options.find(item => item[props.optionValue] === newVal)
      if (selectedItem) {
        // Paksa input teks menampilkan nama siswa (bukan ID)
        setTimeout(() => {
          setInputValue(selectedItem[props.optionLabel])
        }, 50)
      }
    } else if (!newVal) {
      setInputValue('')
    }
  },
  { immediate: true }
)

// Fungsi pencarian teks input
const filterFn = (val, update) => {
  if (val === '') {
    update(() => {
      filteredOptions.value = props.options
    })
    return
  }

  update(() => {
    const needle = val.toLowerCase()
    filteredOptions.value = props.options.filter(item => {
      const labelText = item[props.optionLabel] ? String(item[props.optionLabel]) : ''
      return labelText.toLowerCase().indexOf(needle) > -1
    })
  })
}
</script>
