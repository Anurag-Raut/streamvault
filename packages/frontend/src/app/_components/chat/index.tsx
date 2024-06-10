"use client"

import { useEffect, useState, useRef, useCallback } from "react";
import { toast } from "react-toastify";
import * as api from "~/api";
import ChatBubble from "./chatBubble";

import InfiniteScroll from 'react-infinite-scroll-component';




export type UserDetails = {
    username?: string;
    userId?: string;
    profileImage?: string;
} | null;

type SendChatMessage = {
    message: string;
    streamId: string;
    user: UserDetails;
};

type ReceivedChatMessage = {
    message: string;
    streamId: string;
    user: {
        username: string;
        userId: string;
        profileImage: string;
    };
};

export default function Chat({ streamId }: { streamId: string }) {
    import('ldrs').then(ldrs => {
        const { tailChase } = ldrs;

        // Register components
        tailChase.register();
     
    });
    const [chats, setChats] = useState<ReceivedChatMessage[]>([]);
    const [text, setText] = useState<string>("");
    const [socket, setSocket] = useState<WebSocket | null>(null);
    const [userDetails, setUserDetails] = useState<UserDetails>(null);
    const [finish, setFinish] = useState(false);
    const messagesEndRef = useRef<HTMLDivElement>(null);
    const chatsRef = useRef<HTMLDivElement>(null);

  

    async function getMoreChats(){
        console.log("asdasdasdasd")
        const newChats = await api.post("getChats", JSON.stringify({ videoId: streamId, noOfChats: chats.length }));
        setChats((prev) => [...prev, ...newChats?.chats??[]]);
        setFinish(newChats?.finish??false)
    }




    // const handleScroll = debounce(async (event: any) => {
    //     const { scrollTop, scrollHeight, clientHeight } = event.target;
    //     const scrollRatio = scrollTop / (scrollHeight - clientHeight);
    //     console.log(scrollTop)
    //     if (scrollTop < 100) {
    //         try {
               

    //         }
    //         catch (error) {
    //             setLoadingChats(false)
    //         }
    //     }

    // }, 500)





    // useEffect(() => {
    //     chatsRef.current?.addEventListener("scroll", handleScroll, { passive: true, capture: true })
    //     return () => {
    //         chatsRef.current?.removeEventListener("scroll", handleScroll)

    //     }
    // }, [])

    useEffect(() => {
        async function connect() {
            const data: UserDetails = await api.get("getUserId", {});
            setUserDetails(data);

            if (streamId === "") {
                return;
            }

            const previousChats = await api.post("getChats", JSON.stringify({ videoId: streamId }));
            setChats(previousChats?.chats??[]);
            setFinish(previousChats?.finish??false)
            // setTimeout(() => {

            //     messagesEndRef?.current?.scrollIntoView({ behavior: "smooth" })
            // }, 200)


            const socket = new WebSocket(`${process.env.NEXT_PUBLIC_WS_URL}/chat`, ["streamId", streamId, data?.userId ?? ""]);
            setSocket(socket);
            socket.onopen = () => {
                console.log("WebSocket connection established");
            };

            socket.onmessage = (event) => {
                const message = JSON.parse(event.data);
                console.log(message, "Received message");
                if (message?.error) {
                    toast.error(message?.error);
                    return;
                }
                // scrollToBottom();
                setChats((prevChats) => [message,...prevChats]);


            };


        }

        connect();

        return () => {
            socket?.close(1000);
        };
    }, [streamId]);

    const scrollToBottom = () => {
        const chatContainer = chatsRef.current;
        if (!chatContainer) return;

        const scrollHeight = chatContainer.scrollHeight;
        const scrollTop = chatContainer.scrollTop;
        const clientHeight = chatContainer.clientHeight;

        const distanceFromBottom = scrollHeight - scrollTop - clientHeight;

        if (distanceFromBottom > 100) return;

        messagesEndRef.current?.scrollIntoView({ behavior: "instant" });
    };

    // useEffect(() => {
    //     //3️⃣ bring the last item into view        

    // }, []);

    const sendMessage = () => {
        if (text.trim() === "") return;

        const newMessage: SendChatMessage = {
            message: text,
            streamId: streamId,
            user: userDetails,
        };
        console.log(newMessage, "Sending message");
        // Send message to the server
        socket?.send(JSON.stringify(newMessage));

        // Clear input field
        setText("");
    };

    return (
        <div className="w-full h-full bg-card rounded-md flex flex-col p-3 border border-[#323232]  ">
            <div id="chatScrollDiv"
                 className="   content-start overflow-y-auto  ">
                <InfiniteScroll
                    dataLength={chats?.length}
                    next={getMoreChats}
                    style={{ display: 'flex', flexDirection: 'column-reverse',height:"65vh",overflowY:"auto" }} //To put endMessage and loader to the top.
                    inverse={true} 
                    // className="h-full content-start"
                    
                    hasMore={finish}
                    height={"100%"}
                    endMessage={
                        <p style={{ textAlign: 'center', color:"gray",marginBottom:10 }}>
                          <b>End of Chat</b>
                        </p>
                      }
                    loader={<div className="w-full h-[100px] flex flex-col justify-center items-center mb-3">

                        <l-tail-chase
                            size="40"
                            speed="1.75"
                            color="purple"
                        ></l-tail-chase>
                        <div className=" mt-3 text-white opacity-70 font-bold ">
                            Loading Chat...
                        </div>
                    </div>}
                    scrollableTarget="chatScrollDiv"
                >
                    {chats.map((chat, index) => (
                    <ChatBubble key={index} message={chat.message} user={chat.user} />
                ))}
                </InfiniteScroll>
                <div  ref={messagesEndRef} />


               
            </div>
            <div className="rounded-full">
                <label className="input input-bordered flex items-between gap-2 w-full rounded-full my-2 mt-3 ">
                    <input
                        onChange={(event) => {
                            setText(event.target.value);
                        }}
                        value={text}
                        type="text"
                        className=" flex-1"
                        placeholder="Type your message..."
                    />
                    <button onClick={sendMessage}>Send</button>
                </label>
            </div>
        </div>
    );
}
