"use client"
import { useEffect, useState } from "react"
import { useRecoilState } from "recoil"
import { streamInfo } from "~/recoil/atom/streamInfo"

export default function InfoComponent() {
    const [streamDetails, setStreamDetails] = useRecoilState(streamInfo)
    const [isClient, setIsClient] = useState(false)

    useEffect(() => {
        setIsClient(true)
    }, [])

   

    return (
        <div className="p-6 rounded-lg flex flex-row justify-between w-full">
            <div className="flex flex-col items-between">
                <div className="mb-3">
                    <h1 className="text-[#fff] text-xs text-opacity-60">Title</h1>
                 { isClient&&  <h1 className="text-[#fff] text-xl pl-1">{streamDetails.title}</h1>}
                </div>

                <div className="mb-2">
                    <h1 className="text-[#fff] text-xs text-opacity-60">Description</h1>
                    { isClient&&   <h1 className="text-[#fff] text-xl pl-1">{streamDetails.description}</h1>}
                </div>
                <div className="mb-2">
                    <h1 className="text-[#fff] text-xs text-opacity-60">Privacy</h1>
                    { isClient&&   <h1 className="text-[#fff] text-xl pl-1">{"Public"}</h1>}
                </div>
                <div className="mb-2">
                    <h1 className="text-[#fff] text-xs text-opacity-60">Category</h1>
                    { isClient&&    <h1 className="text-[#fff] text-xl pl-1">{streamDetails.category}</h1>} 
                </div>
            </div>
            <button onClick={() => (document?.getElementById('my_modal_3') as HTMLDialogElement)?.showModal()} className="btn btn-secondary">Edit</button>
        </div>
    )
}
