<!-- ============================================================
  pages/IndexPage.vue  (atau src/pages/HomePage.vue)
  Halaman utama — merakit semua section component

  Struktur komponen:
  IndexPage
  ├── HeroSection          ← hero carousel + stats bar + FAB PPDB
  ├── JurusanSection       ← carousel 4 jurusan + quick nav chips
  ├── TeacherSection       ← sambutan kasek + grid guru
  ├── AgendaSection        ← agenda + info libur
  ├── BkkSection           ← lowongan kerja alumni
  ├── NewsSection          ← berita & prestasi
  ├── OsisSection          ← profil pengurus OSIS
  ├── TestimoniSection     ← carousel testimoni
  ├── CtaSection           ← banner CTA + download brosur
  └── FooterSection        ← footer lengkap
============================================================ -->
<template>
  <q-page>
    <HeroSection />
    <JurusanSection />
    <EkskulSection />
    <TeacherSection />
    <AgendaSection />
    <BkkSection />
    <NewsSection />
    <OsisSection />
    <TestimoniSection />
    <CtaSection />
    <FooterSection />
  </q-page>
</template>

<script setup>
// ── Lazy-import seluruh section agar bundle splitting optimal ──
import { defineAsyncComponent } from 'vue'

const HeroSection = defineAsyncComponent(() => import('@/components/landing/HeroSection.vue'))
const JurusanSection = defineAsyncComponent(() => import('@/components/landing/JurusanSection.vue'))
const TeacherSection = defineAsyncComponent(() => import('@/components/landing/TeacherSection.vue'))
const AgendaSection = defineAsyncComponent(() => import('@/components/landing/AgendaSection.vue'))
const BkkSection = defineAsyncComponent(() => import('@/components/landing/BkkSection.vue'))
const NewsSection = defineAsyncComponent(() => import('@/components/landing/NewsSection.vue'))
const OsisSection = defineAsyncComponent(() => import('@/components/landing/OsisSection.vue'))
const TestimoniSection = defineAsyncComponent(
  () => import('@/components/landing/TestimoniSection.vue'),
)
const CtaSection = defineAsyncComponent(() => import('@/components/landing/CtaSection.vue'))
const FooterSection = defineAsyncComponent(() => import('@/components/landing/FooterSection.vue'))
const EkskulSection = defineAsyncComponent(() => import('@/components/landing/EkskulSection.vue'))
</script>
