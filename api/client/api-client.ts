/**
 * Pixel Game API Client
 * 
 * Example implementation of an API client for the Pixel Game backend.
 * This can be used as a reference for frontend developers.
 */

import type {
  HealthResponse,
  CardsResponse,
  Card,
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  RegisterResponse,
  GameState,
  PlayCardRequest,
  PlayCardResponse,
  ErrorResponse,
  VersionResponse,
} from '../types/api.types';

export class PixelGameAPIClient {
  private baseURL: string;
  private token: string | null = null;

  constructor(baseURL: string = 'http://localhost:8080') {
    this.baseURL = baseURL;
  }

  // Set authentication token
  public setAuthToken(token: string) {
    this.token = token;
  }

  // Helper method for making requests
  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`;
    }

    try {
      const response = await fetch(url, {
        ...options,
        headers,
      });

      if (!response.ok) {
        const error: ErrorResponse = await response.json();
        throw new Error(error.message || 'API request failed');
      }

      return response.json();
    } catch (error) {
      console.error('API request error:', error);
      throw error;
    }
  }

  // Health check
  async checkHealth(): Promise<HealthResponse> {
    return this.request<HealthResponse>('/health');
  }

  // Authentication
  async login(credentials: LoginRequest): Promise<LoginResponse> {
    const response = await this.request<LoginResponse>('/api/v1/auth/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    });
    
    // Automatically set the token after successful login
    this.setAuthToken(response.token);
    
    return response;
  }

  async register(userData: RegisterRequest): Promise<RegisterResponse> {
    const response = await this.request<RegisterResponse>('/api/v1/auth/register', {
      method: 'POST',
      body: JSON.stringify(userData),
    });
    
    // Automatically set the token after successful registration
    this.setAuthToken(response.token);
    
    return response;
  }

  async logout(): Promise<void> {
    await this.request('/api/v1/auth/logout', {
      method: 'POST',
    });
    this.token = null;
  }

  // Cards
  async getCards(): Promise<CardsResponse> {
    return this.request<CardsResponse>('/api/v1/cards');
  }

  async getCard(id: number): Promise<Card> {
    return this.request<Card>(`/api/v1/cards/${id}`);
  }

  // Game
  async startGame(): Promise<GameState> {
    return this.request<GameState>('/api/v1/games/start', {
      method: 'POST',
    });
  }

  async getGameState(gameId: string): Promise<GameState> {
    return this.request<GameState>(`/api/v1/games/${gameId}`);
  }

  async playCard(request: PlayCardRequest): Promise<PlayCardResponse> {
    return this.request<PlayCardResponse>(
      `/api/v1/games/${request.gameId}/actions`,
      {
        method: 'POST',
        body: JSON.stringify(request),
      }
    );
  }

  async endGame(gameId: string): Promise<void> {
    await this.request(`/api/v1/games/${gameId}/end`, {
      method: 'POST',
    });
  }

  // System
  async getVersion(): Promise<VersionResponse> {
    return this.request<VersionResponse>('/api/v1/version');
  }

  // WebSocket connection for real-time updates
  connectWebSocket(gameId: string, handlers: {
    onOpen?: () => void;
    onMessage?: (data: any) => void;
    onError?: (error: Event) => void;
    onClose?: () => void;
  }): WebSocket {
    const wsUrl = this.baseURL.replace('http', 'ws') + `/ws?gameId=${gameId}`;
    const ws = new WebSocket(wsUrl);

    ws.onopen = () => {
      // Send authentication
      if (this.token) {
        ws.send(JSON.stringify({ type: 'auth', token: this.token }));
      }
      handlers.onOpen?.();
    };

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        handlers.onMessage?.(data);
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error);
      }
    };

    ws.onerror = handlers.onError || console.error;
    ws.onclose = handlers.onClose || (() => {});

    return ws;
  }
}

// Export a singleton instance
export const apiClient = new PixelGameAPIClient();

// Usage example:
/*
import { apiClient } from './api-client';

// Login
const loginResponse = await apiClient.login({
  username: 'player1',
  password: 'password123',
  platform: 'web'
});

// Get cards
const cards = await apiClient.getCards();

// Start a game
const gameState = await apiClient.startGame();

// Play a card
const result = await apiClient.playCard({
  gameId: gameState.session.id,
  cardId: 1,
  targetId: 'enemy-1'
});

// Connect WebSocket for real-time updates
const ws = apiClient.connectWebSocket(gameState.session.id, {
  onMessage: (data) => {
    console.log('Game update:', data);
  }
});
*/