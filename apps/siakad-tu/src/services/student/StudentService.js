// apps/siakad-tu/src/services/student/StudentService.js

import { StudentCommandService } from './StudentCommandService.js'
import { StudentQueryService } from './StudentQueryService.js'

/**
 * StudentService — Facade yang menggabungkan Command dan Query.
 *
 * Ini adalah pintu masuk utama untuk semua operasi kesiswaan.
 * UI/Store memanggil Service ini, tidak perlu tahu tentang pemisahan CQRS.
 *
 * Penggunaan:
 *   const service = new StudentService({ engine, adapter, eventDispatcher });
 *   await service.registerStudent(data, context);
 *   await service.getStudents(context, filters);
 */
export class StudentService {
  /**
   * @param {Object} dependencies
   * @param {import('@/engine/student/StudentEngine').StudentEngine} dependencies.engine
   * @param {Object} dependencies.adapter
   * @param {import('@/events/eventDispatcher').EventDispatcher} dependencies.eventDispatcher
   */
  constructor(dependencies) {
    this.command = new StudentCommandService(dependencies)
    this.query = new StudentQueryService(dependencies)
  }

  // === COMMAND METHODS (delegasi ke StudentCommandService) ===

  async registerStudent(commandData, context) {
    return this.command.registerStudent(commandData, context)
  }

  async enrollStudent(commandData, context) {
    return this.command.enrollStudent(commandData, context)
  }

  async updateStudentProfile(commandData, context) {
    return this.command.updateStudentProfile(commandData, context)
  }

  async graduateStudent(commandData, context) {
    return this.command.graduateStudent(commandData, context)
  }

  async transferStudent(commandData, context) {
    return this.command.transferStudent(commandData, context)
  }

  // === QUERY METHODS (delegasi ke StudentQueryService) ===

  async getStudentById(studentId, context) {
    return this.query.getStudentById(studentId, context)
  }

  async getStudents(context, filters = {}) {
    return this.query.getStudents(context, filters)
  }

  async checkNisnAvailability(nisn, context) {
    return this.query.checkNisnAvailability(nisn, context)
  }

  async checkNisAvailability(nis, context) {
    return this.query.checkNisAvailability(nis, context)
  }
}
