/**
 * Shared constants untuk seluruh aplikasi
 */

// Attendance Status
export const ATTENDANCE_STATUS = {
  PRESENT: "PRESENT",
  ABSENT: "ABSENT",
  SICK: "SICK",
  PERMISSION: "PERMISSION",
};

export const ATTENDANCE_STATUS_LABELS = {
  PRESENT: "Hadir",
  ABSENT: "Alpha",
  SICK: "Sakit",
  PERMISSION: "Izin",
};

export const ATTENDANCE_STATUS_COLORS = {
  PRESENT: "positive",
  ABSENT: "negative",
  SICK: "warning",
  PERMISSION: "info",
};

// Session Status
export const SESSION_STATUS = {
  DRAFT: "DRAFT",
  SUBMITTED: "SUBMITTED",
  CLOSED: "CLOSED",
};

export const SESSION_STATUS_LABELS = {
  DRAFT: "Draft",
  SUBMITTED: "Submitted",
  CLOSED: "Closed",
};

// Student Status
export const STUDENT_STATUS = {
  ACTIVE: "ACTIVE",
  INACTIVE: "INACTIVE",
  GRADUATED: "GRADUATED",
  DROPPED: "DROPPED",
};

export const STUDENT_STATUS_LABELS = {
  ACTIVE: "Aktif",
  INACTIVE: "Tidak Aktif",
  GRADUATED: "Lulus",
  DROPPED: "Keluar",
};

// Gender
export const GENDER = {
  MALE: "M",
  FEMALE: "F",
};

export const GENDER_LABELS = {
  M: "Laki-laki",
  F: "Perempuan",
};

// Assessment Types
export const ASSESSMENT_TYPES = {
  EXAM: "EXAM",
  QUIZ: "QUIZ",
  ASSIGNMENT: "ASSIGNMENT",
  PRACTICE: "PRACTICE",
};

export const ASSESSMENT_TYPE_LABELS = {
  EXAM: "Ujian",
  QUIZ: "Kuis",
  ASSIGNMENT: "Tugas",
  PRACTICE: "Praktik",
};

// Pagination
export const DEFAULT_PAGE_SIZE = 10;
export const PAGE_SIZE_OPTIONS = [10, 25, 50, 100];

// Date Formats
export const DATE_FORMAT = {
  DISPLAY: "DD MMMM YYYY",
  INPUT: "YYYY-MM-DD",
  DATETIME: "DD/MM/YYYY HH:mm",
};
