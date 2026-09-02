import katex from 'katex'
import 'katex/dist/katex.min.css'
// import { computed } from 'vue'

export function useLatex() {
  const renderLatex = (text) => {
    if (!text) return ''
    return text.replace(/\$(.+?)\$/g, (match, formula) => {
      try {
        return katex.renderToString(formula.trim(), { throwOnError: false })
      } catch {
        return match
      }
    })
  }

  return { renderLatex }
}
