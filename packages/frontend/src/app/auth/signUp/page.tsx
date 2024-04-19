"use client"
import {motion} from "framer-motion"
import { useRouter } from "next/navigation"
import { useState } from "react"
import { FcGoogle } from "react-icons/fc"
import { toast } from "react-toastify"
import { get, post } from "~/api"



export default function SignUpPage() {
    const [username, setUsername] = useState('')
    const [password, setPassword] = useState('')
    const router=useRouter()

    const handleSignUp = async () => {
        try {
          

            const response = await post('signup', JSON.stringify({ username, password }))
            console.log(response, "ressssss")
            toast.success("Signed up Successfully")
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

    return(
        <div className="w-full h-full flex justify-center items-center">
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
                <button onClick={handleSignUp} className="btn w-full bg-primary text-white font-extrabold text-lg border-4 border-purple hover:border-4 hover:border-purple hover:bg-background2 ">Sign Up</button>
                <div onClick={()=>{
                    router.replace('/auth/signIn')
                }} className="flex items-center justify-center  hover:underline cursor-pointer text-blue-500 font-bold m-3 self-center">Already have an account ? Sign In</div>
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
