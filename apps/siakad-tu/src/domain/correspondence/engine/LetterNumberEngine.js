export class LetterNumberEngine {
  /**
   * Format nomor surat berdasarkan kategori, urutan, kode sekolah, dan tahun.
   * @param {string} categoryCode - Kode kategori (misal: '422.1', '800')
   * @param {number} sequence - Nomor urut
   * @param {string} schoolCode - Kode sekolah (misal: 'SMK-01')
   * @param {number} year - Tahun (misal: 2026)
   * @returns {string} Nomor surat yang terformat
   */
  static format(categoryCode, sequence, schoolCode, year) {
    const paddedSeq = String(sequence).padStart(3, '0')
    return `${categoryCode}/${paddedSeq}/${schoolCode}/${year}`
  }

  /**
   * Validasi format nomor surat.
   * @param {string} letterNumber - Nomor surat yang akan divalidasi
   * @returns {boolean} True jika format valid
   */
  static isValidFormat(letterNumber) {
    // Format: kode/urut(3 digit)/kodeSekolah/tahun(4 digit)
    // Contoh: 422.1/001/SMK-01/2026
    const regex = /^[\d.]+\/\d{3}\/[A-Z0-9-]+\/\d{4}$/
    return regex.test(letterNumber)
  }
}
