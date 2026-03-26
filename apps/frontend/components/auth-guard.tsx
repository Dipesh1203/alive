'use client'

import { useEffect, useState } from 'react'
import { usePathname, useRouter } from 'next/navigation'

interface AuthGuardProps {
    children: React.ReactNode
}

const PUBLIC_ROUTES = ['/login', '/signup']

function isPublicRoute(pathname: string): boolean {
    return PUBLIC_ROUTES.some((route) => pathname === route || pathname.startsWith(`${route}/`))
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

        if (token && onPublicRoute) {
            setCanRender(false)
            router.replace('/')
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
