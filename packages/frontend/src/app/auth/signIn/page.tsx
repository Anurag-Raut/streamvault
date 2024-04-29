import { Suspense } from "react";
import SusSignInPage from "./suspenseComp";


export default function SignInPage() {
    return (
        <Suspense>
        <SusSignInPage/>
        </Suspense>
    )
}