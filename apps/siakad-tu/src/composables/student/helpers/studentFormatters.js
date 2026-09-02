// apps/siakad-tu/src/composables/student/helpers/studentFormatters.js

/**
 * Student Formatters — Pure functions untuk format data siswa ke tampilan UI.
 *
 * Ini adalah UI concern, BUKAN business concern.
 * Semua fungsi ini:
 * - Pure (no side effects)
 * - Testable tanpa Vue/Pinia
 * - Reusable di seluruh aplikasi
 */

/**
 * Format nama lengkap dari object fullName
 * @param {Object} fullName - { firstName, middleName, lastName }
 * @returns {string}
 */
export function formatFullName(fullName) {
  if (!fullName) return '-';
  return [fullName.firstName, fullName.middleName, fullName.lastName]
    .filter(Boolean)
    .join(' ')
    .trim() || '-';
}

/**
 * Format tanggal ISO ke format Indonesia (DD/MM/YYYY)
 * @param {string} isoString - ISO date string
 * @returns {string}
 */
export function formatBirthDate(isoString) {
  if (!isoString) return '-';
  try {
    const date = new Date(isoString);
    if (isNaN(date.getTime())) return '-';
    const day = String(date.getDate()).padStart(2, '0');
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const year = date.getFullYear();
    return `${day}/${month}/${year}`;
  } catch {
    return '-';
  }
}

/**
 * Hitung usia dari tanggal lahir
 * @param {string} isoString - ISO date string
 * @returns {number|null}
 */
export function calculateAge(isoString) {
  if (!isoString) return null;
  try {
    const birthDate = new Date(isoString);
    const now = new Date();
    if (isNaN(birthDate.getTime())) return null;

    let age = now.getFullYear() - birthDate.getFullYear();
    const monthDiff = now.getMonth() - birthDate.getMonth();
    if (monthDiff < 0 || (monthDiff === 0 && now.getDate() < birthDate.getDate())) {
      age--;
    }
    return age;
  } catch {
    return null;
  }
}

/**
 * Format usia untuk ditampilkan
 * @param {string} isoString
 * @returns {string}
 */
export function formatAge(isoString) {
  const age = calculateAge(isoString);
  return age !== null ? `${age} tahun` : '-';
}

/**
 * Format gender code ke label
 * @param {string} code - 'MALE' | 'FEMALE'
 * @returns {string}
 */
export function formatGender(code) {
  const labels = {
    MALE: 'Laki-laki',
    FEMALE: 'Perempuan'
  };
  return labels[code] || '-';
}

/**
 * Format gender code ke singkatan
 * @param {string} code
 * @returns {string}
 */
export function formatGenderShort(code) {
  const labels = { MALE: 'L', FEMALE: 'P' };
  return labels[code] || '-';
}

/**
 * Format student status code ke label
 * @param {string} code - 'ACTIVE' | 'GRADUATED' | 'TRANSFERRED' | 'DROPPED' | 'ALUMNI'
 * @returns {string}
 */
export function formatStudentStatus(code) {
  const labels = {
    ACTIVE: 'Aktif',
    GRADUATED: 'Lulus',
    TRANSFERRED: 'Pindah',
    DROPPED: 'Keluar',
    ALUMNI: 'Alumni'
  };
  return labels[code] || code || '-';
}

/**
 * Format student status ke color (untuk badge)
 * @param {string} code
 * @returns {string} Quasar color name
 */
export function formatStudentStatusColor(code) {
  const colors = {
    ACTIVE: 'positive',
    GRADUATED: 'info',
    TRANSFERRED: 'warning',
    DROPPED: 'negative',
    ALUMNI: 'grey'
  };
  return colors[code] || 'grey';
}

/**
 * Format alamat lengkap
 * @param {Object} address
 * @returns {string}
 */
export function formatAddress(address) {
  if (!address) return '-';
  const parts = [address.street];
  if (address.rtRw) parts.push(`RT/RW ${address.rtRw}`);
  if (address.village) parts.push(address.village);
  if (address.district) parts.push(address.district);
  if (address.city) parts.push(address.city);
  if (address.postalCode) parts.push(address.postalCode);
  return parts.filter(Boolean).join(', ') || '-';
}

/**
 * Format alamat singkat (kota saja)
 * @param {Object} address
 * @returns {string}
 */
export function formatAddressShort(address) {
  if (!address) return '-';
  return address.city || address.village || '-';
}

/**
 * Format contact info
 * @param {Object} contactInfo - { phone, email }
 * @returns {string}
 */
export function formatContact(contactInfo) {
  if (!contactInfo) return '-';
  const parts = [];
  if (contactInfo.phone) parts.push(contactInfo.phone);
  if (contactInfo.email) parts.push(contactInfo.email);
  return parts.join(' | ') || '-';
}

/**
 * Format guardian relation ke label
 * @param {string} relation - 'father' | 'mother' | 'guardian'
 * @returns {string}
 */
export function formatGuardianRelation(relation) {
  const labels = {
    father: 'Ayah',
    mother: 'Ibu',
    guardian: 'Wali'
  };
  return labels[relation] || relation || '-';
}

/**
 * Format enrollment status
 * @param {string} status - 'active' | 'completed' | 'transferred'
 * @returns {string}
 */
export function formatEnrollmentStatus(status) {
  const labels = {
    active: 'Aktif',
    completed: 'Selesai',
    transferred: 'Pindah'
  };
  return labels[status] || status || '-';
}

/**
 * Format NISN dengan spacing untuk readability (XXXX XXXX XX)
 * @param {string} nisn
 * @returns {string}
 */
export function formatNisnDisplay(nisn) {
  if (!nisn) return '-';
  const str = String(nisn);
  if (str.length !== 10) return str;
  return `${str.slice(0, 4)} ${str.slice(4, 8)} ${str.slice(8)}`;
}

/**
 * Format datetime ISO ke tampilan Indonesia
 * @param {string} isoString
 * @returns {string}
 */
export function formatDateTime(isoString) {
  if (!isoString) return '-';
  try {
    const date = new Date(isoString);
    if (isNaN(date.getTime())) return '-';
    const day = String(date.getDate()).padStart(2, '0');
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const year = date.getFullYear();
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    return `${day}/${month}/${year} ${hours}:${minutes}`;
  } catch {
    return '-';
  }
}
