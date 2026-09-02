<template>
  <q-page class="q-pa-md">
    <div class="text-h5">Tagihan Siswa</div>
    <!-- Pencarian siswa -->
    <q-input v-model="keyword" label="Cari siswa (nama/NIS)" @keyup.enter="cariSiswa" />
    <q-list v-if="hasilSiswa.length">
      <q-item clickable v-for="s in hasilSiswa" :key="s.id" @click="pilihSiswa(s)">
        <q-item-section>{{ s.nama }} ({{ s.nis }}) - {{ s.kelas }}</q-item-section>
      </q-item>
    </q-list>
    <div v-if="siswaTerpilih">
      <q-banner class="bg-grey-2 q-my-md">Siswa: {{ siswaTerpilih.nama }}</q-banner>
      <!-- Tombol buat tagihan baru -->
      <q-btn label="Buat Tagihan Baru" color="primary" @click="openDialogTagihan" />
      <!-- Tabel tagihan -->
      <q-table :rows="tagihanStore.items" :columns="tagihanColumns" row-key="id" flat bordered>
        <template v-slot:body-cell-actions="props">
          <q-btn v-if="props.row.status === 'belum'" flat label="Bayar" color="positive"
            @click="bayarTagihan(props.row)" />
        </template>
      </q-table>
    </div>
    <!-- Dialog buat tagihan -->
    <q-dialog v-model="dialogTagihan">
      <q-card style="min-width:400px">
        <q-card-section>
          <div class="text-h6">Buat Tagihan</div>
          <q-select v-model="formTagihan.kategoriId" :options="kategoriOptions" label="Kategori" emit-value
            map-options />
          <q-input v-model.number="formTagihan.nominal" label="Nominal" type="number" />
          <q-input v-model="formTagihan.bulan" label="Bulan (opsional)" />
          <q-input v-model.number="formTagihan.tahun" label="Tahun" type="number" />
          <q-input v-model="formTagihan.tanggalJatuhTempo" label="Jatuh Tempo" type="date" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn flat label="Simpan" color="primary" @click="simpanTagihan" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>
<script setup>
// implementasi menggunakan pinia stores
</script>
