"use client"
import { useRecoilState } from "recoil"
import { streamInfo } from "~/recoil/atom/streamInfo"



export default function InfoComponent() {
    const [streamDetails, setStreamDetails] = useRecoilState(streamInfo)
    return (
        <div className="bg-[#2B2A4C] p-6 rounded-lg flex flex-row justify-between w-full">
            <div>

                <h1 className="text-[#fff] text-xl">Title :  {streamDetails.title}</h1>
                <p className="text-[#fff] text-xl"> Description : {streamDetails.description}</p>
                <p className="text-[#fff] text-xl"> Privacy : {streamDetails.description}</p>
                <p className="text-[#fff] text-xl"> Category : {streamDetails.category}</p>
                <p className="text-[#fff] text-xl"> Viewers : {streamDetails.viewers}</p>

            </div>
            <button onClick={() => (document?.getElementById('my_modal_3') as HTMLDialogElement)?.showModal()} className="btn btn-secondary">Edit</button>
        </div>
    )
}