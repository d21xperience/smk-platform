import { StorageInterface } from "../abstractions/StorageInterface.js";
/*
 * @sdp/storage — Adapter: LocalStorageAdapter
 *
 * Abstraksi untuk browser localStorage.
 * Menyediakan interface yang konsisten dengan automatic JSON serialization.
 *
 * Fitur:
 * - Automatic JSON serialize/deserialize
 * - Error handling untuk quota exceeded
 * - Prefix support untuk namespace isolation
 * - Graceful fallback jika localStorage tidak tersedia
 *
 * Pure JavaScript — TIDAK bergantung pada Vue, Pinia, Quasar, Axios.
 */
export class LocalStorageAdapter extends StorageInterface {
  /*
   * @param {Object} options
   * @param {string} options.prefix - Prefix untuk semua keys (namespace isolation)
   */
  constructor({ prefix = "sdp" } = {}) {
    super();
    this._prefix = prefix;
    this._available = this._checkAvailability();
  }
  /*
   * Memeriksa apakah localStorage tersedia.
   * @private
   * @returns {boolean}
   */
  _checkAvailability() {
    try {
      const testKey = "__sdp_storage_test__";
      window.localStorage.setItem(testKey, "test");
      window.localStorage.removeItem(testKey);
      return true;
    } catch (e) {
      console.warn(
        "[LocalStorageAdapter] localStorage is not available:",
        e.message,
      );
      return false;
    }
  }
  /*
   * Membangun full key dengan prefix.
   * @private
   * @param {string} key
   * @returns {string}
   */
  _buildKey(key) {
    return this._prefix ? `${this._prefix}:${key}` : key;
  }
  /*
   * @override
   */
  getItem(key) {
    if (!this._available) return null;
    try {
      const fullKey = this._buildKey(key);
      const raw = window.localStorage.getItem(fullKey);
      if (raw === null) return null;
      return JSON.parse(raw);
    } catch (e) {
      console.error(
        `[LocalStorageAdapter] Failed to get item "${key}":`,
        e.message,
      );
      return null;
    }
  }
  /*
   * @override
   */
  setItem(key, value) {
    if (!this._available) return false;
    try {
      const fullKey = this._buildKey(key);
      const serialized = JSON.stringify(value);
      window.localStorage.setItem(fullKey, serialized);
      return true;
    } catch (e) {
      if (e.name === "QuotaExceededError" || e.code === 22) {
        console.error(`[LocalStorageAdapter] Quota exceeded for key "${key}"`);
      } else {
        console.error(
          `[LocalStorageAdapter] Failed to set item "${key}":`,
          e.message,
        );
      }
      return false;
    }
  }
  /*
   * @override
   */
  removeItem(key) {
    if (!this._available) return false;
    try {
      const fullKey = this._buildKey(key);
      window.localStorage.removeItem(fullKey);
      return true;
    } catch (e) {
      console.error(
        `[LocalStorageAdapter] Failed to remove item "${key}":`,
        e.message,
      );
      return false;
    }
  }
  /*
   * @override
   */
  clear() {
    if (!this._available) return;
    try {
      if (this._prefix) {
        // Only clear items with our prefix
        const keysToRemove = [];
        for (let i = 0; i < window.localStorage.length; i++) {
          const key = window.localStorage.key(i);
          if (key && key.startsWith(`${this._prefix}:`)) {
            keysToRemove.push(key);
          }
        }
        keysToRemove.forEach((key) => window.localStorage.removeItem(key));
      } else {
        window.localStorage.clear();
      }
    } catch (e) {
      console.error("[LocalStorageAdapter] Failed to clear:", e.message);
    }
  }
  /*
   * @override
   */
  hasItem(key) {
    if (!this._available) return false;
    try {
      const fullKey = this._buildKey(key);
      return window.localStorage.getItem(fullKey) !== null;
    } catch (e) {
      return false;
    }
  }
  /*
   * @override
   */
  keys() {
    if (!this._available) return [];
    try {
      const allKeys = [];
      for (let i = 0; i < window.localStorage.length; i++) {
        const key = window.localStorage.key(i);
        if (this._prefix) {
          if (key && key.startsWith(`${this._prefix}:`)) {
            allKeys.push(key.substring(this._prefix.length + 1));
          }
        } else {
          allKeys.push(key);
        }
      }
      return allKeys;
    } catch (e) {
      return [];
    }
  }
  /*
   * @override
   */
  get length() {
    return this.keys().length;
  }
  /*
   * @override
   */
  get storageType() {
    return "localStorage";
  }
}
