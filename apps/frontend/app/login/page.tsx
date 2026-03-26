import { AuthForm } from '../../components/auth-form'
import { AuthShell } from '../../components/auth-shell'

export default function LoginPage() {
    return (
        <AuthShell
            title="Welcome back"
            subtitle="Log in to continue monitoring your websites in real time."
            footerText="New to Alive?"
            footerLinkLabel="Create an account"
            footerLinkHref="/signup"
        >
            <AuthForm mode="login" />
        </AuthShell>
    )
}
