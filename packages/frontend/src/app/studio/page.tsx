"use client";
import { useRouter } from "next/navigation"
import { useEffect } from "react"


export default function Channel(){
    const router=useRouter();

        useEffect(()=>{
            async function  fetchChannel() {
                try{

                    
                    const res=await fetch(`http://localhost:8080/getUserId`,{
                        method:"GET",
                        headers:{
                            "Content-Type":"application/json"
                        },
                        credentials:"include"
                    },)
                    const data=await res.text();
                    console.log(data,"channel data")
                    router.replace(`studio/${data}`)
                    
                }
                catch(error){
                    console.log(error)
                
                }
            }
            fetchChannel()
            
        },[])
    
        return (
        <div className="w-full h-full">
            Hell
    
        </div>
        )
    
    }