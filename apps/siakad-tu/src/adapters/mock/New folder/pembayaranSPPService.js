// services/pembayaranSPPService.js
const dummyPembayaran = [
  {
    id: 1,
    siswaId: 1,
    bulan: 'Januari',
    tahun: 2025,
    jumlah: 200000,
    tanggalBayar: '2025-01-10',
    petugas: 'Admin',
    buktiUrl: null,
  },
  {
    id: 2,
    siswaId: 1,
    bulan: 'Februari',
    tahun: 2025,
    jumlah: 200000,
    tanggalBayar: '2025-02-12',
    petugas: 'Admin',
    buktiUrl: null,
  },
  {
    id: 3,
    siswaId: 2,
    bulan: 'Januari',
    tahun: 2025,
    jumlah: 200000,
    tanggalBayar: '2025-01-15',
    petugas: 'Admin',
    buktiUrl: null,
  },
]

let nextId = 4
const delay = (ms) => new Promise((resolve) => setTimeout(resolve, ms))

export const pembayaranSPPService = {
  async getAll() {
    await delay(300)
    return [...dummyPembayaran]
  },
  async getBySiswaId(siswaId) {
    await delay(200)
    return dummyPembayaran.filter((p) => p.siswaId === siswaId)
  },
  async create(data) {
    await delay(500)
    const newData = { id: nextId++, ...data, tanggalBayar: new Date().toISOString().split('T')[0] }
    dummyPembayaran.push(newData)
    return newData
  },
  async getById(id) {
    await delay(100)
    return dummyPembayaran.find((p) => p.id === id)
  },
}
