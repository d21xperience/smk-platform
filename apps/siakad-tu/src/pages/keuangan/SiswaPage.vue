<template>
  <q-page class="q-pa-md">
    <div class="row justify-between items-center q-mb-md">
      <div class="text-h5">Data Siswa</div>
      <q-btn label="Tambah Siswa" color="primary" @click="openDialog" />
    </div>

    <q-table :rows="items" :columns="columns" row-key="id" :loading="loading" flat bordered>
      <template v-slot:body-cell-actions="props">
        <q-td auto-width>
          <q-btn flat round dense icon="edit" @click="editItem(props.row)" />
          <q-btn flat round dense icon="delete" color="negative" @click="confirmDelete(props.row)" />
        </q-td>
      </template>
    </q-table>

    <!-- Dialog form -->
    <q-dialog v-model="dialog">
      <q-card style="min-width: 500px">
        <q-card-section>
          <div class="text-h6">{{ isEdit ? 'Edit' : 'Tambah' }} Siswa</div>
        </q-card-section>
        <q-card-section>
          <q-input v-model="form.nis" label="NIS" />
          <q-input v-model="form.nama" label="Nama" />
          <q-select v-model="form.kelasId" :options="kelasOptions" label="Kelas" option-value="id" option-label="nama"
            emit-value map-options />
          <q-select v-model="form.tahunAjaranId" :options="tahunOptions" label="Tahun Ajaran" option-value="id"
            option-label="nama" emit-value map-options />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn flat label="Simpan" color="primary" @click="submitForm" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useSiswa } from '@/composables/useSiswa';
import { useKelas } from '@/composables/useKelas';
import { useTahunAjaran } from '@/composables/useTahunAjaran';
import { useQuasar } from 'quasar';

const $q = useQuasar();
const { items, loading, fetchAll, create, update, remove } = useSiswa();
const { items: kelasList, fetchAll: fetchKelas } = useKelas();
const { items: tahunList, fetchAll: fetchTahun } = useTahunAjaran();

const columns = [
  { name: 'id', label: 'ID', field: 'id', align: 'left' },
  { name: 'nis', label: 'NIS', field: 'nis' },
  { name: 'nama', label: 'Nama', field: 'nama' },
  { name: 'kelasId', label: 'Kelas', field: row => kelasList.find(k => k.id === row.kelasId)?.nama },
  { name: 'actions', label: 'Aksi', align: 'center' }
];

const dialog = ref(false);
const isEdit = ref(false);
const form = ref({ nis: '', nama: '', kelasId: null, tahunAjaranId: null });
const selectedId = ref(null);

const kelasOptions = ref([]);
const tahunOptions = ref([]);

onMounted(async () => {
  await fetchAll();
  await fetchKelas();
  await fetchTahun();
  kelasOptions.value = kelasList.value;
  tahunOptions.value = tahunList.value;
});

function openDialog() {
  isEdit.value = false;
  form.value = { nis: '', nama: '', kelasId: null, tahunAjaranId: null };
  dialog.value = true;
}

function editItem(row) {
  isEdit.value = true;
  selectedId.value = row.id;
  form.value = { ...row };
  dialog.value = true;
}

async function submitForm() {
  try {
    if (isEdit.value) {
      await update(selectedId.value, form.value);
      $q.notify({ type: 'positive', message: 'Data berhasil diupdate' });
    } else {
      await create(form.value);
      $q.notify({ type: 'positive', message: 'Data berhasil ditambahkan' });
    }
    dialog.value = false;
  } catch (err) {
    $q.notify({ type: 'negative', message: err.message });
  }
}

function confirmDelete(row) {
  $q.dialog({
    title: 'Konfirmasi',
    message: `Hapus siswa ${row.nama}?`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    await remove(row.id);
    $q.notify({ type: 'positive', message: 'Data dihapus' });
  });
}
</script>
