
import Sections, { SectionsType } from "./sidebarSections";
import { User } from "../header";
import { get, post } from "~/api";
import { cookies, headers } from "next/headers";
import Link from "next/link";


export default async function  Sidebar({ id }: { id: string }) {
    // const [selected, setSelected] = useState(0);
    
    const user: User = await post('getLoggedUserDetails',{}, {
        Cookie:cookies().toString(),

    })
    console.log(user, "userasdasd")

    const sections: SectionsType[] = [
        {
            name: "Home",
            icon: "",
            path: "/"
        },
        {
            name: "Profile",
            icon: "",
            path: "/channel/" +user.username
    
        },
        {
            name: "Studio",
            icon: "",
            path: "/studio"
        },
        {
            name: "LogOut",
            icon: "",
            path: "/auth/logout"
        }
    ]
    
    

    return (
        <div className="w-[20%] flex flex-col  h-full bg-primaryGrad1  " >

            <Sections pathIndex={1} sections={sections}  />
          

        </div>
    )
}