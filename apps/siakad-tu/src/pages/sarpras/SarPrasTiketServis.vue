<template>
  <q-page class="q-pa-md">
    <!-- PANEL FILTER TAB STATUS -->
    <div class="row items-center justify-between q-mb-md">
      <q-tabs v-model="tabAktif" dense class="text-grey shadow-1 rounded-borders bg-white" active-color="primary"
        indicator-color="primary" align="left" narrow-indicator>
        <q-tab name="Semua" label="Semua Laporan" />
        <q-tab name="Pending" class="text-negative">
          <div class="row items-center no-wrap">
            <span>Baru (Pending)</span>
            <q-badge color="red" floating class="q-ml-xs">{{ hitungTiket('Pending') }}</q-badge>
          </div>
        </q-tab>
        <q-tab name="Proses" class="text-warning" label="Diproses" />
        <q-tab name="Selesai" class="text-positive" label="Selesai" />
      </q-tabs>

      <q-btn color="primary" icon="add_comment" label="Buat Laporan Baru" @click="bukaFormBaru" />
    </div>

    <!-- LIST KARTU TIKET LAPORAN -->
    <div class="row q-col-gutter-md">
      <div v-for="tiket in tiketTersaring" :key="tiket.id" class="col-12 col-md-6">
        <q-card flat bordered class="tiket-card">
          <div :class="`urgensi-bar bg-${getUrgensiColor(tiket.urgency)}`"></div>

          <q-card-section class="q-pl-lg">
            <div class="row items-center no-wrap justify-between">
              <div class="row items-center">
                <q-badge :color="getStatusColor(tiket.status)" class="q-mr-sm text-uppercase">
                  {{ tiket.status }}
                </q-badge>
                <span class="text-caption text-grey-7">{{ tiket.report_number }}</span>
              </div>

              <q-chip dense :color="getUrgensiColor(tiket.urgency) + '-1'"
                :text-color="getUrgensiColor(tiket.urgency) + '-9'" class="text-weight-bold">
                Urgensi: {{ tiket.urgency }}
              </q-chip>
            </div>

            <div class="text-h6 text-grey-9 q-mt-sm">{{ tiket.asset_name }}</div>
            <div class="text-body2 text-grey-7 q-mt-xs">
              <q-icon name="room" size="xs" /> <strong>Lokasi:</strong> {{ tiket.location }}
            </div>
            <div class="text-body2 text-grey-8 q-mt-sm bg-grey-2 q-pa-sm rounded-borders">
              "{{ tiket.description }}"
            </div>

            <q-separator class="q-my-md" />

            <div class="row items-center justify-between">
              <div class="text-caption text-grey-7">
                Pelapor: <strong>{{ tiket.reporter_name }}</strong>
                <div v-if="tiket.cost > 0" class="text-positive text-weight-bold q-mt-xs">
                  Biaya Servis: Rp {{ tiket.cost.toLocaleString('id-ID') }}
                </div>
              </div>

              <q-btn outline dense color="primary" icon="edit_note" label="Update Tindakan"
                @click="bukaDialogTindakan(tiket)" />
            </div>
          </q-card-section>
        </q-card>
      </div>

      <div v-if="tiketTersaring.length === 0" class="col-12 text-center q-pa-xl text-grey-5">
        <q-icon name="assignment_turned_in" size="64px" class="q-mb-sm" />
        <div class="text-h6">Tidak ada tiket laporan dalam kategori ini.</div>
      </div>
    </div>

    <!-- DIALOG BUAT LAPORAN BARU -->
    <q-dialog v-model="dialogBaruTerbuka" persistent>
      <q-card style="width: 500px; max-width: 90vw;">
        <q-card-section class="bg-primary text-white"><div class="text-h6">Buat Laporan Kerusakan</div></q-card-section>
        <q-card-section class="q-gutter-y-md">
          <q-select outlined dense v-model="formBaru.inventory_item_id" :options="store.activeAssets" option-value="id" option-label="name" label="Pilih Aset *" emit-value map-options>
            <template v-slot:option="scope">
              <q-item v-bind="scope.itemProps">
                <q-item-section>
                  <q-item-label>{{ scope.opt.name }}</q-item-label>
                  <q-item-label caption>Kode: {{ scope.opt.id }}</q-item-label>
                </q-item-section>
              </q-item>
            </template>
          </q-select>
          <q-input outlined dense v-model="formBaru.location" label="Lokasi Kerusakan *" />
          <q-input outlined dense v-model="formBaru.reporter_name" label="Nama Pelapor *" />
          <q-select outlined dense v-model="formBaru.urgency" :options="['Rendah', 'Sedang', 'Tinggi']" label="Tingkat Urgensi *" />
          <q-input outlined dense type="textarea" v-model="formBaru.description" label="Deskripsi Masalah *" rows="3" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="primary" label="Kirim Laporan" @click="simpanLaporanBaru" :loading="isSubmitting" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- DIALOG UPDATE TINDAKAN -->
    <q-dialog v-model="dialogTindakanTerbuka">
      <q-card style="width: 450px; max-width: 90vw;">
        <q-card-section class="bg-primary text-white row items-center">
          <div class="text-h6">Tindak Lanjut: {{ FormTindakan.report_number }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <div class="text-weight-bold text-grey-8">{{ FormTindakan.asset_name }}</div>
          <q-select outlined dense v-model="FormTindakan.status" :options="['Pending', 'Proses', 'Selesai']" label="Ubah Status Penanganan *" />
          <q-input outlined dense v-model="FormTindakan.vendor_name" label="Nama Vendor / Teknisi" />
          <q-input outlined dense type="number" v-model.number="FormTindakan.cost" label="Realisasi Biaya Perbaikan (Rp)" prefix="Rp" hint="Isi 0 jika diperbaiki mandiri" />
          <q-input outlined dense type="textarea" v-model="FormTindakan.resolution_notes" label="Catatan Solusi Perbaikan" rows="3" />
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="primary" label="Simpan Pembaruan" @click="simpanTindakan" :loading="isSubmitting" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useQuasar } from 'quasar'
import { useMaintenanceStore } from '@/stores/sarpras/maintenance'
import { useContextStore } from '@/stores/contextStore'

const $q = useQuasar()
const store = useMaintenanceStore()
const contextStore = useContextStore()

const tabAktif = ref('Semua')
const dialogTindakanTerbuka = ref(false)
const dialogBaruTerbuka = ref(false)
const isSubmitting = ref(false)

const formBaru = ref({ inventory_item_id: '', location: '', reporter_name: '', urgency: 'Sedang', description: '' })
const FormTindakan = ref({ id: '', report_number: '', asset_name: '', status: 'Pending', vendor_name: '', cost: 0, resolution_notes: '' })

const tiketTersaring = computed(() => {
  if (tabAktif.value === 'Semua') return store.damageReports
  return store.damageReports.filter(tiket => tiket.status === tabAktif.value)
})

function hitungTiket(status) {
  return store.damageReports.filter(tiket => tiket.status === status).length
}

function getUrgensiColor(urgensi) {
  if (urgensi === 'Tinggi') return 'red'
  if (urgensi === 'Sedang') return 'amber'
  return 'green'
}

function getStatusColor(status) {
  if (status === 'Pending') return 'grey-7'
  if (status === 'Proses') return 'orange'
  return 'green'
}

onMounted(() => {
  store.fetchDamageReports()
  store.fetchActiveAssets()
})

watch(() => contextStore.selectedSemesterId, (newVal) => {
  if (newVal) store.fetchDamageReports()
})

function bukaFormBaru() {
  formBaru.value = { inventory_item_id: '', location: '', reporter_name: '', urgency: 'Sedang', description: '' }
  dialogBaruTerbuka.value = true
}

async function simpanLaporanBaru() {
  if (!formBaru.value.inventory_item_id || !formBaru.value.location || !formBaru.value.reporter_name || !formBaru.value.description) {
    $q.notify({ type: 'negative', message: 'Mohon lengkapi semua field bertanda *' })
    return
  }
  isSubmitting.value = true
  try {
    const asset = store.activeAssets.find(a => a.id === formBaru.value.inventory_item_id)
    const payload = { ...formBaru.value, asset_name: asset ? asset.name : '' }
    await store.createDamageReport(payload)
    $q.notify({ type: 'positive', message: 'Laporan berhasil dikirim' })
    dialogBaruTerbuka.value = false
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Gagal mengirim laporan' })
  } finally {
    isSubmitting.value = false
  }
}

function bukaDialogTindakan(tiket) {
  FormTindakan.value = { ...tiket }
  dialogTindakanTerbuka.value = true
}

async function simpanTindakan() {
  isSubmitting.value = true
  try {
    await store.updateDamageReport(FormTindakan.value.id, {
      status: FormTindakan.value.status,
      vendor: FormTindakan.value.vendor_name,
      cost: FormTindakan.value.cost,
      notes: FormTindakan.value.resolution_notes
    })
    $q.notify({ type: 'positive', message: 'Tiket berhasil diperbarui' })
    dialogTindakanTerbuka.value = false
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Gagal memperbarui tiket' })
  } finally {
    isSubmitting.value = false
  }
}
</script>

<style scoped>
.tiket-card { position: relative; overflow: hidden; transition: transform 0.2s; }
.tiket-card:hover { transform: translateY(-2px); }
.urgensi-bar { position: absolute; top: 0; left: 0; bottom: 0; width: 6px; }
</style>
