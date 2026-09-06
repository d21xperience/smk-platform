<template>
  <q-page class="q-pa-md">
    <q-inner-loading :showing="loading" />

    <div class="row justify-between q-mb-md">
      <div class="text-h5">Mutasi Siswa (Masuk/ Keluar)</div>
      <div class="q-gutter-sm">
        <q-btn color="green-7" icon="table_view" label="XLS" title="Export ke Excel" />
        <q-btn color="deep-orange-7" icon="print" label="PDF" title="Cetak PDF" />
        <q-btn color="primary" icon="add" label="Ajukan Mutasi" @click="showAddStudents = true" />

      </div>
    </div>

    <q-tabs v-model="tabStatus" dense class="q-mb-md" align="left" @update:model-value="loadMutationsByTab">
      <q-tab name="semua" label="Semua" />
      <q-tab name="pending" label="Masuk" />
      <q-tab name="disetujui" label="Keluar" />
      <!-- <q-tab name="ditolak" label="Ditolak" /> -->
    </q-tabs>

    <q-table :rows="mutations" :columns="columns" row-key="id" flat bordered dense :loading="loading">
      <template v-slot:body-cell-jenis="props">
        <q-td :props="props">
          <q-badge :color="props.row.jenis === 'MASUK' ? 'positive' : 'warning'">
            {{ props.row.jenis === 'MASUK' ? 'Mutasi Masuk' : 'Mutasi Keluar' }}
          </q-badge>
        </q-td>
      </template>
      <template v-slot:body-cell-status="props">
        <q-td :props="props">
          <q-badge :color="getStatusColor(props.row.status)">
            {{ props.row.status }}
          </q-badge>
        </q-td>
      </template>
      <template v-slot:body-cell-aksi="props">
        <q-td :props="props">
          <div v-if="props.row.status === 'PENDING' && canApprove">
            <q-btn flat round dense icon="check_circle" color="positive" @click="handleApprove(props.row.id)"
              class="q-mr-sm">
              <q-tooltip>Setujui</q-tooltip>
            </q-btn>
            <q-btn flat round dense icon="cancel" color="negative" @click="handleReject(props.row.id)">
              <q-tooltip>Tolak</q-tooltip>
            </q-btn>
          </div>
          <q-btn flat round dense icon="visibility" color="info" @click="lihatDetail(props.row)" class="q-mr-sm">
            <q-tooltip>Detail</q-tooltip>
          </q-btn>
          <q-btn flat round dense icon="delete" color="negative" @click="hapusMutasi(props.row.id)">
            <q-tooltip>Hapus</q-tooltip>
          </q-btn>
        </q-td>
      </template>
    </q-table>

    <!-- Dialog Form Pengajuan Mutasi (Disederhanakan untuk fokus pada arsitektur) -->
    <!-- <q-dialog v-model="formDialog" persistent>
      <q-card style="min-width: 600px">
        <q-card-section>
          <div class="text-h6">Ajukan Mutasi</div>
        </q-card-section>
        <q-card-section>
          <q-form @submit="submitMutasi" class="q-gutter-md">
            <q-select v-model="form.jenis"
              :options="[{ label: 'Mutasi Masuk (Dari Sekolah Lain)', value: 'MASUK' }, { label: 'Mutasi Keluar (Ke Sekolah Lain)', value: 'KELUAR' }]"
              label="Jenis Mutasi" outlined :rules="[v => !!v]" />

            <template v-if="form.jenis === 'MASUK'">
              <q-input v-model="form.namaCalon" label="Nama Lengkap" outlined :rules="[v => !!v]" />
              <q-input v-model="form.nisnCalon" label="NISN" outlined />
              <q-input v-model="form.asalSekolah" label="Asal Sekolah" outlined />
            </template>

            <template v-else>
              <q-select v-model="form.studentId" :options="siswaOptions" label="Pilih Siswa" option-value="id"
                option-label="nama_dan_kelas" outlined :rules="[v => !!v]" />
              <q-input v-model="form.tujuanSekolah" label="Tujuan Sekolah" outlined :rules="[v => !!v]" />
            </template>

            <q-input v-model="form.alasan" label="Alasan Mutasi" type="textarea" rows="3" outlined
              :rules="[v => !!v]" />

            <div class="q-gutter-sm">
              <q-btn label="Ajukan" type="submit" color="primary" :loading="loading" />
              <q-btn label="Batal" flat color="negative" v-close-popup @click="resetForm" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog> -->

    <!-- Dialog Catatan Penolakan -->
    <!-- <q-dialog v-model="catatanDialog" persistent>
      <q-card>
        <q-card-section>
          <div class="text-h6">Alasan Penolakan</div>
        </q-card-section>
        <q-card-section>
          <q-input v-model="catatanTolak" label="Catatan" type="textarea" outlined autofocus />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" color="negative" v-close-popup />
          <q-btn flat label="Tolak" color="primary" @click="executeReject" :loading="loading" />
        </q-card-actions>
      </q-card>
    </q-dialog> -->



    <!-- Add Students -->
    <q-dialog v-model="showAddStudents">
      <q-card style="min-width: 400px; max-width: 500px; border-radius: 12px;">

        <!-- Bagian Judul Dialog -->
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6 text-weight-bold text-primary">Mutasi Siswa</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <!-- Bagian Isi / Menu Pilihan -->
        <q-card-section class="q-pa-md">
          <q-list bordered separator class="rounded-borders text-grey-9">

            <!-- PILIHAN 1: MANUAL -->
            <q-item clickable v-close-popup @click="goToPage({ name: 'student-import' })">
              <q-item-section avatar>
                <q-icon name="person_add" color="primary" size="md" />
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold text-subtitle1">Mutasi Masuk</q-item-label>
                <q-item-label caption>Input siswa pindahan yang masuk</q-item-label>
              </q-item-section>
            </q-item>

            <!-- PILIHAN 2: IMPORT -->
            <q-item clickable v-close-popup @click="goToPage({ name: 'student-import' })">
              <q-item-section avatar>
                <q-icon name="person_remove" color="primary" size="md" />
              </q-item-section>
              <q-item-section>
                <q-item-label class="text-weight-bold text-subtitle1">Mutasi Keluar</q-item-label>
                <q-item-label caption>Input siswa yang Keluar</q-item-label>
              </q-item-section>
            </q-item>

          </q-list>
        </q-card-section>

      </q-card>
    </q-dialog>



  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useQuasar } from 'quasar';
import { useStudentAffairsStore } from '@/stores/kesiswaan/studentAffairsStore';

const $q = useQuasar();
const store = useStudentAffairsStore();
const showAddStudents = ref(false)
const tabStatus = ref('semua');
// const formDialog = ref(false);
const detailDialog = ref(false);
// const catatanDialog = ref(false);
// let currentRejectId = null;
// const catatanTolak = ref('');
const detail = ref({});

// const form = ref({
//   jenis: null, namaCalon: '', nisnCalon: '', asalSekolah: '',
//   studentId: null, tujuanSekolah: '', alasan: '', dokumen: null
// });

const columns = [
  { name: 'tanggalPengajuan', label: 'Tanggal', field: 'tanggalPengajuan', align: 'left' },
  { name: 'jenis', label: 'Jenis', field: 'jenis', align: 'center' },
  { name: 'nama_siswa', label: 'Nama Siswa', field: 'nama_siswa', align: 'left' },
  { name: 'status', label: 'Status', field: 'status', align: 'center' },
  { name: 'aksi', label: 'Aksi', field: 'aksi', align: 'center' }
];

const loading = computed(() => store.loading);
const mutations = computed(() => store.mutations.map(m => ({
  ...m,
  nama_siswa: m.jenis === 'MASUK' ? m.namaCalon : 'Siswa ID: ' + m.studentId // Simplifikasi untuk mock
})));

// const siswaOptions = computed(() => {
//   // Dalam implementasi nyata, ini diambil dari store siswa yang sudah di-filter 'AKTIF'
//   return [{ id: '1', nama_dan_kelas: 'Ahmad Fauzi (XII RPL 1)' }];
// });

const canApprove = computed(() => {
  // Simulasi role check. Dalam produksi, gunakan useAuthStore
  return true;
});

onMounted(() => {
  loadMutationsByTab('semua');
});

function loadMutationsByTab(status) {
  store.loadMutations(status);
}

function getStatusColor(status) {
  if (status === 'PENDING') return 'warning';
  if (status === 'DISETUJUI') return 'positive';
  return 'negative';
}

// function resetForm() {
//   form.value = { jenis: null, namaCalon: '', nisnCalon: '', asalSekolah: '', studentId: null, tujuanSekolah: '', alasan: '', dokumen: null };
// }

// function openFormDialog() {
//   resetForm();
//   formDialog.value = true;
// }

// async function submitMutasi() {
//   try {
//     await store.submitMutationAction(form.value); // Asumsi method ini ditambahkan ke store
//     $q.notify({ type: 'positive', message: 'Pengajuan mutasi dikirim' });
//     formDialog.value = false;
//     loadMutationsByTab(tabStatus.value);
//   } catch (err) {
//     $q.notify({ type: 'negative', message: err.message });
//   }
// }

function lihatDetail(mutasi) {
  detail.value = mutasi;
  detailDialog.value = true;
}

async function handleApprove(id) {
  $q.dialog({
    title: 'Konfirmasi',
    message: 'Setujui pengajuan mutasi ini? Sistem akan otomatis memperbarui data siswa terkait.',
    cancel: true
  }).onOk(async () => {
    try {
      // BUSINESS LOGIC DIPINDAHKAN KE STORE -> SERVICE -> ENGINE
      await store.approveMutationAction(id, 'Disetujui oleh Admin');
      $q.notify({ type: 'positive', message: 'Mutasi disetujui dan data siswa diperbarui' });
    } catch (err) {
      $q.notify({ type: 'negative', message: err.message });
    }
  });
}

// function handleReject(id) {
//   currentRejectId = id;
//   catatanTolak.value = '';
//   catatanDialog.value = true;
// }

// async function executeReject() {
//   if (!catatanTolak.value.trim()) {
//     $q.notify({ type: 'negative', message: 'Alasan penolakan harus diisi' });
//     return;
//   }
//   try {
//     await store.rejectMutationAction(currentRejectId, catatanTolak.value);
//     $q.notify({ type: 'negative', message: 'Mutasi ditolak' });
//     catatanDialog.value = false;
//     currentRejectId = null;
//   } catch (err) {
//     $q.notify({ type: 'negative', message: err.message });
//   }
// }

function hapusMutasi() {
  $q.dialog({
    title: 'Konfirmasi',
    message: 'Hapus pengajuan mutasi ini?',
    cancel: true
  }).onOk(() => {
    $q.notify({ type: 'positive', message: 'Pengajuan dihapus (Mock)' });
  });
}
</script>
