<template>
  <q-page class="q-pa-md">
    <div class="row items-center q-mb-md">
      <div class="text-h5 text-weight-bold">Manajemen Pengguna</div>
      <q-space />
      <q-btn color="primary" label="Tambah Staf" icon="person_add" @click="openAddDialog" />
    </div>

    <q-table
      :rows="users"
      :columns="columns"
      row-key="UserID"
      flat bordered
      :loading="isLoading"
    >
      <template v-slot:body-cell-IsActive="props">
        <q-td :props="props">
          <q-badge :color="props.row.IsActive ? 'green' : 'red'" class="q-px-sm q-py-xs">
            {{ props.row.IsActive ? 'Aktif' : 'Nonaktif' }}
          </q-badge>
        </q-td>
      </template>

      <template v-slot:body-cell-actions="props">
        <q-td :props="props" class="q-gutter-x-sm">
          <q-btn color="warning" icon="edit" size="sm" unelevated @click="openEditDialog(props.row)">
            <q-tooltip>Ubah Role</q-tooltip>
          </q-btn>
          <q-btn color="info" icon="key" size="sm" unelevated @click="openResetDialog(props.row)">
            <q-tooltip>Reset Password</q-tooltip>
          </q-btn>
          <q-btn 
            :color="props.row.IsActive ? 'negative' : 'positive'" 
            :icon="props.row.IsActive ? 'block' : 'check_circle'" 
            size="sm" unelevated 
            @click="toggleStatus(props.row)"
          >
            <q-tooltip>{{ props.row.IsActive ? 'Nonaktifkan Akun' : 'Aktifkan Akun' }}</q-tooltip>
          </q-btn>
        </q-td>
      </template>
    </q-table>

    <!-- Dialog Tambah / Edit Role -->
    <q-dialog v-model="showDialog" persistent>
      <q-card style="min-width: 400px">
        <q-card-section class="row items-center">
          <div class="text-h6">{{ isEditMode ? 'Ubah Role Pengguna' : 'Tambah Pengguna Baru' }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pt-none q-gutter-y-md">
          <q-input v-model="form.username" label="Username / NIP" outlined dense :readonly="isEditMode" />
          <q-input v-if="!isEditMode" v-model="form.password" type="password" label="Password" outlined dense />
          
          <q-select
            v-model="form.role"
            :options="roleOptions"
            label="Role / Jabatan"
            outlined dense
            emit-value
            map-options
          />
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="primary" label="Simpan" :loading="isSubmitting" @click="submitUser" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog Reset Password -->
    <q-dialog v-model="showResetDialog" persistent>
      <q-card style="min-width: 350px">
        <q-card-section class="row items-center">
          <div class="text-h6">Reset Password</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pt-none">
          <div class="text-caption q-mb-sm">Atur ulang password untuk <strong>{{ resetTarget?.Name }}</strong></div>
          <q-input 
            v-model="newPassword" 
            type="password" 
            label="Password Baru" 
            outlined dense 
            :rules="[val => val.length >= 6 || 'Minimal 6 karakter']"
          />
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="info" label="Reset" :loading="isResetting" @click="submitResetPassword" :disable="newPassword.length < 6" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useQuasar, Notify } from 'quasar'
import { api } from '@/boot/axios'

const $q = useQuasar()

const users = ref([])
const isLoading = ref(false)
const isSubmitting = ref(false)
const isResetting = ref(false)
const showDialog = ref(false)
const showResetDialog = ref(false)
const isEditMode = ref(false)
const resetTarget = ref(null)
const newPassword = ref('')

const form = ref({
  id: '',
  username: '',
  password: '',
  role: 'kesiswaan'
})

const roleOptions = [
  { label: 'Administrator', value: 'admin' },
  { label: 'Kesiswaan', value: 'kesiswaan' },
  { label: 'Kepegawaian', value: 'kepegawaian' },
  { label: 'Sarana Prasarana (Sarpras)', value: 'sarpras' },
  { label: 'Hubungan Industri (Hubin)', value: 'hubin' },
]

const columns = [
  { name: 'Name', label: 'Username / NIP', field: 'Name', align: 'left', sortable: true },
  { name: 'Role', label: 'Role', field: 'Role', align: 'left', sortable: true },
  { name: 'IsActive', label: 'Status', field: 'IsActive', align: 'center' },
  { name: 'actions', label: 'Aksi', align: 'center' }
]

async function fetchUsers() {
  isLoading.value = true
  try {
    const res = await api.get('/users')
    users.value = res.data
  } catch (err) {
    Notify.create({ type: 'negative', message: 'Gagal memuat data pengguna' })
  } finally {
    isLoading.value = false
  }
}

function openAddDialog() {
  isEditMode.value = false
  form.value = { id: '', username: '', password: '', role: 'kesiswaan' }
  showDialog.value = true
}

function openEditDialog(user) {
  isEditMode.value = true
  form.value = {
    id: user.UserID,
    username: user.Name,
    password: '',
    role: user.Role
  }
  showDialog.value = true
}

function openResetDialog(user) {
  resetTarget.value = user
  newPassword.value = ''
  showResetDialog.value = true
}

async function submitUser() {
  isSubmitting.value = true
  try {
    if (isEditMode.value) {
      await api.put(`/users/${form.value.id}`, { role: form.value.role })
      Notify.create({ type: 'positive', message: 'Role berhasil diperbarui' })
    } else {
      await api.post('/users', {
        username: form.value.username,
        password: form.value.password,
        role: form.value.role
      })
      Notify.create({ type: 'positive', message: 'Pengguna berhasil ditambahkan' })
    }
    showDialog.value = false
    await fetchUsers()
  } catch (err) {
    Notify.create({ type: 'negative', message: err.response?.data?.error || 'Gagal menyimpan data' })
  } finally {
    isSubmitting.value = false
  }
}

async function submitResetPassword() {
  isResetting.value = true
  try {
    await api.put(`/users/${resetTarget.value.UserID}/reset-password`, {
      password: newPassword.value
    })
    Notify.create({ type: 'positive', message: 'Password berhasil direset' })
    showResetDialog.value = false
  } catch (err) {
    Notify.create({ type: 'negative', message: err.response?.data?.error || 'Gagal mereset password' })
  } finally {
    isResetting.value = false
  }
}

function toggleStatus(user) {
  $q.dialog({
    title: 'Konfirmasi',
    message: `Apakah Anda yakin ingin ${user.IsActive ? 'menonaktifkan' : 'mengaktifkan'} akun ${user.Name}?`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await api.put(`/users/${user.UserID}/toggle`, { is_active: !user.IsActive })
      Notify.create({ type: 'positive', message: 'Status pengguna berhasil diubah' })
      await fetchUsers()
    } catch (err) {
      Notify.create({ type: 'negative', message: err.response?.data?.error || 'Gagal mengubah status' })
    }
  })
}

onMounted(fetchUsers)
</script>
