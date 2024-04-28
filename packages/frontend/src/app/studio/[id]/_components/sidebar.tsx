
import { headers } from "next/headers";
import { usePathname } from "next/navigation";
import { get, post } from "~/api";
import Avatar from "~/app/_components/avatar";
import { User } from "~/app/_components/header";
import Sections from "~/app/_components/sidebar/sidebarSections";


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
        name: "Customize",
        icon: "",
        path:"/customize"
    },
    {
        name: "Comments",
        icon: "",
        path:"/comments"
    }
]

export default async function Sidebar({id}:{id:string}) {
    // const [selected, setSelected] = useState(0);
    const user: User = await post('getLoggedUserDetails',{}, {}, new Headers(headers()))


    return (
        <div className="w-[20%] min-w-[200px] flex flex-col  h-full bg-primaryGrad1  " >
            
            <div className="my-5 mb-10 justify-center flex  p-6">
                <Avatar size={160} src={user.profileImage} name={user.username} />
            </div>
            <Sections pathIndex={3} sections={sections} />
        </div>
    )
}