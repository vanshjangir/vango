import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  Flex,
  Box,
  Text,
  Button,
  Link,
  IconButton,
  VStack,
  HStack,
  useDisclosure,
  Collapse
} from "@chakra-ui/react";
import { HamburgerIcon, CloseIcon } from "@chakra-ui/icons";
import { FiLogOut } from "react-icons/fi";
import { ReconnectButton } from "./Buttons";

const Navbar: React.FC = () => {
  const { isOpen, onToggle } = useDisclosure();
  const isLoggedIn = localStorage.getItem('isLoggedIn') === 'true';
  const username = localStorage.getItem('username');
  const [logged, setLogged] = useState<boolean>(isLoggedIn);
  const [recon, setRecon] = useState(false);
  const nav = useNavigate();

  useEffect(() => {
    setLogged(isLoggedIn);
  }, [isLoggedIn]);

  const handleLogout = () => {
    localStorage.clear();
    setLogged(false);
  };
  
  const checkOngoing = async () => {
  }

  const reconnect = async () => {
  };
  
  useEffect(() => {
    checkOngoing();
  }, []);

  return (
    <Box
      color="white"
      fontSize="3xl"
      position="sticky"
      top="0"
      zIndex="1000"
      bg="transparent"
    >
      <Box maxW="6xl" mx="auto" px={{ base: 4, sm: 6, lg: 8 }}>
        <Flex justify="space-between" h={20} align="center">
          <Flex align="center">
            <Link href="/" _hover={{ textDecoration: 'none' }}>
              <Text 
                as="h1" 
                fontSize={{ base: "4xl", md: "6xl" }} 
                fontWeight="900" 
                bgGradient="linear(to-r, #f6ad55, #ed8936, #dd6b20)"
                bgClip="text"
                _hover={{ 
                  bgGradient: "linear(to-r, #ed8936, #dd6b20, #c05621)",
                  transform: "scale(1.02)"
                }}
                transition="all 0.3s ease"
                letterSpacing="tight"
                textShadow="0 0 20px rgba(237, 137, 54, 0.3)"
              >
                VanGo
              </Text>
            </Link>
            {recon === true && (
              <ReconnectButton handler={reconnect} />
            )}
          </Flex>

          <HStack spacing={2} display={{ base: "none", md: "flex" }} align="center">
            {logged ? (
              <HStack spacing={2}>
                <Link
                  onClick={() => nav(`/profile/${username}`)}
                  color="orange.400"
                  fontWeight="600"
                  fontSize="2xl"
                  _hover={{ 
                    textDecoration: 'none',
                    color: "orange.100",
                    transform: "translateY(-2px)"
                  }}
                  transition="all 0.3s ease"
                  cursor="pointer"
                >
                  {username}
                </Link>
                <Button
                  onClick={handleLogout}
                  bg="linear-gradient(135deg, rgba(26, 32, 44, 0.8), rgba(45, 55, 72, 0.8))"
                  color="white"
                  _hover={{ 
                    bg: "linear-gradient(135deg, rgba(45, 55, 72, 0.9), rgba(74, 85, 104, 0.9))",
                    transform: "translateY(-2px)",
                    boxShadow: "0 8px 25px rgba(0, 0, 0, 0.3)"
                  }}
                  rounded="xl"
                  px={6}
                  py={3}
                  fontWeight="600"
                  transition="all 0.3s ease"
                  border="2px solid"
                  borderColor="whiteAlpha.300"
                >
                  <FiLogOut/>
                </Button>
              </HStack>
            ) : (
              <>
                <Button
                  as={Link}
                  href="/login"
                  bg="linear-gradient(135deg, rgba(26, 32, 44, 0.8), rgba(45, 55, 72, 0.8))"
                  color="white"
                  _hover={{ 
                    bg: "linear-gradient(135deg, rgba(45, 55, 72, 0.9), rgba(74, 85, 104, 0.9))",
                    textDecoration: "none",
                    transform: "translateY(-2px)",
                    boxShadow: "0 8px 25px rgba(0, 0, 0, 0.3)"
                  }}
                  rounded="xl"
                  px={6}
                  py={3}
                  fontWeight="600"
                  transition="all 0.3s ease"
                >
                  Login
                </Button>
                <Button
                  as={Link}
                  href="/signup"
                  bg="linear-gradient(135deg, #f6ad55, #ed8936)"
                  color="white"
                  _hover={{ 
                    bg: "linear-gradient(135deg, #ed8936, #dd6b20)",
                    textDecoration: "none",
                    transform: "translateY(-2px)",
                    boxShadow: "0 8px 25px rgba(246, 173, 85, 0.4)"
                  }}
                  rounded="xl"
                  px={6}
                  py={3}
                  fontWeight="600"
                  transition="all 0.3s ease"
                >
                  Sign Up
                </Button>
              </>
            )}
          </HStack>

          <IconButton
            size="lg"
            icon={isOpen ? <CloseIcon /> : <HamburgerIcon />}
            aria-label="Open Menu"
            display={{ md: "none" }}
            onClick={onToggle}
            bg="linear-gradient(135deg, rgba(246, 173, 85, 0.2), rgba(237, 137, 54, 0.2))"
            color="white"
            _hover={{ 
              bg: "linear-gradient(135deg, rgba(246, 173, 85, 0.3), rgba(237, 137, 54, 0.3))",
              transform: "rotate(90deg)"
            }}
            transition="all 0.3s ease"
            rounded="xl"
            border="2px solid"
            borderColor="whiteAlpha.300"
          />
        </Flex>
      </Box>

      <Collapse in={isOpen} animateOpacity>
        <Box 
          pb={4} 
          display={{ md: 'none' }}
          bg="linear-gradient(135deg, rgba(26, 32, 44, 0.9), rgba(45, 55, 72, 0.8))"
          backdropFilter="blur(12px)"
          borderTop="2px solid"
          borderColor="whiteAlpha.200"
        >
          <VStack spacing={3} px={4} pt={4}>
            {recon === true && (
              <Box w="full">
                <ReconnectButton handler={reconnect} />
              </Box>
            )}
            {logged ? (
              <>
                <Link
                  onClick={() => nav(`/profile/${username}`)}
                  color="orange.200"
                  fontWeight="600"
                  fontSize="lg"
                  _hover={{ 
                    textDecoration: 'none',
                    color: "orange.100",
                    transform: "translateX(8px)"
                  }}
                  w="full"
                  px={4}
                  py={3}
                  rounded="lg"
                  transition="all 0.3s ease"
                  cursor="pointer"
                >
                  {username}
                </Link>
                <Button
                  onClick={handleLogout}
                  bg="linear-gradient(135deg, rgba(26, 32, 44, 0.8), rgba(45, 55, 72, 0.8))"
                  color="white"
                  _hover={{ 
                    bg: "linear-gradient(135deg, rgba(45, 55, 72, 0.9), rgba(74, 85, 104, 0.9))",
                    transform: "translateY(-2px)"
                  }}
                  rounded="lg"
                  w="full"
                  px={4}
                  py={3}
                  fontWeight="600"
                  transition="all 0.3s ease"
                  border="2px solid"
                  borderColor="whiteAlpha.300"
                >
                  Logout
                </Button>
              </>
            ) : (
              <>
                <Button
                  as={Link}
                  href="/login"
                  bg="linear-gradient(135deg, rgba(26, 32, 44, 0.8), rgba(45, 55, 72, 0.8))"
                  color="white"
                  _hover={{ 
                    bg: "linear-gradient(135deg, rgba(45, 55, 72, 0.9), rgba(74, 85, 104, 0.9))",
                    textDecoration: "none",
                    transform: "translateY(-2px)"
                  }}
                  rounded="lg"
                  w="full"
                  px={4}
                  py={3}
                  fontWeight="600"
                  transition="all 0.3s ease"
                  border="2px solid"
                  borderColor="whiteAlpha.300"
                >
                  Login
                </Button>
                <Button
                  as={Link}
                  href="/signup"
                  bg="linear-gradient(135deg, #f6ad55, #ed8936)"
                  color="white"
                  _hover={{ 
                    bg: "linear-gradient(135deg, #ed8936, #dd6b20)",
                    textDecoration: "none",
                    transform: "translateY(-2px)"
                  }}
                  rounded="lg"
                  w="full"
                  px={4}
                  py={3}
                  fontWeight="600"
                  transition="all 0.3s ease"
                  border="2px solid"
                  borderColor="orange.400"
                >
                  Sign Up
                </Button>
              </>
            )}
          </VStack>
        </Box>
      </Collapse>
    </Box>
  );
};

export default Navbar;
