"use client";

import React, { useEffect, useState } from "react";

import { CheckCircle, Eye, EyeSlash } from "@/shared/assets/icons";
import { TextField as Field } from "@/shared/components";

interface TextFieldProps {
  label: string;
  placeholder: string;
  value: string;
  onChange: (value: string) => void;
  isDisabled?: boolean;
  isLoading?: boolean;
  withSecurity?: boolean;
  maxLength?: number;
  minLength?: number;
  validLabel?: string;
  invalidLabel?: string;
}

const TextField: React.FC<TextFieldProps> = ({
  label,
  placeholder,
  value,
  onChange,
  isDisabled = false,
  isLoading = false,
  withSecurity = false,
  maxLength,
  minLength,
  validLabel,
  invalidLabel,
}: TextFieldProps) => {
  const [visible, setVisible] = useState<boolean>(false);

  useEffect(() => {
    if (isDisabled || isLoading) setVisible(false);
  }, [isDisabled, isLoading]);

  return (
    <div className="w-full flex flex-col justify-center items-start gap-0.5 md:gap-1">
      <div className="w-full flex flex-row justify-between items-center gap-0.5 md:gap-1">
        <p className="w-full font-sans font-normal text-(--primary) text-sm md:text-base text-left truncate">
          {label}
        </p>
        {maxLength !== undefined && (
          <p
            className={`
            font-sans font-normal text-sm md:text-base tracking-wide
            ${value.length === maxLength ? "text-red-500" : "text-(--primary)"}
          `}
          >
            {value.length}/{maxLength}
          </p>
        )}
      </div>
      <Field
        placeholder={isDisabled || isLoading ? "" : placeholder}
        value={value}
        onChange={onChange}
        isDisabled={isDisabled || isLoading}
        isSecure={withSecurity && !visible}
        maxLength={maxLength}
        minLength={minLength}
        variant={{ color: "auth", size: "sm", isValueValid: !invalidLabel }}
        customStyle={withSecurity || validLabel ? "pr-8 md:pr-10" : ""}
      >
        {validLabel && (
          <span className="absolute top-1 md:top-1.5 right-1 md:right-2 h-fit w-fit text-green-500">
            <CheckCircle size={6} strokeWidth={1.4} />
          </span>
        )}
        {withSecurity && (
          <button
            type="button"
            onClick={() => setVisible((prev) => !prev)}
            disabled={isDisabled || isLoading}
            className="
              absolute top-1 md:top-1.5 right-1 md:right-2 h-fit w-fit text-(--dark)
              cursor-pointer transition-all duration-200 ease-in-out
              hover:text-(--primary) disabled:text-gray-500 disabled:cursor-default
            "
          >
            {visible ? (
              <Eye size={6} strokeWidth={1.4} />
            ) : (
              <EyeSlash size={6} strokeWidth={1.4} />
            )}
          </button>
        )}
      </Field>
      {validLabel && (
        <p className="w-full font-sans font-normal text-green-500 text-xs md:text-sm text-left line-clamp-2">
          {validLabel}
        </p>
      )}
      {invalidLabel && (
        <p className="w-full font-sans font-normal text-red-500 text-xs md:text-sm text-left line-clamp-2">
          {invalidLabel}
        </p>
      )}
    </div>
  );
};

export default TextField;
