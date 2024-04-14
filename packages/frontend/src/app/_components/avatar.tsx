

export default function Avatar({ src, size, name, }: { src: string; size: number, name: string }) {


    const hRange:[number,number] = [0, 360];
    const sRange:[number,number]  = [50, 60];
    const lRange:[number,number]  = [45, 55];
    const getHashOfString = (str: string) => {
        let hash = 0;
        for (let i = 0; i < str.length; i++) {
            hash = str.charCodeAt(i) + ((hash << 5) - hash);
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
        return`hsl(${h},${s}%,${l}%)`;
    };


    const getFirstLetter = (name: string) => {
        console.log(name)
        return name.charAt(0).toUpperCase();
    }
console.log(generateHSL(name))
    return (
        <div className={` w-${10} h-${10} rounded-full flex justify-center items-center`}>
            {src ?
                <div className="w-full rounded-full">
                    <img src={src} />
                </div>
                :
                <div className={` w-full h-full flex rounded-full justify-center items-center self-center text-white font-bold text-xl `} style={{
                    backgroundColor: generateHSL(name)
                }}>
                    {getFirstLetter(name)}
                </div>
            }
        </div>
    )
}
