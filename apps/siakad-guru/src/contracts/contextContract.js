// FILE: src/contracts/contextContract.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

import { CONTEXT_MODE } from '@/domain/context/engine/ContextEngine.js'

export const ACTIVE_CONTEXT_STATUS = Object.freeze({
  ACTIVE: 'ACTIVE',
  INACTIVE: 'INACTIVE',
})

export const contextContract = Object.freeze({
  activeContext: Object.freeze({
    schoolId: 'SCH-001',
    academicYearId: 'AY-2026',
    semesterId: '20261',
    userId: 'USR-001',
    role: 'teacher',
    permissions: [],
  }),

  historicalContext: Object.freeze({
    academicYearId: 'AY-2025',
    semesterId: '20252',
    mode: CONTEXT_MODE.HISTORICAL,
  }),

  availableHistory: Object.freeze({
    school: Object.freeze({
      id: 'SCH-001',
      name: 'SMKS Pasundan Jatinangor',
    }),
    academicYears: [],
    entries: [],
    activeAcademicYearId: 'AY-2026',
    activeSemesterId: '20261',
  }),
})

export function normalizeAvailableHistory(payload) {
  const school =
    payload && payload.school
      ? Object.freeze({
          id: payload.school.id ?? null,
          name: payload.school.name ?? '',
        })
      : null

  const academicYears = Array.isArray(payload?.academicYears)
    ? payload.academicYears
        .filter((year) => year && year.id && year.name)
        .map((year) =>
          Object.freeze({
            id: year.id,
            name: year.name,
            semesters: Array.isArray(year.semesters)
              ? year.semesters
                  .filter((semester) => semester && semester.id && semester.name)
                  .map((semester) =>
                    Object.freeze({
                      id: semester.id,
                      name: semester.name,
                    }),
                  )
              : [],
          }),
        )
    : []

  const rawEntries = Array.isArray(payload?.entries)
    ? payload.entries
        .filter((entry) => entry && entry.academicYearId && entry.semesterId)
        .map((entry) =>
          Object.freeze({
            academicYearId: entry.academicYearId,
            semesterId: entry.semesterId,
            label: entry.label ?? `${entry.academicYearId} - ${entry.semesterId}`,
          }),
        )
    : []

  const entries =
    rawEntries.length > 0
      ? rawEntries
      : academicYears.flatMap((year) =>
          year.semesters.map((semester) =>
            Object.freeze({
              academicYearId: year.id,
              semesterId: semester.id,
              label: `${year.name} - ${semester.name}`,
            }),
          ),
        )

  return Object.freeze({
    school,
    academicYears,
    entries,
    activeAcademicYearId: payload?.activeAcademicYearId ?? null,
    activeSemesterId: payload?.activeSemesterId ?? null,
  })
}
