// apps/siakad-tu/scripts/run-pure-tests.js

/**
 * Simple test runner untuk Pure JS tests.
 * Menjalankan semua test yang tidak butuh browser environment.
 *
 * Usage: node scripts/run-pure-tests.js
 */

import { execSync } from 'child_process'
import { fileURLToPath } from 'url'
import { dirname, join } from 'path'

const __filename = fileURLToPath(import.meta.url)
const __dirname = dirname(__filename)
const rootDir = join(__dirname, '..')

const testFiles = [
  '@/domain/student/__tests__/Student.test.js',
  '@/engine/student/__tests__/StudentEngine.test.js',
  '@/composables/student/__tests__/studentFormatters.test.js',
  '@/composables/student/__tests__/studentValidators.test.js',
]

console.log('═══════════════════════════════════════════════════════════')
console.log('  🧪 SDP Pure JS Test Runner (Node.js Polos)')
console.log('═══════════════════════════════════════════════════════════\n')

let totalPassed = 0
let totalFailed = 0
const results = []

for (const file of testFiles) {
  const fullPath = join(rootDir, file)
  console.log(`\n▶ Running: ${file}`)
  console.log('───────────────────────────────────────────────────────────')

  try {
    const output = execSync(`node "${fullPath}"`, {
      encoding: 'utf-8',
      cwd: rootDir,
      stdio: 'pipe',
    })
    console.log(output)

    // Hitung passed/failed dari output
    const passed = (output.match(/✅/g) || []).length
    const failed = (output.match(/❌/g) || []).length
    totalPassed += passed
    totalFailed += failed
    results.push({ file, passed, failed, status: '✅ PASS' })
  } catch (error) {
    console.error(error.stdout || error.message)
    totalFailed += 1
    results.push({ file, passed: 0, failed: 1, status: '❌ FAIL' })
  }
}

console.log('\n═══════════════════════════════════════════════════════════')
console.log('  📊 SUMMARY')
console.log('═══════════════════════════════════════════════════════════')
results.forEach((r) => {
  console.log(`${r.status}  ${r.file}`)
  console.log(`      ✅ ${r.passed} passed | ❌ ${r.failed} failed`)
})
console.log('───────────────────────────────────────────────────────────')
console.log(`TOTAL: ✅ ${totalPassed} passed | ❌ ${totalFailed} failed`)
console.log('═══════════════════════════════════════════════════════════\n')

process.exit(totalFailed > 0 ? 1 : 0)
