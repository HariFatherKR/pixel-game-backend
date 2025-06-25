# API Integration Guide for Frontend Developers

## 📚 Overview

This guide helps frontend developers integrate with the Pixel Game backend API. The backend provides a RESTful API with WebSocket support for real-time game updates.

## 🚀 Quick Start

### 1. API Documentation

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Base URL**: http://localhost:8080
- **API Version**: v1

### 2. TypeScript Types

All API types are defined in `/backend/api/types/api.types.ts`. Copy this file to your frontend project:

```typescript
import type { Card, GameState, LoginRequest } from './api.types';
```

### 3. API Client Example

A complete API client implementation is available in `/backend/api/client/api-client.ts`.

```typescript
import { apiClient } from './api-client';

// Example usage
const cards = await apiClient.getCards();
```

## 🔐 Authentication

The API uses JWT (JSON Web Token) for authentication.

### Login Flow

1. Send credentials to `/api/v1/auth/login`
2. Receive JWT token in response
3. Include token in all subsequent requests

```typescript
// Login
const response = await fetch('http://localhost:8080/api/v1/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    username: 'player1',
    password: 'password123',
    platform: 'web'
  })
});

const { token, user } = await response.json();

// Use token in subsequent requests
const cards = await fetch('http://localhost:8080/api/v1/cards', {
  headers: { 'Authorization': `Bearer ${token}` }
});
```

## 📡 WebSocket Integration

For real-time game updates, connect to the WebSocket endpoint:

```typescript
const ws = new WebSocket('ws://localhost:8080/ws?gameId=123');

ws.onopen = () => {
  // Send authentication
  ws.send(JSON.stringify({ type: 'auth', token: authToken }));
};

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  
  switch (data.type) {
    case 'game_update':
      updateGameState(data.payload);
      break;
    case 'card_played':
      animateCardPlay(data.payload);
      break;
    case 'enemy_action':
      showEnemyAction(data.payload);
      break;
  }
};
```

## 🎮 Game Flow Integration

### 1. Starting a Game

```typescript
// Start new game
const gameState = await apiClient.startGame();

// Connect WebSocket for real-time updates
const ws = apiClient.connectWebSocket(gameState.session.id, {
  onMessage: (data) => updateUI(data)
});
```

### 2. Playing Cards

```typescript
// Play a card
const result = await apiClient.playCard({
  gameId: gameState.session.id,
  cardId: selectedCard.id,
  targetId: selectedEnemy?.id
});

// The WebSocket will receive the update automatically
```

### 3. Handling Game State

```typescript
interface GameState {
  session: GameSession;
  enemies: Enemy[];
  currentTurn: 'player' | 'enemy';
  turnNumber: number;
}

// Update UI based on game state
function updateGameUI(state: GameState) {
  // Update player health
  playerHealthBar.value = state.session.playerHealth;
  
  // Update hand
  renderHand(state.session.hand);
  
  // Update enemies
  state.enemies.forEach(enemy => renderEnemy(enemy));
}
```

## 🎨 Card Rendering

Cards have a specific structure that maps to visual effects:

```typescript
interface Card {
  id: number;
  name: string;
  type: 'action' | 'event' | 'power';
  cost: number;
  description: string;
  effects?: CardEffect[];
}

// Render card based on type
function renderCard(card: Card) {
  const cardElement = document.createElement('div');
  cardElement.className = `card card-${card.type}`;
  
  // Add cost indicator
  if (card.cost > 0) {
    cardElement.dataset.cost = card.cost.toString();
  }
  
  // Add effects
  card.effects?.forEach(effect => {
    cardElement.classList.add(`effect-${effect.type}`);
  });
  
  return cardElement;
}
```

## 🛠️ Development Tools

### Postman Collection

Import `/backend/api/postman/pixel-game-api.postman_collection.json` into Postman for easy API testing.

### Mock Data

For frontend development without the backend:

```typescript
// Mock cards data
export const mockCards: Card[] = [
  {
    id: 1,
    name: "Code Slash",
    type: "action",
    cost: 2,
    description: "Deal 8 damage and apply Vulnerable",
    effects: [
      { type: 'damage', value: 8, target: 'enemy' },
      { type: 'debuff', value: 1, target: 'enemy', duration: 2 }
    ]
  },
  // ... more cards
];

// Mock game state
export const mockGameState: GameState = {
  session: {
    id: "mock-game-123",
    userId: "user-123",
    status: "active",
    currentFloor: 1,
    playerHealth: 80,
    playerMaxHealth: 100,
    // ... etc
  },
  enemies: [
    {
      id: "enemy-1",
      name: "Rogue AI",
      health: 45,
      maxHealth: 45,
      intent: { type: 'attack', value: 12 }
    }
  ],
  currentTurn: 'player',
  turnNumber: 1
};
```

## 🔍 Error Handling

All API errors follow a consistent format:

```typescript
interface ErrorResponse {
  error: string;
  message: string;
  code?: string;
  details?: any;
}

// Handle errors gracefully
try {
  const result = await apiClient.playCard(request);
} catch (error) {
  if (error.code === 'INSUFFICIENT_ENERGY') {
    showNotification('Not enough energy to play this card');
  } else {
    showError(error.message);
  }
}
```

## 📋 API Endpoints Summary

### Public Endpoints (No Auth Required)
- `GET /api/v1/health` - Health check
- `GET /api/v1/version` - API version

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/logout` - Logout (requires auth)
- `POST /api/v1/auth/refresh` - Refresh token
- `GET /api/v1/auth/profile` - Get current user profile (requires auth)

### Cards (Requires Auth)
- `GET /api/v1/cards` - List all cards (with filtering and pagination)
- `GET /api/v1/cards/:id` - Get specific card
- `GET /api/v1/cards/my-collection` - Get user's card collection

### Decks (Requires Auth)
- `POST /api/v1/cards/decks` - Create new deck
- `GET /api/v1/cards/decks` - List user's decks
- `GET /api/v1/cards/decks/:id` - Get specific deck
- `PUT /api/v1/cards/decks/:id` - Update deck
- `DELETE /api/v1/cards/decks/:id` - Delete deck
- `PUT /api/v1/cards/decks/:id/activate` - Set deck as active
- `GET /api/v1/cards/decks/active` - Get active deck

### Game (Requires Auth)
- `POST /api/v1/games/start` - Start new game
- `GET /api/v1/games/current` - Get current active game
- `GET /api/v1/games/:id` - Get specific game
- `POST /api/v1/games/:id/actions` - Play action (card play, etc.)
- `POST /api/v1/games/:id/end-turn` - End turn
- `POST /api/v1/games/:id/surrender` - Surrender game
- `GET /api/v1/games/stats` - Get user's game statistics

### User (Requires Auth)
- `GET /api/v1/users/profile` - Get user profile
- `PUT /api/v1/users/profile` - Update profile
- `GET /api/v1/users/stats` - Get user statistics
- `GET /api/v1/users/collection` - Get card collection
- `POST /api/v1/users/stats/games-played` - Increment games played
- `POST /api/v1/users/stats/games-won` - Increment games won
- `POST /api/v1/users/stats/play-time/:seconds` - Add play time

### Real-time (Future)
- `WS /ws` - WebSocket connection for game updates

## 🔧 CORS Configuration

The backend is configured to accept requests from:
- http://localhost:3000 (default React dev server)
- http://localhost:5173 (default Vite dev server)
- Your production domain

For development, no additional CORS configuration is needed.

## 💡 Best Practices

1. **Type Safety**: Always use the provided TypeScript types
2. **Error Handling**: Implement proper error handling for all API calls
3. **Loading States**: Show loading indicators during API calls
4. **Offline Support**: Cache card data for offline viewing
5. **WebSocket Reconnection**: Implement automatic reconnection logic
6. **Token Management**: Store tokens securely and handle expiration

## 🤝 Need Help?

- Check the Swagger documentation for detailed API specs
- Review the example API client implementation
- Test endpoints using the Postman collection
- Check backend logs for debugging: `docker compose logs backend`

Happy coding! 🎮