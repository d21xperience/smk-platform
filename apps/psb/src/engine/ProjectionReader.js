/**
 * Projection Reader interface.
 * Bertugas mengambil data dari Read Store berdasarkan nama proyeksi.
 * Implementasi konkret akan diinjeksi ke Reporting Engine.
 */
export class ProjectionReader {
  /**
   * @param {Object} readStores - map dari nama proyeksi ke store/getter
   * e.g., { 'classAttendance': () => statisticsCacheStore.classAttendance, ... }
   */
  constructor(readStores) {
    this.readStores = readStores
  }

  /**
   * Baca data dari proyeksi tertentu.
   * @param {string} projectionName - nama proyeksi di ReportDefinition
   * @returns {Array|Object} data mentah dari proyeksi
   */
  read(projectionName) {
    const getter = this.readStores[projectionName]
    if (!getter) throw new Error(`Projection '${projectionName}' tidak ditemukan`)
    return getter()
  }
}
