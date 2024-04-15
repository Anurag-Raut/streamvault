"use client"
import { usePathname, useRouter } from "next/navigation";
import { useState } from "react"
import Sections, { SectionsType } from "./sidebarSections";

const sections: SectionsType[] = [
    {
        name: "Home",
        icon: "",
        path: "/"
    },
    {
        name: "Profile",
        icon: "",
        path: "/profile"

    },
    {
        name: "Settings",
        icon: "",
        path: "/settings"
    },
    {
        name: "LogOut",
        icon: "",
        path: "/logout"
    }
]


export default function Sidebar({ id }: { id: string }) {
    // const [selected, setSelected] = useState(0);



    return (
        <div className="w-[20%] flex flex-col  h-full bg-primaryGrad1  " >

            <Sections pathIndex={1} sections={sections}  />

        </div>
    )
}