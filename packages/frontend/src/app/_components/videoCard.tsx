'use client'
import { useRouter } from "next/navigation"
import React, { useState } from "react"
import VideoJS from "../watch/_components/player"
import Player from "video.js/dist/types/player"
import videojs from "video.js"
import Link from "next/link"


export default function VideoCard({ title, thumbnail, category, id ,user}: {
    title: string,
    thumbnail: string,
    category?: string,
    src:string,
    id: string,
    user:{
        username:string,
        id:string
    }
}) {

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
        <Link href={`/watch/${id}`} onMouseEnter={() => { setHovering(true) }} onMouseLeave={() => { setHovering(false) }} className="h-[200px] ">
            {
                !hovering ?

                    <img src={thumbnail} alt="" className=" w-[340px] h-[191px] rounded-xl" />


                    :
                    <div className=' w-[340px] h-[191px] rounded-xl'>
                        <VideoJS options={videoJsOptions} onReady={handlePlayerReady} />
                    </div>

            }

            <p>{title}</p>
            <p> {user.username}</p>


        </Link>
    )
}