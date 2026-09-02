<template>
  <q-page class="bg-grey-1 q-py-xl">
    <div class="container q-px-md max-width-center">
      <!-- 1. HEADER HALAMAN FAQ -->
      <div class="text-center q-mb-xl">
        <div class="text-h4 text-sm-h3 text-weight-bold text-primary">
          Pusat Bantuan & <span class="text-amber">FAQ</span>
        </div>
        <p class="text-body1 text-grey-7 q-mt-sm">
          Temukan jawaban instan seputar pendaftaran, biaya, kurikulum, dan fasilitas SMK Pasundan
          Jatinangor.
        </p>
      </div>

      <!-- 2. KOTAK PENCARIAN PINTAR -->
      <div class="row justify-center q-mb-xl">
        <div class="col-12 col-sm-8 col-md-6">
          <q-input
            outlined
            v-model="searchQuery"
            placeholder="Ketik kata kunci (misal: Biaya, TBSM, Beasiswa)..."
            class="bg-white shadow-1 rounded-borders"
            clearable
            bg-color="white"
          >
            <template v-slot:prepend>
              <q-icon name="search" color="primary" />
            </template>
          </q-input>
        </div>
      </div>

      <!-- 3. LIST AKORDION FAQ BERDASARKAN KATEGORI -->
      <div class="row justify-center">
        <div class="col-12 col-md-8">
          <div v-for="(group, groupIndex) in filteredFaqGroups" :key="groupIndex" class="q-mb-lg">
            <!-- Nama Kategori Kelompok FAQ -->
            <div
              class="text-subtitle1 text-weight-bold text-primary q-mb-sm q-pl-sm row items-center"
            >
              <q-icon :name="group.icon" class="q-mr-sm" size="xs" />
              {{ group.categoryName }}
            </div>

            <!-- Daftar Pertanyaan di Dalam Kategori -->
            <q-list bordered class="rounded-borders bg-white shadow-1 separator">
              <q-expansion-item
                v-for="(faq, faqIndex) in group.items"
                :key="faqIndex"
                group="faq-accordion"
                :label="faq.question"
                header-class="text-weight-medium text-grey-9 text-subtitle2 text-sm-subtitle1"
                expand-icon-class="text-primary"
              >
                <q-card class="bg-grey-2">
                  <q-card-section class="text-grey-8 text-body2 q-py-md line-height-relaxed">
                    <!-- Menggunakan v-html jika di masa depan ada jawaban yang butuh format link/bold -->
                    <span v-html="faq.answer"></span>
                  </q-card-section>
                </q-card>
              </q-expansion-item>
            </q-list>
          </div>

          <!-- TAMPILAN JIKA PENCARIAN TIDAK DITEMUKAN -->
          <div v-if="filteredFaqGroups.length === 0" class="text-center q-py-xl text-grey-6">
            <q-icon name="find_in_page" size="xl" color="grey-5" />
            <div class="text-h6 q-mt-sm">Pertanyaan tidak ditemukan</div>
            <p class="text-caption">Coba gunakan kata kunci lain atau hubungi CS kami.</p>
            <q-btn
              color="positive"
              icon="chat"
              label="Tanya via WhatsApp"
              class="q-mt-md text-weight-bold"
              @click="hubungiWhatsApp"
            />
          </div>
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'

const searchQuery = ref('')

// Data Master FAQ Sekolah
const faqData = ref([
  {
    categoryName: 'Seputar PPDB & Pendaftaran',
    icon: 'assignment',
    items: [
      {
        question: 'Bagaimana alur pendaftaran siswa baru secara online?',
        answer:
          'Alur pendaftaran terdiri dari 3 langkah: 1) Mengklik tombol <b>Daftar PPDB</b> di halaman utama, 2) Mengisi formulir berkas secara lengkap, 3) Melakukan verifikasi pembayaran registrasi melalui sistem admin.',
      },
      {
        question: 'Apa saja berkas fisik yang harus dipersiapkan saat verifikasi?',
        answer:
          'Berkas yang wajib dibawa meliputi: Fotokopi Ijazah/Surat Keterangan Lulus (SKL), Kartu Keluarga (KK), Akta Kelahiran, dan Pas Foto terbaru ukuran 3x4 (2 lembar).',
      },
    ],
  },
  {
    categoryName: 'Rincian Biaya & Beasiswa',
    icon: 'payments',
    items: [
      {
        question: 'Apakah rincian biaya masuk bisa dicicil?',
        answer:
          'Ya, SMK Pasundan Jatinangor memberikan kemudahan pembayaran Dana Pengembangan Sekolah (Uang Gedung) yang dapat dicicil hingga 3 kali selama semester pertama berjalan.',
      },
      {
        question: 'Apakah tersedia program beasiswa untuk siswa berprestasi?',
        answer:
          'Kami menyediakan 3 jalur beasiswa: Beasiswa Akademik (Peringkat Rapor), Beasiswa Non-Akademik (Juara Olahraga/Seni minimal tingkat kota), dan Beasiswa Afiliasi bagi keluarga kurang mampu.',
      },
    ],
  },
  {
    categoryName: 'Kompetensi Keahlian (TBSM)',
    icon: 'motorcycle',
    items: [
      {
        question: 'Apa keunggulan Jurusan TBSM di SMK ini?',
        answer:
          'Jurusan Teknik & Bisnis Sepeda Motor (TBSM) kami telah terstandarisasi industri resmi dan didukung oleh alat praktek modern serta kesempatan magang langsung di jaringan bengkel resmi.',
      },
    ],
  },
])

// Logika Filter Pencarian Teks secara Real-time tingkat Kelompok dan Item
const filteredFaqGroups = computed(() => {
  if (!searchQuery.value) return faqData.value

  const query = searchQuery.value.toLowerCase()

  return (
    faqData.value
      .map((group) => {
        // Filter item di dalam grup yang cocok dengan kata kunci pertanyaan atau jawaban
        const matchedItems = group.items.filter(
          (item) =>
            item.question.toLowerCase().includes(query) ||
            item.answer.toLowerCase().includes(query),
        )
        // Kembalikan struktur grup baru hanya dengan item yang cocok
        return { ...group, items: matchedItems }
      })
      // Hanya tampilkan grup yang memiliki minimal 1 item yang cocok
      .filter((group) => group.items.length > 0)
  )
})

// Fungsi Darurat Direct Chat WhatsApp Humas Sekolah
function hubungiWhatsApp() {
  window.open('https://wa.me...', '_blank')
}
</script>

<style scoped>
.max-width-center {
  max-width: 1100px;
  margin: 0 auto;
}

.line-height-relaxed {
  line-height: 1.6;
}

/* Efek border halus antar item akordion */
.separator .q-expansion-item {
  border-bottom: 1px solid #e0e0e0;
}

.separator .q-expansion-item:last-child {
  border-bottom: none;
}
</style>
