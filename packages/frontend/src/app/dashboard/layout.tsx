"use client";

import {  ReactNode } from "react";
import { RecoilRoot } from "recoil";


export default function Layout({ children }: {
    children: ReactNode

}) {


    return (
        <RecoilRoot>

            <div className="w-[100%] h-full">
                {children}
            </div>
         </RecoilRoot>

    )

}
