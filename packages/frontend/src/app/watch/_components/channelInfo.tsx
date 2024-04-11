"use client";
import { useState } from "react"
import { toast } from "react-toastify";
import { post } from "~/api";


export default function ChannelInfo({ noOfSubscribers, channelName ,isSubscribed,creatorId}: { noOfSubscribers: number, channelName: string,isSubscribed:boolean,creatorId:string }) {
    const [subscribed, setSubscribed] = useState(isSubscribed)
    const [subscribers,setSubscribers]=useState<number>(noOfSubscribers)

    async function onClick() {
        try {
            if (subscribed) {
                const res = await post('unsubscribe', JSON.stringify(creatorId))
                setSubscribed(false)
                setSubscribers((prev:number)=>prev-1)
            }
            else {
                const res = await post('subscribe', JSON.stringify(creatorId))
                setSubscribed(true)
                setSubscribers((prev:number)=>prev+1)
            }
        }
        catch (error: any) {
            toast.error(error.toString())
        }

    }

    return (
        <div className="flex flex-row">

            <div className="text-lg m-3">{channelName}
                <div className="text-md opacity-80">
                    {subscribers}
                </div>
            </div >
            <div className="flex items-center">
            
            <button onClick={onClick} className={`btn ${subscribed?"btn-primary":"btn-neutral"} rounded-full m-3 `}>{subscribed?"Subscribed":"Subscribe"}</button>
        </div>

        </div>
    )
}