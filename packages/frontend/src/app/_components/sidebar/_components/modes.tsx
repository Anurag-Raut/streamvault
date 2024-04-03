"use client"
import { useState } from "react"


export default function Sections(){

    const [selected,setSelected]=useState(0);
    const sections=[
        {
            name:"Home",
            icon:""
        },
        {
            name:"Profile",
            icon:""
        },
        {
            name:"Settings",
            icon:""
        },
        {
            name:"Logout",
            icon:""
        }
    ]

    return (
        <div className="w-full">
            {
                sections.map((section,index)=>{
                    return (
                        <div onClick={()=>{
                            setSelected(index)
                        }} key={index} className={`w-[95%] cursor-pointer  ${selected===index? "border-2 border-purple p-3 text-purple":"text-gray-400" }  rounded-xl   text-sm mb-5 `}>
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
    )
}