    "use client"
    import { ChangeEvent, useRef } from "react";


    export default function TextArea({
        placeholder,
        value,
        onChange,
        classname,
        height=50

    }:{
        placeholder?: string;
        value?: string;
        onChange?: (value: string) => void;
        classname?: string;
        height?: number;
    }) {
        const ref = useRef<HTMLTextAreaElement>(null);

        const handleInput = (e: ChangeEvent<HTMLTextAreaElement>) => {
            if (ref.current) {
                if(e.target.scrollHeight < height) return;
                ref.current.style.height = "auto";
                ref.current.style.height = `${e.target.scrollHeight - 5}px`;
            }
            console.log(e.target.value, "scrollHeight")
        };


        return (
            <textarea
                ref={ref}
                rows={1}
                onInput={handleInput}
                onChange={(e) => onChange?.(e.target.value)}
                value={value}
                className={`textarea textarea-bordered w-full h-full resize-none max-h-[200px] bg-transparent border-1 focus:border-primary  focus:outline-none ${classname}`}
                placeholder={placeholder}
            ></textarea>
        );
    }