<template>
  <q-dialog v-model="show" persistent>
    <q-card style="width: 500px; max-width: 90vw;">
      <q-card-section class="bg-secondary text-white">
        <div class="text-h6">{{ title }}</div>
      </q-card-section>
      <q-card-section class="q-gutter-y-md">
        <!-- FORM VERIFIKASI KTU -->
        <template v-if="type === 'verify'">
          <q-input outlined dense v-model="form.verification_team" label="Tim Verifikasi (Nama-nama) *" />
          <q-select outlined dense v-model="form.team_recommendation" :options="['SETUJU', 'TOLAK', 'REVISI']" label="Rekomendasi Tim *" />
          <q-input outlined dense type="textarea" v-model="form.team_notes" label="Catatan Tim" />
        </template>

        <!-- FORM APPROVAL KEPSEK/DINAS -->
        <template v-if="type === 'approve'">
          <div class="text-subtitle2">Approval Step: {{ request.current_approval_step + 1 }} / {{ request.approval_level }}</div>
          <q-select outlined dense v-model="form.decision" :options="['APPROVED', 'REJECTED']" label="Keputusan *" />
          <q-input outlined dense type="textarea" v-model="form.notes" label="Catatan (Wajib jika Reject)" />
        </template>

        <!-- FORM EKSEKUSI SARPRAS -->
        <template v-if="type === 'execute'">
          <q-select outlined dense v-model="form.execution_method" :options="['DIMUSNAHKAN', 'DIJUAL', 'DISUMBANGKAN', 'DITUKAR', 'DIMANFAATKAN']" label="Metode Eksekusi *" />
          <q-input outlined dense v-model="form.execution_location" label="Lokasi Eksekusi *" />
          <q-input outlined dense type="textarea" v-model="form.execution_notes" label="Catatan Eksekusi" />
          <q-input outlined dense type="number" v-model.number="form.sale_proceeds" label="Hasil Penjualan (Jika Dijual)" />
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
import { useDisposalStore } from '@/stores/sarpras/disposal'

const props = defineProps({
  modelValue: Boolean,
  request: Object,
  type: String // 'verify', 'approve', 'execute'
})
const emit = defineEmits(['update:modelValue', 'saved'])

const $q = useQuasar()
const store = useDisposalStore()

const show = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const isSubmitting = ref(false)
const form = ref({})

const title = computed(() => {
  if (props.type === 'verify') return 'Verifikasi Usulan (KTU)'
  if (props.type === 'approve') return 'Approval Usulan'
  if (props.type === 'execute') return 'Eksekusi Penghapusan Aset'
  return 'Aksi'
})

const submit = async () => {
  isSubmitting.value = true
  try {
    if (props.type === 'verify') {
      await store.verifyRequest(props.request.id, form.value)
    } else if (props.type === 'approve') {
      const payload = { ...form.value, approval_step: props.request.current_approval_step + 1 }
      await store.approveRequest(props.request.id, payload)
    } else if (props.type === 'execute') {
      await store.executeRequest(props.request.id, form.value)
    }
    $q.notify({ type: 'positive', message: 'Aksi berhasil' })
    show.value = false
    emit('saved')
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Gagal memproses' })
  } finally {
    isSubmitting.value = false
  }
}
</script>
eof
