import { headers } from "next/headers";
import { get, post } from "~/api";
import Avatar from "~/app/_components/avatar";



export default async function Channel({params}:{params:{username:string}}) {
    console.log(params,"params")
    const data:{
        username:string,
        profileImage:string
        userId:string
        subscribers:number
    } = await post('getUserDetailsByUsername',JSON.stringify(params.username),{})
    console.log(data,"dataA Adasd asd")
    return (
        <div className=" h-full  w-full p-6">
            <div className="flex">
                <Avatar size={200} src={data.profileImage} name={data.username} />
                <div className=" ml-8 pt-5">
                    <div className="text-5xl  font-bold ">{data?.username}</div>
                    <div className="opacity-70">{data.subscribers} subscribers</div>
                </div>
            </div>
        </div>
    )
}