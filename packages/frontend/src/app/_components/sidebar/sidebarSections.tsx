"use client"
import Link from "next/link";
import { redirect, usePathname } from "next/navigation";
import { useRouter } from "next/navigation";

export type SectionsType = {
    name: string,
    path: string,
    icon?: string  
}
export default function Sections({sections,pathIndex}:{
    sections:SectionsType[]
    pathIndex:number
    
}){

    const pathname=usePathname();
    console.log(pathname,"pathname" )
    const currentPath=pathname.split('/')[pathIndex]??""
    const previousPath=pathname.split('/').slice(0,pathIndex).join('/')
    console.log(previousPath,"previousPath")

    


   
    return (
        <div className="w-full flex flex-col">
        {
            sections.map((section, index) => {
                return (
                    <Link
                    
                        href={previousPath+sections[index]?.path??""}
                        
                    
                     key={index} className={`w-[100%] cursor-pointer  ${"/"+currentPath === sections[index]?.path ? " border-l-4 w-full  border-l-primary p-3 text-primary font-bold bg-purple3 tracking-wide text-lg rounded-l-md" : "text-gray-400  "} hover:bg-background3    p-3 text-sm  `}>
                        {section.name}
                    </Link>
                )
            })
        }

    </div>
    )
}