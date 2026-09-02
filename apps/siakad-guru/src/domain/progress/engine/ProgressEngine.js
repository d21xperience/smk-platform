import { StudentProgress } from '../models/StudentProgress.js'
import { isValidProgressType } from '../models/ProgressType.js'

export class ProgressEngine {
  static validateProgress(progress) {
    if (!progress || !(progress instanceof StudentProgress)) {
      throw new Error('Invalid student progress object')
    }
    if (!progress.isValid()) {
      throw new Error('Student progress is missing required fields')
    }
    if (!isValidProgressType(progress.type)) {
      throw new Error(`Invalid progress type: ${progress.type}`)
    }
    return true
  }

  static validateProgressInput({ type, title, description }) {
    if (!isValidProgressType(type)) {
      throw new Error(`Invalid progress type: ${type}`)
    }
    if (!title || title.trim().length === 0) {
      throw new Error('Title is required')
    }
    if (!description || description.trim().length === 0) {
      throw new Error('Description is required')
    }
    return true
  }
}
