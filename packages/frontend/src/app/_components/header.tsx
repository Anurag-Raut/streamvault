"use client"

import { redirect } from "next/navigation"
import { get, post } from '~/api';
import { cookies, headers } from 'next/headers';
import Link from 'next/link';
import Avatar from "./avatar";
import { RiVideoAddFill } from "react-icons/ri";
import HeaderDropDown from "./headerDropDown";
import { useEffect, useState } from "react";

export type User = {
    username: string,
    profileImage: string,
    userId: string,
    isLoggedIn: boolean
}
export default async function Header() {
    const user: User = await post('getLoggedUserDetails',{}, {
        Cookie:cookies().toString(),
    })
    console.log(user, "userasdasd")


    return (
        <div className="w-[100%] h-[10%] bg-primaryGrad1 fixed top-0 flex justify-between items-center p-10 " >
            <div>

            </div>
            {
                user.isLoggedIn ?
                    <div className=" flex items-center ">

                        {/* <div className="dropdown dropdown-bottom dropdown-end mx-5">
                            <div tabIndex={0} role="button" className=" m-1"><RiVideoAddFill className="fill-primary" size={28} /></div>
                            <ul tabIndex={0} className="dropdown-content z-[1] menu p-2 shadow bg-background3 rounded-box w-52">
                                <li> <Link className=" " href={'/dashboard'}>
                                        Start Stream
                                    </Link>
                                </li>
                                <li><a>Item 2</a></li>
                            </ul>
                        </div> */}
                        <HeaderDropDown/>

                        {/* </Link> */}
                        <Link className="mx-5" href={'/channel/' + user.username}>
                            <Avatar name={user.username} src={user.profileImage} size={38} />
                        </Link>

                    </div>
                    :
                    <Link href={"/auth/signIn"} className="btn btn-neutral">Sign in</Link>

            }

        </div>
    )

}