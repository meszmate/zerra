import React, { type FormEvent } from "react";
import { Link, useSearchParams } from "react-router-dom";
import { Turnstile, type TurnstileInstance } from "@marsidev/react-turnstile"
import { cn } from "@/lib/utils";
import { Field, FieldDescription, FieldGroup, FieldLabel } from "../ui/field";
import { InputOTP, InputOTPGroup, InputOTPSlot } from "../ui/input-otp";
import { Button } from "../ui/button";
import { Spinner } from "../ui/spinner";
import { TURNSTILE_KEY } from "@/lib/information";
import { Trans, useTranslation } from "react-i18next";

interface CodeFormProps extends React.ComponentProps<"form"> {
    onSubmit: (e: FormEvent<HTMLFormElement>) => Promise<void>;
    setToken: React.Dispatch<string>,
    loading: boolean,
    otp: string,
    setOtp: React.Dispatch<string>,
    tRef: React.RefObject<TurnstileInstance | null>,
}

export default function CodeForm({
    className,
    onSubmit,
    setToken,
    loading,
    otp,
    setOtp,
    tRef,
    ...props
}: CodeFormProps) {
    const { t } = useTranslation();
    const [searchParams] = useSearchParams();

    return (
        <form onSubmit={onSubmit} className={cn("flex flex-col gap-6", className)} {...props}>
            <FieldGroup>
                <Field>
                    <div className="flex flex-col items-center gap-1 text-center">
                        <h1 className="text-2xl font-bold">Email Cím Hitelesítése</h1>
                        <p className="text-muted-foreground text-sm text-balance">
                            <Trans
                                i18nKey="auth.code.sentCode"
                                values={{
                                    val: searchParams.get("to") ? searchParams.get("to") : "ismeretlen",
                                }}
                                components={{
                                    email: <span className="font-bold" />
                                }}
                            />
                        </p>
                    </div>
                    <FieldLabel htmlFor="otp"></FieldLabel>
                    <Field className="flex justify-center">
                        <InputOTP maxLength={6} id="otp" required value={otp} onChange={(e) => setOtp(e)}>
                            <InputOTPGroup className="gap-2.5 *:data-[slot=input-otp-slot]:rounded-md *:data-[slot=input-otp-slot]:border">
                                <InputOTPSlot index={0} />
                                <InputOTPSlot index={1} />
                                <InputOTPSlot index={2} />
                                <InputOTPSlot index={3} />
                                <InputOTPSlot index={4} />
                                <InputOTPSlot index={5} />
                            </InputOTPGroup>
                        </InputOTP>
                    </Field>
                    <FieldDescription>
                        {t("auth.code.footerDescription")}
                    </FieldDescription>
                </Field>
                <FieldGroup>
                    <Button type="submit" disabled={loading}>
                        {loading ? <Spinner /> : "Hitelesítés"}
                    </Button>
                    <FieldDescription className="text-center">
                        <Trans
                            i18nKey="auth.code.noCode"
                            components={{
                                l: <Link to="/auth/login"></Link>
                            }}
                        />
                    </FieldDescription>
                </FieldGroup>
                <Field>
                    <Turnstile
                        ref={tRef}
                        siteKey={TURNSTILE_KEY}
                        onSuccess={(token) => setToken(token)}
                        onExpire={() => {
                            setToken("");
                            tRef.current?.reset();
                        }}
                        options={{
                            theme: "light",
                            size: "invisible",
                            language: "hu",
                        }}
                    />
                </Field>
            </FieldGroup>
        </form >
    )
}
