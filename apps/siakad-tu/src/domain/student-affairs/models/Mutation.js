export class Mutation {
  constructor(data) {
    this.id = data.id
    this.schoolId = data.schoolId
    this.studentId = data.studentId || null // Null jika mutasi masuk (siswa belum ada di DB)
    this.jenis = data.jenis // 'MASUK' | 'KELUAR'
    this.namaCalon = data.namaCalon || null // Untuk mutasi masuk
    this.nisnCalon = data.nisnCalon || null
    this.tujuanSekolah = data.tujuanSekolah || null // Untuk mutasi keluar
    this.alasan = data.alasan
    this.dokumenUrl = data.dokumenUrl || null
    this.status = data.status // 'PENDING', 'DISETUJUI', 'DITOLAK'
    this.catatanVerifikasi = data.catatanVerifikasi || null
    this.academicYearId = data.academicYearId
    this.tanggalPengajuan = data.tanggalPengajuan
  }

  isPending() {
    return this.status === 'PENDING'
  }
}
