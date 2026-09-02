<template>
  <q-page padding>
    <q-tabs v-model="tab" class="text-grey" active-color="primary" indicator-color="primary" align="left">
      <q-tab name="kode" label="Kode Klasifikasi" />
      <q-tab name="rak" label="Lokasi Fisik" />
    </q-tabs>

    <q-tab-panels v-model="tab" animated>
      <!-- Panel Kode Klasifikasi -->
      <q-tab-panel name="kode">
        <div class="text-h6 q-mb-md">Struktur Kode Klasifikasi</div>
        <q-btn color="primary" icon="add" label="Tambah Kategori Utama" class="q-mb-md" @click="tambahRoot" />

        <q-tree :nodes="treeData" node-key="id" default-expand-all label-key="label" children-key="children">
          <template v-slot:default-header="prop">
            <div class="row items-center full-width">
              <div class="text-weight-bold q-mr-sm text-primary">{{ prop.node.code }}</div>
              <div>{{ prop.node.label }}</div>
              <q-space />
              <q-btn flat round dense icon="edit" size="sm" @click.stop="editNode(prop.node)" />
              <q-btn flat round dense icon="delete" color="negative" size="sm" @click.stop="hapusNode(prop.node)" />
            </div>
          </template>
        </q-tree>

        <q-dialog v-model="dialogNode" persistent>
          <q-card style="min-width: 400px">
            <q-card-section class="row items-center q-pb-none">
              <div class="text-h6">{{ isEditingNode ? 'Edit' : 'Tambah' }} Kode</div>
              <q-space />
              <q-btn icon="close" flat round dense v-close-popup />
            </q-card-section>
            <q-card-section>
              <q-form @submit.prevent="simpanNode">
                <q-input v-model="nodeForm.kode" label="Kode (misal 421.3)" required />
                <q-input v-model="nodeForm.label" label="Nama Klasifikasi" required />
                <div class="q-mt-md text-right">
                  <q-btn label="Batal" flat v-close-popup />
                  <q-btn label="Simpan" type="submit" color="primary" :loading="isLoading" />
                </div>
              </q-form>
            </q-card-section>
          </q-card>
        </q-dialog>
      </q-tab-panel>

      <!-- Panel Lokasi Fisik -->
      <q-tab-panel name="rak">
        <div class="text-h6 q-mb-md">Daftar Rak / Ordner</div>
        <q-btn color="primary" icon="add" label="Tambah Lokasi" class="q-mb-md" @click="tambahLokasi" />

        <q-table :rows="lokasiList" :columns="kolomLokasi" row-key="id" :loading="isLoading">
          <template v-slot:body-cell-aksi="props">
            <q-td :props="props">
              <q-btn flat icon="edit" size="sm" @click="editLokasi(props.row)" />
              <q-btn flat icon="delete" color="negative" size="sm" @click="hapusLokasi(props.row.id)" />
            </q-td>
          </template>
        </q-table>

        <q-dialog v-model="dialogRak" persistent>
          <q-card style="min-width: 400px">
            <q-card-section class="row items-center q-pb-none">
              <div class="text-h6">{{ isEditLokasi ? 'Edit' : 'Tambah' }} Lokasi</div>
              <q-space />
              <q-btn icon="close" flat round dense v-close-popup />
            </q-card-section>
            <q-card-section>
              <q-form @submit.prevent="simpanLokasi">
                <q-input v-model="lokasiForm.nama" label="Nama (misal: Laci 1, Ordner KUR-2026)" required />
                <q-input v-model="lokasiForm.deskripsi" label="Deskripsi" type="textarea" />
                <div class="q-mt-md text-right">
                  <q-btn label="Batal" flat v-close-popup />
                  <q-btn label="Simpan" type="submit" color="primary" :loading="isLoading" />
                </div>
              </q-form>
            </q-card-section>
          </q-card>
        </q-dialog>
      </q-tab-panel>
    </q-tab-panels>
  </q-page>
</template>

<script setup>
import { onMounted } from 'vue'
import { useClassification } from '@/composables/useClassification.js'

const {
  tab,
  treeData,
  lokasiList,
  isLoading,
  dialogNode,
  isEditingNode,
  nodeForm,
  dialogRak,
  isEditLokasi,
  lokasiForm,
  loadTree,
  loadLocations,
  tambahRoot,
  editNode,
  simpanNode,
  hapusNode,
  tambahLokasi,
  editLokasi,
  simpanLokasi,
  hapusLokasi
} = useClassification()

// Kolom tabel (UI specific)
const kolomLokasi = [
  { name: 'nama', label: 'Nama Lokasi', field: 'name', align: 'left' },
  { name: 'deskripsi', label: 'Deskripsi', field: 'description', align: 'left' },
  { name: 'aksi', label: 'Aksi', align: 'center' }
]

onMounted(() => {
  loadTree()
  loadLocations()
})
</script>
