<template>
  <q-dialog :model-value="modelValue" @update:model-value="$emit('update:modelValue', $event)" persistent>
    <q-card style="min-width: 500px; border-radius: 12px;">
      <q-card-section class="row items-center q-pb-none">
        <div class="text-h6 text-weight-bold text-primary">Import {{ title }}</div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup />
      </q-card-section>
      
      <q-card-section class="q-pt-md">
        <!-- Banner Template -->
        <q-banner class="bg-blue-1 text-primary rounded-borders q-mb-md">
          <template v-slot:avatar>
            <q-icon name="info" color="primary" />
          </template>
          Gunakan template standar agar format kolom sesuai.
          <template v-slot:action>
            <q-btn 
              flat 
              color="primary" 
              label="Unduh Template" 
              icon="download" 
              @click="$emit('download-template')" 
            />
          </template>
        </q-banner>

        <!-- Input File -->
        <q-file 
          v-model="selectedFile" 
          :label="`Pilih File Excel (.xlsx, .xls)`" 
          outlined 
          dense 
          accept=".xlsx, .xls" 
          class="q-mb-md"
        >
          <template v-slot:prepend>
            <q-icon name="attach_file" />
          </template>
        </q-file>
        
        <div class="text-caption text-grey-7 q-mb-md">
          *Pilih file excel yang sudah diisi sesuai template.
        </div>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat label="Batal" color="grey" v-close-popup />
        <q-btn 
          color="primary" 
          label="Mulai Import" 
          icon="upload" 
          :disabled="!selectedFile" 
          :loading="loading" 
          @click="handleSubmit" 
        />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  modelValue: Boolean,
  title: {
    type: String,
    default: 'Data'
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'submit', 'download-template'])

const selectedFile = ref(null)

// Reset file saat modal ditutup
watch(() => props.modelValue, (isOpen) => {
  if (!isOpen) {
    selectedFile.value = null
  }
})

const handleSubmit = () => {
  if (selectedFile.value) {
    emit('submit', selectedFile.value)
  }
}
</script>
