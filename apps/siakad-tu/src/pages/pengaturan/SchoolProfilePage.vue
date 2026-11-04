<template>
  <q-page class="q-pa-md">
    <div class="row items-center q-mb-md">
      <div class="text-h5 text-weight-bold">Profil Sekolah</div>
    </div>

    <q-card flat bordered class="q-pa-md">
      <q-form @submit="saveProfile" class="q-gutter-md">
        <div class="row q-col-gutter-md">
          <!-- Kolom Kiri -->
          <div class="col-12 col-md-6 q-gutter-y-md">
            <div class="text-subtitle1 text-weight-bold text-primary">Identitas Utama</div>
            <q-input v-model="form.Name" label="Nama Sekolah" outlined dense :rules="[val => !!val || 'Wajib diisi']" />
            <q-input v-model="form.NPSN" label="NPSN" outlined dense />
            <q-input v-model="form.NSS" label="NSS" outlined dense />
            <q-select v-model="form.Jenjang" :options="['SD', 'SMP', 'SMA', 'SMK']" label="Jenjang" outlined dense />
            <q-input v-model.number="form.BentukPendidikan" type="number"
              label="Bentuk Pendidikan ID (Cth: 24 untuk SMK)" outlined dense />
            <q-input v-model.number="form.StatusKepemilikan" type="number" label="Status Kepemilikan ID (Cth: 1)"
              outlined dense />
          </div>

          <!-- Kolom Kanan -->
          <div class="col-12 col-md-6 q-gutter-y-md">
            <div class="text-subtitle1 text-weight-bold text-primary">Kepala Sekolah</div>
            <q-input v-model="form.PrincipalName" label="Nama Kepala Sekolah" outlined dense />
            <q-input v-model="form.PrincipalNIP" label="NIP Kepala Sekolah" outlined dense />
            <q-input v-model="form.PrincipalNIY" label="NIY Kepala Sekolah" outlined dense />

            <div class="text-subtitle1 text-weight-bold text-primary q-mt-lg">Kontak & Alamat</div>
            <q-input v-model="form.Address" label="Alamat Jalan" type="textarea" autogrow outlined dense />
            <div class="row q-col-gutter-sm">
              <div class="col-6"><q-input v-model="form.Kelurahan" label="Kelurahan/Desa" outlined dense /></div>
              <div class="col-6"><q-input v-model="form.Kecamatan" label="Kecamatan" outlined dense /></div>
              <div class="col-6"><q-input v-model="form.KabKota" label="Kabupaten/Kota" outlined dense /></div>
              <div class="col-6"><q-input v-model="form.Propinsi" label="Provinsi" outlined dense /></div>
              <div class="col-6"><q-input v-model="form.PostalCode" label="Kode Pos" outlined dense /></div>
            </div>
            <div class="row q-col-gutter-sm">
              <div class="col-6"><q-input v-model="form.Phone" label="Telepon" outlined dense /></div>
              <div class="col-6"><q-input v-model="form.Fax" label="Fax" outlined dense /></div>
              <div class="col-6"><q-input v-model="form.Email" type="email" label="Email" outlined dense /></div>
              <div class="col-6"><q-input v-model="form.Website" label="Website" outlined dense /></div>
            </div>
          </div>
        </div>

        <div class="row justify-end q-mt-lg">
          <q-btn label="Simpan Perubahan" type="submit" color="primary" icon="save" :loading="isSaving" />
        </div>
      </q-form>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Notify } from 'quasar'
import { api } from '@/boot/axios'

const isLoading = ref(false)
const isSaving = ref(false)

const form = ref({
  SchoolID: '',
  Name: '',
  NPSN: '',
  NSS: '',
  Address: '',
  PostalCode: '',
  Phone: '',
  Fax: '',
  Kelurahan: '',
  Kecamatan: '',
  KabKota: '',
  Propinsi: '',
  Website: '',
  Email: '',
  PrincipalName: '',
  PrincipalNIP: '',
  PrincipalNIY: '',
  StatusKepemilikan: 1,
  KodeAktivasi: '',
  Jenjang: '',
  BentukPendidikan: 0
})

async function fetchProfile() {
  isLoading.value = true
  try {
    const res = await api.get('/school')
    form.value = res.data
  } catch (err) {
    Notify.create({ type: 'negative', message: `Gagal memuat profil sekolah ${err}` })
  } finally {
    isLoading.value = false
  }
}

async function saveProfile() {
  isSaving.value = true
  try {
    await api.put('/school', form.value)
    Notify.create({ type: 'positive', message: 'Profil sekolah berhasil diperbarui' })
    // Memicu reload halaman agar nama sekolah di header langsung ter-update jika aplikasi di-refresh
    window.location.reload()
  } catch (err) {
    Notify.create({ type: 'negative', message: err.response?.data?.error || 'Gagal menyimpan profil' })
  } finally {
    isSaving.value = false
  }
}

onMounted(fetchProfile)
</script>
