import { atom } from "recoil";
import { AtomEffect } from 'recoil';

const store = typeof window !== 'undefined' ? window.localStorage : null;

export const localStorageEffect: (key: string) => AtomEffect<any> =
  (key) =>
  ({ setSelf, onSet }) => {
    if (store) {
      const savedValue = store.getItem(key);
      if (savedValue != null) {
        setSelf(JSON.parse(savedValue));
      }

      onSet((newValue, _, isReset) => {
        isReset ? store.removeItem(key) : store.setItem(key, JSON.stringify(newValue));
      });
    }
  };

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