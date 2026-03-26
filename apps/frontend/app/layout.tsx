import type { Metadata, Viewport } from 'next'
import { Cascadia_Mono } from 'next/font/google'
import { Analytics } from '@vercel/analytics/next'
import { AuthGuard } from '../components/auth-guard'
import './globals.css'

const cascadiaMono = Cascadia_Mono({
  subsets: ["latin"],
  variable: '--font-cascadia-mono',
});

export const metadata: Metadata = {
  title: 'Alive - Uptime Monitoring',
  description: 'Monitor your websites and services with real-time uptime tracking, latency monitoring, and instant alerts.',
  icons: {
    apple: '/alive_logo.png',
  },
}

export const viewport: Viewport = {
  themeColor: [
    { media: '(prefers-color-scheme: light)', color: '#fafafa' },
    { media: '(prefers-color-scheme: dark)', color: '#0a0a0f' },
  ],
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="en" className="dark">
      <body className={`${cascadiaMono.variable} font-sans antialiased`}>
        <AuthGuard>{children}</AuthGuard>
        <Analytics />
      </body>
    </html>
  )
}
