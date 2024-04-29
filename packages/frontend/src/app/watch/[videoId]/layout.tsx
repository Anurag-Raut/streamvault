
import { post } from "~/api"
import LikeComponent from "../_components/likeComponent"
import { cookies, headers } from "next/headers"
import { LikeState } from "../types"
import ChannelInfo from "../_components/channelInfo"
import Chat from "~/app/_components/chat"



type Video = {
    title: string,
    u: string,
    description: string,
    category: string,
    comments: number,
    createdAt: string,
    thumbnail: string,
    likes: number,
    dislikes: number,
    user: {
        id: string,
        username: string,
        profileImage: string

    }
    likeState: LikeState,
    subscribers: number,
    isSubscribed: boolean

}



export default async function Watch({ params, children }: { params: { videoId: string }, children: React.ReactNode }) {


    const data: Video = await post('getVideoData', JSON.stringify(params.videoId), {
        Cookie:cookies().toString(),

    })


    console.log(data, "name")
    return (
        <div className='w-full h-full p-9 flex justify-between'>
            <div className='   h-full ' >
                {children}
                <div className='mt-3 w-full'>
                    <div className='text-3xl '>{data.title}</div>
                    <div className="flex w-full m-3 items-center">
                        <ChannelInfo profileImage={data.user.profileImage} isSubscribed={data.isSubscribed} noOfSubscribers={data.subscribers} creatorId={data.user.id} channelName={data.user.username} />
                        <LikeComponent videoId={params.videoId} likes={data.likes} dislikes={data.dislikes} likeState={data.likeState} />
                    </div>



                </div>
            </div>
            <div className=" w-[500px]">
                <Chat streamId={params.videoId} />

            </div>


        </div>
    )

}