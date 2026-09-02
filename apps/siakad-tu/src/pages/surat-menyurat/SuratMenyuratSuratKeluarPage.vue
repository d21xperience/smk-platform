<template>
  <q-page padding>
    <div class="text-h5 q-mb-md">Surat Keluar</div>

    <q-tabs
      v-model="activeTab"
      inline-label
      class="text-grey"
      active-color="primary"
      indicator-color="primary"
      align="left"
    >
      <q-tab name="siswa" icon="school" label="Surat Siswa" />
      <q-tab name="ptk" icon="people" label="Guru & Staf" />
      <q-tab name="kedinasan" icon="business" label="Kedinasan" />
    </q-tabs>

    <q-tab-panels v-model="activeTab" animated>
      <q-tab-panel name="siswa">
        <SuratKeluarSiswaPanel />
      </q-tab-panel>
      <q-tab-panel name="ptk">
        <div class="text-grey-7 q-pa-md">Panel PTK (Akan diimplementasikan dengan pola yang sama)</div>
      </q-tab-panel>
      <q-tab-panel name="kedinasan">
        <div class="text-grey-7 q-pa-md">Panel Kedinasan (Akan diimplementasikan dengan pola yang sama)</div>
      </q-tab-panel>
    </q-tab-panels>

    <!-- Tombol Buku Ekspedisi (Global) -->
    <q-page-sticky position="bottom-right" :offset="[18, 18]">
      <q-btn fab icon="local_shipping" color="accent" @click="dialogEkspedisi = true" />
    </q-page-sticky>

    <!-- Dialog Cetak (Global untuk semua panel) -->
    <SuratKeluarPrintDialog />
  </q-page>
</template>

<script setup>
import { onMounted, watch } from 'vue'
import { useOutgoingLetter } from '@/composables/useOutgoingLetter.js'

import SuratKeluarSiswaPanel from '@/components/correspondence/outgoing/panels/SuratKeluarSiswaPanel.vue'
import SuratKeluarPrintDialog from '@/components/correspondence/outgoing/SuratKeluarPrintDialog.vue'

const {
  activeTab,
  dialogEkspedisi,
  loadList,
  loadAllLists,
  loadOptions
} = useOutgoingLetter()

watch(activeTab, (newTab) => {
  loadList(newTab)
})

onMounted(() => {
  loadAllLists()
  loadOptions()
})
</script>
