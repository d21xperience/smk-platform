<template>

  <!-- ===== DENGAN CHILDREN → q-expansion-item ===== -->
  <q-expansion-item v-if="hasChildren" :icon="icon" :label="title" :caption="caption" :group="group"
    header-class="text-weight-medium" expand-separator>
    <q-item v-for="child in children" :key="child.link" clickable v-ripple :to="child.link" active-class="text-primary">
      <q-item-section avatar>
        <q-icon :name="child.icon" />
      </q-item-section>
      <q-item-section>{{ child.title }}</q-item-section>
    </q-item>
  </q-expansion-item>

  <!-- ===== TANPA CHILDREN → q-item biasa ===== -->
  <q-item v-else clickable v-ripple :to="link" active-class="text-primary">
    <q-item-section v-if="icon" avatar>
      <q-icon :name="icon" />
    </q-item-section>
    <q-item-section>
      <q-item-label>{{ title }}</q-item-label>
      <q-item-label v-if="caption" caption>{{ caption }}</q-item-label>
    </q-item-section>
  </q-item>

</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  title: { type: String, required: true },
  icon: { type: String, default: '' },
  caption: { type: String, default: '' },
  link: { type: String, default: '#' },
  group: { type: String, default: '' },     // untuk accordion group
  children: { type: Array, default: () => [] }
})

const hasChildren = computed(() => props.children.length > 0)
</script>
