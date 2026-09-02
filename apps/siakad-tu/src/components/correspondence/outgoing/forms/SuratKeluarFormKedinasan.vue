<template>
  <q-dialog v-model="dialogKedinasan" persistent>
    <q-card style="min-width: 700px; max-width: 90vw;">
      <q-card-section class="row items-center q-pb-none">
        <div class="text-h6">Formulir Surat Kedinasan</div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup />
      </q-card-section>
      <q-card-section>
        <q-form @submit.prevent="submitCreateKedinasan">
          <q-select v-model="formKedinasan.jenis" :options="jenisSuratKedinasan" label="Jenis Surat" required
            @update:model-value="val => handleJenisChange('kedinasan', val)" />

          <template v-if="formKedinasan.jenis === 'undangan'">
            <q-input v-model="formKedinasan.tujuan" label="Tujuan (misal: Orang Tua/Wali Murid)" required />
            <q-input v-model="formKedinasan.acara" label="Acara" required />
            <q-input v-model="formKedinasan.tanggalRapat" label="Tanggal Rapat" type="date" required />
            <q-input v-model="formKedinasan.tempat" label="Tempat" required />
          </template>

          <template v-else-if="formKedinasan.jenis === 'mou'">
            <q-input v-model="formKedinasan.namaMitra" label="Nama Mitra/DUDI" required />
            <q-input v-model="formKedinasan.bidangKerjasama" label="Bidang Kerjasama" required />
            <q-input v-model="formKedinasan.tujuanKerjasama" label="Tujuan Kerjasama" type="textarea" />
          </template>

          <template v-else-if="formKedinasan.jenis === 'balasan'">
            <q-input v-model="formKedinasan.noSuratDituju" label="Nomor Surat yang Dibalas" required />
            <q-input v-model="formKedinasan.isiBalasan" label="Isi Balasan" type="textarea" required />
          </template>

          <template v-else-if="formKedinasan.jenis === 'sk'">
            <q-input v-model="formKedinasan.tentang" label="Tentang" required />
            <q-input v-model="formKedinasan.uraian" label="Uraian Singkat" type="textarea" />
          </template>

          <q-input v-model="generatedNomorKedinasan" label="Nomor Surat (otomatis)" readonly class="q-mt-md" />

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
  dialogKedinasan,
  formKedinasan,
  isLoading,
  jenisSuratKedinasan,
  generatedNomorKedinasan,
  handleJenisChange,
  submitCreateKedinasan
} = useOutgoingLetter()
</script>
