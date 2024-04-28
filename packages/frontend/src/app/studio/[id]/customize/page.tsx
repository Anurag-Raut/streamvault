import Card from "~/app/_components/card";
import TextInput from "../_components/textInput";
import CustomizeComponent from "./customizeComponent";
import { User } from "~/app/_components/header";
import { headers } from "next/headers";
import { get, post } from "~/api";

export default async function Customize() {
    const user:User=await post('getLoggedUserDetails',{},{},headers())
    return (
        <div className="w-full h-[90%] p-5">
            <h1 className="text-xl my-3">
                Channel Customize
            </h1>
          <CustomizeComponent username={user.username} profileImage={user.profileImage}  />
        </div>
    )
}