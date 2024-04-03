import { atom } from "recoil";
const store = typeof window !== 'undefined' ? window.localStorage : null;

const localStorageEffect = (key: string) => ({ setSelf, onSet }: any) => {

    if(store){

        const savedValue = store.getItem(key)
        if (savedValue != null) {
            setSelf(JSON.parse(savedValue));
        }
    
        onSet((newValue: any, _: any,) => {
            store.setItem(key, JSON.stringify(newValue));
            
        });

    }
  
}
export const streamInfo = atom({
    key: "streamInfo",
    default: {
        title: "New Stream",
        description: "",
        viewers: 0,
        game: "",
        thumbnail: "",
        likes: 0,
        category: ""
    },
    effects: [localStorageEffect("streamInfo")],
});