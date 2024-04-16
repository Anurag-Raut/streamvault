import { headers } from "next/headers"
import { get } from "~/api"
import { Content } from "~/app/studio/[id]/content/page"


export default async function VODS({ params }: {
    params: {
        username: string
    }
}) {

    const contents = await get('getContent?isVOD=true', { }, headers())
    console.log(contents, "content")
    return (
        <div>
                {
                    contents.map((content:Content,index:number)=>(
                        <div key={index} className="flex flex-row gap-3">
                            <div className="w-1/4">
                                <img src={content.thumbnail} alt="thumbnail" className="w-full h-20" />
                            </div>
                            <div className="w-3/4">
                                <div className="text-lg font-bold">{content.title}</div>
                                <div className="text-sm">{content.createdAt}</div>
                                <div className="text-sm">{content.category}</div>
                                <div className="text-sm">{content.likes} likes</div>
                                <div className="text-sm">{content.comments} comments</div>
                            </div>
                        </div>

                    ))

                }
        </div>
    )




}