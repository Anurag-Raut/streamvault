"use client"

import { useRouter } from "next/navigation"
import { useEffect } from "react"
import { get } from "~/api"



export default function LogOut(){
    const router = useRouter()

    useEffect(()=>{
        async function logout(){
            const data = await get("signOut",{})
            console.log(data)
            router.replace("http://localhost:3000/")
            
        }
        logout()
    },[])
    return (
        <div>
            loggin out out
        </div>
    )
}