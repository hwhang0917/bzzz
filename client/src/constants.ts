export interface Emoji {
  key: string;
  value: string;
}

export const EMOJIS: Emoji[] = [
  { key: "smile", value: "😃" },
  { key: "laugh", value: "😁" },
  { key: "lol", value: "🤣" },
  { key: "wink", value: "😜" },
  { key: "shock", value: "😱" },
  { key: "bored", value: "😒" },
  { key: "heart", value: "❤️" },
];

export const EMOJI_MAP = new Map(EMOJIS.map((e) => [e.key, e.value]));
