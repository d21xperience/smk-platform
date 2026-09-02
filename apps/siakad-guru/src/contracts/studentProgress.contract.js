export const StudentProgressContract = {
  getClassProgress: {
    request: { classId: 'string', semesterId: 'string', academicYearId: 'string' },
    response: { success: 'boolean', data: 'array<StudentProgress>', errors: 'array' },
  },
  getStudentProgress: {
    request: { studentId: 'string', semesterId: 'string', academicYearId: 'string' },
    response: { success: 'boolean', data: 'StudentProgress', errors: 'array' },
  },
  addTeacherNote: {
    request: { studentId: 'string', content: 'string', type: 'string', teacherId: 'string' },
    response: { success: 'boolean', data: 'TeacherNote', errors: 'array' },
  },
  syncFromAssessment: {
    request: {
      studentId: 'string',
      averageScore: 'number',
      coreSubjectsGrades: 'object',
      promotionStatus: 'string',
    },
    response: { success: 'boolean', data: 'StudentProgress', errors: 'array' },
  },
}
