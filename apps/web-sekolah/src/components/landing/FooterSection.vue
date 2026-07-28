<!-- ============================================================
  components/FooterSection.vue
============================================================ -->
<template>
  <footer class="footer-section q-pt-xl">
    <div class="max-width-center q-mx-auto q-px-md">
      <div class="row q-col-gutter-xl q-pb-xl">

        <!-- ── Identitas sekolah ──────────────────── -->
        <div class="col-12 col-md-4">
          <div class="row items-center no-wrap q-mb-md">
            <div class="footer-logo-mark q-mr-sm">
              <!-- <q-avatar size="40px" class="bg-transparent"> -->
                <img src="~@/assets/logo-smk.png" alt="Logo" onerror="this.style.display='none'" width="40px"/>
              <!-- </q-avatar> -->
            </div>
            <div class="text-h6 text-weight-bold text-white">SMK Pasundan Jatinangor</div>
          </div>
          <p class="text-caption text-grey-5 line-height-relaxed q-mb-lg">
            Sekolah vokasi unggulan yang berkomitmen mencetak generasi kompeten,
            berkarakter, dan siap kerja melalui program penyelarasan kurikulum
            standar industri internasional.
          </p>
          <!-- Sosmed -->
          <div class="row q-gutter-sm">
            <q-btn v-for="s in socmedLinks" :key="s.platform" round dense flat :color="s.color" :icon="s.icon"
              class="footer-sosmed-btn shadow-1" @click="bukaMedsos(s.platform)">
              <q-tooltip>{{ s.label }}</q-tooltip>
            </q-btn>
          </div>
        </div>

        <!-- ── Kontak ─────────────────────────────── -->
        <div class="col-12 col-sm-6 col-md-4">
          <div class="footer-heading">Informasi Kontak</div>
          <q-list dense class="footer-list">
            <q-item class="q-px-none items-start">
              <q-item-section avatar style="min-width:28px">
                <q-icon name="place" color="amber-4" size="xs" />
              </q-item-section>
              <q-item-section class="text-caption text-grey-5">
                {{ SCHOOL_INFO.address }}
              </q-item-section>
            </q-item>
            <q-item class="q-px-none">
              <q-item-section avatar style="min-width:28px">
                <q-icon name="phone" color="amber-4" size="xs" />
              </q-item-section>
              <q-item-section class="text-caption text-grey-5">{{ SCHOOL_INFO.phone }}</q-item-section>
            </q-item>
            <q-item class="q-px-none">
              <q-item-section avatar style="min-width:28px">
                <q-icon name="mail" color="amber-4" size="xs" />
              </q-item-section>
              <q-item-section class="text-caption text-grey-5">{{ SCHOOL_INFO.email }}</q-item-section>
            </q-item>
            <q-item class="q-px-none">
              <q-item-section avatar style="min-width:28px">
                <q-icon name="chat" color="green-4" size="xs" />
              </q-item-section>
              <q-item-section class="text-caption text-grey-5">
                WhatsApp: 0812-3456-7890
              </q-item-section>
            </q-item>
          </q-list>
        </div>

        <!-- ── Tautan cepat ───────────────────────── -->
        <div class="col-12 col-sm-6 col-md-4">
          <div class="footer-heading">Tautan Cepat</div>
          <div v-for="link in quickLinks" :key="link.label" class="footer-link" @click="$router.push(link.to)">
            <q-icon name="chevron_right" size="xs" class="q-mr-xs" />
            {{ link.label }}
          </div>
        </div>
      </div>

      <q-separator dark class="q-mb-md" />

      <!-- Copyright -->
      <div class="row justify-between items-center q-pb-lg text-caption text-grey-6 flex-wrap gap-xs">
        <div>Made with
          <q-icon name="favorite" color="red-10" size="xs" class="q-mx-xs" />
          by. Deden Moh Jaenudin.
        </div>
        <div class="row items-center gap-xs">
          © {{ currentYear }}. All Rights Reserved.
        </div>
      </div>
    </div>
  </footer>
</template>

<script setup>
import { SCHOOL_INFO } from '../../mocks/data/schoolData'
import { useSchool } from '../../composables/useSchool'
import { useUtil } from '@/composables/helper/useUtils'

const { bukaMedsos } = useSchool()
const { currentYear } = useUtil()
const socmedLinks = [
  { platform: 'instagram', icon: 'bi-instagram', color: 'pink-4', label: 'Instagram Resmi' },
  { platform: 'facebook', icon: 'bi-facebook', color: 'blue-5', label: 'Facebook Fanspage' },
  { platform: 'youtube', icon: 'bi-youtube', color: 'red-5', label: 'YouTube Channel' },
  { platform: 'whatsapp', icon: 'bi-whatsapp', color: 'green-5', label: 'WhatsApp' },
  { platform: 'tiktok', icon: 'bi-tiktok', color: 'green-5', label: 'Tiktok' },
]

const quickLinks = [
  { label: 'Portal PPDB Online', to: '/ppdb' },
  { label: 'e-Rapor Siswa', to: '/erapor' },
  { label: 'Mading Digital', to: '/mading' },
  { label: 'Bursa Kerja (BKK)', to: '/bkk' },
  { label: 'Kalender Akademik', to: '/agenda' },
  { label: 'Hubungi Kami', to: '/kontak' },
]
</script>

<style scoped>
.footer-section {
  background: #0d1220;
  color: rgba(255, 255, 255, .7);
}

.footer-logo-mark {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background: var(--q-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 900;
  color: white;
  font-size: 13px;
  flex-shrink: 0;
}

.footer-heading {
  color: white;
  font-weight: 700;
  font-size: 13px;
  margin-bottom: 14px;
  padding-bottom: 8px;
  border-bottom: 1px solid rgba(255, 255, 255, .08);
}

.footer-link {
  color: rgba(255, 255, 255, .55);
  font-size: 12px;
  padding: 4px 0;
  cursor: pointer;
  display: flex;
  align-items: center;
  transition: color .15s;
}

.footer-link:hover {
  color: #ffca28;
}

.footer-sosmed-btn {
  background: rgba(255, 255, 255, .07) !important;
  transition: background .15s, transform .15s !important;
}

.footer-sosmed-btn:hover {
  background: rgba(255, 255, 255, .15) !important;
  transform: translateY(-2px);
}

.line-height-relaxed {
  line-height: 1.7;
}

.max-width-center {
  max-width: 1100px;
}

.gap-xs {
  gap: 8px;
}
</style>
