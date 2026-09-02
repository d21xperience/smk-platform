/*
 * @sdp/storage — Abstraction: StorageInterface
 *
 * Mendefinisikan contract/interface untuk seluruh storage adapter.
 * Semua adapter (LocalStorage, SessionStorage, Memory) WAJIB
 * mengimplementasikan interface ini.
 *
 * Pure JavaScript — TIDAK bergantung pada Vue, Pinia, Quasar, Axios.
 */
/*
 * @typedef {Object} StorageAdapter
 * @property {Function} getItem - Mengambil value berdasarkan key
 * @property {Function} setItem - Menyimpan value berdasarkan key
 * @property {Function} removeItem - Menghapus value berdasarkan key
 * @property {Function} clear - Menghapus semua data
 * @property {Function} hasItem - Memeriksa apakah key ada
 * @property {Function} keys - Mendapatkan semua keys
 * @property {Function} length - Mendapatkan jumlah items
 */
export class StorageInterface {
  /*
   * Mengambil value berdasarkan key.
   * @param {string} key
   * @returns {*} Parsed value atau null jika tidak ada
   */
  getItem(key) {
    throw new Error(
      "StorageInterface.getItem() must be implemented by adapter",
    );
  }
  /*
   * Menyimpan value berdasarkan key.
   * @param {string} key
   * @param {*} value - Value akan di-serialize ke JSON
   * @returns {boolean} True jika berhasil
   */
  setItem(key, value) {
    throw new Error(
      "StorageInterface.setItem() must be implemented by adapter",
    );
  }
  /*
   * Menghapus value berdasarkan key.
   * @param {string} key
   * @returns {boolean} True jika berhasil
   */
  removeItem(key) {
    throw new Error(
      "StorageInterface.removeItem() must be implemented by adapter",
    );
  }
  /*
   * Menghapus semua data dari storage.
   * @returns {void}
   */
  clear() {
    throw new Error("StorageInterface.clear() must be implemented by adapter");
  }
  /*
   * Memeriksa apakah key ada di storage.
   * @param {string} key
   * @returns {boolean}
   */
  hasItem(key) {
    throw new Error(
      "StorageInterface.hasItem() must be implemented by adapter",
    );
  }
  /*
   * Mendapatkan semua keys yang tersimpan.
   * @returns {Array<string>}
   */
  keys() {
    throw new Error("StorageInterface.keys() must be implemented by adapter");
  }
  /*
   * Mendapatkan jumlah items yang tersimpan.
   * @returns {number}
   */
  get length() {
    throw new Error("StorageInterface.length must be implemented by adapter");
  }
  /*
   * Mendapatkan nama storage type.
   * @returns {string}
   */
  get storageType() {
    throw new Error(
      "StorageInterface.storageType must be implemented by adapter",
    );
  }
}
