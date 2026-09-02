<template>
  <span ref="formulaRef" class="math-formula"></span>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import katex from 'katex'
import 'katex/dist/katex.min.css'

const props = defineProps({
  formula: {
    type: String,
    required: true
  },
  displayMode: {
    type: Boolean,
    default: false
  }
})

const formulaRef = ref(null)

const renderFormula = () => {
  if (formulaRef.value && props.formula) {
    try {
      katex.render(props.formula, formulaRef.value, {
        throwOnError: false,
        displayMode: props.displayMode,
        strict: false
      })
    } catch (error) {
      console.error('KaTeX error:', error)
      formulaRef.value.textContent = props.formula // Fallback
    }
  }
}

onMounted(renderFormula)
watch(() => props.formula, renderFormula)
</script>

<style scoped>
.math-formula {
  font-size: 1.1em;
  overflow-x: auto;
  overflow-y: hidden;
}
</style>
