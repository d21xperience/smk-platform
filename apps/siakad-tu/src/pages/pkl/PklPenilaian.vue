<template>
  <q-page class="q-pa-md">
    <div class="text-h5 q-mb-md">Penilaian PKL Terintegrasi</div>

    <div class="row q-col-gutter-md q-mb-md">
      <div class="col-12 col-md-4">
        <q-select v-model="selectedPenempatan" :options="penempatanOptions" label="Pilih Penempatan" option-value="id"
          option-label="label" outlined dense />
      </div>
    </div>

    <div v-if="selectedPenempatan">
      <q-table :rows="nilaiList" :columns="columns" row-key="id" flat bordered dense>
        <template v-slot:body-cell-aksi="props">
          <q-td :props="props">
            <q-btn flat round dense icon="edit" color="info" @click="openForm(props.row)" />
          </q-td>
        </template>
      </q-table>
    </div>

    <!-- Form penilaian -->
    <q-dialog v-model="formDialog" persistent>
      <q-card style="min-width: 400px">
        <q-card-section>
          <div class="text-h6">Input Nilai</div>
        </q-card-section>
        <q-card-section>
          <q-form @submit="submitNilai" class="q-gutter-md">
            <q-select v-model="form.penilai" :options="['guru', 'industri']" label="Penilai" outlined />
            <q-input v-model.number="form.nilai_hard" label="Nilai Hard Skill (0-100)" type="number" min="0" max="100"
              outlined />
            <q-input v-model.number="form.nilai_soft" label="Nilai Soft Skill (0-100)" type="number" min="0" max="100"
              outlined />
            <q-input v-model="form.catatan" label="Catatan" type="textarea" rows="2" outlined />
            <div class="q-gutter-sm">
              <q-btn label="Simpan" type="submit" color="primary" />
              <q-btn label="Batal" flat color="negative" v-close-popup />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { usePenilaianStore } from '@/stores/pkl/penilaian'
import { usePenempatanStore } from '@/stores/pkl/penempatan'
import { useSiswaStore } from '@/stores/kesiswaan/siswa'

const $q = useQuasar()
const penilaianStore = usePenilaianStore()
const penempatanStore = usePenempatanStore()
const siswaStore = useSiswaStore()

const selectedPenempatan = ref(null)
const formDialog = ref(false)
const editId = ref(null)
const form = ref({ penilai: 'guru', nilai_hard: 0, nilai_soft: 0, catatan: '' })

const penempatanOptions = computed(() => {
  return penempatanStore.list.map(p => ({
    id: p.id,
    label: `${p.siswa_nama} - ${p.mitra_nama} (${p.tanggal_mulai} s.d. ${p.tanggal_selesai})`,
    ...p
  }))
})

const nilaiList = computed(() => {
  if (!selectedPenempatan.value) return []
  return penilaianStore.list.filter(n => n.penempatan_id === selectedPenempatan.value.id)
})

const columns = [
  { name: 'penilai', label: 'Penilai', field: 'penilai', align: 'left' },
  { name: 'nilai_hard', label: 'Hard Skill', field: 'nilai_hard', align: 'center' },
  { name: 'nilai_soft', label: 'Soft Skill', field: 'nilai_soft', align: 'center' },
  { name: 'catatan', label: 'Catatan', field: 'catatan', align: 'left' },
  { name: 'aksi', label: 'Aksi', field: 'aksi', align: 'center' }
]

function openForm(row = null) {
  if (row) {
    editId.value = row.id
    form.value = { ...row }
  } else {
    editId.value = null
    form.value = { penilai: 'guru', nilai_hard: 0, nilai_soft: 0, catatan: '' }
  }
  formDialog.value = true
}

function submitNilai() {
  const data = {
    penempatan_id: selectedPenempatan.value.id,
    siswa_id: selectedPenempatan.value.siswa_id,
    ...form.value
  }
  if (editId.value) {
    penilaianStore.update(editId.value, data)
    $q.notify({ type: 'positive', message: 'Nilai diupdate' })
  } else {
    penilaianStore.tambah(data)
    $q.notify({ type: 'positive', message: 'Nilai ditambahkan' })
  }
  formDialog.value = false
}

onMounted(() => {
  penempatanStore.loadData()
  penilaianStore.loadData()
  siswaStore.loadData()
})
</script>
