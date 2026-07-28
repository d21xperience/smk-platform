/**
 * Shared validators untuk form validation
 * Kompatibel dengan Quasar validation rules
 */

export const required = (val) => !!val || "Field ini wajib diisi";

export const email = (val) => {
  const pattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return pattern.test(val) || "Format email tidak valid";
};

export const minLength = (min) => (val) =>
  (val && val.length >= min) || `Minimal ${min} karakter`;

export const maxLength = (max) => (val) =>
  (val && val.length <= max) || `Maksimal ${max} karakter`;

export const numeric = (val) => !isNaN(val) || "Hanya angka yang diperbolehkan";

export const phoneNumber = (val) => {
  const pattern = /^[0-9+\-\s()]{10,20}$/;
  return pattern.test(val) || "Format nomor telepon tidak valid";
};

export const nis = (val) => {
  const pattern = /^[0-9]{5,10}$/;
  return pattern.test(val) || "NIS harus 5-10 digit angka";
};

export const nisn = (val) => {
  const pattern = /^[0-9]{10}$/;
  return pattern.test(val) || "NISN harus 10 digit angka";
};
