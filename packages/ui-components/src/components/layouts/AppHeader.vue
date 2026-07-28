<!-- packages/ui-components/src/components/AppHeader.vue -->
<template>
  <q-header elevated class="bg-primary text-white">
    <q-toolbar>
      <!-- Toggle Drawer (Mobile) -->
      <q-btn flat dense round icon="menu" aria-label="Menu" @click="emit('toggle-drawer')" class="lt-md" />

      <q-toolbar-title class="row items-center no-wrap">
        <q-avatar size="32px" class="q-mr-sm">
          <img :src="logo" alt="Logo" />
        </q-avatar>
        <span class="text-weight-bold">{{ appTitle }}</span>
      </q-toolbar-title>

      <!-- Slot untuk custom actions (misal: tombol tambahan) -->
      <slot name="header-actions" />

      <!-- User Info & Actions (Desktop) -->
      <div class="gt-sm row items-center q-gutter-sm">
        <q-btn flat round dense icon="notifications" />
        <q-chip outline color="white" text-color="white" icon="person">
          {{ userName }}
        </q-chip>
        <q-btn flat round dense icon="logout" @click="emit('logout')">
          <q-tooltip>Keluar</q-tooltip>
        </q-btn>
      </div>
    </q-toolbar>
  </q-header>
</template>

<script setup>
defineProps({
  appTitle: {
    type: String,
    required: true,
  },
  logo: {
    type: String,
    default: '~@/assets/logo-smk.png', // fallback, di-override tiap app
  },
  userName: {
    type: String,
    default: 'Pengguna',
  },
})

const emit = defineEmits(['toggle-drawer', 'logout'])
</script>
