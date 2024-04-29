"use client"
import { useRouter, useSearchParams } from "next/navigation"
import { useEffect, useState } from "react"
import { toast } from "react-toastify"
import { get, post } from "~/api"
import { FcGoogle } from "react-icons/fc";
import { motion } from "framer-motion"


export default function SusSignInPage() {
    const [username, setUsername] = useState('')
    const [password, setPassword] = useState('')
    const searchParams = useSearchParams()
    const router = useRouter()
    useEffect(() => {
        const code = searchParams.get('code')
        async function loginWithGoogle() {
            await post("loginWithGoogle", JSON.stringify(code))
            router.replace("/")
            router.refresh()

        }
        if (code) {
            loginWithGoogle()
        }
    }, [])
    async function handleSignIn() {
        try {
       

            const response = await post('signin', JSON.stringify({ username, password }))
            toast.success("Signed in Successfully")
            router.replace("/")
            router.refresh()

        }
        catch (error) {
            console.log(error)
            toast.error(error?.toString()??"an error occured")
        }


    }

    async function loginWithGoogle() {
        const url: string = await get('getGoogleUrl')
        console.log(url)
        router.replace(url)
    }

    return (
        <div className="w-full h-full h-[calc(100vh-81px)] min-h-[calc(100vh-81px)] flex justify-center items-center">
            <div className=" w-[500px] h-fit p-5 bg-background3 rounded-xl ">
                <div className="my-3 mb-5 ">
                    <p className="text-2xl font-bold ">    Welcome to Streamvault</p>
                    <p className="opacity-70">Welcome to streamvault and start exploring</p>

                </div>
                <div className="mb-6 ">
                    <h2 className="mb-1 pl-1">Username</h2>
                    <input onChange={(e) => setUsername(e.target.value)} type="text" placeholder="Enter your username" className="input input-bordered w-full " />
                </div>
                <div className="mb-6">

                    <h2 className="mb-1 pl-1" >Password</h2>
                    <input onChange={(e) => setPassword(e.target.value)} type="password" placeholder="Enter your password" className="input input-bordered w-full  " />
                </div>
                <button onClick={handleSignIn} className="btn w-full bg-primary text-white font-extrabold text-lg border-4 border-purple hover:border-4 hover:border-purple hover:bg-background2 ">Sign in</button>
                <div onClick={()=>{
                    router.replace('/auth/signUp')
                }} className="flex items-center justify-center  hover:underline cursor-pointer text-blue-500 font-bold m-3 self-center">Dont have an Account ? Sign Up</div>
                <div className="divider divider-secondary">OR</div>

                <motion.button whileHover={{
                    scale: 1.05

                }}
                    whileTap={{
                        scale: 0.95

                    }}


                    className="w-full h-12 bg-white rounded-xl justify-center items-center flex hover:opacity-60" onClick={loginWithGoogle} >
                    <FcGoogle size={23} />
                    <span className="ml-3 text-background font-bold text-lg">Sign in with Google</span>
                </motion.button>

            </div>

        </div>
    )
}
