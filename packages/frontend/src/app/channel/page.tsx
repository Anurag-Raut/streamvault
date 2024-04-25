"use client"
import { redirect } from "next/navigation"
import { useEffect } from "react"
import { toast } from "react-toastify"




export default function Channel(){
    console.log("elloooo")
    useEffect(()=>{
        toast.warning("Not logged in")
        redirect("/auth/signIn")
    },[])
   
       
        
 return (
    <div>

    </div>
 )

}