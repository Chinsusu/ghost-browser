import { Link, useLocation } from 'react-router-dom'
import { Users, Globe, Bot, Settings, Ghost } from 'lucide-react'
import { clsx } from 'clsx'

const navItems = [
  { path: '/profiles', icon: Users, label: 'Profiles' },
  { path: '/proxies', icon: Globe, label: 'Proxies' },
  { path: '/ai', icon: Bot, label: 'AI Personalities' },
  { path: '/settings', icon: Settings, label: 'Settings' },
]

export default function Sidebar() {
  const location = useLocation()

  return (
    <aside className="w-64 bg-gray-800 border-r border-gray-700 flex flex-col">
      <div className="p-4 border-b border-gray-700">
        <div className="flex items-center gap-2">
          <Ghost className="w-8 h-8 text-purple-500" />
          <span className="text-xl font-bold">Ghost Browser</span>
        </div>
      </div>

      <nav className="flex-1 p-4">
        <ul className="space-y-2">
          {navItems.map((item) => {
            const Icon = item.icon
            const isActive = location.pathname === item.path || 
                           (item.path === '/profiles' && location.pathname === '/')
            
            return (
              <li key={item.path}>
                <Link
                  to={item.path}
                  className={clsx(
                    'flex items-center gap-3 px-4 py-2 rounded-lg transition-colors',
                    isActive
                      ? 'bg-purple-600 text-white'
                      : 'text-gray-400 hover:bg-gray-700 hover:text-white'
                  )}
                >
                  <Icon className="w-5 h-5" />
                  <span>{item.label}</span>
                </Link>
              </li>
            )
          })}
        </ul>
      </nav>

      <div className="p-4 border-t border-gray-700 text-xs text-gray-500">
        <p>Version 1.0.0</p>
        <p>© 2024 Ghost Browser</p>
      </div>
    </aside>
  )
}
