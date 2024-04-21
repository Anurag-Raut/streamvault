"use client"

import Image from "next/image";
import Card from "../_components/card";
import TextArea from "./_components/textarea";
import Uploadvideo from "./_components/uploadVideo";
import UploadThumbnail from "./_components/uploadThumbnail";
import { useState } from "react";
import { toast } from "react-toastify";
import UploadVideo from "./_components/uploadVideo";

export default function UploadPage() {
    enum Visibility {
        Public = 0,
        Private = 1
    }
    const [data, setData] = useState<{
        title: string,
        description: string,
        category: string,
        visibility: Visibility,
        thumbnail: string | null,
        videoId: string | null
    
    }>({
        title: '',
        description: '',
        category: '',
        visibility: Visibility.Public,
        thumbnail: null,
        videoId: null

    })

    const validate = () => {
        if (!data.title || !data.description || !data.category  || !data.thumbnail || !data.videoId) {
            return false
        }
        if (data.title.length < 5 || data.description.length < 5 || data.category.length < 5) {
            return false
        }

        return true
    }

    const upload = async() => {
        if (!validate()) {
            console.log(data)
            toast.error("Please fill all fields")
            return
        }

        console.log(data, 'data')
        await fetch('http://localhost:8080/saveVod', {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(data),
            credentials: 'include'
        }).then((res) => {
            if (res.ok) {
                toast.success("videoId uploaded successfully")
            }
            else {
                toast.error("Error uploading videoId")
            }
        }).catch((e) => {
            toast.error("Error uploading videoId")
        })
    }
    console.log(data, 'data')

    return (
        <div className="w-full h-full p-6">
            <div className="mb-4 flex justify-between items-center ">
                <div className="text-xl ml-2 ">
                    Upload videoId
                </div>
                <button onClick={upload} className="btn btn-primary text-md ">Save</button>

            </div>
            <div className=" flex">

                <div className="w-full mr-3">
                    <Card>
                        <div className="mb-3">
                            <p className="text-md mb-2">Title <span className="text-sm opacity-60 ml-2">(required)</span></p>
                            <TextArea value={data.title} onChange={(val) => {
                                setData((prev) => (
                                    {
                                        ...prev,
                                        title: val
                                    }
                                ))

                            }}

                            />
                        </div>
                        <div className="mb-3 h-auto min-h-[200px]">
                            <p className="text-md mb-2">Description <span className="text-sm opacity-60 ml-2">(required)</span></p>
                            <TextArea onChange={(val) => {
                                setData((prev) => (
                                    {
                                        ...prev,
                                        description: val
                                    }
                                ))

                            }} height={200} classname="max-h-[600px] min-h-[200px]" />
                        </div>
                        <div className="mb-3">
                            <p className="text-md mb-2">Category <span className="text-sm opacity-60 ml-2">(required)</span></p>
                            <TextArea onChange={(val) => {
                                setData((prev) => (
                                    {
                                        ...prev,
                                        category: val
                                    }
                                ))

                            }} />
                        </div>
                        <div className="mb-3">
                            <p className="text-md mb-2">Visibility <span className="text-sm opacity-60 ml-2">(required)</span></p>
                            <div className="flex mb-3">
                                <input onChange={(e) => setData((prev) => (
                                    {
                                        ...prev,
                                        visibility: Visibility.Public
                                    }
                                ))} type="radio" name="radio-2" className="radio radio-primary" />
                                <div className="ml-2">Public</div>
                            </div>
                            <div className="flex mb-3">
                                <input onChange={(e) => setData((prev) => (
                                    {
                                        ...prev,
                                        visibility: Visibility.Private
                                    }
                                ))} type="radio" name="radio-2" className="radio radio-primary" />
                                <div className="ml-2">Private</div>
                            </div>
                            <div>
                            </div>


                        </div>
                    </Card>
                </div>
                <div className="w-full ml-3">
                    <Card>
                        <UploadVideo setVideo={(videoId: any) => setData((prev) => (
                            {
                                ...prev,
                                videoId: videoId
                            }
                        ))} />
                    </Card>
                    <Card>
                        <UploadThumbnail  setThumbnail={(val:string)=>{
                            setData((prev)=>({
                                ...prev,
                                thumbnail:val
                            }))
                        }} />
                    </Card>

                </div>

            </div>
        </div>
    );
}   