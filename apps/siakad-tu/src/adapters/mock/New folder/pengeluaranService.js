let dummyPengeluaran = [
  {
    id: 1,
    tanggal: '2025-02-01',
    nominal: 500000,
    kategoriId: 1,
    keterangan: 'Pembelian kertas dan tinta',
    bukti: 'bukti1.png',
  },
  {
    id: 2,
    tanggal: '2025-02-05',
    nominal: 2000000,
    kategoriId: 2,
    keterangan: 'Perbaikan atap kelas',
    bukti: null,
  },
  {
    id: 3,
    tanggal: '2025-02-10',
    nominal: 750000,
    kategoriId: 3,
    keterangan: 'Tagihan internet bulan Januari',
    bukti: 'bukti3.pdf',
  },
]
let nextId = 4
const delay = (ms) => new Promise((resolve) => setTimeout(resolve, ms))

export const pengeluaranService = {
  async getAll() {
    await delay(300)
    return [...dummyPengeluaran]
  },
  async create(data) {
    await delay(500)
    const baru = { id: nextId++, ...data, bukti: data.bukti || null }
    dummyPengeluaran.push(baru)
    return baru
  },
  async update(id, data) {
    await delay(400)
    const index = dummyPengeluaran.findIndex((p) => p.id === id)
    if (index !== -1) dummyPengeluaran[index] = { ...dummyPengeluaran[index], ...data }
    return dummyPengeluaran[index]
  },
  async delete(id) {
    await delay(400)
    const index = dummyPengeluaran.findIndex((p) => p.id === id)
    if (index !== -1) dummyPengeluaran.splice(index, 1)
    return true
  },
}
