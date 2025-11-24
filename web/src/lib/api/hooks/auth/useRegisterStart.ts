import { useMutation } from "@tanstack/react-query";
import type RegisterStart from "../../models/auth/RegisterStart";
import registerStart from "../../client/auth/registerStart";

export default function useRegisterStart() {
    return useMutation({
        mutationFn: (data: RegisterStart) => registerStart(data)
    })
}
