import { useState } from 'react'
import { Plus, Play, Trash2, Copy, Download, RefreshCw } from 'lucide-react'

interface Profile {
  id: string
  name: string
  tags: string[]
  createdAt: string
  lastUsedAt?: string
  isRunning?: boolean
}

export default function ProfilesPage() {
  const [profiles, setProfiles] = useState<Profile[]>([])
  const [loading, setLoading] = useState(false)

  const loadProfiles = async () => {
    setLoading(true)
    try {
      // @ts-ignore - Wails bindings
      const data = await window.go.app.App.GetProfiles()
      setProfiles(data || [])
    } catch (e) {
      console.error('Failed to load profiles:', e)
    }
    setLoading(false)
  }

  const createProfile = async () => {
    try {
      // @ts-ignore
      await window.go.app.App.GenerateRandomProfile()
      loadProfiles()
    } catch (e) {
      console.error('Failed to create profile:', e)
    }
  }

  const launchBrowser = async (id: string) => {
    try {
      // @ts-ignore
      await window.go.app.App.LaunchBrowser(id)
      loadProfiles()
    } catch (e) {
      console.error('Failed to launch browser:', e)
    }
  }

  const deleteProfile = async (id: string) => {
    if (!confirm('Are you sure you want to delete this profile?')) return
    try {
      // @ts-ignore
      await window.go.app.App.DeleteProfile(id)
      loadProfiles()
    } catch (e) {
      console.error('Failed to delete profile:', e)
    }
  }

  return (
    <div className="p-6">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">Browser Profiles</h1>
        <div className="flex gap-2">
          <button
            onClick={loadProfiles}
            className="flex items-center gap-2 px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors"
          >
            <RefreshCw className="w-4 h-4" />
            Refresh
          </button>
          <button
            onClick={createProfile}
            className="flex items-center gap-2 px-4 py-2 bg-purple-600 hover:bg-purple-500 rounded-lg transition-colors"
          >
            <Plus className="w-4 h-4" />
            New Profile
          </button>
        </div>
      </div>

      {loading ? (
        <div className="text-center py-12 text-gray-400">Loading...</div>
      ) : profiles.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-gray-400 mb-4">No profiles yet. Create your first profile to get started.</p>
          <button
            onClick={createProfile}
            className="px-6 py-3 bg-purple-600 hover:bg-purple-500 rounded-lg transition-colors"
          >
            Create Profile
          </button>
        </div>
      ) : (
        <div className="grid gap-4">
          {profiles.map((profile) => (
            <div
              key={profile.id}
              className="flex items-center justify-between p-4 bg-gray-800 rounded-lg border border-gray-700"
            >
              <div>
                <h3 className="font-medium">{profile.name}</h3>
                <p className="text-sm text-gray-400">
                  Created: {new Date(profile.createdAt).toLocaleDateString()}
                  {profile.lastUsedAt && (
                    <> · Last used: {new Date(profile.lastUsedAt).toLocaleDateString()}</>
                  )}
                </p>
                {profile.tags.length > 0 && (
                  <div className="flex gap-1 mt-2">
                    {profile.tags.map((tag) => (
                      <span
                        key={tag}
                        className="px-2 py-0.5 text-xs bg-gray-700 rounded"
                      >
                        {tag}
                      </span>
                    ))}
                  </div>
                )}
              </div>
              <div className="flex gap-2">
                <button
                  onClick={() => launchBrowser(profile.id)}
                  className="p-2 bg-green-600 hover:bg-green-500 rounded-lg transition-colors"
                  title="Launch Browser"
                >
                  <Play className="w-4 h-4" />
                </button>
                <button
                  className="p-2 bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors"
                  title="Duplicate"
                >
                  <Copy className="w-4 h-4" />
                </button>
                <button
                  className="p-2 bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors"
                  title="Export"
                >
                  <Download className="w-4 h-4" />
                </button>
                <button
                  onClick={() => deleteProfile(profile.id)}
                  className="p-2 bg-red-600 hover:bg-red-500 rounded-lg transition-colors"
                  title="Delete"
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
