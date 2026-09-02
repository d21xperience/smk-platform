<template>
  <span :class="['custom-text', fontClass]" :dir="direction" :style="customStyle">
    <slot>{{ text }}</slot>
  </span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  text: { type: String, default: '' },
  script: {
    type: String,
    default: 'latin',
    validator: (v) => ['latin', 'arabic', 'sundanese', 'javanese', 'bali', 'batak'].includes(v),
  },
})

const fontClass = computed(() => {
  const map = {
    latin: 'font-latin',
    arabic: 'font-arabic',
    sundanese: 'font-sundanese',
    javanese: 'font-javanese',
    bali: 'font-bali',
    batak: 'font-batak',
  }
  return map[props.script] || 'font-latin'
})

const direction = computed(() => (props.script === 'arabic' ? 'rtl' : 'ltr'))

const customStyle = computed(() => ({
  fontFamily: getFontFamily(props.script),
  fontSize: props.script != 'latin' ? '2em' : '1em',
  lineHeight: props.script === 'arabic' ? '2' : '1.6',
}))

const getFontFamily = (script) => {
  const fonts = {
    latin: "'Roboto', sans-serif",
    arabic: "'Amiri', 'Noto Sans Arabic', serif",
    sundanese: "'Noto Sans Sundanese', sans-serif",
    javanese: "'Noto Sans Javanese', sans-serif",
    bali: "'Noto Sans Balinese', sans-serif",
    batak: "'Noto Sans Batak', sans-serif",
  }
  return fonts[script] || fonts.latin
}
</script>

<style scoped>
.custom-text {
  display: inline-block;
}

.font-arabic {
  text-align: right;
}
</style>
