'use client'

import { useMemo, useState } from 'react'
import { useRouter } from 'next/navigation'
import { Loader2 } from 'lucide-react'
import { login, signup, type AuthPayload } from '../lib/api'
import { Button } from './ui/button'
import { Input } from './ui/input'
import { Label } from './ui/label'

interface AuthFormProps {
    mode: 'login' | 'signup'
}

export function AuthForm({ mode }: AuthFormProps) {
    const router = useRouter()
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const [isSubmitting, setIsSubmitting] = useState(false)
    const [error, setError] = useState<string | null>(null)

    const labels = useMemo(() => {
        if (mode === 'login') {
            return {
                action: 'Log in',
                loading: 'Logging in...',
            }
        }

        return {
            action: 'Create account',
            loading: 'Creating account...',
        }
    }, [mode])

    const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault()
        setError(null)

        if (!email.trim() || !password.trim()) {
            setError('Email and password are required')
            return
        }

        if (mode === 'signup' && password.length < 6) {
            setError('Password must be at least 6 characters long')
            return
        }

        try {
            setIsSubmitting(true)
            const payload: AuthPayload = mode === 'login'
                ? await login(email.trim(), password)
                : await signup(email.trim(), password)

            localStorage.setItem('auth_token', payload.token)
            localStorage.setItem('auth_user', JSON.stringify(payload.user))

            router.push('/')
            router.refresh()
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Authentication failed')
        } finally {
            setIsSubmitting(false)
        }
    }

    return (
        <form className="space-y-5" onSubmit={handleSubmit}>
            <div className="space-y-2">
                <Label htmlFor="email">Email</Label>
                <Input
                    id="email"
                    type="email"
                    autoComplete="email"
                    placeholder="name@company.com"
                    value={email}
                    onChange={(event) => setEmail(event.target.value)}
                    disabled={isSubmitting}
                    required
                />
            </div>

            <div className="space-y-2">
                <Label htmlFor="password">Password</Label>
                <Input
                    id="password"
                    type="password"
                    autoComplete={mode === 'login' ? 'current-password' : 'new-password'}
                    placeholder="Enter your password"
                    value={password}
                    onChange={(event) => setPassword(event.target.value)}
                    disabled={isSubmitting}
                    required
                    minLength={mode === 'signup' ? 6 : undefined}
                />
            </div>

            {error && (
                <p className="rounded-md border border-down/50 bg-down/10 px-3 py-2 text-sm text-down">
                    {error}
                </p>
            )}

            <Button type="submit" className="w-full" disabled={isSubmitting}>
                {isSubmitting ? (
                    <>
                        <Loader2 className="size-4 animate-spin" />
                        {labels.loading}
                    </>
                ) : (
                    labels.action
                )}
            </Button>
        </form>
    )
}
