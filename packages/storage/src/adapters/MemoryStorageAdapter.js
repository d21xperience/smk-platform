import { StorageInterface } from "../abstractions/StorageInterface.js";
/*
 * @sdp/storage — Adapter: MemoryStorageAdapter
 *
 * In-memory storage implementation.
 * Data hilang ketika aplikasi/page di-reload.
 *
 * Cocok untuk:
 * - Unit testing (tidak perlu browser)
 * - SSR (Server-Side Rendering)
 * - Fallback ketika localStorage/sessionStorage tidak tersedia
 * - Temporary state yang tidak perlu persist
 *
 * Pure JavaScript — TIDAK bergantung pada Vue, Pinia, Quasar, Axios, atau browser API.
 */
export class MemoryStorageAdapter extends StorageInterface {
  /*
   * @param {Object} options
   * @param {string} options.prefix - Prefix untuk semua keys (untuk konsistensi API)
   */
  constructor({ prefix = "sdp" } = {}) {
    super();
    this._prefix = prefix;
    this._store = new Map();
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
    try {
      const fullKey = this._buildKey(key);
      const value = this._store.get(fullKey);
      if (value === undefined) return null;
      return JSON.parse(value);
    } catch (e) {
      console.error(
        `[MemoryStorageAdapter] Failed to get item "${key}":`,
        e.message,
      );
      return null;
    }
  }
  /*
   * @override
   */
  setItem(key, value) {
    try {
      const fullKey = this._buildKey(key);
      const serialized = JSON.stringify(value);
      this._store.set(fullKey, serialized);
      return true;
    } catch (e) {
      console.error(
        `[MemoryStorageAdapter] Failed to set item "${key}":`,
        e.message,
      );
      return false;
    }
  }
  /*
   * @override
   */
  removeItem(key) {
    try {
      const fullKey = this._buildKey(key);
      return this._store.delete(fullKey);
    } catch (e) {
      console.error(
        `[MemoryStorageAdapter] Failed to remove item "${key}":`,
        e.message,
      );
      return false;
    }
  }
  /*
   * @override
   */
  clear() {
    if (this._prefix) {
      const keysToRemove = [];
      for (const key of this._store.keys()) {
        if (key.startsWith(`${this._prefix}:`)) {
          keysToRemove.push(key);
        }
      }
      keysToRemove.forEach((key) => this._store.delete(key));
    } else {
      this._store.clear();
    }
  }
  /*
   * @override
   */
  hasItem(key) {
    const fullKey = this._buildKey(key);
    return this._store.has(fullKey);
  }
  /*
   * @override
   */
  keys() {
    const allKeys = [];
    for (const key of this._store.keys()) {
      if (this._prefix) {
        if (key.startsWith(`${this._prefix}:`)) {
          allKeys.push(key.substring(this._prefix.length + 1));
        }
      } else {
        allKeys.push(key);
      }
    }
    return allKeys;
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
    return "memory";
  }
}
