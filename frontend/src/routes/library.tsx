import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { FileUp, Trash2, Download, File, Loader2, BookOpen } from 'lucide-react'

export const Route = createFileRoute('/library')({
  component: LibraryView,
})

type LibraryDocument = {
  id: number
  original_name: string
  size: number
  mime_type: string
  created_at: string
}

function formatBytes(bytes: number) {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

function LibraryView() {
  const queryClient = useQueryClient()
  const navigate = useNavigate()
  const [isDragging, setIsDragging] = useState(false)
  const [uploading, setUploading] = useState(false)
  const [ingestingId, setIngestingId] = useState<number | null>(null)

  const { data: documents, isLoading } = useQuery<LibraryDocument[]>({
    queryKey: ['library'],
    queryFn: async () => {
      const res = await fetch('/api/library')
      if (!res.ok) throw new Error('Failed to fetch')
      return res.json()
    },
  })

  const uploadMutation = useMutation({
    mutationFn: async (file: window.File) => {
      const formData = new FormData()
      formData.append('file', file)
      const res = await fetch('/api/library/upload', {
        method: 'POST',
        body: formData,
      })
      if (!res.ok) throw new Error('Upload failed')
      return res.json()
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['library'] })
    },
  })

  const deleteMutation = useMutation({
    mutationFn: async (id: number) => {
      const res = await fetch(`/api/library/${id}`, { method: 'DELETE' })
      if (!res.ok) throw new Error('Delete failed')
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['library'] })
    },
  })

  const ingestMutation = useMutation({
    mutationFn: async (id: number) => {
      setIngestingId(id)
      const res = await fetch(`/api/library/${id}/ingest`, { method: 'POST' })
      if (!res.ok) throw new Error('Ingest failed')
      
      const reader = res.body?.getReader()
      const decoder = new TextDecoder()
      let finalData = null

      if (reader) {
        let buffer = ''
        while (true) {
          const { done, value } = await reader.read()
          if (done) break
          
          buffer += decoder.decode(value, { stream: true })
          const lines = buffer.split('\n')
          buffer = lines.pop() || ''
          
          for (const line of lines) {
            if (line.startsWith('data: ')) {
              try {
                const event = JSON.parse(line.substring(6))
                if (event.status === 'processing') {
                  // Connection is alive
                  console.log('Ingestion processing...')
                } else if (event.status === 'complete') {
                  finalData = event.wiki
                } else if (event.status === 'error') {
                  throw new Error(event.detail || 'Ingestion error')
                }
              } catch (e) {
                // Ignore JSON parse errors for non-event lines or fragments
              }
            }
          }
        }
      }
      
      if (!finalData) throw new Error('Stream ended without completion data')
      return finalData
    },
    onSuccess: (data) => {
      setIngestingId(null)
      // Redirect to the newly created wiki page
      navigate({ to: `/wiki/${data.slug}` })
    },
    onError: (err: any) => {
      setIngestingId(null)
      alert(`Failed to ingest document: ${err.message}`)
    }
  })

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault()
    setIsDragging(true)
  }

  const handleDragLeave = () => {
    setIsDragging(false)
  }

  const handleDrop = async (e: React.DragEvent) => {
    e.preventDefault()
    setIsDragging(false)
    const files = Array.from(e.dataTransfer.files)
    if (files.length > 0) {
      setUploading(true)
      for (const file of files) {
        await uploadMutation.mutateAsync(file)
      }
      setUploading(false)
    }
  }

  const handleFileSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(e.target.files || [])
    if (files.length > 0) {
      setUploading(true)
      for (const file of files) {
        await uploadMutation.mutateAsync(file)
      }
      setUploading(false)
    }
  }

  const handleDownload = async (id: number, filename: string) => {
    try {
      const res = await fetch(`/api/library/${id}/download`)
      if (!res.ok) throw new Error('Download failed')
      const blob = await res.blob()
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', filename)
      document.body.appendChild(link)
      link.click()
      link.parentNode?.removeChild(link)
    } catch (err) {
      console.error(err)
    }
  }

  return (
    <div className="p-8 max-w-6xl mx-auto">
      <h1 className="text-3xl font-bold mb-6 text-gray-800 dark:text-gray-100">Local Library</h1>
      
      {/* Upload Zone */}
      <div
        className={`border-2 border-dashed rounded-xl p-10 text-center transition-colors ${
          isDragging ? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20' : 'border-gray-300 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800/50'
        }`}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
      >
        <FileUp className="mx-auto h-12 w-12 text-gray-400 mb-4" />
        <p className="text-lg text-gray-600 dark:text-gray-300 mb-2">
          Drag and drop ebooks, PDFs, or documents here
        </p>
        <p className="text-sm text-gray-500 dark:text-gray-400 mb-4">
          Files are stored securely and never published publicly.
        </p>
        <label className="cursor-pointer inline-flex items-center px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg transition-colors">
          Browse Files
          <input type="file" className="hidden" multiple onChange={handleFileSelect} disabled={uploading} />
        </label>
        {uploading && (
          <div className="mt-4 flex items-center justify-center text-blue-600">
            <Loader2 className="animate-spin mr-2 h-5 w-5" />
            <span>Uploading...</span>
          </div>
        )}
      </div>

      {/* Document List */}
      <div className="mt-10">
        <h2 className="text-xl font-semibold mb-4 text-gray-800 dark:text-gray-100">Your Documents</h2>
        {isLoading ? (
          <div className="text-center text-gray-500 py-8">Loading...</div>
        ) : documents && documents.length > 0 ? (
          <div className="bg-white dark:bg-gray-800 rounded-xl shadow overflow-hidden">
            <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
              <thead className="bg-gray-50 dark:bg-gray-900/50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">File Name</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Size</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Uploaded</th>
                  <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Actions</th>
                </tr>
              </thead>
              <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
                {documents.map((doc) => (
                  <tr key={doc.id} className="hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center">
                        <File className="flex-shrink-0 h-5 w-5 text-gray-400 mr-3" />
                        <span className="text-sm font-medium text-gray-900 dark:text-gray-100">{doc.original_name}</span>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                      {formatBytes(doc.size)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                      {new Date(doc.created_at).toLocaleDateString()}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium space-x-3">
                      <button
                        onClick={() => ingestMutation.mutate(doc.id)}
                        disabled={ingestingId === doc.id}
                        className="text-green-600 hover:text-green-900 dark:hover:text-green-400 transition-colors disabled:opacity-50"
                        title="Ingest to Global Wiki"
                      >
                        {ingestingId === doc.id ? (
                          <Loader2 className="h-5 w-5 inline animate-spin" />
                        ) : (
                          <BookOpen className="h-5 w-5 inline" />
                        )}
                      </button>
                      <button
                        onClick={() => handleDownload(doc.id, doc.original_name)}
                        className="text-blue-600 hover:text-blue-900 dark:hover:text-blue-400 transition-colors"
                        title="Download"
                      >
                        <Download className="h-5 w-5 inline" />
                      </button>
                      <button
                        onClick={() => {
                          if (confirm('Are you sure you want to delete this document?')) {
                            deleteMutation.mutate(doc.id)
                          }
                        }}
                        className="text-red-600 hover:text-red-900 dark:hover:text-red-400 transition-colors"
                        title="Delete"
                      >
                        <Trash2 className="h-5 w-5 inline" />
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <div className="text-center py-12 bg-white dark:bg-gray-800 rounded-xl shadow">
            <p className="text-gray-500 dark:text-gray-400">No documents in your library yet.</p>
          </div>
        )}
      </div>
    </div>
  )
}
