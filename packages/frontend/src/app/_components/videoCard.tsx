'use client'
import { useRouter } from "next/navigation"
import React, { useState } from "react"
import VideoJS from "../watch/_components/player"
import Player from "video.js/dist/types/player"
import videojs from "video.js"
import Link from "next/link"
import Avatar from "./avatar"
import millify from "millify";

export function timeAgo(date: string) {
    const seconds = Math.floor((new Date().getTime() - new Date(date).getTime()) / 1000);
    const intervals = {
        year: 31536000,
        month: 2592000,
        week: 604800,
        day: 86400,
        hour: 3600,
        minute: 60
    };

    if (seconds < intervals.minute) {
        return "Just now";
    } else if (seconds < intervals.hour) {
        const minutes = Math.floor(seconds / intervals.minute);
        return `${minutes} minute${minutes > 1 ? 's' : ''} ago`;
    } else if (seconds < intervals.day) {
        const hours = Math.floor(seconds / intervals.hour);
        return `${hours} hour${hours > 1 ? 's' : ''} ago`;
    } else if (seconds < intervals.week) {
        const days = Math.floor(seconds / intervals.day);
        return `${days} day${days > 1 ? 's' : ''} ago`;
    } else if (seconds < intervals.month) {
        const weeks = Math.floor(seconds / intervals.week);
        return `${weeks} week${weeks > 1 ? 's' : ''} ago`;
    } else if (seconds < intervals.year) {
        const months = Math.floor(seconds / intervals.month);
        return `${months} month${months > 1 ? 's' : ''} ago`;
    } else {
        const years = Math.floor(seconds / intervals.year);
        return `${years} year${years > 1 ? 's' : ''} ago`;
    }
}



export default function VideoCard({ title, thumbnail, id, user ,createdAt,views}: {
    title: string,
    thumbnail: string,
    category?: string,
    id: string,
    createdAt: string,
    user: {
        username: string,
        id: string,
        profileImage: string,
    },
    views: number
}) {
console.log(user.profileImage, "profile image")
    const [hovering, setHovering] = useState(false)
    const router = useRouter()



    const playerRef = React.useRef<Player | null>(null);

    const videoJsOptions = {
        autoplay: true,
        // controls: true,
        responsive: true,
        aspectRatio: '16:9',


        fluid: true,
        sources: [{
            src: `http://localhost:8080/hls/${id}/${id}.m3u8`,

        }]
    };

    const handlePlayerReady = (player: any) => {
        playerRef.current = player;

        // You can handle player events here, for example:
        player.on('waiting', () => {
            videojs.log('player is waiting');
        });

        player.on('dispose', () => {
            videojs.log('player will dispose');
        });
    };

    

    return (
        <Link href={`/watch/${id}`} onMouseEnter={() => { setHovering(true) }} onMouseLeave={() => { setHovering(false) }} className="h-[200px] m-3 ">
            {
                !hovering ?

                    <img src={thumbnail} alt="" className=" w-[340px] h-[191px] rounded-xl" />


                    :
                    <div className=' w-[340px] h-[191px] rounded-xl'>
                        <VideoJS options={videoJsOptions} onReady={handlePlayerReady} />
                    </div>

            }
            <div className="flex mt-3 ">
                <div className="mr-3">
                    <Avatar size={35} name={user.username} src={user.profileImage} />
                </div>
                <div>

                <p className="font-bold text-lg">{title}</p>
                <p className="text-md opacity-65 font-medium"> {user.username}</p>
                <div className="flex">
                    <p className="text-md opacity-65 font-medium ">{`${ millify(views)} views`}</p>
                    <p className="text-md opacity-65 ml-3 font-medium"> {timeAgo(createdAt)}</p>
                </div>
                </div>

            </div>

        </Link>
    )
}