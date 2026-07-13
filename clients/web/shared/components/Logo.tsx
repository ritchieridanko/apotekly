import React from "react";

interface LogoProps {
  size: "xl" | "2xl" | "3xl" | "4xl" | "5xl" | "6xl" | "7xl" | "8xl";
}

const textSize: Record<LogoProps["size"], string> = {
  xl: "text-lg md:text-xl",
  "2xl": "text-xl md:text-2xl",
  "3xl": "text-2xl md:text-3xl",
  "4xl": "text-3xl md:text-4xl",
  "5xl": "text-4xl md:text-3xl",
  "6xl": "text-5xl md:text-4xl",
  "7xl": "text-6xl md:text-5xl",
  "8xl": "text-7xl md:text-6xl",
};

const Logo: React.FC<LogoProps> = ({ size }: LogoProps) => {
  return (
    <h1
      className={`
        w-full font-mono font-bold md:font-extrabold text-(--primary) text-center select-none
        ${textSize[size]}
      `}
    >
      {process.env.NEXT_PUBLIC_APP_NAME}
    </h1>
  );
};

export default Logo;
