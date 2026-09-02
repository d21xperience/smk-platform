// src/adapters/mock/assessment.mock.js
import { AssessmentSession } from '@/models/AssessmentSession'
import { AssessmentComponent } from '@/models/AssessmentComponent'
import { AssessmentFinalResult } from '@/models/AssessmentFinalResult'

// Data siswa per kelas (sama dengan attendance mock)
const studentsByClass = {
  'XII RPL 1': [
    { id: 1, name: 'Ahmad Fauzi' },
    { id: 2, name: 'Bunga Lestari' },
    { id: 3, name: 'Cahya Pratama' },
    { id: 4, name: 'Dewi Sartika' },
    { id: 5, name: 'Eko Saputro' },
  ],
  'XII RPL 2': [
    { id: 11, name: 'Kartika Sari' },
    { id: 12, name: 'Lutfi Hakim' },
    { id: 13, name: 'Mega Permata' },
  ],
}

let sessions = []
let nextSessionId = 1
let nextComponentId = 1

function createDefaultComponents(sessionId) {
  return [
    new AssessmentComponent({
      id: nextComponentId++,
      sessionId,
      type: 'PH',
      name: 'PH 1',
      weight: 25,
      maxScore: 100,
      scores: [],
    }),
    new AssessmentComponent({
      id: nextComponentId++,
      sessionId,
      type: 'PTS',
      name: 'PTS',
      weight: 30,
      maxScore: 100,
      scores: [],
    }),
    new AssessmentComponent({
      id: nextComponentId++,
      sessionId,
      type: 'PAS',
      name: 'PAS',
      weight: 45,
      maxScore: 100,
      scores: [],
    }),
  ]
}

export const assessmentMockAdapter = {
  async fetchOrCreateAssessmentSession(teachingSessionId, options) {
    await new Promise((resolve) => setTimeout(resolve, 300))
    let session = sessions.find((s) => s.teachingSessionId === teachingSessionId)
    if (!session) {
      const components = createDefaultComponents(nextSessionId)
      session = new AssessmentSession({
        id: nextSessionId++,
        teachingSessionId,
        className: options.className,
        subject: options.subject,
        academicYearId: options.academicYearId,
        semesterId: options.semesterId,
        components,
        status: 'draft',
      })
      sessions.push(session)
    }
    return new AssessmentSession({
      ...session,
      components: session.components.map(
        (c) => new AssessmentComponent({ ...c, scores: [...c.scores] }),
      ),
    })
  },

  async saveAssessmentDraft(sessionData) {
    await new Promise((resolve) => setTimeout(resolve, 200))
    const idx = sessions.findIndex((s) => s.id === sessionData.id)
    if (idx < 0) throw new Error('Sesi tidak ditemukan')
    sessions[idx] = new AssessmentSession({
      ...sessionData,
      components: sessionData.components.map((c) => new AssessmentComponent({ ...c })),
    })
    return new AssessmentSession({
      ...sessions[idx],
      components: sessions[idx].components.map(
        (c) => new AssessmentComponent({ ...c, scores: [...c.scores] }),
      ),
    })
  },

  async submitAssessment(sessionId) {
    await new Promise((resolve) => setTimeout(resolve, 250))
    const session = sessions.find((s) => s.id === sessionId)
    if (!session) throw new Error('Sesi tidak ditemukan')
    session.status = 'submitted'
    return new AssessmentSession({
      ...session,
      components: session.components.map(
        (c) => new AssessmentComponent({ ...c, scores: [...c.scores] }),
      ),
    })
  },

  async getAssessmentComponents(sessionId) {
    await new Promise((resolve) => setTimeout(resolve, 150))
    const session = sessions.find((s) => s.id === sessionId)
    return session
      ? session.components.map((c) => new AssessmentComponent({ ...c, scores: [...c.scores] }))
      : []
  },

  async saveComponentScores(componentId, scores) {
    await new Promise((resolve) => setTimeout(resolve, 150))
    for (const session of sessions) {
      const comp = session.components.find((c) => c.id === componentId)
      if (comp) {
        comp.scores = scores.map((s) => ({ ...s }))
        return new AssessmentComponent({ ...comp, scores: [...comp.scores] })
      }
    }
    throw new Error('Komponen tidak ditemukan')
  },

  async getAssessmentFinalResult(sessionId) {
    await new Promise((resolve) => setTimeout(resolve, 200))
    const session = sessions.find((s) => s.id === sessionId)
    if (!session || session.status !== 'submitted') return null
    return new AssessmentFinalResult({
      sessionId: session.id,
      className: session.className,
      subject: session.subject,
      grades: [], // akan diisi oleh engine
      generatedAt: new Date().toISOString(),
    })
  },

  async getStudentsByClass(className) {
    await new Promise((resolve) => setTimeout(resolve, 200))
    return (studentsByClass[className] || []).map((s) => ({ studentId: s.id, studentName: s.name }))
  },
}
