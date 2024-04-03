"use client"
import Sections from "./_components/modes";



export default function Sidebar(){
    return (
        <div className="w-[20%] flex  h-full bg-[radial-gradient(circle_at_top_right,_var(--tw-gradient-stops))] from-primaryGrad2   to-primaryGrad1 border-r-2 border-border p-6 " >
          
            <Sections   />
         
        </div>
    )
}