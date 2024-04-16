"use client"
import { useRouter, useSearchParams } from "next/navigation"
import { useEffect, useState } from "react"
import { toast } from "react-toastify"
import { get, post } from "~/api"



export default function SignInPage() {
    const [username, setUsername] = useState('')
    const [password, setPassword] = useState('')
    const searchParams = useSearchParams()
    const router = useRouter()
    useEffect(() => {
        const code = searchParams.get('code')
        async function loginWithGoogle() {
            await post("loginWithGoogle", JSON.stringify(code))
            router.replace("http://localhost:3000/")
        }
        if (code) {
            loginWithGoogle()
        }
    }, [])
    async function handleSignIn() {
        try {
            // console.log("asd")

            // const res = await fetch('http://localhost:8080/signup', {
            //     method: 'POST',

            //     body: JSON.stringify({
            //         username: username,
            //         password: password
            //     }),
            //     credentials:'include'
            // })
            // const response = await res.text()

            const response = await post('signup', JSON.stringify({ username, password }))
            // toast.success(response)
            console.log(response, "ressssss")

        }
        catch (error) {
            console.log(error)
            toast.error('An error occurred')
        }


    }

    async function loginWithGoogle() {
        const url: string = await get('getGoogleUrl')
        console.log(url)
        router.replace(url)
    }

    return (
        <div className="w-full h-full flex justify-center items-center">
            <div className="w-fit h-fit p-5 bg-primary rounded-xl ">
                <h2>Username</h2>
                <input onChange={(e) => setUsername(e.target.value)} type="text" placeholder="Enter your username" className="input input-bordered w-full max-w-xs m-3" />
                <h2 className="m-2" >Password</h2>
                <input onChange={(e) => setPassword(e.target.value)} type="password" placeholder="Enter your password" className="input input-bordered w-full max-w-xs m-3" />
                <button onClick={handleSignIn} className="btn">Sign in</button>

                <button onClick={loginWithGoogle} >Log in with google</button>

            </div>

        </div>
    )
}
