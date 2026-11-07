<template>
  <q-page class="q-pa-md bg-grey-2">
    <div class="row justify-between items-center q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Data GTK</h1>
        <p class="text-caption text-grey-7 q-mb-none">Guru & Tenaga Kependidikan (DAPODIK)</p>
      </div>
      <q-btn color="primary" icon="add" label="Tambah Guru" @click="openDialog()" />
    </div>

    <q-card class="shadow-1">
      <q-table :rows="guruStore.list" :columns="columns" row-key="id" :loading="guruStore.isLoading" flat bordered
        dense>
        <template v-slot:body-cell-foto="props">
          <q-td :props="props" class="text-center">
            <q-avatar size="40px">
              <img :src="props.row.fotoUrl || 'https://cdn.quasar.dev/img/boy-avatar.png'"
                @error="(e) => e.target.src = 'https://cdn.quasar.dev/img/boy-avatar.png'">
            </q-avatar>
          </q-td>
        </template>
        <template v-slot:body-cell-gender="props">
          <q-td :props="props" class="text-center">
            <q-badge :color="props.value === 'L' ? 'blue-6' : 'pink-6'">{{ props.value === 'L' ? 'L' : 'P' }}</q-badge>
          </q-td>
        </template>
        <template v-slot:body-cell-status="props">
          <q-td :props="props" class="text-center">
            <q-badge :color="props.row.isActiveInPeriod ? 'green-6' : 'grey-6'">{{ props.row.isActiveInPeriod ? 'Aktif'
              : 'Tidak Aktif' }}</q-badge>
          </q-td>
        </template>
        <template v-slot:body-cell-aksi="props">
          <q-td :props="props" class="text-center">
            <q-btn flat round dense icon="edit" color="info" @click="openDialog(props.row)" class="q-mr-sm" />
            <q-btn flat round dense icon="delete" color="negative" @click="confirmDelete(props.row)" />
          </q-td>
        </template>
      </q-table>
    </q-card>

    <!-- Dialog Form Tambah/Edit -->
    <q-dialog v-model="dialogOpen" persistent maximized transition-show="slide-up" transition-hide="slide-down">
      <q-card class="column no-wrap">
        <q-card-section class="row items-center q-pb-none bg-primary text-white">
          <div class="text-h6">{{ isEdit ? 'Edit Data Guru' : 'Tambah Data Guru' }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="col scroll q-pa-none">
          <q-form @submit="onSubmit" class="q-gutter-none">
            <q-tabs v-model="tab" class="text-primary" align="left">
              <q-tab name="profil" icon="person" label="Profil & Kepegawaian" />
              <q-tab name="pendidikan" icon="school" label="Riwayat Pendidikan" />
              <q-tab name="sertifikasi" icon="workspace_premium" label="Sertifikasi" />
            </q-tabs>
            <q-separator />

            <q-tab-panels v-model="tab" animated class="bg-white">

              <!-- TAB 1: PROFIL -->
              <q-tab-panel name="profil">
                <div class="row justify-center q-mb-md">
                  <q-avatar size="100px" class="shadow-2">
                    <img :src="previewFoto || form.fotoUrl || 'https://cdn.quasar.dev/img/boy-avatar.png'">
                    <q-badge floating color="primary" rounded class="cursor-pointer"
                      @click="$refs.fileInput.pickFiles()">
                      <q-icon name="camera_alt" size="18px" />
                    </q-badge>
                  </q-avatar>
                  <q-file ref="fileInput" v-model="fileFoto" @update:model-value="onFilePicked" class="hidden"
                    accept="image/*" />
                </div>

                <div class="text-subtitle1 text-weight-bold text-primary q-mt-md border-bottom">A. Identitas Pribadi
                </div>
                <div class="row q-col-gutter-sm">
                  <div class="col-12 col-md-4"><q-input v-model="form.nik" label="NIK" filled mask="################" />
                  </div>
                  <div class="col-12 col-md-4"><q-input v-model="form.name" label="Nama Lengkap" filled
                      :rules="[val => !!val || 'Wajib diisi']" /></div>
                  <div class="col-12 col-md-2"><q-input v-model="form.gelarDepan" label="Gelar Depan" filled /></div>
                  <div class="col-12 col-md-2"><q-input v-model="form.gelarBelakang" label="Gelar Belakang" filled />
                  </div>
                </div>
                <div class="row q-col-gutter-sm">
                  <div class="col-12 col-md-3"><q-select v-model="form.gender"
                      :options="[{ label: 'Laki-laki', value: 'L' }, { label: 'Perempuan', value: 'P' }]" emit-value
                      map-options label="JK" filled /></div>
                  <div class="col-12 col-md-3"><q-select v-model="form.agamaId" :options="agamaOptions" emit-value
                      map-options label="Agama" filled /></div>
                  <div class="col-12 col-md-3"><q-select v-model="form.statusPerkawinan" :options="kawinOptions"
                      emit-value map-options label="Status Perkawinan" filled /></div>
                  <div class="col-12 col-md-3"><q-input v-model="form.namaIbuKandung" label="Nama Ibu Kandung" filled />
                  </div>
                </div>
                <div class="row q-col-gutter-sm">
                  <div class="col-12 col-md-6"><q-input v-model="form.birthPlace" label="Tempat Lahir" filled /></div>
                  <div class="col-12 col-md-6"><q-input v-model="form.birthDate" label="Tanggal Lahir" filled
                      mask="####-##-##" hint="YYYY-MM-DD" /></div>
                </div>
                <div class="row q-col-gutter-sm">
                  <div class="col-12 col-md-6"><q-input v-model="form.noHp" label="No. HP" filled
                      mask="##############" />
                  </div>
                  <div class="col-12 col-md-6"><q-input v-model="form.email" label="Email" filled type="email" /></div>
                </div>
                <q-input v-model="form.address" label="Alamat Lengkap" type="textarea" filled />

                <div class="text-subtitle1 text-weight-bold text-primary q-mt-lg border-bottom">B. Data Kepegawaian
                </div>
                <div class="row q-col-gutter-sm">
                  <div class="col-12 col-md-3"><q-select v-model="form.jenisPtkId" :options="jenisPtkOptions" emit-value
                      map-options label="Jenis PTK" filled /></div>
                  <div class="col-12 col-md-3"><q-select v-model="form.statusKeaktifan" :options="statusAktifOptions"
                      emit-value map-options label="Status Aktif" filled /></div>
                  <div class="col-12 col-md-3"><q-select v-model="form.statusKepegawaian"
                      :options="statusKepegawaianOptions" emit-value map-options label="Status Kepegawaian" filled />
                  </div>
                  <div class="col-12 col-md-3"><q-input v-model="form.nip" label="NIP" filled
                      mask="##################" />
                  </div>
                </div>
                <div class="row q-col-gutter-sm">
                  <div class="col-12 col-md-6"><q-input v-model="form.nuptk" label="NUPTK" filled
                      mask="################" />
                  </div>
                </div>
              </q-tab-panel>

              <!-- TAB 2: RIWAYAT PENDIDIKAN -->
              <q-tab-panel name="pendidikan">
                <q-banner v-if="!isEdit" class="bg-orange-2 text-orange-9 q-mb-md" rounded>
                  <template v-slot:avatar><q-icon name="info" color="orange-9" /></template>
                  Silakan isi dan <strong>Simpan Profil</strong> terlebih dahulu.
                </q-banner>

                <div class="row items-center q-mb-md">
                  <div class="text-subtitle1 text-weight-bold text-primary">Riwayat Pendidikan Formal</div>
                  <q-space />
                  <q-btn color="secondary" icon="add" label="Tambah" size="sm" @click="openEduDialog"
                    :disable="!isEdit" />
                </div>
                <q-table :rows="educations" :columns="eduColumns" row-key="id" flat bordered dense>
                  <template v-slot:body-cell-dokumen="props">
                    <q-td :props="props" class="text-center">
                      <q-btn v-if="props.row.dokumenUrl" flat round dense icon="visibility" color="primary" type="a"
                        :href="props.row.dokumenUrl" target="_blank" />
                      <span v-else class="text-grey-5">-</span>
                    </q-td>
                  </template>
                  <template v-slot:body-cell-aksi="props">
                    <q-td :props="props" class="text-center">
                      <q-btn flat round dense icon="delete" color="negative" @click="deleteEducation(props.row.id)" />
                    </q-td>
                  </template>
                </q-table>
              </q-tab-panel>

              <!-- TAB 3: SERTIFIKASI -->
              <q-tab-panel name="sertifikasi">
                <q-banner v-if="!isEdit" class="bg-orange-2 text-orange-9 q-mb-md" rounded>
                  <template v-slot:avatar><q-icon name="info" color="orange-9" /></template>
                  Silakan isi dan <strong>Simpan Profil</strong> terlebih dahulu.
                </q-banner>

                <div class="row items-center q-mb-md">
                  <div class="text-subtitle1 text-weight-bold text-primary">Sertifikasi</div>
                  <q-space />
                  <q-btn color="secondary" icon="add" label="Tambah" size="sm" @click="openCertDialog"
                    :disable="!isEdit" />
                </div>
                <q-table :rows="certifications" :columns="certColumns" row-key="id" flat bordered dense>
                  <template v-slot:body-cell-dokumen="props">
                    <q-td :props="props" class="text-center">
                      <q-btn v-if="props.row.dokumenUrl" flat round dense icon="visibility" color="primary" type="a"
                        :href="props.row.dokumenUrl" target="_blank" />
                      <span v-else class="text-grey-5">-</span>
                    </q-td>
                  </template>
                  <template v-slot:body-cell-aksi="props">
                    <q-td :props="props" class="text-center">
                      <q-btn flat round dense icon="delete" color="negative"
                        @click="deleteCertification(props.row.id)" />
                    </q-td>
                  </template>
                </q-table>
              </q-tab-panel>

            </q-tab-panels>

            <div class="row justify-end q-pa-md bg-grey-1">
              <q-btn label="Batal" color="grey" flat v-close-popup class="q-mr-sm" />
              <q-btn v-if="tab === 'profil'" label="Simpan Profil" color="primary" type="submit"
                :loading="guruStore.isLoading" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Dialog Form Pendidikan -->
    <q-dialog v-model="eduDialogOpen" persistent>
      <q-card style="min-width: 500px">
        <q-card-section class="bg-secondary text-white">
          <div class="text-h6">Tambah Riwayat Pendidikan</div>
        </q-card-section>
        <q-card-section class="q-gutter-y-md">
          <q-input v-model="eduForm.jenjang" label="Jenjang (SMA/D3/S1/S2)" filled dense required />
          <q-input v-model="eduForm.jurusan" label="Jurusan" filled dense />
          <q-input v-model="eduForm.namaSekolah" label="Nama Sekolah/Universitas" filled dense />
          <div class="row q-col-gutter-sm">
            <div class="col-6"><q-input v-model.number="eduForm.tahunLulus" label="Tahun Lulus" type="number" filled
                dense /></div>
            <div class="col-6"><q-input v-model="eduForm.gelar" label="Gelar Didapat" filled dense /></div>
          </div>
          <q-file v-model="eduFile" label="Upload Ijazah (PDF/Gambar)" filled dense accept=".pdf,image/*" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="secondary" label="Simpan" @click="submitEducation" :loading="isSubmittingEdu" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog Form Sertifikasi -->
    <q-dialog v-model="certDialogOpen" persistent>
      <q-card style="min-width: 500px">
        <q-card-section class="bg-secondary text-white">
          <div class="text-h6">Tambah Sertifikasi</div>
        </q-card-section>
        <q-card-section class="q-gutter-y-md">
          <q-input v-model="certForm.jenisSertifikasi" label="Jenis Sertifikasi" filled dense required />
          <q-input v-model="certForm.nomorSertifikat" label="Nomor Sertifikat" filled dense />
          <q-input v-model.number="certForm.tahunSertifikasi" label="Tahun Sertifikasi" type="number" filled dense />
          <q-file v-model="certFile" label="Upload Sertifikat (PDF/Gambar)" filled dense accept=".pdf,image/*" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="secondary" label="Simpan" @click="submitCertification" :loading="isSubmittingCert" />
        </q-card-actions>
      </q-card>
    </q-dialog>

  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useGuruStore } from '@/stores/kepegawaian/guru'
import { api } from '@/boot/axios'

const $q = useQuasar()
const guruStore = useGuruStore()

const dialogOpen = ref(false)
const isEdit = ref(false)
const tab = ref('profil')
const fileFoto = ref(null)
const previewFoto = ref(null)

// State Riwayat
const educations = ref([])
const certifications = ref([])

// State Dialog Pendidikan
const eduDialogOpen = ref(false)
const isSubmittingEdu = ref(false)
const eduForm = ref({ jenjang: 'S1', jurusan: '', namaSekolah: '', tahunLulus: new Date().getFullYear(), gelar: '', dokumenUrl: '' })
const eduFile = ref(null)

// State Dialog Sertifikasi
const certDialogOpen = ref(false)
const isSubmittingCert = ref(false)
const certForm = ref({ jenisSertifikasi: '', nomorSertifikat: '', tahunSertifikasi: new Date().getFullYear(), dokumenUrl: '' })
const certFile = ref(null)

const agamaOptions = [
  { label: 'Islam', value: 1 }, { label: 'Kristen', value: 2 }, { label: 'Katolik', value: 3 },
  { label: 'Hindu', value: 4 }, { label: 'Buddha', value: 5 }, { label: 'Konghucu', value: 6 }
]
const kawinOptions = [{ label: 'Belum Menikah', value: 1 }, { label: 'Menikah', value: 2 }, { label: 'Cerai Hidup', value: 3 }, { label: 'Cerai Mati', value: 4 }]
const jenisPtkOptions = [{ label: 'Guru Mata Pelajaran', value: 1 }, { label: 'Guru Kelas', value: 2 }, { label: 'Kepala Sekolah', value: 3 }, { label: 'Tenaga Administrasi', value: 4 }]
const statusAktifOptions = [{ label: 'Aktif', value: 1 }, { label: 'Cuti', value: 2 }, { label: 'Pensiun', value: 3 }, { label: 'Non-Aktif', value: 4 }]
const statusKepegawaianOptions = [{ label: 'PNS', value: 1 }, { label: 'PPPK', value: 2 }, { label: 'GTY/PTY', value: 3 }, { label: 'Honorer', value: 4 }]

const form = ref({ id: '', name: '', gelarDepan: '', gelarBelakang: '', nip: '', nuptk: '', gender: 'L', birthPlace: '', birthDate: '', address: '', jenisPtkId: 1, statusKeaktifan: 1, nik: '', agamaId: 1, statusKepegawaian: 1, noHp: '', email: '', statusPerkawinan: 1, namaIbuKandung: '', fotoUrl: '' })

const columns = [
  { name: 'foto', label: 'Foto', align: 'center' },
  { name: 'name', label: 'Nama Lengkap & Gelar', field: 'name', align: 'left', sortable: true },
  { name: 'nik', label: 'NIK', field: 'nik', align: 'left' },
  { name: 'nip', label: 'NIP', field: 'nip', align: 'left' },
  { name: 'gender', label: 'JK', field: 'gender', align: 'center' },
  { name: 'status', label: 'Status Jadwal', align: 'center' },
  { name: 'aksi', label: 'Aksi', align: 'center' }
]

const eduColumns = [
  { name: 'jenjang', label: 'Jenjang', field: 'jenjang', align: 'left' },
  { name: 'jurusan', label: 'Jurusan', field: 'jurusan', align: 'left' },
  { name: 'namaSekolah', label: 'Nama Sekolah', field: 'namaSekolah', align: 'left' },
  { name: 'tahunLulus', label: 'Tahun Lulus', field: 'tahunLulus', align: 'center' },
  { name: 'gelar', label: 'Gelar', field: 'gelar', align: 'left' },
  { name: 'dokumen', label: 'Ijazah', align: 'center' },
  { name: 'aksi', label: 'Aksi', align: 'center' }
]

const certColumns = [
  { name: 'jenisSertifikasi', label: 'Jenis Sertifikasi', field: 'jenisSertifikasi', align: 'left' },
  { name: 'nomorSertifikat', label: 'Nomor Sertifikat', field: 'nomorSertifikat', align: 'left' },
  { name: 'tahunSertifikasi', label: 'Tahun', field: 'tahunSertifikasi', align: 'center' },
  { name: 'dokumen', label: 'Dokumen', align: 'center' },
  { name: 'aksi', label: 'Aksi', align: 'center' }
]

onMounted(() => { guruStore.fetchTeachers() })

const onFilePicked = (file) => { if (file) previewFoto.value = URL.createObjectURL(file) }

const openDialog = async (row = null) => {
  fileFoto.value = null
  previewFoto.value = null
  tab.value = 'profil'

  if (row) {
    isEdit.value = true
    const tgl = row.birthDate ? new Date(row.birthDate).toISOString().split('T')[0] : ''
    let namaBersih = row.name
    if (row.gelarBelakang) namaBersih = namaBersih.replace(`, ${row.gelarBelakang}`, '')
    if (row.gelarDepan) namaBersih = namaBersih.replace(`${row.gelarDepan} `, '')

    form.value = { ...row, name: namaBersih, birthDate: tgl }
    await fetchEducations(row.id)
    await fetchCertifications(row.id)
  } else {
    isEdit.value = false
    educations.value = []
    certifications.value = []
    form.value = { id: '', name: '', gelarDepan: '', gelarBelakang: '', nip: '', nuptk: '', gender: 'L', birthPlace: '', birthDate: '', address: '', jenisPtkId: 1, statusKeaktifan: 1, nik: '', agamaId: 1, statusKepegawaian: 1, noHp: '', email: '', statusPerkawinan: 1, namaIbuKandung: '', fotoUrl: '' }
  }
  dialogOpen.value = true
}

const fetchEducations = async (ptkID) => {
  try {
    const res = await api.get(`/teachers/history/${ptkID}/educations`)
    educations.value = res.data || []
  } catch (err) { console.error(err) }
}

const fetchCertifications = async (ptkID) => {
  try {
    const res = await api.get(`/teachers/history/${ptkID}/certifications`)
    certifications.value = res.data || []
  } catch (err) { console.error(err) }
}

// --- LOGIKA PENDIDIKAN ---
const openEduDialog = () => {
  eduForm.value = { jenjang: 'S1', jurusan: '', namaSekolah: '', tahunLulus: new Date().getFullYear(), gelar: '', dokumenUrl: '' }
  eduFile.value = null
  eduDialogOpen.value = true
}

const submitEducation = async () => {
  isSubmittingEdu.value = true
  try {
    if (eduFile.value) {
      const formData = new FormData()
      formData.append('document', eduFile.value)
      const uploadRes = await api.post('/teachers/history/upload', formData, { headers: { 'Content-Type': 'multipart/form-data' } })
      eduForm.value.dokumenUrl = uploadRes.data.url
    }
    await api.post(`/teachers/history/${form.value.id}/educations`, eduForm.value)
    await fetchEducations(form.value.id)
    eduDialogOpen.value = false
    $q.notify({ color: 'positive', message: 'Pendidikan ditambahkan' })
  } catch (e) {
    $q.notify({ color: 'negative', message: `Gagal menambah pendidikan ${e}` })
  } finally {
    isSubmittingEdu.value = false
  }
}

const deleteEducation = async (id) => {
  try {
    await api.delete(`/teachers/history/${form.value.id}/educations/${id}`)
    await fetchEducations(form.value.id)
  } catch (e) { $q.notify({ color: 'negative', message: `Gagal menghapus ${e}` }) }
}

// --- LOGIKA SERTIFIKASI ---
const openCertDialog = () => {
  certForm.value = { jenisSertifikasi: '', nomorSertifikat: '', tahunSertifikasi: new Date().getFullYear(), dokumenUrl: '' }
  certFile.value = null
  certDialogOpen.value = true
}

const submitCertification = async () => {
  isSubmittingCert.value = true
  try {
    if (certFile.value) {
      const formData = new FormData()
      formData.append('document', certFile.value)
      const uploadRes = await api.post('/teachers/history/upload', formData, { headers: { 'Content-Type': 'multipart/form-data' } })
      certForm.value.dokumenUrl = uploadRes.data.url
    }
    await api.post(`/teachers/history/${form.value.id}/certifications`, certForm.value)
    await fetchCertifications(form.value.id)
    certDialogOpen.value = false
    $q.notify({ color: 'positive', message: 'Sertifikasi ditambahkan' })
  } catch (e) {
    $q.notify({ color: 'negative', message: `Gagal menambah sertifikasi ${e}` })
  } finally {
    isSubmittingCert.value = false
  }
}

const deleteCertification = async (id) => {
  try {
    await api.delete(`/teachers/history/${form.value.id}/certifications/${id}`)
    await fetchCertifications(form.value.id)
  } catch (e) { $q.notify({ color: 'negative', message: `Gagal menghapus ${e}` }) }
}

const onSubmit = async () => {
  try {
    const payload = { ...form.value }
    let savedId = form.value.id

    if (isEdit.value) {
      await guruStore.updateTeacher(form.value.id, payload)
      $q.notify({ color: 'positive', message: 'Data Guru diperbarui' })
    } else {
      const res = await guruStore.createTeacher(payload)
      savedId = res.id
      $q.notify({ color: 'positive', message: 'Guru baru ditambahkan' })
    }

    if (fileFoto.value && savedId) {
      const formData = new FormData()
      formData.append('photo', fileFoto.value)
      await api.post(`/teachers/photo/${savedId}`, formData, { headers: { 'Content-Type': 'multipart/form-data' } })
    }

    dialogOpen.value = false
    await guruStore.fetchTeachers()
  } catch (err) {
    $q.notify({ color: 'negative', message: err.response?.data?.error || 'Terjadi kesalahan' })
  }
}

const confirmDelete = (row) => {
  $q.dialog({ title: 'Konfirmasi Hapus', message: `Hapus data guru ${row.name}?`, cancel: true, persistent: true })
    .onOk(async () => {
      try {
        await guruStore.deleteTeacher(row.id)
        $q.notify({ color: 'positive', message: 'Data guru dihapus' })
      } catch (err) { $q.notify({ color: 'negative', message: `Gagal: ${err}` }) }
    })
}
</script>

<style scoped>
.border-bottom {
  border-bottom: 2px solid #e0e0e0;
  padding-bottom: 8px;
  margin-bottom: 12px;
}

.hidden {
  display: none;
}
</style>
