import { pembayaranSPPService } from './pembayaranSPPService.js'
import { pembayaranLainService } from './pembayaranLainService.js'
import { pengeluaranService } from './pengeluaranService.js'

export const laporanService = {
  async getPemasukan() {
    const spp = await pembayaranSPPService.getAll()
    const lain = await pembayaranLainService.getAll()
    const totalPemasukan = [...spp, ...lain].reduce((sum, item) => sum + item.jumlah, 0)
    return { data: [...spp, ...lain], total: totalPemasukan }
  },
  async getPengeluaran() {
    const items = await pengeluaranService.getAll()
    const total = items.reduce((sum, item) => sum + item.jumlah, 0)
    return { data: items, total }
  },
  async getBukuKas() {
    const pemasukan = await this.getPemasukan()
    const pengeluaran = await this.getPengeluaran()
    const saldo = pemasukan.total - pengeluaran.total
    return { pemasukan: pemasukan.data, pengeluaran: pengeluaran.data, saldo }
  },
}
