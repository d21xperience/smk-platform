import { StorageInterface } from "../abstractions/StorageInterface.js";
/*
 * @sdp/storage — Adapter: SessionStorageAdapter
 *
 * Abstraksi untuk browser sessionStorage.
 * Data hilang ketika browser tab/window ditutup.
 *
 * Cocok untuk:
 * - Session-specific data
 * - Temporary state
 * - Data yang tidak perlu persist setelah browser ditutup
 *
 * Pure JavaScript — TIDAK bergantung pada Vue, Pinia, Quasar, Axios.
 */
export class SessionStorageAdapter extends StorageInterface {
  /*
   * @param {Object} options
   * @param {string} options.prefix - Prefix untuk semua keys
   */
  constructor({ prefix = "sdp" } = {}) {
    super();
    this._prefix = prefix;
    this._available = this._checkAvailability();
  }
  /*
   * Memeriksa apakah sessionStorage tersedia.
   * @private
   * @returns {boolean}
   */
  _checkAvailability() {
    try {
      const testKey = "__sdp_storage_test__";
      window.sessionStorage.setItem(testKey, "test");
      window.sessionStorage.removeItem(testKey);
      return true;
    } catch (e) {
      console.warn(
        "[SessionStorageAdapter] sessionStorage is not available:",
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
      const raw = window.sessionStorage.getItem(fullKey);
      if (raw === null) return null;
      return JSON.parse(raw);
    } catch (e) {
      console.error(
        `[SessionStorageAdapter] Failed to get item "${key}":`,
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
      window.sessionStorage.setItem(fullKey, serialized);
      return true;
    } catch (e) {
      if (e.name === "QuotaExceededError" || e.code === 22) {
        console.error(
          `[SessionStorageAdapter] Quota exceeded for key "${key}"`,
        );
      } else {
        console.error(
          `[SessionStorageAdapter] Failed to set item "${key}":`,
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
      window.sessionStorage.removeItem(fullKey);
      return true;
    } catch (e) {
      console.error(
        `[SessionStorageAdapter] Failed to remove item "${key}":`,
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
        const keysToRemove = [];
        for (let i = 0; i < window.sessionStorage.length; i++) {
          const key = window.sessionStorage.key(i);
          if (key && key.startsWith(`${this._prefix}:`)) {
            keysToRemove.push(key);
          }
        }
        keysToRemove.forEach((key) => window.sessionStorage.removeItem(key));
      } else {
        window.sessionStorage.clear();
      }
    } catch (e) {
      console.error("[SessionStorageAdapter] Failed to clear:", e.message);
    }
  }
  /*
   * @override
   */
  hasItem(key) {
    if (!this._available) return false;
    try {
      const fullKey = this._buildKey(key);
      return window.sessionStorage.getItem(fullKey) !== null;
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
      for (let i = 0; i < window.sessionStorage.length; i++) {
        const key = window.sessionStorage.key(i);
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
    return "sessionStorage";
  }
}
