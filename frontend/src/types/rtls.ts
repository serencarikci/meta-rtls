export type Site = {
  id: string
  code: string
  name: string
  timezone: string
  status?: string
}

export type Building = {
  id: string
  code: string
  name: string
  siteId: string
}

export type Floor = {
  id: string
  buildingId?: string
  code: string
  name: string
  widthM: number
  heightM: number
  levelIndex?: number
}

export type Zone = {
  id: string
  code: string
  name: string
  minX: number
  minY: number
  maxX: number
  maxY: number
}

export type LivePosition = {
  tagId: string
  tagCode: string
  floorId: string
  x: number
  y: number
  zoneCode?: string
  updatedAt: string
}
