export function useTimer() {
  const formatTime = (sec) => {
    const mm = Math.floor(sec / 60)
      .toString()
      .padStart(2, '0')
    const ss = (sec % 60).toString().padStart(2, '0')
    return `${mm}:${ss}`
  }

  return { formatTime }
}
