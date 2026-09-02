export const ProgressContract = {
  fetchStudentsByClass: {
    request: {
      classId: 'string',
    },
    response: [
      {
        studentId: 'string',
        studentName: 'string',
        classId: 'string',
        className: 'string',
      },
    ],
  },

  loadProgressByStudent: {
    request: {
      studentId: 'string',
      academicYearId: 'string',
      semesterId: 'string',
    },
    response: [
      {
        id: 'string',
        studentId: 'string',
        studentName: 'string',
        classId: 'string',
        className: 'string',
        type: 'string', // note, achievement, violation, counseling
        title: 'string',
        description: 'string',
        date: 'string',
        teacherId: 'string',
        teacherName: 'string',
        academicYearId: 'string',
        semesterId: 'string',
        schoolId: 'string',
      },
    ],
  },

  createProgress: {
    request: {
      studentId: 'string',
      studentName: 'string',
      classId: 'string',
      className: 'string',
      type: 'string',
      title: 'string',
      description: 'string',
      date: 'string',
      teacherId: 'string',
      teacherName: 'string',
      schoolId: 'string',
      academicYearId: 'string',
      semesterId: 'string',
    },
    response: {
      id: 'string',
      studentId: 'string',
      studentName: 'string',
      classId: 'string',
      className: 'string',
      type: 'string',
      title: 'string',
      description: 'string',
      date: 'string',
      teacherId: 'string',
      teacherName: 'string',
      academicYearId: 'string',
      semesterId: 'string',
      schoolId: 'string',
    },
  },
}
