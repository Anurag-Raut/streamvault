import { headers } from "next/headers"
import Image from "next/image"
import { post } from "~/api"
import Avatar from "~/app/_components/avatar"
import Card from "~/app/_components/card"

type Comment = {
    user: {
        username: string,
        userId: string,
        profileImage: string

    },
    message: string,
    createdAt: string,
    video: {
        title: string,
        thumbnail: string
    }

}
export default async function Comments({
    params
}: {
    params: {
        id: string
    }
}) {
    const data = await post('/getCommmentsForChannel', JSON.stringify(params.id), {}, new Headers(headers()))
    // console.log(data, "dddaaaa")
    return (
        <div className="w-full h-full p-5">
            <h1 className="text-xl my-3">
                Channel Comments
            </h1>
            <Card>
                <div className="overflow-x-auto">
                    <table className="table">
                        {/* head */}
                        <thead>
                            <tr>
                                <th>User</th>
                                <th>Message</th>
                                <th>Date</th>
                                <th>Video</th>


                                <th></th>
                            </tr>
                        </thead>
                        <tbody>
                            {
                                data?.map((item: Comment, index: Number) => {
                                    return (
                                        <tr className="hoveropacity-50 hover:bg-black rounded-xl">
                                            <td>
                                                <Avatar name={item.user.username} src={item.user.profileImage} size={34} />
                                                <div className="mt-2 text-md font-bold opacity-70">
                                                    {item.user.username}
                                                </div>
                                            </td>
                                            <td>
                                                {item.message}
                                            </td>

                                            <td>
                                                {new Date(item.createdAt).toDateString()}
                                            </td>
                                            <td>
                                                <div className="flex items-center gap-3">
                                                    <div>
                                                        <Image
                                                            src={item.video.thumbnail}
                                                            width={80}
                                                            height={70}
                                                            alt="thumbnail"
                                                            placeholder="blur"
                                                            blurDataURL={"https://img.freepik.com/free-photo/woman-holding-leafy-white-pillow-mockup_53876-128613.jpg?size=626&ext=jpg"}


                                                        />
                                                    </div>
                                                    <div>
                                                        <div className="font-bold">{item.video.title}</div>
                                                        {/* <div className="text-sm opacity-50">{item.description}</div> */}
                                                    </div>
                                                </div>
                                            </td>
                                        </tr>
                                    )
                                })

                            }

                        </tbody>
                        {/* foot */}


                    </table>
                </div>
            </Card>
        </div>
    )
}