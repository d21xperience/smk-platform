<template>
  <q-dialog v-model="dialogPtk" persistent>
    <q-card style="min-width: 700px; max-width: 90vw;">
      <q-card-section class="row items-center q-pb-none">
        <div class="text-h6">Formulir Surat Guru & Staf</div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup />
      </q-card-section>
      <q-card-section>
        <q-form @submit.prevent="submitCreatePtk">
          <q-select v-model="formPtk.jenis" :options="jenisSuratPtk" label="Jenis Surat" required
            @update:model-value="val => handleJenisChange('ptk', val)" />

          <template v-if="formPtk.jenis === 'aktif_tugas'">
            <q-input v-model="formPtk.namaGuru" label="Nama Guru/Staf" required />
            <q-input v-model="formPtk.nip" label="NIP" />
            <q-input v-model="formPtk.jabatan" label="Jabatan" required />
            <q-input v-model="formPtk.keperluan" label="Keperluan (sertifikasi/kuliah/dll)" required />
          </template>

          <template v-else-if="formPtk.jenis === 'sppd'">
            <q-input v-model="formPtk.namaGuru" label="Nama Pelaksana" required />
            <q-input v-model="formPtk.tujuan" label="Tempat Tujuan" required />
            <q-input v-model="formPtk.tanggalBerangkat" label="Tanggal Berangkat" type="date" required />
            <q-input v-model="formPtk.tanggalKembali" label="Tanggal Kembali" type="date" required />
            <q-uploader label="Unggah Draft SPJ (opsional)" accept=".pdf,.xlsx" max-files="1" auto-upload
              class="q-mt-md" />
          </template>

          <template v-else-if="formPtk.jenis === 'mgmp'">
            <q-input v-model="formPtk.namaGuru" label="Nama Guru" required />
            <q-input v-model="formPtk.kegiatan" label="Nama Kegiatan MGMP" required />
            <q-input v-model="formPtk.tanggalKegiatan" label="Tanggal Kegiatan" type="date" required />
            <q-uploader label="Unggah Surat Undangan (dari WA/email)" accept="image/*,.pdf" max-files="1" auto-upload
              required class="q-mt-md" />
          </template>

          <q-input v-model="generatedNomorPtk" label="Nomor Surat (otomatis)" readonly class="q-mt-md" />

          <div class="q-mt-md text-right">
            <q-btn label="Batal" flat color="grey" v-close-popup />
            <q-btn label="Simpan" type="submit" color="primary" :loading="isLoading" />
          </div>
        </q-form>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { useOutgoingLetter } from '@/composables/useOutgoingLetter.js'

const {
  dialogPtk,
  formPtk,
  isLoading,
  jenisSuratPtk,
  generatedNomorPtk,
  handleJenisChange,
  submitCreatePtk
} = useOutgoingLetter()
</script>
