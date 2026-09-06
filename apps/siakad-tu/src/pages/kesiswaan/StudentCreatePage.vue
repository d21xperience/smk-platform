<template>
  <q-page padding>
    <DynamicHeader :title="isEditMode ? 'Edit Data Siswa' : 'Tambah Siswa Baru'" />
    <q-card>
      <q-card-section>
        <q-form @submit="onSubmit">
          <!-- Data Pribadi -->
          <div class="text-h6 q-mb-md">Data Pribadi</div>

          <div class="row q-gutter-md">
            <q-input v-model="formData.nisn" label="NISN" :rules="[val => !!val || 'NISN wajib diisi']"
              class="col-12 col-sm-4" />
            <q-input v-model="formData.nis" label="NIS" :rules="[val => !!val || 'NIS wajib diisi']"
              class="col-12 col-sm-4" />
            <q-select v-model="formData.gender" :options="genderOptions" label="Jenis Kelamin" emit-value map-options
              :rules="[val => !!val || 'Jenis kelamin wajib dipilih']" class="col-12 col-sm-4" />
          </div>

          <div class="row q-gutter-md q-mt-md">
            <q-input v-model="formData.fullName.firstName" label="Nama Depan"
              :rules="[val => !!val || 'Nama depan wajib diisi']" class="col-12 col-sm-4" />
            <q-input v-model="formData.fullName.middleName" label="Nama Tengah" class="col-12 col-sm-4" />
            <q-input v-model="formData.fullName.lastName" label="Nama Belakang" class="col-12 col-sm-4" />
          </div>

          <q-input v-model="formData.birthDate" label="Tanggal Lahir" type="date"
            :rules="[val => !!val || 'Tanggal lahir wajib diisi']" class="q-mt-md" />

          <!-- Alamat -->
          <div class="text-h6 q-mt-xl q-mb-md">Alamat</div>

          <q-input v-model="formData.address.street" label="Jalan" class="q-mb-md" />

          <div class="row q-gutter-md">
            <q-input v-model="formData.address.rtRw" label="RT/RW" class="col-12 col-sm-3" />
            <q-input v-model="formData.address.village" label="Kelurahan" class="col-12 col-sm-3" />
            <q-input v-model="formData.address.district" label="Kecamatan" class="col-12 col-sm-3" />
            <q-input v-model="formData.address.postalCode" label="Kode Pos" class="col-12 col-sm-3" />
          </div>

          <q-input v-model="formData.address.city" label="Kota/Kabupaten" class="q-mt-md" />

          <!-- Kontak -->
          <div class="text-h6 q-mt-xl q-mb-md">Kontak</div>

          <div class="row q-gutter-md">
            <q-input v-model="formData.contactInfo.phone" label="No. Telepon" class="col-12 col-sm-6" />
            <q-input v-model="formData.contactInfo.email" label="Email" type="email" class="col-12 col-sm-6" />
          </div>

          <!-- Wali -->
          <div class="text-h6 q-mt-xl q-mb-md">Data Wali</div>

          <div class="row q-gutter-md">
            <q-input v-model="formData.guardianInfo.name" label="Nama Wali" class="col-12 col-sm-6" />
            <q-input v-model="formData.guardianInfo.relation" label="Hubungan" class="col-12 col-sm-6" />
          </div>

          <div class="row q-gutter-md q-mt-md">
            <q-input v-model="formData.guardianInfo.phone" label="No. Telepon Wali" class="col-12 col-sm-6" />
            <q-input v-model="formData.guardianInfo.occupation" label="Pekerjaan Wali" class="col-12 col-sm-6" />
          </div>

          <!-- Actions -->
          <div class="row q-gutter-sm justify-end q-mt-xl">
            <q-btn flat label="Batal" color="grey" @click="onCancel" />
            <q-btn type="submit" :label="isEditMode ? 'Update' : 'Simpan'" color="primary" :loading="isLoading" />
          </div>
        </q-form>
      </q-card-section>
    </q-card>

    <!-- Notification Dialog -->
    <q-dialog v-model="showNotification">
      <q-card>
        <q-card-section :class="notificationType === 'success' ? 'bg-positive text-white' : 'bg-negative text-white'">
          <div class="text-h6">{{ notificationType === 'success' ? 'Berhasil' : 'Gagal' }}</div>
          <div>{{ notificationMessage }}</div>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="OK" color="primary" @click="showNotification = false" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useStudentCreate } from '@/composables/kesiswaan/useStudentCreate'
import DynamicHeader from '@/components/DynamicHeader.vue'

const route = useRoute()

const {
  isLoading,
  showNotification,
  notificationMessage,
  notificationType,
  formData,
  genderOptions,
  isEditMode,
  onSubmit,
  onCancel,
  loadData,
} = useStudentCreate()

onMounted(() => {
  const studentId = route.params.id
  if (studentId) {
    loadData(studentId)
  }
})
</script>
