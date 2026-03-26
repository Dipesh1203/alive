import { AuthForm } from '../../components/auth-form'
import { AuthShell } from '../../components/auth-shell'

export default function SignupPage() {
    return (
        <AuthShell
            title="Create your account"
            subtitle="Start tracking uptime, latency, and incidents in minutes."
            footerText="Already have an account?"
            footerLinkLabel="Log in"
            footerLinkHref="/login"
        >
            <AuthForm mode="signup" />
        </AuthShell>
    )
}
