import { useState } from "react"
import { cookies } from 'next/headers'



export default function SignInPage() {
    const [username,setUsername]=useState('')
    const [password,setPassword]=useState('')
    async function handleSignIn(){
        const res=await fetch('/api/auth/signIn', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                username: 'username',
                password: 'password'
            })
        })
        const response=await res.text()
        cookies().set({
            name: 'token',
            value: response,
            httpOnly: true,
            path: 'http://localhost:',
            
          })
        

    }

    return(
        <div className="w-full h-full flex justify-center items-center">
            <div className="w-fit h-fit p-5 bg-primary rounded-xl ">
            <h2>Username</h2>
            <input onChange={(e)=>setUsername(e.target.value)} type="text" placeholder="Enter your username" className="input input-bordered w-full max-w-xs m-3" />
            <h2 className="m-2" >Password</h2>
            <input onChange={(e)=>setPassword(e.target.value)} type="password" placeholder="Enter your password" className="input input-bordered w-full max-w-xs m-3" />
            <button onClick={handleSignIn} className="btn">Sign in</button>


            </div>

        </div>
    )
}
