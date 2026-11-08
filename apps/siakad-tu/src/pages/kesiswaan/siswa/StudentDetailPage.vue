<template>
  <q-page class="q-pa-md bg-grey-2">
    <div v-if="isLoading" class="flex flex-center q-pa-xl">
      <q-spinner color="primary" size="3em" />
    </div>

    <div v-else-if="student">
      <!-- Header Card -->
      <q-card class="shadow-1 q-mb-md q-pa-md border-radius-xl">
        <div class="row items-center justify-between">
          <div class="row items-center">
            <q-avatar size="64px" color="primary" text-color="white" icon="person" class="q-mr-md" />
            <div>
              <div class="text-h5 text-weight-bold text-primary">{{ student.fullName }}</div>
              <div class="text-subtitle1 text-grey-8">
                NIS: <strong>{{ student.nis }}</strong> | NISN: <strong>{{ student.nisn }}</strong>
              </div>
              <div class="q-mt-xs">
                <q-badge :color="getStatusColor(student.status)" class="text-subtitle2 q-px-sm">
                  {{ student.status }}
                </q-badge>
                <q-badge color="blue-grey" class="text-subtitle2 q-px-sm q-ml-sm">
                  Kelas: {{ student.className || 'Belum ada kelas' }}
                </q-badge>
              </div>
            </div>
          </div>
          <div class="q-gutter-sm">
            <q-btn color="primary" icon="edit" label="Edit Data" @click="editStudent" />
            <q-btn color="grey-7" icon="arrow_back" label="Kembali" @click="goBack" />
          </div>
        </div>
      </q-card>

      <!-- Tabs Informasi -->
      <q-card class="shadow-1 border-radius-xl">
        <q-tabs v-model="tab" dense class="text-grey" active-color="primary" indicator-color="primary"
          align="justify" narrow-indicator>
          <q-tab name="profil" icon="badge" label="Profil Siswa" />
          <q-tab name="orangtua" icon="family_restroom" label="Data Orang Tua" />
          <q-tab name="alamat" icon="home" label="Alamat & Kontak" />
        </q-tabs>

        <q-separator />

        <q-tab-panels v-model="tab" animated>
          <!-- TAB PROFIL -->
          <q-tab-panel name="profil">
            <div class="row q-col-gutter-lg q-mt-sm">
              <div class="col-12 col-sm-6">
                <div class="text-caption text-grey-7">Nama Lengkap</div>
                <div class="text-body1 text-weight-medium">{{ student.fullName }}</div>
              </div>
              <div class="col-12 col-sm-6">
                <div class="text-caption text-grey-7">Jenis Kelamin</div>
                <div class="text-body1 text-weight-medium">{{ student.gender === 'L' ? 'Laki-laki' : 'Perempuan' }}</div>
              </div>
              <div class="col-12 col-sm-6">
                <div class="text-caption text-grey-7">Tempat Lahir</div>
                <div class="text-body1 text-weight-medium">{{ student.birthPlace || '-' }}</div>
              </div>
              <div class="col-12 col-sm-6">
                <div class="text-caption text-grey-7">Tanggal Lahir</div>
                <div class="text-body1 text-weight-medium">{{ formatDate(student.birthDate) }}</div>
              </div>
            </div>
          </q-tab-panel>

          <!-- TAB ORANG TUA -->
          <q-tab-panel name="orangtua">
            <div class="row q-col-gutter-lg q-mt-sm">
              <div class="col-12 col-sm-6">
                <div class="text-caption text-grey-7">Nama Ayah</div>
                <div class="text-body1 text-weight-medium">{{ student.parentName || '-' }}</div>
              </div>
              <div class="col-12 col-sm-6">
                <div class="text-caption text-grey-7">Nama Ibu</div>
                <div class="text-body1 text-weight-medium">{{ student.motherName || '-' }}</div>
              </div>
            </div>
          </q-tab-panel>

          <!-- TAB ALAMAT & KONTAK -->
          <q-tab-panel name="alamat">
            <div class="row q-col-gutter-lg q-mt-sm">
              <div class="col-12">
                <div class="text-caption text-grey-7">Alamat Lengkap</div>
                <div class="text-body1 text-weight-medium">{{ student.address || 'Belum diisi' }}</div>
              </div>
              <div class="col-12 col-sm-6">
                <div class="text-caption text-grey-7">No. Telepon Siswa</div>
                <div class="text-body1 text-weight-medium">{{ student.phone || '-' }}</div>
              </div>
            </div>
          </q-tab-panel>
        </q-tab-panels>
      </q-card>
    </div>

    <!-- State Jika Data Tidak Ditemukan -->
    <div v-else class="text-center q-pa-xl">
      <q-icon name="search_off" size="4em" color="grey-5" />
      <div class="text-h6 text-grey-6 q-mt-md">Data siswa tidak ditemukan</div>
      <q-btn color="primary" label="Kembali ke Daftar" class="q-mt-md" @click="goBack" />
    </div>
  </q-page>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useStudentDetail } from '@/composables/student/useStudentDetail'

const tab = ref('profil')

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

<style scoped>
.border-radius-xl {
  border-radius: 12px;
}
</style>
