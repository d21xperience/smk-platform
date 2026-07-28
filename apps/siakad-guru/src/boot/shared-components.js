import { boot } from 'quasar/wrappers'
// Impor semua komponen langsung dari package workspace pnpm Anda
import * as SharedUi from '@smk-platform/ui-components'

export default boot(({ app }) => {
  // Loop semua objek yang di-ekspor oleh ui-components
  Object.entries(SharedUi).forEach(([name, component]) => {
    // Daftarkan hanya objek yang merupakan komponen Vue (berawalan 'App')
    if (name.startsWith('App')) {
      app.component(name, component)
    }
  })
})
