/**
 * @typedef {Object} StudentFilter
 * @property {string} [search]
 * @property {string} [status]
 * @property {string} [jurusan]
 * @property {string} [tingkat]
 */

export class StudentAffairsAdapterContract {
  // Dashboard
  async getDashboardMetrics() {
    throw new Error('Method not implemented')
  }

  // Student
  async getStudents() {
    throw new Error('Method not implemented')
  }
  async getStudentById() {
    throw new Error('Method not implemented')
  }
  async createStudent() {
    throw new Error('Method not implemented')
  }
  async updateStudent() {
    throw new Error('Method not implemented')
  }

  // Mutation
  async getMutations() {
    throw new Error('Method not implemented')
  }
  async createMutation() {
    throw new Error('Method not implemented')
  }
  async approveMutation() {
    throw new Error('Method not implemented')
  }
  async rejectMutation() {
    throw new Error('Method not implemented')
  }
  async deleteMutation() {
    throw new Error('Method not implemented')
  }
}
