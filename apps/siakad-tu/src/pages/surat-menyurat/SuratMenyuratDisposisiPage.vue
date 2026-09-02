<template>
  <q-page padding>
    <!-- Role Switcher untuk Simulasi (Bisa dihapus di production jika auth sudah final) -->
    <div class="q-mb-md row items-center">
      <div class="text-subtitle2 q-mr-md">Simulasi Role:</div>
      <q-btn-toggle
        v-model="currentRole"
        toggle-color="primary"
        :options="[
          { label: 'Kepala Sekolah', value: 'kepsek' },
          { label: 'Waka/Guru', value: 'waka' }
        ]"
        @update:model-value="switchRole"
      />
    </div>

    <template v-if="currentRole === 'kepsek'">
      <!-- Kotak Masuk Disposisi Kepsek -->
      <div class="text-h5 q-mb-md">Kotak Masuk Disposisi</div>

      <q-card v-if="pendingIncomingLetters.length === 0" class="q-mb-md">
        <q-card-section class="text-center text-grey-6">
          Tidak ada surat masuk menunggu disposisi.
        </q-card-section>
      </q-card>

      <q-card v-for="surat in pendingIncomingLetters" :key="surat.id" class="q-mb-md">
        <q-card-section>
          <div class="text-h6">{{ surat.perihalSurat }}</div>
          <p><strong>Asal:</strong> {{ surat.asal }} | <strong>Tgl Diterima:</strong> {{ surat.tanggalDiterima }}</p>
          <q-btn label="Lihat & Disposisi" color="primary" @click="openDialogDisposisi(surat)" />
        </q-card-section>
      </q-card>

      <!-- Dialog Buat Disposisi -->
      <q-dialog v-model="dialogDisposisi" persistent>
        <q-card style="min-width: 500px">
          <q-card-section class="row items-center q-pb-none">
            <div class="text-h6">Buat Disposisi</div>
            <q-space />
            <q-btn icon="close" flat round dense v-close-popup />
          </q-card-section>
          <q-card-section>
            <q-form @submit.prevent="submitDisposisi">
              <q-select
                v-model="formDisposisi.tujuan"
                :options="stafOptions"
                label="Tujuan Staf"
                required
                emit-value
                map-options
              />
              <q-input v-model="formDisposisi.instruksi" label="Instruksi Tindak Lanjut" type="textarea" required />
              <div class="q-mt-md text-right">
                <q-btn label="Batal" flat v-close-popup />
                <q-btn label="Kirim Disposisi" type="submit" color="primary" :loading="isLoading" />
              </div>
            </q-form>
          </q-card-section>
        </q-card>
      </q-dialog>
    </template>

    <template v-if="currentRole === 'waka' || currentRole === 'guru'">
      <!-- Kotak Tugas Staf -->
      <div class="text-h5 q-mb-md">Instruksi Disposisi</div>

      <q-card v-if="myTasks.length === 0" class="q-mb-md">
        <q-card-section class="text-center text-grey-6">
          Tidak ada tugas disposisi.
        </q-card-section>
      </q-card>

      <q-card v-for="item in myTasks" :key="item.id" class="q-mb-md">
        <q-card-section>
          <div class="text-h6">{{ item.instruksi }}</div>
          <p>Dari: {{ item.dari }} | Tanggal: {{ item.tanggal }}</p>
          <q-badge v-if="item.status === 'SELESAI'" color="green" label="Selesai" class="q-mr-sm" />
          <q-btn v-else label="Tindak Lanjut" color="orange" @click="openDialogTindakLanjut(item)" />
        </q-card-section>
      </q-card>

      <!-- Dialog Tindak Lanjut -->
      <q-dialog v-model="dialogTindakLanjut" persistent>
        <q-card style="min-width: 500px">
          <q-card-section class="row items-center q-pb-none">
            <div class="text-h6">Laporan Tindak Lanjut</div>
            <q-space />
            <q-btn icon="close" flat round dense v-close-popup />
          </q-card-section>
          <q-card-section>
            <q-form @submit.prevent="submitTindakLanjut">
              <q-input v-model="formTindakLanjut.laporan" label="Laporan Singkat" type="textarea" required />
              <q-uploader
                v-model="formTindakLanjut.bukti"
                label="Unggah Bukti (foto/dokumen)"
                accept="image/*,.pdf"
                max-files="3"
                auto-upload
                :factory="uploadBukti"
                class="q-mt-md"
              />
              <div class="q-mt-md text-right">
                <q-btn label="Batal" flat v-close-popup />
                <q-btn label="Selesai" type="submit" color="green" :loading="isLoading" />
              </div>
            </q-form>
          </q-card-section>
        </q-card>
      </q-dialog>
    </template>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useDisposition } from '@/composables/useDisposition.js'
import { useOperationalContext, setOperationalContext } from '@/composables/useOperationalContext.js'

const {
  dialogDisposisi,
  dialogTindakLanjut,
  formDisposisi,
  formTindakLanjut,
  pendingIncomingLetters,
  myTasks,
  stafOptions,
  isLoading,
  loadKepsekInbox,
  loadStafTasks,
  loadStafOptions,
  openDialogDisposisi,
  submitDisposisi,
  openDialogTindakLanjut,
  submitTindakLanjut,
  uploadBukti
} = useDisposition()

const context = useOperationalContext()
const currentRole = ref(context.role)

// Simulasi perpindahan role untuk testing UI
function switchRole(role) {
  const userId = role === 'kepsek' ? 'USER-KEPSEK' : 'USER-WAKUR'
  setOperationalContext({ ...context, role: role, userId: userId })
  loadData()
}

function loadData() {
  if (currentRole.value === 'kepsek') {
    loadKepsekInbox()
  } else {
    loadStafTasks()
  }
}

onMounted(() => {
  loadStafOptions()
  loadData()
})
</script>
