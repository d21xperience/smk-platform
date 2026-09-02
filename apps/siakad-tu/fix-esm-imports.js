// apps/siakad-tu/fix-esm-imports.js

import { readdirSync, readFileSync, writeFileSync, statSync } from 'fs'
import { join } from 'path'
import { fileURLToPath } from 'url'

const __dirname = fileURLToPath(new URL('.', import.meta.url))
const srcDir = join(__dirname, 'src')

function fixImports(dir) {
  const files = readdirSync(dir)
  for (const file of files) {
    const fullPath = join(dir, file)
    if (statSync(fullPath).isDirectory()) {
      // Skip node_modules atau folder tersembunyi
      if (file === 'node_modules' || file.startsWith('.')) continue
      fixImports(fullPath)
    } else if (file.endsWith('.js')) {
      let content = readFileSync(fullPath, 'utf-8')

      // Regex: Cari "from './something'" atau "from '../something'"
      // yang TIDAK diakhiri dengan '.js'
      const updatedContent = content.replace(
        /(from\s+['"])(\.[^'"]+)(['"])/g,
        (match, prefix, path, suffix) => {
          if (!path.endsWith('.js')) {
            return `${prefix}${path}.js${suffix}`
          }
          return match
        },
      )

      if (content !== updatedContent) {
        writeFileSync(fullPath, updatedContent, 'utf-8')
        console.log(`✅ Fixed: ${fullPath.replace(__dirname + '/', '')}`)
      }
    }
  }
}

console.log('🔧 Memulai perbaikan ekstensi import ESM...\n')
fixImports(srcDir)
console.log('\n🎉 Selesai! Semua relative import sekarang memiliki ekstensi .js')
