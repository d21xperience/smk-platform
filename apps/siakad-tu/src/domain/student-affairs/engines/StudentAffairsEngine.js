export class StudentAffairsEngine {
  /**
   * Memvalidasi dan menghasilkan command yang harus dieksekusi Service
   * saat mutasi disetujui.
   */
  static processMutationApproval(mutation) {
    if (mutation.jenis === 'MASUK') {
      // Business Rule: Mutasi masuk yang disetujui harus membuat entitas Student baru
      return {
        action: 'CREATE_STUDENT',
        payload: {
          nis: `M${Date.now()}`, // Generate NIS sementara
          nisn: mutation.nisnCalon,
          nama: mutation.namaCalon,
          asalSekolah: mutation.asalSekolah,
          status: 'AKTIF',
          academicYearId: mutation.academicYearId,
          schoolId: mutation.schoolId
        }
      };
    }

    if (mutation.jenis === 'KELUAR') {
      // Business Rule: Mutasi keluar yang disetujui harus mengubah status siswa
      return {
        action: 'UPDATE_STUDENT_STATUS',
        payload: {
          studentId: mutation.studentId,
          newStatus: 'MUTASI_KELUAR',
          alasanPindah: mutation.alasan,
          tujuanSekolah: mutation.tujuanSekolah
        }
      };
    }

    throw new Error('Jenis mutasi tidak dikenali');
  }

  static validateMutationSubmission(data) {
    if (!data.jenis) throw new Error('Jenis mutasi wajib diisi');
    if (!data.alasan) throw new Error('Alasan mutasi wajib diisi');

    if (data.jenis === 'MASUK') {
      if (!data.namaCalon || !data.nisnCalon) {
        throw new Error('Data calon siswa (Nama, NISN) wajib diisi untuk mutasi masuk');
      }
    } else if (data.jenis === 'KELUAR') {
      if (!data.studentId || !data.tujuanSekolah) {
        throw new Error('Siswa dan tujuan sekolah wajib dipilih untuk mutasi keluar');
      }
    }
    return true;
  }
}
