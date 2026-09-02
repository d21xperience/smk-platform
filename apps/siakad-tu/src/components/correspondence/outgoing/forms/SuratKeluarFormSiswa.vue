<template>
  <q-dialog v-model="dialogSiswa" persistent>
    <q-card style="min-width: 700px; max-width: 90vw;">
      <q-card-section class="row items-center q-pb-none">
        <div class="text-h6">Formulir Surat Siswa</div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup @click="closeDialogSiswa" />
      </q-card-section>

      <q-card-section>
        <q-form @submit.prevent="submitCreateSiswa">
          <q-select v-model="formSiswa.jenis" :options="jenisSuratSiswa" label="Jenis Surat" required
            @update:model-value="val => handleJenisChange('siswa', val)" />

          <template v-if="formSiswa.jenis?.value === 'pkl'">
            <q-select v-model="formSiswa.dudi" :options="dudiOptions" label="Mitra DUDI" required />
            <q-select v-model="formSiswa.siswaIds" :options="daftarSiswaOptions" label="Nama Siswa (bisa multi)"
              multiple use-chips required />
          </template>

          <template v-else-if="formSiswa.jenis?.value === 'aktif'">
            <q-input v-model="formSiswa.namaSiswa" label="Nama Siswa" required />
            <q-input v-model="formSiswa.kelas" label="Kelas" required />
            <q-input v-model="formSiswa.keperluan" label="Keperluan (beasiswa/BPJS/dll)" required />
          </template>

          <template v-else-if="formSiswa.jenis?.value === 'kesalahan_nama'">
            <q-input v-model="formSiswa.namaSalah" label="Nama Salah di Dokumen" required />
            <q-input v-model="formSiswa.namaBenar" label="Nama Benar" required />
            <q-input v-model="formSiswa.keterangan" label="Keterangan Perbedaan" type="textarea" />
          </template>

          <template v-else-if="formSiswa.jenis?.value === 'tunggakan'">
            <q-input v-model="formSiswa.namaSiswa" label="Nama Siswa" required />
            <q-input v-model="formSiswa.nominal" label="Nominal Tunggakan (Rp)" type="number" required />
            <q-input v-model="formSiswa.rincian" label="Rincian" type="textarea" />
          </template>

          <q-input v-model="generatedNomorSiswa" label="Nomor Surat (otomatis)" readonly class="q-mt-md" />

          <div class="q-mt-md text-right">
            <q-btn label="Batal" flat color="grey" v-close-popup @click="closeDialogSiswa" />
            <q-btn label="Simpan & Cetak" type="submit" color="primary" :loading="isLoading" />
          </div>
        </q-form>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { useOutgoingLetter } from '@/composables/useOutgoingLetter.js'

const {
  dialogSiswa,
  formSiswa,
  isLoading,
  jenisSuratSiswa,
  dudiOptions,
  daftarSiswaOptions,
  generatedNomorSiswa,
  handleJenisChange,
  submitCreateSiswa,
  closeDialogSiswa
} = useOutgoingLetter()
</script>
