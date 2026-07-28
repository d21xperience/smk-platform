<!-- ============================================================
  components/NewsSection.vue
  Berita & Prestasi Terbaru — card grid modern
============================================================ -->
<template>
  <section id="berita" class="q-py-xl bg-body">
    <div class="max-width-center q-mx-auto q-px-md">
      <div class="row items-end justify-between q-mb-lg">
        <SectionHeader
          title="Berita & Prestasi Terbaru"
          subtitle="Kabar terkini dari kegiatan akademik, prestasi siswa, dan informasi sekolah."
        />
        <q-btn
          flat color="primary" no-caps
          label="Semua Berita" icon-right="arrow_forward"
          to="/berita" class="gt-xs text-weight-bold q-mb-sm"
        />
      </div>

      <div class="row q-col-gutter-md">
        <div
          v-for="(n, i) in NEWS"
          :key="n.id"
          class="col-12 col-sm-4"
        >
          <q-card
            flat bordered
            class="news-card"
            :class="i === 0 ? 'news-card--featured' : ''"
            clickable v-ripple
            @click="$router.push(n.to)"
          >
            <q-img
              :src="n.img"
              loading="lazy"
              :ratio="i === 0 ? 16/9 : 4/3"
              class="news-img"
            >
              <div v-if="i === 0" class="absolute-top-left q-ma-sm">
                <q-badge color="red-7" class="text-weight-bold q-px-sm">
                  🔥 Featured
                </q-badge>
              </div>
            </q-img>
            <q-card-section class="q-pa-md">
              <div class="text-caption text-grey-5 row items-center q-mb-sm">
                <q-icon name="calendar_today" size="xs" class="q-mr-xs" />
                {{ n.date }}
                <span class="q-mx-xs">•</span>
                <q-icon name="person" size="xs" class="q-mr-xs" />
                {{ n.author }}
              </div>
              <div class="text-subtitle2 text-weight-bold news-title">
                {{ n.title }}
              </div>
            </q-card-section>
            <q-card-actions class="q-pt-none q-px-md q-pb-md">
              <q-btn
                flat dense no-caps
                color="primary"
                label="Baca Selengkapnya"
                icon-right="arrow_forward"
                class="text-caption text-weight-bold"
                :to="n.to"
              />
            </q-card-actions>
          </q-card>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { NEWS } from '../../mocks/data/schoolData'
import SectionHeader from './SectionHeader.vue'
</script>

<style scoped>
.news-card {
  border-radius: 14px;
  overflow: hidden;
  transition: transform .2s, box-shadow .2s;
  height: 100%;
}
.news-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0,0,0,.1) !important;
}
.news-title {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.5;
}
.news-img { transition: transform .4s; }
.news-card:hover .news-img { transform: scale(1.03); }
.max-width-center { max-width: 1100px; }
</style>
