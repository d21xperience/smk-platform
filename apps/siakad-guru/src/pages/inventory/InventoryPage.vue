<!--
FILE: src/pages/inventory/InventoryPage.vue
STATUS: MODIFY
STATUS IMPLEMENTASI: COMPLETE
-->

<template>
  <q-page padding>
    <ReadOnlyContextBanner />

    <div class="row items-center justify-between q-mb-md">
      <div class="text-h5">Inventory Reporting</div>

      <q-btn v-if="!isReadOnlyContext" color="primary" label="Scan QR" icon="qr_code_scanner"
        @click="showQRDialog = true" />

      <q-badge v-else color="grey">
        Read-only
      </q-badge>
    </div>

    <q-banner v-if="error" class="bg-negative text-white q-mb-md">
      {{ error }}
    </q-banner>

    <q-card class="q-mb-md">
      <q-card-section>
        <q-select v-model="selectedClass" :options="classOptions" label="Pilih Kelas" outlined dense emit-value
          map-options @update:model-value="handleClassChange" />
      </q-card-section>
    </q-card>

    <q-card v-if="isLoading" class="q-mb-md">
      <q-card-section class="text-center">
        <q-spinner color="primary" size="2em" />

        <div class="q-mt-sm">Loading...</div>
      </q-card-section>
    </q-card>

    <template v-else>
      <q-card v-if="items.length === 0">
        <q-card-section class="text-center text-grey">
          Pilih kelas untuk melihat daftar aset.
        </q-card-section>
      </q-card>

      <q-list v-else separator>
        <q-item v-for="item in items" :key="item.id" clickable v-ripple :disable="isReadOnlyContext"
          @click="handleManualReport(item)">
          <q-item-section avatar>
            <q-icon :name="item.hasQRCode ? 'qr_code' : 'category'" :color="item.hasQRCode ? 'primary' : 'grey'"
              size="2em" />
          </q-item-section>

          <q-item-section>
            <q-item-label>{{ item.name }}</q-item-label>
            <q-item-label caption>
              {{ item.category }} | {{ item.location }}
            </q-item-label>
          </q-item-section>

          <q-item-section side>
            <q-badge :color="item.hasQRCode ? 'positive' : 'warning'">
              {{ item.hasQRCode ? 'QR' : 'Manual' }}
            </q-badge>
          </q-item-section>
        </q-item>
      </q-list>
    </template>

    <q-dialog v-model="showQRDialog">
      <q-card style="min-width: 400px">
        <q-card-section>
          <div class="text-h6">Scan QR Code</div>
        </q-card-section>

        <q-card-section>
          <q-input v-model="qrInput" label="Masukkan Item ID (simulasi scan)" outlined dense />

          <div class="text-caption text-grey q-mt-sm">
            Contoh: INV-001, INV-002, INV-003
          </div>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Batal" color="primary" v-close-popup />

          <q-btn flat label="Scan" color="positive" @click="handleQRScan" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import { useQuasar } from 'quasar'
import { useInventory } from '../../composables/useInventory.js'
import { useContext } from '../../composables/useContext.js'
import ReadOnlyContextBanner from '../../components/common/ReadOnlyContextBanner.vue'

const $q = useQuasar()

const { isReadOnlyContext } = useContext()

const {
  items,
  isLoading,
  error,
  loadItems,
  navigateToReport,
} = useInventory()

const selectedClass = ref('')
const qrInput = ref('')
const showQRDialog = ref(false)

const classOptions = [
  { label: 'X-A', value: 'CLS-X-A' },
  { label: 'XI-B', value: 'CLS-XI-B' },
  { label: 'XII-C', value: 'CLS-XII-C' },
]

const handleClassChange = async (classId) => {
  await loadItems({ classId })
}

const handleQRScan = async () => {
  if (isReadOnlyContext.value) {
    $q.notify({
      type: 'warning',
      message: 'Data historis bersifat read-only',
    })

    return
  }

  if (!qrInput.value) {
    $q.notify({
      type: 'warning',
      message: 'Masukkan Item ID terlebih dahulu',
    })

    return
  }

  navigateToReport(qrInput.value)

  showQRDialog.value = false
  qrInput.value = ''
}

const handleManualReport = (item) => {
  if (isReadOnlyContext.value) {
    $q.notify({
      type: 'warning',
      message: 'Data historis bersifat read-only',
    })

    return
  }

  navigateToReport(item.id)
}
</script>
