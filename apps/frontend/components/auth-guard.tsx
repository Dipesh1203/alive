'use client'

import { useEffect, useState } from 'react'
import { usePathname, useRouter } from 'next/navigation'

interface AuthGuardProps {
    children: React.ReactNode
}

const PUBLIC_ROUTES = ['/', '/login', '/signup']
const AUTH_ROUTES = ['/login', '/signup']

function isPublicRoute(pathname: string): boolean {
    if (pathname === '/') {
        return true
    }

    return PUBLIC_ROUTES
        .filter((route) => route !== '/')
        .some((route) => pathname === route || pathname.startsWith(`${route}/`))
}

function isAuthRoute(pathname: string): boolean {
    return AUTH_ROUTES.some((route) => pathname === route || pathname.startsWith(`${route}/`))
}

export function AuthGuard({ children }: AuthGuardProps) {
    const pathname = usePathname()
    const router = useRouter()
    const [canRender, setCanRender] = useState(false)

    useEffect(() => {
        const token = localStorage.getItem('auth_token')
        const onPublicRoute = isPublicRoute(pathname)

        if (!token && !onPublicRoute) {
            setCanRender(false)
            router.replace('/login')
            return
        }

        if (token && isAuthRoute(pathname)) {
            setCanRender(false)
            router.replace('/dashboard')
            return
        }

        setCanRender(true)
    }, [pathname, router])

    if (!canRender) {
        return (
            <div className="flex min-h-screen items-center justify-center bg-background px-4 text-sm text-muted-foreground">
                Checking your session...
            </div>
        )
    }

    return <>{children}</>
}
