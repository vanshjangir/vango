import { Button } from "@chakra-ui/react";

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
      h="230px"
      bg="white"
      color="gray.900"
      rounded="2px"
      border="1px solid"
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
