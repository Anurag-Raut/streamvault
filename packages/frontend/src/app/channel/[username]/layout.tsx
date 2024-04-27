import { headers } from "next/headers";
import { get, post } from "~/api";
import Avatar from "~/app/_components/avatar";
import Tabs from "../_components/tabs";
import { ReactNode } from "react";
import { User } from "~/app/_components/header";
import Link from "next/link";



export default async function Layout({ params, children }: { params: { username: string }, children: ReactNode }) {
    console.log(params, "params")
    const username = decodeURIComponent(params.username)

    const data: {
        username: string,
        profileImage: string
        userId: string
        subscribers: number
    } = await post('getUserDetailsByUsername', JSON.stringify(username), {}, new Headers(headers()))
    const loggedUserDetails: User = await get('getLoggedUserDetails', {}, new Headers(headers()))
    const isCurrentUser = loggedUserDetails.isLoggedIn && loggedUserDetails.username === data.username
    return (
        <div className=" h-full  w-full p-6">
            <div className="flex">
                <Avatar size={200} src={data.profileImage} name={data.username} />
                <div className=" ml-8 pt-5">
                    <div className="text-5xl  font-bold ">{data?.username}</div>
                    <div className="opacity-70">{data.subscribers} subscribers</div>
                    <div className="mt-5">
                        {

                            isCurrentUser && <Link href={`/studio/${loggedUserDetails.userId}/customize`} className=" btn p-3 px-5 bg-background4 hover:bg-background3 text-white p-2 rounded-full">Edit Profile</Link>
                            }
                    </div>
                </div>

            </div>
            <Tabs username={params.username} />
            {
                children
            }
        </div>
    )
}