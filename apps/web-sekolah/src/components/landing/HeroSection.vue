<!-- ============================================================
  components/HeroSection.vue
  Hero dengan QCarousel full-viewport + CTA & stats bar
============================================================ -->
<template>
  <section id="hero">
    <!-- ── Carousel ───────────────────────────────────────── -->
    <q-carousel
      v-model="slide"
      animated arrows infinite
      :autoplay="5000"
      transition-prev="slide-right"
      transition-next="slide-left"
      height="90vh"
      class="hero-carousel"
    >
      <q-carousel-slide
        v-for="s in HERO_SLIDES"
        :key="s.id"
        :name="s.id"
        class="hero-slide q-pa-none"
      >
        <!-- Background image via CSS untuk performa -->
        <div class="hero-bg" :style="{ backgroundImage: `url(${s.img})` }" />
        <div class="hero-overlay" />

        <!-- Content -->
        <div class="hero-content absolute-full flex flex-center">
          <div class="text-center text-white q-px-md" style="max-width:700px">
            <transition appear enter-active-class="animated fadeInDown">
              <div class="hero-eyebrow text-amber-4 text-weight-bold text-caption text-uppercase letter-spacing-wide q-mb-sm">
                {{ s.eyebrow }}
              </div>
            </transition>
            <transition appear enter-active-class="animated fadeInUp">
              <h1 class="text-h3 text-sm-h2 text-weight-black q-mb-md hero-title">
                {{ s.title }}
              </h1>
            </transition>
            <transition appear enter-active-class="animated fadeIn">
              <p class="text-subtitle1 text-grey-3 q-mb-xl hero-sub">
                {{ s.subtitle }}
              </p>
            </transition>
            <q-btn
              :color="s.cta.highlight ? 'amber' : 'white'"
              :text-color="s.cta.highlight ? 'black' : 'primary'"
              :label="s.cta.label"
              unelevated
              size="lg"
              no-caps
              class="hero-cta text-weight-bold q-px-xl"
              @click="handleCta(s.cta)"
            />
          </div>
        </div>
      </q-carousel-slide>

      <!-- Custom navigation dots -->
      <template #control>
        <q-carousel-control position="bottom" :offset="[0, 16]" class="flex flex-center">
          <div class="hero-dots">
            <button
              v-for="s in HERO_SLIDES" :key="s.id"
              :class="['hero-dot', { active: slide === s.id }]"
              @click="slide = s.id"
            />
          </div>
        </q-carousel-control>
      </template>
    </q-carousel>

    <!-- ── Stats bar ──────────────────────────────────────── -->
    <div class="stats-bar bg-primary text-white">
      <div class="max-width-center q-mx-auto q-px-md">
        <div class="row q-col-gutter-xs text-center">
          <div
            v-for="(stat, i) in STATS"
            :key="i"
            class="col-6 col-sm-3 stats-item"
          >
            <q-icon :name="stat.icon" color="amber-4" size="sm" class="q-mb-xs" />
            <div class="text-h5 text-weight-black text-amber-3">{{ stat.value }}</div>
            <div class="text-caption text-blue-2 text-weight-medium">{{ stat.label }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- QPageSticky: tombol PPDB melayang -->
    <!-- <q-page-sticky position="bottom-right" :offset="[20, 20]" style="z-index:900">
      <q-btn
        fab no-caps
        color="red-7"
        icon="edit_note"
        label="Daftar PPDB"
        class="text-weight-bold ppdb-fab shadow-10"
        to="/ppdb"
      >
        <q-tooltip anchor="top middle" self="bottom middle">
          Buka Formulir PPDB 2025/2026
        </q-tooltip>
      </q-btn>
    </q-page-sticky> -->
  </section>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { HERO_SLIDES, STATS } from '../../mocks/data/schoolData'
import { useSchool } from '../../composables/useSchool'

const slide = ref(HERO_SLIDES[0].id)
const router = useRouter()
const { scrollToSection } = useSchool()

function handleCta(cta) {
  if (cta.to?.startsWith('#')) scrollToSection(cta.to.slice(1))
  else router.push(cta.to)
}
</script>

<style scoped>
.hero-carousel { overflow: hidden; }

.hero-bg {
  position: absolute; inset: 0;
  background-size: cover;
  background-position: center;
  transition: transform 8s ease;
  transform: scale(1.05);
}
.hero-slide:hover .hero-bg { transform: scale(1); }

.hero-overlay {
  position: absolute; inset: 0;
  background: linear-gradient(
    160deg,
    rgba(13, 27, 62, .85) 0%,
    rgba(13, 27, 62, .55) 60%,
    rgba(0, 0, 0, .4) 100%
  );
}

.hero-content { z-index: 2; }

.hero-eyebrow {
  letter-spacing: 2px;
  text-shadow: 0 1px 8px rgba(0,0,0,.5);
}
.hero-title {
  text-shadow: 0 2px 20px rgba(0,0,0,.5);
  line-height: 1.1;
}
.hero-sub {
  text-shadow: 0 1px 8px rgba(0,0,0,.4);
  line-height: 1.6;
}
.hero-cta {
  border-radius: 50px;
  letter-spacing: .5px;
  box-shadow: 0 8px 24px rgba(0,0,0,.3);
  transition: transform .2s, box-shadow .2s;
}
.hero-cta:hover { transform: translateY(-2px); box-shadow: 0 12px 32px rgba(0,0,0,.4); }

/* Dots navigation */
.hero-dots { display: flex; gap: 8px; }
.hero-dot {
  width: 8px; height: 8px; border-radius: 50%;
  background: rgba(255,255,255,.45); border: none; cursor: pointer;
  transition: all .3s;
}
.hero-dot.active {
  width: 28px; border-radius: 4px;
  background: #ffca28;
}

/* Stats bar */
.stats-bar { padding: 20px 0; }
.stats-item {
  padding: 10px 8px;
  border-right: 1px solid rgba(255,255,255,.12);
}
.stats-item:last-child { border-right: none; }

.ppdb-fab { border-radius: 50px !important; padding: 0 20px !important; }

.max-width-center { max-width: 1100px; }
</style>
