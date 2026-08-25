"use client";

import React from "react";
import { tv, VariantProps } from "tailwind-variants";

const textField = tv({
  base: `
    w-full flex flex-row justify-start items-center border-0
    font-sans font-normal text-sm md:text-base text-left
    outline-2 cursor-text transition-all duration-200 ease-in-out
    placeholder:font-sans placeholder:font-normal placeholder:text-sm md:placeholder:text-base placeholder:text-left
    disabled:cursor-default
  `,
  variants: {
    color: {
      auth: `
        bg-gray-100 text-(--primary) outline-gray-100 placeholder:text-gray-400
        focus:bg-(--light) focus:outline-(--primary) focus:ring-4 focus:ring-(--primary)/30
        filled:bg-(--light) filled:outline-(--primary)
        disabled:font-medium disabled:text-gray-500
        filled:disabled:bg-gray-100 filled:disabled:outline-gray-100
      `,
    },
    size: {
      sm: "px-1.5 md:px-2 py-1 md:py-1.5 rounded-sm",
    },
    isValueValid: {
      true: "",
      false: "outline-red-500 filled:outline-red-500",
    },
  },
  defaultVariants: {
    isValueValid: true,
  },
});

interface TextFieldProps {
  placeholder: string;
  value: string;
  onChange: (value: string) => void;
  isDisabled?: boolean;
  isSecure?: boolean;
  maxLength?: number;
  minLength?: number;
  variant: VariantProps<typeof textField>;
  customStyle?: string;
  children?: React.ReactNode;
}

const TextField: React.FC<TextFieldProps> = ({
  placeholder,
  value,
  onChange,
  isDisabled = false,
  isSecure = false,
  maxLength,
  minLength,
  variant,
  customStyle,
  children,
}: TextFieldProps) => {
  return (
    <div className="relative w-full">
      <input
        type={isSecure ? "password" : "text"}
        inputMode="text"
        placeholder={placeholder}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        maxLength={maxLength}
        minLength={minLength}
        disabled={isDisabled}
        className={textField({ ...variant, className: customStyle })}
      />
      {children}
    </div>
  );
};

export default TextField;
