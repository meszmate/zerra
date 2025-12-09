import type { FormEvent } from "react";
import React from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import type { TurnstileInstance } from "@marsidev/react-turnstile";
import toast from "react-hot-toast";
import type { AppError } from "@/lib/api/client/normalizeError";
import buildError from "@/lib/helper/buildError";
import useRegisterConfirm from "@/lib/api/hooks/auth/useRegisterConfirm";
import CodeForm from "@/components/auth/code-form";
import { useTranslation } from "react-i18next";

export function RegisterConfirm() {
    const { t } = useTranslation();

    const [otp, setOtp] = React.useState<string>("");
    const [turnstileToken, setTurnstileToken] = React.useState<string>("");
    const [loading, setLoading] = React.useState<boolean>(false);
    const [searchParams] = useSearchParams();
    const navigate = useNavigate();
    const registerConfirm = useRegisterConfirm();

    const tRef = React.useRef<TurnstileInstance | null>(null)

    async function submit(e: FormEvent<HTMLFormElement>) {
        e.preventDefault();
        setLoading(true);

        try {
            const sessionToken = searchParams.get("session") ?? "";
            if (!sessionToken) {
                toast.error(t("auth.session.notfoundRedirect"))
                navigate("/auth/register")
            };

            await toast.promise(
                registerConfirm.mutateAsync({
                    code: otp,
                    turnstile: turnstileToken,
                    session: sessionToken,
                }),
                {
                    success: t("auth.register.success"),
                    loading: t("auth.register.loading"),
                    error: (err: AppError) => buildError(err),
                }
            );

            navigate("/auth/login")
        } finally {
            tRef.current?.reset();
            setLoading(false);
        }

    }

    return (
        <CodeForm
            onSubmit={submit}
            setToken={setTurnstileToken}
            loading={loading}
            otp={otp}
            setOtp={setOtp}
            tRef={tRef}
        />
    )
}
