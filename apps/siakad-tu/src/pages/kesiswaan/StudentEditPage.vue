<template>
  <q-page padding>
    <q-inner-loading :showing="loading" />

    <div class="q-mb-md flex items-center">
      <q-btn flat round icon="arrow_back" @click="goBack" class="q-mr-sm" />
      <h5 class="q-my-none text-primary">Ubah Data Siswa</h5>
      <q-badge v-if="isLocked" color="warning" class="q-ml-md" label="Terkunci (Dapodik)" />
    </div>

    <q-card>
      <q-card-section v-if="isLocked" class="bg-warning text-dark flex items-center">
        <q-icon name="warning" class="q-mr-sm" size="sm" />
        <div class="text-caption">Data identitas inti terkunci karena berasal dari sinkronisasi Dapodik. Hanya data
          kontak/alamat yang dapat diperbarui.</div>
      </q-card-section>

      <q-card-section>
        <q-form @submit="onSubmit" class="q-gutter-md">
          <div class="row q-col-gutter-md">
            <div class="col-12 col-md-6">
              <q-input v-model="form.nisn" label="NISN" outlined disable />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="form.nama" label="Nama Lengkap" outlined :disable="isLocked"
                :rules="[v => !!v || 'Nama Wajib Diisi']" />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="form.tempatLahir" label="Tempat Lahir" outlined :disable="isLocked" />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="form.tanggalLahir" label="Tanggal Lahir" type="date" outlined :disable="isLocked" />
            </div>
            <div class="col-12 col-md-6">
              <q-select v-model="form.jenisKelamin"
                :options="[{ label: 'Laki-laki', value: 'L' }, { label: 'Perempuan', value: 'P' }]" label="Jenis Kelamin"
                outlined emit-value map-options :disable="isLocked" />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="form.agama" label="Agama" outlined :disable="isLocked" />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="form.noHp" label="No. HP" outlined />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="form.email" label="Email" type="email" outlined />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="form.asalSekolah" label="Asal Sekolah" outlined :disable="isLocked" />
            </div>
            <div class="col-12 col-md-6">
              <q-select v-model="form.jurusan" :options="['RPL', 'TKRO', 'AKL', 'ULP', 'DPIB']" label="Jurusan" outlined
                emit-value map-options :disable="isLocked" />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="form.kelas" label="Kelas" outlined :disable="isLocked" />
            </div>
            <div class="col-12 col-12">
              <q-input v-model="form.alamat" label="Alamat Lengkap" type="textarea" rows="2" outlined />
            </div>
          </div>

          <div class="flex justify-end q-gutter-sm q-mt-md">
            <q-btn flat label="Batal" color="negative" @click="goBack" />
            <q-btn type="submit" label="Simpan Perubahan" color="primary" :loading="loading" />
          </div>
        </q-form>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useQuasar } from 'quasar';
import { useStudentForm } from '@/composables/student/useStudentForm';

const router = useRouter();
const route = useRoute();
const $q = useQuasar();

const studentId = route.params.id;
const { form, loading, isLocked, loadStudent, saveStudent } = useStudentForm(studentId);

onMounted(() => {
  loadStudent();
});

function goBack() {
  router.push({ name: 'student-detail', params: { id: studentId } });
}

async function onSubmit() {
  const success = await saveStudent();
  if (success) {
    $q.notify({ type: 'positive', message: 'Data siswa berhasil diperbarui' });
    goBack();
  } else {
    $q.notify({ type: 'negative', message: 'Gagal memperbarui data siswa' });
  }
}
</script>
