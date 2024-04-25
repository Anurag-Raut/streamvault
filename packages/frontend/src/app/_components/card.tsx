
"use client"
import { CSSProperties, ReactNode } from "react"



export default function Card({children,classname}:{
    children: ReactNode,
    classname?:String

}){
    return (
        <div  className={`bg-card p-6 rounded-lg mb-5 ${classname} items-center `}>
            {children}
        </div>
    )
}