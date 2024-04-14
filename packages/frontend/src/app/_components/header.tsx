"use server"

import { redirect } from "next/navigation"
import { get } from '~/api';
import { headers } from 'next/headers';
import Link from 'next/link';
import Avatar from "./avatar";
type User = {
    username: string,
    profileImage: string,
    userId: string,
    isLoggedIn: boolean
}
export default async function Header() {
    const user: User = await get('getUserDetails', {}, headers())
    console.log(user, "user")


    return (
        <div className="w-[100%] h-[10%] bg-primaryGrad1 fixed top-0 flex justify-between items-center p-10 " >
            <div>

            </div>
            {
                user.isLoggedIn ?
                   <Avatar name={user.username} src={user.profileImage} size={10} />
                    :
                    <Link href={"/auth/signIn"} className="btn btn-neutral">Sign in</Link>

            }

        </div>
    )

}