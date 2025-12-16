// import React from 'react'
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import Sidebar from './components/Sidebar'
import ProfilesPage from './pages/ProfilesPage'
import ProxiesPage from './pages/ProxiesPage'
import SettingsPage from './pages/SettingsPage'
import AIPage from './pages/AIPage'

const queryClient = new QueryClient()

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <div className="flex h-screen bg-gray-900 text-white">
          <Sidebar />
          <main className="flex-1 overflow-auto">
            <Routes>
              <Route path="/" element={<ProfilesPage />} />
              <Route path="/profiles" element={<ProfilesPage />} />
              <Route path="/proxies" element={<ProxiesPage />} />
              <Route path="/ai" element={<AIPage />} />
              <Route path="/settings" element={<SettingsPage />} />
            </Routes>
          </main>
        </div>
      </BrowserRouter>
    </QueryClientProvider>
  )
}
