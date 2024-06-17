"use client"
import { ChangeEvent, useState } from "react";
import { constSelector, useRecoilState } from "recoil";
import { streamInfo } from "~/recoil/atom/streamInfo";
import { toast } from 'react-toastify';
import axios from "axios";
import { post } from "~/api";

export default function Modal() {
    const [streamInfoState, setStreamInfoState] = useRecoilState(streamInfo)

    const [newStreamInfo, setNewStreamInfo] = useState(streamInfoState)
    const [files, setFiles] = useState<FileList | null>(null)

    const uploadThumbnail = async () => {
        try {
            const formdata = new FormData()
            formdata.append('thumbnail', files?.[0] as Blob)

            // const res=await axios.post('${process.env.NEXT_PUBLIC_BACKEND_URL}/uploadThumbnail',formdata,{withCredentials:true})
            const res:{
                thumbnailPath:string
            } = await post('uploadThumbnail', formdata, {
                // "Content-Type":"multipart/form-data"

            })
            console.log(res.thumbnailPath, "res")
            const thumbnail = "${process.env.NEXT_PUBLIC_BACKEND_URL}/hls/" + res.thumbnailPath;
            console.log(thumbnail, "thumbnail")
            setNewStreamInfo((prev:any) => ({ ...prev, thumbnail: thumbnail }))
            toast.success('Thumbnail uploaded successfully')
        }
        catch (err) {
            console.log(err)
            toast.error('err')
        }



    }

    const onStreamInfoEdit = async () => {



        setStreamInfoState(newStreamInfo)
        const modal = document.getElementById('my_modal_3') as HTMLDialogElement | null;
        if (modal) {
            modal.close();
        }
    }

    return (
        <dialog id="my_modal_3" className="modal w-full">
            <div className="modal-box w-full min-w-[50vw]">
                <form method="dialog">
                    {/* if there is a button in form, it will close the modal */}
                    <button className="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">✕</button>
                </form>
                <div className="w-[100%] ">
                    <div className="flex items-center my-5 justify-between  ">
                        <h1 className="text-xl text-center mr-3 text-nowrap">Title :</h1>
                        <input onChange={(event: ChangeEvent<HTMLInputElement>) => {
                            setNewStreamInfo((prev:any) => ({ ...prev, title: event.target.value }))
                        }} type="text" placeholder="Type here" className="input input-bordered w-full max-w-[75%] " />

                    </div>
                    <div className="flex items-center my-5 justify-between ">
                        <h1 className="text-xl text-center mr-3 text-nowrap ">Description :</h1>
                        <textarea onChange={(event: ChangeEvent<HTMLTextAreaElement>) => {
                            setNewStreamInfo((prev:any) => ({ ...prev, description: event.target.value }))
                        }} className="textarea textarea-bordered w-full  max-w-[75%]" placeholder="Bio"></textarea>
                    </div>
                    <div className="flex items-center my-5 justify-between ">
                        <h1 className="text-xl text-center mr-3 text-nowrap">Category :</h1>
                        <input onChange={(event: ChangeEvent<HTMLInputElement>) => {
                            setNewStreamInfo((prev:any) => ({ ...prev, category: event.target.value }))
                        }} type="text" placeholder="Type here" className="input input-bordered w-full max-w-[75%]" />

                    </div>
                    <div className="flex items-center my-5 justify-between ">
                        <h1 className="text-xl text-center mr-3 text-nowrap">Thumbnail :</h1>
                        <div className=" max-w-[75%] w-full ">

                        <input onChange={(event: ChangeEvent<HTMLInputElement>) => {
                            setFiles(event.target.files)
                            // setNewStreamInfo((prev)=>({...prev,thumbnail:event.target.files?.[0]}))
                        }} type="file" className="file-input file-input-bordered file-input-accent w-full max-w-[75%]" />
                        <button onClick={uploadThumbnail} className="btn btn-primary m-3">upload</button>
                        </div>

                    </div>
                    <button onClick={onStreamInfoEdit} className="btn btn-primary">Save</button>

                </div>
            </div>
        </dialog>
    )
}