// apps/siakad-tu/src/domain/student/value-objects/NISN.js

/**
 * Value Object: NISN (Nomor Induk Siswa Nasional)
 *
 * Aturan:
 * - 10 digit angka
 * - Tidak boleh kosong
 * - Immutable setelah dibuat
 * - Self-validating
 */
export class NISN {
  #value;

  constructor(value) {
    this.#validate(value);
    this.#value = value;
    Object.freeze(this);
  }

  #validate(value) {
    if (!value) {
      throw new Error('NISN tidak boleh kosong.');
    }

    const strValue = String(value).trim();

    if (!/^\d{10}$/.test(strValue)) {
      throw new Error(
        `NISN harus terdiri dari 10 digit angka. Diterima: "${value}".`
      );
    }
  }

  get value() {
    return this.#value;
  }

  equals(other) {
    if (!(other instanceof NISN)) return false;
    return this.#value === other.#value;
  }

  toString() {
    return this.#value;
  }

  toJSON() {
    return this.#value;
  }

  static fromJSON(value) {
    return new NISN(value);
  }
}
