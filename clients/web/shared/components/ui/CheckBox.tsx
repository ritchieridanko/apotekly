"use client";

import React from "react";
import { tv, VariantProps } from "tailwind-variants";

const checkBox = tv({
  base: `
    flex justify-center items-center
    transition-all duration-200 ease-in-out
  `,
  variants: {
    color: {
      primary: "border-2",
    },
    size: {
      sm: "size-3.5 md:size-4 rounded-sm",
    },
    isChecked: {
      true: "bg-(--primary) border-(--primary)",
      false: "bg-(--light) border-gray-400",
    },
  },
});

interface CheckBoxProps {
  isChecked: boolean;
  onChange: (value: boolean) => void;
  isDisabled?: boolean;
  variant: VariantProps<typeof checkBox>;
  customStyle?: string;
  children?: React.ReactNode;
}

const CheckBox: React.FC<CheckBoxProps> = ({
  isChecked,
  onChange,
  isDisabled = false,
  variant,
  customStyle,
  children,
}: CheckBoxProps) => {
  return (
    <label
      className={`
        flex flex-row justify-center items-center gap-2
        ${!isDisabled && "cursor-pointer"}
      `}
    >
      <input
        type="checkbox"
        className="hidden"
        checked={isChecked}
        onChange={(e) => {
          if (!isDisabled) {
            e.stopPropagation();
            onChange(e.target.checked);
          }
        }}
        disabled={isDisabled}
      />
      <div className={checkBox({ ...variant, className: customStyle })}>
        {isChecked && (
          <svg
            viewBox="0 0 24 24"
            className="size-4 text-(--light)"
            fill="none"
            stroke="currentColor"
            strokeWidth="4"
          >
            <path
              d="M5 13l4 4L19 7"
              strokeLinecap="round"
              strokeLinejoin="round"
            />
          </svg>
        )}
      </div>
      {children}
    </label>
  );
};

export default CheckBox;
