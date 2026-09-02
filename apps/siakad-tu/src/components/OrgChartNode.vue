<!--
  OrgChartNode.vue
  Satu simpul (kotak jabatan) pada struktur organisasi, dipanggil secara rekursif
  untuk merender anak-anaknya. Garis penghubung sepenuhnya dibuat dengan CSS
  (::before/::after pada <li>) - tanpa library diagram eksternal.
-->
<template>
  <li>
    <div class="org-node" :style="nodeStyle" @click="$emit('node-click', node)">
      <q-icon :name="typeMeta.icon" :style="{ color: iconColor }" size="20px" class="org-node__icon" />
      <div class="org-node__jabatan" :style="{ color: nodeTextColor }">{{ node.jabatan }}</div>
      <div v-if="node.nama" class="org-node__nama" :style="{ color: nodeTextColor }">{{ node.nama }}</div>
    </div>

    <ul v-if="node.children && node.children.length">
      <OrgChartNode
        v-for="child in node.children"
        :key="child.id"
        :node="child"
        @node-click="$emit('node-click', $event)"
      />
    </ul>
  </li>
</template>

<script setup>
import { computed } from 'vue'
import { ORG_NODE_TYPES } from '@/data/organizationData'

const props = defineProps({
  node: { type: Object, required: true }
})
defineEmits(['node-click'])

const typeMeta = computed(() => ORG_NODE_TYPES[props.node.tipe] || ORG_NODE_TYPES.staff)

const nodeStyle = computed(() => ({
  backgroundColor: typeMeta.value.fill,
  borderColor: typeMeta.value.border,
  borderStyle: props.node.tipe === 'koordinatif' ? 'dashed' : 'solid'
}))

const nodeTextColor = computed(() => typeMeta.value.text)
const iconColor = computed(() => typeMeta.value.text)
</script>

<style scoped>
li {
  position: relative;
  padding: 24px 8px 0;
  text-align: center;
  list-style: none;
}

/* garis vertikal & horizontal penghubung antar level, teknik CSS org-chart klasik */
li::before,
li::after {
  content: '';
  position: absolute;
  top: 0;
  right: 50%;
  width: 50%;
  height: 24px;
  border-top: 2px solid #cfd8dc;
}
li::after {
  right: auto;
  left: 50%;
  border-left: 2px solid #cfd8dc;
}
li:only-child::before,
li:only-child::after {
  display: none;
}
li:only-child {
  padding-top: 0;
}
li:first-child::before,
li:last-child::after {
  border: 0 none;
}
li:last-child::before {
  border-right: 2px solid #cfd8dc;
  border-radius: 0 6px 0 0;
}
li:first-child::after {
  border-radius: 6px 0 0 0;
}

ul {
  display: flex;
  padding-top: 24px;
  position: relative;
  margin: 0;
}
ul::before {
  content: '';
  position: absolute;
  top: 0;
  left: 50%;
  border-left: 2px solid #cfd8dc;
  width: 0;
  height: 24px;
}

.org-node {
  display: inline-flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  min-width: 148px;
  max-width: 180px;
  padding: 10px 12px;
  border-radius: 8px;
  border-width: 2px;
  cursor: pointer;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.org-node:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.12);
}

.org-node__icon {
  margin-bottom: 2px;
}

.org-node__jabatan {
  font-size: 0.78rem;
  font-weight: 700;
  line-height: 1.25;
}

.org-node__nama {
  font-size: 0.72rem;
  opacity: 0.85;
  line-height: 1.25;
}
</style>
