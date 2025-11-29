import type { FormEvent } from "react";
import React from "react";
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import {
    Field,
    FieldDescription,
    FieldGroup,
    FieldLabel,
} from "@/components/ui/field"
import {
    InputOTP,
    InputOTPGroup,
    InputOTPSlot,
} from "@/components/ui/input-otp"
import { useNavigate, useSearchParams } from "react-router-dom";

export default function LoginConfirm({
    className,
    ...props
}: React.ComponentProps<"form">) {
    const [otp, setOtp] = React.useState<string>("");
    const [turnstileToken, setTurnstileToken] = React.useState<string>("");
    const [searchParams] = useSearchParams();
    const navigate = useNavigate();

    async function submit(e: FormEvent<HTMLFormElement>) {
        e.preventDefault();

    }

    React.useEffect(() => {
        if (!searchParams.get("to")) {
            navigate("/auth/login")
        }
    }, [searchParams, navigate])

    return (
        <form onSubmit={submit} className={cn("flex flex-col gap-6", className)} {...props}>
            <FieldGroup>
                <Field>

                    <div className="flex flex-col items-center gap-1 text-center">
                        <h1 className="text-2xl font-bold">Email Cím Hitelesítése</h1>
                        <p className="text-muted-foreground text-sm text-balance">
                            A kód el lett küldve a <span>{searchParams.get("to") ? searchParams.get("to") : "ismeretlen"} címre</span>
                        </p>
                    </div>
                    <FieldLabel htmlFor="otp">Verification code</FieldLabel>
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
                    <FieldDescription>
                        Enter the 6-digit code sent to your email.
                    </FieldDescription>
                </Field>
                <FieldGroup>
                    <Button type="submit">Verify</Button>
                    <FieldDescription className="text-center">
                        Didn&apos;t receive the code? <a href="#">Resend</a>
                    </FieldDescription>
                </FieldGroup>
            </FieldGroup>

        </form>
    )
}
