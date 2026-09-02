<template>
  <q-btn-dropdown flat stretch no-caps :label="props.label" content-class="bg-white text-black shadow-5"
    class="font-tab-match">
    <q-list :style="{ minWidth: props.minWidth }" separator>
      <template v-for="menu in props.items" :key="menu.title">

        <!-- KONDISI A: JIKA MENU ADALAH SEBUAH GRUP (Mempunyai Anak Menu) -->
        <q-expansion-item v-if="menu.children && menu.children.length > 0" :icon="menu.icon" :label="menu.title"
          header-class="text-weight-bold text-primary" expand-separator>
          <q-list class="bg-grey-1">
            <!-- Loop Anak Menu Menggunakan Helper Function untuk Deteksi Link -->
            <component :is="getLinkComponent(child)" v-for="child in menu.children" :key="child.title"
              v-bind="getLinkAttributes(child)" clickable v-ripple>
              <!-- FIX: Menyuntikkan focus helper milik Quasar untuk efek hover pada tag 'a' -->
              <div v-if="child.url" class="q-focus-helper"></div>

              <q-item-section avatar v-if="child.icon">
                <q-icon :name="child.icon" size="xs" color="grey-7" />
              </q-item-section>
              <q-item-section class="text-body2 text-grey-9">
                <div class="row items-center no-wrap">
                  <span>{{ child.title }}</span>
                  <!-- Visual Anchor: Icon kecil penanda jika link membuka tab baru -->
                  <q-icon v-if="child.url && child.target === '_blank'" name="open_in_new" size="12px" color="grey-5"
                    class="q-ml-xs" />
                </div>
              </q-item-section>
            </component>
          </q-list>
        </q-expansion-item>

        <!-- KONDISI B: JIKA MENU MANDIRI BIASA (Tanpa Anak Menu) -->
        <component v-else :is="getLinkComponent(menu)" v-bind="getLinkAttributes(menu)" clickable v-ripple>
          <!-- FIX: Menyuntikkan focus helper milik Quasar untuk efek hover pada tag 'a' -->
          <div v-if="menu.url" class="q-focus-helper"></div>

          <q-item-section avatar v-if="menu.icon">
            <q-icon :name="menu.icon" size="sm" color="primary" />
          </q-item-section>
          <q-item-section class="text-weight-medium">
            <div class="row items-center no-wrap">
              <span>{{ menu.title }}</span>
              <q-icon v-if="menu.url && menu.target === '_blank'" name="open_in_new" size="12px" color="grey-5"
                class="q-ml-xs" />
            </div>
          </q-item-section>
        </component>

      </template>
    </q-list>
  </q-btn-dropdown>
</template>

<script setup>
import { QItem } from 'quasar'

const props = defineProps({
  label: { type: String, default: 'Persuratan' },
  minWidth: { type: String, default: '240px' },
  items: {
    type: Array,
    default: () => [
      { title: 'Dashboard', to: '/persuratan/dashboard', icon: 'dashboard' },
      {
        title: 'Alur Surat',
        icon: 'mail',
        children: [
          { title: 'Surat Masuk', to: '/persuratan/surat-masuk', icon: 'call_received' },
          { title: 'Surat Keluar', to: '/persuratan/surat-keluar', icon: 'call_made' },
        ]
      },
      // CONTOH GRUP DENGAN LINK EKSTERNAL (Buka Tab Baru / Href)
      {
        title: 'Layanan Eksternal',
        icon: 'cloud_queue',
        children: [
          { title: 'E-Open Perpustakaan', url: 'https://smk.sch.id', target: '_blank', icon: 'menu_book' },
          { title: 'Panduan Administrasi', url: '/files/panduan_tu.pdf', target: '_blank', icon: 'description' }
        ]
      }
    ]
  }
})

/**
 * Deteksi komponen dasar yang akan digunakan:
 * Jika memiliki properti 'url', maka render sebagai tag 'a' HTML biasa.
 * Jika tidak, gunakan 'q-item' bawaan Quasar untuk routing internal vue-router.
 */
const getLinkComponent = (item) => {
  return item.url ? 'a' : QItem
}

/**
 * Distribusi atribut secara dinamis:
 * Jika link eksternal: pasang href, target, style reset agar tag 'a' tidak berwarna biru khas browser.
 * Jika link internal: pasang atribut 'to' bawaan Quasar routing.
 */
// const getLinkAttributes = (item) => {
//   if (item.url) {
//     return {
//       href: item.url,
//       target: item.target || '_self',
//       class: 'q-item row no-wrap items-center text-black text-decoration-none',
//       style: { textDecoration: 'none', color: 'inherit' }
//     }
//   }
//   return {
//     to: item.to
//   }
// }const getLinkComponent = (item) => {
//   return item.url ? 'a' : QItem
// }

const getLinkAttributes = (item) => {
  if (item.url) {
    return {
      href: item.url,
      target: item.target || '_self',
      // FIX: Menambahkan class 'q-item q-item--type q-hoverable' agar tag 'a' mewarisi style hover murni Quasar
      class: 'q-item q-item--type q-hoverable row no-wrap items-center text-black custom-href-link',
    }
  }
  return {
    to: item.to
  }
}
</script>

<style scoped>
.font-tab-match {
  font-weight: 500;
}

/* FIX: Memastikan tag 'a' bersih dari style underline default browser dan mewarisi pointer */
.custom-href-link {
  text-decoration: none;
  color: inherit;
  cursor: pointer;
  position: relative;
  /* Wajib ada agar q-focus-helper membatasi areanya hanya di dalam baris */
}
</style>
