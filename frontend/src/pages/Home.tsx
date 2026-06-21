import { useState } from "react";
import { useNavigate } from "react-router-dom";
import Navbar from "../components/Navbar";
import { PlayButton, SpectatorButton } from "../components/Buttons";
import { Flex, Box, Text, Image, VStack, HStack } from "@chakra-ui/react";

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

      <Flex flex="1" direction="column">
        {/* Main content area */}
        <Box maxW="6xl" mx="auto" w="full" px={{ base: 4, sm: 6, lg: 4 }} py={{ base: 8, md: 0 }} flex="1" display="flex" alignItems="center">
          <Flex
            w="full"
            direction={{ base: "column", md: "row" }}
            align={{ base: "stretch", md: "flex-end" }}
            gap={0}
          >
            {/* Left column — text + buttons */}
            <VStack
              w={{ base: "100%", md: "50%" }}
              align="stretch"
              spacing={0}
              pr={{ base: 0, md: 8 }}
              justify="flex-end"
            >
              {/* Large headline block */}
              <Box mb={6}>
                <Text
                  fontSize={{ base: "5xl", md: "7xl" }}
                  fontWeight="700"
                  letterSpacing="tighter"
                  lineHeight="0.95"
                >
                  Play
                </Text>
                <Text
                  fontSize={{ base: "5xl", md: "7xl" }}
                  fontWeight="700"
                  letterSpacing="tighter"
                  lineHeight="0.95"
                >
                  Go,
                </Text>
                <Text
                  fontSize={{ base: "5xl", md: "7xl" }}
                  fontWeight="700"
                  letterSpacing="tighter"
                  lineHeight="0.95"
                  color="gray.500"
                >
                  simply.
                </Text>
              </Box>

              {/* Description */}
              <Text fontSize="sm" color="gray.500" lineHeight="1.7" mb={5}>
                Fast online games, no clutter. Just you,
                <br />
                your opponent, and the board.
              </Text>

              {/* Feature tags */}
              <HStack spacing={4} mb={6}>
                <Text fontSize="xs" color="gray.600" fontWeight="500" letterSpacing="widest" textTransform="uppercase">
                  9×9
                </Text>
                <Text fontSize="xs" color="gray.700">—</Text>
                <Text fontSize="xs" color="gray.600" fontWeight="500" letterSpacing="widest" textTransform="uppercase">
                  13×13
                </Text>
                <Text fontSize="xs" color="gray.700">—</Text>
                <Text fontSize="xs" color="gray.600" fontWeight="500" letterSpacing="widest" textTransform="uppercase">
                  19×19
                </Text>
              </HStack>

              {/* Status */}
              <Box mb={4}>
                {matchStatus === "pending" ? (
                  <HStack spacing={3} align="center">
                    <HStack spacing={1}>
                      <Box w="6px" h="6px" bg="gray.600" />
                      <Box w="6px" h="6px" bg="gray.500" />
                      <Box w="6px" h="6px" bg="gray.400" />
                    </HStack>
                    <Text fontSize="xs" color="gray.500" fontWeight="500" letterSpacing="wide" textTransform="uppercase">
                      Finding opponent
                    </Text>
                  </HStack>
                ) : (
                  <Text fontSize="xs" color="gray.600" fontWeight="500" letterSpacing="wide" textTransform="uppercase">
                    {token ? "Ready to play" : "Log in to save games & rating"}
                  </Text>
                )}
              </Box>

              {/* Action buttons */}
              <VStack spacing={3} align="stretch">
                <PlayButton
                  label={"Play Online"}
                  handler={findGame}
                />
                <SpectatorButton
                  label={"Spectate"}
                  handler={spectateGame}
                />
              </VStack>
            </VStack>

            {/* Right column — board image */}
            <Box
              w={{ base: "100%", md: "50%" }}
              display={{ base: "none", md: "block" }}
              overflow="hidden"
            >
              <Image
                src="/boardbg.png"
                alt="Go board"
                objectFit="cover"
                w="100%"
                h="100%"
                maxH="520px"
              />
            </Box>
          </Flex>
        </Box>

        {/* Footer line */}
        <Box maxW="6xl" mx="auto" w="full" px={{ base: 4, sm: 6, lg: 4 }} pb={6}>
          <Box borderTop="1px solid" borderColor="gray.900" pt={4}>
            <Text fontSize="xs" color="gray.700" letterSpacing="wide">
              VANGO — Open source Go
            </Text>
          </Box>
        </Box>
      </Flex>
    </Flex>
  );
};

export default Home;
