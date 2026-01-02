import { useEffect, useRef } from "react";
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
} from "../utils/board";
import {
  Box,
  Flex,
  VStack,
} from "@chakra-ui/react";

const Game: React.FC = () => {
  const ctxRef = useRef<CanvasRenderingContext2D | null>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const socketRef = useRef<WebSocket | null>(null);
  const gameStateRef = useRef<GameState | null>(null);
  const msgRef = useRef<HTMLDivElement>(null);
  const messages = useRef<ChatMessage[]>([]);

  let pRemTime = 900;
  let opRemTime = 900;
  let started = false;

  const getGameState = async () => {
    const socket = socketRef.current;
    if (!socket) return;
    socket.send(JSON.stringify({ type: "syncstate", }));
  };

  const updateState = (state: string, move: string) => {
    if (!gameStateRef.current)
      return;
    
    const gameState = gameStateRef.current;
    const newMoves = decodeState(gameStateRef, state);
    
    newMoves.forEach((item) => {
      gameState.state[item.x][item.y] = item.c;
      if (item.c === EMPTY_CELL) {
        return;
      }
      showMoveStatus(move);
    });

    if (!gameState.history) gameState.history = [];
    gameState.history.push(move);
    redrawCanvas(canvasRef, gameStateRef, ctxRef);
  }

  const handleSocketRecv = async (data: any) => {
    const msg: MsgMove|MsgMoveStatus|MsgSync|MsgGameover|MsgChat|MsgStart = await JSON.parse(data);
    switch (msg.type) {
      case "start":
        await getGameState();
        break;
      
      case "movestatus":
        if (msg.code !== "VALID") {
          showMoveStatus(msg.code);
          redrawCanvas(canvasRef, gameStateRef, ctxRef);
          break;
        }
        if (gameStateRef.current) {
          const gameState = gameStateRef.current;
          if (msg.move === "ps") {
            if (!gameState.history) gameState.history = [];
            gameState.history.push(msg.move);
            showMoveStatus("Pass");
          } else {
            updateState(msg.state, msg.move);
          }
          pRemTime = 900 - Math.round(msg.pRemTime / 1000);
          opRemTime = 900 - Math.round(msg.opRemTime / 1000);
        }
        break;

      case "move":
        if (gameStateRef.current) {
          const gameState = gameStateRef.current;
          if (msg.move === "ps") {
            gameState.history.push(msg.move);
            showMoveStatus("Pass");
          } else {
            updateState(msg.state, msg.move);
          }
        }
        break;

      case "syncstate":
        if (gameStateRef.current) {
          const gameState = gameStateRef.current;
          gameState.pname = gameState.color === WHITE_CELL ? msg.whitename : msg.blackname;
          gameState.opname = gameState.color === WHITE_CELL ? msg.blackname : msg.whitename;

          updateState(msg.state, "");

          gameState.history = msg.history;
          pRemTime = gameState.color === WHITE_CELL ? msg.whiteRemTime : msg.blackRemTime;
          opRemTime = gameState.color === WHITE_CELL ? msg.blackRemTime : msg.whiteRemTime;
          redrawCanvas(canvasRef, gameStateRef, ctxRef);
        }
        break;

      case "gameover":
        showEndMessage(msg.winner === gameStateRef.current?.color ? "You won" : "You lost");
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

  const showMoveStatus = (msg: string) => {
    if (msgRef.current) {
      msgRef.current.className = "text-center text-xl font-bold h-[40px]";
      if (!isNaN(Number(msg.slice(1)))) {
        msgRef.current.innerText = toShow(msg);
      } else {
        msgRef.current.innerText = msg === "ps" ? "Pass" : msg;
      }
    }
  };

  const showEndMessage = (msg: string) => {
    if (msgRef.current) {
      msgRef.current.className = "text-center text-3xl font-bold h-[40px]";
      msgRef.current.innerText = msg;
    }
  };

  const toShow = (move: string) => {
    if (move) {
      return move === "ps" ? "Pass" :
        move[0].toUpperCase() + String(Number(move.slice(1)) + 1);
    }
    return "";
  };

  useEffect(() => {
    redrawCanvas(canvasRef, gameStateRef, ctxRef);
  }, [gameStateRef.current]);


  const startSetup = async () => {
    if (started) return;
    started = true;

    gameStateRef.current = {
      gameId: "",
      pname: "",
      opname: "",
      color: EMPTY_CELL,
      state: Array.from({length: 19}, () => new Array(19).fill(EMPTY_CELL)),
      history: [],
    }
    
    const clickHandler = (e: MouseEvent) => {
      handleCanvasClick(canvasRef, gameStateRef, ctxRef, socketRef, e);
    };
    const resizeHandler = () => {
      redrawCanvas(canvasRef, gameStateRef, ctxRef);
    };

    canvasRef.current?.addEventListener("click", clickHandler);
    window.addEventListener("resize", resizeHandler);
    redrawCanvas(canvasRef, gameStateRef, ctxRef);
    
    await setupSocket();
  }

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
      <Box
        position="absolute"
        top="0"
        left="0"
        right="0"
        bottom="0"
        opacity="0.05"
        backgroundSize="100px 100px"
      />

      <Navbar />
      <Flex 
        id="game-container" 
        w="100%" 
        flexDirection={{ base: "column", lg: "row" }} 
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
          <Box position="relative">
            <Box
              position="absolute"
              top="-15px"
              left="-15px"
              right="-15px"
              bottom="-15px"
              borderRadius="3xl"
              bg="radial-gradient(circle, rgba(246, 173, 85, 0.15), rgba(237, 137, 54, 0.05))"
              filter="blur(20px)"
            />
            <canvas
              ref={canvasRef}
              id="canvas"
              width={gridSize + 2 * cellSize}
              height={gridSize + 2 * cellSize}
              style={{
                position: "relative",
                zIndex: 2,
                backgroundColor: "#fef3c7",
                borderRadius: "24px",
                border: "3px solid rgba(246, 173, 85, 0.3)",
                boxShadow: "0 25px 50px rgba(0, 0, 0, 0.5)"
              }}
            />
          </Box>
        </VStack>
      </Flex>
    </Box>
  );
};

export default Game;
