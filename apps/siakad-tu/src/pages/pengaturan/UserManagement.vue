<template>
  <q-page class="q-pa-md">
    <div class="row justify-between q-mb-md">
      <div class="text-h5">Manajemen Pengguna & Role</div>
      <q-btn color="primary" icon="add" label="Tambah Pengguna" @click="openForm" />
    </div>

    <q-table :rows="userStore.list" :columns="columns" row-key="id" flat bordered dense>
      <template v-slot:body-cell-role="props">
        <q-td :props="props">
          <q-badge :color="getRoleColor(props.row.role)" :label="props.row.role" />
        </q-td>
      </template>
      <template v-slot:body-cell-is_active="props">
        <q-td :props="props">
          <q-badge :color="props.row.is_active ? 'positive' : 'negative'"
            :label="props.row.is_active ? 'Aktif' : 'Nonaktif'" />
        </q-td>
      </template>
      <template v-slot:body-cell-aksi="props">
        <q-td :props="props">
          <q-btn flat round dense icon="edit" color="info" @click="openForm(props.row)" class="q-mr-sm" />
          <q-btn flat round dense icon="lock_reset" color="warning" @click="resetPassword(props.row.id)"
            class="q-mr-sm">
            <q-tooltip>Reset Password</q-tooltip>
          </q-btn>
          <q-btn flat round dense icon="delete" color="negative" @click="confirmDelete(props.row.id)" />
        </q-td>
      </template>
    </q-table>

    <q-dialog v-model="formDialog" persistent>
      <q-card style="min-width: 500px">
        <q-card-section>
          <div class="text-h6">{{ editMode ? 'Edit Pengguna' : 'Tambah Pengguna' }}</div>
        </q-card-section>
        <q-card-section>
          <q-form @submit="submitForm" class="q-gutter-md">
            <q-input v-model="form.username" label="Username" outlined :rules="[v => !!v]" />
            <q-input v-model="form.nama" label="Nama Lengkap" outlined :rules="[v => !!v]" />
            <q-input v-model="form.email" label="Email" type="email" outlined />
            <q-input v-if="!editMode" v-model="form.password" label="Password" type="password" outlined
              :rules="[v => !!v]" />
            <q-select v-model="form.role" :options="roleOptions" label="Role" outlined />
            <q-select v-model="form.is_active" :options="[{ label: 'Aktif', value: true }, { label: 'Nonaktif', value: false }]"
              label="Status" outlined />
            <div class="q-gutter-sm">
              <q-btn label="Simpan" type="submit" color="primary" />
              <q-btn label="Batal" flat color="negative" v-close-popup @click="resetForm" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useUserStore } from 'stores/user'

const $q = useQuasar()
const userStore = useUserStore()

const formDialog = ref(false)
const editMode = ref(false)
const editId = ref(null)
const form = ref({ username: '', nama: '', email: '', password: '', role: 'guru', is_active: true })

const roleOptions = ['admin', 'operator', 'kepala_sekolah', 'guru', 'siswa', 'orangtua', 'kaprodi', 'wakasek']
const columns = [
  { name: 'username', label: 'Username', field: 'username' },
  { name: 'nama', label: 'Nama', field: 'nama' },
  { name: 'email', label: 'Email', field: 'email' },
  { name: 'role', label: 'Role', field: 'role' },
  { name: 'is_active', label: 'Status', field: 'is_active' },
  { name: 'aksi', label: 'Aksi', field: 'aksi' }
]

function getRoleColor(role) {
  const map = { admin: 'red', operator: 'blue', kepala_sekolah: 'purple', guru: 'green', siswa: 'orange', orangtua: 'brown', kaprodi: 'teal', wakasek: 'cyan' }
  return map[role] || 'grey'
}

function openForm(row = null) {
  if (row) {
    editMode.value = true
    editId.value = row.id
    form.value = { ...row, password: '' }
  } else {
    editMode.value = false
    resetForm()
  }
  formDialog.value = true
}
function resetForm() {
  form.value = { username: '', nama: '', email: '', password: '', role: 'guru', is_active: true }
}
function submitForm() {
  if (editMode.value) {
    userStore.update(editId.value, form.value)
    $q.notify({ type: 'positive', message: 'Pengguna diupdate' })
  } else {
    userStore.tambah(form.value)
    $q.notify({ type: 'positive', message: 'Pengguna ditambahkan' })
  }
  formDialog.value = false
}
function resetPassword(id) {
  const newPass = Math.random().toString(36).substring(2, 8)
  userStore.update(id, { password: newPass })
  $q.notify({ type: 'info', message: `Password baru: ${newPass}` })
}
function confirmDelete(id) {
  $q.dialog({ title: 'Konfirmasi', message: 'Hapus pengguna ini?', cancel: true }).onOk(() => {
    userStore.hapus(id)
    $q.notify({ type: 'positive', message: 'Pengguna dihapus' })
  })
}
onMounted(() => {
  userStore.loadData()
})
</script>
