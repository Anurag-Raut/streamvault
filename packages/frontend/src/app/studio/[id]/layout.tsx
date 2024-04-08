import Sidebar from "./_components/sidebar"



export default function Studio({children,params}:{
    children:React.ReactNode,
    params:{
        id:string
    }
}){

    return(
        <div className="w-full h-full flex flex-row ">
            <Sidebar id={params.id}/>
            
            {children}
        </div>
    )
}