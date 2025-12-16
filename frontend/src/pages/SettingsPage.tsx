import { useState } from 'react'
import { Save, FolderOpen, Download, Upload } from 'lucide-react'

export default function SettingsPage() {
  const [ollamaUrl, setOllamaUrl] = useState('http://localhost:11434')
  const [ollamaModel, setOllamaModel] = useState('llama3.2')
  const [browserPath, setBrowserPath] = useState('')

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-6">Settings</h1>

      <div className="space-y-8 max-w-2xl">
        {/* Browser Settings */}
        <section className="bg-gray-800 rounded-lg p-6">
          <h2 className="text-lg font-semibold mb-4">Browser Settings</h2>
          
          <div className="space-y-4">
            <div>
              <label className="block text-sm text-gray-400 mb-2">Browser Executable Path</label>
              <div className="flex gap-2">
                <input
                  type="text"
                  value={browserPath}
                  onChange={(e) => setBrowserPath(e.target.value)}
                  placeholder="Auto-detect Edge browser"
                  className="flex-1 bg-gray-700 rounded px-4 py-2"
                />
                <button className="px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded">
                  <FolderOpen className="w-4 h-4" />
                </button>
              </div>
              <p className="text-xs text-gray-500 mt-1">Leave empty to auto-detect Microsoft Edge</p>
            </div>
          </div>
        </section>

        {/* AI Settings */}
        <section className="bg-gray-800 rounded-lg p-6">
          <h2 className="text-lg font-semibold mb-4">AI Settings (Ollama)</h2>
          
          <div className="space-y-4">
            <div>
              <label className="block text-sm text-gray-400 mb-2">Ollama API URL</label>
              <input
                type="text"
                value={ollamaUrl}
                onChange={(e) => setOllamaUrl(e.target.value)}
                className="w-full bg-gray-700 rounded px-4 py-2"
              />
            </div>
            
            <div>
              <label className="block text-sm text-gray-400 mb-2">Model</label>
              <select
                value={ollamaModel}
                onChange={(e) => setOllamaModel(e.target.value)}
                className="w-full bg-gray-700 rounded px-4 py-2"
              >
                <option value="llama3.2">Llama 3.2</option>
                <option value="llama3.1">Llama 3.1</option>
                <option value="mistral">Mistral</option>
                <option value="mixtral">Mixtral</option>
                <option value="phi3">Phi-3</option>
                <option value="gemma2">Gemma 2</option>
              </select>
            </div>
          </div>
        </section>

        {/* Data Management */}
        <section className="bg-gray-800 rounded-lg p-6">
          <h2 className="text-lg font-semibold mb-4">Data Management</h2>
          
          <div className="flex gap-4">
            <button className="flex items-center gap-2 px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded">
              <Download className="w-4 h-4" />
              Export All Data
            </button>
            <button className="flex items-center gap-2 px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded">
              <Upload className="w-4 h-4" />
              Import Data
            </button>
          </div>
        </section>

        {/* About */}
        <section className="bg-gray-800 rounded-lg p-6">
          <h2 className="text-lg font-semibold mb-4">About</h2>
          <div className="text-gray-400 text-sm space-y-1">
            <p>Ghost Browser v1.0.0</p>
            <p>Antidetect browser with AI personalities</p>
            <p className="mt-4">Built with Go, Wails, React, and Edge CDP</p>
          </div>
        </section>

        <button className="flex items-center gap-2 px-6 py-3 bg-purple-600 hover:bg-purple-500 rounded-lg">
          <Save className="w-4 h-4" />
          Save Settings
        </button>
      </div>
    </div>
  )
}
