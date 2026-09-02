<template>
  <q-dialog :model-value="modelValue" @update:model-value="$emit('update:modelValue', $event)" persistent>
    <q-card style="min-width: 500px">
      <q-card-section class="row items-center q-pb-none">
        <div class="text-h6">Buku Ekspedisi</div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup />
      </q-card-section>
      <q-card-section>
        <q-form @submit.prevent="catatEkspedisi">
          <q-select v-model="formEkspedisi.selectedSurat" :options="semuaSuratOptions" label="Pilih Surat Keluar" />
          <q-input v-model="formEkspedisi.tanggalKirim" label="Tanggal Pengiriman" type="date" required />
          <q-input v-model="formEkspedisi.kurir" label="Kurir/Ekspedisi" required />
          <q-input v-model="formEkspedisi.tandaTerima" label="Tanda Terima Digital (link/upload)" />
          <div class="q-mt-md text-right">
            <q-btn label="Batal" flat color="grey" v-close-popup />
            <q-btn label="Catat" type="submit" color="primary" />
          </div>
        </q-form>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { useQuasar } from 'quasar'
import { useOutgoingLetter } from '@/composables/useOutgoingLetter.js'

defineProps({
  modelValue: { type: Boolean, required: true }
})

defineEmits(['update:modelValue'])

const $q = useQuasar()
const { formEkspedisi, semuaSuratOptions } = useOutgoingLetter()

function catatEkspedisi() {
  $q.notify({ type: 'positive', message: 'Ekspedisi tercatat' })
}
</script>
