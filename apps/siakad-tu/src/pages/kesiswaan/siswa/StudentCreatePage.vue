<template>
  <q-page padding>
    <div class="row items-center q-mb-md">
      <q-btn flat round icon="arrow_back" color="primary" @click="$router.back()" />
      <h5 class="text-h5 text-weight-bold q-my-none q-ml-sm text-primary">
        {{ isEditMode ? 'Edit Data Siswa' : 'Tambah Siswa Baru' }}
      </h5>
    </div>

    <q-card>
      <q-card-section>
        <q-form @submit="onSubmit" class="q-gutter-md">
          <div class="text-h6 q-mb-md text-secondary">Data Pribadi</div>

          <div class="row q-col-gutter-md">
            <div class="col-12 col-md-4">
              <q-input v-model="form.nisn" label="NISN" outlined :rules="[v => !!v || 'Wajib diisi']" />
            </div>
            <div class="col-12 col-md-4">
              <q-input v-model="form.nis" label="NIS" outlined :rules="[v => !!v || 'Wajib diisi']" />
            </div>
            <div class="col-12 col-md-4">
              <q-select v-model="form.jenisKelamin" :options="[{label:'Laki-laki', value:'L'}, {label:'Perempuan', value:'P'}]"
                label="Jenis Kelamin" outlined emit-value map-options :rules="[v => !!v || 'Wajib dipilih']" />
            </div>
          </div>

          <q-input v-model="form.nama" label="Nama Lengkap" outlined :rules="[v => !!v || 'Wajib diisi']" />

          <div class="row q-col-gutter-md">
            <div class="col-12 col-md-6">
              <q-input v-model="form.tempatLahir" label="Tempat Lahir" outlined />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="form.tanggalLahir" label="Tanggal Lahir" type="date" outlined />
            </div>
          </div>

          <div class="row q-col-gutter-md">
            <div class="col-12 col-md-6">
              <q-input v-model="form.noHp" label="No. HP" outlined />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="form.kelas" label="Kelas (Contoh: X RPL 1)" outlined />
            </div>
          </div>

          <q-input v-model="form.alamat" label="Alamat Lengkap" type="textarea" outlined rows="3" />

          <div class="row q-gutter-sm justify-end q-mt-xl">
            <q-btn flat label="Batal" color="grey" @click="$router.back()" />
            <q-btn type="submit" :label="isEditMode ? 'Simpan Perubahan' : 'Simpan Data Siswa'" color="primary" :loading="loading" />
          </div>
        </q-form>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuasar } from 'quasar'
import { useStudentForm } from '@/composables/student/useStudentForm'

const route = useRoute()
const router = useRouter()
const $q = useQuasar()
const studentId = route.params.id

const { form, loading, isEditMode, loadStudent, saveStudent } = useStudentForm(studentId)

onMounted(() => {
  if (studentId) loadStudent()
})

async function onSubmit() {
  const success = await saveStudent()
  if (success) {
    $q.notify({ type: 'positive', message: 'Data siswa berhasil disimpan' })
    router.push({ name: 'student-list' })
  } else {
    $q.notify({ type: 'negative', message: 'Gagal menyimpan data siswa' })
  }
}
</script>
