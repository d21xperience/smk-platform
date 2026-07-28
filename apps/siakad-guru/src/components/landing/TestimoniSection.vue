<!-- ============================================================
  components/TestimoniSection.vue
============================================================ -->
<template>
  <section id="testimoni" :class="['q-py-xl', $q.dark.isActive ? 'bg-grey-9' : 'bg-blue-grey-1']">
    <div class="max-width-center q-mx-auto q-px-md text-center">
      <SectionHeader title="Apa Kata Mereka?" align="center" />

      <q-carousel v-model="activeIdx" transition-prev="scale" transition-next="scale" swipeable animated
        control-color="primary" navigation padding arrows infinite :autoplay="6000" height="auto"
        class="bg-transparent testi-carousel">
        <q-carousel-slide v-for="(t, i) in TESTIMONIES" :key="i" :name="i" class="column no-wrap flex-center q-pa-sm">
          <q-card flat :class="['testi-card q-pa-xl', $q.dark.isActive ? 'bg-grey-8' : 'bg-white']"
            style="max-width:680px;margin:0 auto">
            <!-- Quote icon -->
            <q-icon name="format_quote" size="48px" color="primary" class="q-mb-md opacity-30" />

            <q-avatar size="88px" class="shadow-3 q-mb-lg">
              <img :src="t.avatar" />
            </q-avatar>

            <p class="text-body1 text-grey-7 text-italic line-height-relaxed q-mb-lg">
              "{{ t.statement }}"
            </p>

            <div class="text-subtitle1 text-weight-bold text-primary">{{ t.name }}</div>
            <div class="text-caption text-amber-8 text-weight-bold q-mt-xs">{{ t.role }}</div>

            <!-- Rating bintang -->
            <div class="q-mt-md">
              <q-icon v-for="s in 5" :key="s" name="star" color="amber-5" size="xs" />
            </div>
          </q-card>
        </q-carousel-slide>
      </q-carousel>
    </div>
  </section>
</template>

<script setup>
import { ref } from 'vue'
import { useQuasar } from 'quasar'
import { TESTIMONIES } from '../../mocks/data/schoolData'
import SectionHeader from './SectionHeader.vue'

const $q = useQuasar()
const activeIdx = ref(0)
</script>

<style scoped>
.testi-carousel {
  border-radius: 16px;
  overflow: hidden;
}

.testi-card {
  border-radius: 20px;
}

.line-height-relaxed {
  line-height: 1.8;
}

.max-width-center {
  max-width: 1100px;
}
</style>
