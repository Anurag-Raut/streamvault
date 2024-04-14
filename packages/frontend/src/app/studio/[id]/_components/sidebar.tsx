"use client"
import { usePathname, useRouter } from "next/navigation";
import { useState } from "react"


export default function Sidebar({id}:{id:string}) {
    const router = useRouter();
    const [selected, setSelected] = useState(0);
    const sections = [
        {
            name: "Home",
            icon: "",
            path:"/"
        },
        {
            name: "Content",
            icon: "",
            path:"/content"

        },
        {
            name: "Analytics",
            icon: "",
            path:"/analytics"
        },
        {
            name: "Comments",
            icon: "",
            path:"/comments"
        }
    ]

    return (
        <div className="w-[20%] flex flex-col  h-full bg-primaryGrad1   p-6 " >
            <div className="my-5 mb-10">
                <img src="/logo.png" className="w-10 h-10" />   
            </div>
            <div className="w-full">
                {
                    sections.map((section, index) => {
                        return (
                            <div onClick={() => {
                                // console.log(path+"/content")
                                router.replace(`/studio/${id}${section.path}`)
                                setSelected(index)
                            }} key={index} className={`w-[95%] cursor-pointer  ${selected === index ? "border-2 border-purple p-3 text-purple" : "text-gray-400"}  rounded-xl   text-sm mb-5 `}>
                                {section.name}
                            </div>
                        )
                    })
                }
                {/* // <div className="w-[95%] border-2 border-purple p-3 rounded-xl text-purple  text-sm mb-5 ">
            //     Home 
            // </div>
            // <div className="w-[95%] p-3 text-white  text-sm ">
            //     Home 
            // </div> */}
            </div>
        </div>
    )
}