/* eslint-disable no-unused-vars */
// PURE JS
import { Projection } from './Projection'

export class ProjectionManager {
  constructor() {
    this.projections = new Map()
  }

  register({ name, initialState, updateFn }) {
    if (this.projections.has(name)) {
      throw new Error(`Projection "${name}" already registered.`)
    }
    const projection = new Projection({ name, initialState, updateFn })
    this.projections.set(name, projection)
    return projection
  }

  getProjection(name) {
    const projection = this.projections.get(name)
    if (!projection) {
      throw new Error(`Projection "${name}" not found.`)
    }
    return projection
  }

  getState(name) {
    return this.getProjection(name).getState()
  }

  applyEvent(event) {
    // Apply event to all projections?
    // Biasanya hanya projection yang terdaftar untuk event tertentu.
    // Di sini kita akan scan semua projection dan jika ada yang ingin menerima event,
    // kita bisa lakukan, tetapi lebih baik menggunakan listener terpisah.
    // Kita akan expose method untuk individual update.
  }

  updateProjection(name, event) {
    const projection = this.getProjection(name)
    return projection.applyEvent(event)
  }

  resetAll() {
    for (const [name, projection] of this.projections) {
      projection.reset(projection.initialState)
    }
  }
}

// Singleton
export const projectionManager = new ProjectionManager()
