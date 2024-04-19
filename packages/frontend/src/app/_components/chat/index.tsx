"use client"

import { useEffect, useState, useRef } from "react";
import { toast } from "react-toastify";
import * as api from "~/api";
import ChatBubble from "./chatBubble";

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
    const [chats, setChats] = useState<ReceivedChatMessage[]>([]);
    const [text, setText] = useState<string>("");
    const [socket, setSocket] = useState<WebSocket | null>(null);
    const [userDetails, setUserDetails] = useState<UserDetails>(null);
    const messagesEndRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        async function connect() {
            const data: UserDetails = await api.get("getUserId", {});
            setUserDetails(data);

            if (streamId === "") {
                return;
            }

            const previousChats = await api.post("getChats", JSON.stringify({ videoId: streamId }));
            setChats([...previousChats.reverse()]);

            const socket = new WebSocket("ws://localhost:8080/chat", ["streamId", streamId, data?.userId ?? ""]);

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
                scrollToBottom();
                setChats((prevChats) => [...prevChats, message]);
              

            };

            setSocket(socket);
        }

        connect();

        return () => {
            socket?.close(1000);
        };
    }, [streamId]);

    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: "instant" });
    };
    useEffect(() => {
        //3️⃣ bring the last item into view        
        messagesEndRef?.current?.scrollIntoView({behavior: "smooth"})
    }, [chats]);

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
            <div className=" flex-col  content-start overflow-y-auto  h-full ">
                {chats.map((chat, index) => (
                    <ChatBubble key={index} message={chat.message} user={chat.user} />
                ))}
                <div ref={messagesEndRef} />
            </div>
            <div className="rounded-full">
                <label className="input input-bordered flex items-between gap-2 w-full rounded-full my-2 mt-3 ">
                    <input
                        onChange={(event) => {
                            setText(event.target.value);
                        }}
                        value={text}
                        type="text"
                        className="bg-red-400 flex-1"
                        placeholder="Type your message..."
                    />
                    <button onClick={sendMessage}>Send</button>
                </label>
            </div>
        </div>
    );
}
