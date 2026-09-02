<template>
  <q-page padding>
    <q-inner-loading :showing="loading" />

    <div class="q-mb-md flex items-center">
      <q-btn flat round icon="arrow_back" @click="goBack" class="q-mr-sm" />
      <h5 class="q-my-none text-primary">Tambah Siswa Baru</h5>
    </div>

    <q-card>
      <q-card-section>
        <q-form @submit="onSubmit" class="">
          <div class="row q-col-gutter-md">
            <div class="col-12 col-md-6">
              <q-input v-model="form.nisn" label="NISN" outlined :rules="[v => !!v || 'NISN Wajib Diisi']" />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="form.nama" label="Nama Lengkap" outlined :rules="[v => !!v || 'Nama Wajib Diisi']" />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="form.tempatLahir" label="Tempat Lahir" outlined />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="form.tanggalLahir" label="Tanggal Lahir" type="date" outlined />
            </div>
            <div class="col-12 col-md-6">
              <q-select v-model="form.jenisKelamin"
                :options="[{ label: 'Laki-laki', value: 'L' }, { label: 'Perempuan', value: 'P' }]" label="Jenis Kelamin"
                outlined emit-value map-options />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="form.agama" label="Agama" outlined />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="form.noHp" label="No. HP" outlined />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="form.email" label="Email" type="email" outlined />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="form.asalSekolah" label="Asal Sekolah" outlined />
            </div>
            <div class="col-12 col-md-6">
              <q-select v-model="form.jurusan" :options="['RPL', 'TKRO', 'AKL', 'ULP', 'DPIB']" label="Jurusan" outlined
                emit-value map-options />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="form.kelas" label="Kelas (Contoh: X RPL 1)" outlined />
            </div>
            <div class="col-12 col-12">
              <q-input v-model="form.alamat" label="Alamat Lengkap" type="textarea" rows="2" outlined />
            </div>
          </div>

          <div class="flex justify-end q-gutter-sm q-mt-md">
            <q-btn flat label="Batal" color="negative" @click="goBack" />
            <q-btn type="submit" label="Simpan" color="primary" :loading="loading" />
          </div>
        </q-form>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { useRouter } from 'vue-router';
import { useQuasar } from 'quasar';
import { useStudentForm } from '@/composables/student/useStudentForm';

const router = useRouter();
const $q = useQuasar();
const { form, loading, saveStudent } = useStudentForm();

function goBack() {
  router.push({ name: 'student-list' });
}

async function onSubmit() {
  const success = await saveStudent();
  if (success) {
    $q.notify({ type: 'positive', message: 'Siswa berhasil ditambahkan' });
    router.push({ name: 'student-list' });
  } else {
    $q.notify({ type: 'negative', message: 'Gagal menambahkan siswa' });
  }
}
</script>
