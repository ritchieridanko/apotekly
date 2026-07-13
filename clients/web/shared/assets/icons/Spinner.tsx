import React from "react";

import { setIconSize } from "@/shared/utils";

const Spinner: React.FC<IconProps> = ({ size, strokeWidth }: IconProps) => {
  return (
    <svg
      className={`${setIconSize(size)} animate-spin`}
      viewBox="0 0 24 24"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      <circle
        className="opacity-25"
        cx="12"
        cy="12"
        r="10"
        stroke="currentColor"
        strokeWidth={strokeWidth}
      />
      <circle
        className="opacity-75"
        cx="12"
        cy="12"
        r="10"
        stroke="currentColor"
        strokeWidth={strokeWidth}
        strokeDasharray="60"
        strokeDashoffset="40"
      />
    </svg>
  );
};

export default Spinner;
