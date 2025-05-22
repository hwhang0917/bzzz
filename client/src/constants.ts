export interface Emoji {
  descriptor: string;
  value: string;
}

export const EMOJIS: Emoji[] = [
  { descriptor: "Grinning face", value: "😃" },
  { descriptor: "Beaming face with smiling eyes", value: "😁" },
  { descriptor: "Rolling on the floor laughing", value: "🤣" },
  { descriptor: "Winking face with tongue", value: "😜" },
  { descriptor: "Face screaming in fear", value: "😱" },
  { descriptor: "Unamused face", value: "😒" },
  { descriptor: "Loudly Crying Face", value: "😭" },
  { descriptor: "Red heard", value: "❤️" },
];

export const EMOJI_MAP = new Map(EMOJIS.map((e) => [e.descriptor, e.value]));
