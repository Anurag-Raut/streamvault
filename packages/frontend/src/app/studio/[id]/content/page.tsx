
import Card from "~/app/_components/card";
import { cookies, headers } from "next/headers";
import { get } from "~/api";

export type Content={
    id:string,
    thumbnail:string,
    title:string,
    createdAt:string,
    likes?:number,
    comments?:number,
    description?:string,
    category:string,
    views:number
}
export default async function Content({ params }: { params: { id: string } }) {

    

    

    const data= await get('getDashboardContent',{
        Cookie:cookies().toString(),

    })
    console.log(data,"dashboard content")
    return (
        <div className="w-full h-[90%] p-5">
            <h1 className="text-xl my-3">
                Channel Content
            </h1>
            <Card classname={"h-[90%] w-[1000px] "}>
                <div className=" overflow-y-auto h-full w-full ">
                    <table className="table w-full ">
                        {/* head */}
                        <thead>
                            <tr>
                                <th>
                                    <label>
                                        <input type="checkbox" className="checkbox" />
                                    </label>
                                </th>
                                <th>Video</th>
                                <th>Date</th>
                                <th>Category</th>
                                <th>Likes </th>
                                <th>Comments</th>
                                <th>Views</th>

                                <th></th>
                            </tr>
                        </thead>
                        <tbody>
                            {
                                data?.map((item:Content,index:Number)=>{
                                    return (
                                        <tr className="hover:opacity-50 hover:bg-black">
                                        <th>
                                            <label>
                                                <input type="checkbox" className="checkbox" />
                                            </label>
                                        </th>
                                        <td>
                                            <div className="flex items-center gap-3">
                                                <div>
                                                    <img src={item.thumbnail} className="w-20 " alt="" />
                                                </div>
                                                <div>
                                                    <div className="font-bold">{item.title}</div>
                                                    <div className="text-sm opacity-50">{item.description}</div>
                                                </div>
                                            </div>
                                        </td>
                                        <td>
                                            {new Date(item.createdAt).toDateString()}
                                        </td>
                                        <td>
                                            {item.category}
                                        </td>
                                        <td>{item.likes}</td>
                                        <th>
                                            <button className="btn btn-ghost btn-xs">{item.comments}</button>
                                        </th>
                                        <th>
                                                {item.views}
                                        </th>
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