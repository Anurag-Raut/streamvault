


export default function SignUpPage() {
    return(
        <div className="w-full h-full flex justify-center items-center">
            <div className="w-fit h-fit p-5 bg-primary rounded-xl ">
            <h2>Username</h2>
            <input type="text" placeholder="Enter your username" className="input input-bordered w-full max-w-xs m-3" />
            <h2 className="m-2" >Password</h2>
            <input type="password" placeholder="Enter your password" className="input input-bordered w-full max-w-xs m-3" />
            <button className="btn">Sign in</button>


            </div>

        </div>
    )
}
