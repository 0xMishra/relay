import Link from "next/link";
import React from "react";
import { FaGithub } from "react-icons/fa";

const Navbar = () => {
  return (
    <div className="fixed bg-[#171717] top-0 z-[5] flex w-screen items-center justify-center border-b-[1px] border-solid border-gray-600 p-3 md:ml-3">
      <div className="flex w-[90vw] max-w-[700px] items-center justify-between">
        <h1 className="bg-gradient-to-b to-gray-500 from-white bg-clip-text pl-2 pr-1 text-3xl font-black tracking-tighter text-transparent">
          Relay
        </h1>
        <div>
          <Link
            href={"https://github.com/0xMishra/relay"}
            target="_blank"
            className="my-3 flex items-center justify-center rounded-lg bg-[#1d2021] px-2 py-1 text-lg font-semibold text-white hover:bg-gray-900"
          >
            <FaGithub />
            <span className="ml-1">Github</span>
          </Link>
        </div>
      </div>{" "}
    </div>
  );
};

export default Navbar;
