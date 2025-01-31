"use client";
import { LogTerminal } from "@/components/LogTerminal";
import Navbar from "@/components/Navbar";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { toast } from "@/hooks/use-toast";
import axios from "axios";
import { Loader2 } from "lucide-react";
import Link from "next/link";
import { generate } from "random-words";
import { useEffect, useState } from "react";

export default function Home() {
  const [isLoading, setIsloading] = useState<boolean>(false);
  const [gitUrl, setGitUrl] = useState<string>("");
  const [isTermOpen, setIsTermOpen] = useState<boolean>(false);
  const [projectId, setProjectId] = useState<string>(generate() as string);
  const [logs, setLogs] = useState<string[]>(["setting up the server..."]);
  const [deploymenUrl, setDeploymenUrl] = useState<string>("");

  async function deployProject() {
    try {
      setIsloading(true);

      const pId = localStorage.getItem(gitUrl.trim());

      if (pId) {
        let res = await axios.post(`http://localhost:3000/project`, {
          url: gitUrl,
          pid: pId,
        });
        console.log(res.data);
        setProjectId(res.data.url.split(".")[0]);
        setDeploymenUrl(res.data.url);
        return;
      }

      let res = await axios.post(`http://localhost:3000/project`, {
        url: gitUrl,
      });
      console.log(res.data);
      setProjectId(res.data.url.split(".")[0]);
      setDeploymenUrl(res.data.url);

      localStorage.setItem(gitUrl.trim(), projectId);
    } catch (error) {
      console.log(error);
      setIsloading(false);
      toast({ variant: "destructive", description: `${error}` });
    }
  }

  useEffect(() => {
    const socket = new WebSocket(`ws://localhost:3000/ws/${projectId}`);

    socket.onopen = () => {
      console.log("connected to the server");
    };

    socket.onmessage = (msg) => {
      console.log(msg.data);
      if (msg.data === "deployment successful") {
        toast({ variant: "default", description: `deployment successful` });
        setIsloading(false);
      }
      setLogs((logs) => [...logs, msg.data]);
    };

    return () => socket.close();
  }, []);

  return (
    <main className="bg-[#171717] h-screen w-screen">
      <Navbar />
      <LogTerminal logs={logs} isOpen={isTermOpen} setIsOpen={setIsTermOpen} />
      <div className="flex h-[100vh] w-[100vw] items-center justify-center overflow-hidden">
        <div className="z-[2] flex h-[100%] w-[100%] items-center justify-center">
          <div className="relative -top-10 z-[20] flex w-[90vw] max-w-[500px] flex-col items-center bg-opacity-[70%] bg-blur-[16rem] justify-center rounded-lg bg-[#1d2021] p-5">
            <h1 className="pr-1 text-3xl font-semibold tracking-tight text-white">
              Welcome To
              <span className="bg-gradient-to-b to-gray-500 from-white bg-clip-text pl-2 pr-1 text-3xl font-black tracking-tighter text-transparent">
                Relay
              </span>
            </h1>
            <p className="mx-auto max-w-xs pb-4 pt-4 text-center text-lg font-light text-gray-400">
              Deploy your react apps with just a click
            </p>

            <div className="flex w-[100%] flex-col items-center justify-center p-4">
              {deploymenUrl.length > 0 ? (
                <Link
                  href={`http://${deploymenUrl}`}
                  target="_blank"
                  className="text-[#989898] cursor-pointer text-xl pb-4"
                >
                  Queued deployment at
                  <span className="hover:underline text-[#989898] ml-1 cursor-pointer text-xl">
                    {deploymenUrl}
                  </span>
                </Link>
              ) : (
                <div className="text-white h-7 w-[100%] p-2 underline cursor-pointer text-xl"></div>
              )}
              <Input
                placeholder="Git repo URL"
                className="w-[100%] text-white text-lg"
                style={{ fontSize: "1.2rem" }}
                value={gitUrl}
                onChange={(e) => setGitUrl(e.target.value)}
                readOnly={isLoading}
              />
              <Button
                onClick={() => deployProject()}
                className="my-3 flex w-[100%] items-center font-semibold justify-center text-black rounded-lg bg-white px-2 py-1 text-lg hover:bg-gray-300"
                disabled={isLoading}
              >
                {isLoading ? (
                  <Loader2 size={80} className="animate-spin w-10 h-10" />
                ) : (
                  "Deploy"
                )}
              </Button>

              {isLoading ? (
                <div
                  onClick={() => setIsTermOpen(!isTermOpen)}
                  className="text-white underline cursor-pointer text-xl"
                >
                  View logs
                </div>
              ) : (
                <div className="text-white h-7 w-[100%] underline cursor-pointer text-xl"></div>
              )}
            </div>
          </div>
        </div>
        <div className="absolute bottom-0 z-[10] h-80 w-[40rem] rounded-t-full bg-gradient-to-t from-white to-gray-500 "></div>
      </div>
    </main>
  );
}
