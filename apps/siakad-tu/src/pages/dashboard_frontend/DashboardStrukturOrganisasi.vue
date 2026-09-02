<template>
  <q-page class="org-page q-pa-md">
    <!-- HEADER -->
    <div class="row items-center justify-between q-mb-md wrap q-gutter-sm no-print">
      <div>
        <div class="text-h6 text-weight-bold text-grey-9">Struktur Organisasi Sekolah</div>
        <div class="text-caption text-grey-6">{{ schoolName }} &middot; Tahun Ajaran {{ academicYear }}</div>
      </div>
      <q-btn no-caps outline color="primary" icon="print" label="Cetak" @click="handlePrint" />
    </div>

    <!-- KARTU BAGAN -->
    <q-card flat bordered class="q-pa-md q-mb-md">
      <OrgChart :data="organizationTree" :side-node="sideNode" @node-click="showNodeDetail" />
    </q-card>

    <!-- LEGENDA -->
    <q-card flat bordered class="q-pa-md">
      <div class="text-subtitle2 text-weight-bold q-mb-sm">Keterangan</div>
      <div class="row q-gutter-md">
        <div v-for="(meta, key) in ORG_NODE_TYPES" :key="key" class="row items-center q-gutter-xs">
          <div class="legend-swatch" :class="{ 'legend-swatch--dashed': key === 'koordinatif' }"
            :style="{ backgroundColor: meta.fill, borderColor: meta.border }" />
          <span class="text-caption text-grey-8">{{ meta.label }}</span>
        </div>
      </div>
      <div class="text-caption text-grey-6 q-mt-sm">
        Garis putus-putus menandakan hubungan koordinatif/kemitraan, bukan garis komando langsung.
      </div>
    </q-card>

    <!-- DIALOG DETAIL SIMPUL -->
    <q-dialog v-model="detailDialogOpen">
      <q-card style="min-width: 300px; max-width: 400px" class="full-width">
        <q-card-section class="row items-center justify-between">
          <div class="text-subtitle1 text-weight-bold">Detail Jabatan</div>
          <q-btn flat round dense icon="close" v-close-popup />
        </q-card-section>
        <q-separator />
        <q-card-section v-if="selectedNode">
          <div class="text-weight-bold">{{ selectedNode.jabatan }}</div>
          <div class="text-body2 text-grey-8 q-mt-xs">{{ selectedNode.nama || 'Belum diisi' }}</div>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import OrgChart from '@/components/OrgChart.vue'
import { organizationTree, sideNode, ORG_NODE_TYPES } from '@/data/organizationData.js'

// TODO: ganti dengan data profil sekolah dari GET /api/sekolah/profile (Go backend)
const schoolName = ref('SMK Pasundan Jatinangor')
const academicYear = ref('2025/2026')

const detailDialogOpen = ref(false)
const selectedNode = ref(null)

function showNodeDetail(node) {
  selectedNode.value = node
  detailDialogOpen.value = true
}

function handlePrint() {
  window.print()
}
</script>

<style scoped>
.org-page {
  max-width: 1300px;
  margin: 0 auto;
}

.legend-swatch {
  width: 16px;
  height: 16px;
  border-radius: 3px;
  border: 2px solid;
}

.legend-swatch--dashed {
  border-style: dashed;
}

@media print {
  .no-print {
    display: none !important;
  }
}
</style>
