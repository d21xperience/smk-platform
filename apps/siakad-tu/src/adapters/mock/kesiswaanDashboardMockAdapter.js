import { mockStorage } from './mockStorage.js'
import { dashboardMockData } from './data/dashboardKesiswaanData.js'

export const kesiswaanDashboardMockAdapter = {
  /**
   * Simulasi latency untuk meniru network delay
   * @param {number} ms - milidetik
   */
  async _simulateLatency(ms = 150) {
    return new Promise((resolve) => setTimeout(resolve, ms))
  },

  /**
   * Helper untuk response error
   */
  _error(code, message, details = null) {
    return {
      success: false,
      data: null,
      error: { code, message, details },
    }
  },

  /**
   * Helper untuk response sukses
   */
  _success(data) {
    return { success: true, data, error: null }
  },

  /**
   * Memuat data dashboard kesiswaan berdasarkan context
   * @param {OperationalContext} context - Harus memiliki schoolId yang valid
   * @returns {Promise<{ success: boolean, data?: object, error?: object }>}
   */
  async loadDashboard(context) {
    await this._simulateLatency()

    // Validasi context
    if (!context || !context.schoolId) {
      return this._error('INVALID_CONTEXT', 'Context tidak valid atau schoolId tidak ditemukan.', {
        context,
      })
    }

    const { schoolId } = context

    try {
      // 1. Ambil data siswa dari mockStorage
      const students = mockStorage.read(schoolId, 'students') || []

      // 2. Hitung total siswa aktif (asumsi status === 'ACTIVE')
      const totalSiswaAktif = students.filter((s) => s.status === 'ACTIVE').length

      // 3. Gabungkan dengan data statis dari dashboardMockData
      //    Data statis yang akan tetap digunakan:
      //    - totalPelanggaranBulanIni
      //    - totalBerprestasi
      //    - totalEkskul
      //    - siswaPerhatian
      //    - trenPelanggaran
      //    - poinPerTingkat
      //    - kasusTerbaru
      //    - ekstrakurikuler
      //    - prestasiList
      //    - agendaList
      const {
        totalPelanggaranBulanIni,
        totalBerprestasi,
        totalEkskul,
        siswaPerhatian,
        trenPelanggaran,
        poinPerTingkat,
        kasusTerbaru,
        ekstrakurikuler,
        prestasiList,
        agendaList,
      } = dashboardMockData

      // 4. Bentuk response sesuai contract
      const dashboardData = {
        totalSiswaAktif,
        totalPelanggaranBulanIni,
        totalBerprestasi,
        totalEkskul,
        siswaPerhatian,
        trenPelanggaran,
        poinPerTingkat,
        kasusTerbaru,
        ekstrakurikuler,
        prestasiList,
        agendaList,
      }

      return this._success(dashboardData)
    } catch (error) {
      return this._error('LOAD_DASHBOARD_FAILED', 'Gagal memuat data dashboard.', {
        error: error.message,
      })
    }
  },
}
