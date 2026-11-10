<template>
  <q-page class="q-pa-md">
    <div class="row q-col-gutter-md q-mb-md">
      <div class="col-12 col-md-4">
        <q-card class="bg-blue-1 text-blue-10">
          <q-card-section class="row items-center justify-between">
            <div>
              <div class="text-subtitle2">Total Jenis Logistik</div>
              <div class="text-h5 text-weight-bold">{{ totalItems }} Item</div>
            </div>
            <q-icon name="inventory_2" size="2rem" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-md-4">
        <q-card class="bg-red-1 text-red-10">
          <q-card-section class="row items-center justify-between">
            <div>
              <div class="text-subtitle2">Stok Kritis</div>
              <div class="text-h5 text-weight-bold">{{ criticalItems }} Item</div>
            </div>
            <q-icon name="warning" size="2rem" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-md-4">
        <q-card class="bg-green-1 text-green-10">
          <q-card-section class="row items-center justify-between">
            <div>
              <div class="text-subtitle2">Kondisi Stok Aman</div>
              <div class="text-h5 text-weight-bold">{{ safeItems }} Item</div>
            </div>
            <q-icon name="check_circle" size="2rem" />
          </q-card-section>
        </q-card>
      </div>
    </div>

    <q-card>
      <q-card-section class="row items-center justify-between">
        <div class="text-h6">Katalog Barang Habis Pakai</div>
        <div class="row q-gutter-sm">
          <q-input v-model="search" dense outlined placeholder="Cari ATK/Logistik..." style="width: 250px">
            <template v-slot:append><q-icon name="search" /></template>
          </q-input>
          <q-btn color="primary" icon="add" label="TAMBAH ITEM BARU" @click="showAddDialog = true" no-caps />
        </div>
      </q-card-section>

      <q-table
        :rows="filteredItems"
        :columns="columns"
        row-key="id"
        :loading="store.isLoading"
        flat
        bordered
      >
        <template v-slot:body-cell-stok="props">
          <q-td :props="props">
            <div class="row items-center q-gutter-sm">
              <span class="text-weight-bold">{{ props.row.stok }} {{ props.row.satuan }}</span>
              <q-badge :color="props.row.stok <= props.row.stok_minimum ? 'red' : 'green'">
                {{ props.row.stok <= props.row.stok_minimum ? 'Kritis' : 'Aman' }}
              </q-badge>
            </div>
          </q-td>
        </template>
        
        <template v-slot:body-cell-aksi="props">
          <q-td :props="props" class="q-gutter-x-sm">
            <q-btn color="green" label="RESTOCK" size="sm" @click="openRestockDialog(props.row)" no-caps />
            <q-btn color="red" label="GUNAKAN" size="sm" @click="openUseDialog(props.row)" no-caps />
          </q-td>
        </template>
      </q-table>
    </q-card>

    <q-dialog v-model="showAddDialog" persistent>
      <q-card style="min-width: 400px">
        <q-card-section class="text-h6">Tambah Item Baru</q-card-section>
        <q-card-section class="q-gutter-y-md">
          <q-input v-model="formAdd.kode" label="ID / Kode Logistik (Cth: LOG001)" outlined dense />
          <q-input v-model="formAdd.nama" label="Nama Barang Habis Pakai" outlined dense />
          <q-input v-model="formAdd.kategori" label="Kategori (Cth: ATK, Tinta)" outlined dense />
          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <q-input v-model="formAdd.satuan" label="Satuan (Cth: Rim, Pcs)" outlined dense />
            </div>
            <div class="col-6">
              <q-input v-model.number="formAdd.stok_minimum" type="number" label="Batas Stok Kritis" outlined dense />
            </div>
          </div>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="primary" label="Simpan" @click="submitAdd" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <q-dialog v-model="showMutateDialog">
      <q-card style="min-width: 350px">
        <q-card-section class="text-h6">
          {{ mutationType === 'IN' ? 'Restock Barang' : 'Gunakan Barang' }} : {{ selectedItem?.nama }}
        </q-card-section>
        <q-card-section class="q-gutter-y-md">
          <q-input v-model.number="formMutate.jumlah" type="number" :label="`Jumlah ${mutationType === 'IN' ? 'Masuk' : 'Keluar'}`" outlined dense />
          <q-input v-model="formMutate.keterangan" label="Keterangan (Opsional)" outlined dense />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn :color="mutationType === 'IN' ? 'green' : 'red'" label="Konfirmasi" @click="submitMutation" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useQuasar } from 'quasar';
import { useLogisticStore } from '@/stores/sarpras/logistic'
import { useContextStore } from '@/stores/contextStore';

const $q = useQuasar();
const store = useLogisticStore()
const contextStore = useContextStore();

const search = ref('');
const showAddDialog = ref(false);
const showMutateDialog = ref(false);

const mutationType = ref('IN');
const selectedItem = ref(null);

const formAdd = ref({
  kode: '',
  nama: '',
  kategori: '',
  satuan: '',
  stok_minimum: 5
});

const formMutate = ref({
  jumlah: 1,
  keterangan: ''
});

const columns = [
  { name: 'kode', label: 'ID Logistik', field: 'kode', align: 'left', sortable: true },
  { name: 'nama', label: 'Nama Barang Habis Pakai', field: 'nama', align: 'left', sortable: true },
  { name: 'stok', label: 'Sisa Stok Saat Ini', field: 'stok', align: 'left' },
  { name: 'aksi', label: 'Aset Kontrol Cepat', field: 'id', align: 'center' },
];

const filteredItems = computed(() => {
  if (!search.value) return store.items;
  return store.items.filter(item => 
    item.nama.toLowerCase().includes(search.value.toLowerCase()) || 
    item.kode.toLowerCase().includes(search.value.toLowerCase())
  );
});

const totalItems = computed(() => store.items.length);
const criticalItems = computed(() => store.items.filter(i => i.stok <= i.stok_minimum).length);
const safeItems = computed(() => store.items.filter(i => i.stok > i.stok_minimum).length);

onMounted(() => {
  store.fetchItems();
});

watch(() => contextStore.selectedSemesterId, (newVal) => {
  if (newVal) store.fetchItems()
});

const submitAdd = async () => {
  try {
    await store.createItem(formAdd.value);
    $q.notify({ type: 'positive', message: 'Item berhasil ditambahkan' });
    showAddDialog.value = false;
    formAdd.value = { kode: '', nama: '', kategori: '', satuan: '', stok_minimum: 5 };
  } catch (error) {
    $q.notify({ type: 'negative', message: error.response?.data?.error || 'Terjadi kesalahan' });
  }
};

const openRestockDialog = (item) => {
  selectedItem.value = item;
  mutationType.value = 'IN';
  formMutate.value = { jumlah: 1, keterangan: '' };
  showMutateDialog.value = true;
};

const openUseDialog = (item) => {
  selectedItem.value = item;
  mutationType.value = 'OUT';
  formMutate.value = { jumlah: 1, keterangan: '' };
  showMutateDialog.value = true;
};

const submitMutation = async () => {
  try {
    if (mutationType.value === 'IN') {
      await store.restockItem(selectedItem.value.id, formMutate.value);
      $q.notify({ type: 'positive', message: 'Restock berhasil' });
    } else {
      await store.useItem(selectedItem.value.id, formMutate.value);
      $q.notify({ type: 'positive', message: 'Penggunaan berhasil dicatat' });
    }
    showMutateDialog.value = false;
  } catch (error) {
    $q.notify({ type: 'negative', message: error.response?.data?.error || 'Terjadi kesalahan' });
  }
};
</script>
