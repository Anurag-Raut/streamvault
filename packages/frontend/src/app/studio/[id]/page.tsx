import { headers } from "next/headers";
import { get } from "~/api";
import Card from "~/app/_components/card";
import { Content } from "./content/page";
import { AiFillDislike, AiFillLike } from "react-icons/ai";
import { MdComment } from "react-icons/md";
import Link from "next/link";



// id:string,
// thumbnail:string,
// title:string,
// createdAt:string,
// likes?:number,
// comments?:number,
// description?:string,
// category:string,

export default async function Channel({ params }: { params: { id: string } }) {
    const data: {
        subscribers: number,
        subscribersLast7Days: number,
        totalVideos: number

    } = await get('getChannelSummary', {}, headers())

    const content = await get('getContent', {}, headers())
    console.log(content, "dataaaaa")
    return (
        <div className="w-full h-[90%] p-5   ">
            <h1 className="text-2xl opacity-90 tracking-wide  font-bold my-3">
                Channel Dashboard
            </h1>
            <div className="flex  w-full">

                <Card classname={'m-3 min-w-[300px] min-h-[400px]'} >
                    <div className="text-xl ">

                        Latest Videos
                    </div>
                    <div className="">
                        {
                            content?.slice(0, 3).map((item: Content, index: number) => {
                                return (
                                    <div key={index} className="flex items-center gap-3 my-2  border border-background bg-[#0F0F0F] p-3 rounded-lg hover:bg-purple3  ">
                                        <img src={item.thumbnail} className="w-20 h-20 rounded-md" />
                                        <div className=" flex">

                                            <div>
                                                <div className="text-lg font-bold mb-2">
                                                    {item.title}
                                                </div>
                                                <div className="text-sm opacity-60">
                                                    {new Date(item.createdAt).toDateString()}
                                                </div>
                                            </div>
                                            <div className=" ml-7 flex flex-col  " >
                                                <div className="flex mb-2">
                                                    <div className="flex mb-1 items-center  ">
                                                        <div className="text-sm opacity-70  mx-2">
                                                            {item.likes}
                                                        </div>
                                                        <AiFillLike className="fill-green-400" />
                                                    </div>
                                                    <div className="flex mb-1 items-center">
                                                        <div className="text-sm opacity-70 mx-2">
                                                            {item.likes}
                                                        </div>
                                                        <AiFillDislike className="fill-red-500" />
                                                    </div>

                                                </div>

                                                <div className="flex mb-1 self-end items-center">
                                                    <div className="text-sm opacity-70 mx-3">
                                                        {item.comments}
                                                    </div>
                                                    <MdComment className="fill-yellow-500" />
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                )
                            })
                        }
                        <Link href={params.id+'/content'} className="w-full p-3 items-center justify-center flex bg-background3 rounded-xl">
                        View All
                        </Link>
                    </div>
                </Card>
                <Card classname={'m-3 min-w-[300px] h-fit'}>
                    <div className="">
                        <div className="mb-5">
                            <div className="text-xl ">

                                Current Subscribers
                            </div>
                            <div className="text-xl mt-2  font-extrabold">
                                {data.subscribers}
                            </div>
                        </div>
                        <div>
                            <div className="mb-3 text-lg ">
                                Summary
                            </div>
                            <div className="mb-2 text-sm flex justify-between">
                                <h3 className="opacity-60">

                                    {" subscribers(last 7 days)"}
                                </h3>
                                <h3 className="font-bold">
                                    {data.subscribersLast7Days}
                                </h3>
                            </div>
                            <div className="mb-2 text-sm flex justify-between">
                                <h3 className="opacity-60">

                                    Total Videos
                                </h3>
                                <h3 className="font-bold">
                                    {data.totalVideos}
                                </h3>
                            </div>
                        </div>

                    </div>

                </Card>
            </div>
        </div>
    )

}