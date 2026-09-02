<!-- ============================================================
  components/TeacherSection.vue
  Grid kartu profil guru + sambutan Kepala Sekolah
============================================================ -->
<template>
  <section id="guru" class="q-py-xl bg-body">
    <div class="max-width-center q-mx-auto q-px-md">
      <!-- ── Sambutan Kepala Sekolah ─────────────────────── -->
      <q-card flat bordered class="principal-card q-mb-xl overflow-hidden">
        <div class="row items-center no-wrap">
          <div
            class="principal-accent col-auto"
            :style="{ background: 'linear-gradient(135deg, var(--q-primary), #1976d2)' }"
          />
          <div class="row items-center q-pa-lg full-width q-col-gutter-lg">
            <div class="col-12 col-sm-auto text-center">
              <q-avatar size="110px" class="shadow-4 principal-avatar">
                <img
                  src="https://api.dicebear.com/7.x/avataaars/svg?seed=principal&backgroundColor=d1f4e0"
                  alt="Kepala Sekolah"
                />
              </q-avatar>
            </div>
            <div class="col">
              <div
                class="text-caption text-weight-bold text-primary text-uppercase letter-spacing-wide q-mb-xs"
              >
                Sambutan Kepala Sekolah
              </div>
              <blockquote class="principal-quote text-body1 text-grey-8 q-mb-md">
                "Kami berkomitmen menyediakan pendidikan kejuruan yang relevan dengan industri
                global. Bersama staf pengajar kompeten dan fasilitas modern, kami siap mengantarkan
                putra-putri Anda menuju gerbang kesuksesan karier."
              </blockquote>
              <div class="text-subtitle2 text-weight-bold text-primary">
                Dr. H. Ahmad Subarjo, M.Pd.
              </div>
              <div class="text-caption text-grey-5">Kepala SMK Pasundan Jatinangor</div>
            </div>
          </div>
        </div>
      </q-card>

      <!-- ── Header grid guru ───────────────────────────── -->
      <div class="row items-end justify-between q-mb-lg">
        <SectionHeader
          title="Tenaga Pendidik & Instruktur Ahli"
          subtitle="Dibimbing langsung oleh praktisi industri dan akademisi berpengalaman."
        />
        <q-btn
          outline
          color="primary"
          no-caps
          label="Semua Pengajar"
          icon="groups"
          class="gt-xs text-weight-bold q-mb-sm"
          to="/pengajar"
        />
      </div>

      <!-- ── Grid kartu guru ────────────────────────────── -->
      <div class="row q-col-gutter-md">
        <div v-for="(t, i) in TEACHERS" :key="i" class="col-12 col-sm-6 col-md-3">
          <q-card class="teacher-card text-center shadow-1" flat bordered>
            <div class="teacher-header">
              <div class="teacher-avatar-wrap">
                <q-avatar size="96px" class="teacher-avatar shadow-3">
                  <img :src="t.photo" :alt="t.name" />
                </q-avatar>
              </div>
            </div>
            <q-card-section class="q-pt-none q-pb-lg">
              <div class="text-subtitle1 text-weight-bold q-mt-sm">{{ t.name }}</div>
              <div
                class="text-caption text-primary text-weight-bold text-uppercase letter-spacing-wide q-mt-xs"
              >
                {{ t.role }}
              </div>
              <q-separator class="q-my-md" />
              <div class="text-caption text-grey-6 q-mb-sm">
                <q-icon name="school" size="xs" class="q-mr-xs" />
                {{ t.expertise }}
              </div>
              <q-badge v-if="t.cert" color="amber-8" class="text-weight-bold text-caption q-px-sm">
                <q-icon name="verified" size="xs" class="q-mr-xs" />
                {{ t.cert }}
              </q-badge>
            </q-card-section>
          </q-card>
        </div>
      </div>

      <!-- Mobile: tombol lihat semua -->
      <div class="row justify-center q-mt-lg lt-sm">
        <q-btn
          outline
          color="primary"
          no-caps
          label="Lihat Semua Pengajar"
          icon="groups"
          class="full-width text-weight-bold"
          to="/pengajar"
        />
      </div>
    </div>
  </section>
</template>

<script setup>
import { TEACHERS } from '../../mocks/data/schoolData'
import SectionHeader from './SectionHeader.vue'
</script>

<style scoped>
/* Principal card */
.principal-card {
  border-radius: 16px;
}
.principal-accent {
  width: 6px;
  align-self: stretch;
  min-height: 120px;
}
.principal-quote {
  font-style: italic;
  line-height: 1.8;
  border-left: 3px solid var(--q-primary);
  padding-left: 16px;
  margin: 0 0 12px;
  color: #555;
}
.principal-avatar {
  border: 3px solid white;
}

/* Teacher card */
.teacher-card {
  border-radius: 16px;
  overflow: hidden;
  transition:
    transform 0.25s,
    box-shadow 0.25s;
  height: 100%;
}
.teacher-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.1) !important;
}

.teacher-header {
  background: linear-gradient(135deg, var(--q-primary) 0%, #1976d2 100%);
  height: 80px;
  display: flex;
  align-items: flex-end;
  justify-content: center;
  padding-bottom: 0;
}
.teacher-avatar-wrap {
  transform: translateY(48px);
}
.teacher-avatar {
  border: 3px solid white;
  background: white;
}

.letter-spacing-wide {
  letter-spacing: 1px;
}
.max-width-center {
  max-width: 1100px;
}
</style>
