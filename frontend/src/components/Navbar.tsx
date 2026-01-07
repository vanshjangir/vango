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
      color="gray.100"
      fontSize="3xl"
      position="sticky"
      top="0"
      zIndex="1000"
      bg="black"
      borderColor="gray.800"
    >
      <Box maxW="6xl" mx="auto" px={{ base: 4, sm: 6, lg: 8 }}>
        <Flex justify="space-between" h={20} align="center">
          <Flex align="center">
            <Link href="/" _hover={{ textDecoration: 'none' }}>
              <Text
                as="h1"
                fontSize={{ base: "4xl", md: "6xl" }}
                color="white"
              >
                VanGo
              </Text>
            </Link>
            {recon === true && (
              <ReconnectButton handler={reconnect} />
            )}
          </Flex>

          <HStack spacing={3} display={{ base: "none", md: "flex" }} align="center">
            {logged ? (
              <HStack spacing={3}>
                <Link
                  onClick={() => nav(`/profile/${username}`)}
                  color="gray.300"
                  fontWeight="500"
                  fontSize="2xl"
                  cursor="pointer"
                  _hover={{ color: "white" }}
                >
                  {username}
                </Link>
                <Button
                  onClick={handleLogout}
                  bg="gray.800"
                  color="gray.200"
                  border="1px solid"
                  borderColor="gray.700"
                  px={4}
                  py={2}
                  fontWeight="500"
                  _hover={{ borderColor: "gray.600" }}
                >
                  <FiLogOut/>
                </Button>
              </HStack>
            ) : (
              <>
                <Button
                  as={Link}
                  href="/login"
                  bg="gray.800"
                  color="gray.200"
                  border="1px solid"
                  borderColor="gray.700"
                  rounded="2px"
                  px={5}
                  py={2}
                  fontWeight="500"
                  _hover={{ borderColor: "gray.600" }}
                >
                  Login
                </Button>
                <Button
                  as={Link}
                  href="/signup"
                  bg="white"
                  color="gray.900"
                  rounded="2px"
                  px={5}
                  py={2}
                  fontWeight="500"
                  _hover={{ bg: "gray.100" }}
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
            bg="gray.800"
            color="gray.200"
            border="1px solid"
            borderColor="gray.700"
            _hover={{ borderColor: "gray.600" }}
            rounded="2px"
          />
        </Flex>
      </Box>

      <Collapse in={isOpen} animateOpacity>
        <Box
          pb={4}
          display={{ md: 'none' }}
          bg="gray.900"
          borderBottom="1px solid"
          borderColor="gray.800"
        >
          <VStack spacing={2} px={4} pt={4}>
            {recon === true && (
              <Box w="full">
                <ReconnectButton handler={reconnect} />
              </Box>
            )}
            {logged ? (
              <>
                <Link
                  onClick={() => nav(`/profile/${username}`)}
                  color="gray.300"
                  fontWeight="500"
                  fontSize="lg"
                  _hover={{ color: "white" }}
                  w="full"
                  px={4}
                  py={3}
                  cursor="pointer"
                >
                  {username}
                </Link>
                <Button
                  onClick={handleLogout}
                  bg="gray.800"
                  color="gray.200"
                  border="1px solid"
                  borderColor="gray.700"
                  _hover={{ borderColor: "gray.600" }}
                  rounded="2px"
                  w="full"
                  px={4}
                  py={3}
                  fontWeight="500"
                >
                  Logout
                </Button>
              </>
            ) : (
              <>
                <Button
                  as={Link}
                  href="/login"
                  bg="gray.800"
                  color="gray.200"
                  border="1px solid"
                  borderColor="gray.700"
                  _hover={{ borderColor: "gray.600", textDecoration: "none" }}
                  rounded="2px"
                  w="full"
                  px={4}
                  py={3}
                  fontWeight="500"
                >
                  Login
                </Button>
                <Button
                  as={Link}
                  href="/signup"
                  bg="white"
                  color="gray.900"
                  _hover={{ bg: "gray.100", textDecoration: "none" }}
                  rounded="2px"
                  w="full"
                  px={4}
                  py={3}
                  fontWeight="500"
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
