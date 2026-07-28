<!-- ============================================================
  components/OsisSection.vue
  Highlight pengurus OSIS
============================================================ -->
<template>
  <section id="osis" :class="['q-py-xl', $q.dark.isActive ? 'bg-dark' : 'bg-grey-1']">
    <div class="max-width-center q-mx-auto q-px-md">
      <SectionHeader
        title="Kreativitas & Kepemimpinan Siswa"
        subtitle="Mengenal pilar penggerak kegiatan kesiswaan dan kolaborator Mading Digital SMK."
        align="center"
      />

      <div class="row q-col-gutter-lg justify-center">
        <div
          v-for="(o, i) in OSIS"
          :key="i"
          class="col-12 col-sm-4"
        >
          <q-card flat bordered class="osis-card text-center" :class="$q.dark.isActive ? 'bg-grey-9' : 'bg-white'">
            <!-- Header gradient -->
            <div class="osis-header" />

            <q-card-section class="q-pt-none osis-body">
              <q-avatar
                size="88px"
                class="osis-avatar shadow-3"
              >
                <img :src="o.photo" :alt="o.name" />
              </q-avatar>

              <div class="text-subtitle1 text-weight-bold q-mt-md q-mb-xs">{{ o.name }}</div>
              <q-badge color="blue-7" class="text-weight-bold q-px-sm q-mb-sm">
                {{ o.role }}
              </q-badge>
              <div class="text-caption text-grey-5 row items-center justify-center">
                <q-icon name="school" size="xs" class="q-mr-xs" />
                {{ o.major }}
              </div>
            </q-card-section>
          </q-card>
        </div>
      </div>

      <div class="row justify-center q-mt-xl">
        <q-btn
          outline color="primary" no-caps
          label="Lihat Program Kerja & Kegiatan OSIS"
          icon="collections"
          to="/mading"
          class="text-weight-bold"
        />
      </div>
    </div>
  </section>
</template>

<script setup>
import { useQuasar } from 'quasar'
import { OSIS } from '../../mocks/data/schoolData'
import SectionHeader from './SectionHeader.vue'

const $q = useQuasar()
</script>

<style scoped>
.osis-card {
  border-radius: 16px;
  overflow: hidden;
  transition: transform .2s, box-shadow .2s;
}
.osis-card:hover {
  transform: translateY(-5px) scale(1.02);
  box-shadow: 0 10px 28px rgba(0,0,0,.1) !important;
}
.osis-header {
  height: 72px;
  background: linear-gradient(135deg, var(--q-primary) 0%, #1976d2 100%);
}
.osis-body {
  padding-top: 0;
}
.osis-avatar {
  transform: translateY(-44px);
  border: 3px solid white;
  background: white;
  margin-bottom: -44px;
}
.max-width-center { max-width: 1100px; }
</style>
