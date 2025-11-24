import type LoginStart from "../../models/auth/LoginStart";
import type Session from "@/lib/api/models/auth/Session";
import Request from "../Request";

export default async function loginStart(data: LoginStart): Promise<Session> {
    return await Request<Session>({
        method: "POST",
        url: "/auth/login/start",
        data,
    });
}
