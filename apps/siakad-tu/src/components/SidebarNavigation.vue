<template>
  <q-drawer v-model="drawerModel" show-if-above bordered :mini="miniState" @mouseover="miniState = false"
    @mouseout="miniState = true" mini-to-overlay :width="280" :breakpoint="500" class="bg-grey-1">
    <q-scroll-area class="fit">
      <!-- Header Sidebar / Profil Pengguna Singkat -->
      <!-- <div class="q-pa-sm q-mb-md flex flex-center" v-if="!miniState">
        <div class="text-center">
          <q-avatar size="56px" class="q-mb-sm bg-primary text-white">
            <q-icon name="admin_panel_settings" size="32px" />
          </q-avatar>
          <div class="text-weight-bold text-subtitle2 text-primary">TU KESISWAAN</div>
          <div class="text-caption text-grey-7">SMK Negeri Indonesia</div>
        </div>
      </div> -->

      <!-- Daftar Menu Navigasi -->
      <q-list padding class="menu-list">
        <template v-for="(menu, index) in menuList" :key="index">

          <!-- Jika menu memiliki sub-menu (Expansion Item) -->
          <q-expansion-item v-if="menu.children" :icon="menu.icon" :label="menu.label"
            :header-class="isActiveGroup(menu.children) ? 'text-primary text-weight-bold' : 'text-grey-8'"
            default-opened>
            <q-item v-for="(sub, subIndex) in menu.children" :key="subIndex" clickable v-ripple :to="sub.to" exact
              active-class="bg-blue-1 text-primary text-weight-bold" class="q-pl-lg">
              <q-item-section avatar>
                <q-icon :name="sub.icon" size="20px" />
              </q-item-section>
              <q-item-section>
                {{ sub.label }}
              </q-item-section>
            </q-item>
          </q-expansion-item>

          <!-- Jika menu tunggal (Single Item) -->
          <q-item v-else clickable v-ripple :to="menu.to" exact active-class="bg-blue-1 text-primary text-weight-bold"
            class="text-grey-8">
            <q-item-section avatar>
              <q-icon :name="menu.icon" />
            </q-item-section>
            <q-item-section>
              {{ menu.label }}
            </q-item-section>
          </q-item>

        </template>
      </q-list>
    </q-scroll-area>
  </q-drawer>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'

// Definisi Props dan Emits standar di <script setup>
const props = defineProps({
  modelValue: {
    type: Boolean,
    required: true
  }
})

const emit = defineEmits(['update:modelValue'])

const route = useRoute()

// Sinkronisasi state drawer menggunakan computed setter/getter
const drawerModel = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

// State untuk efek mini/expand sewaktu di-hover
const miniState = ref(false)

// Data Menu Kesiswaan SMK
const menuList = [
  {
    label: 'Dashboard',
    icon: 'dashboard',
    to: { name: 'kesiswaan-dashboard' }
  },
  {
    label: 'PSB (Siswa Baru)',
    icon: 'person_add_alt_1',
    to: { name: 'psb-siswa-baru' }
  },
  {
    label: 'Data Induk Siswa',
    icon: 'face',
    children: [
      { label: 'Direktori Siswa', icon: 'groups', to: { name: 'direktori-siswa' } },
      { label: 'Buku Induk & Klaper', icon: 'menu_book', to: { name: 'buku-induk-siswa' } },
      { label: 'Manajemen Mutasi', icon: 'compare_arrows', to: { name: 'manajemen-mutasi-siswa' } }
    ]
  },
  {
    label: 'Hubin & PKL',
    icon: 'business_center',
    children: [
      { label: 'Plotting & Berkas PKL', icon: 'assignment_ind', to: { name: 'ploting-berkas-siswa' } },
      { label: 'Sertifikasi & UKK', icon: 'workspace_premium', to: { name: 'sertifikat-ukk-siswa' } },
      { label: 'Asuransi Kerja', icon: 'gavel', to: { name: 'asuransi-kerja-siswa' } }
    ]
  },
  {
    label: 'Layanan Persuratan',
    icon: 'mail',
    children: [
      { label: 'Surat Keterangan Aktif', icon: 'description', to: { name: 'surat-aktif-siswa' } },
      { label: 'Beasiswa & PIP', icon: 'payments', to: { name: 'beasiswa-pip-siswa' } },
      { label: 'Legalisir & Alumni', icon: 'history_edu', to: { name: 'legalisir-alumni' } }
    ]
  },
  {
    label: 'Bursa Kerja (BKK)',
    icon: 'work',
    children: [
      { label: 'Tracer Study (BMW)', icon: 'analytics', to: { name: 'tracer-study-alumni' } },
      { label: 'Lowongan Kerja', icon: 'campaign', to: { name: 'loker' } }
    ]
  },
  {
    label: 'Rekap & Pelaporan',
    icon: 'summarize',
    to: { name: 'rekap-pelaporan' }
  }
]

// Fungsi mengecek sub-menu aktif agar grup menu otomatis berubah warna
const isActiveGroup = (children) => {
  return children.some(child => route.path === child.to)
}
</script>

<style scoped></style>
