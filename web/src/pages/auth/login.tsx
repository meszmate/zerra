import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import {
    Field,
    FieldDescription,
    FieldGroup,
    FieldLabel,
    FieldSeparator,
} from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { Link, useNavigate } from "react-router-dom"
import ExternalAuth from "@/components/ExternalAuth"
import React, { type FormEvent } from "react"
import { Turnstile, type TurnstileInstance } from "@marsidev/react-turnstile"
import useLoginStart from "@/lib/api/hooks/auth/useLoginStart"
import toast from "react-hot-toast"
import type { AppError } from "@/lib/api/client/normalizeError"
import buildError from "@/lib/helper/buildError"
import { TURNSTILE_KEY } from "@/lib/information"

export function LoginForm({
    className,
    ...props
}: React.ComponentProps<"form">) {
    const [email, setEmail] = React.useState<string>("");
    const [password, setPassword] = React.useState<string>("");
    const [turnstileToken, setTurnstileToken] = React.useState<string>("");
    const [loading, setLoading] = React.useState<boolean>(false);
    const register = useLoginStart();
    const tref = React.useRef<TurnstileInstance | null>(null)
    const navigate = useNavigate();

    async function submit(e: FormEvent<HTMLFormElement>) {
        e.preventDefault();

        if (loading) return;

        setLoading(true)
        try {
            const cmail = email;
            const resp = await toast.promise(
                register.mutateAsync({
                    email: cmail,
                    password,
                    turnstile: turnstileToken,
                }),
                {
                    success: "A hitelesítőkód sikeresen el lett küldve az email címre.",
                    loading: "Hitelesítés...",
                    error: (err: AppError) => buildError(err),
                }
            )

            navigate(`/auth/login/confirm?session=${resp.session}&to=${cmail}`)
        } finally {
            setLoading(false)
        }
    }

    return (
        <form className={cn("flex flex-col gap-6", className)} {...props} onSubmit={submit}>
            <FieldGroup>
                <div className="flex flex-col items-center gap-1 text-center">
                    <h1 className="text-2xl font-bold">Login to your account</h1>
                    <p className="text-muted-foreground text-sm text-balance">
                        Enter your email below to login to your account
                    </p>
                </div>
                <Field>
                    <FieldLabel htmlFor="email">Email</FieldLabel>
                    <Input value={email} onChange={(e) => setEmail(e.target.value)} id="email" type="email" placeholder="m@example.com" required />
                </Field>
                <Field>
                    <div className="flex items-center">
                        <FieldLabel htmlFor="password">Jelszó</FieldLabel>
                        <a
                            href="/auth/reset-password"
                            className="ml-auto text-sm underline-offset-4 hover:underline"
                        >
                            Elfelejtetted Jelszavad?
                        </a>
                    </div>
                    <Input value={password} onChange={(e) => setPassword(e.target.value)} id="password" type="password" required />
                </Field>
                <Field>
                    <Button type="submit">Bejelentkezés</Button>
                </Field>
                <FieldSeparator>Vagy</FieldSeparator>
                <Field>
                    <ExternalAuth />
                    <FieldDescription className="text-center">
                        Nincs fiókod?{" "}
                        <Link to="/auth/register" className="underline underline-offset-4">
                            Regisztráció
                        </Link>
                    </FieldDescription>
                </Field>
                <Field>
                    <Turnstile
                        ref={tref}
                        siteKey={TURNSTILE_KEY}
                        onSuccess={(token) => setTurnstileToken(token)}
                        onExpire={() => {
                            setTurnstileToken("");
                            tref.current?.reset();
                        }}
                        options={{
                            theme: "light",
                            size: "invisible",
                            language: "hu",
                        }}
                    />
                </Field>
            </FieldGroup>
        </form>
    )
}
