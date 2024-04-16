
import Sections, { SectionsType } from "./sidebarSections";
import { User } from "../header";
import { get } from "~/api";
import { headers } from "next/headers";


export default async function  Sidebar({ id }: { id: string }) {
    // const [selected, setSelected] = useState(0);
    
    const user: User = await get('getLoggedUserDetails', {},headers() )
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
    
    

    return (
        <div className="w-[20%] flex flex-col  h-full bg-primaryGrad1  " >

            <Sections pathIndex={1} sections={sections}  />

        </div>
    )
}