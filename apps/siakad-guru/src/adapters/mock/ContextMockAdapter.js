// FILE: src/adapters/mock/ContextMockAdapter.js
// STATUS: MODIFY
// STATUS IMPLEMENTASI: COMPLETE

import { OperationalContext } from '../../domain/context/models/OperationalContext.js'
import { CONTEXT_MODE } from '../../domain/context/engine/ContextEngine.js'

const SCHOOL = Object.freeze({
  id: 'SCH-001',
  name: 'SMKS Pasundan Jatinangor',
})

const LEGACY_USER = Object.freeze({
  id: 'USR-001',
  name: 'Budi Santoso, S.Pd.',
  role: 'teacher',
  permissions: ['attendance.create', 'assessment.create', 'journal.write', 'inventory.report'],
})

function delay(ms) {
  return new Promise((resolve) => {
    setTimeout(resolve, ms)
  })
}

function ensureUser(user) {
  if (!user || !user.id || !user.role) {
    throw new Error('User tidak valid untuk context request')
  }
}

function buildSemesters(year) {
  return [
    { id: `${year}1`, name: 'Semester 1' },
    { id: `${year}2`, name: 'Semester 2' },
  ]
}

function buildAcademicYear(year) {
  return {
    id: `AY-${year}`,
    name: `${year}/${year + 1}`,
    semesters: buildSemesters(year),
  }
}

function getActiveAcademicPeriod(date = new Date()) {
  const calendarYear = date.getFullYear()
  const month = date.getMonth() + 1

  if (month >= 7) {
    return {
      year: calendarYear,
      academicYearId: `AY-${calendarYear}`,
      semesterId: `${calendarYear}1`,
    }
  }

  const activeYear = calendarYear - 1

  return {
    year: activeYear,
    academicYearId: `AY-${activeYear}`,
    semesterId: `${activeYear}2`,
  }
}

function buildHistoryAcademicYears(activeYear) {
  const years = []

  for (let year = activeYear - 1; year >= activeYear - 3; year -= 1) {
    years.push(buildAcademicYear(year))
  }

  return years
}

function buildHistoryEntries(academicYears) {
  return academicYears.flatMap((year) =>
    year.semesters.map((semester) =>
      Object.freeze({
        academicYearId: year.id,
        semesterId: semester.id,
        label: `${year.name} - ${semester.name}`,
      }),
    ),
  )
}

export class ContextMockAdapter {
  async fetchAvailableContexts() {
    await delay(500)

    const activePeriod = getActiveAcademicPeriod()
    const academicYears = [
      buildAcademicYear(activePeriod.year),
      ...buildHistoryAcademicYears(activePeriod.year),
    ]

    return {
      school: SCHOOL,
      academicYears,
      user: LEGACY_USER,
    }
  }

  async selectContext({ schoolId, academicYearId, semesterId, user }) {
    ensureUser(user)
    await delay(300)

    return new OperationalContext({
      schoolId,
      academicYearId,
      semesterId,
      userId: user.id,
      role: user.role,
      permissions: Array.isArray(user.permissions) ? user.permissions : [],
    })
  }

  async fetchActiveOperationalContext({ user }) {
    ensureUser(user)
    await delay(300)

    const activePeriod = getActiveAcademicPeriod()

    return new OperationalContext({
      schoolId: SCHOOL.id,
      academicYearId: activePeriod.academicYearId,
      semesterId: activePeriod.semesterId,
      userId: user.id,
      role: user.role,
      permissions: Array.isArray(user.permissions) ? user.permissions : [],
    })
  }

  async fetchAvailableHistory({ user }) {
    ensureUser(user)
    await delay(400)

    const activePeriod = getActiveAcademicPeriod()
    const academicYears = buildHistoryAcademicYears(activePeriod.year)
    const entries = buildHistoryEntries(academicYears)

    return {
      school: SCHOOL,
      academicYears,
      entries,
      activeAcademicYearId: activePeriod.academicYearId,
      activeSemesterId: activePeriod.semesterId,
    }
  }

  async fetchHistoricalContext({ user, academicYearId, semesterId }) {
    ensureUser(user)

    if (!academicYearId || !semesterId) {
      throw new Error('Tahun pelajaran dan semester historis wajib dipilih')
    }

    await delay(300)

    const activePeriod = getActiveAcademicPeriod()
    const historyYears = buildHistoryAcademicYears(activePeriod.year)

    const selectedYear = historyYears.find((year) => year.id === academicYearId)

    if (!selectedYear) {
      throw new Error('Tahun pelajaran historis tidak tersedia')
    }

    const selectedSemester = selectedYear.semesters.find((semester) => semester.id === semesterId)

    if (!selectedSemester) {
      throw new Error('Semester historis tidak tersedia')
    }

    return Object.freeze({
      academicYearId,
      semesterId,
      mode: CONTEXT_MODE.HISTORICAL,
    })
  }
}
