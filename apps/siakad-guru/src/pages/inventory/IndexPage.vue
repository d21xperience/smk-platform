<template>
  <div class="q-pa-md">
    <h6>Inventaris Sekolah</h6>

    <!-- Stats Cards -->
    <div class="row q-col-gutter-md q-mb-md">
      <div class="col-md-3 col-sm-6 col-xs-12">
        <q-card flat bordered>
          <q-card-section>
            <div class="text-grey-6 text-caption">Total Aset</div>
            <div class="text-h5">{{ assets.length }}</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-md-3 col-sm-6 col-xs-12">
        <q-card flat bordered>
          <q-card-section>
            <div class="text-grey-6 text-caption">Tersedia</div>
            <div class="text-h5 text-positive">{{ availableAssets.length }}</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-md-3 col-sm-6 col-xs-12">
        <q-card flat bordered>
          <q-card-section>
            <div class="text-grey-6 text-caption">Dipinjam</div>
            <div class="text-h5 text-orange">{{ borrowedAssets.length }}</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-md-3 col-sm-6 col-xs-12">
        <q-card flat bordered>
          <q-card-section>
            <div class="text-grey-6 text-caption">Rusak</div>
            <div class="text-h5 text-negative">{{ damagedAssets.length }}</div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Filter & Actions -->
    <div class="row q-col-gutter-md q-mb-md">
      <div class="col-md-4 col-sm-6 col-xs-12">
        <q-input v-model="searchQuery" label="Cari Aset..." outlined dense />
      </div>
      <div class="col-md-3 col-sm-6 col-xs-12">
        <q-select v-model="filterCategory" :options="categoryOptions" label="Kategori" outlined dense clearable />
      </div>
      <div class="col-md-3 col-sm-6 col-xs-12">
        <q-select v-model="filterStatus" :options="statusOptions" label="Status" outlined dense clearable />
      </div>
      <div class="col-md-2 col-sm-6 col-xs-12">
        <q-btn color="primary" label="Refresh" @click="loadData" icon="refresh" />
      </div>
    </div>

    <!-- Table Assets -->
    <q-table
      :rows="filteredAssets"
      :columns="columns"
      row-key="id"
      :loading="isLoading"
      flat
      bordered
      v-model:pagination="pagination"
    >
      <template v-slot:body-cell-status="props">
        <q-td>
          <q-badge :color="getStatusColor(props.row.status)">
            {{ props.row.status }}
          </q-badge>
        </q-td>
      </template>
      <template v-slot:body-cell-action="props">
        <q-td>
          <div class="q-gutter-xs">
            <!-- Borrow -->
            <q-btn
              v-if="props.row.isAvailable()"
              color="positive"
              icon="assignment"
              size="sm"
              @click="openBorrowDialog(props.row)"
            />
            <!-- Return -->
            <q-btn
              v-else-if="props.row.isBorrowed()"
              color="primary"
              icon="assignment_return"
              size="sm"
              @click="openReturnDialog(props.row)"
            />
            <!-- Damage Report -->
            <q-btn
              color="negative"
              icon="report_problem"
              size="sm"
              @click="openDamageDialog(props.row)"
              :disable="props.row.isDamaged()"
            />
          </div>
        </q-td>
      </template>
    </q-table>

    <!-- Borrow Dialog -->
    <q-dialog v-model="showBorrowDialog" persistent>
      <q-card style="min-width: 400px">
        <q-card-section>
          <div class="text-h6">Pinjam Aset</div>
          <div class="text-grey-6">{{ selectedAsset?.name }}</div>
        </q-card-section>
        <q-card-section>
          <q-form @submit.prevent="submitBorrow">
            <q-input v-model="borrowForm.borrowerName" label="Nama Peminjam" outlined dense />
            <q-input v-model="borrowForm.purpose" label="Tujuan Peminjaman" outlined dense />
            <q-input v-model="borrowForm.expectedReturnDate" label="Tanggal Kembali" type="date" outlined dense />
            <q-input v-model="borrowForm.notes" label="Catatan" outlined dense autogrow />
            <div v-if="borrowError" class="text-negative q-mt-sm">{{ borrowError }}</div>
            <q-card-actions align="right">
              <q-btn flat label="Batal" v-close-popup />
              <q-btn color="primary" label="Pinjam" type="submit" :loading="isLoading" />
            </q-card-actions>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Return Dialog -->
    <q-dialog v-model="showReturnDialog" persistent>
      <q-card style="min-width: 400px">
        <q-card-section>
          <div class="text-h6">Kembalikan Aset</div>
          <div class="text-grey-6">{{ selectedAsset?.name }}</div>
          <div>Dipinjam oleh: {{ selectedIssue?.borrowerName }}</div>
        </q-card-section>
        <q-card-section>
          <q-form @submit.prevent="submitReturn">
            <div v-if="returnError" class="text-negative">{{ returnError }}</div>
            <q-card-actions align="right">
              <q-btn flat label="Batal" v-close-popup />
              <q-btn color="primary" label="Kembalikan" type="submit" :loading="isLoading" />
            </q-card-actions>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Damage Dialog -->
    <q-dialog v-model="showDamageDialog" persistent>
      <q-card style="min-width: 400px">
        <q-card-section>
          <div class="text-h6">Laporkan Kerusakan</div>
          <div class="text-grey-6">{{ selectedAsset?.name }}</div>
        </q-card-section>
        <q-card-section>
          <q-form @submit.prevent="submitDamage">
            <q-input v-model="damageDescription" label="Deskripsi Kerusakan" outlined dense autogrow />
            <div v-if="damageError" class="text-negative">{{ damageError }}</div>
            <q-card-actions align="right">
              <q-btn flat label="Batal" v-close-popup />
              <q-btn color="negative" label="Laporkan" type="submit" :loading="isLoading" />
            </q-card-actions>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useInventoryStore } from '@/stores/inventory.store';
import { useQuasar } from 'quasar';
import { ASSET_STATUS } from '@/domain/inventory/models/InventoryItem';

const $q = useQuasar();
const store = useInventoryStore();

// State
const searchQuery = ref('');
const filterCategory = ref(null);
const filterStatus = ref(null);
const isLoading = computed(() => store.isLoading);
const assets = computed(() => store.assets);
const availableAssets = computed(() => store.availableAssets);
const borrowedAssets = computed(() => store.borrowedAssets);
const damagedAssets = computed(() => store.damagedAssets);

// Dialog states
const showBorrowDialog = ref(false);
const showReturnDialog = ref(false);
const showDamageDialog = ref(false);
const selectedAsset = ref(null);
const selectedIssue = ref(null);
const borrowForm = ref({ borrowerName: '', purpose: '', expectedReturnDate: '', notes: '' });
const borrowError = ref('');
const returnError = ref('');
const damageDescription = ref('');
const damageError = ref('');

// Filter options
const categoryOptions = ['Elektronik', 'Furniture', 'ATK', 'Lainnya'];
const statusOptions = Object.values(ASSET_STATUS);

// Table
const pagination = ref({ rowsPerPage: 10 });
const columns = [
  { name: 'code', label: 'Kode', field: 'code', align: 'left' },
  { name: 'name', label: 'Nama Aset', field: 'name', align: 'left' },
  { name: 'category', label: 'Kategori', field: 'category', align: 'left' },
  { name: 'room', label: 'Ruangan', field: 'room', align: 'left' },
  { name: 'status', label: 'Status', field: 'status', align: 'center' },
  { name: 'action', label: 'Aksi', field: 'action', align: 'center' },
];

const filteredAssets = computed(() => {
  let result = assets.value;
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase();
    result = result.filter(a => a.name.toLowerCase().includes(q) || a.code.toLowerCase().includes(q));
  }
  if (filterCategory.value) result = result.filter(a => a.category === filterCategory.value);
  if (filterStatus.value) result = result.filter(a => a.status === filterStatus.value);
  return result;
});

const getStatusColor = (status) => {
  const map = {
    [ASSET_STATUS.AVAILABLE]: 'positive',
    [ASSET_STATUS.BORROWED]: 'orange',
    [ASSET_STATUS.DAMAGED]: 'negative',
    [ASSET_STATUS.MAINTENANCE]: 'blue',
  };
  return map[status] || 'grey';
};

// Load data
const loadData = async () => {
  try {
    await store.fetchAssets();
    await store.fetchIssues();
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Gagal memuat data: ' + err.message });
  }
};

// Borrow
const openBorrowDialog = (asset) => {
  selectedAsset.value = asset;
  borrowForm.value = { borrowerName: '', purpose: '', expectedReturnDate: '', notes: '' };
  borrowError.value = '';
  showBorrowDialog.value = true;
};

const submitBorrow = async () => {
  try {
    await store.borrowAsset(selectedAsset.value.id, borrowForm.value);
    showBorrowDialog.value = false;
    $q.notify({ type: 'positive', message: `${selectedAsset.value.name} berhasil dipinjam.` });
    await loadData();
  } catch (err) {
    borrowError.value = err.message;
  }
};

// Return
const openReturnDialog = async (asset) => {
  selectedAsset.value = asset;
  // Cari issue yang sedang dipinjam untuk aset ini
  const issue = store.issues.find(i => i.assetId === asset.id && i.isBorrowed());
  if (!issue) {
    $q.notify({ type: 'warning', message: 'Tidak ada transaksi peminjaman aktif.' });
    return;
  }
  selectedIssue.value = issue;
  returnError.value = '';
  showReturnDialog.value = true;
};

const submitReturn = async () => {
  try {
    await store.returnAsset(selectedIssue.value.id);
    showReturnDialog.value = false;
    $q.notify({ type: 'positive', message: `${selectedAsset.value.name} berhasil dikembalikan.` });
    await loadData();
  } catch (err) {
    returnError.value = err.message;
  }
};

// Damage
const openDamageDialog = (asset) => {
  selectedAsset.value = asset;
  damageDescription.value = '';
  damageError.value = '';
  showDamageDialog.value = true;
};

const submitDamage = async () => {
  try {
    await store.reportDamage(selectedAsset.value.id, damageDescription.value);
    showDamageDialog.value = false;
    $q.notify({ type: 'warning', message: `Kerusakan pada ${selectedAsset.value.name} dilaporkan.` });
    await loadData();
  } catch (err) {
    damageError.value = err.message;
  }
};

// Mount
onMounted(loadData);
</script>
