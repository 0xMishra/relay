"use client";
import { SetStateAction } from "react";

export const LogTerminal = ({
  isOpen,
  setIsOpen,
}: {
  isOpen: boolean;
  setIsOpen: (value: SetStateAction<boolean>) => void;
}) => {
  const toggleSidebar = () => {
    setIsOpen(!isOpen);
  };

  return (
    <div className="relative bg-[#171717] text-white w-[100%]">
      <div
        className={`fixed right-0 top-0 z-40 h-screen md:w-[70%] w-[95%] max-w-[700px] transform bg-[#171717] shadow-2xl ${
          isOpen ? "translate-x-0" : "translate-x-full"
        } transition-transform duration-500 ease-in-out`}
      >
        <div className="p-6"></div>
      </div>

      {isOpen && (
        <div
          onClick={toggleSidebar}
          className="fixed inset-0 z-10 bg-black bg-opacity-50"
        ></div>
      )}
    </div>
  );
};
