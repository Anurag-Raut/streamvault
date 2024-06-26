import { headers } from "next/headers"
import { get } from "~/api"
import { UserDetails } from "~/app/_components/chat"
import VideoCard from "../_components/VideoCard"
import { Content } from "~/app/studio/[id]/content/page"


export default async function LiveStreams({ params }: {
    params: {
        username: string
    }
}) {
    const username = decodeURIComponent(params.username)
    const contents = await get(`getContent?isVOD=false&username=${username}`, {}, new Headers(headers()))
    console.log(contents, "contentwasa")
    return (
        <div className="w-[100%] h-full p-5 grid xl:grid-cols-3 lg:grid-cols-2 md:grid-cols-1 gap-y-[90px] ">
            {
                contents?.map((content: {
                    title: string, thumbnail: string, id: string, createdAt: string, user: {
                        username: string,
                        user: string,
                        profileImage: string

                    },
                    views: number

                }, index: number) => (
                    <VideoCard title={content.title} thumbnail={content.thumbnail} id={content.id} createdAt={content.createdAt} views={content.views} />

                ))

            }
        </div>
    )



}