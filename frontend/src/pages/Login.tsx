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
      bg="linear-gradient(135deg, #1a202c 0%, #2d3748 25%, #4a5568 50%, #2d3748 75%, #1a202c 100%)"
      color="white" 
      direction="column"
      position="relative"
      overflow="hidden"
    >
      <Flex
        direction="column"
        alignItems="center"
        justifyContent="center"
        flex="1"
        px={6}
        py={12}
        position="relative"
        zIndex={1}
      >
        {/* Error Message */}
        {error && (
          <Box
            bg="red.900"
            borderColor="red.500"
            border="1px solid"
            color="red.200"
            px={4}
            py={3}
            borderRadius="lg"
            mb={6}
            maxW="md"
            w="full"
          >
            <Text textAlign="center">{error}</Text>
          </Box>
        )}

        {/* Form Container */}
        <Box 
          w="full" 
          maxW="md" 
          backdropFilter="blur(12px)"
          p={8} 
          boxShadow="0 8px 32px rgba(0, 0, 0, 0.3)"
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
            
            <Divider borderColor="whiteAlpha.300" flex="1" />

            <Button
              onClick={handleGuestLogin}
              w="176px"
              py={3}
              h="auto"
              bg="linear-gradient(135deg, #f6ad55, #ed8936)"
              color="white"
              borderRadius="4px"
              fontWeight="700"
              fontSize="lg"
              transition="all 0.3s ease"
              _hover={{ 
                transform: "translateY(-2px)",
                boxShadow: "0 8px 25px rgba(237, 137, 54, 0.3)",
                bg: "linear-gradient(135deg, #ed8936, #dd6b20)"
              }}
              _active={{ transform: "translateY(0)" }}
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
