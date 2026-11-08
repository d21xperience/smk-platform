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
            <!-- Data Identitas -->
            <div class="col-12 col-md-6">
              <q-input v-model="form.nis" label="NIS" outlined :disable="isLocked" />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="form.nisn" label="NISN" outlined :disable="isLocked" />
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
            
            <!-- Data Kelas & Jurusan -->
            <div class="col-12 col-md-6">
              <q-select v-model="form.kelas" :options="classOptions" label="Kelas" outlined emit-value map-options
                :disable="isLocked" :loading="loadingClasses" option-label="label" option-value="value"
                @update:model-value="onClassChange">
                <template v-slot:no-option>
                  <q-item>
                    <q-item-section class="text-grey">
                      Tidak ada data kelas
                    </q-item-section>
                  </q-item>
                </template>
              </q-select>
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="form.jurusan" label="Jurusan" outlined readonly disable bg-color="grey-2" />
            </div>

            <!-- Data Kontak & Orang Tua -->
            <div class="col-12 col-md-6">
              <q-input v-model="form.noHp" label="No. HP Siswa" outlined />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="form.email" label="Email Siswa" type="email" outlined />
            </div>
            
            <div class="col-12 col-md-6">
              <q-input v-model="form.parentName" label="Nama Orang Tua (Ayah)" outlined />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="form.parentPhone" label="No. HP Orang Tua" outlined />
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
const { 
  form, 
  loading, 
  isLocked, 
  classOptions, 
  loadingClasses,
  onClassChange,
  loadStudent, 
  saveStudent 
} = useStudentForm(studentId);

onMounted(async () => {
  await loadStudent();
});

function goBack() {
  router.push({ name: 'kesiswaan-siswa-detail', params: { id: studentId } });
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
