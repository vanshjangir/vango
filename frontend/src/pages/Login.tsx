import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { CredentialResponse, GoogleLogin } from "@react-oauth/google";
import {
  Box,
  Flex,
  Text,
  Button,
  VStack,
  Divider,
} from "@chakra-ui/react";

const Login: React.FC = () => {
  const [error, setError] = useState("");
  const nav = useNavigate();
  const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;

  const handleGoogleLogin = async (credentialResponse: CredentialResponse) => {
    if (!credentialResponse.credential) {
      setError("Google Login unsuccessful");
      console.log("Google Login unsuccessful");
      return;
    }

    const response = await fetch(BACKEND_URL + '/login', {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        "type": "google",
        "credentials": credentialResponse.credential,
      }),
    });

    if (response.status === 200) {
      const json = await response.json();
      console.log("Login successful");
      localStorage.setItem("token", json.token);
      localStorage.setItem("tokenType", "google");
      localStorage.setItem("username", json.username);
      localStorage.setItem('isLoggedIn', 'true');
      nav("/");
    } else {
      setError("Google Login unsuccessful");
      console.log("Google Login unsuccessful");
    }
  };

  const handleGuestLogin = async () => {
    const response = await fetch(BACKEND_URL + '/login', {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        "type": "guest",
        "credentials": "",
      }),
    });

    if (response.status === 200) {
      const json = await response.json();
      console.log("Login successful");
      localStorage.setItem("token", json.token);
      localStorage.setItem("tokenType", "guest");
      localStorage.setItem("username", json.username);
      localStorage.setItem('isLoggedIn', 'true');
      nav("/");
    } else {
      setError("Login unsuccessful");
      console.log("Login unsuccessful");
    }
  }

  return (
    <Flex 
      minH="100vh" 
      bg="black"
      color="white" 
      direction="column"
    >
      <Flex
        direction="column"
        alignItems="center"
        justifyContent="center"
        flex="1"
        px={6}
        py={12}
      >
        {/* Error Message */}
        {error && (
          <Box
            bg="gray.900"
            borderColor="red.800"
            border="1px solid"
            color="red.400"
            px={4}
            py={3}
            borderRadius="2px"
            mb={6}
            maxW="md"
            w="full"
          >
            <Text textAlign="center" fontWeight="500">{error}</Text>
          </Box>
        )}

        {/* Form Container */}
        <Box 
          w="full" 
          maxW="md" 
          bg="gray.900"
          border="1px solid"
          borderColor="gray.800"
          p={8}
          borderRadius="2px"
        >
          <VStack spacing={6}>
            <GoogleLogin
              onSuccess={(credentialResponse) => {
                handleGoogleLogin(credentialResponse);
              }}
              onError={() => {
                console.log("Google Login Failed");
                setError("Google Login unsuccessful");
              }}
              useOneTap
            />
            
            <Divider borderColor="gray.700" flex="1" />

            <Button
              onClick={handleGuestLogin}
              w="176px"
              py={3}
              h="auto"
              bg="white"
              color="gray.900"
              borderRadius="2px"
              fontWeight="200"
              fontSize="lg"
              _hover={{ 
                bg: "gray.100"
              }}
            >
              Login as Guest
            </Button>
          </VStack>
        </Box>
      </Flex>
    </Flex>
  );
};

export default Login;
