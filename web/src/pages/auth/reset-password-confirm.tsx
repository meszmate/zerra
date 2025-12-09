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
import { Link } from "react-router-dom"
import ExternalAuth from "@/components/auth/ExternalAuth"
import React, { type FormEvent } from "react"
import { Turnstile, type TurnstileInstance } from "@marsidev/react-turnstile"
import toast from "react-hot-toast"
import type { AppError } from "@/lib/api/client/normalizeError"
import buildError from "@/lib/helper/buildError"
import { TURNSTILE_KEY } from "@/lib/information"
import { Spinner } from "@/components/ui/spinner"
import useResetPasswordStart from "@/lib/api/hooks/auth/useResetPasswordStart"
import { Trans, useTranslation } from "react-i18next"

export function ResetPasswordConfirm({
    className,
    ...props
}: React.ComponentProps<"form">) {
    const { t } = useTranslation();

    const [email, setEmail] = React.useState<string>("");
    const [turnstileToken, setTurnstileToken] = React.useState<string>("");
    const [loading, setLoading] = React.useState<boolean>(false);
    const resetPassword = useResetPasswordStart();
    const tref = React.useRef<TurnstileInstance | null>(null)
    const [sent, setSent] = React.useState<boolean>(false);

    async function submit(e: FormEvent<HTMLFormElement>) {
        e.preventDefault();

        if (loading || sent) return;

        setLoading(true)
        try {
            const cmail = email;
            await toast.promise(
                resetPassword.mutateAsync({
                    email: cmail,
                    turnstile: turnstileToken,
                }),
                {
                    success: t("auth.reset_password.sentEmail"),
                    loading: t("auth.reset_password.loading"),
                    error: (err: AppError) => buildError(err),
                }
            )

            setSent(true)
        } finally {
            setLoading(false)
        }
    }

    return (
        <form className={cn("flex flex-col gap-6", className)} {...props} onSubmit={submit}>
            <FieldGroup>
                <div className="flex flex-col items-center gap-1 text-center">
                    <h1 className="text-2xl font-bold">{t("auth.reset_password.title")}</h1>
                    <p className="text-muted-foreground text-sm text-balance">
                        {t("auth.reset_password.description")}
                    </p>
                </div>
                <Field>
                    <FieldLabel htmlFor="email">{t("auth.email")}</FieldLabel>
                    <Input value={email} onChange={(e) => setEmail(e.target.value)} id="email" type="email" placeholder={t("auth.emailPlaceholder")} required />
                </Field>
                <Field>
                    <Button type="submit" disabled={loading}>
                        {loading ? <Spinner /> : t("auth.reset_password.buttonLabel")}
                    </Button>
                </Field>
                <FieldSeparator>{t("common.or")}</FieldSeparator>
                <Field>
                    <ExternalAuth />
                    <FieldDescription className="text-center">
                        <Trans
                            i18nKey="auth.reset_password.footerText"
                            components={{
                                l: <Link to="/auth/register" />
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
