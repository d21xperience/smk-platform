/**
 * Dashboard Kesiswaan Contract
 *
 * Mendefinisikan interface untuk mengambil data dashboard kesiswaan.
 * Semua adapter (mock/real) HARUS mengikuti contract ini.
 */

export const DashboardContract = {
  /**
   * Memuat data dashboard kesiswaan untuk konteks operasional tertentu.
   *
   * @param {OperationalContext} context - Konteks operasional (schoolId, periodId, dll)
   * @returns {Promise<{
   *   success: boolean,
   *   data?: {
   *     totalSiswaAktif: number,
   *     totalPelanggaranBulanIni: number,
   *     totalBerprestasi: number,
   *     totalEkskul: number,
   *     siswaPerhatian: Array<{
   *       nama: string,
   *       kelas: string,
   *       catatan: string,
   *       poin: number
   *     }>,
   *     trenPelanggaran: Array<{
   *       label: string,
   *       jumlah: number
   *     }>,
   *     poinPerTingkat: Array<{
   *       label: string,
   *       total: number,
   *       warna: string
   *     }>,
   *     kasusTerbaru: Array<{
   *       id: number,
   *       judul: string,
   *       siswa: string,
   *       kelas: string,
   *       waktu: string,
   *       kategori: 'Ringan' | 'Sedang' | 'Berat'
   *     }>,
   *     ekstrakurikuler: Array<{
   *       nama: string,
   *       pembina: string,
   *       peserta: number,
   *       icon: string,
   *       warna: string
   *     }>,
   *     prestasiList: Array<{
   *       judul: string,
   *       siswa: string,
   *       tanggal: string
   *     }>,
   *     agendaList: Array<{
   *       tanggal: string,
   *       bulan: string,
   *       judul: string,
   *       waktu: string
   *     }>
   *   },
   *   error?: {
   *     code: string,
   *     message: string,
   *     details?: any,
   *     source?: string
   *   }
   * }>}
   */
  loadDashboard: async () => {
    /**
     * Implementasi:
     * - MockAdapter: agregasi dari mockStorage + data statis
     * - RealAdapter: panggil endpoint API
     */
    console.warn('DashboardContract.loadDashboard() — implementasi belum di-override!')
  },
}
