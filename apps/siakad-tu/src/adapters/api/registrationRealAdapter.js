// apps/siakad-tu/src/adapters/api/registrationRealAdapter.js
import { api } from '@/boot/axios.js'

const mapBackendToRegistration = (backendDto) => {
  if (!backendDto) return null
  return {
    id: backendDto.ID,
    schoolId: backendDto.SchoolID,
    academicPeriodId: backendDto.AcademicPeriodID,
    calonNisn: backendDto.CalonNISN,
    nis: backendDto.NIS,
    calonNama: backendDto.CalonNama,
    calonJenisKelamin: backendDto.CalonJenisKelamin,
    calonTempatLahir: backendDto.CalonTempatLahir,
    calonTanggalLahir: backendDto.CalonTanggalLahir,
    asalSekolah: backendDto.AsalSekolah,
    jurusanDipilih: backendDto.JurusanDipilih,
    biayaStandar: backendDto.BiayaStandar,
    biayaPendaftaran: backendDto.BiayaPendaftaran,
    berkasIjazah: backendDto.BerkasIjazah,
    berkasKK: backendDto.BerkasKK,
    berkasAkta: backendDto.BerkasAkta,
    status: backendDto.Status,
    createdAt: backendDto.CreatedAt,
    updatedAt: backendDto.UpdatedAt,
  }
}

const handleApiError = (error) => ({
  success: false,
  data: null,
  error: {
    code: error.response?.status || 'API_ERROR',
    message: error.response?.data?.error || error.message || 'Terjadi kesalahan.',
  },
})

const handleSuccess = (data) => ({ success: true, data, error: null })

export const registrationRealAdapter = {
  async getRegistrations(context, filters = {}) {
    try {
      const params = {}
      if (filters.status) params.status = filters.status
      if (filters.search) params.search = filters.search
      const response = await api.get('/registrations', { params })
      
      const items = response.data.data || []
      return handleSuccess({
        items: items.map(mapBackendToRegistration),
        total: response.data.total || items.length,
      })
    } catch (error) { return handleApiError(error) }
  },

  async createRegistration(command) {
    try {
      // Payload sesuai service.CreateRegistrationRequest
      const payload = {
        CalonNISN: command.calonNisn || '',
        CalonNama: command.calonNama || '',
        CalonJenisKelamin: command.calonJenisKelamin || 'L',
        CalonTempatLahir: command.calonTempatLahir || '',
        CalonTanggalLahir: command.calonTanggalLahir || null,
        AsalSekolah: command.asalSekolah || '',
        JurusanDipilih: command.jurusanDipilih || '',
        BiayaStandar: command.biayaStandar || 0,
        BiayaPendaftaran: command.biayaPendaftaran || 0,
        BerkasIjazah: command.berkasIjazah || false,
        BerkasKK: command.berkasKK || false,
        BerkasAkta: command.berkasAkta || false,
      }
      const response = await api.post('/registrations', payload)
      return handleSuccess(mapBackendToRegistration(response.data))
    } catch (error) { return handleApiError(error) }
  },

  async generateNisMassal(prefix) {
    try {
      const response = await api.post('/registrations/generate-nis', { prefix })
      return handleSuccess(response.data)
    } catch (error) { return handleApiError(error) }
  },

  async activateRegistration(regId) {
    try {
      const response = await api.post(`/registrations/${regId}/activate`)
      return handleSuccess(response.data)
    } catch (error) { return handleApiError(error) }
  },
}
