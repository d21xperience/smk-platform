<template>
  <q-page padding>
    <div v-if="isLoading" class="flex flex-center q-pa-xl">
      <q-spinner color="primary" size="3em" />
    </div>

    <div v-else-if="student">
      <div class="row items-center justify-between q-mb-md">
        <div>
          <div class="text-h5 text-weight-bold">
            {{ student.fullName?.firstName }} {{ student.fullName?.lastName }}
          </div>
          <div class="text-grey-7">NISN: {{ student.nisn }} | NIS: {{ student.nis }}</div>
        </div>
        <div>
          <q-btn color="primary" icon="edit" label="Edit" class="q-mr-sm" @click="editStudent" />
          <q-btn color="grey" icon="arrow_back" label="Kembali" @click="goBack" />
        </div>
      </div>

      <!-- Status Badge -->
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
              <div>{{ formatGender(student.gender) }}</div>
            </div>
            <div class="col-12 col-sm-6">
              <div class="text-caption text-grey-7">Tanggal Lahir</div>
              <div>{{ formatDate(student.birthDate) }}</div>
            </div>
          </div>
        </q-card-section>
      </q-card>

      <!-- Alamat -->
      <q-card class="q-mb-md">
        <q-card-section>
          <div class="text-h6 q-mb-md">Alamat</div>
          <div v-if="student.address">
            <div>{{ student.address.street }}</div>
            <div v-if="student.address.rtRw">RT/RW: {{ student.address.rtRw }}</div>
            <div>
              {{ student.address.village }}{{ student.address.district ? ', ' + student.address.district : '' }}
            </div>
            <div>
              {{ student.address.city }}{{ student.address.postalCode ? ' ' + student.address.postalCode : '' }}
            </div>
          </div>
          <div v-else class="text-grey-6">-</div>
        </q-card-section>
      </q-card>

      <!-- Kontak -->
      <q-card class="q-mb-md">
        <q-card-section>
          <div class="text-h6 q-mb-md">Kontak</div>
          <div class="row q-gutter-md">
            <div class="col-12 col-sm-6">
              <div class="text-caption text-grey-7">No. Telepon</div>
              <div>{{ student.contactInfo?.phone || '-' }}</div>
            </div>
            <div class="col-12 col-sm-6">
              <div class="text-caption text-grey-7">Email</div>
              <div>{{ student.contactInfo?.email || '-' }}</div>
            </div>
          </div>
        </q-card-section>
      </q-card>

      <!-- Wali -->
      <q-card>
        <q-card-section>
          <div class="text-h6 q-mb-md">Data Wali</div>
          <div v-if="student.guardianInfo">
            <div class="row q-gutter-md">
              <div class="col-12 col-sm-6">
                <div class="text-caption text-grey-7">Nama Wali</div>
                <div>{{ student.guardianInfo.name || '-' }}</div>
              </div>
              <div class="col-12 col-sm-6">
                <div class="text-caption text-grey-7">Hubungan</div>
                <div>{{ student.guardianInfo.relation || '-' }}</div>
              </div>
            </div>
            <div class="row q-gutter-md q-mt-md">
              <div class="col-12 col-sm-6">
                <div class="text-caption text-grey-7">No. Telepon</div>
                <div>{{ student.guardianInfo.phone || '-' }}</div>
              </div>
              <div class="col-12 col-sm-6">
                <div class="text-caption text-grey-7">Pekerjaan</div>
                <div>{{ student.guardianInfo.occupation || '-' }}</div>
              </div>
            </div>
          </div>
          <div v-else class="text-grey-6">-</div>
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
import { useStudentDetail } from '@/composables/kesiswaan/useStudentDetail'

const {
  student,
  isLoading,
  formatGender,
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
