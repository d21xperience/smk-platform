<template>
  <q-dialog :model-value="modelValue" @update:model-value="$emit('update:modelValue', $event)" persistent maximized>
    <q-card class="column full-height">
      <q-card-section class="row items-center q-pb-none bg-primary text-white">
        <div class="text-h6">Kelola Anggota: {{ classData?.name }}</div>
        <div class="q-ml-md text-subtitle2">
          ({{ kelasStore.members.length }} / {{ classData?.capacity }} Siswa)
        </div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup />
      </q-card-section>

      <q-card-section class="col q-pa-md row q-col-gutter-md">
        <!-- Kiri: Siswa Tersedia -->
        <div class="col-12 col-md-5 column">
          <q-input v-model="searchAvailable" placeholder="Cari siswa..." dense outlined @update:model-value="loadAvailable">
            <template v-slot:append><q-icon name="search" /></template>
          </q-input>
          <q-list bordered separator class="col q-mt-sm" style="overflow-y: auto; max-height: 60vh;">
            <q-item v-for="s in kelasStore.availableStudents" :key="s.id" clickable @click="selectedAvailable = s" :class="selectedAvailable?.id === s.id ? 'bg-blue-1' : ''">
              <q-item-section>
                <q-item-label>{{ s.name }}</q-item-label>
                <q-item-label caption>NIS: {{ s.nis }} | JK: {{ s.gender }}</q-item-label>
              </q-item-section>
            </q-item>
            <q-item v-if="kelasStore.availableStudents.length === 0 && !kelasStore.isLoadingAvailable">
              <q-item-section class="text-center text-grey">Tidak ada siswa tersedia</q-item-section>
            </q-item>
          </q-list>
        </div>

        <!-- Tengah: Tombol -->
        <div class="col-12 col-md-2 column flex-center q-gutter-y-md">
          <q-btn color="primary" icon="arrow_forward" label="Tambahkan" :disable="!selectedAvailable" @click="addSelected" :loading="loadingAdd" />
          <q-btn color="negative" icon="arrow_back" label="Hapus" :disable="!selectedMember" @click="removeSelected" :loading="loadingRemove" />
        </div>

        <!-- Kanan: Anggota Kelas -->
        <div class="col-12 col-md-5 column">
          <q-input v-model="searchMembers" placeholder="Cari anggota..." dense outlined @update:model-value="loadMembers">
            <template v-slot:append><q-icon name="search" /></template>
          </q-input>
          <q-list bordered separator class="col q-mt-sm" style="overflow-y: auto; max-height: 60vh;">
            <q-item v-for="s in filteredMembers" :key="s.studentId" clickable @click="selectedMember = s" :class="selectedMember?.studentId === s.studentId ? 'bg-red-1' : ''">
              <q-item-section>
                <q-item-label>{{ s.name }}</q-item-label>
                <q-item-label caption>NIS: {{ s.nis }}</q-item-label>
              </q-item-section>
            </q-item>
            <q-item v-if="filteredMembers.length === 0 && !kelasStore.isLoadingMembers">
              <q-item-section class="text-center text-grey">Belum ada anggota</q-item-section>
            </q-item>
          </q-list>
        </div>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useQuasar } from 'quasar'
import { useKelasStore } from '@/stores/kesiswaan/kelasStore'

const props = defineProps({
  modelValue: Boolean,
  classData: Object
})
const emit = defineEmits(['update:modelValue'])

const $q = useQuasar()
const kelasStore = useKelasStore()

const searchAvailable = ref('')
const searchMembers = ref('')
const selectedAvailable = ref(null)
const selectedMember = ref(null)
const loadingAdd = ref(false)
const loadingRemove = ref(false)

const filteredMembers = computed(() => {
  if (!searchMembers.value) return kelasStore.members
  const lower = searchMembers.value.toLowerCase()
  return kelasStore.members.filter(m => m.name.toLowerCase().includes(lower) || m.nis.includes(lower))
})

const loadMembers = () => {
  if (props.classData) kelasStore.fetchClassMembers(props.classData.id)
}

const loadAvailable = () => {
  kelasStore.fetchAvailableStudents({ keyword: searchAvailable.value, limit: 100 })
}

watch(() => props.modelValue, (isOpen) => {
  if (isOpen && props.classData) {
    selectedAvailable.value = null
    selectedMember.value = null
    searchAvailable.value = ''
    searchMembers.value = ''
    loadMembers()
    loadAvailable()
  }
})

const addSelected = async () => {
  if (!selectedAvailable.value) return

  // Validasi Frontend (Opsional, untuk mencegah request jika sudah terlihat penuh)
  if (props.classData.capacity > 0 && kelasStore.members.length >= props.classData.capacity) {
    $q.notify({ color: 'warning', message: 'Kapasitas kelas sudah penuh!', icon: 'warning' })
    return
  }

  loadingAdd.value = true
  try {
    await kelasStore.addMember(props.classData.id, selectedAvailable.value.id)
    $q.notify({ color: 'positive', message: 'Siswa berhasil ditambahkan', icon: 'check' })
    selectedAvailable.value = null
    loadMembers()
    loadAvailable() 
  } catch (err) {
    // TANGKAP ERROR 400 DARI BACKEND (Kapasitas Penuh)
    if (err.response?.status === 400 && err.response?.data?.error?.includes('capacity')) {
      $q.notify({ color: 'warning', message: 'Gagal: Kapasitas kelas sudah penuh!', icon: 'warning' })
    } else {
      $q.notify({ color: 'negative', message: err.response?.data?.error || 'Gagal menambah anggota' })
    }
  } finally {
    loadingAdd.value = false
  }
}

const removeSelected = async () => {
  if (!selectedMember.value) return
  loadingRemove.value = true
  try {
    await kelasStore.removeMember(props.classData.id, selectedMember.value.studentId)
    $q.notify({ color: 'positive', message: 'Siswa dikeluarkan dari kelas', icon: 'check' })
    selectedMember.value = null
    loadMembers()
    loadAvailable()
  } catch (err) {
    $q.notify({ color: 'negative', message: err.response?.data?.error || 'Gagal menghapus anggota' })
  } finally {
    loadingRemove.value = false
  }
}
</script>
