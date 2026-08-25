"use client";

import React from "react";
import { tv, VariantProps } from "tailwind-variants";

const link = tv({
  base: "font-sans underline-offset-2 cursor-pointer hover:underline",
  variants: {
    color: {
      primary: "font-medium text-(--primary) text-sm md:text-base",
    },
    isDisabled: {
      true: "text-gray-500 cursor-default",
      false: "",
    },
  },
});

interface LinkProps {
  onClick: () => void;
  inNewTab?: boolean;
  isDisabled?: boolean;
  variant: VariantProps<typeof link>;
  customStyle?: string;
  children?: React.ReactNode;
}

const Link: React.FC<LinkProps> = ({
  onClick,
  inNewTab = false,
  isDisabled = false,
  variant,
  customStyle,
  children,
}: LinkProps) => {
  return (
    <a
      target={inNewTab ? "_blank" : "_self"}
      onClick={(e) => {
        if (!isDisabled) {
          e.preventDefault();
          e.stopPropagation();
          onClick();
        }
      }}
      className={link({ ...variant, className: customStyle })}
    >
      {children}
    </a>
  );
};

export default Link;
