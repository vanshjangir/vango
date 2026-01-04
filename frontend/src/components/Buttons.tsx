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
      fontWeight="bold"
      h="230px"
      bg="linear-gradient(135deg, #f6ad55, #ed8936)"
      color="white"
      rounded="2px"
      boxShadow="lg"
      _hover={{
        bg: "linear-gradient(135deg, #ed8936, #dd6b20)",
        boxShadow: "xl",
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
      fontWeight="bold"
      py={2}
      marginX={2}
      h="auto"
      bg="linear-gradient(135deg, #f6ad55, #ed8936)"
      color="white"
      rounded="2px"
      boxShadow="md"
      border="2px solid transparent"
      _hover={{
        bg: "linear-gradient(135deg, #ed8936, #dd6b20)",
        boxShadow: "lg",
        borderColor: "orange.300"
      }}
      transition="all 0.2s ease"
    >
      Reconnect
    </Button>
  );
};
