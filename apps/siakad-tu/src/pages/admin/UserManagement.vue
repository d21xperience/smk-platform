<template>
  <q-page class="q-pa-md">
    <div class="text-h5 q-mb-md">Manajemen Pengguna</div>

    <q-tabs v-model="tab" dense class="q-mb-md" align="left">
      <q-tab name="siswa" label="Siswa" />
      <q-tab name="guru" label="Guru" />
      <q-tab name="admin" label="Admin" />
    </q-tabs>

    <q-tab-panels v-model="tab" animated>
      <!-- Tab Siswa -->
      <q-tab-panel name="siswa" class="q-pa-none">
        <UserTable role="siswa" :data="siswaList" :columns="siswaColumns" @add="openForm('siswa', null)"
          @edit="openForm('siswa', $event)" @delete="confirmDelete('siswa', $event)" />
      </q-tab-panel>

      <!-- Tab Guru -->
      <q-tab-panel name="guru" class="q-pa-none">
        <UserTable role="guru" :data="guruList" :columns="guruColumns" @add="openForm('guru', null)"
          @edit="openForm('guru', $event)" @delete="confirmDelete('guru', $event)" />
      </q-tab-panel>

      <!-- Tab Admin -->
      <q-tab-panel name="admin" class="q-pa-none">
        <UserTable role="admin" :data="adminList" :columns="adminColumns" @add="openForm('admin', null)"
          @edit="openForm('admin', $event)" @delete="confirmDelete('admin', $event)" />
      </q-tab-panel>
    </q-tab-panels>

    <!-- Dialog Form Tambah/Edit -->
    <q-dialog v-model="formDialog" persistent>
      <q-card style="min-width: 500px;">
        <q-card-section>
          <div class="text-h6">{{ formTitle }}</div>
        </q-card-section>
        <q-card-section>
          <q-form ref="userForm" @submit="submitForm" class="q-gutter-md">
            <q-input v-model="formData.nama" label="Nama Lengkap" lazy-rules
              :rules="[val => !!val || 'Nama harus diisi']" outlined />
            <q-input v-model="formData.email" label="Email" type="email" lazy-rules
              :rules="[val => !!val || 'Email harus diisi', val => /.+@.+\..+/.test(val) || 'Email tidak valid']"
              outlined />
            <q-input v-model="formData.username" label="Username" lazy-rules
              :rules="[val => !!val || 'Username harus diisi']" outlined />
            <q-input v-if="!isEdit" v-model="formData.password" label="Password" type="password" lazy-rules
              :rules="[val => !!val || 'Password harus diisi']" outlined />
            <q-input v-if="!isEdit" v-model="formData.confirmPassword" label="Konfirmasi Password" type="password"
              lazy-rules :rules="[val => val === formData.password || 'Password tidak cocok']" outlined />
            <!-- Field khusus untuk siswa -->
            <q-select v-if="currentRole === 'siswa'" v-model="formData.kelas" :options="kelasOptions" label="Kelas"
              outlined :rules="[val => !!val || 'Kelas harus dipilih']" />
            <q-input v-if="currentRole === 'siswa'" v-model="formData.nis" label="NIS" outlined />
            <!-- Field khusus untuk guru -->
            <q-input v-if="currentRole === 'guru'" v-model="formData.nip" label="NIP" outlined />
            <q-select v-if="currentRole === 'guru'" v-model="formData.mataPelajaran" :options="mapelOptions"
              label="Mata Pelajaran" multiple outlined use-chips />
            <!-- Field khusus untuk admin -->
            <q-input v-if="currentRole === 'admin'" v-model="formData.level" label="Level Admin (super/admin)"
              outlined />
          </q-form>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" color="negative" v-close-popup />
          <q-btn flat label="Simpan" color="primary" type="submit" @click="submitForm" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog Konfirmasi Hapus -->
    <q-dialog v-model="deleteDialog" persistent>
      <q-card>
        <q-card-section>
          <div class="text-h6">Konfirmasi Hapus</div>
        </q-card-section>
        <q-card-section>
          Apakah Anda yakin ingin menghapus {{ deleteTarget ? deleteTarget.nama : '' }}?
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" color="negative" v-close-popup />
          <q-btn flat label="Hapus" color="primary" @click="deleteUser" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import UserTable from 'components/UserTable.vue' // komponen tabel reusable
import { useQuasar } from 'quasar'

const $q = useQuasar()

// Data dummy
const siswaList = ref([
  { id: 1, nama: 'Ahmad Faizal', email: 'ahmad@example.com', username: 'ahmad', kelas: '12 IPA 1', nis: '12345' },
  { id: 2, nama: 'Siti Nurhaliza', email: 'siti@example.com', username: 'siti', kelas: '11 IPS 2', nis: '12346' }
])
const guruList = ref([
  { id: 1, nama: 'Budi Santoso', email: 'budi@example.com', username: 'budi', nip: '98765', mataPelajaran: ['Matematika', 'Fisika'] },
  { id: 2, nama: 'Dewi Lestari', email: 'dewi@example.com', username: 'dewi', nip: '98766', mataPelajaran: ['Bahasa Inggris'] }
])
const adminList = ref([
  { id: 1, nama: 'Admin Utama', email: 'admin@example.com', username: 'admin', level: 'super' }
])

const kelasOptions = ['10 IPA 1', '10 IPA 2', '10 IPS 1', '11 IPA 1', '11 IPA 2', '11 IPS 1', '12 IPA 1', '12 IPA 2', '12 IPS 1']
const mapelOptions = ['Matematika', 'Fisika', 'Kimia', 'Biologi', 'Bahasa Indonesia', 'Bahasa Inggris', 'Sejarah', 'Geografi', 'Ekonomi']

// State untuk tab dan form
const tab = ref('siswa')
const currentRole = ref('siswa')
const formDialog = ref(false)
const isEdit = ref(false)
const editId = ref(null)
const formData = ref({
  nama: '',
  email: '',
  username: '',
  password: '',
  confirmPassword: '',
  kelas: null,
  nis: '',
  nip: '',
  mataPelajaran: [],
  level: ''
})
const deleteDialog = ref(false)
const deleteTarget = ref(null)
const deleteRole = ref('')

// Definisi kolom untuk tabel (bisa juga di dalam UserTable, tapi lebih fleksibel di sini)
const siswaColumns = [
  { name: 'nama', label: 'Nama', field: 'nama', align: 'left' },
  { name: 'email', label: 'Email', field: 'email', align: 'left' },
  { name: 'username', label: 'Username', field: 'username', align: 'left' },
  { name: 'kelas', label: 'Kelas', field: 'kelas', align: 'left' },
  { name: 'nis', label: 'NIS', field: 'nis', align: 'left' },
  { name: 'actions', label: 'Aksi', field: 'actions', align: 'center' }
]
const guruColumns = [
  { name: 'nama', label: 'Nama', field: 'nama', align: 'left' },
  { name: 'email', label: 'Email', field: 'email', align: 'left' },
  { name: 'username', label: 'Username', field: 'username', align: 'left' },
  { name: 'nip', label: 'NIP', field: 'nip', align: 'left' },
  { name: 'mataPelajaran', label: 'Mata Pelajaran', field: 'mataPelajaran', align: 'left' },
  { name: 'actions', label: 'Aksi', field: 'actions', align: 'center' }
]
const adminColumns = [
  { name: 'nama', label: 'Nama', field: 'nama', align: 'left' },
  { name: 'email', label: 'Email', field: 'email', align: 'left' },
  { name: 'username', label: 'Username', field: 'username', align: 'left' },
  { name: 'level', label: 'Level', field: 'level', align: 'left' },
  { name: 'actions', label: 'Aksi', field: 'actions', align: 'center' }
]

function openForm(role, user = null) {
  currentRole.value = role
  if (user) {
    isEdit.value = true
    editId.value = user.id
    // Salin data user ke formData (tanpa password)
    formData.value = {
      nama: user.nama || '',
      email: user.email || '',
      username: user.username || '',
      password: '',
      confirmPassword: '',
      kelas: user.kelas || null,
      nis: user.nis || '',
      nip: user.nip || '',
      mataPelajaran: user.mataPelajaran || [],
      level: user.level || ''
    }
  } else {
    isEdit.value = false
    editId.value = null
    formData.value = {
      nama: '',
      email: '',
      username: '',
      password: '',
      confirmPassword: '',
      kelas: null,
      nis: '',
      nip: '',
      mataPelajaran: [],
      level: ''
    }
  }
  formDialog.value = true
}

function submitForm() {
  // Validasi sederhana, bisa pakai ref form
  // Di sini kita asumsikan validasi sudah lewat q-form rules
  // Simpan data sesuai role
  const newUser = {
    id: isEdit.value ? editId.value : Date.now(),
    nama: formData.value.nama,
    email: formData.value.email,
    username: formData.value.username,
    ...(currentRole.value === 'siswa' && { kelas: formData.value.kelas, nis: formData.value.nis }),
    ...(currentRole.value === 'guru' && { nip: formData.value.nip, mataPelajaran: formData.value.mataPelajaran }),
    ...(currentRole.value === 'admin' && { level: formData.value.level })
  }
  if (!isEdit.value) {
    // Tambah
    if (currentRole.value === 'siswa') siswaList.value.push(newUser)
    else if (currentRole.value === 'guru') guruList.value.push(newUser)
    else adminList.value.push(newUser)
    $q.notify({ type: 'positive', message: 'Pengguna berhasil ditambahkan' })
  } else {
    // Edit
    const index = currentRole.value === 'siswa' ? siswaList.value.findIndex(u => u.id === editId.value) :
      (currentRole.value === 'guru' ? guruList.value.findIndex(u => u.id === editId.value) :
        adminList.value.findIndex(u => u.id === editId.value))
    if (index !== -1) {
      if (currentRole.value === 'siswa') siswaList.value[index] = { ...siswaList.value[index], ...newUser }
      else if (currentRole.value === 'guru') guruList.value[index] = { ...guruList.value[index], ...newUser }
      else adminList.value[index] = { ...adminList.value[index], ...newUser }
      $q.notify({ type: 'positive', message: 'Pengguna berhasil diupdate' })
    }
  }
  formDialog.value = false
}

function confirmDelete(role, user) {
  deleteRole.value = role
  deleteTarget.value = user
  deleteDialog.value = true
}

function deleteUser() {
  if (deleteRole.value === 'siswa') {
    const index = siswaList.value.findIndex(u => u.id === deleteTarget.value.id)
    if (index !== -1) siswaList.value.splice(index, 1)
  } else if (deleteRole.value === 'guru') {
    const index = guruList.value.findIndex(u => u.id === deleteTarget.value.id)
    if (index !== -1) guruList.value.splice(index, 1)
  } else if (deleteRole.value === 'admin') {
    const index = adminList.value.findIndex(u => u.id === deleteTarget.value.id)
    if (index !== -1) adminList.value.splice(index, 1)
  }
  $q.notify({ type: 'positive', message: 'Pengguna berhasil dihapus' })
  deleteDialog.value = false
  deleteTarget.value = null
}
</script>
