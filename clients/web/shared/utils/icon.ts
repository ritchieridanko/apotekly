const ICON_SIZE: Record<number, string> = {
  2: "size-1 md:size-2",
  3: "size-2 md:size-3",
  4: "size-3 md:size-4",
  5: "size-4 md:size-5",
  6: "size-5 md:size-6",
  7: "size-6 md:size-7",
  8: "size-7 md:size-8",
  9: "size-8 md:size-9",
  10: "size-9 md:size-10",
  11: "size-10 md:size-11",
  12: "size-11 md:size-12",
  20: "size-19 md:size-20",
};

export const setIconSize = (size: number): string => ICON_SIZE[size];
