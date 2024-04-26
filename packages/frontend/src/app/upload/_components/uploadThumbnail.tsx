"use client";
import axios from "axios";
import Image from "next/image";
import { useState } from "react";
import { MdOutlineVideoCameraBack } from "react-icons/md";


export default function UploadThumbnail({ setThumbnail }: {
    
    setThumbnail: any;
}) {

    const [file, setFile] = useState<any>(null)

    const onUpload = async (file: any) => {
        setFile(file)
        const formData = new FormData()
        formData.append("thumbnail", file)

        const res=await axios.post("${process.env.NEXT_PUBLIC_BACKEND_URL}/uploadThumbnail",formData,{
            onUploadProgress:(progressEvent)=>{
                console.log(progressEvent, "progress")
            },
            withCredentials:true
        })
        console.log(res.data, "res")
        if(res.data.thumbnailPath){

            setThumbnail("${process.env.NEXT_PUBLIC_BACKEND_URL}/hls/"+res.data.thumbnailPath)
        }



    }

    return (
        <div>
            <div className="mb-4 ">Upload thumbnail</div>
            <div className="flex flex items-center    ">
                {
                    file ?

                        <Image
                            src={URL?.createObjectURL(file)}
                            width={200}
                            height={120}
                            alt="thumbnail"
                            className="object-scale-down w-[200px] h-[120px]"
                        />

                        :
                        <div className="w-[200px] h-[120px] items-center bg-gray-800 p-3 justify-center flex mx-3 mr-8 rounded-lg ">
                            <MdOutlineVideoCameraBack size={49} className="fill-gray-400" />

                        </div>

                }



                <input onChange={(e) => onUpload(e?.target?.files?.[0])} type="file" className="file-input file-input-bordered file-input-primary w-full max-w-xs" />

            </div>
        </div>
    )
}