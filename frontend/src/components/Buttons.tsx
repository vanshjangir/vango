import { Box, Button, Input } from "@chakra-ui/react";
import { useState } from "react";

export const PlayButton = (
  {label, handler}:{
    label: string;
    handler: () => void;
  }
) => {
  return (
    <Button
      onClick={() => handler()}
      w="full"
      size="lg"
      fontSize="md"
      fontWeight="500"
      letterSpacing="wide"
      textTransform="uppercase"
      h="56px"
      bg="white"
      color="black"
      rounded="0px"
      border="none"
      _hover={{
        bg: "gray.200",
      }}
      transition="background 0.15s ease"
    >
      {label}
    </Button>
  );
};

export const SpectatorButton = (
  {label, handler}:{
    label: string;
    handler: (gameId: string) => void;
  }
) => {
  const [gameId, setGameId] = useState<string>("");
  return (
    <Box
      w="full"
      display="flex"
      flexDirection="row"
      alignItems="stretch"
    >
      <Input
        flex="1"
        rounded="0px"
        border="1px solid"
        borderColor="gray.800"
        bg="transparent"
        color="white"
        h="56px"
        fontSize="sm"
        px={4}
        placeholder="Game ID or link"
        _placeholder={{ color: "gray.600" }}
        _hover={{ borderColor: "gray.700" }}
        _focus={{ borderColor: "gray.600", boxShadow: "none" }}
        onChange={(e) => {
          const value = e.target.value;
          let id = "";
          for (let i = value.length-1; i >= 0; i--) {
            if (!isNaN(Number(value[i]))) id = value[i] + id;
            else break;
          }
          setGameId(id)
        }}
      />
      <Button
        onClick={() => handler(gameId)}
        h="56px"
        px={6}
        fontSize="md"
        fontWeight="500"
        letterSpacing="wide"
        textTransform="uppercase"
        rounded="0px"
        bg="transparent"
        color="gray.400"
        border="1px solid"
        borderColor="gray.800"
        borderLeft="none"
        _hover={{
          color: "white",
          borderColor: "gray.700",
        }}
        transition="all 0.15s ease"
      >
        {label}
      </Button>
    </Box>
  );
};

export const ReconnectButton = (
  {handler}:{
    handler: () => void;
  }
) => {
  return (
    <Button
      onClick={handler}
      size="md"
      fontSize="lg"
      fontWeight="500"
      py={2}
      marginX={2}
      h="auto"
      bg="gray.800"
      color="white"
      rounded="0px"
      border="1px solid"
      borderColor="gray.700"
      _hover={{
        borderColor: "gray.600"
      }}
    >
      Reconnect
    </Button>
  );
};
