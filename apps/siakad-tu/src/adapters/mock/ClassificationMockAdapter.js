import {
  createClassificationNode,
  createArchiveLocationItem,
} from '../../contracts/classificationContract.js'

let treeDb = [
  {
    id: '1',
    parentId: null,
    code: '400',
    label: 'Pendidikan',
    children: [
      {
        id: '2',
        parentId: '1',
        code: '421',
        label: 'Sekolah',
        children: [
          { id: '3', parentId: '2', code: '421.3', label: 'Kurikulum', children: [] },
          { id: '4', parentId: '2', code: '421.4', label: 'Kesiswaan', children: [] },
        ],
      },
    ],
  },
]

let locationDb = [
  { id: '1', schoolId: 'SCH-1', name: 'Laci 1', description: 'Surat Masuk 2026' },
  { id: '2', schoolId: 'SCH-1', name: 'Ordner KUR-2026', description: 'Arsip Kurikulum' },
]

const delay = (ms) => new Promise((resolve) => setTimeout(resolve, ms))

export const ClassificationMockAdapter = {
  async getTree(context) {
    console.log(context)
    await delay(300)
    return treeDb.map(mapToContractNode)
  },

  async createNode(context, payload) {
    await delay(400)
    const newNode = {
      id: String(Date.now()),
      parentId: payload.parentId,
      code: payload.code,
      label: payload.label,
      children: [],
    }
    if (!payload.parentId) {
      treeDb.push(newNode)
    } else {
      const parent = findNode(treeDb, payload.parentId)
      if (parent) parent.children.push(newNode)
    }
    return newNode
  },

  async updateNode(context, id, payload) {
    await delay(300)
    const node = findNode(treeDb, id)
    if (node) {
      node.code = payload.code
      node.label = payload.label
    }
    return node
  },

  async deleteNode(context, id) {
    await delay(300)
    deleteNode(treeDb, id)
    return true
  },

  async listLocations(context) {
    await delay(200)
    return locationDb.filter((l) => l.schoolId === context.schoolId).map(createArchiveLocationItem)
  },

  async createLocation(context, payload) {
    await delay(300)
    const newLoc = { id: String(Date.now()), schoolId: context.schoolId, ...payload }
    locationDb.push(newLoc)
    return createArchiveLocationItem(newLoc)
  },

  // BARU: Update lokasi fisik
  async updateLocation(context, id, payload) {
    await delay(300)
    const loc = locationDb.find((l) => l.id === id)
    if (loc) {
      loc.name = payload.name
      loc.description = payload.description
    }
    return createArchiveLocationItem(loc)
  },

  // BARU: Hapus lokasi fisik
  async deleteLocation(context, id) {
    await delay(300)
    const index = locationDb.findIndex((l) => l.id === id)
    if (index !== -1) locationDb.splice(index, 1)
    return true
  },

  async getKlasifikasiOptions(context) {
    console.log(context)
    await delay(100)
    const flat = []
    const flatten = (nodes) =>
      nodes.forEach((n) => {
        flat.push(`${n.code} - ${n.label}`)
        if (n.children) flatten(n.children)
      })
    flatten(treeDb)
    return flat
  },
}

function mapToContractNode(node) {
  return createClassificationNode({
    id: node.id,
    parentId: node.parentId,
    code: node.code,
    label: node.label,
    children: node.children.map(mapToContractNode),
  })
}

function findNode(nodes, id) {
  for (const node of nodes) {
    if (node.id === id) return node
    if (node.children) {
      const found = findNode(node.children, id)
      if (found) return found
    }
  }
  return null
}

function deleteNode(nodes, id) {
  const index = nodes.findIndex((n) => n.id === id)
  if (index !== -1) {
    nodes.splice(index, 1)
    return true
  }
  for (const node of nodes) {
    if (deleteNode(node.children, id)) return true
  }
  return false
}
