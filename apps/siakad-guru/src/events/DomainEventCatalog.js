/**
 * Domain Event Catalog
 * Spesifikasi lengkap setiap event untuk kontrak antar modul.
 */
export const EventCatalog = {
  TEACHING_SESSION_STARTED: {
    producer: 'TeachingEngine',
    payload: {
      sessionId: 'number',
      date: 'string (YYYY-MM-DD)',
      startTime: 'string (HH:mm)',
      subject: 'string',
      className: 'string',
      teacherId: 'number',
    },
  },

  TEACHING_SESSION_COMPLETED: {
    producer: 'TeachingEngine',
    payload: {
      sessionId: 'number',
      date: 'string',
      endTime: 'string',
      subject: 'string',
      className: 'string',
      teacherId: 'number',
      durationMinutes: 'number', // durasi aktual
      isJournalFilled: 'boolean',
      isAttendanceSubmitted: 'boolean',
    },
  },

  TEACHING_SESSION_LOCKED: {
    producer: 'TeachingEngine',
    payload: {
      sessionId: 'number',
      lockedAt: 'string (ISO)',
      lockedBy: 'string (user role)',
    },
  },

  TEACHER_CHECKED_IN: {
    producer: 'TeachingEngine',
    payload: {
      sessionId: 'number',
      teacherId: 'number',
      checkInTime: 'string (ISO)',
      status: 'string (on_time | late)',
    },
  },

  TEACHER_CHECKED_OUT: {
    producer: 'TeachingEngine',
    payload: {
      sessionId: 'number',
      teacherId: 'number',
      checkOutTime: 'string (ISO)',
    },
  },

  TEACHER_LATE: {
    producer: 'TeachingEngine',
    payload: {
      sessionId: 'number',
      teacherId: 'number',
      scheduledStart: 'string (ISO)',
      actualCheckIn: 'string (ISO)',
      lateMinutes: 'number',
    },
  },

  ATTENDANCE_DRAFT_SAVED: {
    producer: 'AttendanceEngine',
    payload: {
      sessionId: 'number',
      className: 'string',
      savedAt: 'string (ISO)',
      studentCount: 'number',
    },
  },

  ATTENDANCE_SUBMITTED: {
    producer: 'AttendanceEngine',
    payload: {
      sessionId: 'number',
      className: 'string',
      submittedAt: 'string (ISO)',
      summary: {
        total: 'number',
        hadir: 'number',
        izin: 'number',
        sakit: 'number',
        alpha: 'number',
      },
    },
  },

  JOURNAL_SAVED: {
    producer: 'TeachingEngine',
    payload: {
      sessionId: 'number',
      subject: 'string',
      className: 'string',
      savedAt: 'string (ISO)',
    },
  },

  AUDIT_ENTRY_CREATED: {
    producer: 'AuditService',
    payload: {
      sessionId: 'number',
      action: 'string',
      timestamp: 'string (ISO)',
      details: 'string',
    },
  },
}
