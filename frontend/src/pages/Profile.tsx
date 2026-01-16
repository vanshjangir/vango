import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { UserProfileData } from "../types/game";
import Navbar from "../components/Navbar";
import {
  Flex,
  Box,
  Heading,
  Button,
  Input,
  Image,
  Text,
  Grid,
  TableContainer,
  Table,
  Thead,
  Tbody,
  Tr,
  Td,
  Th,
  Link,
} from "@chakra-ui/react"

const Profile = () => {
  const nav = useNavigate();
  const { username } = useParams();
  const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;
  const [userData, setUserData] = useState<UserProfileData | null>(null);
  const token = localStorage.getItem("token") || "";
  const tokenType = localStorage.getItem("tokenType") || "";
  const [textAreaVis, setTextAreaVis] = useState<boolean>(false);
  const [newUsername, setNewUsername] = useState<string>("");

  const getData = async () => {
    try {
      const response = await fetch(BACKEND_URL + `/profile?username=${username}`, {
        method: "GET",
      });
      const data = await response.json();
      setUserData(data);
    } catch (error) {
      console.error("Failed to fetch user data", error);
    }
  };

  const changeUsername = async () => {
    const response = await fetch(BACKEND_URL + `/changeusername`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `${tokenType} ${token}`,
      },
      body: JSON.stringify({
        "username": username,
        "newusername": newUsername,
      }),
    });

    if(response.status === 200){
      localStorage.setItem('username', newUsername);
    }
    setTextAreaVis(false)
  }

  const onButtonClick = () => {
    setTextAreaVis(true);
  }

  useEffect(() => {
    getData();
  }, []);

if (!userData) {
  return (
    <Flex 
      h="100vh" 
      bg="black"
      flexDir="column" 
      color="white"
    >
      <Flex alignItems="center" justifyContent="center" h="100vh">
        <Box textAlign="center">
          <Text 
            fontSize="xl" 
            fontWeight="500"
            color="gray.400"
          >
            Loading...
          </Text>
        </Box>
      </Flex>
    </Flex>
  );
}

  return (
    <Flex 
      minH="100vh" 
      bg="black"
      flexDir="column" 
      color="white"
    >
      <Navbar />
      <Box 
        maxW="6xl" 
        mx="auto" 
        w="full" 
        px={{ base: 4, sm: 6, lg: 4 }} 
        py={8}
      >
        <Box
          bg="gray.900"
          borderRadius="2px"
          p={{ base: 6, md: 8 }}
          border="1px solid"
          borderColor="gray.800"
          mb={8}
        >
          <Flex alignItems="center" gap="6" mb="8">
            <Box>
              <Flex alignItems="center" gap={4}>
                {textAreaVis === false ? (
                  <>
                    <Heading 
                      as="h1" 
                      size="2xl" 
                      fontWeight="600"
                      color="white"
                    >
                      {userData.name || username}
                    </Heading>
                    <Button
                      onClick={onButtonClick}
                      bg="gray.800"
                      _hover={{
                        bg: "gray.700",
                      }}
                      borderRadius="2px"
                      p={3}
                      border="1px solid"
                      borderColor="gray.700"
                    >
                      <Image src="/editpencil.png" width={"24px"} />
                    </Button>
                  </>
                ) : (
                    <>
                      <Input
                        bg="gray.800"
                        border="1px solid"
                        borderColor="gray.700"
                        _focus={{
                          borderColor: "gray.600",
                        }}
                        placeholder="username"
                        autoFocus
                        onChange={(e) => setNewUsername(e.target.value)}
                        onKeyDown={(e) => {
                          if (e.key === "Enter") changeUsername();
                        }}
                        borderRadius="2px"
                        fontSize="lg"
                        py={3}
                        color="white"
                        _placeholder={{ color: "gray.500" }}
                      />
                      <Button
                        onClick={changeUsername}
                        bg="green.600"
                        color="white"
                        _hover={{
                          bg: "green.500",
                        }}
                        ml="3"
                        borderRadius="2px"
                        p={3}
                        border="1px solid"
                        borderColor="green.500"
                      >
                        <Image src="/tick.png" w="24px" />
                      </Button>
                    </>
                  )}
              </Flex>
              <Text 
                fontSize={"2xl"} 
                color="gray.300"
                fontWeight="500"
                mt={3}
              >
                Rating: {userData.rating}
              </Text>
            </Box>
          </Flex>
        </Box>

        <Box mb="8">
          <Heading 
            as="h2" 
            size="xl" 
            fontWeight="600" 
            mb="6"
            color="white"
          >
            Statistics
          </Heading>
          <Grid
            templateColumns={{ base: "repeat(2, 1fr)", sm: "repeat(4, 1fr)" }}
            gap="6"
          >
            <Box 
              textAlign="center" 
              bg="gray.900"
              p="6" 
              borderRadius="2px"
              border="1px solid"
              borderColor="gray.800"
            >
              <Text display="block" fontSize="3xl" fontWeight="600" color="white">
                {userData.gamesPlayed}
              </Text>
              <Text color="gray.400" fontWeight="500" mt={2}>Games Played</Text>
            </Box>
            <Box 
              textAlign="center" 
              bg="gray.900"
              p="6" 
              borderRadius="2px"
              border="1px solid"
              borderColor="green.600"
            >
              <Text display="block" fontSize="3xl" fontWeight="600" color="white">
                {userData.wins}
              </Text>
              <Text color="gray.400" fontWeight="500" mt={2}>Wins</Text>
            </Box>
            <Box 
              textAlign="center" 
              bg="gray.900"
              p="6" 
              borderRadius="2px"
              border="1px solid"
              borderColor="red.600"
            >
              <Text display="block" fontSize="3xl" fontWeight="600" color="white">
                {userData.losses}
              </Text>
              <Text color="gray.400" fontWeight="500" mt={2}>Losses</Text>
            </Box>
            <Box 
              textAlign="center" 
              bg="gray.900"
              p="6" 
              borderRadius="2px"
              border="1px solid"
              borderColor="gray.800"
            >
              <Text display="block" fontSize="3xl" fontWeight="600" color="white">
                {userData.highestrating}
              </Text>
              <Text color="gray.400" fontWeight="500" mt={2}>Highest Rating</Text>
            </Box>
          </Grid>
        </Box>

        <Box
          bg="gray.900"
          borderRadius="2px"
          p={{ base: 6, md: 8 }}
          border="1px solid"
          borderColor="gray.800"
        >
          <Heading 
            as="h2" 
            size="xl" 
            fontWeight="600" 
            mb="6"
            color="white"
          >
            Recent Games
          </Heading>
          {userData.recentGames && userData.recentGames.length > 0 ? (
            <TableContainer>
              <Table variant="simple">
                <Thead>
                  <Tr>
                    <Th 
                      fontSize={"md"} 
                      color="gray.400" 
                      fontWeight="500"
                      textTransform="none"
                      letterSpacing="normal"
                      borderColor="gray.800"
                    >
                      Result
                    </Th>
                    <Th 
                      fontSize={"md"} 
                      color="gray.400" 
                      fontWeight="500"
                      textTransform="none"
                      letterSpacing="normal"
                      borderColor="gray.800"
                    >
                      Opponent
                    </Th>
                    <Th 
                      fontSize={"md"} 
                      color="gray.400" 
                      fontWeight="500"
                      textTransform="none"
                      letterSpacing="normal"
                      borderColor="gray.800"
                    >
                      Date
                    </Th>
                    <Th 
                      fontSize={"md"} 
                      color="gray.400" 
                      textAlign="center"
                      fontWeight="500"
                      textTransform="none"
                      letterSpacing="normal"
                      borderColor="gray.800"
                    >
                      Review
                    </Th>
                  </Tr>
                </Thead>
                <Tbody fontSize={"md"}>
                  {userData.recentGames.map((game, index) => (
                    <Tr 
                      key={index}
                      borderColor="gray.800"
                      _hover={{
                        bg: "gray.800",
                      }}
                    >
                      <Td 
                        color={game.result === "Lost" ? "red.400" : "green.400"}
                        fontWeight="500"
                        borderColor="gray.800"
                      >
                        {game.result}
                      </Td>
                      <Td borderColor="gray.800">
                        <Link
                          href={"/profile/" + game.opponent}
                          color="gray.300"
                          _hover={{ 
                            textDecoration: "none",
                            color: "white",
                          }}
                          fontWeight="500"
                        >
                          {game.opponent}
                        </Link>
                      </Td>
                      <Td color="gray.400" fontWeight="500" borderColor="gray.800">
                        {
                          new Date(game.created_at.split(" ")[0])
                          .toLocaleDateString('en-US', {
                            day: 'numeric',
                            month: 'short',
                            year: 'numeric'
                          }).replace(/(\d+)/, (day: any) => {
                            const s = ["th", "st", "nd", "rd"];
                            const v = day % 100;
                            return day + (s[(v - 20) % 10] || s[v] || s[0]);
                          })
                        }
                      </Td>
                      <Td textAlign="center" borderColor="gray.800">
                        <Button
                          onClick={() => {
                            nav(`/review/${game.gameid}`);
                          }}
                          size="sm"
                          bg="gray.800"
                          color="white"
                          _hover={{ 
                            bg: "gray.700",
                          }}
                          borderRadius="2px"
                          px={4}
                          py={2}
                          fontWeight="500"
                          border="1px solid"
                          borderColor="gray.700"
                        >
                          Review
                        </Button>
                      </Td>
                    </Tr>
                  ))}
                </Tbody>
              </Table>
            </TableContainer>
          ) : (
              <Text 
                color="gray.500" 
                fontSize="md" 
                textAlign="center"
                py={8}
                fontWeight="500"
              >
                No recent games played.
              </Text>
            )}
        </Box>
      </Box>
    </Flex>
  );
}

export default Profile;
