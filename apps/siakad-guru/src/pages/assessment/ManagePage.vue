<!-- eslint-disable no-unused-vars -->
<template>
  <div class="q-pa-md">
    <h6>Manajemen Penilaian</h6>

    <!-- Form Pilih Konteks (Jika belum ada session) -->
    <div v-if="!currentSession" class="row q-col-gutter-md">
      <div class="col-md-3 col-sm-6 col-xs-12">
        <q-select v-model="selectedClass" :options="classOptions" label="Kelas" outlined dense />
      </div>
      <div class="col-md-3 col-sm-6 col-xs-12">
        <q-select v-model="selectedSubject" :options="subjectOptions" label="Mata Pelajaran" outlined dense />
      </div>
      <div class="col-md-3 col-sm-6 col-xs-12">
        <q-select v-model="selectedPeriod" :options="periodOptions" label="Periode" outlined dense />
      </div>
      <div class="col-md-3 col-sm-6 col-xs-12">
        <q-btn color="primary" label="Buat / Buka Sesi" @click="loadOrCreate" :loading="loading" />
      </div>
    </div>

    <!-- Detail Sesi -->
    <div v-else>
      <!-- Header Sesi -->
      <div class="row items-center q-mb-md">
        <div class="col">
          <div class="text-subtitle1">
            {{ currentSession.classId }} - {{ currentSession.subjectId }}
          </div>
          <div class="text-caption text-grey-6">
            {{ currentSession.period }} | {{ currentSession.date }}
            <q-badge :color="isFinalized ? 'positive' : 'warning'" class="q-ml-sm">
              {{ isFinalized ? 'Final' : 'Draft' }}
            </q-badge>
          </div>
        </div>
        <div class="col-auto">
          <q-btn v-if="!isFinalized" color="positive" label="Finalisasi" @click="finalizeAssessment"
            :loading="loading" />
          <q-btn color="grey" label="Tutup" @click="closeSession" flat />
        </div>
      </div>

      <!-- Manajemen Komponen -->
      <q-card flat bordered class="q-mb-md">
        <q-card-section>
          <div class="row items-center">
            <div class="col">
              <div class="text-subtitle2">Komponen Penilaian</div>
              <div class="text-caption text-grey-6">
                Total Bobot: {{ totalWeight }}% {{ totalWeight === 100 ? '✅' : '⚠️' }}
              </div>
            </div>
            <div class="col-auto">
              <q-btn v-if="!isFinalized" color="primary" label="Tambah Komponen" icon="add" size="sm"
                @click="showAddComponent = true" />
            </div>
          </div>
        </q-card-section>
        <q-card-section>
          <q-table :rows="components" :columns="componentColumns" row-key="id" flat dense
            v-model:pagination="pagination">
            <template v-slot:body-cell-weight="props">
              <q-td>
                <q-input v-model="props.row.weight" type="number" dense outlined :disable="isFinalized"
                  style="max-width: 80px" @blur="onComponentUpdate(props.row)" />
              </q-td>
            </template>
            <template v-slot:body-cell-action="props">
              <q-td>
                <q-btn v-if="!isFinalized" color="negative" icon="delete" size="sm" flat
                  @click="removeComponent(props.rowIndex)" />
              </q-td>
            </template>
          </q-table>
        </q-card-section>
      </q-card>

      <!-- Input Nilai Siswa -->
      <q-card flat bordered>
        <q-card-section>
          <div class="row items-center">
            <div class="col">
              <div class="text-subtitle2">Input Nilai Siswa</div>
            </div>
            <div class="col-auto">
              <q-btn v-if="!isFinalized" color="secondary" label="Hitung Nilai" @click="calculateGrades"
                :loading="calculating" :disable="components.length === 0" />
            </div>
          </div>
        </q-card-section>
        <q-card-section>
          <q-table :rows="studentsWithScores" :columns="scoreColumns" row-key="id" flat bordered :loading="loading">
            <template v-slot:body-cell="props">
              <q-td v-if="props.col.name === 'name' || props.col.name === 'nis'">
                {{ props.value }}
              </q-td>
              <q-td v-else-if="props.col.name === 'finalScore'">
                <strong>{{ props.row.finalScore !== undefined ? props.row.finalScore.toFixed(2) : '-' }}</strong>
              </q-td>
              <q-td v-else-if="props.col.name === 'predicate'">
                <q-badge :color="getPredicateColor(props.row.predicate)">
                  {{ props.row.predicate || '-' }}
                </q-badge>
              </q-td>
              <q-td v-else-if="props.col.name === 'passed'">
                <q-icon :name="props.row.passed ? 'check_circle' : 'cancel'"
                  :color="props.row.passed ? 'positive' : 'negative'" />
              </q-td>
              <q-td v-else>
                <q-input v-model="props.row[props.col.name]" type="number" dense outlined
                  :disable="isFinalized || calculating" style="max-width: 80px" @blur="onScoreChange(props.row)" />
              </q-td>
            </template>
          </q-table>
        </q-card-section>
      </q-card>
    </div>

    <!-- Dialog Tambah Komponen -->
    <q-dialog v-model="showAddComponent" persistent>
      <q-card style="min-width: 400px">
        <q-card-section>
          <div class="text-h6">Tambah Komponen</div>
        </q-card-section>
        <q-card-section>
          <q-input v-model="newComponent.name" label="Nama Komponen" outlined dense />
          <q-input v-model="newComponent.weight" label="Bobot (%)" type="number" outlined dense />
          <q-input v-model="newComponent.maxScore" label="Nilai Maksimum" type="number" outlined dense />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="primary" label="Tambah" @click="addComponent" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useAssessmentStore } from '@/stores/assessment.store';
import { useQuasar } from 'quasar';

const $q = useQuasar();
const assessmentStore = useAssessmentStore();

const currentSession = computed(() => assessmentStore.currentSession);
const loading = computed(() => assessmentStore.isLoading);
const calculating = ref(false);
const isFinalized = computed(() => assessmentStore.isFinalized);
const components = computed(() => assessmentStore.components || []);
const results = computed(() => assessmentStore.results || []);

// Form untuk memilih konteks
const selectedClass = ref('cls-xi-rpl-1');
const selectedSubject = ref('sub-mathematics');
const selectedPeriod = ref('mid');
const classOptions = [
  { label: 'XI RPL 1', value: 'cls-xi-rpl-1' },
  { label: 'XI TKJ 1', value: 'cls-xi-tkj-1' },
];
const subjectOptions = [
  { label: 'Matematika', value: 'sub-mathematics' },
  { label: 'Bahasa Inggris', value: 'sub-english' },
];
const periodOptions = [
  { label: 'Mid Semester', value: 'mid' },
  { label: 'Semester Final', value: 'final' },
];

// Dialog tambah komponen
const showAddComponent = ref(false);
const newComponent = ref({ name: '', weight: 0, maxScore: 100 });

// Tabel komponen
const componentColumns = [
  { name: 'name', label: 'Nama', field: 'name', align: 'left' },
  { name: 'weight', label: 'Bobot (%)', field: 'weight', align: 'center' },
  { name: 'maxScore', label: 'Nilai Maks', field: 'maxScore', align: 'center' },
  { name: 'action', label: 'Aksi', field: 'action', align: 'center' },
];

// Tabel nilai siswa
const pagination = ref({ rowsPerPage: 0 });
const studentsWithScores = ref([]);

// Kolom untuk tabel nilai
const scoreColumns = computed(() => {
  const cols = [
    { name: 'name', label: 'Nama', field: 'name', align: 'left' },
    { name: 'nis', label: 'NIS', field: 'nis', align: 'left' },
  ];
  // Tambahkan kolom untuk setiap komponen
  for (const comp of components.value) {
    cols.push({
      name: comp.id,
      label: comp.name,
      field: comp.id,
      align: 'center',
    });
  }
  cols.push(
    { name: 'finalScore', label: 'Nilai Akhir', field: 'finalScore', align: 'center' },
    { name: 'predicate', label: 'Predikat', field: 'predicate', align: 'center' },
    { name: 'passed', label: 'Lulus', field: 'passed', align: 'center' },
  );
  return cols;
});

const totalWeight = computed(() => {
  return components.value.reduce((sum, c) => sum + (c.weight || 0), 0);
});

const loadOrCreate = async () => {
  try {
    await assessmentStore.loadOrCreateSession({
      classId: selectedClass.value,
      subjectId: selectedSubject.value,
      period: selectedPeriod.value,
    });
    // Load siswa dari mock
    await loadStudents();
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Gagal: ' + err.message });
  }
};

const loadStudents = async () => {
  try {
    // Ambil siswa dari API mock
    const response = await fetch(`/api/v1/students/class/${selectedClass.value}`);
    const data = await response.json();
    const students = data.data || [];

    // Inisialisasi studentsWithScores
    studentsWithScores.value = students.map(student => {
      const row = {
        id: student.id,
        name: student.name,
        nis: student.nis,
      };
      // Isi nilai per komponen (jika ada di results)
      const result = results.value.find(r => r.studentId === student.id);
      if (result) {
        row.finalScore = result.finalScore;
        row.predicate = result.predicate;
        row.passed = result.passed;
        // Isi nilai per komponen dari details
        if (result.details) {
          for (const comp of components.value) {
            row[comp.id] = result.details[comp.id]?.score || 0;
          }
        }
      } else {
        // Default 0
        for (const comp of components.value) {
          row[comp.id] = 0;
        }
      }
      return row;
    });
  } catch (err) {
    console.error('Failed to load students:', err);
  }
};

const addComponent = () => {
  if (!newComponent.value.name || newComponent.value.weight <= 0) {
    $q.notify({ type: 'warning', message: 'Nama dan bobot harus diisi' });
    return;
  }
  // Tambahkan ke store (via API)
  // Untuk mock, kita update local dulu
  const comp = {
    id: 'comp-' + Date.now(),
    name: newComponent.value.name,
    weight: parseFloat(newComponent.value.weight),
    maxScore: parseFloat(newComponent.value.maxScore) || 100,
    order: components.value.length,
  };
  assessmentStore.currentSession?.addComponent(comp);
  // Refresh student table untuk tambah kolom
  loadStudents();
  showAddComponent.value = false;
  newComponent.value = { name: '', weight: 0, maxScore: 100 };
};

const removeComponent = (index) => {
  assessmentStore.currentSession?.removeComponent(index);
  loadStudents();
};

const onComponentUpdate = (row) => {
  console.log(row)

  // Update weight di store
  // (tidak full implementasi untuk mock)
};

const onScoreChange = (row) => {
  console.log(row)
  // Simpan perubahan score di row
  // Tapi kita tidak langsung update ke store, nanti di calculate
};

const calculateGrades = async () => {
  calculating.value = true;
  try {
    // Kumpulkan semua nilai per siswa
    const allScores = {};
    for (const student of studentsWithScores.value) {
      const scores = {};
      for (const comp of components.value) {
        scores[comp.id] = parseFloat(student[comp.id]) || 0;
      }
      allScores[student.id] = scores;
    }
    await assessmentStore.calculateGrade(allScores);
    // Refresh data
    await loadStudents();
    $q.notify({ type: 'positive', message: 'Perhitungan nilai berhasil' });
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Gagal hitung: ' + err.message });
  } finally {
    calculating.value = false;
  }
};

const finalizeAssessment = async () => {
  try {
    await assessmentStore.finalize();
    $q.notify({ type: 'positive', message: 'Penilaian difinalisasi' });
    await loadStudents();
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Gagal finalisasi: ' + err.message });
  }
};

const closeSession = () => {
  assessmentStore.clear();
  studentsWithScores.value = [];
};

const getPredicateColor = (predicate) => {
  const map = {
    A: 'positive',
    B: 'blue',
    C: 'orange',
    D: 'negative',
    E: 'grey',
  };
  return map[predicate] || 'grey';
};
</script>
