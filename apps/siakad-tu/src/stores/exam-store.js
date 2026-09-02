import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useExamStore = defineStore('exam', () => {
  const user = ref(null)
  const isLoggedIn = ref(false)
  const role = ref('')
  const timer = ref(3600) // 60 menit
  let timerInterval = null

  const mockStudents = ref([
    { id: 'S001', name: 'Andi Pratama', status: 'Belum Login', progress: 0, lastHeartbeat: null },
    {
      id: 'S002',
      name: 'Budi Santoso',
      status: 'Mengerjakan',
      progress: 45,
      lastHeartbeat: Date.now(),
    },
    { id: 'S003', name: 'Citra Dewi', status: 'Selesai', progress: 100, lastHeartbeat: Date.now() },
    {
      id: 'S004',
      name: 'Dina Sari',
      status: 'Terputus',
      progress: 60,
      lastHeartbeat: Date.now() - 40000,
    },
  ])

  const questions = ref([
    {
      id: 1,
      type: 'pg',
      text: 'Ibukota Indonesia saat ini adalah?',
      script: 'latin',
      options: ['Jakarta', 'Bandung', 'Surabaya', 'Medan'],
      answer: 'Jakarta',
    },
    { id: 2, type: 'essay', text: 'Jelaskan 3 manfaat fotosintesis bagi manusia!', answer: '' },
    {
      id: 3,
      type: 'pg',
      text: '2 + 2 × 3 = ?',
      script: 'latin',
      options: ['8', '10', '12', '6'],
      answer: '8',
    },
    // Matematika
    {
      id: 4,
      type: 'pg',
      text: 'Hasil dari $\\frac{3}{4} + \\frac{2}{5}$ adalah...',
      options: ['$\\frac{23}{20}$', '$\\frac{5}{9}$', '$\\frac{6}{20}$', '$\\frac{1}{2}$'],
      answer: '$\\frac{23}{20}$',
    },

    // Bahasa Arab
    {
      id: 5,
      type: 'pg',
      text: 'لمّا كان الاعتراف بالكرامة المتأصلة في جميع',
      script: 'arabic',
      options: [
        'قُلْ هُوَ اللّٰهُ اَحَدٌ',
        'وَلَمْ يَكُنْ لَهٗ كُفُوًا اَحَدٌ',
        'كِتَابٌ',
        'مَدْرَسَةٌ',
      ],
      answer: 'كِتَابٌ',
    },

    // Aksara Sunda
    {
      id: 6,
      type: 'pg',
      text: 'ᮞᮊᮙᮔ...',
      script: 'sundanese',
      options: ['ᮊᮥᮛ᮪ᮓᮤᮔ᮪', 'ᮃᮊ᮪ᮞᮛ', 'ᮞᮥᮔ᮪ᮓ', 'ᮞᮥᮔ᮪ᮓᮞ'],
      answer: 'ᮊᮛ᮪ᮓᮤᮔ',
    },

    // Aksara Jawa
    {
      id: 7,
      type: 'essay',
      text: 'ꦲꦏ꧀ꦱꦫꦗꦮ꧈ꦠꦸꦭꦶꦱꦏꦺꦴꦤ꧀ꦝꦁꦲꦶꦁꦲꦏ꧀ꦫꦗꦮ',
      script: 'javanese',
      answer: '',
    },

    // Campuran (Arab + Latin)
    {
      id: 8,
      type: 'pg',
      text: 'Surat Al-Fatihah terdiri dari ... ayat',
      script: 'latin',
      options: ['5', '6', '7', '8'],
      answer: '7',
      additionalInfo: 'بِسْمِ اللَّهِ الرَّحْمَٰنِ الرَّحِيمِ',
    },
    {
      id: 9,
      type: 'pg',

      text: 'Diketahui persamaan kuadrat $x^2 - 5x + 6 = 0$. Nilai $x$ adalah?',
      options: [
        '$x = 2$ atau $x = 3$',
        '$x = -2$ atau $x = -3$',
        '$x = 1$ atau $x = 6$',
        '$x = -1$ atau $x = -6$',
      ],
      answer: '$x = 2$ atau $x = 3$',
    },
    {
      id: 10,
      type: 'pg',
      text: 'Hasil dari $\\int_0^1 x^2 \\, dx$ adalah...',
      options: ['$\\frac{1}{3}$', '$\\frac{1}{2}$', '$1$', '$\\frac{2}{3}$'],
      answer: '$\\frac{1}{3}$',
    },
    {
      id: 11,
      type: 'essay',
      text: 'Buktikan teorema Pythagoras: $a^2 + b^2 = c^2$',
      answer: '',
    },
    // src/data/questions.js

    {
      id: 12,
      type: 'hotspot',
      text: 'Perhatikan peta berikut! Klik pada area yang menunjukkan letak Pulau Kalimantan!',
      script: 'latin',
      image:
        'https://img.magnific.com/premium-vector/peta-indonesia-map-indonesia-indonesia-map-illustration_61558-320.jpg?semt=ais_hybrid&w=740&q=80',
      zones: [
        {
          id: 'z1',
          label: 'Pulau Sumatera',
          correct: false,
          // Koordinat polygon dalam viewBox 0 0 100 100
          points: '5,20 25,18 28,45 8,48',
        },
        {
          id: 'z2',
          label: 'Pulau Kalimantan',
          correct: true, // ← jawaban benar
          points: '28,15 55,13 58,50 30,52',
        },
        {
          id: 'z3',
          label: 'Pulau Jawa',
          correct: false,
          points: '28,52 55,50 57,62 30,64',
        },
        {
          id: 'z4',
          label: 'Pulau Sulawesi',
          correct: false,
          points: '58,15 75,13 78,50 60,52',
        },
        {
          id: 'z5',
          label: 'Pulau Papua',
          correct: false,
          points: '78,15 98,13 98,55 78,55',
        },
      ],
      answer: 'z2', // id zona yang benar
    },
    {
      id: 13,
      type: 'matching',
      text: 'Pasangkan negara dengan ibu kotanya yang benar!',
      leftItems: [
        { id: 'l1', text: 'Indonesia' },
        { id: 'l2', text: 'Jepang' },
        { id: 'l3', text: 'Prancis' },
        { id: 'l4', text: 'Brazil' },
      ],
      rightItems: [
        { id: 'r1', text: 'Paris' },
        { id: 'r2', text: 'Brasília' },
        { id: 'r3', text: 'Tokyo' },
        { id: 'r4', text: 'Jakarta' },
      ],
      // Kunci jawaban: id kiri → id kanan
      answer: {
        l1: 'r4',
        l2: 'r3',
        l3: 'r1',
        l4: 'r2',
      },
    },
  ])

  const answers = ref({})
  const currentQIndex = ref(0)

  const login = (serial, token) => {
    if (serial === 'S001' && token === '1234') {
      user.value = { id: 'S001', name: 'Andi Pratama' }
      role.value = 'siswa'
      isLoggedIn.value = true
      startTimer()
      return { success: true, redirect: '/exam' }
    }
    if (serial === 'GURU01' && token === 'admin') {
      user.value = { id: 'G001', name: 'Pak Guru Matematika' }
      role.value = 'guru'
      isLoggedIn.value = true
      startMonitoringSimulation()
      return { success: true, redirect: '/monitor' }
    }
    return { success: false, message: 'Serial atau Token salah' }
  }

  const startTimer = () => {
    if (timerInterval) clearInterval(timerInterval)
    timerInterval = setInterval(() => {
      if (timer.value > 0) timer.value--
      else clearInterval(timerInterval)
    }, 1000)
  }

  const startMonitoringSimulation = () => {
    setInterval(() => {
      mockStudents.value.forEach((s) => {
        if (Math.random() > 0.6) {
          const statuses = ['Belum Login', 'Login', 'Mengerjakan', 'Selesai', 'Terputus']
          s.status = statuses[Math.floor(Math.random() * statuses.length)]
          s.progress =
            s.status === 'Selesai'
              ? 100
              : Math.min(100, s.progress + Math.floor(Math.random() * 10))
          s.lastHeartbeat = Date.now()
        }
      })
    }, 5000)
  }

  const saveAnswer = (qId, value) => {
    answers.value[qId] = value
  }

  const logout = () => {
    user.value = null
    isLoggedIn.value = false
    role.value = ''
    timer.value = 3600
    answers.value = {}
    currentQIndex.value = 0
    if (timerInterval) clearInterval(timerInterval)
  }

  return {
    user,
    isLoggedIn,
    role,
    mockStudents,
    questions,
    answers,
    currentQIndex,
    timer,
    login,
    saveAnswer,
    logout,
  }
})
