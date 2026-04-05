"use client"

import { useEffect, useMemo, useRef, useState } from 'react'
import Link from 'next/link'
import { gsap } from 'gsap'
import { ScrollTrigger } from 'gsap/ScrollTrigger'
import {
  ArrowRight,
  Bolt,
  CheckCircle2,
  Clock3,
  Globe2,
  ShieldCheck,
  Sparkles,
  Star,
} from 'lucide-react'
import {
  fetchLandingFaqs,
  fetchLandingOverview,
  fetchLandingPricing,
  fetchLandingTestimonials,
  type LandingFAQ,
  type LandingFeature,
  type LandingPricingPlan,
  type LandingStat,
  type LandingTestimonial,
} from '../lib/api'
import { Button } from '../components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../components/ui/card'
import { Badge } from '../components/ui/badge'

gsap.registerPlugin(ScrollTrigger)

type BillingCycle = 'monthly' | 'yearly'

export default function LandingPage() {
  const rootRef = useRef<HTMLDivElement | null>(null)
  const [hasSession, setHasSession] = useState(false)
  const [billing, setBilling] = useState<BillingCycle>('monthly')
  const [stats, setStats] = useState<LandingStat[]>([])
  const [features, setFeatures] = useState<LandingFeature[]>([])
  const [pricing, setPricing] = useState<LandingPricingPlan[]>([])
  const [testimonials, setTestimonials] = useState<LandingTestimonial[]>([])
  const [faqs, setFaqs] = useState<LandingFAQ[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [isPricingLoading, setIsPricingLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const token = localStorage.getItem('auth_token')
    setHasSession(Boolean(token))
  }, [])

  useEffect(() => {
    let isMounted = true

    const loadLandingData = async () => {
      setIsLoading(true)
      setError(null)

      try {
        const [overview, testimonialsData, faqData] = await Promise.all([
          fetchLandingOverview(billing),
          fetchLandingTestimonials(),
          fetchLandingFaqs(),
        ])

        if (!isMounted) return

        setStats(overview.stats)
        setFeatures(overview.features)
        setPricing(overview.pricing)
        setTestimonials(testimonialsData.length ? testimonialsData : overview.testimonials)
        setFaqs(faqData.length ? faqData : overview.faqs)
      } catch (err) {
        if (!isMounted) return
        setError(err instanceof Error ? err.message : 'Failed to load landing page')
      } finally {
        if (isMounted) {
          setIsLoading(false)
        }
      }
    }

    void loadLandingData()

    return () => {
      isMounted = false
    }
  }, [billing])

  useEffect(() => {
    if (!rootRef.current || isLoading) return

    const ctx = gsap.context(() => {
      gsap.fromTo('[data-hero-chip]', { opacity: 0, y: -16 }, { opacity: 1, y: 0, duration: 0.7, ease: 'power2.out' })
      gsap.fromTo('[data-hero-title]', { opacity: 0, y: 26 }, { opacity: 1, y: 0, duration: 0.9, delay: 0.1, ease: 'power3.out' })
      gsap.fromTo('[data-hero-copy]', { opacity: 0, y: 18 }, { opacity: 1, y: 0, duration: 0.8, delay: 0.25, ease: 'power2.out' })
      gsap.fromTo('[data-hero-cta]', { opacity: 0, y: 14 }, { opacity: 1, y: 0, duration: 0.75, delay: 0.35, ease: 'power2.out' })
      gsap.fromTo('[data-hero-card]', { opacity: 0, x: 20, rotate: 1 }, { opacity: 1, x: 0, rotate: 0, duration: 0.9, delay: 0.2, ease: 'power3.out' })

      gsap.utils.toArray<HTMLElement>('[data-reveal]').forEach((element) => {
        gsap.fromTo(
          element,
          { opacity: 0, y: 30 },
          {
            opacity: 1,
            y: 0,
            duration: 0.85,
            ease: 'power3.out',
            scrollTrigger: {
              trigger: element,
              start: 'top 88%',
            },
          },
        )
      })
    }, rootRef)

    return () => {
      ctx.revert()
    }
  }, [isLoading, stats.length, features.length, testimonials.length, pricing.length, faqs.length])

  const yearlySavingNote = useMemo(() => {
    if (billing !== 'yearly') return null
    return 'Yearly billing applies two months free equivalent.'
  }, [billing])

  const refreshPricing = async (nextBilling: BillingCycle) => {
    setBilling(nextBilling)
    setIsPricingLoading(true)
    try {
      const plans = await fetchLandingPricing(nextBilling)
      setPricing(plans)
    } catch {
      // Keep previously loaded pricing when refresh fails.
    } finally {
      setIsPricingLoading(false)
    }
  }

  return (
    <div ref={rootRef} className="relative min-h-screen overflow-x-clip bg-background text-foreground">
      <div className="pointer-events-none absolute inset-0 opacity-60" style={{ backgroundImage: 'radial-gradient(circle at 20% 15%, color-mix(in oklch, var(--primary) 22%, transparent), transparent 32%), radial-gradient(circle at 80% 10%, color-mix(in oklch, var(--warning) 20%, transparent), transparent 38%)' }} />
      <div className="pointer-events-none absolute -left-24 top-48 h-72 w-72 rounded-full bg-primary/15 blur-3xl" />
      <div className="pointer-events-none absolute right-0 top-96 h-72 w-72 rounded-full bg-uptime/10 blur-3xl" />

      <header className="sticky top-0 z-30 border-b border-border/60 bg-background/75 backdrop-blur-md">
        <div className="mx-auto flex h-16 w-full max-w-6xl items-center justify-between px-4 sm:px-6 lg:px-8">
          <Link href="/" className="inline-flex items-center gap-2 text-sm font-semibold tracking-wide">
            <span className="inline-flex size-8 items-center justify-center rounded-md bg-uptime text-uptime-foreground">
              <ShieldCheck className="size-4" />
            </span>
            Alive
          </Link>
          <div className="flex items-center gap-2">
            <Button variant="ghost" asChild>
              <Link href="/login">Log in</Link>
            </Button>
            <Button asChild>
              <Link href={hasSession ? '/dashboard' : '/signup'}>
                {hasSession ? 'Open dashboard' : 'Start free'}
              </Link>
            </Button>
          </div>
        </div>
      </header>

      <main className="mx-auto flex w-full max-w-6xl flex-col gap-14 px-4 pb-20 pt-14 sm:px-6 lg:px-8 lg:pt-16">
        <section className="grid gap-8 lg:grid-cols-[1.2fr_0.8fr] lg:items-center">
          <div>
            <Badge data-hero-chip variant="outline" className="rounded-full border-primary/30 bg-primary/10 px-3 py-1 text-xs text-primary">
              <Sparkles className="mr-1 size-3.5" />
              Built for modern uptime teams
            </Badge>

            <h1 data-hero-title className="mt-4 text-4xl font-bold tracking-tight sm:text-5xl lg:text-6xl">
              Monitoring that feels
              <span className="block text-primary">calm when traffic spikes.</span>
            </h1>

            <p data-hero-copy className="mt-5 max-w-xl text-base text-muted-foreground sm:text-lg">
              Alive watches your websites across global regions, detects incidents fast, and gives your team the context needed to fix issues before users notice.
            </p>

            <div data-hero-cta className="mt-7 flex flex-wrap items-center gap-3">
              <Button size="lg" asChild>
                <Link href={hasSession ? '/dashboard' : '/signup'}>
                  {hasSession ? 'Go to dashboard' : 'Create account'}
                  <ArrowRight className="size-4" />
                </Link>
              </Button>
              <Button size="lg" variant="outline" asChild>
                <Link href="/login">Live product preview</Link>
              </Button>
            </div>
          </div>

          <Card data-hero-card className="border-primary/20 bg-card/70 shadow-2xl shadow-primary/10">
            <CardHeader>
              <CardTitle className="text-xl">Live reliability snapshot</CardTitle>
              <CardDescription>Updated from the landing content API</CardDescription>
            </CardHeader>
            <CardContent className="grid gap-3">
              {isLoading && <p className="text-sm text-muted-foreground">Loading metrics...</p>}
              {error && <p className="text-sm text-down">{error}</p>}
              {!isLoading && !error && stats.slice(0, 3).map((metric) => (
                <div key={metric.label} className="rounded-lg border border-border/60 bg-background/40 p-3">
                  <p className="text-xs uppercase tracking-wide text-muted-foreground">{metric.label}</p>
                  <p className="mt-1 text-2xl font-semibold text-foreground">{metric.value}</p>
                  <p className="mt-1 text-xs text-muted-foreground">{metric.hint}</p>
                </div>
              ))}
            </CardContent>
          </Card>
        </section>

        <section data-reveal>
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
            {stats.map((metric) => (
              <Card key={metric.label} className="border-border/60 bg-card/70">
                <CardHeader className="pb-2">
                  <CardDescription className="text-xs uppercase tracking-wide">{metric.label}</CardDescription>
                </CardHeader>
                <CardContent>
                  <p className="text-3xl font-semibold">{metric.value}</p>
                  <p className="mt-1 text-xs text-muted-foreground">{metric.hint}</p>
                </CardContent>
              </Card>
            ))}
          </div>
        </section>

        <section data-reveal className="space-y-5">
          <div className="max-w-xl">
            <p className="text-sm uppercase tracking-wide text-primary">Why teams choose Alive</p>
            <h2 className="mt-2 text-3xl font-semibold tracking-tight">Designed for diagnosis, not dashboards-for-show.</h2>
          </div>
          <div className="grid gap-4 md:grid-cols-2">
            {features.map((feature, index) => (
              <Card key={feature.title} className="border-border/70 bg-card/75">
                <CardHeader>
                  <CardTitle className="flex items-center gap-2 text-lg">
                    {index % 2 === 0 ? <Bolt className="size-4 text-primary" /> : <Globe2 className="size-4 text-warning" />}
                    {feature.title}
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-sm text-muted-foreground">{feature.description}</p>
                </CardContent>
              </Card>
            ))}
          </div>
        </section>

        <section data-reveal className="space-y-5">
          <div className="flex flex-wrap items-center justify-between gap-3">
            <div>
              <p className="text-sm uppercase tracking-wide text-primary">Pricing</p>
              <h2 className="mt-1 text-3xl font-semibold tracking-tight">Simple plans with room to scale</h2>
            </div>
            <div className="inline-flex rounded-full border border-border/70 bg-secondary/50 p-1">
              <button
                type="button"
                onClick={() => void refreshPricing('monthly')}
                className={`rounded-full px-4 py-1.5 text-sm transition ${billing === 'monthly' ? 'bg-primary text-primary-foreground' : 'text-muted-foreground hover:text-foreground'}`}
              >
                Monthly
              </button>
              <button
                type="button"
                onClick={() => void refreshPricing('yearly')}
                className={`rounded-full px-4 py-1.5 text-sm transition ${billing === 'yearly' ? 'bg-primary text-primary-foreground' : 'text-muted-foreground hover:text-foreground'}`}
              >
                Yearly
              </button>
            </div>
          </div>

          {yearlySavingNote && <p className="text-sm text-muted-foreground">{yearlySavingNote}</p>}
          {isPricingLoading && <p className="text-sm text-muted-foreground">Refreshing pricing...</p>}

          <div className="grid gap-4 lg:grid-cols-3">
            {pricing.map((plan) => (
              <Card key={plan.name} className={`relative border-border/70 bg-card/80 ${plan.popular ? 'border-primary/70 shadow-lg shadow-primary/20' : ''}`}>
                {plan.popular && (
                  <Badge className="absolute right-4 top-4 rounded-full bg-primary text-primary-foreground">Most popular</Badge>
                )}
                <CardHeader>
                  <CardTitle>{plan.name}</CardTitle>
                  <CardDescription>{plan.description}</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="mb-4">
                    <span className="text-4xl font-semibold">${plan.price}</span>
                    <span className="ml-1 text-sm text-muted-foreground">/{plan.interval === 'yearly' ? 'year' : 'month'}</span>
                  </div>
                  <ul className="space-y-2">
                    {plan.features.map((feature) => (
                      <li key={feature} className="flex items-start gap-2 text-sm text-muted-foreground">
                        <CheckCircle2 className="mt-0.5 size-4 text-uptime" />
                        {feature}
                      </li>
                    ))}
                  </ul>
                  <Button className="mt-5 w-full" variant={plan.popular ? 'default' : 'outline'} asChild>
                    <Link href={hasSession ? '/dashboard' : '/signup'}>{plan.cta}</Link>
                  </Button>
                </CardContent>
              </Card>
            ))}
          </div>
        </section>

        <section data-reveal className="space-y-5">
          <div>
            <p className="text-sm uppercase tracking-wide text-primary">Testimonials</p>
            <h2 className="mt-1 text-3xl font-semibold tracking-tight">Teams running critical workloads trust Alive</h2>
          </div>
          <div className="grid gap-4 lg:grid-cols-3">
            {testimonials.map((testimonial) => (
              <Card key={`${testimonial.name}-${testimonial.company}`} className="border-border/70 bg-card/75">
                <CardHeader>
                  <CardTitle className="text-lg">{testimonial.name}</CardTitle>
                  <CardDescription>{testimonial.role} at {testimonial.company}</CardDescription>
                </CardHeader>
                <CardContent>
                  <p className="text-sm text-muted-foreground">"{testimonial.quote}"</p>
                  <div className="mt-3 flex items-center justify-between">
                    <div className="flex items-center gap-1 text-warning">
                      {Array.from({ length: testimonial.rating }).map((_, index) => (
                        <Star key={`${testimonial.name}-star-${index}`} className="size-3.5 fill-current" />
                      ))}
                    </div>
                    <span className="text-xs text-muted-foreground">{testimonial.location}</span>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </section>

        <section data-reveal className="space-y-5">
          <div>
            <p className="text-sm uppercase tracking-wide text-primary">FAQ</p>
            <h2 className="mt-1 text-3xl font-semibold tracking-tight">Questions from teams evaluating Alive</h2>
          </div>
          <div className="grid gap-4 md:grid-cols-2">
            {faqs.map((faq) => (
              <Card key={faq.question} className="border-border/70 bg-card/75">
                <CardHeader>
                  <CardTitle className="text-base">{faq.question}</CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-sm text-muted-foreground">{faq.answer}</p>
                </CardContent>
              </Card>
            ))}
          </div>
        </section>

        <section data-reveal className="relative overflow-hidden rounded-2xl border border-primary/30 bg-gradient-to-r from-primary/20 via-uptime/10 to-warning/15 p-8">
          <div className="max-w-2xl">
            <p className="inline-flex items-center gap-2 text-sm text-primary">
              <Clock3 className="size-4" />
              Start in minutes
            </p>
            <h2 className="mt-2 text-3xl font-semibold tracking-tight">Get real-time confidence in your production uptime.</h2>
            <p className="mt-3 text-sm text-muted-foreground">
              Connect your first endpoint, configure alerts, and keep everyone aligned with incident timelines and regional latency insight.
            </p>
            <div className="mt-5 flex flex-wrap gap-3">
              <Button size="lg" asChild>
                <Link href={hasSession ? '/dashboard' : '/signup'}>
                  {hasSession ? 'Open dashboard' : 'Try Alive now'}
                </Link>
              </Button>
              <Button size="lg" variant="outline" asChild>
                <Link href="/login">Talk to your account</Link>
              </Button>
            </div>
          </div>
        </section>
      </main>
    </div>
  )
}
