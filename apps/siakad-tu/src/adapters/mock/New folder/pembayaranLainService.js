let dummyPembayaranLain = [
  {
    id: 1,
    siswaId: 1,
    kategoriId: 1,
    jumlah: 150000,
    tanggalBayar: '2025-03-10',
    keterangan: 'Praktikum IPA',
    petugas: 'Admin',
  },
  {
    id: 2,
    siswaId: 1,
    kategoriId: 2,
    jumlah: 500000,
    tanggalBayar: '2025-03-15',
    keterangan: 'Gedung Tahap 1',
    petugas: 'Admin',
  },
  {
    id: 3,
    siswaId: 2,
    kategoriId: 3,
    jumlah: 100000,
    tanggalBayar: '2025-03-20',
    keterangan: 'Ujian Tengah Semester',
    petugas: 'Admin',
  },
]
let nextId = 4
const delay = (ms) => new Promise((resolve) => setTimeout(resolve, ms))

export const pembayaranLainService = {
  async getAll() {
    await delay(300)
    return [...dummyPembayaranLain]
  },
  async getBySiswaId(siswaId) {
    await delay(200)
    return dummyPembayaranLain.filter((p) => p.siswaId === siswaId)
  },
  async create(data) {
    await delay(500)
    const newData = { id: nextId++, ...data, tanggalBayar: new Date().toISOString().split('T')[0] }
    dummyPembayaranLain.push(newData)
    return newData
  },
  async delete(id) {
    await delay(400)
    const index = dummyPembayaranLain.findIndex((p) => p.id === id)
    if (index !== -1) {
      dummyPembayaranLain.splice(index, 1)
      return true
    }
    return false
  },
}
