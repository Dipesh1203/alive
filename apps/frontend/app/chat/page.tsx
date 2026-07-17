'use client'

import { FormEvent, useEffect, useMemo, useRef, useState } from 'react'
import { ArrowUpRight, Loader2, MessageSquareText, Sparkles } from 'lucide-react'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { DashboardLayout } from '../../components/dashboard-layout'
import { Button } from '../../components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../../components/ui/card'
import { Textarea } from '../../components/ui/textarea'
import { sendChatMessage } from '../../lib/api'

type ChatRole = 'user' | 'assistant'

interface ChatMessage {
  id: string
  role: ChatRole
  content: string
}

const starterMessages: ChatMessage[] = [
  {
    id: 'assistant-hello',
    role: 'assistant',
    content: 'Ask about a log sample, incident pattern, or suspected root cause. I will send it to the RAG-backed analyzer.',
  },
]

export default function ChatPage() {
  const [messages, setMessages] = useState<ChatMessage[]>(starterMessages)
  const [query, setQuery] = useState('')
  const [isSending, setIsSending] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const bottomRef = useRef<HTMLDivElement | null>(null)

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth', block: 'end' })
  }, [messages, isSending])

  const canSubmit = useMemo(() => query.trim().length > 0 && !isSending, [isSending, query])

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    const trimmedQuery = query.trim()
    if (!trimmedQuery || isSending) return

    const userMessage: ChatMessage = {
      id: `user-${Date.now()}`,
      role: 'user',
      content: trimmedQuery,
    }

    setMessages((current) => [...current, userMessage])
    setQuery('')
    setError(null)
    setIsSending(true)

    try {
      const response = await sendChatMessage(trimmedQuery)
      console.log('Chat response:', response)
      setMessages((current) => [
        ...current,
        {
          id: `assistant-${Date.now()}`,
          role: 'assistant',
          content: response.llm_output || 'No response was returned from the analyzer.',
        },
      ])
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Chat request failed'
      setError(message)
      setMessages((current) => [
        ...current,
        {
          id: `assistant-error-${Date.now()}`,
          role: 'assistant',
          content: `I could not complete that request: ${message}`,
        },
      ])
    } finally {
      setIsSending(false)
    }
  }

  return (
    <DashboardLayout
      title="AI Chat"
      breadcrumbs={[{ label: 'AI Chat' }]}
    >
      <div className="grid gap-6 xl:grid-cols-[1.15fr_0.85fr]">
        <Card className="border-border/70 bg-card/80 shadow-sm shadow-black/5">
          <CardHeader className="border-b border-border/60">
            <CardTitle className="flex items-center gap-2 text-xl">
              <MessageSquareText className="size-5 text-primary" />
              RAG chat
            </CardTitle>
            <CardDescription>
              Send a question or log sample to the backend analyzer and review the retrieved answer here.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4 pt-6">
            <div className="max-h-[60vh] space-y-4 overflow-auto pr-1">
              {messages.map((message) => (
                <div
                  key={message.id}
                  className={[
                    'flex',
                    message.role === 'user' ? 'justify-end' : 'justify-start',
                  ].join(' ')}
                >
                  <div
                    className={[
                      'max-w-[85%] rounded-2xl border px-4 py-3 text-sm leading-6 shadow-sm',
                      message.role === 'user'
                        ? 'border-primary/20 bg-primary text-primary-foreground'
                        : 'border-border/70 bg-muted/60 text-foreground',
                    ].join(' ')}
                  >
                    {/* Updated to conditionally render markdown style safely */}
                    {message.role === 'assistant' ? (
                      <div className="prose prose-sm dark:prose-invert max-w-none break-words space-y-2
                        [&_ul]:list-disc [&_ol]:list-decimal [&_ul]:ml-4 [&_ol]:ml-4
                        [&_table]:w-full [&_table]:border-collapse [&_table]:my-2
                        [&_th]:border [&_th]:border-border [&_th]:p-2 [&_th]:bg-muted/50
                        [&_td]:border [&_td]:border-border [&_td]:p-2
                        [&_code]:bg-muted [&_code]:px-1 [&_code]:py-0.5 [&_code]:rounded [&_code]:font-mono [&_code]:text-xs
                        [&_pre]:bg-zinc-950 [&_pre]:text-zinc-50 [&_pre]:p-3 [&_pre]:rounded-lg [&_pre]:overflow-x-auto [&_pre]:font-mono">
                        <ReactMarkdown remarkPlugins={[remarkGfm]}>
                          {message.content}
                        </ReactMarkdown>
                      </div>
                    ) : (
                      <p className="whitespace-pre-wrap">{message.content}</p>
                    )}
                  </div>
                </div>
              ))}
              {isSending && (
                <div className="flex justify-start">
                  <div className="inline-flex items-center gap-2 rounded-2xl border border-border/70 bg-muted/60 px-4 py-3 text-sm text-muted-foreground">
                    <Loader2 className="size-4 animate-spin" />
                    Analyzing logs...
                  </div>
                </div>
              )}
              <div ref={bottomRef} />
            </div>

            <form onSubmit={handleSubmit} className="space-y-3 border-t border-border/60 pt-4">
              <Textarea
                value={query}
                onChange={(event) => setQuery(event.target.value)}
                placeholder="Ask about an incident, paste a log snippet, or request a root cause summary..."
                className="min-h-32 resize-none"
              />
              {error && <p className="text-sm text-destructive">{error}</p>}
              <div className="flex flex-wrap items-center justify-between gap-3">
                <p className="text-xs text-muted-foreground">
                  The backend currently expects a search query and returns the model output.
                </p>
                <Button type="submit" disabled={!canSubmit}>
                  <Sparkles className="size-4" />
                  Send
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>

        <Card className="border-border/70 bg-card/60">
          <CardHeader>
            <CardTitle className="text-lg">How to use it</CardTitle>
            <CardDescription>
              This is the minimal chat shell for your current RAG pipeline.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4 text-sm text-muted-foreground">
            <div className="rounded-xl border border-border/60 bg-background/70 p-4">
              <p className="font-medium text-foreground">Input</p>
              <p className="mt-1">Type a log symptom, incident description, or a raw sample snippet.</p>
            </div>
            <div className="rounded-xl border border-border/60 bg-background/70 p-4">
              <p className="font-medium text-foreground">Backend</p>
              <p className="mt-1">The UI calls /api/chat?search=... and renders llm_output.</p>
            </div>
            <div className="rounded-xl border border-border/60 bg-background/70 p-4">
              <p className="font-medium text-foreground">Next step</p>
              <p className="mt-1">If you want multi-turn memory, the backend route needs a chat history contract, not just a single search string.</p>
            </div>
            <Button variant="outline" className="w-full justify-between" asChild>
              <a href="/dashboard">
                Back to dashboard
                <ArrowUpRight className="size-4" />
              </a>
            </Button>
          </CardContent>
        </Card>
      </div>
    </DashboardLayout>
  )
}