"use client"
import React from "react";
import { AnimatePresence, motion } from "framer-motion";

export default function TextInput({ value, setValue, placeholder = "Enter", type = "text" }: {
    value: string,
    setValue: React.Dispatch<React.SetStateAction<string>>,
    placeholder: string,
    type: string
}) {
    const [focus, setFocus] = React.useState(false);

    return (
        <div className="m-3 min-w-[400px]">
            <motion.input

                onFocus={() => {
                    setFocus(true); // Update focus state to true when input is focused
                }}
                onBlur={() => {
                    setFocus(false); // Update focus state to false when input loses focus
                }}
                value={value}
                onChange={(e) => setValue(e.target.value)}
                type="text"
                placeholder={placeholder}
                className={`focus:outline-0 p-2 w-full  bg-transparent ${focus ? 'border-blue-500' : 'border-purple-500'}`}
            />
            <div className="flex">

            <AnimatePresence>

                {focus && <motion.div className="rounded-xl " initial={{
                    width: 0,
                    height: 1.5,
                    backgroundColor: "#cc00ff"
                    
                }}
                
                animate={{
                    width: "100%",
                    height: 1.5,
                    backgroundColor: "#cc00ff"
                    
                }}
                
                exit={{
                    width: 0,
                    height: 1.5,
                    backgroundColor: "#cc00ff"
                    
                    
                }
            }
            
            
            /> 
            // :<div className="w-full h-0.5 rounded-xl bg-white opacity-50" />
                }
                    </AnimatePresence>
                    <AnimatePresence>

                {!focus && <motion.div className="rounded-xl " initial={{
                    width: 0,
                    height: 1.5,
                    backgroundColor: "white",
                    alignSelf:"flex-end"
                    
                }}
                
                animate={{
                    width: "100%",
                    height: 1.5,
                    backgroundColor: "white"
                    
                }}

                    exit={{
                        width: 0,
                        height: 1.5,
                        backgroundColor: "white"
                        
                        
                    }
                }
                
                
                /> 
            }

            </AnimatePresence>
            </div>

        </div>
    );
}
