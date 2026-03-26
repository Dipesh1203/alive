import type { ReactNode } from 'react'
import Link from 'next/link'
import { Activity } from 'lucide-react'

interface AuthShellProps {
    title: string
    subtitle: string
    children: ReactNode
    footerText: string
    footerLinkLabel: string
    footerLinkHref: string
}

export function AuthShell({
    title,
    subtitle,
    children,
    footerText,
    footerLinkLabel,
    footerLinkHref,
}: AuthShellProps) {
    return (
        <main className="relative flex min-h-screen items-center justify-center overflow-hidden bg-background px-4 py-10 sm:px-6">
            <div className="pointer-events-none absolute inset-0">
                <div className="absolute left-1/2 top-0 h-72 w-72 -translate-x-1/2 rounded-full bg-primary/20 blur-3xl" />
                <div className="absolute bottom-0 right-0 h-72 w-72 rounded-full bg-warning/20 blur-3xl" />
                <div className="absolute bottom-10 left-0 h-56 w-56 rounded-full bg-uptime/15 blur-3xl" />
            </div>

            <section className="relative z-10 w-full max-w-md rounded-2xl border border-border/70 bg-card/95 p-6 shadow-xl backdrop-blur-sm sm:p-8">
                <div className="mb-8 flex flex-col items-center text-center">
                    <Link href="/" className="mb-4 inline-flex items-center gap-2 rounded-full border border-border/70 bg-secondary/60 px-3 py-1.5 text-sm text-muted-foreground transition-colors hover:text-foreground">
                        <span className="flex size-7 items-center justify-center rounded-full bg-uptime text-uptime-foreground">
                            <Activity className="size-4" />
                        </span>
                        Alive Monitor
                    </Link>
                    <h1 className="text-2xl font-bold tracking-tight sm:text-3xl">{title}</h1>
                    <p className="mt-2 text-sm text-muted-foreground">{subtitle}</p>
                </div>

                {children}

                <p className="mt-6 text-center text-sm text-muted-foreground">
                    {footerText}{' '}
                    <Link href={footerLinkHref} className="font-semibold text-primary hover:underline">
                        {footerLinkLabel}
                    </Link>
                </p>
            </section>
        </main>
    )
}
