import { api as apiClient } from '@/boot/axios'
import { AcademicYear } from '@/models/AcademicYear'
import { Semester } from '@/models/Semester'
import { AcademicPeriod } from '@/models/AcademicPeriod'

export const academicApiAdapter = {
  async getAcademicYears() {
    const { data } = await apiClient.get('/academic-years')
    return data.map((item) => new AcademicYear(item))
  },

  async getSemesters(academicYearId) {
    const { data } = await apiClient.get(`/academic-years/${academicYearId}/semesters`)
    return data.map((item) => new Semester(item))
  },

  async getAcademicPeriods(semesterId) {
    const { data } = await apiClient.get(`/semesters/${semesterId}/periods`)
    return data.map((item) => new AcademicPeriod(item))
  },
}
