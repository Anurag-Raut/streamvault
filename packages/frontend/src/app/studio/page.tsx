"use client";
import { useRouter } from "next/navigation"
import { useEffect } from "react"
import { toast } from "react-toastify";
import { get } from "~/api";

import Loader from "../_components/loading";

// Manually defined



export default function Channel() {
    const router = useRouter();

    useEffect(() => {
        async function fetchChannel() {
            try {


                // const res=await fetch(`${process.env.NEXT_PUBLIC_BACKEND_URL}/getUserId`,{
                //     method:"GET",
                //     headers:{
                //         "Content-Type":"application/json"
                //     },
                //     credentials:"include"
                // },)
                // const data=await res.text();

                const data: {
                    userId: string

                } = await get('getUserId', {})
                if (!data.userId) {
                    toast.warning("Not logged in")
                    router.replace("/auth/signIn")
                    return
                }
                console.log(data.userId, "channel data")
                router.replace(`studio/${data.userId}`)

            }
            catch (error) {
                console.log(error)

            }
        }
        fetchChannel()

    }, [])

    return (
        <div className="w-full h-full h-[calc(100vh-81px)] min-h-[calc(100vh-81px)] justify-center items-center flex ">

            <Loader />

        </div>
    )

}