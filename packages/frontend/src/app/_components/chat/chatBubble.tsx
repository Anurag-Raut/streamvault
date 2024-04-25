"use client";
import { UserDetails } from ".";
import { motion } from "framer-motion";
import Avatar from "../avatar";

export default function ChatBubble({ message, user }: {
    message: string,
    user: {
        username: string,
        userId: string,
        profileImage: string

    }
}) {
    const hRange: [number, number] = [0, 360];
    const sRange: [number, number] = [60, 70];
    const lRange: [number, number] = [55, 65];
    const getHashOfString = (str: string) => {
        let hash = 0;
        for (let i = 0; i < str.length; i++) {
            hash = str.charCodeAt(i) + ((hash << 7) - hash);
        }
        hash = Math.abs(hash);
        return hash;
    };
    const normalizeHash = (hash: number, min: number, max: number) => {
        return Math.floor((hash % (max - min)) + min);
    };
    const generateHSL = (name: string) => {
        const hash = getHashOfString(name);
        const h = normalizeHash(hash, hRange[0], hRange[1]);
        const s = normalizeHash(hash, sRange[0], sRange[1]);
        const l = normalizeHash(hash, lRange[0], lRange[1]);
        return `hsl(${h},${s}%,${l}%)`;
    };
    return (
        <motion.div initial={{
            // x: -10,
            y: 20,
            opacity: 0.2
        }}
            animate={{
                x: 0,
                y: 0,
                opacity: 1
            }}
            transition={{
                type: "spring", stiffness: 100 ,
                duration: 0.5
            }}
            className="flex flex-row items-start justify-start w-fit h-auto  rounded-lg p-2 "
        >


            <Avatar size={25} src={user?.profileImage} name={user?.username} />

            <div style={{
                color: generateHSL("user.username")
            }} className="font-bold text-sm ml-2">
                {user.username}:
            </div>
            <div className="ml-2 bg-red-400  text-wrap text-ellipsis  flex-wrap w-fit max-w-full break-all   ">
                {message}
            </div>


        </motion.div>
    )

}