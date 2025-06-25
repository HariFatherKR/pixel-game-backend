/**
 * Pixel Game API Types
 * 
 * This file contains TypeScript type definitions for the Pixel Game API.
 * These types can be imported and used in the frontend application.
 */

// Health Check
export interface HealthResponse {
  status: 'healthy' | 'unhealthy';
  service: string;
  timestamp: number;
}

// Card Types
export type CardType = 'action' | 'event' | 'power';

export interface Card {
  id: number;
  name: string;
  type: CardType;
  cost: number;
  description: string;
  // Additional properties for game mechanics
  damage?: number;
  shield?: number;
  draw?: number;
  energy?: number;
  effects?: CardEffect[];
}

export interface CardEffect {
  type: 'damage' | 'shield' | 'draw' | 'energy' | 'debuff' | 'buff';
  value: number;
  target: 'self' | 'enemy' | 'all_enemies';
  duration?: number; // For buffs/debuffs
}

export interface CardsResponse {
  cards: Card[];
  total: number;
}

// Game Types
export interface GameSession {
  id: string;
  userId: string;
  status: 'active' | 'completed' | 'abandoned';
  currentFloor: number;
  playerHealth: number;
  playerMaxHealth: number;
  playerEnergy: number;
  playerMaxEnergy: number;
  deck: Card[];
  hand: Card[];
  discardPile: Card[];
  score: number;
  createdAt: string;
  updatedAt: string;
}

export interface Enemy {
  id: string;
  name: string;
  health: number;
  maxHealth: number;
  intent: EnemyIntent;
  status: StatusEffect[];
}

export interface EnemyIntent {
  type: 'attack' | 'defend' | 'buff' | 'debuff' | 'unknown';
  value?: number;
}

export interface StatusEffect {
  type: 'vulnerable' | 'weak' | 'poison' | 'strength' | 'shield';
  value: number;
  duration: number;
}

// User Types
export interface User {
  id: string;
  username: string;
  email: string;
  platform: 'android' | 'ios' | 'web';
  stats: UserStats;
  createdAt: string;
  updatedAt: string;
}

export interface UserStats {
  gamesPlayed: number;
  gamesWon: number;
  highestScore: number;
  totalPlayTime: number; // in seconds
  cardsCollected: number;
  favoriteCard?: string;
}

// API Request/Response Types
export interface LoginRequest {
  username: string;
  password: string;
  platform: 'android' | 'ios' | 'web';
}

export interface LoginResponse {
  user: User;
  token: string;
  refreshToken: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  platform: 'android' | 'ios' | 'web';
}

export interface RegisterResponse {
  user: User;
  token: string;
  refreshToken: string;
}

export interface PlayCardRequest {
  gameId: string;
  cardId: number;
  targetId?: string; // For targeted effects
}

export interface PlayCardResponse {
  success: boolean;
  gameState: GameState;
  message?: string;
}

export interface GameState {
  session: GameSession;
  enemies: Enemy[];
  currentTurn: 'player' | 'enemy';
  turnNumber: number;
}

// Error Response
export interface ErrorResponse {
  error: string;
  message: string;
  code?: string;
  details?: any;
}

// Version
export interface VersionResponse {
  version: string;
  build: string;
}

// WebSocket Messages
export interface WebSocketMessage {
  type: 'game_update' | 'card_played' | 'enemy_action' | 'game_over' | 'error';
  payload: any;
  timestamp: number;
}

export interface GameUpdateMessage extends WebSocketMessage {
  type: 'game_update';
  payload: GameState;
}

export interface CardPlayedMessage extends WebSocketMessage {
  type: 'card_played';
  payload: {
    playerId: string;
    cardId: number;
    targetId?: string;
    effects: CardEffect[];
  };
}

// Utility Types
export type ApiResponse<T> = {
  success: boolean;
  data?: T;
  error?: ErrorResponse;
};

export type PaginatedResponse<T> = {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
  hasMore: boolean;
};