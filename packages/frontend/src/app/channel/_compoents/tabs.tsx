"use client"
import { motion } from "framer-motion"
import Link from "next/link";
import { usePathname } from "next/navigation";
export default function Tabs({ username }: {
    username: string
}) {
    const pathname = usePathname();
    const tabs = [
        {
            pathname: `/channel/${username}`,
            name: "Live"

        },
        {
            pathname: `/channel/${username}/vod`,
            name: "Videos"
        },

    ]
    const currentPath = pathname.split('/')[3] ?? ""
    const previousPath = '/' + pathname.split('/').slice(1).join('/')
    console.log(previousPath, "previousPath")

    console.log(currentPath, "currentPath")
    return (

        <div className="flex my-5 flex-col w-full ">
            <div className="flex flex-row w-fit">

                {
                    tabs.map((tab, index) => (
                        <Link href={tab.pathname} className={`mr-5 h-full w-fit flex flex-col justify-end min-w-[100px] items-center ${(previousPath) !== (tab.pathname) && "hover:border-b-4 hover:border-gray-400" }  `}>
                            <div className={`py-3 text-xl font-extrabold ${(previousPath) !== (tab.pathname) && "opacity-50"}`}>

                                {tab.name}
                            </div>
                            {(previousPath) === (tab.pathname) && (<motion.div layoutId="underline" className="w-full h-1 rounded-xl  bg-purple ">

                            </motion.div>)}


                        </Link>
                    ))
                }
            </div>
            <div className="w-full bg-white opacity-50 h-[0.5px] rounded-xl" />



        </div>

    )
}