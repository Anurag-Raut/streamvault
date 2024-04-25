"use client"

import { toast } from "react-toastify";
import { post } from "~/api";
import { AiFillDislike, AiFillLike, AiOutlineDislike, AiOutlineLike } from "react-icons/ai";
import { useEffect, useState } from "react";
import { LikeState } from "../types";
import { set } from "video.js/dist/types/tech/middleware";



export default function LikeComponent({ likes, dislikes, videoId, likeState }: { likes: number; videoId: string, dislikes: number, likeState: LikeState }) {
    const [likeCount, setLikeCount] = useState(likes)
    const [dislikeCount, setDislikeCount] = useState(dislikes)
    const [curretnLikeState, setCurretnLikeState] = useState<LikeState>(likeState)

    useEffect(() => {
        setLikeCount(likes)
        setDislikeCount(dislikes)
        setCurretnLikeState(likeState)
    }, [likes, dislikes, videoId, likeState])
    console.log(likeState, "like state")

    async function like() {
        try {
            if (curretnLikeState === LikeState.Liked) {


                const res = await post('removeLike', JSON.stringify(videoId))
                setLikeCount(likeCount - 1)
                setCurretnLikeState(LikeState.Neutral)
            }
            else {
                const res = await post('like', JSON.stringify(videoId))
                setLikeCount(likeCount + 1)
                if (curretnLikeState === LikeState.Disliked) {
                    setDislikeCount(dislikeCount - 1)
                }
                setCurretnLikeState(LikeState.Liked)
            }
        }
        catch (error: any) {
            toast.error(error.toString())
        }

    }
    async function dislike() {
        try {
            if (curretnLikeState === LikeState.Disliked) {

                const res = await post('removeLike', JSON.stringify(videoId))
                setDislikeCount(dislikeCount - 1)
                setCurretnLikeState(LikeState.Neutral)

            }
            else {
                const res = await post('dislike', JSON.stringify(videoId))
                setDislikeCount(dislikeCount + 1)
                if (curretnLikeState === LikeState.Liked) {
                    setLikeCount(likeCount - 1)
                }
                setCurretnLikeState(LikeState.Disliked)
            }
        }
        catch (error: any) {
            toast.error(error.toString())
        }

    }

    return (
        <div className=" rounded-full p-3 bg-neutral w-auto flex justify-between ">
            <button onClick={like} className=" flex  ">
                {
                    curretnLikeState === LikeState.Liked ?
                        <AiFillLike size={25} className="mx-2 fill-primary " />
                        :
                        <AiOutlineLike size={25} className="mx-2" />

                }
                {likeCount}
            </button>
            <div className="divider divider-horizontal" />

            <button onClick={dislike} className="flex" >
                {
                    curretnLikeState === LikeState.Disliked ?
                        <AiFillDislike size={25} className="mx-2 fill-red" />
                        :
                        <AiOutlineDislike size={25} className="mx-2 " />
                }
                {dislikeCount}
            </button>

        </div>)

}