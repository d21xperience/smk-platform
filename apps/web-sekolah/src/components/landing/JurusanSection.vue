<!-- ============================================================
  components/JurusanSection.vue
  Carousel jurusan dengan card modern — full refactor
============================================================ -->
<template>
  <section
    id="jurusan"
    :class="['jurusan-section q-py-xl', $q.dark.isActive ? 'bg-dark' : 'bg-grey-1']"
  >
    <div class="max-width-center q-mx-auto q-px-md">
      <SectionHeader
        title="Kompetensi Keahlian Unggulan"
        subtitle="Geser untuk menjelajahi program keahlian yang diselaraskan langsung dengan kebutuhan industri global."
        align="center"
      />

      <!-- Carousel jurusan -->
      <q-carousel
        v-model="activeSlide"
        transition-prev="slide-right"
        transition-next="slide-left"
        swipeable
        animated
        control-color="primary"
        navigation
        arrows
        infinite
        :autoplay="6000"
        height="auto"
        class="bg-transparent no-shadow jurusan-carousel"
      >
        <q-carousel-slide v-for="j in JURUSAN" :key="j.id" :name="j.id" class="q-pa-xs q-pa-sm-md">
          <q-card flat bordered class="jurusan-card">
            <div class="row no-wrap items-stretch q-col-jurusan">
              <!-- ── Gambar kiri ──────────────────────── -->
              <div class="col-12 col-sm-5 jurusan-img-wrap">
                <q-img :src="j.img" loading="lazy" height="100%" class="jurusan-img" fit="cover">
                  <div class="absolute-top-left q-ma-md">
                    <q-badge
                      :color="j.badgeColor"
                      class="text-weight-bold q-px-sm q-py-xs text-caption"
                    >
                      {{ j.badge }}
                    </q-badge>
                  </div>
                  <div class="absolute-bottom q-pa-md jurusan-img-overlay">
                    <div class="text-white text-caption text-weight-bold">
                      <q-icon name="handshake" size="xs" class="q-mr-xs" />
                      {{ j.mitra }}
                    </div>
                  </div>
                </q-img>
              </div>

              <!-- ── Konten kanan ────────────────────── -->
              <q-card-section class="col-12 col-sm-7 q-pa-lg flex flex-column justify-between">
                <div>
                  <!-- Judul & Akreditasi -->
                  <div class="row items-start justify-between no-wrap q-mb-md">
                    <div>
                      <div
                        class="text-caption text-weight-bold text-grey-5 letter-spacing-wide text-uppercase q-mb-xs"
                      >
                        Program Keahlian
                      </div>
                      <div class="text-h5 text-weight-black text-grey-9 dark-text">
                        {{ j.name }}
                      </div>
                      <div class="text-subtitle2 text-grey-6 q-mt-xs">{{ j.fullName }}</div>
                    </div>
                    <q-chip
                      dense
                      square
                      :color="j.color"
                      text-color="white"
                      class="text-weight-bold q-ml-sm"
                      icon="verified"
                    >
                      Akreditasi {{ j.akreditasi }}
                    </q-chip>
                  </div>

                  <!-- Deskripsi -->
                  <p class="text-body2 text-grey-7 line-height-relaxed q-mb-md">
                    {{ j.desc }}
                  </p>

                  <!-- Keunggulan chips -->
                  <div class="q-gutter-xs q-mb-md">
                    <q-chip
                      v-for="(k, i) in j.keunggulan"
                      :key="i"
                      dense
                      outline
                      :color="j.color"
                      class="text-caption"
                      icon="check_circle"
                    >
                      {{ k }}
                    </q-chip>
                  </div>
                </div>

                <!-- CTA -->
                <div class="row items-center justify-between q-pt-sm border-top">
                  <div class="text-caption text-grey-5 text-weight-medium">
                    <q-icon name="groups" size="xs" class="q-mr-xs" />
                    Tersedia kursi baru TA 2025/2026
                  </div>
                  <q-btn
                    unelevated
                    no-caps
                    :color="j.color"
                    :label="`Eksplor ${j.name}`"
                    :icon="j.icon"
                    class="text-weight-bold"
                    :to="j.to"
                  />
                </div>
              </q-card-section>
            </div>
          </q-card>
        </q-carousel-slide>
      </q-carousel>

      <!-- Quick nav tabs -->
      <div class="row justify-center q-mt-lg q-gutter-sm">
        <q-chip
          v-for="j in JURUSAN"
          :key="j.id"
          clickable
          :selected="activeSlide === j.id"
          :color="activeSlide === j.id ? j.color : 'grey-3'"
          :text-color="activeSlide === j.id ? 'white' : 'grey-7'"
          class="text-weight-bold jurusan-tab-chip"
          @click="activeSlide = j.id"
        >
          <q-icon :name="j.icon" class="q-mr-xs" size="xs" />
          {{ j.name }}
        </q-chip>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref } from 'vue'
import { useQuasar } from 'quasar'
import { JURUSAN } from '../../mocks/data/schoolData'
import SectionHeader from './SectionHeader.vue'

const $q = useQuasar()
const activeSlide = ref(JURUSAN[0].id)
</script>

<style scoped>
.jurusan-carousel {
  border-radius: 16px;
}

.jurusan-card {
  border-radius: 16px;
  overflow: hidden;
  min-height: 340px;
}

.jurusan-img-wrap {
  min-height: 240px;
}

.jurusan-img {
  height: 100%;
  min-height: 240px;
}

.jurusan-img-overlay {
  background: linear-gradient(transparent, rgba(0, 0, 0, 0.6));
}

.jurusan-tab-chip {
  transition: all 0.2s;
  cursor: pointer;
}
.jurusan-tab-chip:hover {
  transform: translateY(-2px);
}

.border-top {
  border-top: 1px solid rgba(0, 0, 0, 0.08);
}
.dark-text {
  color: var(--q-grey-9) !important;
}
.body--dark .dark-text {
  color: #e0e0e0 !important;
}
.body--dark .border-top {
  border-color: rgba(255, 255, 255, 0.08);
}

.letter-spacing-wide {
  letter-spacing: 1px;
}
.line-height-relaxed {
  line-height: 1.7;
}
.max-width-center {
  max-width: 1100px;
}

@media (max-width: 599px) {
  .q-col-jurusan {
    flex-direction: column;
  }
}
</style>
