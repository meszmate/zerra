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
import useRegisterStart from "@/lib/api/hooks/auth/useRegisterStart"
import { Turnstile, type TurnstileInstance } from "@marsidev/react-turnstile"
import { TURNSTILE_KEY } from "@/lib/information"
import toast from "react-hot-toast"
import type { AppError } from "@/lib/api/client/normalizeError"
import buildError from "@/lib/helper/buildError"

export function RegisterForm({
    className,
    ...props
}: React.ComponentProps<"form">) {
    const [email, setEmail] = React.useState<string>("");
    const [password, setPassword] = React.useState<string>("");
    const [passwordConfirm, setPasswordConfirm] = React.useState<string>("");
    const [turnstileToken, setTurnstileToken] = React.useState<string>("");
    const [loading, setLoading] = React.useState<boolean>(false);
    const register = useRegisterStart();
    const tref = React.useRef<TurnstileInstance | null>(null)
    const navigate = useNavigate();

    async function submit(e: FormEvent<HTMLFormElement>) {
        e.preventDefault();

        if (loading) return;

        if (password !== passwordConfirm) {
            toast.error("Nem egyeznek a jelszavak.")
            return
        }

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

            navigate(`/auth/register/confirm?session=${resp.session}&to=${cmail}`)
        } finally {
            setLoading(false)
        }

    }

    return (
        <form onSubmit={submit} className={cn("flex flex-col gap-6", className)} {...props}>
            <FieldGroup>
                <div className="flex flex-col items-center gap-1 text-center">
                    <h1 className="text-2xl font-bold">Regisztráld a fiókod</h1>
                    <p className="text-muted-foreground text-sm text-balance">
                        Add meg az információkot alul, a fiókod regisztrálásához
                    </p>
                </div>
                <Field>
                    <FieldLabel htmlFor="email">Email</FieldLabel>
                    <Input value={email} onChange={(e) => setEmail(e.target.value)} id="email" type="email" placeholder="m@example.com" required />
                </Field>
                <Field>
                    <FieldLabel htmlFor="password">Jelszó</FieldLabel>
                    <Input value={password} onChange={(e) => setPassword(e.target.value)} id="password" type="password" required />
                </Field>
                <Field>
                    <FieldLabel htmlFor="confirm-password">Jelszó Megerősítése</FieldLabel>
                    <Input value={passwordConfirm} onChange={(e) => setPasswordConfirm(e.target.value)} id="confirm-password" type="password" required />
                </Field>
                <Field>
                    <Button type="submit">Regisztráció</Button>
                </Field>
                <FieldSeparator>Vagy</FieldSeparator>
                <Field>
                    <ExternalAuth />
                    <FieldDescription className="text-center">
                        Van már fiókod?{" "}
                        <Link to="/auth/login" className="underline underline-offset-4">
                            Bejelentkezés
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
