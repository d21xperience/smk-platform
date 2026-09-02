export class Student {
  constructor(data) {
    this.id = data.id;
    this.schoolId = data.schoolId;
    this.dapodikId = data.dapodikId || null;
    this.nis = data.nis;
    this.nisn = data.nisn;
    this.nama = data.nama;
    this.tempatLahir = data.tempatLahir;
    this.tanggalLahir = data.tanggalLahir;
    this.jenisKelamin = data.jenisKelamin;
    this.agama = data.agama;
    this.alamat = data.alamat;
    this.noHp = data.noHp;
    this.email = data.email;
    this.asalSekolah = data.asalSekolah;
    this.jurusan = data.jurusan;
    this.kelas = data.kelas;
    this.status = data.status; // 'AKTIF', 'MUTASI_KELUAR', 'LULUS', 'DO'
    this.academicYearId = data.academicYearId;
    this.isLocked = data.isLocked || false; // True jika berasal dari Dapodik
    this.createdAt = data.createdAt;
    this.updatedAt = data.updatedAt;
  }

  isActive() {
    return this.status === 'AKTIF';
  }
}
