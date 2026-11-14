<template>
  <q-page class="q-pa-md">
    <q-linear-progress v-if="submitting" indeterminate color="teal" class="q-mb-sm" />

    <div class="row items-center justify-between q-mb-md">
      <div class="text-h5">
        <q-icon name="add_circle" class="q-mr-sm" />Input Penerimaan Iuran
      </div>
    </div>

    <div class="row q-col-gutter-md">
      <div class="col-12 col-md-7">
        <q-card>
          <q-card-section>
            <div class="text-subtitle1 text-weight-medium q-mb-sm">Data Pembayaran</div>

            <q-input v-model="form.nis" label="NIS Siswa *" outlined dense :rules="[v => !!v || 'Wajib diisi']"
              @blur="validateSiswa" :error="!!siswaError" :error-message="siswaError"
              hint="Masukkan NIS siswa aktif (validasi otomatis)">
              <template v-slot:append>
                <q-icon v-if="siswaInfo" name="check_circle" color="positive" />
                <q-icon v-else-if="siswaError" name="error" color="negative" />
              </template>
            </q-input>
            <div v-if="siswaInfo" class="q-mb-md text-caption text-positive">
              <q-icon name="person" /> {{ siswaInfo.nama_siswa }} — {{ siswaInfo.kelas_nama }}
            </div>

            <q-select v-model="form.jenis_iuran_id" :options="jenisOptions" label="Jenis Iuran *" outlined dense
              emit-value map-options :rules="[v => !!v || 'Wajib diisi']" @update:model-value="onJenisChange" />

            <div class="row q-col-gutter-sm q-mt-sm">
              <div class="col-6">
                <q-select v-model="form.periode_bulan" :options="bulanOptions" label="Bulan *" outlined dense
                  emit-value map-options :rules="[v => !!v || 'Wajib diisi']" />
              </div>
              <div class="col-6">
                <q-input v-model.number="form.periode_tahun" label="Tahun *" outlined dense type="number" />
              </div>
            </div>

            <q-input v-model.number="form.jumlah" label="Jumlah Bayar *" outlined dense type="number" prefix="Rp"
              :rules="[v => v > 0 || 'Harus > 0']" hint="Nominal harus sesuai tarif (jika ada)" />

            <q-select v-model="form.metode_bayar" :options="metodeOptions" label="Metode Bayar *" outlined dense
              emit-value map-options :rules="[v => !!v || 'Wajib diisi']" />

            <q-input v-if="form.metode_bayar !== 'TUNAI'" v-model="form.no_referensi"
              label="No. Referensi (transfer/QRIS)" outlined dense />

            <q-input v-model="form.keterangan" label="Keterangan" outlined dense type="textarea" rows="2" />

            <q-input v-model="form.tanggal" label="Tanggal Bayar (kosong = hari ini)" outlined dense type="date" />
          </q-card-section>

          <q-card-actions align="right">
            <q-btn flat label="Reset" @click="resetForm" />
            <q-btn color="teal" icon="save" label="Proses Pembayaran" :loading="submitting" @click="submit" />
          </q-card-actions>
        </q-card>
      </div>

      <div class="col-12 col-md-5">
        <q-card v-if="result" class="bg-green-1 q-mb-md">
          <q-card-section>
            <div class="text-subtitle1 text-weight-medium text-positive q-mb-sm">
              <q-icon name="check_circle" /> Pembayaran Berhasil
            </div>
            <q-list dense separator>
              <q-item>
                <q-item-section><div class="text-caption">No. BKM</div>
                  <div class="text-subtitle2 text-weight-bold">{{ result.penerimaan.no_bkm }}</div>
                </q-item-section>
              </q-item>
              <q-item>
                <q-item-section><div class="text-caption">Jumlah</div>
                  <div class="text-subtitle2 text-weight-bold">{{ formatRp(result.penerimaan.jumlah) }}</div>
                </q-item-section>
              </q-item>
              <q-item>
                <q-item-section><div class="text-caption">Status Tagihan</div>
                  <q-chip dense :color="result.tagihan.status === 'LUNAS' ? 'positive' : 'warning'" text-color="white">
                    {{ result.tagihan.status }}
                  </q-chip>
                </q-item-section>
              </q-item>
              <q-item v-if="result.bukti_kas">
                <q-item-section>
                  <div class="text-caption">BKU Entry</div>
                  <div class="text-caption">Debet: {{ formatRp(result.bukti_kas.debet) }} ·
                    Saldo: {{ formatRp(result.bukti_kas.saldo_berjalan) }}</div>
                </q-item-section>
              </q-item>
            </q-list>
            <q-btn class="q-mt-md full-width" outline color="positive" label="Input Pembayaran Baru" @click="resetForm" />
          </q-card-section>
        </q-card>

        <q-card v-if="!result">
          <q-card-section>
            <div class="text-subtitle2 text-grey-7 q-mb-sm">ℹ️ Info</div>
            <ul class="text-caption text-grey-8 q-pl-md">
              <li>Sistem otomatis generate <strong>BKM</strong></li>
              <li>Otomatis insert ke <strong>Buku Kas Umum</strong> (debet)</li>
              <li>Otomatis update status tagihan → <strong>LUNAS</strong> / <strong>SEBAGIAN</strong></li>
              <li>Nominal divalidasi terhadap tarif (jika ada)</li>
            </ul>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useTransaksiStore } from '@/stores/finance/transaksiStore'
import { useMasterStore } from '@/stores/finance/masterStore'
import { api } from '@/boot/axios'

const $q = useQuasar()
const store = useTransaksiStore()
const masterStore = useMasterStore()

const submitting = ref(false)
const siswaInfo = ref(null)
const siswaError = ref('')
const result = ref(null)

const now = new Date()
const form = reactive({
  nis: '',
  jenis_iuran_id: '',
  periode_bulan: now.getMonth() + 1,
  periode_tahun: now.getFullYear(),
  jumlah: 0,
  metode_bayar: 'TUNAI',
  no_referensi: '',
  keterangan: '',
  tanggal: '',
})

const bulanOptions = [
  'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
  'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
].map((label, i) => ({ label, value: i + 1 }))

const metodeOptions = [
  { label: 'Tunai', value: 'TUNAI' },
  { label: 'Transfer Bank', value: 'TRANSFER' },
  { label: 'QRIS', value: 'QRIS' },
]

const jenisOptions = computed(() =>
  masterStore.jenisIuran.items.filter(j => j.is_aktif).map(j => ({
    label: `${j.kode} — ${j.nama}`,
    value: j.id,
  }))
)

function formatRp(n) {
  return n === null || n === undefined ? '-' : 'Rp ' + Number(n).toLocaleString('id-ID')
}

async function validateSiswa() {
  siswaInfo.value = null
  siswaError.value = ''
  if (!form.nis) return
  try {
    // Endpoint lookup siswa — pakai backend tu-core
    const res = await api.get('/students', { params: { nis: form.nis, limit: 1 } })
    const list = res.data?.data || res.data || []
    const found = Array.isArray(list) ? list[0] : list
    if (found) {
      siswaInfo.value = {
        nama_siswa: found.nm_siswa || found.nama || found.name || '(tanpa nama)',
        kelas_nama: found.kelas_nama || found.nm_kelas || '',
      }
    } else {
      siswaError.value = 'NIS tidak ditemukan'
    }
  } catch {
    // Silent — backend akan validasi saat submit
  }
}

function onJenisChange() {
  // Auto-fill nominal dari tarif jika tersedia
  if (!form.jenis_iuran_id) return
  const tarif = masterStore.tarifIuran.items.find(
    t => t.jenis_iuran_id === form.jenis_iuran_id && (t.kelas === '' || t.kelas === siswaInfo.value?.kelas_nama)
  )
  if (tarif) form.jumlah = tarif.nominal
}

function resetForm() {
  form.nis = ''
  form.jenis_iuran_id = ''
  form.periode_bulan = now.getMonth() + 1
  form.periode_tahun = now.getFullYear()
  form.jumlah = 0
  form.metode_bayar = 'TUNAI'
  form.no_referensi = ''
  form.keterangan = ''
  form.tanggal = ''
  siswaInfo.value = null
  siswaError.value = ''
  result.value = null
}

async function submit() {
  if (!form.nis || !form.jenis_iuran_id || !form.jumlah) {
    $q.notify({ type: 'warning', message: 'Lengkapi: NIS, Jenis Iuran, Jumlah' })
    return
  }
  submitting.value = true
  try {
    const payload = {
      nis: form.nis,
      jenis_iuran_id: form.jenis_iuran_id,
      periode_bulan: form.periode_bulan,
      periode_tahun: form.periode_tahun,
      jumlah: form.jumlah,
      metode_bayar: form.metode_bayar,
      no_referensi: form.no_referensi || undefined,
      keterangan: form.keterangan || undefined,
      tanggal: form.tanggal || undefined,
    }
    const res = await store.createPenerimaan(payload)
    result.value = res.data
    $q.notify({ type: 'positive', message: `BKM: ${res.data?.penerimaan?.no_bkm || 'OK'}` })
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || e.message, timeout: 5000 })
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  if (masterStore.jenisIuran.items.length === 0) {
    await masterStore.jenisIuran.list({ limit: 100 })
  }
  if (masterStore.tarifIuran.items.length === 0) {
    await masterStore.tarifIuran.list({ limit: 100 })
  }
})
</script>
