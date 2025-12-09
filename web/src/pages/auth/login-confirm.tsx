import type { FormEvent } from "react";
import React from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import type { TurnstileInstance } from "@marsidev/react-turnstile";
import toast from "react-hot-toast";
import type { AppError } from "@/lib/api/client/normalizeError";
import buildError from "@/lib/helper/buildError";
import CodeForm from "@/components/auth/code-form";
import useLoginConfirm from "@/lib/api/hooks/auth/useLoginConfirm";
import setToken from "@/lib/helper/setToken";
import { useTranslation } from "react-i18next";

export function LoginConfirm() {
    const { t } = useTranslation();

    const [otp, setOtp] = React.useState<string>("");
    const [turnstileToken, setTurnstileToken] = React.useState<string>("");
    const [loading, setLoading] = React.useState<boolean>(false);
    const [searchParams] = useSearchParams();
    const navigate = useNavigate();
    const loginConfirm = useLoginConfirm();

    const tRef = React.useRef<TurnstileInstance | null>(null)

    async function submit(e: FormEvent<HTMLFormElement>) {
        e.preventDefault();
        setLoading(true);

        try {
            const sessionToken = searchParams.get("session") ?? "";
            if (!sessionToken) {
                toast.error(t("auth.login.session.notfoundRedirect"))
                navigate("/auth/login")
            };

            const resp = await toast.promise(
                loginConfirm.mutateAsync({
                    code: otp,
                    turnstile: turnstileToken,
                    session: sessionToken,
                }),
                {
                    success: t("auth.register.success"),
                    loading: t("auth.code.loading"),
                    error: (err: AppError) => buildError(err),
                }
            );

            setToken(resp);
            navigate("/app");
        } finally {
            tRef.current?.reset();
            setLoading(false);
        }

    }

    React.useEffect(() => {
        if (!searchParams.get("to")) {
            navigate("/auth/register")
        }
    }, [searchParams, navigate])

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
