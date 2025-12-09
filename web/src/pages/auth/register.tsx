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
import { Link, useNavigate, useSearchParams } from "react-router-dom"
import ExternalAuth from "@/components/auth/ExternalAuth"
import React, { type FormEvent } from "react"
import useRegisterStart from "@/lib/api/hooks/auth/useRegisterStart"
import { Turnstile, type TurnstileInstance } from "@marsidev/react-turnstile"
import { TURNSTILE_KEY } from "@/lib/information"
import toast from "react-hot-toast"
import type { AppError } from "@/lib/api/client/normalizeError"
import buildError from "@/lib/helper/buildError"
import { Trans, useTranslation } from "react-i18next"

export function RegisterForm({
    className,
    ...props
}: React.ComponentProps<"form">) {
    const { t } = useTranslation();

    const [email, setEmail] = React.useState<string>("");
    const [password, setPassword] = React.useState<string>("");
    const [passwordConfirm, setPasswordConfirm] = React.useState<string>("");
    const [turnstileToken, setTurnstileToken] = React.useState<string>("");
    const [loading, setLoading] = React.useState<boolean>(false);
    const register = useRegisterStart();
    const tref = React.useRef<TurnstileInstance | null>(null)
    const navigate = useNavigate();
    const [searchParams] = useSearchParams();

    async function submit(e: FormEvent<HTMLFormElement>) {
        e.preventDefault();

        if (loading) return;

        if (password !== passwordConfirm) {
            toast.error(t("auth.passwordMatch"))
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
                    success: t("auth.code.sentCodeToast"),
                    loading: t("auth.code.loading"),
                    error: (err: AppError) => buildError(err),
                }
            )
            navigate(`/auth/register/confirm?session=${resp.session}&to=${cmail}`)
        } finally {
            setLoading(false)
        }

    }

    React.useEffect(() => {
        if (!searchParams.get("to")) {
            navigate("/auth/register")
        }
    }, [searchParams, navigate])

    return (
        <form onSubmit={submit} className={cn("flex flex-col gap-6", className)} {...props}>
            <FieldGroup>
                <div className="flex flex-col items-center gap-1 text-center">
                    <h1 className="text-2xl font-bold">{t("auth.register.title")}</h1>
                    <p className="text-muted-foreground text-sm text-balance">
                        {t("auth.register.description")}
                    </p>
                </div>
                <Field>
                    <FieldLabel htmlFor="email">{t("auth.email")}</FieldLabel>
                    <Input value={email} onChange={(e) => setEmail(e.target.value)} id="email" type="email" placeholder={t("auth.emailPlaceholder")} required />
                </Field>
                <Field>
                    <FieldLabel htmlFor="password">{t("auth.password")}</FieldLabel>
                    <Input value={password} onChange={(e) => setPassword(e.target.value)} id="password" type="password" placeholder={t("auth.passwordPlaceholder")} required />
                </Field>
                <Field>
                    <FieldLabel htmlFor="confirm-password">{t("auth.passwordConfirm")}</FieldLabel>
                    <Input value={passwordConfirm} onChange={(e) => setPasswordConfirm(e.target.value)} id="confirm-password" type="password" placeholder={t("auth.passwordConfirmPlaceholder")} required />
                </Field>
                <Field>
                    <Button type="submit">{t("auth.register.name")}</Button>
                </Field>
                <FieldSeparator>{t("common.or")}</FieldSeparator>
                <Field>
                    <ExternalAuth />
                    <FieldDescription className="text-center">
                        <Trans
                            i18nKey="auth.register.footerText"
                            components={{
                                l: <Link to={"/auth/login"} />
                            }}
                        />
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
