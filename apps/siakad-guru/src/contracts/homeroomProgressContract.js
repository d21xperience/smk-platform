// FILE: src/contracts/homeroomProgressContract.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

/**
 * Homeroom Progress Contract
 * Menentukan bentuk data yang dibutuhkan oleh UI Rekap Nilai Wali Kelas.
 * Frontend-First: hanya field yang diperlukan UI yang dipetakan.
 */

export const PROGRESS_STATUS = {
  PASS: 'PASS',
  FAIL: 'FAIL',
}

export const PROGRESS_TREND = {
  UP: 'UP',
  DOWN: 'DOWN',
  STABLE: 'STABLE',
  NEW: 'NEW',
}

export const PROGRESS_PREDICATE = {
  A: 'A',
  B: 'B',
  C: 'C',
  D: 'D',
}

/**
 * Memetakan hasil olahan Domain Engine ke kontrak UI.
 * Engine sudah memproses, contract hanya memastikan shape konsisten.
 * @param {Object} engineResult - Hasil dari HomeroomProgressEngine.processClassProgress().
 * @returns {Object} Data sesuai kontrak UI.
 */
export function mapHomeroomProgress(engineResult) {
  if (!engineResult) {
    return {
      classId: null,
      className: null,
      summary: {
        totalStudents: 0,
        classAverage: 0,
        highestAverage: 0,
        lowestAverage: 0,
        passRate: 0,
      },
      students: [],
    }
  }

  return {
    classId: engineResult.classId || null,
    className: engineResult.className || null,
    summary: {
      totalStudents: engineResult.summary?.totalStudents || 0,
      classAverage: engineResult.summary?.classAverage || 0,
      highestAverage: engineResult.summary?.highestAverage || 0,
      lowestAverage: engineResult.summary?.lowestAverage || 0,
      passRate: engineResult.summary?.passRate || 0,
    },
    students: Array.isArray(engineResult.students)
      ? engineResult.students.map(mapStudentProgress)
      : [],
  }
}

/**
 * Memetakan data satu siswa dari hasil engine ke kontrak UI.
 * @param {Object} student - Data siswa dari engine.
 * @returns {Object} Data siswa sesuai kontrak UI.
 */
function mapStudentProgress(student) {
  return {
    studentId: student.studentId || '',
    studentName: student.studentName || '',
    seatNumber: student.seatNumber || 0,
    subjects: Array.isArray(student.subjects) ? student.subjects.map(mapSubjectProgress) : [],
    overallAverage: student.overallAverage || 0,
    classRank: student.classRank || 0,
    totalSubjects: student.totalSubjects || 0,
    passCount: student.passCount || 0,
    failCount: student.failCount || 0,
  }
}

/**
 * Memetakan data satu mata pelajaran dari hasil engine ke kontrak UI.
 * @param {Object} subject - Data mata pelajaran dari engine.
 * @returns {Object} Data mata pelajaran sesuai kontrak UI.
 */
function mapSubjectProgress(subject) {
  return {
    subjectId: subject.subjectId || '',
    subjectName: subject.subjectName || '',
    formativeScore: subject.formativeScore || 0,
    summativeScore: subject.summativeScore || 0,
    finalScore: subject.finalScore || 0,
    kkm: subject.kkm || 0,
    status: subject.status || PROGRESS_STATUS.FAIL,
    trend: subject.trend || PROGRESS_TREND.NEW,
    predicate: subject.predicate || PROGRESS_PREDICATE.D,
  }
}
