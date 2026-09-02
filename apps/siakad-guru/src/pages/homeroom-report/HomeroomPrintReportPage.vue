<!-- FILE: src/pages/homeroom-report/HomeroomPrintReportPage.vue -->
<!-- STATUS: NEW -->
<!-- STATUS IMPLEMENTASI: COMPLETE -->

<template>
  <q-page class="q-pa-md">
    <div class="no-print row items-center justify-between q-mb-md">
      <div class="text-h5">Laporan Wali Kelas</div>
      <q-btn color="primary" icon="print" label="Cetak Laporan" :loading="isAllLoading" @click="printReport" />
    </div>

    <q-banner v-if="!isAllLoading && students.length === 0" class="bg-warning text-dark q-mb-md no-print">
      Tidak ada data untuk ditampilkan. Pastikan context sudah dipilih dan Anda adalah wali kelas.
    </q-banner>

    <div v-if="isAllLoading" class="flex flex-center no-print" style="min-height: 200px;">
      <q-spinner-dots size="60px" color="primary" />
    </div>

    <div v-if="!isAllLoading" class="print-report">
      <div class="report-header text-center q-mb-lg">
        <h2 class="report-title">LAPORAN WALI KELAS</h2>
        <h3 class="report-subtitle">{{ reportHeader.schoolName }}</h3>
        <div class="report-meta">
          <span>Tahun Ajaran: {{ reportHeader.academicYear }}</span>
          <span class="q-mx-md">|</span>
          <span>Semester: {{ reportHeader.semester }}</span>
        </div>
        <div class="report-meta">
          <span>Kelas: {{ reportHeader.className }}</span>
          <span class="q-mx-md">|</span>
          <span>Wali Kelas: {{ reportHeader.teacherName }}</span>
        </div>
      </div>

      <div class="report-section">
        <h4 class="section-title">I. DATA SISWA ({{ students.length }} siswa)</h4>
        <table class="report-table">
          <thead>
            <tr>
              <th class="text-center" style="width: 40px;">No</th>
              <th style="width: 120px;">NISN</th>
              <th style="width: 80px;">NIS</th>
              <th>Nama Lengkap</th>
              <th class="text-center" style="width: 40px;">JK</th>
              <th class="text-center" style="width: 80px;">Status</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="student in students" :key="student.id">
              <td class="text-center">{{ student.seatNumber }}</td>
              <td>{{ student.nisn }}</td>
              <td>{{ student.nis }}</td>
              <td>{{ student.fullName }}</td>
              <td class="text-center">{{ student.gender }}</td>
              <td class="text-center">{{ student.status }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="report-section">
        <h4 class="section-title">II. REKAP NILAI</h4>
        <p v-if="progressSummary" class="summary-text">
          Rata-rata Kelas: {{ progressSummary.classAverage }} |
          Ketuntasan: {{ progressSummary.passRate }}% |
          Tertinggi: {{ progressSummary.highestAverage }} |
          Terendah: {{ progressSummary.lowestAverage }}
        </p>
        <table class="report-table">
          <thead>
            <tr>
              <th class="text-center" style="width: 40px;">No</th>
              <th>Nama Lengkap</th>
              <th class="text-center" style="width: 80px;">Rata-rata</th>
              <th class="text-center" style="width: 60px;">Rank</th>
              <th class="text-center" style="width: 60px;">Tuntas</th>
              <th class="text-center" style="width: 60px;">Belum</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="student in progressStudents" :key="student.studentId">
              <td class="text-center">{{ student.seatNumber }}</td>
              <td>{{ student.studentName }}</td>
              <td class="text-center">{{ student.overallAverage }}</td>
              <td class="text-center">{{ student.classRank }}</td>
              <td class="text-center">{{ student.passCount }}</td>
              <td class="text-center">{{ student.failCount }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="report-section">
        <h4 class="section-title">III. REKAP KEHADIRAN</h4>
        <p v-if="attendanceSummary" class="summary-text">
          Rata-rata Kehadiran: {{ attendanceSummary.classAveragePercentage }}% |
          Baik: {{ attendanceSummary.goodCount }} |
          Perhatian: {{ attendanceSummary.warningCount }} |
          Kritis: {{ attendanceSummary.criticalCount }}
        </p>
        <table class="report-table">
          <thead>
            <tr>
              <th class="text-center" style="width: 40px;">No</th>
              <th>Nama Lengkap</th>
              <th class="text-center" style="width: 60px;">Hadir</th>
              <th class="text-center" style="width: 60px;">Sakit</th>
              <th class="text-center" style="width: 60px;">Izin</th>
              <th class="text-center" style="width: 60px;">Alpha</th>
              <th class="text-center" style="width: 80px;">%</th>
              <th class="text-center" style="width: 80px;">Status</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="student in attendanceStudents" :key="student.studentId">
              <td class="text-center">{{ student.seatNumber }}</td>
              <td>{{ student.studentName }}</td>
              <td class="text-center">{{ student.semesterTotal.present }}</td>
              <td class="text-center">{{ student.semesterTotal.sick }}</td>
              <td class="text-center">{{ student.semesterTotal.permitted }}</td>
              <td class="text-center">{{ student.semesterTotal.absent }}</td>
              <td class="text-center">{{ student.semesterTotal.percentage }}%</td>
              <td class="text-center">{{ getAttendanceStatusLabel(student.status) }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="report-section">
        <h4 class="section-title">IV. CATATAN SISWA ({{ notes.length }} catatan)</h4>
        <table class="report-table">
          <thead>
            <tr>
              <th class="text-center" style="width: 90px;">Tanggal</th>
              <th style="width: 140px;">Siswa</th>
              <th class="text-center" style="width: 80px;">Kategori</th>
              <th>Judul</th>
              <th class="text-center" style="width: 80px;">Status</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="note in notes" :key="note.noteId">
              <td class="text-center">{{ note.noteDate }}</td>
              <td>{{ note.studentName }}</td>
              <td class="text-center">{{ getNoteCategoryLabel(note.category) }}</td>
              <td>{{ note.title }}</td>
              <td class="text-center">{{ note.status === 'SUBMITTED' ? 'Final' : 'Draft' }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="report-section">
        <h4 class="section-title">V. LOG KOMUNIKASI ORANG TUA ({{ communications.length }} log)</h4>
        <table class="report-table">
          <thead>
            <tr>
              <th class="text-center" style="width: 90px;">Tanggal</th>
              <th style="width: 140px;">Siswa</th>
              <th class="text-center" style="width: 90px;">Kategori</th>
              <th class="text-center" style="width: 70px;">Arah</th>
              <th>Judul</th>
              <th class="text-center" style="width: 90px;">Status</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="comm in communications" :key="comm.communicationId">
              <td class="text-center">{{ comm.communicationDate }}</td>
              <td>{{ comm.studentName }}</td>
              <td class="text-center">{{ getCommCategoryLabel(comm.category) }}</td>
              <td class="text-center">{{ getCommDirectionLabel(comm.direction) }}</td>
              <td>{{ comm.title }}</td>
              <td class="text-center">{{ getCommStatusLabel(comm.status) }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="report-signature">
        <div class="signature-block">
          <p>Mengetahui,</p>
          <p>Kepala Sekolah</p>
          <div class="signature-space"></div>
          <p>_________________________</p>
          <p>NIP.</p>
        </div>
        <div class="signature-block">
          <p>&nbsp;</p>
          <p>Wali Kelas {{ reportHeader.className }}</p>
          <div class="signature-space"></div>
          <p>{{ reportHeader.teacherName }}</p>
          <p>NIP.</p>
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { onMounted } from 'vue'
import { useHomeroomPrintReport } from '@/composables/useHomeroomPrintReport.js'

const {
  students,
  progressSummary,
  progressStudents,
  attendanceSummary,
  attendanceStudents,
  notes,
  communications,
  reportHeader,
  isAllLoading,
  loadAllData,
  printReport,
} = useHomeroomPrintReport()

function getAttendanceStatusLabel(status) {
  switch (status) {
    case 'GOOD': return 'Baik'
    case 'WARNING': return 'Perhatian'
    case 'CRITICAL': return 'Kritis'
    default: return status
  }
}

function getNoteCategoryLabel(category) {
  switch (category) {
    case 'POSITIVE': return 'Positif'
    case 'NEGATIVE': return 'Negatif'
    case 'NEUTRAL': return 'Netral'
    case 'INCIDENT': return 'Insiden'
    default: return category
  }
}

function getCommCategoryLabel(category) {
  switch (category) {
    case 'PHONE_CALL': return 'Telepon'
    case 'WHATSAPP': return 'WhatsApp'
    case 'MEETING': return 'Pertemuan'
    case 'HOME_VISIT': return 'Kunjungan'
    case 'OTHER': return 'Lainnya'
    default: return category
  }
}

function getCommDirectionLabel(direction) {
  switch (direction) {
    case 'INBOUND': return 'Masuk'
    case 'OUTBOUND': return 'Keluar'
    case 'BIDIRECTIONAL': return 'Dua Arah'
    default: return direction
  }
}

function getCommStatusLabel(status) {
  switch (status) {
    case 'PENDING_FOLLOW_UP': return 'Pending'
    case 'FOLLOWED_UP': return 'Followed Up'
    case 'CLOSED': return 'Selesai'
    default: return status
  }
}

onMounted(() => {
  loadAllData()
})
</script>

<style>
.print-report {
  background: white;
  padding: 20px;
  max-width: 900px;
  margin: 0 auto;
}

.report-title {
  font-size: 20px;
  font-weight: bold;
  margin-bottom: 4px;
  text-transform: uppercase;
}

.report-subtitle {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 8px;
}

.report-meta {
  font-size: 13px;
  margin-bottom: 4px;
}

.section-title {
  font-size: 15px;
  font-weight: bold;
  margin-top: 24px;
  margin-bottom: 8px;
  border-bottom: 2px solid #333;
  padding-bottom: 4px;
}

.summary-text {
  font-size: 12px;
  margin-bottom: 8px;
  color: #555;
}

.report-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 12px;
  margin-bottom: 16px;
}

.report-table th,
.report-table td {
  border: 1px solid #333;
  padding: 6px 8px;
}

.report-table th {
  background-color: #f0f0f0;
  font-weight: bold;
}

.report-signature {
  display: flex;
  justify-content: space-between;
  margin-top: 48px;
  padding-top: 16px;
}

.signature-block {
  text-align: center;
  width: 250px;
}

.signature-space {
  height: 60px;
}

.signature-block p {
  margin: 4px 0;
  font-size: 13px;
}

@media print {

  .no-print,
  .q-header,
  .q-drawer,
  .q-footer {
    display: none !important;
  }

  .q-page-container {
    padding: 0 !important;
    margin: 0 !important;
  }

  .q-page {
    padding: 0 !important;
    margin: 0 !important;
  }

  .print-report {
    padding: 0;
    max-width: 100%;
  }

  .report-section {
    page-break-inside: avoid;
  }

  .report-table {
    page-break-inside: auto;
  }

  .report-table tr {
    page-break-inside: avoid;
    page-break-after: auto;
  }

  .report-signature {
    page-break-inside: avoid;
  }

  @page {
    size: A4;
    margin: 15mm;
  }
}
</style>
