"use client"

import { useEffect, useState } from "react";
import { toast } from "react-toastify";
import * as api from "~/api";

type ChatMessage = {
    message: string;
    userId: string | null;
    streamId: string
};

export default function Chat({ streamId }: {
    streamId: string
}) {
    const [chats, setChats] = useState<ChatMessage[]>([]);
    const [text, setText] = useState<string>("");
    const [socket, setSocket] = useState<WebSocket | null>(null);
    const [userId, setUserId] = useState<string | null>(null)

    useEffect(() => {


        async function connect() {
            const data: {
                userId: string

            } = await api.get('getUserId', {})
            setUserId(data.userId)
            console.log(data.userId, "userId", streamId)
            if (streamId === "") {
                return
            }
            const previousChats = await api.post('getChats', JSON.stringify({ videoId: streamId }))
            setChats([...previousChats.reverse()])
            console.log(previousChats, "previousChats")
            console.log(streamId, "streamId", data.userId)
            const socket = new WebSocket("ws://localhost:8080/chat", ["streamId", streamId, data.userId]);

            socket.onopen = () => {
                console.log("WebSocket connection established");
            };

            socket.onmessage = (event) => {
                const message = JSON.parse(event.data);
                console.log(message, "Received message")
                if (message?.error) {
                    toast.error(message?.error)
                    return
                }
                setChats((prevChats) => [...prevChats, message]);
            };
            setSocket(socket);
        }
        connect()

        return () => {
            socket?.close();
        }

    }, [streamId]);

    const sendMessage = () => {
        if (text.trim() === "") return;

        const newMessage: ChatMessage = {
            message: text,
            userId: userId, // Replace with actual user ID    
            streamId: streamId
        };
        console.log(newMessage, "Sending message");

        // Send message to the server
        socket?.send(JSON.stringify(newMessage));

        // Clear input field
        setText("");
    };

    return (
        <div className="w-full h-full bg-card rounded-md flex flex-col   p-3 border border-[#323232] ">
            <div className="flex-1 flex flex-col justify-end ">
                {chats.map((chat, index) => (
                    <div key={index} className="text-white">
                        {chat.userId}: {chat.message}
                    </div>
                ))}
            </div>
            <div className=" rounded-full ">
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
