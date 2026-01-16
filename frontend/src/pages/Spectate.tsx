import { useEffect, useRef, useState } from "react";
import { ChatMessage, MsgMove, MsgSync, MsgGameover, MsgStart } from "../types/game";
import { GameState } from "../types/game";
import Navbar from "../components/Navbar";
import {
  cellSize,
  gridSize,
  EMPTY_CELL,
  redrawCanvas,
  decodeState,
  WHITE_CELL,
  BLACK_CELL,
} from "../utils/board";
import {
  Box,
  Flex,
  VStack,
  Text,
  Input,
  Button,
  HStack,
} from "@chakra-ui/react";

const Spectate: React.FC = () => {
  const ctxRef = useRef<CanvasRenderingContext2D | null>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const socketRef = useRef<WebSocket | null>(null);
  const gsRef = useRef<GameState | null>(null);
  const messages = useRef<ChatMessage[]>([]);
  const [pRemTime, setPRemTime] = useState<number>(60000);
  const [opRemTime, setOpRemTime] = useState<number>(60000);
  const [history, setHistory] = useState<string[]>([]);
  const [chatMessages, setChatMessages] = useState<ChatMessage[]>([]);
  const [chatInput, setChatInput] = useState<string>("");
  const tickRef = useRef<boolean>(false);
  const turnRef = useRef<number>(BLACK_CELL);
  const [msg, setMsg] = useState<string>("Starting...");
  const historyBoxRef = useRef<HTMLDivElement>(null);
  const chatBoxRef = useRef<HTMLDivElement>(null);

  let started = false;

  const getGameState = async () => {
    const socket = socketRef.current;
    if (!socket) return;
    socket.send(JSON.stringify({ type: "syncstate", }));
  };

  const updateState = (state: string, move: string) => {
    if (!gsRef.current)
      return;
    
    const gameState = gsRef.current;
    const newMoves = decodeState(gsRef, state);
    
    newMoves.forEach((item) => {
      gameState.state[item.x][item.y] = item.c;
      if (item.c === EMPTY_CELL) {
        return;
      }
      setMsg(move);
    });

    if (!gameState.history) gameState.history = [];
    gameState.history.push(move);
    redrawCanvas(canvasRef, gsRef, ctxRef);
  }

  const handleSocketRecv = async (data: any) => {
    const msg: MsgMove|MsgSync|MsgGameover|MsgStart = await JSON.parse(data);
    switch (msg.type) {
      case "start":
        await afterStart(msg);
        break;

      case "move":
        if (gsRef.current) {
          const gameState = gsRef.current;
          if (msg.move === "ps") {
            gameState.history.push(msg.move);
            setMsg("Pass")
          } else {
            updateState(msg.state, msg.move);
          }
          setHistory([...gameState.history]);
          turnRef.current = 1 - turnRef.current;
        }
        break;

      case "syncstate":
        if (gsRef.current) {
          const gs = gsRef.current;
          gs.pname = gs.color === WHITE_CELL ? msg.whitename : msg.blackname;
          gs.opname = gs.color === WHITE_CELL ? msg.blackname : msg.whitename;
          updateState(msg.state, "");

          gs.history = msg.history;
          setHistory([...msg.history]);
          setPRemTime(gs.color === WHITE_CELL ? msg.whiteRemTime : msg.blackRemTime);
          setOpRemTime(gs.color === WHITE_CELL ? msg.blackRemTime : msg.whiteRemTime);
          redrawCanvas(canvasRef, gsRef, ctxRef);
        }
        break;

      case "gameover":
        tickRef.current = false;
        setMsg(msg.winner === gsRef.current?.color ? "You won" : "You lost");
        destSocket();
        break;
    }
  };
  
  const destSocket = async () => {
    if (socketRef.current) {
      socketRef.current.close();
      socketRef.current = null;
    }
  }

  const setupSocket = async () => {
    const wsurl = localStorage.getItem('wsurl') ?? '';
    const token = localStorage.getItem('token');
    const tokenType = localStorage.getItem('tokenType');
    socketRef.current = new WebSocket(`${wsurl}/play`, `${tokenType}.${token}`);
    socketRef.current.onmessage = async (event: MessageEvent) => {
      const data = await event.data;
      await handleSocketRecv(data);
    };
  };

  const afterStart = async (msg: MsgStart) => {
    await getGameState();
    if (gsRef.current) gsRef.current.color = msg.color;
    tickRef.current = true;
    setMsg("Started")
    const url = new URL(window.location.href);
    url.searchParams.set("id", `${msg.gameid}`);
    window.history.replaceState({}, "", url);
    tickClock();
  }

  const updateClock = async () => {
    if (!gsRef.current) return;
    if (gsRef.current.color == turnRef.current) {
      setPRemTime(t => t - 1000);
    } else {
      setOpRemTime(t => t - 1000);
    }
  }

  const tickClock = async () => {
    setTimeout(() => {
      if (!tickRef.current) return;
      updateClock();
      tickClock();
    }, 1*1000);
  }

  const formatTime = (ms: number) => {
    const totalSec = Math.max(0, Math.floor(ms / 1000));
    const m = Math.floor(totalSec / 60);
    const s = totalSec % 60;
    return `${m}:${s.toString().padStart(2, "0")}`;
  };

  useEffect(() => {
    redrawCanvas(canvasRef, gsRef, ctxRef);
  }, [gsRef.current]);


  const startSetup = async () => {
    if (started) return;
    started = true;

    gsRef.current = {
      gameId: "",
      pname: "",
      opname: "",
      color: EMPTY_CELL,
      state: Array.from({length: 19}, () => new Array(19).fill(EMPTY_CELL)),
      history: [],
    }
    
    const resizeHandler = () => {
      redrawCanvas(canvasRef, gsRef, ctxRef);
    };

    window.addEventListener("resize", resizeHandler);
    redrawCanvas(canvasRef, gsRef, ctxRef);
    
    await setupSocket();
  }

  useEffect(() => {
    if (historyBoxRef.current) {
      historyBoxRef.current.scrollTop = historyBoxRef.current.scrollHeight
    }
  })

  useEffect(() => {
    if (chatBoxRef.current) {
      chatBoxRef.current.scrollTop = chatBoxRef.current.scrollHeight;
    }
  }, [chatMessages]);

  const sendChatMessage = () => {
    const socket = socketRef.current;
    const text = chatInput.trim();
    if (!socket || !text) return;
    socket.send(JSON.stringify({ type: "chat", text: text }));
    messages.current = [
      ...messages.current,
      { type: "sent", text }
    ];
    setChatMessages([...messages.current]);
    setChatInput("");
  };

  const sendPassMove = () => {
    const socket = socketRef.current;
    if (!socket) return;
    socket.send(JSON.stringify({ type: "move", move: "ps" }));
  };

  const sendAbort = () => {
    const socket = socketRef.current;
    if (!socket) return;
    socket.send(JSON.stringify({ type: "abort" }));
  };

  useEffect(() => {
    startSetup();
  }, []);

  return (
    <Box 
      minH="100vh" 
      display="flex" 
      flexDirection="column" 
      bg="black"
      overflowY="auto"
    >
      <Navbar />
      <Flex 
        id="game-container" 
        w="100%" 
        flexDirection={{ base: "column", lg: "row", sm: "column" }} 
        alignItems={{ base: "center", lg: "flex-start" }} 
        justifyContent="center"
        gap={8}
        p={6}
      >
        <VStack 
          id="game-board" 
          spacing={6} 
          align="center"
          order={{ base: 2, lg: 2 }}
        >
          <canvas
            ref={canvasRef}
            id="canvas"
            width={gridSize + 2 * cellSize}
            height={gridSize + 2 * cellSize}
            style={{
              backgroundColor: "#fef3c7",
              borderRadius: "2px",
              border: "1px solid #1f2937"
            }}
          />
        </VStack>
        <VStack 
          id="game-board" 
          spacing={6} 
          align="center"
          order={{ base: 2, lg: 2 }}
        >
          <Box
            px={4}
            py={3}
            borderRadius="2px"
            bg={gsRef.current?.color === BLACK_CELL ? "gray.900" : "gray.100"}
            color={gsRef.current?.color === BLACK_CELL ? "white" : "black"}
            border="4px solid"
            borderColor={gsRef.current?.color === turnRef.current ? "green.500" : "gray.700"}
            minW="320px"
            textAlign="center"
          >
            <Box fontSize="sm" opacity={0.7}>
              {gsRef.current?.opname ? gsRef.current?.opname : "Loading..."}
              {gsRef.current?.color !== turnRef.current ? "(turn)" : ""}
            </Box>
            <Box fontSize="lg" fontWeight="600">{formatTime(opRemTime)}</Box>
            <Box fontSize="lg" fontWeight="600">{formatTime(pRemTime)}</Box>
            <Box fontSize="sm" opacity={0.7}>
              {gsRef.current?.opname ? gsRef.current?.opname : "Loading..."}
              {gsRef.current?.color === turnRef.current ? "(turn)" : ""}
            </Box>
          </Box>
          <Box
            fontSize={"4xl"}
            color={"white"}
            minW="320px"
            minH="54px"
            textAlign="center"
            fontWeight="500"
          >
            {msg}
          </Box>
          <HStack w="320px" spacing={3}>
            <Button
              flex={1}
              size="sm"
              onClick={sendPassMove}
              bg="green.600"
              color="white"
              rounded={"2px"}
              _hover={{ bg: "green.500" }}
            >
              Pass
            </Button>
            <Button
              flex={1}
              size="sm"
              onClick={sendAbort}
              bg="red.600"
              color="white"
              rounded={"2px"}
              _hover={{ bg: "red.500" }}
            >
              Abort
            </Button>
          </HStack>
          <Box
            w="320px"
            h="220px"
            overflowY="auto"
            bg="gray.900"
            ref={historyBoxRef}
            borderRadius="2px"
            border="1px solid"
            borderColor="gray.800"
            css={{
              "&::-webkit-scrollbar": { width: "6px" },
              "&::-webkit-scrollbar-thumb": { background: "#4b5563", borderRadius: "4px" },
            }}
          >
            <VStack align="stretch" spacing={0}>
              {Array.from({ length: Math.ceil(history.length / 2) }).map((_, rowIndex) => {
                const moveIndex = rowIndex * 2;
                const blackMove = history[moveIndex];
                const whiteMove = history[moveIndex + 1];
                
                return (
                  <Flex key={rowIndex}>
                    <Box
                      flex={1}
                      px={3}
                      py={2}
                      bg={blackMove ? "gray.800" : "gray.900"}
                      textAlign="center"
                      display="flex"
                      alignItems="center"
                      justifyContent="center"
                    >
                      <Text
                        color={blackMove ? "white" : "gray.600"}
                        fontWeight="500"
                        fontSize="sm"
                      >
                        {moveIndex + 1}. {blackMove === "ps" ? "Pass" : blackMove || ""}
                      </Text>
                    </Box>
                    <Box
                      flex={1}
                      px={3}
                      py={2}
                      bg={whiteMove ? "gray.100" : "gray.900"}
                      textAlign="center"
                      display="flex"
                      alignItems="center"
                      justifyContent="center"
                    >
                      <Text
                        color={whiteMove ? "black" : "gray.600"}
                        fontWeight="500"
                        fontSize="sm"
                      >
                        {whiteMove ? `${moveIndex + 2}. ${whiteMove === "ps" ? "Pass" : whiteMove}` : ""}
                      </Text>
                    </Box>
                  </Flex>
                );
              })}
              {history.length === 0 && (
                <Text color="gray.500" textAlign="center" py={4} fontSize="sm">
                  No moves yet
                </Text>
              )}
            </VStack>
          </Box>
          <Box
            w="320px"
            h="228px"
            bg="gray.900"
            borderRadius="2px"
            border="1px solid"
            borderColor="gray.800"
            overflow="hidden"
            display="flex"
            flexDirection="column"
            gap={3}
            p={3}
          >
            <Text fontSize="md" fontWeight="600" color="white" mb={1}>
              Chat
            </Text>
            <Box
              ref={chatBoxRef}
              flex="1"
              overflowY="auto"
              css={{
                "&::-webkit-scrollbar": { width: "6px" },
                "&::-webkit-scrollbar-thumb": { background: "#4b5563", borderRadius: "4px" },
              }}
            >
              <VStack align="stretch" spacing={2}>
                {chatMessages.map((m, idx) => (
                  <Flex key={idx} justify={m.type === "sent" ? "flex-end" : "flex-start"}>
                    <Box
                      maxW="80%"
                      px={3}
                      py={2}
                      borderRadius="md"
                      color={m.type === "sent" ? "white" : "gray.400"}
                      fontSize="sm"
                      boxShadow="sm"
                    >
                      {m.text}
                    </Box>
                  </Flex>
                ))}
                {chatMessages.length === 0 && (
                  <Text color="gray.500" fontSize="sm" textAlign="center">
                    No messages yet
                  </Text>
                )}
              </VStack>
            </Box>
            <HStack>
              <Input
                value={chatInput}
                onChange={(e) => setChatInput(e.target.value)}
                placeholder="Type a message..."
                bg="gray.800"
                borderColor="gray.700"
                color="white"
                _placeholder={{ color: "gray.500" }}
              />
              <Button
                onClick={sendChatMessage}
                colorScheme="blue"
              >
                Send
              </Button>
            </HStack>
          </Box>
        </VStack>
      </Flex>
    </Box>
  );
};

export default Spectate;
