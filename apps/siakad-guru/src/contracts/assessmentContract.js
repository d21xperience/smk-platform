export const AssessmentContract = {
  loadAssessmentByClass: {
    request: {
      classId: 'string',
      subjectId: 'string',
      academicYearId: 'string',
      semesterId: 'string'
    },
    response: {
      id: 'string',
      classId: 'string',
      className: 'string',
      subjectId: 'string',
      subjectName: 'string',
      date: 'string',
      components: [
        {
          id: 'string',
          name: 'string',
          type: 'string',
          weight: 'number'
        }
      ],
      scores: 'object', // { componentId: score }
      status: 'string', // draft, finalized
      academicYearId: 'string',
      semesterId: 'string',
      schoolId: 'string'
    }
  },

  createAssessment: {
    request: {
      classId: 'string',
      className: 'string',
      subjectId: 'string',
      subjectName: 'string',
      date: 'string',
      components: [
        {
          id: 'string',
          name: 'string',
          type: 'string',
          weight: 'number'
        }
      ],
      schoolId: 'string',
      academicYearId: 'string',
      semesterId: 'string'
    },
    response: {
      id: 'string',
      classId: 'string',
      className: 'string',
      subjectId: 'string',
      subjectName: 'string',
      date: 'string',
      components: 'array',
      scores: 'object',
      status: 'string',
      academicYearId: 'string',
      semesterId: 'string',
      schoolId: 'string'
    }
  },

  updateScore: {
    request: {
      assessmentId: 'string',
      componentId: 'string',
      score: 'number'
    },
    response: {
      success: 'boolean',
      assessment: 'AssessmentSession'
    }
  },

  finalizeAssessment: {
    request: {
      assessmentId: 'string'
    },
    response: {
      success: 'boolean',
      assessment: 'AssessmentSession',
      result: {
        finalScore: 'number',
        grade: 'string',
        predicate: 'string'
      }
    }
  }
}
