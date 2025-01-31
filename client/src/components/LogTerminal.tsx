"use client";
import { SetStateAction } from "react";

export const LogTerminal = ({
  isOpen,
  logs,
  setIsOpen,
}: {
  isOpen: boolean;
  setIsOpen: (value: SetStateAction<boolean>) => void;
  logs: string[];
}) => {
  const toggleSidebar = () => {
    setIsOpen(!isOpen);
  };

  return (
    <div className="relative bg-[#171717] text-white w-[100%]">
      <div
        className={`fixed md:overflow-y-scroll md:pb-[5.2rem] overflow-x-hidden bottom-0 h-[50vh] md:right-0 md:top-[5.2rem] z-40 md:h-screen md:w-[70%] w-[100%] max-w-[768px] transform bg-[#171717] shadow-2xl ${
          isOpen ? "translate-x-0 " : "translate-x-full "
        } transition-transform duration-100 ease-in-out`}
      >
        <div className="m-3 border-b-[1px] border-gray-500 border-solid text-3xl w-[100%] h-12 text-[#989898] font-bold">
          <h2 className="border-b-[4px] border-gray-500 border-solid w-20 pb-[0.6rem] ">
            Logs
          </h2>
        </div>
        <div className="p-3">
          {logs.map((log, idx) => {
            return (
              <p
                key={idx}
                className="text-[#989898] my-3 md:text-xl text-lg font-semibold"
              >
                $ {log}
              </p>
            );
          })}
        </div>
      </div>

      {isOpen && (
        <div
          onClick={toggleSidebar}
          className="fixed cursor-pointer inset-0 z-10 top-[5.2rem] bg-black bg-opacity-50"
        ></div>
      )}
    </div>
  );
};
