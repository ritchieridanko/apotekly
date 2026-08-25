"use client";

import React from "react";

import { CheckBox as CB } from "@/shared/components";

interface CheckBoxProps {
  label: string;
  isChecked: boolean;
  onChange: (value: boolean) => void;
  isDisabled?: boolean;
  isLoading?: boolean;
}

const CheckBox: React.FC<CheckBoxProps> = ({
  label,
  isChecked,
  onChange,
  isDisabled = false,
  isLoading = false,
}: CheckBoxProps) => {
  return (
    <CB
      isChecked={isChecked}
      onChange={(v) => onChange(v)}
      isDisabled={isDisabled || isLoading}
      variant={{ color: "primary", size: "sm", isChecked: isChecked }}
    >
      <p className="font-sans font-normal text-(--dark) text-sm md:text-base">
        {label}
      </p>
    </CB>
  );
};

export default CheckBox;
