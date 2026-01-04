import { useState } from "react";
import { useNavigate } from "react-router-dom";
import Navbar from "../components/Navbar";
import { PlayButton } from "../components/Buttons"
import { Flex, Box, Text, Image, VStack, HStack, Container } from "@chakra-ui/react";

const Home = () => {
  const nav = useNavigate();
  const [matchStatus, setMatchStatus] = useState("");
  const token = localStorage.getItem("token") || "";
  const tokenType = localStorage.getItem("tokenType") || "";
  const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;
  let finding = false;

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
      bg="linear-gradient(135deg, #1a202c 0%, #2d3748 25%, #4a5568 50%, #2d3748 75%, #1a202c 100%)"
      direction="column" 
      color="white"
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
      <Container maxW="1300px" flex="1" mt={"160px"} py={6}>
        <HStack 
          justifyContent="center" 
          spacing={{ base: 0, lg: 4 }} 
          alignItems="center" 
          w="full" 
          h="full"
        >
          {/* Board Section */}
          <Box display={{ base: "none", lg: "block" }}
            h={{ base: "auto", lg: "542px" }} 
            w={{ base: "auto", lg: "542px" }} 
          >
            <Box position="relative">
              <Image 
                src="/boardbg.png" 
                boxShadow="0 20px 40px rgba(0, 0, 0, 0.4)" 
                w="542px"
                h="542px"
                rounded="4px"
              />
            </Box>
          </Box>

          <VStack 
            spacing={4} 
            align="center" 
            justify="center" 
            h={{ base: "auto", lg: "542px" }} 
            w={{ base: "auto", lg: "542px" }} 
          >
            <Box
              rounded="2px"
              h="60px"
              w="full"
            >
              <Box
                position="absolute"
                top="0"
                left="0"
                right="0"
                bg="linear-gradient(90deg, #f6ad55, #ed8936, #dd6b20, #c05621)"
              />

              {matchStatus === "pending" ? (
                <VStack spacing={4}>
                  <HStack spacing={2} justify="center">
                    <Box w="12px" h="12px" bg="orange.400" borderRadius="full"/>
                    <Box w="12px" h="12px" bg="orange.300" borderRadius="full"/>
                    <Box w="12px" h="12px" bg="orange.200" borderRadius="full"/>
                  </HStack>
                  <Text 
                    fontSize="xl" 
                    color="orange.200"
                    fontWeight="600"
                    textAlign="center"
                  >
                    Finding worthy opponent...
                  </Text>
                </VStack>
              ) : (
                  <VStack spacing={3}>
                    <Box
                      w="8px"
                      h="8px"
                      bg="green.400"
                      borderRadius="full"
                      boxShadow="0 0 10px rgba(72, 187, 120, 0.6)"
                    />
                    <Text 
                      fontSize="lg"
                      color="gray.200"
                      fontWeight="600"
                      textAlign="center"
                    >
                      {token ? "Ready for play" : "Log In to Play"}
                    </Text>
                  </VStack>
                )}
            </Box>
            <PlayButton
              label={"Play Online"}
              handler={findGame}
            />
            <PlayButton
              label={"Spectate"}
              handler={findGame}
            />
          </VStack>
        </HStack>
      </Container>
    </Flex>
  );
};

export default Home;
