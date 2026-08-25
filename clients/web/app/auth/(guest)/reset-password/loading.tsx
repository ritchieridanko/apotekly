import { Spinner } from "@/shared/assets/icons";

const ResetPasswordLoading: React.FC = () => {
  return (
    <div className="mt-2 w-full flex flex-col justify-center items-center gap-2 text-(--primary)">
      <Spinner size={12} strokeWidth={3.5} />
      <p className="font-sans font-medium text-sm md:text-base text-center select-none">
        Verifying Your Token...
      </p>
    </div>
  );
};

export default ResetPasswordLoading;
