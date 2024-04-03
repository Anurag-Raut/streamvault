
"use client"

import { useRouter } from "next/navigation"

export default function Header(){
    const router = useRouter()

    function handleSignIn(){

        router.push('/auth/signIn')

    }

    return (
        <div className="w-[100%] h-[10%] bg-gradient-to-l from-primaryGrad1 from-30% via-primaryGrad2 to-primaryGrad1 fixed top-0 flex justify-between items-center p-10 border-b-2 border-border" >
                <div>

                </div>
                <button onClick={handleSignIn} className="btn btn-neutral">Sign in</button>

        </div>
    )

}