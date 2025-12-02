<?php

namespace SatBridge;

class AuthMiddleware
{
    public function authenticate(): bool
    {
        $apiKey = $_ENV['API_KEY'] ?? '';
        
        if (empty($apiKey)) {
            return false;
        }

        // Check X-API-Key header
        $headers = getallheaders();
        $providedKey = $headers['X-API-Key'] ?? $headers['x-api-key'] ?? '';

        return hash_equals($apiKey, $providedKey);
    }
}

