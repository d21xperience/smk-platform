/*
 * @sdp/context — Domain Model: School
 *
 * Merepresentasikan entitas Sekolah dalam konteks operasional.
 * Digunakan oleh seluruh aplikasi SDP (siakad-tu, siakad-guru, psb, website).
 *
 * Pure JavaScript — TIDAK bergantung pada Vue, Pinia, Quasar, Axios.
 */
export class School {
  /*
   * @param {Object} params
   * @param {string} params.id - Unique identifier
   * @param {string} params.name - Nama sekolah
   * @param {string} params.npsn - Nomor Pokok Sekolah Nasional
   * @param {string} params.address - Alamat sekolah
   * @param {string|null} params.logoUrl - URL logo sekolah (opsional)
   * @param {boolean} params.isActive - Status aktif
   */
  constructor({ id, name, npsn, address, logoUrl, isActive }) {
    if (!id) throw new Error('School ID is required')

    this.id = id
    this.name = name||''
    this.npsn = npsn||''
    this.address = address||''
    this.logoUrl = logoUrl||null
    this.isActive = isActive !== false
  }
  /*
   * Mendapatkan label display untuk sekolah.
   * @returns {string}
   */
  getDisplayLabel() {
    return this.name
  }
  /*
   * Serialisasi ke plain object.
   * @returns {Object}
   */
  toJSON() {
    return {
      id: this.id,
      name: this.name,
      npsn: this.npsn,
      address: this.address,
      logoUrl: this.logoUrl,
      isActive: this.isActive
    }
  }
  /*
   * Deserialisasi dari plain object.
   * @param {Object} data
   * @returns {School}
   */
  static fromJSON(data) {
    if (!data) return null
    return new School({
      id: data.id,
      name: data.name,
      npsn: data.npsn,
      address: data.address,
      logoUrl: data.logoUrl,
      isActive: data.isActive
    })
  }
}
