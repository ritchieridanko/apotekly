"use client";

import React from "react";
import { tv, VariantProps } from "tailwind-variants";

const button = tv({
  base: `
    w-full flex flex-row justify-center items-center border-0
    font-sans font-normal text-sm md:text-base text-left
    outline-2 cursor-pointer transition-all duration-200 ease-in-out
    disabled:cursor-default
  `,
  variants: {
    color: {
      primary: `
        bg-(--primary) text-(--light) outline-(--primary)
        hover:bg-(--primary)/75 hover:outline-(--primary)/75
        active:bg-(--primary)/90 active:outline-(--primary)/90
        disabled:bg-gray-400 disabled:outline-gray-400
      `,
      secondary: `
        bg-(--light) text-(--primary) outline-(--primary)
        hover:bg-(--primary)/5
        active:bg-(--primary)/10
        disabled:bg-gray-200 disabled:text-gray-500 disabled:outline-gray-200
      `,
    },
    size: {
      sm: "px-1.5 md:px-2 py-1 md:py-1.5 rounded-sm",
    },
  },
});

interface ButtonProps {
  onClick: () => void;
  isDisabled?: boolean;
  variant: VariantProps<typeof button>;
  customStyle?: string;
  children?: React.ReactNode;
}

const Button: React.FC<ButtonProps> = ({
  onClick,
  isDisabled = false,
  variant,
  customStyle,
  children,
}: ButtonProps) => {
  return (
    <button
      type="button"
      onClick={(e) => {
        if (!isDisabled) {
          e.preventDefault();
          e.stopPropagation();
          onClick();
        }
      }}
      disabled={isDisabled}
      className={button({ ...variant, className: customStyle })}
    >
      {children}
    </button>
  );
};

export default Button;
