export const JournalContract = {
  loadJournalByTeachingSession: {
    request: {
      teachingSessionId: 'string'
    },
    response: {
      id: 'string',
      teachingSessionId: 'string',
      classId: 'string',
      className: 'string',
      subjectId: 'string',
      subjectName: 'string',
      date: 'string',
      material: 'string',
      activities: 'string',
      reflection: 'string',
      status: 'string', // draft, submitted
      academicYearId: 'string',
      semesterId: 'string',
      schoolId: 'string'
    }
  },

  createJournal: {
    request: {
      teachingSessionId: 'string',
      classId: 'string',
      className: 'string',
      subjectId: 'string',
      subjectName: 'string',
      date: 'string',
      schoolId: 'string',
      academicYearId: 'string',
      semesterId: 'string'
    },
    response: {
      id: 'string',
      teachingSessionId: 'string',
      classId: 'string',
      className: 'string',
      subjectId: 'string',
      subjectName: 'string',
      date: 'string',
      material: 'string',
      activities: 'string',
      reflection: 'string',
      status: 'string',
      academicYearId: 'string',
      semesterId: 'string',
      schoolId: 'string'
    }
  },

  updateJournal: {
    request: {
      journalId: 'string',
      material: 'string',
      activities: 'string',
      reflection: 'string'
    },
    response: {
      success: 'boolean',
      journal: 'TeachingJournal'
    }
  },

  submitJournal: {
    request: {
      journalId: 'string'
    },
    response: {
      success: 'boolean',
      journal: 'TeachingJournal'
    }
  }
}
