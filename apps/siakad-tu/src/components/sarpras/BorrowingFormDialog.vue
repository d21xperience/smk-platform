<template>
  <q-dialog v-model="show" persistent>
    <q-card style="width: 600px; max-width: 90vw;">
      <q-card-section class="bg-primary text-white">
        <div class="text-h6">Form Peminjaman Aset</div>
      </q-card-section>

      <q-card-section class="q-gutter-y-md">
        <q-select outlined dense v-model="form.inventory_item_id" :options="store.activeAssets" option-value="id" option-label="name"
          label="Pilih Aset *" emit-value map-options>
          <template v-slot:option="scope">
            <q-item v-bind="scope.itemProps">
              <q-item-section>
                <q-item-label>{{ scope.opt.name }}</q-item-label>
                <q-item-label caption>Kode: {{ scope.opt.id }} | Merk: {{ scope.opt.brand || '-' }}</q-item-label>
              </q-item-section>
            </q-item>
          </template>
        </q-select>

        <div class="row q-col-gutter-sm">
          <div class="col-6">
            <q-input outlined dense v-model="form.borrower_name" label="Nama Peminjam *" />
          </div>
          <div class="col-6">
            <q-select outlined dense v-model="form.purpose_type" 
              :options="[
                { label: 'Pembelajaran', value: 'PEMBELAJARAN' },
                { label: 'Rapat', value: 'RAPAT' },
                { label: 'Kegiatan Sekolah', value: 'KEGIATAN' },
                { label: 'Lomba / Kompetisi', value: 'LOMBA' },
                { label: 'Dinas Luar', value: 'DINAS' },
                { label: 'Lainnya', value: 'LAINNYA' }
              ]" 
              label="Tujuan Peminjaman *" emit-value map-options />
          </div>
        </div>

        <q-input outlined dense type="textarea" v-model="form.purpose_description" label="Deskripsi Tujuan *" />

        <div class="row q-col-gutter-sm">
          <div class="col-6">
            <q-input outlined dense v-model="form.borrow_date" type="date" label="Tanggal Pinjam *" />
          </div>
          <div class="col-6">
            <q-input outlined dense v-model="form.expected_return_date" type="date" label="Estimasi Tanggal Kembali *" />
          </div>
        </div>

        <q-input outlined dense v-model="form.usage_location" label="Lokasi Penggunaan *" />
      </q-card-section>

      <q-card-actions align="right" class="q-pa-md bg-grey-2">
        <q-btn flat label="Batal" v-close-popup />
        <q-btn color="primary" label="Ajukan Peminjaman" @click="submit" :loading="isSubmitting" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useBorrowingStore } from '@/stores/sarpras/borrowing'

const props = defineProps({ modelValue: Boolean })
const emit = defineEmits(['update:modelValue', 'saved'])

const $q = useQuasar()
const store = useBorrowingStore()

const show = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const isSubmitting = ref(false)
const form = ref({
  inventory_item_id: '',
  borrower_name: '',
  purpose_type: 'PEMBELAJARAN',
  purpose_description: '',
  borrow_date: new Date().toISOString().split('T')[0],
  expected_return_date: new Date(Date.now() + 86400000).toISOString().split('T')[0], // Besok
  usage_location: ''
})

onMounted(() => {
  store.fetchActiveAssets()
})

async function submit() {
  if (!form.value.inventory_item_id || !form.value.borrower_name || !form.value.purpose_description || !form.value.usage_location) {
    $q.notify({ type: 'negative', message: 'Mohon lengkapi semua field bertanda *' })
    return
  }
  
  if (form.value.expected_return_date < form.value.borrow_date) {
    $q.notify({ type: 'negative', message: 'Tanggal kembali tidak boleh sebelum tanggal pinjam' })
    return
  }

  isSubmitting.value = true
  try {
    await store.createBorrowing(form.value)
    $q.notify({ type: 'positive', message: 'Pengajuan peminjaman berhasil' })
    show.value = false
    emit('saved')
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Gagal mengajukan peminjaman' })
  } finally {
    isSubmitting.value = false
  }
}
</script>
