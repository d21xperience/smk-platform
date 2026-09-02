<template>
  <q-page padding>
    <q-inner-loading :showing="loading" />

    <div class="q-mb-md flex items-center justify-between">
      <div class="flex items-center">
        <q-btn flat round icon="arrow_back" @click="goBack" class="q-mr-sm" />
        <h5 class="q-my-none text-primary">Detail Siswa</h5>
        <q-badge v-if="student?.isLocked" color="warning" class="q-ml-md" label="Dapodik" />
        <q-badge :color="student?.status === 'AKTIF' ? 'positive' : 'negative'" class="q-ml-sm"
          :label="student?.status" />
      </div>
      <div class="q-gutter-sm">
        <q-btn color="primary" icon="edit" label="Ubah" @click="goToEdit" />
        <q-btn color="orange" icon="compare_arrows" label="Mutasi" @click="goToMutate" />
      </div>
    </div>

    <q-card v-if="student">
      <q-card-section>
        <div class="text-h6 q-mb-md">Data Pribadi</div>
        <q-list bordered separator class="rounded-borders">
          <q-item>
            <q-item-section side><q-item-label caption>NISN</q-item-label></q-item-section>
            <q-item-section>{{ student.nisn }}</q-item-section>
          </q-item>
          <q-item>
            <q-item-section side><q-item-label caption>Nama Lengkap</q-item-label></q-item-section>
            <q-item-section>{{ student.nama }}</q-item-section>
          </q-item>
          <q-item>
            <q-item-section side><q-item-label caption>Tempat, Tanggal Lahir</q-item-label></q-item-section>
            <q-item-section>{{ student.tempatLahir }}, {{ student.tanggalLahir }}</q-item-section>
          </q-item>
          <q-item>
            <q-item-section side><q-item-label caption>Jenis Kelamin</q-item-label></q-item-section>
            <q-item-section>{{ student.jenisKelamin === 'L' ? 'Laki-laki' : 'Perempuan' }}</q-item-section>
          </q-item>
          <q-item>
            <q-item-section side><q-item-label caption>Agama</q-item-label></q-item-section>
            <q-item-section>{{ student.agama }}</q-item-section>
          </q-item>
        </q-list>
      </q-card-section>

      <q-separator />

      <q-card-section>
        <div class="text-h6 q-mb-md">Data Akademik & Kontak</div>
        <q-list bordered separator class="rounded-borders">
          <q-item>
            <q-item-section side><q-item-label caption>Asal Sekolah</q-item-label></q-item-section>
            <q-item-section>{{ student.asalSekolah }}</q-item-section>
          </q-item>
          <q-item>
            <q-item-section side><q-item-label caption>Jurusan / Kelas</q-item-label></q-item-section>
            <q-item-section>{{ student.jurusan }} - {{ student.kelas }}</q-item-section>
          </q-item>
          <q-item>
            <q-item-section side><q-item-label caption>No. HP</q-item-label></q-item-section>
            <q-item-section>{{ student.noHp || '-' }}</q-item-section>
          </q-item>
          <q-item>
            <q-item-section side><q-item-label caption>Email</q-item-label></q-item-section>
            <q-item-section>{{ student.email || '-' }}</q-item-section>
          </q-item>
          <q-item>
            <q-item-section side><q-item-label caption>Alamat</q-item-label></q-item-section>
            <q-item-section>{{ student.alamat || '-' }}</q-item-section>
          </q-item>
        </q-list>
      </q-card-section>
    </q-card>

    <q-card v-else-if="!loading" class="q-pa-md text-center text-grey-6">
      Data siswa tidak ditemukan.
    </q-card>
  </q-page>
</template>

<script setup>
import { onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useStudentDetail } from '@/composables/student/useStudentDetail';

const router = useRouter();
const route = useRoute();

const studentId = route.params.id;
const { student, loading, loadDetail } = useStudentDetail(studentId);

onMounted(() => {
  loadDetail();
});

function goBack() {
  router.push({ name: 'student-list' });
}

function goToEdit() {
  router.push({ name: 'student-edit', params: { id: studentId } });
}

function goToMutate() {
  router.push({ name: 'manajemen-mutasi-siswa', query: { studentId } });
}
</script>
