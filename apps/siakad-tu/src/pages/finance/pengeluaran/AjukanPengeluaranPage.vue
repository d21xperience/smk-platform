<template>
  <q-page class="q-pa-md">
    <q-linear-progress v-if="store.pengajuanLoading" indeterminate color="teal" class="q-mb-sm" />

    <div class="row items-center justify-between q-mb-md">
      <div class="text-h5"><q-icon name="add_shopping_cart" class="q-mr-sm" />Ajukan Pengeluaran</div>
    </div>

    <div class="row q-col-gutter-md">
      <div class="col-12 col-md-7">
        <q-card>
          <q-card-section>
            <div class="text-subtitle1 text-weight-medium q-mb-sm">Form Pengajuan</div>
            <q-input v-model="form.keperluan" label="Keperluan *" outlined dense type="textarea" rows="3"
              :rules="[v => !!v || 'Wajib diisi']" hint="Jelaskan detail keperluan pengeluaran" />
            <q-input v-model.number="form.jumlah" label="Jumlah *" outlined dense type="number" prefix="Rp"
              :rules="[v => v > 0 || 'Harus > 0']" class="q-mt-sm" />
            <q-select v-model="form.prioritas" :options="prioritasOptions" label="Prioritas *" outlined dense
              emit-value map-options :rules="[v => !!v || 'Wajib diisi']" class="q-mt-sm" />
            <q-input v-model="form.catatan_pemohon" label="Catatan Pemohon" outlined dense type="textarea"
              rows="2" class="q-mt-sm" />
            <q-input v-model="form.tanggal" label="Tanggal (kosong = hari ini)" outlined dense type="date"
              class="q-mt-sm" />
          </q-card-section>
          <q-card-actions align="right">
            <q-btn flat label="Reset" @click="resetForm" />
            <q-btn color="teal" icon="send" label="Ajukan" :loading="store.pengajuanLoading" @click="submit" />
          </q-card-actions>
        </q-card>

        <q-card v-if="lastSubmitted" class="bg-green-1 q-mt-md">
          <q-card-section>
            <div class="text-positive text-weight-medium q-mb-sm">
              <q-icon name="check_circle" /> Pengajuan Berhasil Dikirim
            </div>
            <div class="text-caption">
              No. Pengajuan: <strong>{{ lastSubmitted.nomor_pengajuan }}</strong><br />
              Status: <q-chip dense color="orange" text-color="white">{{ lastSubmitted.status }}</q-chip>
            </div>
            <div class="text-caption q-mt-sm text-grey-7">
              Selanjutnya: menunggu verifikasi bendahara.
            </div>
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-md-5">
        <q-card>
          <q-card-section>
            <div class="text-subtitle2 text-grey-7 q-mb-sm">📋 Pengajuan Saya (Terakhir)</div>
            <q-list v-if="myRecent.length > 0" separator dense>
              <q-item v-for="p in myRecent" :key="p.id">
                <q-item-section>
                  <q-item-label>{{ p.nomor_pengajuan }}</q-item-label>
                  <q-item-label caption>{{ p.keperluan.substring(0, 40) }}...</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-chip dense :color="statusColor(p.status)" text-color="white" class="text-caption">
                    {{ statusLabel(p.status) }}
                  </q-chip>
                </q-item-section>
              </q-item>
            </q-list>
            <div v-else class="text-caption text-grey-6">Belum ada pengajuan.</div>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { usePengeluaranStore } from '@/stores/finance/pengeluaranStore'

const $q = useQuasar()
const store = usePengeluaranStore()

const lastSubmitted = ref(null)
const form = reactive({
  keperluan: '',
  jumlah: 0,
  prioritas: 'SEDANG',
  catatan_pemohon: '',
  tanggal: '',
})

const prioritasOptions = [
  { label: 'Rendah', value: 'RENDAH' },
  { label: 'Sedang', value: 'SEDANG' },
  { label: 'Tinggi', value: 'TINGGI' },
  { label: 'Darurat', value: 'DARURAT' },
]

const myRecent = computed(() => store.pengajuan.slice(0, 5))

function statusColor(s) {
  return {
    MENUNGGU_VERIFIKASI_BENDAHARA: 'orange',
    MENUNGGU_PERSETUJUAN_KEPALA: 'blue',
    DISETUJUI: 'positive',
    DITOLAK_BENDAHARA: 'negative',
    DITOLAK_KEPALA: 'negative',
    PERLU_REVISI: 'warning',
    CAIR: 'teal',
  }[s] || 'grey'
}
function statusLabel(s) {
  return {
    MENUNGGU_VERIFIKASI_BENDAHARA: 'Menunggu Bendahara',
    MENUNGGU_PERSETUJUAN_KEPALA: 'Menunggu Kepsek',
    DISETUJUI: 'Disetujui',
    DITOLAK_BENDAHARA: 'Ditolak Bendahara',
    DITOLAK_KEPALA: 'Ditolak Kepsek',
    PERLU_REVISI: 'Perlu Revisi',
    CAIR: 'Cair',
  }[s] || s
}

function resetForm() {
  form.keperluan = ''
  form.jumlah = 0
  form.prioritas = 'SEDANG'
  form.catatan_pemohon = ''
  form.tanggal = ''
  lastSubmitted.value = null
}

async function submit() {
  if (!form.keperluan || !form.jumlah) {
    $q.notify({ type: 'warning', message: 'Lengkapi keperluan & jumlah' })
    return
  }
  try {
    const res = await store.createPengajuan({ ...form })
    lastSubmitted.value = res.data
    $q.notify({ type: 'positive', message: `Pengajuan ${res.data.nomor_pengajuan} berhasil` })
    resetForm()
    await loadMyRecent()
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || e.message })
  }
}

async function loadMyRecent() {
  try {
    await store.listPengajuan({ pemohon_user_id: store.userID, limit: 5 })
  } catch {
    console.log('AjukanPengeluaran...')
  }
}

onMounted(loadMyRecent)
</script>
