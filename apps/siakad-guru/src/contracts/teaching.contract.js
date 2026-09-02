/**
 * FRONTEND DRIVEN CONTRACT
 * Teaching Journal API Contract
 */
export const TeachingContract = {
  // ... existing ...

  loadOrCreateJournal: {
    request: {
      sessionId: 'string',
    },
    response: {
      success: 'boolean',
      data: 'TeachingJournal',
      errors: 'array',
    },
  },
  saveJournalDraft: {
    request: {
      material: { title: 'string', content: 'string', attachments: 'array' },
      learningActivities:
        'array<{ step: string, description: string, duration: number, order: number }>',
      reflection: {
        studentUnderstanding: 'string',
        teacherSelfReflection: 'string',
        improvementPlan: 'string',
        notes: 'string',
      },
    },
    response: {
      success: 'boolean',
      data: 'TeachingJournal',
      errors: 'array',
    },
  },
  finalizeJournal: {
    request: {},
    response: {
      success: 'boolean',
      data: 'TeachingJournal',
      errors: 'array',
    },
  },
}
