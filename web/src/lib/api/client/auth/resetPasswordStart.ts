import type ResetPasswordStart from "../../models/auth/ResetPasswordStart";
import Request from "../Request";

export default async function resetPasswordStart(data: ResetPasswordStart): Promise<void> {
    return await Request<void>({
        method: "POST",
        url: "/auth/reset-password/start",
        data,
    })
}
