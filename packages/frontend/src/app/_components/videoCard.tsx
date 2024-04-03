

export default function VideoCard({title,thumbnail,category}:{
    title:string,
    thumbnail:string,
    category?:string
}){

    return (
        <div className="">
            <img src={thumbnail} alt="" className=" w-[340px] h-[200px] rounded-xl" />
            <p>{title}</p>

        </div>
    )
}