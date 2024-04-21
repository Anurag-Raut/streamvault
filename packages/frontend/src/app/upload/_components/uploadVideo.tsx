import axios from "axios";
import Image from "next/image";
import { useState } from "react";
import ProgressBar from "@ramonak/react-progress-bar";



export default function UploadVideo({ setVideo }: { setVideo: any }) {
    const [currentVideo, setCurrentVideo] = useState<any>(null)
    const [uploading, setUploading] = useState(false)
    const [progress, setProgress] = useState(0)
    const [uploaded, setUploaded] = useState(false)
    const onUpload = async (file: any) => {
        setCurrentVideo(file)
        setUploading(true)
        const formData = new FormData()
        formData.append("file", file)
        const res = await axios.post("http://localhost:8080/uploadVideo", formData, {
            onUploadProgress: (progressEvent) => {
                console.log(progressEvent, "progress")
                setProgress(Math.round((progressEvent?.progress ?? 0) * 100))
            },
            withCredentials: true
        })
        setVideo(res.data)
        setUploading(false)
        setUploaded(true)
        console.log(res.data, "res")



    }




    return (
        <div className=" w-full min-w-[500px]">
            <div className="mb-4  ">Upload Video</div>
            <div className="w-full flex flex-col  items-center">

                <div className="relative p-4 py-6 m-3 flex flex-col justify-center items-center  outline-dashed rounded-xl outline-gray-500 min-h-[30vh] max-w-md w-[28vw] hover:bg-background3 min-w-[400px]  ">
                    <input onChange={(e) => {
                        onUpload(e?.target?.files?.[0])
                    }} type="file" className="absolute w-full h-full bg-transparent opacity-0  " />
                    {
                        uploaded ?
                            <video controls className="w-[200px] h-[100px]" >
                                <source src={URL.createObjectURL(currentVideo)} type="video/mp4" />
                            </video>
                            :
                            uploading ?
                                <div className="w-[350px] flex flex-col items-center  ">
                                    <div className="text-lg mb-3   ">Uploading...</div>
                                    <ProgressBar className="w-full bg-primary" isLabelVisible={false} completed={progress} />

                                </div>

                                :
                                <div className=" flex flex-col w-full  justify-center items-center">

                                    <Image
                                        src={"/thumbnailUpload.png"}
                                        width={100}
                                        height={100}
                                        alt="thumbnail"
                                    />
                                    <div className="text-sm font-bold opacity-60 mt-6">
                                        Click to select a video
                                    </div>
                                    <div className="text-sm  opacity-60 my-2">
                                        OR
                                    </div>
                                    <div className="text-sm font-bold opacity-60 mt-0">
                                        Drag and Drop your video here
                                    </div>
                                </div>
                    }


                </div>
            </div>
        </div>
    )
}