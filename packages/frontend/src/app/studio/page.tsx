"use client";
    import { useRouter } from "next/navigation"
import { useEffect } from "react"
import { get } from "~/api";


export default function Channel(){
    const router=useRouter();

        useEffect(()=>{
            async function  fetchChannel() {
                try{

                    
                    // const res=await fetch(`http://localhost:8080/getUserId`,{
                    //     method:"GET",
                    //     headers:{
                    //         "Content-Type":"application/json"
                    //     },
                    //     credentials:"include"
                    // },)
                    // const data=await res.text();

                    const data:{
                        userId:string
                    
                    }=await get('getUserId',{})
                    console.log(data.userId,"channel data")
                    router.replace(`studio/${data.userId}`)
                    
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