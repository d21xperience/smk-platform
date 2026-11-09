<template>
  <div>
    <q-tabs v-model="tab" dense class="text-grey" active-color="primary" indicator-color="primary" align="left"
      narrow-indicator>
      <q-tab name="dokumen" label="Dokumen Aset" icon="folder" />
      <q-tab name="riwayat" label="Riwayat Mutasi" icon="history" />
    </q-tabs>

    <q-separator />

    <q-tab-panels v-model="tab" animated>
      <!-- TAB DOKUMEN -->
      <q-tab-panel name="dokumen">
        <div class="row items-center justify-between q-mb-md">
          <div class="text-subtitle1 text-weight-bold">Dokumen Aset: {{ assetCode }}</div>
          <q-btn color="primary" icon="cloud_upload" label="Upload Dokumen" @click="bukaFormUpload" size="sm" no-caps />
        </div>

        <q-table :rows="store.documents" :columns="columns" row-key="id" flat bordered :loading="store.isLoading"
          :rows-per-page-options="[5, 10]">
          <template v-slot:body-cell-document_type="props">
            <q-td :props="props"><q-chip dense outline color="primary" class="text-weight-bold">{{ props.value
                }}</q-chip></q-td>
          </template>
          <template v-slot:body-cell-status="props">
            <q-td :props="props" class="text-center"><q-badge :color="getStatusColor(props.value)">{{ props.value
                }}</q-badge></q-td>
          </template>
          <template v-slot:body-cell-aksi="props">
            <q-td :props="props" class="text-center">
              <div class="row justify-center q-gutter-xs">
                <q-btn dense round size="sm" color="primary" icon="download"
                  @click="unduhDokumen(props.row.id)"><q-tooltip>Unduh</q-tooltip></q-btn>
                <q-btn v-if="props.row.status === 'UPLOADED'" dense round size="sm" color="positive" icon="check"
                  @click="verifikasi(props.row.id)"><q-tooltip>Verifikasi</q-tooltip></q-btn>
                <q-btn v-if="props.row.status === 'UPLOADED'" dense round size="sm" color="negative" icon="close"
                  @click="bukaTolak(props.row.id)"><q-tooltip>Tolak</q-tooltip></q-btn>
              </div>
            </q-td>
          </template>
        </q-table>
      </q-tab-panel>

      <!-- TAB RIWAYAT MUTASI -->
      <q-tab-panel name="riwayat">
        <div class="row items-center justify-between q-mb-md">
          <div class="text-subtitle1 text-weight-bold">Riwayat Perubahan Aset</div>
          <q-btn color="secondary" icon="swap_horiz" label="Mutasi Aset" @click="bukaFormMutasi" size="sm" no-caps />
        </div>

        <q-timeline color="primary">
          <q-timeline-entry v-for="(h, idx) in kirStore.history" :key="idx" :title="h.action"
            :subtitle="new Date(h.changed_at).toLocaleString('id-ID')"
            :icon="h.action === 'CREATED' ? 'add_circle' : h.action === 'DIHAPUS' ? 'delete' : 'edit'"
            :color="h.action === 'CREATED' ? 'green' : h.action === 'DIHAPUS' ? 'red' : 'blue'">
            <div>{{ h.description || 'Perubahan data oleh user' }}</div>
          </q-timeline-entry>
        </q-timeline>
        <div v-if="kirStore.history.length === 0" class="text-center text-grey q-pa-md">Belum ada riwayat mutasi.</div>
      </q-tab-panel>
    </q-tab-panels>

    <!-- Dialog Upload Dokumen -->
    <q-dialog v-model="showForm" persistent>
      <q-card style="min-width: 500px">
        <q-card-section class="bg-primary text-white">
          <div class="text-h6">Upload Dokumen Baru</div>
        </q-card-section>
        <q-card-section class="q-gutter-y-md">
          <div class="row q-col-gutter-sm">
            <div class="col-6"><q-select outlined dense v-model="form.document_type" :options="docTypes"
                label="Jenis Dokumen *" /></div>
            <div class="col-6"><q-select outlined dense v-model="form.visibility"
                :options="['PUBLIC', 'INTERNAL', 'RESTRICTED']" label="Hak Akses *" /></div>
          </div>
          <q-input outlined dense v-model="form.title" label="Judul Dokumen *" />
          <q-input outlined dense v-model="form.document_number" label="No. Dokumen (Opsional)" />
          <q-input outlined dense v-model="form.description" type="textarea" label="Keterangan (Opsional)" />
          <q-file v-model="file" label="Pilih File Dokumen *" outlined dense
            accept=".pdf,.jpg,.jpeg,.png,.doc,.docx,.xls,.xlsx">
            <template v-slot:prepend><q-icon name="attach_file" /></template>
          </q-file>
        </q-card-section>
        <q-card-actions align="right" class="q-pa-md bg-grey-2">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="primary" label="Upload" @click="submitUpload" :loading="isSubmitting" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog Mutasi Aset -->
    <q-dialog v-model="showMutasi">
      <q-card style="width: 450px">
        <q-card-section class="bg-secondary text-white">
          <div class="text-h6">Form Mutasi Aset</div>
        </q-card-section>
        <q-card-section class="q-gutter-y-md">
          <q-select outlined dense v-model="mutasiForm.action"
            :options="['PINDAH_RUANG', 'UBAH_KONDISI', 'PENGHAPUSAN']" label="Jenis Mutasi *" />

          <q-select v-if="mutasiForm.action === 'PINDAH_RUANG'" outlined dense v-model="mutasiForm.new_room_id"
            :options="kirStore.rooms" option-value="id" option-label="nama" label="Pilih Ruang Baru" emit-value
            map-options />

          <q-select v-if="mutasiForm.action === 'UBAH_KONDISI' || mutasiForm.action === 'PENGHAPUSAN'" outlined dense
            v-model="mutasiForm.new_condition" :options="['Bagus', 'Rusak Ringan', 'Rusak Berat']"
            label="Kondisi Baru" />

          <q-input outlined dense v-model="mutasiForm.reason" type="textarea" label="Alasan / Keterangan *" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="secondary" label="Simpan Mutasi" @click="submitMutasi" :loading="isSubmitting" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog Tolak Dokumen -->
    <q-dialog v-model="showTolak">
      <q-card style="width: 400px">
        <q-card-section class="text-h6">Tolak Dokumen</q-card-section>
        <q-card-section><q-input outlined dense type="textarea" v-model="rejectReason"
            label="Alasan Penolakan *" /></q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="negative" label="Tolak" @click="submitTolak" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useInventoryDocumentStore } from '@/stores/sarpras/inventoryDocument'
import { useKirStore } from '@/stores/sarpras/kir'

const props = defineProps({
  assetId: { type: String, required: true },
  assetCode: { type: String, required: true }
})

const $q = useQuasar()
const store = useInventoryDocumentStore()
const kirStore = useKirStore()

const tab = ref('dokumen')
const showForm = ref(false)
const showMutasi = ref(false)
const showTolak = ref(false)
const isSubmitting = ref(false)
const file = ref(null)
const rejectReason = ref('')
const selectedDocId = ref('')

const docTypes = ['FAKTUR', 'BAST', 'SURAT_JALAN', 'GARANSI', 'FOTO_ASET', 'MANUAL_BOOK', 'DOKUMEN_HIBAH', 'LAINNYA']
const form = ref({ document_type: 'FOTO_ASET', visibility: 'INTERNAL', title: '', document_number: '', description: '' })
const mutasiForm = ref({ action: 'PINDAH_RUANG', new_room_id: '', new_condition: '', reason: '' })

const columns = [
  { name: 'document_type', label: 'Jenis', field: 'document_type', align: 'left' },
  { name: 'title', label: 'Judul', field: 'title', align: 'left' },
  { name: 'status', label: 'Status', field: 'status', align: 'center' },
  { name: 'aksi', label: 'Aksi', field: 'id', align: 'center' }
]

const getStatusColor = (s) => ({ 'UPLOADED': 'amber', 'VERIFIED': 'green', 'REJECTED': 'red', 'ARCHIVED': 'grey' }[s] || 'grey')

onMounted(() => {
  store.fetchDocuments('INVENTORY_ITEM', props.assetId)
  kirStore.fetchHistory(props.assetId)
  kirStore.fetchRooms() // Untuk pilihan mutasi ruang
})

function bukaFormUpload() {
  form.value = { document_type: 'FOTO_ASET', visibility: 'INTERNAL', title: '', document_number: '', description: '' }
  file.value = null
  showForm.value = true
}

async function submitUpload() {
  if (!form.value.title || !file.value) { $q.notify({ type: 'negative', message: 'Judul dan File wajib diisi' }); return }
  isSubmitting.value = true
  try {
    const fd = new FormData()
    fd.append('file', file.value)
    fd.append('document_type', form.value.document_type)
    fd.append('visibility', form.value.visibility)
    fd.append('title', form.value.title)
    fd.append('document_number', form.value.document_number)
    fd.append('description', form.value.description)
    fd.append('entity_type', 'INVENTORY_ITEM')
    fd.append('entity_id', props.assetId)

    await store.uploadDocument(fd)
    $q.notify({ type: 'positive', message: 'Dokumen berhasil diupload' })
    showForm.value = false
    store.fetchDocuments('INVENTORY_ITEM', props.assetId)
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Gagal upload dokumen' })
  } finally { isSubmitting.value = false }
}

async function unduhDokumen(id) {
  try {
    const url = await store.getDownloadUrl(id)
    window.open(url, '_blank')
  } catch (e) { $q.notify({ type: 'negative', message: `Gagal mendapatkan link dokumen ${e}` }) }
}

async function verifikasi(id) {
  try {
    await store.verifyDocument(id)
    $q.notify({ type: 'positive', message: 'Dokumen diverifikasi' })
    store.fetchDocuments('INVENTORY_ITEM', props.assetId)
  } catch (e) { $q.notify({ type: 'negative', message: `Gagal verifikasi${e}` }) }
}

function bukaTolak(id) {
  selectedDocId.value = id
  rejectReason.value = ''
  showTolak.value = true
}

async function submitTolak() {
  if (!rejectReason.value) { $q.notify({ type: 'negative', message: 'Alasan wajib diisi' }); return }
  try {
    await store.rejectDocument(selectedDocId.value, rejectReason.value)
    $q.notify({ type: 'info', message: 'Dokumen ditolak' })
    showTolak.value = false
    store.fetchDocuments('INVENTORY_ITEM', props.assetId)
  } catch (e) { $q.notify({ type: 'negative', message: `Gagal menolak ${e}` }) }
}

function bukaFormMutasi() {
  mutasiForm.value = { action: 'PINDAH_RUANG', new_room_id: '', new_condition: '', reason: '' }
  showMutasi.value = true
}

async function submitMutasi() {
  if (!mutasiForm.value.reason) { $q.notify({ type: 'negative', message: 'Alasan wajib diisi' }); return }
  isSubmitting.value = true
  try {
    const payload = {
      action: mutasiForm.value.action,
      new_room_id: mutasiForm.value.action === 'PINDAH_RUANG' ? mutasiForm.value.new_room_id : '',
      new_condition: mutasiForm.value.new_condition,
      status_aset: mutasiForm.value.action === 'PENGHAPUSAN' ? 'DIHAPUS' : 'AKTIF',
      reason: mutasiForm.value.reason
    }
    await kirStore.mutateItem(props.assetId, payload)
    $q.notify({ type: 'positive', message: 'Mutasi berhasil dicatat' })
    showMutasi.value = false
    kirStore.fetchHistory(props.assetId)
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Gagal mutasi' })
  } finally { isSubmitting.value = false }
}
</script>
