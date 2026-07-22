export type VersionInfo = {
  service: string
  version: string
}

export type FrontendConfig = {
  service: string
  version: string
  uiPort: number
  apiBase: string
  apiProxyTarget: string
}

async function fetchData<T>(url: string): Promise<T> {
  const res = await fetch(url)
  const body = await res.json().catch(() => ({}))
  if (!res.ok) {
    throw new Error(body.error ?? `request failed: ${res.status}`)
  }
  if (body.data === undefined) {
    throw new Error(body.error ?? 'empty response')
  }
  return body.data as T
}

export class Services {
  async getVersion(): Promise<VersionInfo> {
    return fetchData<VersionInfo>('http://localhost:5173?func=getversion')
  }

  async getConfig(): Promise<FrontendConfig> {
    return fetchData<FrontendConfig>('http://localhost:5173?func=getconfig')
  }
}
