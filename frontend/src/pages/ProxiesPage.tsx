import { useState } from 'react'
import { Plus, Trash2, RefreshCw, CheckCircle, XCircle, Clock } from 'lucide-react'

interface Proxy {
  id: string
  name: string
  type: string
  host: string
  port: number
  country?: string
  lastCheckStatus: string
  lastCheckLatency: number
}

export default function ProxiesPage() {
  const [proxies, setProxies] = useState<Proxy[]>([])
  const [, setLoading] = useState(false)
  const [showAddModal, setShowAddModal] = useState(false)
  const [importText, setImportText] = useState('')

  const loadProxies = async () => {
    setLoading(true)
    try {
      // @ts-ignore
      const data = await window.go.app.App.GetProxies()
      setProxies(data || [])
    } catch (e) {
      console.error('Failed to load proxies:', e)
    }
    setLoading(false)
  }

  const checkAllProxies = async () => {
    setLoading(true)
    try {
      // @ts-ignore
      await window.go.app.App.CheckAllProxies()
      loadProxies()
    } catch (e) {
      console.error('Failed to check proxies:', e)
    }
    setLoading(false)
  }

  const importProxies = async () => {
    try {
      // @ts-ignore
      const count = await window.go.app.App.ImportProxies(importText, 'http')
      alert(`Imported ${count} proxies`)
      setImportText('')
      setShowAddModal(false)
      loadProxies()
    } catch (e) {
      console.error('Failed to import proxies:', e)
    }
  }

  const deleteProxy = async (id: string) => {
    if (!confirm('Delete this proxy?')) return
    try {
      // @ts-ignore
      await window.go.app.App.DeleteProxy(id)
      loadProxies()
    } catch (e) {
      console.error('Failed to delete proxy:', e)
    }
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'working':
        return <CheckCircle className="w-4 h-4 text-green-500" />
      case 'failed':
        return <XCircle className="w-4 h-4 text-red-500" />
      default:
        return <Clock className="w-4 h-4 text-gray-500" />
    }
  }

  return (
    <div className="p-6">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">Proxy Manager</h1>
        <div className="flex gap-2">
          <button onClick={loadProxies} className="flex items-center gap-2 px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded-lg">
            <RefreshCw className="w-4 h-4" /> Refresh
          </button>
          <button onClick={checkAllProxies} className="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-500 rounded-lg">
            <CheckCircle className="w-4 h-4" /> Check All
          </button>
          <button onClick={() => setShowAddModal(true)} className="flex items-center gap-2 px-4 py-2 bg-purple-600 hover:bg-purple-500 rounded-lg">
            <Plus className="w-4 h-4" /> Import
          </button>
        </div>
      </div>

      {showAddModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-gray-800 p-6 rounded-lg w-full max-w-lg">
            <h2 className="text-xl font-bold mb-4">Import Proxies</h2>
            <textarea
              value={importText}
              onChange={(e) => setImportText(e.target.value)}
              className="w-full h-48 bg-gray-700 rounded p-3 text-sm font-mono"
              placeholder="host:port&#10;host:port:user:pass&#10;http://user:pass@host:port"
            />
            <div className="flex justify-end gap-2 mt-4">
              <button onClick={() => setShowAddModal(false)} className="px-4 py-2 bg-gray-700 rounded-lg">Cancel</button>
              <button onClick={importProxies} className="px-4 py-2 bg-purple-600 rounded-lg">Import</button>
            </div>
          </div>
        </div>
      )}

      {proxies.length === 0 ? (
        <div className="text-center py-12 text-gray-400">No proxies. Click Import to add.</div>
      ) : (
        <table className="w-full">
          <thead>
            <tr className="text-left text-gray-400 border-b border-gray-700">
              <th className="pb-3">Status</th>
              <th className="pb-3">Name</th>
              <th className="pb-3">Type</th>
              <th className="pb-3">Host:Port</th>
              <th className="pb-3">Latency</th>
              <th className="pb-3">Actions</th>
            </tr>
          </thead>
          <tbody>
            {proxies.map((p) => (
              <tr key={p.id} className="border-b border-gray-700/50">
                <td className="py-3">{getStatusIcon(p.lastCheckStatus)}</td>
                <td className="py-3">{p.name}</td>
                <td className="py-3 uppercase text-xs">{p.type}</td>
                <td className="py-3 font-mono text-sm">{p.host}:{p.port}</td>
                <td className="py-3">{p.lastCheckLatency > 0 ? `${p.lastCheckLatency}ms` : '-'}</td>
                <td className="py-3">
                  <button onClick={() => deleteProxy(p.id)} className="p-1.5 bg-red-600 hover:bg-red-500 rounded">
                    <Trash2 className="w-3.5 h-3.5" />
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  )
}
