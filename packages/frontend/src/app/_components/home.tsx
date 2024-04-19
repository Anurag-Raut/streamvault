"use client";
import { useEffect, useState } from "react";
import { toast } from "react-toastify";
import { array } from "zod";
import VideoCard from "./videoCard";
import { get } from "~/api";




export default function Home() {

    const [data, setData] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        async function fetchData() {
            try {
                // const res = await fetch('http://localhost:8080/streams', {
                //     method: "GET",
                //     headers: {
                //         "Content-Type": "application/json"
                //     },
                //     credentials: 'include'
                // },
                // );
                // const data = await res.json();

                const data = await get('getVideos')


                console.log(data, "data")
                setData(data);
                setLoading(false);


            }
            catch (error) {
                console.log(error)
                setLoading(false);
                toast.error("Error fetching data")
            }

        }

        fetchData()

    }, [])



    return (
        <div className="w-[100%] p-5 grid xl:grid-cols-3 lg:grid-cols-2 md:grid-cols-1 gap-y-[90px] overflow-y-auto">
            {loading ? Array.from({ length: 21 }, (_, i) => (
                <div data-theme="mytheme" key={i} className="skeleton w-[340px] h-[200px]  "></div>
            ))
                :

                data.map(({
                    title,
                    thumbnail,
                    description,
                    category,
                    id,
                    user,
                    src,
                    createdAt,
                    views

                }: {
                    title: string,
                    thumbnail: string,
                    description: string,
                    category: string,
                    id: string,
                    src: string,
                    createdAt: string,
                    views: number,
                    user: {
                        username: string,
                        id: string
                        profileImage: string
                    }


                }, index: number) => (



                    <VideoCard title={title} thumbnail={thumbnail} id={id} user={user} createdAt={createdAt} views={views} />


                ))

            }


        </div >
    )
}