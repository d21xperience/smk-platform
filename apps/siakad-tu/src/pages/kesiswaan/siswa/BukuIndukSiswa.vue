<template>
  <q-page class="q-pa-md bg-grey-2">
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Buku Induk Siswa</h1>
        <p class="text-caption text-grey-7 q-mb-none">Pencatatan data kronologis arsip statis kesiswaan SMK</p>
      </div>
      <q-btn icon="refresh" label="Muat Ulang" color="primary" outline :loading="store.loading" @click="loadList" />
    </div>

    <q-banner v-if="store.error" class="bg-red-1 text-red-9 q-mb-md" rounded>
      <q-icon name="error" class="q-mr-sm" />{{ store.error }}
    </q-banner>

    <div class="row q-col-gutter-md">
      <!-- Panel Kiri: Klaper -->
      <div class="col-12 col-md-4">
        <q-card class="shadow-1">
          <q-card-section class="q-pb-none">
            <div class="text-subtitle2 text-weight-bold text-grey-9 q-mb-sm">Indeks Klaper (Abjad)</div>
            <div class="flex q-gutter-xs justify-center q-mb-md">
              <q-btn
                v-for="letter in alphabets"
                :key="letter"
                :color="selectedLetter === letter ? 'primary' : 'grey-3'"
                :text-color="selectedLetter === letter ? 'white' : 'dark'"
                size="xs" padding="xs" dense
                style="min-width: 26px"
                @click="selectLetter(letter)"
              >{{ letter }}</q-btn>
            </div>
            <q-input v-model="searchKlaper" dense outlined clearable placeholder="Cari nama / NIS..." class="q-mb-sm">
              <template #prepend><q-icon name="search" /></template>
            </q-input>
          </q-card-section>
          <q-card-section class="q-pa-none" style="max-height: 480px; overflow-y: auto;">
            <q-list separator>
              <q-item
                v-for="siswa in filteredList"
                :key="siswa.id"
                clickable v-ripple
                :active="store.activeStudent?.id === siswa.id"
                active-class="bg-blue-1 text-primary text-weight-bold"
                @click="selectStudent(siswa)"
              >
                <q-item-section avatar>
                  <q-avatar color="primary" text-color="white" size="32px">
                    {{ (siswa.name || '?').charAt(0) }}
                  </q-avatar>
                </q-item-section>
                <q-item-section>
                  <q-item-label>{{ siswa.name }}</q-item-label>
                  <q-item-label caption>NIS: {{ siswa.nis }} | NISN: {{ siswa.nisn || '-' }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-icon name="chevron_right" />
                </q-item-section>
              </q-item>
              <q-item v-if="filteredList.length === 0 && !store.loading" class="text-center text-grey-6 q-pa-md">
                Tidak ada siswa berawalan "{{ selectedLetter }}"
              </q-item>
              <q-item v-if="store.loading" class="justify-center q-pa-md">
                <q-spinner color="primary" size="2em" />
              </q-item>
            </q-list>
          </q-card-section>
        </q-card>
      </div>

      <!-- Panel Kanan: Detail -->
      <div class="col-12 col-md-8">
        <q-card v-if="!store.activeStudent" class="shadow-1 fit flex flex-center text-center q-pa-xl">
          <div>
            <q-icon name="menu_book" size="64px" color="grey-4" />
            <div class="text-subtitle1 text-grey-6 q-mt-md">
              Silakan pilih nama siswa di panel Klaper untuk melihat Lembar Buku Induk.
            </div>
          </div>
        </q-card>

        <q-card v-else class="shadow-1">
          <q-card-section class="bg-primary text-white row items-center justify-between">
            <div class="row items-center q-gutter-md">
              <q-avatar size="60px" class="bg-white text-primary text-weight-bold">
                {{ (store.activeStudent.name || '?').charAt(0) }}
              </q-avatar>
              <div>
                <div class="text-h6 text-weight-bold">{{ store.activeStudent.name }}</div>
                <div class="text-caption text-blue-2">
                  No. Buku Induk: BI/{{ store.activeStudent.nis }}/SMK
                </div>
              </div>
            </div>
            <div class="row q-gutter-sm items-center">
              <q-badge :color="completenessColor" class="q-pa-sm">
                Kelengkapan: {{ store.activeStudent.skorKelengkapan }}% — {{ store.activeStudent.statusKelengkapan }}
              </q-badge>
              <q-btn
                icon="print" color="white" text-color="primary"
                label="Cetak Lembar Induk"
                :loading="printing"
                @click="printBukuInduk"
              />
            </div>
          </q-card-section>

          <q-tabs v-model="activeTab" dense class="text-grey" active-color="primary" indicator-color="primary" align="justify" narrow-indicator>
            <q-tab name="pribadi" label="Data Pribadi" icon="person" />
            <q-tab name="pendidikan" label="Pendidikan" icon="school" />
            <q-tab name="keluarga" label="Orang Tua/Wali" icon="family_restroom" />
            <q-tab name="pelengkap" label="Pelengkap" icon="assignment_ind" />
            <q-tab name="nilai" label="Nilai Rapor" icon="grading" />
          </q-tabs>
          <q-separator />

          <q-card-section class="q-pa-md">
            <q-tab-panels v-model="activeTab" animated>
              <!-- TAB PRIBADI -->
              <q-tab-panel name="pribadi" class="q-pa-none">
                <div class="row q-col-gutter-sm">
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Nama Lengkap" :model-value="store.activeStudent.name || '-'" dense /></div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="NIS" :model-value="store.activeStudent.nis || '-'" dense /></div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="NISN" :model-value="store.activeStudent.nisn || '-'" dense /></div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Jenis Kelamin" :model-value="genderLabel" dense /></div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Tempat Lahir" :model-value="store.activeStudent.birthPlace || '-'" dense /></div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Tanggal Lahir" :model-value="birthDateFormatted" dense /></div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Agama" :model-value="store.activeStudent.agama || '-'" dense /></div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Telepon" :model-value="store.activeStudent.phone || '-'" dense /></div>
                  <div class="col-12"><q-input readonly stack-label type="textarea" rows="2" label="Alamat" :model-value="store.activeStudent.address || '-'" dense /></div>
                </div>
              </q-tab-panel>

              <!-- TAB PENDIDIKAN -->
              <q-tab-panel name="pendidikan" class="q-pa-none">
                <div class="row q-col-gutter-sm">
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Sekolah Asal (SMP/MTs)" :model-value="store.activeStudent.sekolahAsal || '-'" dense /></div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Diterima di Kelas" :model-value="store.activeStudent.diterimaKelas || '-'" dense /></div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Tanggal Diterima" :model-value="diterimaTanggalFormatted" dense /></div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Skor Kelengkapan" :model-value="`${store.activeStudent.skorKelengkapan}%`" dense /></div>
                </div>
              </q-tab-panel>

              <!-- TAB KELUARGA -->
              <q-tab-panel name="keluarga" class="q-pa-none">
                <div class="row q-col-gutter-sm">
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Nama Ayah" :model-value="store.activeStudent.namaAyah || '-'" dense /></div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Pekerjaan Ayah" :model-value="store.activeStudent.pekerjaanAyah || '-'" dense /></div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Nama Ibu" :model-value="store.activeStudent.namaIbu || '-'" dense /></div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Pekerjaan Ibu" :model-value="store.activeStudent.pekerjaanIbu || '-'" dense /></div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Nama Wali" :model-value="store.activeStudent.namaWali || '-'" dense /></div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Pekerjaan Wali" :model-value="store.activeStudent.pekerjaanWali || '-'" dense /></div>
                  <div class="col-12"><q-input readonly stack-label type="textarea" rows="2" label="Alamat Orang Tua" :model-value="store.activeStudent.alamatOrtu || '-'" dense /></div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Telepon Orang Tua" :model-value="store.activeStudent.teleponOrtu || '-'" dense /></div>
                  <div class="col-12"><q-input readonly stack-label type="textarea" rows="2" label="Alamat Wali" :model-value="store.activeStudent.alamatWali || '-'" dense /></div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Telepon Wali" :model-value="store.activeStudent.teleponWali || '-'" dense /></div>
                </div>
              </q-tab-panel>

              <!-- TAB PELENGKAP -->
              <q-tab-panel name="pelengkap" class="q-pa-none">
                <div class="row q-col-gutter-sm">
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Status dalam Keluarga" :model-value="store.activeStudent.statusDalamKel || '-'" dense /></div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Anak Ke-" :model-value="String(store.activeStudent.anakKe || '-')" dense /></div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Status Kelengkapan" :model-value="store.activeStudent.statusKelengkapan || '-'" dense /></div>
                  <div class="col-12 q-mt-sm">
                    <q-banner class="bg-blue-1 text-primary rounded-borders dense">
                      <q-icon name="info" class="q-mr-xs" />
                      Data pelengkap dapat dilengkapi melalui menu edit siswa.
                    </q-banner>
                  </div>
                </div>
              </q-tab-panel>

              <!-- TAB NILAI RAPOR -->
              <q-tab-panel name="nilai" class="q-pa-none">
                <div v-if="store.loadingDetail" class="text-center q-pa-md">
                  <q-spinner color="primary" size="2em" />
                </div>
                <q-markup-table v-else dense flat bordered>
                  <thead>
                    <tr>
                      <th class="text-center">No</th>
                      <th class="text-left">Mata Pelajaran</th>
                      <th class="text-center">Smt</th>
                      <th class="text-center">Harian</th>
                      <th class="text-center">UTS</th>
                      <th class="text-center">UAS</th>
                      <th class="text-center">Akhir</th>
                      <th class="text-center">Predikat</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="(row, idx) in raporGroups" :key="`${row.semester}-${row.namaMataPelajaran}`">
                      <td class="text-center">{{ idx + 1 }}</td>
                      <td>{{ row.namaMataPelajaran }}</td>
                      <td class="text-center">{{ row.semester }}</td>
                      <td class="text-center">{{ row.daily ?? '-' }}</td>
                      <td class="text-center">{{ row.midterm ?? '-' }}</td>
                      <td class="text-center">{{ row.final ?? '-' }}</td>
                      <td class="text-center"><strong>{{ row.avg }}</strong></td>
                      <td class="text-center">
                        <q-badge :color="predikatColor(row.predikat)">{{ row.predikat }}</q-badge>
                      </td>
                    </tr>
                    <tr v-if="raporGroups.length === 0">
                      <td colspan="8" class="text-center text-grey-6">
                        Belum ada data nilai untuk siswa ini pada periode aktif.
                      </td>
                    </tr>
                  </tbody>
                </q-markup-table>
              </q-tab-panel>
            </q-tab-panels>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useStudentBukuIndukStore } from '@/stores/kesiswaan/studentBukuIndukStore.js'

const $q = useQuasar()
const store = useStudentBukuIndukStore()

const alphabets = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'.split('')
const selectedLetter = ref('A')
const searchKlaper = ref('')
const activeTab = ref('pribadi')
const printing = ref(false)

onMounted(loadList)

async function loadList() {
  await store.fetchList()
}

const filteredList = computed(() => {
  return store.students.filter(siswa => {
    const name = (siswa.name || '').toUpperCase()
    const nis = siswa.nis || ''
    const matchLetter = name.startsWith(selectedLetter.value)
    const searchLower = searchKlaper.value.toLowerCase()
    const matchSearch = !searchKlaper.value ||
      name.toLowerCase().includes(searchLower) ||
      nis.includes(searchKlaper.value)
    return matchLetter && matchSearch
  })
})

const genderLabel = computed(() => {
  const g = store.activeStudent?.gender
  if (!g) return '-'
  if (g === 'MALE' || g === 'L') return 'Laki-laki'
  if (g === 'FEMALE' || g === 'P') return 'Perempuan'
  return g
})

const birthDateFormatted = computed(() => formatDate(store.activeStudent?.birthDate))
const diterimaTanggalFormatted = computed(() => formatDate(store.activeStudent?.diterimaTanggal))

const completenessColor = computed(() => {
  const skor = store.activeStudent?.skorKelengkapan || 0
  if (skor >= 81) return 'green'
  if (skor >= 51) return 'orange'
  return 'red'
})

const raporGroups = computed(() => {
  const scores = store.activeStudent?.academicScores || []
  const map = new Map()
  for (const sc of scores) {
    const key = `${sc.semester}-${sc.mataPelajaranId}`
    if (!map.has(key)) {
      map.set(key, {
        semester: sc.semester,
        namaMataPelajaran: sc.namaMataPelajaran,
        daily: null, midterm: null, final: null,
      })
    }
    const item = map.get(key)
    if (sc.assessmentType === 'DAILY') item.daily = sc.score
    if (sc.assessmentType === 'MIDTERM') item.midterm = sc.score
    if (sc.assessmentType === 'FINAL') item.final = sc.score
  }
  return Array.from(map.values()).map(item => {
    const vals = [item.daily, item.midterm, item.final].filter(v => v !== null && v !== undefined)
    const avg = vals.length > 0 ? vals.reduce((a, b) => a + b, 0) / vals.length : 0
    let predikat = '-'
    if (avg >= 90) predikat = 'A'
    else if (avg >= 80) predikat = 'B'
    else if (avg >= 70) predikat = 'C'
    else if (avg > 0) predikat = 'D'
    return { ...item, avg: Math.round(avg), predikat }
  })
})

function formatDate(input) {
  if (!input) return '-'
  const d = new Date(input)
  if (isNaN(d.getTime())) return '-'
  return d.toLocaleDateString('id-ID', { day: '2-digit', month: 'long', year: 'numeric' })
}

function selectLetter(letter) {
  selectedLetter.value = letter
  searchKlaper.value = ''
}

function selectStudent(siswa) {
  activeTab.value = 'pribadi'
  store.fetchDetail(siswa.id)
}

function predikatColor(predikat) {
  if (predikat === 'A') return 'green'
  if (predikat === 'B') return 'blue'
  if (predikat === 'C') return 'orange'
  return 'grey'
}

async function printBukuInduk() {
  if (!store.activeStudent) return
  printing.value = true
  try {
    const result = await store.downloadPDF(store.activeStudent.id)
    if (!result.success) {
      $q.notify({ color: 'negative', message: result.error || 'Gagal mengunduh PDF', icon: 'error' })
      return
    }
    const url = URL.createObjectURL(result.data)
    const link = document.createElement('a')
    link.href = url
    link.download = `buku-induk-${store.activeStudent.nis}.pdf`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
    $q.notify({ color: 'positive', message: 'PDF Buku Induk berhasil diunduh', icon: 'check_circle' })
  } finally {
    printing.value = false
  }
}
</script>

<style scoped>
.q-item--active {
  border-left: 4px solid var(--q-primary);
}
</style>
