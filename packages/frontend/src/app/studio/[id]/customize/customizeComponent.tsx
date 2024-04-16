



"use client"
import Card from "~/app/_components/card"
import TextInput from "../_components/textInput"
import { useState } from "react"
import Image from "next/image"
import Avatar from "~/app/_components/avatar"
import { post } from "~/api"
import { toast } from "react-toastify"



export default function CustomizeComponent({ username, profileImage }: {
    username: string,
    profileImage: string
}) {

    const [newUsername, setNewUsername] = useState<string>(username)
    const [src, setSrc] = useState<string>(profileImage)
    const [files, setFiles] = useState<FileList | null>(null)

    function onFileChange(event: any) {
        setFiles(event.target.files)
        const file = event.target.files[0]
        const objectUrl = URL.createObjectURL(file)
        setSrc(objectUrl)

    }

    async function onSubmit() {
        //upload profile image and get the url
        try {

            const formdata = new FormData()
            formdata.append('profileImage', files?.[0] as Blob)

            // const res=await axios.post('http://localhost:8080/uploadThumbnail',formdata,{withCredentials:true})
            const res: {
                profileImagePath: string
            } = await post('uploadProfileImage', formdata, {
                // "Content-Type":"multipart/form-data"

            })
            console.log(res.profileImagePath, "res")
            const profileImageUrl = "http://localhost:8080/hls/" + res.profileImagePath;
            console.log(profileImageUrl, "thumbnail")

            //update username and profile image
            const data = {
                username: newUsername,
                profileImage: profileImageUrl
            }
            await post('updateUserDetails', JSON.stringify(data))
            toast.success('Profile updated successfully')

        }
        catch (err) {
            console.log(err)
            toast.error(err?.toString() ?? "")


        }

    }

    return (

        <Card>
            <div>
                <div className="mb-14">
                    <p className="text-lg">Username </p>
                    <TextInput value={newUsername} setValue={setNewUsername} type="text" placeholder={"Enter username"} />
                </div>
                <div >
                    <p className="mb-8 text-lg  ">Profile Image </p>
                    <div className=" flex items-center  ">
                        <Avatar name="Anurag" src={src} size={250} />
                        <input onChange={onFileChange} type="file" className="file-input file-input-bordered w-full max-w-xs ml-14" />

                    </div>
                </div>
                <div className="w-full justify-end flex flex-row">
                    <button onClick={onSubmit} className="btn bg-primary text-white font-extrabold hover:bg-primary hover:opacity-60 text-md">Save</button>
                </div>

            </div>
        </Card>
    )

}