// Mock students data (same as attendance)
const mockStudents = {
  'CLS-X-A': [
    { studentId: 'STU-001', studentName: 'Ahmad Fauzi', classId: 'CLS-X-A', className: 'X-A' },
    { studentId: 'STU-002', studentName: 'Budi Hartono', classId: 'CLS-X-A', className: 'X-A' },
    { studentId: 'STU-003', studentName: 'Citra Dewi', classId: 'CLS-X-A', className: 'X-A' },
    { studentId: 'STU-004', studentName: 'Diana Putri', classId: 'CLS-X-A', className: 'X-A' },
    { studentId: 'STU-005', studentName: 'Eko Prasetyo', classId: 'CLS-X-A', className: 'X-A' },
    { studentId: 'STU-006', studentName: 'Fitri Handayani', classId: 'CLS-X-A', className: 'X-A' },
    { studentId: 'STU-007', studentName: 'Gunawan Wijaya', classId: 'CLS-X-A', className: 'X-A' },
    { studentId: 'STU-008', studentName: 'Hana Safitri', classId: 'CLS-X-A', className: 'X-A' },
    { studentId: 'STU-009', studentName: 'Irfan Hakim', classId: 'CLS-X-A', className: 'X-A' },
    { studentId: 'STU-010', studentName: 'Joko Susilo', classId: 'CLS-X-A', className: 'X-A' },
  ],
  'CLS-XI-B': [
    { studentId: 'STU-011', studentName: 'Kartika Sari', classId: 'CLS-XI-B', className: 'XI-B' },
    { studentId: 'STU-012', studentName: 'Lukman Hakim', classId: 'CLS-XI-B', className: 'XI-B' },
    { studentId: 'STU-013', studentName: 'Maya Anggraini', classId: 'CLS-XI-B', className: 'XI-B' },
    { studentId: 'STU-014', studentName: 'Nanda Pratama', classId: 'CLS-XI-B', className: 'XI-B' },
    {
      studentId: 'STU-015',
      studentName: 'Oscar Firmansyah',
      classId: 'CLS-XI-B',
      className: 'XI-B',
    },
  ],
}

// Mock progress records storage
let mockProgressRecords = [
  {
    id: 'PRG-001',
    studentId: 'STU-001',
    studentName: 'Ahmad Fauzi',
    classId: 'CLS-X-A',
    className: 'X-A',
    type: 'achievement',
    title: 'Juara 1 Olimpiade Matematika',
    description: 'Meraih juara 1 dalam Olimpiade Matematika tingkat kota.',
    date: '2026-07-10',
    teacherId: 'USR-001',
    teacherName: 'Budi Santoso, S.Pd.',
    academicYearId: 'AY-2026',
    semesterId: '20261',
    schoolId: 'SCH-001',
  },
  {
    id: 'PRG-002',
    studentId: 'STU-003',
    studentName: 'Citra Dewi',
    classId: 'CLS-X-A',
    className: 'X-A',
    type: 'note',
    title: 'Peningkatan Kinerja',
    description: 'Menunjukkan peningkatan signifikan dalam tugas-tugas matematika.',
    date: '2026-07-12',
    teacherId: 'USR-001',
    teacherName: 'Budi Santoso, S.Pd.',
    academicYearId: 'AY-2026',
    semesterId: '20261',
    schoolId: 'SCH-001',
  },
]

let progressCounter = 3

export class ProgressMockAdapter {
  async fetchStudentsByClass({ classId }) {
    await new Promise((resolve) => setTimeout(resolve, 500))
    return mockStudents[classId] || []
  }

  async loadProgressByStudent({ studentId, academicYearId, semesterId }) {
    await new Promise((resolve) => setTimeout(resolve, 400))

    return mockProgressRecords.filter(
      (record) =>
        record.studentId === studentId &&
        record.academicYearId === academicYearId &&
        record.semesterId === semesterId,
    )
  }

  async createProgress({
    studentId,
    studentName,
    classId,
    className,
    type,
    title,
    description,
    date,
    teacherId,
    teacherName,
    schoolId,
    academicYearId,
    semesterId,
  }) {
    await new Promise((resolve) => setTimeout(resolve, 400))

    const newProgress = {
      id: `PRG-${String(progressCounter++).padStart(3, '0')}`,
      studentId,
      studentName,
      classId,
      className,
      type,
      title,
      description,
      date,
      teacherId,
      teacherName,
      academicYearId,
      semesterId,
      schoolId,
    }

    mockProgressRecords.push(newProgress)
    return newProgress
  }
}
