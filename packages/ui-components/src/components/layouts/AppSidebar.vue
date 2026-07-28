<!-- packages/ui-components/src/components/AppSidebar.vue -->
<template>
  <q-drawer v-model="internalDrawer" show-if-above bordered class="bg-grey-1" :width="260">
    <q-scroll-area class="fit">
      <q-list padding>
        <!-- Loop Menu Items -->
        <div v-for="(group, i) in menuGroups" :key="i">
          <q-item-label v-if="group.label" header class="text-grey-8 text-weight-bold">
            {{ group.label }}
          </q-item-label>

          <q-item v-for="item in group.items" :key="item.name" clickable v-ripple :to="item.to" exact
            active-class="text-primary bg-blue-1">
            <q-item-section avatar>
              <q-icon :name="item.icon" />
            </q-item-section>
            <q-item-section>{{ item.label }}</q-item-section>
          </q-item>

          <q-separator v-if="group.separatorAfter" class="q-my-md" />
        </div>

        <!-- Footer Menu (Logout) -->
        <q-separator class="q-my-md" />
        <q-item clickable v-ripple @click="emit('logout')">
          <q-item-section avatar><q-icon name="logout" /></q-item-section>
          <q-item-section>Keluar</q-item-section>
        </q-item>
      </q-list>
    </q-scroll-area>
  </q-drawer>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  modelValue: {
    type: Boolean,
    required: true,
  },
  menuGroups: {
    type: Array,
    required: true,
    // Contoh structure:
    // [
    //   { label: 'MENU UTAMA', items: [{ label: 'Dashboard', icon: 'dashboard', to: { name: 'dashboard' } }] },
    //   { label: 'LAINNYA', items: [...], separatorAfter: true }
    // ]
  },
})

const emit = defineEmits(['update:modelValue', 'logout'])

const internalDrawer = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val),
})
</script>
