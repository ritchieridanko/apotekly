import React from "react";

interface AuthFeedbackTileProps {
  message: string;
  icon: React.FC<IconProps>;
  variant: "success" | "error";
}

const AuthFeedbackTile: React.FC<AuthFeedbackTileProps> = ({
  message,
  icon: Icon,
  variant,
}: AuthFeedbackTileProps) => {
  return (
    <div className="w-full flex flex-col justify-center items-center gap-3">
      <span
        className={variant === "success" ? "text-green-500" : "text-red-500"}
      >
        <Icon size={20} strokeWidth={1.5} />
      </span>
      <p
        className={`
          px-2.5 md:px-3 rounded-2xl font-sans font-medium text-(--light) text-sm md:text-base text-center
          select-none
          ${variant === "success" ? "bg-green-500" : "bg-red-500"}
        `}
      >
        {message}
      </p>
    </div>
  );
};

export default AuthFeedbackTile;
