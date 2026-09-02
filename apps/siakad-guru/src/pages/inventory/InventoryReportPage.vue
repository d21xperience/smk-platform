<!--
FILE: src/pages/inventory/InventoryReportPage.vue
STATUS: MODIFY
STATUS IMPLEMENTASI: COMPLETE
-->

<template>
  <q-page padding>
    <div class="row items-center q-mb-md">
      <q-btn flat dense icon="arrow_back" label="Kembali" :to="{ name: 'inventory' }" class="q-mr-sm" />

      <div class="text-h5">Lapor Kondisi Aset</div>
    </div>

    <ReadOnlyContextBanner />

    <q-banner v-if="error" class="bg-negative text-white q-mb-md">
      {{ error }}
    </q-banner>

    <q-card v-if="isLoading" class="q-mb-md">
      <q-card-section class="text-center">
        <q-spinner color="primary" size="2em" />

        <div class="q-mt-sm">Loading...</div>
      </q-card-section>
    </q-card>

    <template v-else-if="selectedItem">
      <q-card class="q-mb-md">
        <q-card-section>
          <div class="text-h6">{{ selectedItem.name }}</div>

          <q-separator class="q-my-md" />

          <q-list dense>
            <q-item>
              <q-item-section>
                <q-item-label>Kategori</q-item-label>
              </q-item-section>

              <q-item-section side>
                <q-item-label caption>{{ selectedItem.category }}</q-item-label>
              </q-item-section>
            </q-item>

            <q-item>
              <q-item-section>
                <q-item-label>Lokasi</q-item-label>
              </q-item-section>

              <q-item-section side>
                <q-item-label caption>{{ selectedItem.location }}</q-item-label>
              </q-item-section>
            </q-item>

            <q-item>
              <q-item-section>
                <q-item-label>Kelas</q-item-label>
              </q-item-section>

              <q-item-section side>
                <q-item-label caption>{{ selectedItem.className }}</q-item-label>
              </q-item-section>
            </q-item>

            <q-item>
              <q-item-section>
                <q-item-label>Metode Laporan</q-item-label>
              </q-item-section>

              <q-item-section side>
                <q-badge :color="selectedItem.hasQRCode ? 'positive' : 'warning'">
                  {{ selectedItem.hasQRCode ? 'QR Scan' : 'Manual' }}
                </q-badge>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>
      </q-card>

      <q-card class="q-mb-md">
        <q-card-section>
          <div class="text-h6 q-mb-md">Form Laporan</div>

          <q-select v-model="condition" :options="conditionOptions" label="Kondisi Aset" outlined dense emit-value
            map-options class="q-mb-md" :disable="isReadOnlyContext" />

          <q-input v-model="description" label="Deskripsi Kerusakan" type="textarea" outlined rows="3" class="q-mb-md"
            :disable="isReadOnlyContext" />

          <q-input v-model="reportDate" label="Tanggal Laporan" type="date" outlined dense
            :disable="isReadOnlyContext" />
        </q-card-section>
      </q-card>

      <q-card v-if="isReadOnlyContext">
        <q-card-section class="text-center text-grey">
          Pelaporan pada context historis tidak tersedia. Data hanya dapat dilihat.
        </q-card-section>
      </q-card>

      <q-card v-else>
        <q-card-section class="row q-gutter-sm">
          <q-btn color="positive" label="Kirim Laporan" :loading="isLoading" @click="handleSubmit" />
        </q-card-section>
      </q-card>
    </template>

    <q-card v-else>
      <q-card-section class="text-center text-grey">
        Aset tidak ditemukan.
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuasar } from 'quasar'
import { useInventory } from '../../composables/useInventory.js'
import { useContext } from '../../composables/useContext.js'
import { REPORT_METHOD } from '../../domain/inventory/models/InventoryCondition.js'
import ReadOnlyContextBanner from '../../components/common/ReadOnlyContextBanner.vue'

const route = useRoute()
const router = useRouter()
const $q = useQuasar()

const { isReadOnlyContext } = useContext()

const {
  selectedItem,
  isLoading,
  error,
  loadItemById,
  submitDamageReport,
} = useInventory()

const itemId = computed(() => route.params.itemId)

const condition = ref('')
const description = ref('')
const reportDate = ref(new Date().toISOString().split('T')[0])

const conditionOptions = [
  { label: 'Baik', value: 'good' },
  { label: 'Rusak Ringan', value: 'damaged' },
  { label: 'Rusak Berat', value: 'broken' },
  { label: 'Hilang', value: 'missing' },
]

onMounted(async () => {
  try {
    await loadItemById({ itemId: itemId.value })
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.message,
    })

    router.push({ name: 'inventory' })
  }
})

const handleSubmit = async () => {
  if (isReadOnlyContext.value) {
    $q.notify({
      type: 'warning',
      message: 'Data historis bersifat read-only',
    })

    return
  }

  if (!condition.value || !description.value) {
    $q.notify({
      type: 'warning',
      message: 'Lengkapi semua field terlebih dahulu',
    })

    return
  }

  try {
    const reportMethod = selectedItem.value.hasQRCode
      ? REPORT_METHOD.QR
      : REPORT_METHOD.MANUAL

    await submitDamageReport({
      itemId: selectedItem.value.id,
      itemName: selectedItem.value.name,
      condition: condition.value,
      description: description.value,
      reportMethod,
      date: reportDate.value,
      classId: selectedItem.value.classId,
      className: selectedItem.value.className,
    })

    $q.notify({
      type: 'positive',
      message: 'Laporan kerusakan berhasil dikirim',
    })

    router.push({ name: 'inventory' })
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.message,
    })
  }
}
</script>
