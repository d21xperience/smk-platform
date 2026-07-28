<!-- packages/ui-components/src/components/AppLayout.vue -->
<template>
  <q-layout view="hHh lpR fFf">
    <!-- Header -->
    <AppHeader :app-title="appTitle" :logo="logo" :user-name="userName" @toggle-drawer="toggleDrawer"
      @logout="emit('logout')">
      <template #header-actions>
        <slot name="header-actions" />
      </template>
    </AppHeader>

    <!-- Sidebar -->
    <AppSidebar v-model="drawerOpen" :menu-groups="menuGroups" @logout="emit('logout')" />

    <!-- Main Content -->
    <q-page-container class="bg-grey-2">
      <slot />
    </q-page-container>
  </q-layout>
</template>

<script setup>
import { ref } from 'vue'
import AppHeader from './AppHeader.vue'
import AppSidebar from './AppSidebar.vue'

const props = defineProps({
  appTitle: {
    type: String,
    required: true,
  },
  logo: {
    type: String,
    default: '~@/assets/logo-smk.png',
  },
  userName: {
    type: String,
    default: 'Pengguna',
  },
  menuGroups: {
    type: Array,
    required: true,
  },
})

const emit = defineEmits(['logout'])

const drawerOpen = ref(false)

function toggleDrawer() {
  drawerOpen.value = !drawerOpen.value
}
</script>
