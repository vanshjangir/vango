import { useEffect, useRef, useState } from "react";
import { ChatMessage, MsgChat, MsgMove, MsgMoveStatus, MsgSync, MsgGameover, MsgStart } from "../types/game";
import { GameState } from "../types/game";
import Navbar from "../components/Navbar";
import {
  cellSize,
  gridSize,
  EMPTY_CELL,
  handleCanvasClick,
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
} from "@chakra-ui/react";

const Game: React.FC = () => {
  const ctxRef = useRef<CanvasRenderingContext2D | null>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const socketRef = useRef<WebSocket | null>(null);
  const gsRef = useRef<GameState | null>(null);
  const messages = useRef<ChatMessage[]>([]);
  const [pRemTime, setPRemTime] = useState<number>(60000);
  const [opRemTime, setOpRemTime] = useState<number>(60000);
  const [history, setHistory] = useState<string[]>([]);
  const tickRef = useRef<boolean>(false);
  const turnRef = useRef<number>(BLACK_CELL);
  const [msg, setMsg] = useState<string>("Starting...");
  const historyBoxRef = useRef<HTMLDivElement>(null);

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
    const msg: MsgMove|MsgMoveStatus|MsgSync|MsgGameover|MsgChat|MsgStart = await JSON.parse(data);
    switch (msg.type) {
      case "start":
        await afterStart(msg);
        break;
      
      case "movestatus":
        if (msg.code !== "VALID") {
          setMsg(msg.code);
          redrawCanvas(canvasRef, gsRef, ctxRef);
          break;
        }
        if (gsRef.current) {
          const gameState = gsRef.current;
          if (msg.move === "ps") {
            gameState.history.push(msg.move);
            setMsg("Pass")
          } else {
            updateState(msg.state, msg.move);
          }
          setHistory([...gameState.history]);
          setPRemTime(gameState.color == BLACK_CELL ? msg.blackRemTime : msg.whiteRemTime);
          setOpRemTime(gameState.color == WHITE_CELL ? msg.blackRemTime : msg.whiteRemTime);
          turnRef.current = 1 - turnRef.current;
        }
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

      case "chat":
        messages.current = [
          ...messages.current,
          { type: 'received', text: msg.message.trim() }
        ];
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
    
    const clickHandler = (e: MouseEvent) => {
      handleCanvasClick(canvasRef, gsRef, ctxRef, socketRef, e);
    };
    const resizeHandler = () => {
      redrawCanvas(canvasRef, gsRef, ctxRef);
    };

    canvasRef.current?.addEventListener("click", clickHandler);
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
    startSetup();
  }, []);

  return (
    <Box 
      minH="100vh" 
      display="flex" 
      flexDirection="column" 
      bg="linear-gradient(135deg, #1a202c 0%, #2d3748 25%, #4a5568 50%, #2d3748 75%, #1a202c 100%)"
      overflowY="auto"
      position="relative"
      overflow="hidden"
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
        position="relative"
        zIndex={1}
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
              position: "relative",
              zIndex: 2,
              backgroundColor: "#fef3c7",
              borderRadius: "4px",
              border: "3px solid rgba(246, 173, 85, 0.3)",
              boxShadow: "0 25px 50px rgba(0, 0, 0, 0.5)"
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
            borderRadius="md"
            bg={gsRef.current?.color === BLACK_CELL ? "gray.900" : "gray.100"}
            color={gsRef.current?.color === BLACK_CELL ? "white" : "black"}
            border={gsRef.current?.color === turnRef.current ? "4px solid" : "4px dashed"}
            borderColor={gsRef.current?.color === turnRef.current ? "green.400" : "gray.400"}
            minW="320px"
            textAlign="center"
          >
            <Box fontSize="sm" opacity={0.7}>{gsRef.current?.opname}</Box>
            <Box fontSize="lg" fontWeight="bold">{formatTime(opRemTime)}</Box>
            <Box fontSize="lg" fontWeight="bold">{formatTime(pRemTime)}</Box>
            <Box fontSize="sm" opacity={0.7}>
              {gsRef.current?.pname} {gsRef.current?.color === turnRef.current ? "(Your turn)" : ""}
            </Box>
          </Box>
          <Box
            fontSize={"4xl"}
            color={"white"}
            minW="320px"
            textAlign="center"
          >
            {msg}
          </Box>
          <Box
            w="320px"
            h="300px"
            overflowY="auto"
            bg="gray.800"
            ref={historyBoxRef}
            borderRadius="4px"
            css={{
              "&::-webkit-scrollbar": { display: "none" },
              scrollbarWidth: "none"
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
                      bg={blackMove ? "gray.900" : "gray.800"}
                      textAlign="center"
                      display="flex"
                      alignItems="center"
                      justifyContent="center"
                    >
                      <Text
                        color={blackMove ? "white" : "gray.500"}
                        fontWeight="semibold"
                        fontSize="sm"
                      >
                        {moveIndex + 1}. {blackMove === "ps" ? "Pass" : blackMove || ""}
                      </Text>
                    </Box>
                    <Box
                      flex={1}
                      px={3}
                      py={2}
                      bg={whiteMove ? "gray.100" : "gray.800"}
                      textAlign="center"
                      display="flex"
                      alignItems="center"
                      justifyContent="center"
                    >
                      <Text
                        color={whiteMove ? "black" : "gray.500"}
                        fontWeight="semibold"
                        fontSize="sm"
                      >
                        {whiteMove ? `${moveIndex + 2}. ${whiteMove === "ps" ? "Pass" : whiteMove}` : ""}
                      </Text>
                    </Box>
                  </Flex>
                );
              })}
              {history.length === 0 && (
                <Text color="gray.400" textAlign="center" py={4} fontSize="sm">
                  No moves yet
                </Text>
              )}
            </VStack>
          </Box>
        </VStack>
      </Flex>
    </Box>
  );
};

export default Game;
