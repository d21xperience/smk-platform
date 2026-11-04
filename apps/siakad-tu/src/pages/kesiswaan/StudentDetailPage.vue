<template>
  <q-page padding>
    <div v-if="isLoading" class="flex flex-center q-pa-xl">
      <q-spinner color="primary" size="3em" />
    </div>

    <div v-else-if="student">
      <div class="row items-center justify-between q-mb-md">
        <div>
          <!-- Perbaikan: fullName sekarang String, bukan Object -->
          <div class="text-h5 text-weight-bold">{{ student.fullName }}</div>
          <div class="text-grey-7">NISN: {{ student.nisn }} | NIS: {{ student.nis }}</div>
        </div>
        <div>
          <q-btn color="primary" icon="edit" label="Edit" class="q-mr-sm" @click="editStudent" />
          <q-btn color="grey" icon="arrow_back" label="Kembali" @click="goBack" />
        </div>
      </div>

      <q-badge :color="getStatusColor(student.status)" class="q-mb-md">
        {{ student.status }}
      </q-badge>

      <!-- Data Pribadi -->
      <q-card class="q-mb-md">
        <q-card-section>
          <div class="text-h6 q-mb-md">Data Pribadi</div>
          <div class="row q-gutter-md">
            <div class="col-12 col-sm-6">
              <div class="text-caption text-grey-7">Jenis Kelamin</div>
              <div>{{ student.gender === 'L' ? 'Laki-laki' : 'Perempuan' }}</div>
            </div>
            <div class="col-12 col-sm-6">
              <div class="text-caption text-grey-7">Tempat, Tanggal Lahir</div>
              <div>{{ student.birthPlace || '-' }}, {{ formatDate(student.birthDate) }}</div>
            </div>
            <div class="col-12 col-sm-6">
              <div class="text-caption text-grey-7">Kelas Saat Ini</div>
              <div>{{ student.className || '-' }}</div>
            </div>
            <div class="col-12 col-sm-6">
              <div class="text-caption text-grey-7">Wali Kelas</div>
              <div>{{ student.waliKelas || '-' }}</div>
            </div>
          </div>
        </q-card-section>
      </q-card>

      <!-- Alamat -->
      <q-card class="q-mb-md">
        <q-card-section>
          <div class="text-h6 q-mb-md">Alamat</div>
          <div>{{ student.address || '-' }}</div>
        </q-card-section>
      </q-card>

      <!-- Kontak & Orang Tua -->
      <q-card>
        <q-card-section>
          <div class="text-h6 q-mb-md">Kontak & Data Orang Tua</div>
          <div class="row q-gutter-md">
            <div class="col-12 col-sm-6">
              <div class="text-caption text-grey-7">No. Telepon Siswa</div>
              <div>{{ student.phone || '-' }}</div>
            </div>
            <div class="col-12 col-sm-6">
              <div class="text-caption text-grey-7">Nama Ayah</div>
              <div>{{ student.parentName || '-' }}</div>
            </div>
            <div class="col-12 col-sm-6">
              <div class="text-caption text-grey-7">Nama Ibu</div>
              <div>{{ student.motherName || '-' }}</div>
            </div>
          </div>
        </q-card-section>
      </q-card>
    </div>

    <div v-else class="text-center q-pa-xl">
      <div class="text-h6 text-grey-6">Data siswa tidak ditemukan</div>
      <q-btn color="primary" label="Kembali ke Daftar" class="q-mt-md" @click="goBack" />
    </div>
  </q-page>
</template>

<script setup>
import { onMounted } from 'vue'
import { useStudentDetail } from '@/composables/student/useStudentDetail'

const {
  student,
  isLoading,
  formatDate,
  getStatusColor,
  editStudent,
  goBack,
  loadData,
} = useStudentDetail()

onMounted(() => {
  loadData()
})
</script>
