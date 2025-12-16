import { useState } from 'react'
import { Bot, Send, RefreshCw } from 'lucide-react'

export default function AIPage() {
  const [selectedProfile, setSelectedProfile] = useState('')
  const [message, setMessage] = useState('')
  const [chatHistory, setChatHistory] = useState<{role: string, content: string}[]>([])
  const [loading, setLoading] = useState(false)

  const sendMessage = async () => {
    if (!message.trim() || !selectedProfile) return
    
    setChatHistory(prev => [...prev, { role: 'user', content: message }])
    setLoading(true)
    
    try {
      // @ts-ignore
      const response = await window.go.app.App.Chat(selectedProfile, message)
      setChatHistory(prev => [...prev, { role: 'assistant', content: response }])
    } catch (e) {
      console.error('Failed to send message:', e)
      setChatHistory(prev => [...prev, { role: 'assistant', content: 'Error: Failed to get response' }])
    }
    
    setMessage('')
    setLoading(false)
  }

  return (
    <div className="p-6 h-full flex flex-col">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">AI Personalities</h1>
        <select
          value={selectedProfile}
          onChange={(e) => setSelectedProfile(e.target.value)}
          className="bg-gray-700 rounded-lg px-4 py-2"
        >
          <option value="">Select a profile</option>
        </select>
      </div>

      <div className="flex-1 bg-gray-800 rounded-lg p-4 mb-4 overflow-auto">
        {chatHistory.length === 0 ? (
          <div className="h-full flex items-center justify-center text-gray-400">
            <div className="text-center">
              <Bot className="w-12 h-12 mx-auto mb-2 opacity-50" />
              <p>Select a profile and start chatting</p>
              <p className="text-sm mt-2">Make sure Ollama is running locally</p>
            </div>
          </div>
        ) : (
          <div className="space-y-4">
            {chatHistory.map((msg, i) => (
              <div key={i} className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}>
                <div className={`max-w-[70%] p-3 rounded-lg ${
                  msg.role === 'user' ? 'bg-purple-600' : 'bg-gray-700'
                }`}>
                  {msg.content}
                </div>
              </div>
            ))}
            {loading && (
              <div className="flex justify-start">
                <div className="bg-gray-700 p-3 rounded-lg">
                  <RefreshCw className="w-4 h-4 animate-spin" />
                </div>
              </div>
            )}
          </div>
        )}
      </div>

      <div className="flex gap-2">
        <input
          type="text"
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          onKeyDown={(e) => e.key === 'Enter' && sendMessage()}
          placeholder="Type a message..."
          className="flex-1 bg-gray-700 rounded-lg px-4 py-2"
          disabled={!selectedProfile}
        />
        <button
          onClick={sendMessage}
          disabled={!selectedProfile || loading}
          className="px-4 py-2 bg-purple-600 hover:bg-purple-500 rounded-lg disabled:opacity-50"
        >
          <Send className="w-5 h-5" />
        </button>
      </div>
    </div>
  )
}
