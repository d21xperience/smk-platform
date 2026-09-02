<template>
  <q-dialog v-model="isPrintDialogOpen" maximized persistent>
    <q-card class="column full-height">
      <!-- Header Dialog -->
      <q-bar class="bg-grey-3">
        <div class="text-subtitle1 q-ml-sm">Preview Cetak Surat</div>
        <q-space />
        <q-btn
          dense
          flat
          icon="print"
          label="Cetak / Save PDF"
          color="primary"
          @click="triggerPrint"
        />
        <q-btn
          dense
          flat
          icon="close"
          v-close-popup
          @click="closePrintDialog"
        />
      </q-bar>

      <!-- Area Preview -->
      <q-card-section class="col bg-grey-2 scroll" style="height: calc(100vh - 50px);">
        <div class="row justify-center">
          <!-- Render Template PKL jika jenisnya 'pkl' -->
          <SuratPklPrintView
            v-if="letterToPrint && letterToPrint.jenis === 'pkl'"
            :letter="letterToPrint"
          />

          <!-- Fallback untuk jenis surat lain -->
          <div v-else class="text-center q-pa-xl text-grey-6 bg-white shadow-1 rounded-borders" style="min-width: 210mm; min-height: 297mm;">
            <q-icon name="description" size="48px" class="q-mb-md" />
            <div class="text-h6">Template Belum Tersedia</div>
            <div class="text-body2 q-mt-sm">
              Template untuk jenis surat "{{ letterToPrint?.jenis || 'tidak diketahui' }}" belum diimplementasikan.
            </div>
          </div>
        </div>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { useOutgoingLetter } from '@/composables/useOutgoingLetter.js'
import SuratPklPrintView from '@/components/correspondence/templates/SuratPklPrintView.vue'

const {
  isPrintDialogOpen,
  letterToPrint,
  closePrintDialog,
  triggerPrint
} = useOutgoingLetter()
</script>

<style scoped>
/* Memastikan area scroll berperilaku baik */
.scroll {
  overflow-y: auto;
  overflow-x: hidden;
}
</style>
