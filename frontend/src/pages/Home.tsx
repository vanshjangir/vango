import { useState } from "react";
import { useNavigate } from "react-router-dom";
import Navbar from "../components/Navbar";
import { PlayButton, SpectatorButton } from "../components/Buttons";
import { Flex, Box, Text, Image, VStack, HStack, Container } from "@chakra-ui/react";

const Home = () => {
  const nav = useNavigate();
  const [matchStatus, setMatchStatus] = useState("");
  const token = localStorage.getItem("token") || "";
  const tokenType = localStorage.getItem("tokenType") || "";
  const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;
  let finding = false;

  const spectateGame = async (gameId: string) => {
    nav(`/spectate/${gameId}`);
  }

  const findGame = async () => {
    if (finding) return;
    finding = true;
    setMatchStatus("pending");
    const response = await fetch(BACKEND_URL + "/findgame", {
      headers: {
        "Authorization": `${tokenType} ${token}`,
      },
    });
    const json = await response.json();
    if (response.status === 200) {
      const wsurl = json.wsurl;
      localStorage.setItem('wsurl', wsurl);
      nav('/game')
    } else {
    setMatchStatus("");
      console.log(`Error occured while finding a game ${json}`);
    }
  };

  return (
    <Flex
      minH="100vh"
      bg="black"
      direction="column"
      color="white"
    >
      <Navbar />
      <Flex flex="1" align="center" justify="center" px={6} py={10}>
        <Container maxW="6xl" mx="auto">
          <Box>
            <Text fontSize="4xl" fontWeight="600" letterSpacing="tight">
              Play Go, simply.
            </Text>
            <Text mt={3} fontSize="md" color="gray.400">
              Fast online games, no clutter. Just you, your opponent,
              and the board.
            </Text>
            {matchStatus === "pending" ? (
              <HStack spacing={3} align="center">
                <HStack spacing={1}>
                  <Box w="8px" h="8px" bg="gray.500" borderRadius="full" />
                  <Box w="8px" h="8px" bg="gray.400" borderRadius="full" />
                  <Box w="8px" h="8px" bg="gray.300" borderRadius="full" />
                </HStack>
                <Text fontSize="sm" color="gray.400">
                  Finding an opponent...
                </Text>
              </HStack>
            ) : (
              <Text fontSize="sm" color="gray.500">
                {token ? "You’re ready to play." : "Log in to save your games and rating."}
              </Text>
            )}
          </Box>

          <HStack
            spacing={{ base: 10, md: 16 }}
            marginTop={"24px"}
            align="flex-start"
            gap={4}
            justify="space-between"
            flexDir={{ base: "column", md: "row" }}
          >
            <VStack w={{ base: "100%", md: "50%" }} spacing={4}>
              <PlayButton
                label={"Play Online"}
                handler={findGame}
              />
              <SpectatorButton
                label={"Spectate"}
                handler={spectateGame}
              />
            </VStack>

            <Box
              w={{ base: "100%", md: "50%" }}
              h={{ base: "100%", md: "50%" }}
              borderRadius="2px"
              overflow="hidden"
            >
              <Image
                src="/boardbg.png"
                alt="Go board"
                objectFit="cover"
                w="100%"
                h="100%"
              />
            </Box>
          </HStack>
        </Container>
      </Flex>
    </Flex>
  );
};

export default Home;
