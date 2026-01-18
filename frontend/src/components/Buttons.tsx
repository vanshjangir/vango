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
      fontSize="3xl"
      fontWeight="600"
      h="200px"
      bg="white"
      color="gray.900"
      rounded="2px"
      borderColor="gray.300"
      _hover={{
        bg: "gray.100",
        borderColor: "gray.400"
      }}
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
      flexDirection="column"
      alignItems="center"
      justifyContent="center"
    >
      <Input
        w="full"
        rounded="2px 2px 0px 0px"
        border="1px solid"
        textAlign="center"
        placeholder="Game Id or Link"
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
        w="full"
        h="200px"
        fontSize="3xl"
        fontWeight="600"
        rounded="0px 0px 2px 2px"
        bg="white"
        color="gray.900"
        borderColor="gray.300"
        _hover={{
          bg: "gray.100",
          borderColor: "gray.400"
        }}
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
      rounded="2px"
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
