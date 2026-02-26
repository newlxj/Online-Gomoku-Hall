import { ref } from 'vue'

// 音效管理
export function useSound() {
  const isMuted = ref(false)
  const volume = ref(0.5)

  // 创建音频上下文
  const audioContext = typeof window !== 'undefined' ? new (window.AudioContext || (window as any).webkitAudioContext)() : null

  // 生成简单的音效
  function playTone(frequency: number, duration: number, type: OscillatorType = 'sine') {
    if (!audioContext || isMuted.value) return

    try {
      const oscillator = audioContext.createOscillator()
      const gainNode = audioContext.createGain()

      oscillator.connect(gainNode)
      gainNode.connect(audioContext.destination)

      oscillator.frequency.value = frequency
      oscillator.type = type

      gainNode.gain.setValueAtTime(volume.value * 0.3, audioContext.currentTime)
      gainNode.gain.exponentialRampToValueAtTime(0.01, audioContext.currentTime + duration)

      oscillator.start(audioContext.currentTime)
      oscillator.stop(audioContext.currentTime + duration)
    } catch (e) {
      // 忽略音频错误
    }
  }

  // 落子音效
  function playPlace() {
    playTone(800, 0.1, 'sine')
    setTimeout(() => playTone(600, 0.05, 'sine'), 50)
  }

  // 胜利音效
  function playWin() {
    const notes = [523, 659, 784, 1047] // C5, E5, G5, C6
    notes.forEach((freq, i) => {
      setTimeout(() => playTone(freq, 0.3, 'sine'), i * 150)
    })
  }

  // 失败音效
  function playLose() {
    playTone(300, 0.3, 'sawtooth')
    setTimeout(() => playTone(200, 0.5, 'sawtooth'), 200)
  }

  // 悔棋音效
  function playUndo() {
    playTone(400, 0.1, 'triangle')
    setTimeout(() => playTone(300, 0.1, 'triangle'), 80)
  }

  // 按钮点击音效
  function playClick() {
    playTone(1000, 0.05, 'square')
  }

  // 倒计时警告音效 - 步时小于10秒
  function playTimeWarning() {
    // 紧急的滴答声
    playTone(880, 0.1, 'square')
    setTimeout(() => playTone(880, 0.1, 'square'), 150)
    setTimeout(() => playTone(880, 0.1, 'square'), 300)
  }

  // 倒计时紧急音效 - 步时小于5秒
  function playTimeCritical() {
    // 快速连续的滴声
    playTone(1200, 0.08, 'square')
    setTimeout(() => playTone(1200, 0.08, 'square'), 100)
    setTimeout(() => playTone(1200, 0.08, 'square'), 200)
    setTimeout(() => playTone(1200, 0.08, 'square'), 300)
  }

  // 切换静音
  function toggleMute() {
    isMuted.value = !isMuted.value
  }

  // 设置音量
  function setVolume(v: number) {
    volume.value = Math.max(0, Math.min(1, v))
  }

  return {
    isMuted,
    volume,
    playPlace,
    playWin,
    playLose,
    playUndo,
    playClick,
    playTimeWarning,
    playTimeCritical,
    toggleMute,
    setVolume
  }
}
