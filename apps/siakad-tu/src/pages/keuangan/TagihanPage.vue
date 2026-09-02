<template>
  <q-page class="q-pa-md">
    <div class="text-h5 q-mb-md">Manajemen Tagihan Siswa</div>

    <!-- ========== 1. PENCARIAN SISWA ========== -->
    <q-card class="q-mb-md">
      <q-card-section>
        <div class="text-subtitle1">Cari Siswa</div>
        <q-input v-model="keyword" label="Nama atau NIS" outlined dense @keyup.enter="cariSiswa">
          <template v-slot:append>
            <q-icon name="search" class="cursor-pointer" @click="cariSiswa" />
          </template>
        </q-input>
      </q-card-section>
    </q-card>

    <!-- Hasil pencarian -->
    <div v-if="hasilPencarian.length > 0" class="q-mb-md">
      <div class="text-subtitle1">Hasil Pencarian:</div>
      <q-list bordered separator>
        <q-item v-for="siswa in hasilPencarian" :key="siswa.id" clickable @click="pilihSiswa(siswa)"
          :active="siswaTerpilih?.id === siswa.id" active-class="bg-primary text-white">
          <q-item-section>
            <q-item-label>{{ siswa.nama }}</q-item-label>
            <q-item-label caption>NIS: {{ siswa.nis }} | Kelas: {{ siswa.kelas }}</q-item-label>
          </q-item-section>
        </q-item>
      </q-list>
    </div>

    <!-- ========== 2. INFORMASI SISWA TERPILIH & TAGIHAN ========== -->
    <div v-if="siswaTerpilih">
      <q-banner class="bg-grey-2 q-mb-md">
        <template v-slot:avatar>
          <q-icon name="school" color="primary" />
        </template>
        <div>
          <strong>{{ siswaTerpilih.nama }}</strong> (NIS: {{ siswaTerpilih.nis }}) - {{ siswaTerpilih.kelas }}
          <q-btn flat round dense icon="refresh" @click="refreshTagihan" class="float-right" />
        </div>
      </q-banner>

      <!-- Tombol buat tagihan baru -->
      <div class="q-mb-md">
        <q-btn label="Buat Tagihan Baru" color="primary" icon="add" @click="openDialogTagihan" />
      </div>

      <!-- Tabel tagihan -->
      <q-table :rows="tagihanStore.items" :columns="tagihanColumns" row-key="id" :loading="tagihanStore.loading" flat
        bordered dense>
        <template v-slot:body-cell-kategori="props">
          <q-td>
            {{ getNamaKategori(props.row.kategoriId) }}
          </q-td>
        </template>

        <template v-slot:body-cell-nominal="props">
          <q-td align="right">
            Rp {{ props.row.nominal?.toLocaleString() }}
          </q-td>
        </template>

        <template v-slot:body-cell-status="props">
          <q-td>
            <q-badge :color="props.row.status === 'lunas' ? 'positive' : 'negative'">
              {{ props.row.status === 'lunas' ? 'LUNAS' : 'BELUM BAYAR' }}
            </q-badge>
          </q-td>
        </template>

        <template v-slot:body-cell-actions="props">
          <q-td auto-width>
            <q-btn v-if="props.row.status === 'belum'" label="Bayar" color="positive" dense size="sm"
              @click="bayarTagihan(props.row)" />
            <q-btn flat round dense icon="delete" color="negative" @click="hapusTagihan(props.row.id)"
              class="q-ml-sm" />
          </q-td>
        </template>
      </q-table>

      <!-- Info ringkasan -->
      <div class="q-mt-md text-right">
        <strong>Total Tagihan Belum Lunas:</strong> Rp {{ totalBelumLunas.toLocaleString() }}
      </div>
    </div>

    <!-- ========== 3. DIALOG BUAT / EDIT TAGIHAN ========== -->
    <q-dialog v-model="dialogTagihan" persistent>
      <q-card style="min-width: 450px">
        <q-card-section>
          <div class="text-h6">{{ isEdit ? 'Edit Tagihan' : 'Buat Tagihan Baru' }}</div>
        </q-card-section>

        <q-card-section class="q-gutter-md">
          <q-select v-model="formTagihan.kategoriId" :options="kategoriOptions" label="Kategori Tagihan"
            option-value="id" option-label="nama" emit-value map-options outlined dense
            :rules="[val => !!val || 'Kategori harus dipilih']" />

          <q-input v-model.number="formTagihan.nominal" label="Nominal (Rp)" type="number" outlined dense prefix="Rp "
            :rules="[val => val > 0 || 'Nominal harus lebih dari 0']" />

          <div class="row q-col-gutter-md">
            <div class="col-6">
              <q-input v-model="formTagihan.bulan" label="Bulan (opsional)" outlined dense />
            </div>
            <div class="col-6">
              <q-input v-model.number="formTagihan.tahun" label="Tahun" type="number" outlined dense />
            </div>
          </div>

          <q-input v-model="formTagihan.tanggalJatuhTempo" label="Jatuh Tempo" type="date" outlined dense />

          <q-input v-model="formTagihan.keterangan" label="Keterangan (opsional)" outlined dense autogrow />
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn flat label="Simpan" color="primary" @click="simpanTagihan" :loading="saving" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useRouter } from 'vue-router'
import { useTagihanStore } from 'stores/tagihanStore'
import { useKategoriTagihanStore } from 'stores/kategoriTagihanStore'
import { siswaService } from '../../services/siswaService'
// import { dummyTagihan } from '../../services/tagihanService'
const $q = useQuasar()
const router = useRouter()

// Pinia stores
const tagihanStore = useTagihanStore()
const kategoriStore = useKategoriTagihanStore()

// State lokal
const keyword = ref('')
const hasilPencarian = ref([])
const siswaTerpilih = ref(null)

const dialogTagihan = ref(false)
const isEdit = ref(false)
const saving = ref(false)

// Form tagihan
const formTagihan = ref({
  id: null,
  siswaId: null,
  kategoriId: null,
  nominal: null,
  bulan: '',
  tahun: new Date().getFullYear(),
  tanggalJatuhTempo: '',
  keterangan: '',
  status: 'belum'
})

// Columns tabel tagihan
const tagihanColumns = [
  { name: 'kategori', label: 'Kategori', field: 'kategoriId', align: 'left' },
  { name: 'bulan', label: 'Bulan', field: 'bulan', align: 'left' },
  { name: 'tahun', label: 'Tahun', field: 'tahun', align: 'center' },
  { name: 'nominal', label: 'Nominal', field: 'nominal', align: 'right' },
  { name: 'tanggalJatuhTempo', label: 'Jatuh Tempo', field: 'tanggalJatuhTempo', align: 'center' },
  { name: 'status', label: 'Status', field: 'status', align: 'center' },
  { name: 'actions', label: 'Aksi', align: 'center', field: 'actions' }
]

// Computed
const kategoriOptions = computed(() => kategoriStore.items)
const totalBelumLunas = computed(() => {
  return tagihanStore.items
    .filter(t => t.status === 'belum')
    .reduce((sum, t) => sum + t.nominal, 0)
})

// Helper
function getNamaKategori(kategoriId) {
  const kat = kategoriStore.items.find(k => k.id === kategoriId)
  return kat ? kat.nama : '-'
}

// Pencarian siswa
async function cariSiswa() {
  if (!keyword.value.trim()) {
    hasilPencarian.value = []
    return
  }
  try {
    const result = await siswaService.search(keyword.value)
    hasilPencarian.value = result
  } catch (err) {
    $q.notify({ type: 'negative', message: `Gagal mencari siswa ${err}` })
  }
}

// Pilih siswa
async function pilihSiswa(siswa) {
  siswaTerpilih.value = siswa
  await tagihanStore.fetchBySiswa(siswa.id)
}

// Refresh tagihan
async function refreshTagihan() {
  if (siswaTerpilih.value) {
    await tagihanStore.fetchBySiswa(siswaTerpilih.value.id)
    $q.notify({ type: 'info', message: 'Data tagihan diperbarui' })
  }
}

// Buka dialog buat tagihan
function openDialogTagihan() {
  if (!siswaTerpilih.value) {
    $q.notify({ type: 'warning', message: 'Pilih siswa terlebih dahulu' })
    return
  }
  isEdit.value = false
  formTagihan.value = {
    id: null,
    siswaId: siswaTerpilih.value.id,
    kategoriId: null,
    nominal: null,
    bulan: '',
    tahun: new Date().getFullYear(),
    tanggalJatuhTempo: '',
    keterangan: '',
    status: 'belum'
  }
  dialogTagihan.value = true
}

// Simpan tagihan (create)
async function simpanTagihan() {
  if (!formTagihan.value.kategoriId || !formTagihan.value.nominal || formTagihan.value.nominal <= 0) {
    $q.notify({ type: 'warning', message: 'Lengkapi data wajib (kategori & nominal)' })
    return
  }

  saving.value = true
  try {
    await tagihanStore.buatTagihan({
      siswaId: formTagihan.value.siswaId,
      kategoriId: formTagihan.value.kategoriId,
      nominal: formTagihan.value.nominal,
      bulan: formTagihan.value.bulan,
      tahun: formTagihan.value.tahun,
      tanggalJatuhTempo: formTagihan.value.tanggalJatuhTempo,
      keterangan: formTagihan.value.keterangan,
      status: 'belum'
    })
    $q.notify({ type: 'positive', message: 'Tagihan berhasil dibuat' })
    dialogTagihan.value = false
    await refreshTagihan() // refresh daftar tagihan
  } catch (err) {
    $q.notify({ type: 'negative', message: `Gagal menyimpan tagihan ${err}` })
  } finally {
    saving.value = false
  }
}

// Aksi bayar tagihan: redirect ke halaman pembayaran dengan membawa data tagihan
function bayarTagihan(tagihan) {
  // Simpan data tagihan ke store sementara atau query params
  // Redirect ke halaman pembayaran, misal /keuangan/pembayaran?tagihanId=...
  router.push({
    path: '/keuangan/pembayaran',
    query: { tagihanId: tagihan.id, siswaId: siswaTerpilih.value.id }
  })
}

// Hapus tagihan (hanya jika belum dibayar)
async function hapusTagihan(id) {
  const tagihan = tagihanStore.items.find(t => t.id === id)
  if (tagihan.status === 'lunas') {
    $q.notify({ type: 'warning', message: 'Tagihan yang sudah lunas tidak bisa dihapus' })
    return
  }
  $q.dialog({
    title: 'Konfirmasi',
    message: 'Hapus tagihan ini?',
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await tagihanStore.hapusTagihan(id)
      $q.notify({ type: 'positive', message: 'Tagihan dihapus' })
      await refreshTagihan()
    } catch (err) {
      $q.notify({ type: 'negative', message: `Gagal hapus tagihan ${err}` })
    }
  })
}

// Lifecycle
onMounted(async () => {
  await kategoriStore.fetchAll()
  // Jika ada query parameter siswaId, bisa langsung pilih siswa (opsional)
})
</script>

<style scoped>
/* optional */
</style>
