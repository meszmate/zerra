import { useMutation } from "@tanstack/react-query";
import type LoginStart from "../../models/auth/LoginStart";
import loginStart from "../../client/auth/loginStart";

export default function useLoginStart() {
    return useMutation({
        mutationFn: (data: LoginStart) => loginStart(data)
    })
}
