<!--
  OrgChart.vue
  Komponen struktur organisasi yang reusable. Cukup diberi `data` (simpul akar + children
  bertingkat) untuk dipakai di konteks manapun - struktur sekolah, struktur panitia PPDB,
  struktur OSIS, dsb. `sideNode` opsional untuk entitas koordinatif (mis. Komite Sekolah)
  yang digambar di samping akar dengan garis putus-putus, bukan sebagai bawahan.
-->
<template>
  <div class="org-chart">
    <div class="org-chart__scroll">
      <div class="org-chart__inner">
        <div class="org-chart__root-row" :class="{ 'org-chart__root-row--with-side': sideNode }">
          <div v-if="sideNode" class="org-chart__side">
            <div class="org-node org-node--side" :style="sideNodeStyle">
              <q-icon :name="sideTypeMeta.icon" :style="{ color: sideTypeMeta.text }" size="20px" />
              <div class="org-node__jabatan" :style="{ color: sideTypeMeta.text }">{{ sideNode.jabatan }}</div>
              <div v-if="sideNode.nama" class="org-node__nama" :style="{ color: sideTypeMeta.text }">
                {{ sideNode.nama }}
              </div>
            </div>
            <div class="org-chart__side-connector" />
          </div>

          <ul class="org-chart__tree">
            <OrgChartNode :node="data" @node-click="$emit('node-click', $event)" />
          </ul>
        </div>
      </div>
    </div>

    <div class="org-chart__hint text-caption text-grey-6 gt-xs">
      Gunakan scroll horizontal bila struktur lebih lebar dari layar.
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import OrgChartNode from './OrgChartNode.vue'
import { ORG_NODE_TYPES } from '@/data/organizationData'

const props = defineProps({
  data: { type: Object, required: true },
  sideNode: { type: Object, default: null }
})
defineEmits(['node-click'])

const sideTypeMeta = computed(() => {
  if (!props.sideNode) return ORG_NODE_TYPES.staff
  return ORG_NODE_TYPES[props.sideNode.tipe] || ORG_NODE_TYPES.staff
})

const sideNodeStyle = computed(() => ({
  backgroundColor: sideTypeMeta.value.fill,
  borderColor: sideTypeMeta.value.border,
  borderStyle: 'dashed'
}))
</script>

<style scoped>
.org-chart__scroll {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
  padding: 8px 4px 16px;
}

.org-chart__inner {
  display: inline-flex;
  min-width: 100%;
  justify-content: center;
}

.org-chart__root-row {
  display: flex;
  align-items: center;
}

.org-chart__side {
  position: relative;
  margin-right: 48px;
  flex-shrink: 0;
}

.org-chart__side-connector {
  position: absolute;
  top: 50%;
  right: -48px;
  width: 48px;
  border-top: 2px dashed #d4af37;
}

.org-node--side {
  display: inline-flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  min-width: 148px;
  max-width: 180px;
  padding: 10px 12px;
  border-radius: 8px;
  border-width: 2px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
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

.org-chart__tree {
  display: flex;
  padding: 0;
  margin: 0;
}

.org-chart__hint {
  text-align: center;
  display: none;
}

/* Tampilkan petunjuk scroll hanya di layar sempit tempat tree kemungkinan melebar */
@media (max-width: 1023px) {
  .org-chart__hint {
    display: block;
  }
}
</style>
