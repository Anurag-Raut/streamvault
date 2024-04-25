import Sidebar from "./_components/sidebar"



export default function Studio({children,params}:{
    children:React.ReactNode,
    params:{
        id:string,
        
    }
}){

    return(
        <div className="w-full h-[calc(100vh-81px)] flex flex-row min-h-[calc(100vh-81px)]  ">
            <Sidebar id={params.id}   />
            
            {children}
        </div>
    )
}