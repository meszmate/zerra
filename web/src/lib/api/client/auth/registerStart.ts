import type RegisterStart from "@/lib/api/models/auth/RegisterStart";
import type Session from "@/lib/api/models/auth/Session";
import Request from "../Request";

export default async function registerStart(data: RegisterStart): Promise<Session> {
    return await Request<Session>({
        method: "POST",
        url: "/auth/register/start",
        data,
    })
}
