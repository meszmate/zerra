import { useMutation } from "@tanstack/react-query";
import type ResetPasswordStart from "../../models/auth/ResetPasswordStart";
import resetPasswordStart from "../../client/auth/resetPasswordStart";

export default function useResetPasswordStart() {
    return useMutation({
        mutationFn: (data: ResetPasswordStart) => resetPasswordStart(data)
    })
}
