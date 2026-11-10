<template>
  <q-dialog v-model="show" persistent>
    <q-card style="width: 500px; max-width: 90vw;">
      <q-card-section class="bg-secondary text-white">
        <div class="text-h6">{{ title }}</div>
      </q-card-section>

      <q-card-section class="q-gutter-y-md">
        <!-- FORM APPROVAL -->
        <template v-if="type === 'approve'">
          <div class="text-subtitle2">Step Approval: {{ request.current_approval_step + 1 }} dari {{ request.approval_level }}</div>
          <q-select outlined dense v-model="form.decision" :options="['APPROVED', 'REJECTED']" label="Keputusan *" />
          <q-input outlined dense type="textarea" v-model="form.notes" label="Catatan (Wajib jika Reject)" />
        </template>

        <!-- FORM PICKUP (SERAH TERIMA) -->
        <template v-if="type === 'pickup'">
          <q-select outlined dense v-model="form.condition" :options="['BAIK', 'RUSAK_RINGAN']" label="Kondisi Saat Diserahkan *" />
          <q-input outlined dense type="textarea" v-model="form.notes" label="Catatan (Opsional)" />
          <!-- Placeholder Upload Foto (Mock) -->
          <q-banner class="bg-grey-2 text-caption rounded-borders">Fitur Upload Foto dapat ditambahkan di iterasi berikutnya.</q-banner>
        </template>

        <!-- FORM RETURN (PENGEMBALIAN) -->
        <template v-if="type === 'return'">
          <q-select outlined dense v-model="form.condition" :options="['BAIK', 'RUSAK_RINGAN', 'RUSAK_BERAT']" label="Kondisi Saat Dikembalikan *" />
          <q-input outlined dense type="textarea" v-model="form.notes" label="Catatan (Wajib jika Rusak)" />
          <q-banner class="bg-grey-2 text-caption rounded-borders">Fitur Upload Foto dapat ditambahkan di iterasi berikutnya.</q-banner>
        </template>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn flat label="Batal" v-close-popup />
        <q-btn color="primary" label="Konfirmasi" @click="submit" :loading="isSubmitting" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'
import { useBorrowingStore } from '@/stores/sarpras/borrowing'

const props = defineProps({
  modelValue: Boolean,
  request: Object,
  type: String // 'approve', 'pickup', 'return'
})
const emit = defineEmits(['update:modelValue', 'saved'])

const $q = useQuasar()
const store = useBorrowingStore()

const show = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const isSubmitting = ref(false)
const form = ref({
  decision: 'APPROVED',
  notes: '',
  condition: 'BAIK',
  photos: '[]' // Backend expects JSON string
})

const title = computed(() => {
  if (props.type === 'approve') return 'Approval Peminjaman'
  if (props.type === 'pickup') return 'Serah Terima Aset (Pickup)'
  if (props.type === 'return') return 'Pengembalian Aset'
  return 'Aksi'
})

async function submit() {
  isSubmitting.value = true
  try {
    if (props.type === 'approve') {
      const payload = {
        approval_step: props.request.current_approval_step + 1,
        decision: form.value.decision,
        notes: form.value.notes
      }
      await store.approveBorrowing(props.request.id, payload)
    } else if (props.type === 'pickup') {
      await store.pickupBorrowing(props.request.id, form.value)
    } else if (props.type === 'return') {
      await store.returnBorrowing(props.request.id, form.value)
    }
    $q.notify({ type: 'positive', message: 'Aksi berhasil diproses' })
    show.value = false
    emit('saved')
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Gagal memproses aksi' })
  } finally {
    isSubmitting.value = false
  }
}
</script>
