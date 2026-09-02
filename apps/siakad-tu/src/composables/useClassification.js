import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'
import { useClassificationStore } from '../stores/ClassificationStore.js'
import { useOperationalContext } from './useOperationalContext.js'
import {
  createClassificationCreatePayload,
  createLocationCreatePayload,
  createLocationUpdatePayload,
} from '../contracts/classificationContract.js'

export function useClassification() {
  const $q = useQuasar()
  const store = useClassificationStore()
  const context = useOperationalContext()

  const tab = ref('kode')

  // Tree states
  const dialogNode = ref(false)
  const isEditingNode = ref(false)
  const selectedNodeId = ref(null)
  const nodeForm = ref({ kode: '', label: '', parentId: null })

  // Location states
  const dialogRak = ref(false)
  const isEditLokasi = ref(false)
  const editIdLokasi = ref(null)
  const lokasiForm = ref({ nama: '', deskripsi: '' })

  const treeData = computed(() => store.tree)
  const lokasiList = computed(() => store.locations)
  const isLoading = computed(() => store.isLoading)

  async function loadTree() {
    await store.loadTree(context)
  }

  async function loadLocations() {
    await store.loadLocations(context)
  }

  function tambahRoot() {
    isEditingNode.value = false
    selectedNodeId.value = null
    nodeForm.value = { kode: '', label: '', parentId: null }
    dialogNode.value = true
  }

  function editNode(node) {
    isEditingNode.value = true
    selectedNodeId.value = node.id
    nodeForm.value = { kode: node.code, label: node.label, parentId: node.parentId }
    dialogNode.value = true
  }

  async function simpanNode() {
    try {
      const payload = createClassificationCreatePayload({
        parentId: nodeForm.value.parentId,
        code: nodeForm.value.kode,
        label: nodeForm.value.label,
      })

      if (isEditingNode.value) {
        await store.updateNode(context, selectedNodeId.value, payload)
        $q.notify({ type: 'positive', message: 'Kode klasifikasi berhasil diperbarui' })
      } else {
        await store.createNode(context, payload)
        $q.notify({ type: 'positive', message: 'Kode klasifikasi berhasil ditambahkan' })
      }
      dialogNode.value = false
    } catch (e) {
      $q.notify({ type: 'negative', message: e.message || 'Gagal menyimpan kode' })
    }
  }

  function hapusNode(node) {
    $q.dialog({
      title: 'Konfirmasi',
      message: `Hapus kode "${node.code} - ${node.label}" dan semua sub-kodenya?`,
      cancel: true,
      persistent: true,
    }).onOk(async () => {
      try {
        await store.deleteNode(context, node.id)
        $q.notify({ type: 'positive', message: 'Kode berhasil dihapus' })
      } catch (e) {
        $q.notify({ type: 'negative', message: e.message || 'Gagal menghapus kode' })
      }
    })
  }

  function tambahLokasi() {
    isEditLokasi.value = false
    editIdLokasi.value = null
    lokasiForm.value = { nama: '', deskripsi: '' }
    dialogRak.value = true
  }

  function editLokasi(row) {
    isEditLokasi.value = true
    editIdLokasi.value = row.id
    lokasiForm.value = { nama: row.name, deskripsi: row.description }
    dialogRak.value = true
  }

  async function simpanLokasi() {
    try {
      if (isEditLokasi.value) {
        const payload = createLocationUpdatePayload({
          name: lokasiForm.value.nama,
          description: lokasiForm.value.deskripsi,
        })
        await store.updateLocation(context, editIdLokasi.value, payload)
        $q.notify({ type: 'positive', message: 'Lokasi berhasil diperbarui' })
      } else {
        const payload = createLocationCreatePayload({
          name: lokasiForm.value.nama,
          description: lokasiForm.value.deskripsi,
        })
        await store.createLocation(context, payload)
        $q.notify({ type: 'positive', message: 'Lokasi berhasil ditambahkan' })
      }
      dialogRak.value = false
    } catch (e) {
      $q.notify({ type: 'negative', message: e.message || 'Gagal menyimpan lokasi' })
    }
  }

  function hapusLokasi(id) {
    $q.dialog({
      title: 'Konfirmasi',
      message: 'Hapus lokasi fisik ini?',
      cancel: true,
      persistent: true,
    }).onOk(async () => {
      try {
        await store.deleteLocation(context, id)
        $q.notify({ type: 'positive', message: 'Lokasi berhasil dihapus' })
      } catch (e) {
        $q.notify({ type: 'negative', message: e.message || 'Gagal menghapus lokasi' })
      }
    })
  }

  return {
    tab,
    treeData,
    lokasiList,
    isLoading,
    dialogNode,
    isEditingNode,
    nodeForm,
    dialogRak,
    isEditLokasi,
    lokasiForm,
    loadTree,
    loadLocations,
    tambahRoot,
    editNode,
    simpanNode,
    hapusNode,
    tambahLokasi,
    editLokasi,
    simpanLokasi,
    hapusLokasi,
  }
}
