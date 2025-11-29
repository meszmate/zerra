import Logo from "@/components/Logo"
import { NAME } from "@/lib/information"
import { Outlet } from "react-router-dom"

export default function AuthLayout() {
    return (
        <div className="grid min-h-svh lg:grid-cols-2">
            <div className="flex flex-col gap-4 p-6 md:p-10">
                <div className="flex justify-center gap-2 md:justify-start">
                    <a href="#" className="flex items-center gap-2 font-medium">
                        <Logo className="w-5.5" />
                        {NAME}
                    </a>
                </div>
                <div className="flex flex-1 items-center justify-center">
                    <div className="w-full max-w-md">
                        <Outlet />
                    </div>
                </div>
            </div>
            <div className="bg-muted relative hidden lg:block">

            </div>
        </div>
    )
}
