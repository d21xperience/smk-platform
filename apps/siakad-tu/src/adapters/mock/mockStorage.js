// apps/siakad-tu/src/adapters/mock/mockStorage.js

/**
 * MockStorage — Wrapper localStorage dengan scoping per school.
 *
 * Struktur key:
 * `sdp_mock:{schoolId}:{domain}` → JSON data
 *
 * Contoh:
 * `sdp_mock:smk-pasundan-001:students` → array of student data
 *
 * Scoping ini memastikan data antar sekolah terisolasi,
 * sesuai dengan Operational Context.
 */
export class MockStorage {
  #prefix = 'sdp_mock'

  /**
   * Build key dengan scoping school
   * @param {string} schoolId
   * @param {string} domain - 'students', 'enrollments', dll
   * @returns {string}
   */
  _buildKey(schoolId, domain) {
    if (!schoolId) throw new Error('MockStorage: schoolId wajib untuk scoping.')
    return `${this.#prefix}:${schoolId}:${domain}`
  }

  /**
   * Baca data dari storage
   * @param {string} schoolId
   * @param {string} domain
   * @returns {Array} Array of data (empty array jika belum ada)
   */
  read(schoolId, domain) {
    try {
      const key = this._buildKey(schoolId, domain)
      const raw = localStorage.getItem(key)
      if (!raw) return []
      return JSON.parse(raw)
    } catch (error) {
      console.error(`[MockStorage] Error membaca ${domain}:`, error)
      return []
    }
  }

  /**
   * Tulis data ke storage
   * @param {string} schoolId
   * @param {string} domain
   * @param {Array} data
   */
  write(schoolId, domain, data) {
    try {
      const key = this._buildKey(schoolId, domain)
      localStorage.setItem(key, JSON.stringify(data))
    } catch (error) {
      console.error(`[MockStorage] Error menulis ${domain}:`, error)
      throw new Error(`Gagal menyimpan data ${domain}.`)
    }
  }

  /**
   * Hapus semua data untuk school tertentu
   * @param {string} schoolId
   */
  clearSchool(schoolId) {
    if (!schoolId) {
      // Hapus semua data mock
      Object.keys(localStorage)
        .filter((k) => k.startsWith(this.#prefix))
        .forEach((k) => localStorage.removeItem(k))
      return
    }

    Object.keys(localStorage)
      .filter((k) => k.startsWith(`${this.#prefix}:${schoolId}:`))
      .forEach((k) => localStorage.removeItem(k))
  }

  /**
   * Reset semua data mock (untuk testing)
   */
  resetAll() {
    this.clearSchool(null)
  }

  /**
   * Cek apakah storage kosong untuk domain tertentu
   * @param {string} schoolId
   * @param {string} domain
   * @returns {boolean}
   */
  isEmpty(schoolId, domain) {
    return this.read(schoolId, domain).length === 0
  }
}

// Singleton instance
export const mockStorage = new MockStorage()
