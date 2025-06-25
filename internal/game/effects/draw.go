package effects

import (
	"fmt"
	"math/rand"
)

// DrawEffect implements card drawing
type DrawEffect struct {
	cardCount int
}

// NewDrawEffect creates a new draw effect
func NewDrawEffect(count int) *DrawEffect {
	return &DrawEffect{
		cardCount: count,
	}
}

// Execute draws cards
func (e *DrawEffect) Execute(ctx *EffectContext) (*EffectResult, error) {
	result := &EffectResult{
		Success:    true,
		Messages:   []string{},
		CardsDrawn: []string{},
	}

	// Draw cards
	drawn := e.drawCards(ctx, e.cardCount)
	result.CardsDrawn = drawn
	
	if len(drawn) > 0 {
		result.Messages = append(result.Messages, 
			fmt.Sprintf("Drew %d cards", len(drawn)))
	} else {
		result.Messages = append(result.Messages, "No cards to draw")
	}

	return result, nil
}

// CanExecute checks if cards can be drawn
func (e *DrawEffect) CanExecute(ctx *EffectContext) (bool, string) {
	if len(ctx.PlayerState.Hand) >= 10 {
		return false, "hand is full"
	}
	if len(ctx.PlayerState.DrawPile) == 0 && len(ctx.PlayerState.DiscardPile) == 0 {
		return false, "no cards to draw"
	}
	return true, ""
}

// GetType returns the effect type
func (e *DrawEffect) GetType() string {
	return "draw"
}

// GetDescription returns the effect description
func (e *DrawEffect) GetDescription() string {
	return fmt.Sprintf("Draw %d cards", e.cardCount)
}

// drawCards draws cards from the draw pile
func (e *DrawEffect) drawCards(ctx *EffectContext, count int) []string {
	drawn := []string{}
	
	for i := 0; i < count && len(ctx.PlayerState.Hand) < 10; i++ {
		// If draw pile is empty, shuffle discard pile
		if len(ctx.PlayerState.DrawPile) == 0 {
			if len(ctx.PlayerState.DiscardPile) == 0 {
				break // No more cards to draw
			}
			e.shuffleDiscardIntoDraw(ctx)
		}
		
		if len(ctx.PlayerState.DrawPile) > 0 {
			// Draw top card
			card := ctx.PlayerState.DrawPile[0]
			ctx.PlayerState.DrawPile = ctx.PlayerState.DrawPile[1:]
			ctx.PlayerState.Hand = append(ctx.PlayerState.Hand, card)
			drawn = append(drawn, card)
		}
	}
	
	return drawn
}

// shuffleDiscardIntoDraw shuffles discard pile into draw pile
func (e *DrawEffect) shuffleDiscardIntoDraw(ctx *EffectContext) {
	ctx.PlayerState.DrawPile = ctx.PlayerState.DiscardPile
	ctx.PlayerState.DiscardPile = []string{}
	
	// Shuffle the draw pile
	for i := len(ctx.PlayerState.DrawPile) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		ctx.PlayerState.DrawPile[i], ctx.PlayerState.DrawPile[j] = 
			ctx.PlayerState.DrawPile[j], ctx.PlayerState.DrawPile[i]
	}
}

// ScryEffect lets player look at top cards and discard some
type ScryEffect struct {
	scryCount int
}

// NewScryEffect creates a scry effect
func NewScryEffect(count int) *ScryEffect {
	return &ScryEffect{
		scryCount: count,
	}
}

// Execute performs scry
func (e *ScryEffect) Execute(ctx *EffectContext) (*EffectResult, error) {
	result := &EffectResult{
		Success:  true,
		Messages: []string{},
	}

	// Get top cards
	topCards := []string{}
	for i := 0; i < e.scryCount && i < len(ctx.PlayerState.DrawPile); i++ {
		topCards = append(topCards, ctx.PlayerState.DrawPile[i])
	}

	// For now, we'll discard the bottom half (in a real game, player would choose)
	discardCount := len(topCards) / 2
	for i := 0; i < discardCount; i++ {
		card := ctx.PlayerState.DrawPile[len(topCards)-1-i]
		ctx.PlayerState.DrawPile = append(
			ctx.PlayerState.DrawPile[:len(topCards)-1-i],
			ctx.PlayerState.DrawPile[len(topCards)-i:]...)
		ctx.PlayerState.DiscardPile = append(ctx.PlayerState.DiscardPile, card)
	}

	result.Messages = append(result.Messages, 
		fmt.Sprintf("Scryed %d cards, discarded %d", len(topCards), discardCount))

	return result, nil
}

// CanExecute checks if scry can be performed
func (e *ScryEffect) CanExecute(ctx *EffectContext) (bool, string) {
	if len(ctx.PlayerState.DrawPile) == 0 {
		return false, "no cards in draw pile"
	}
	return true, ""
}

// GetType returns the effect type
func (e *ScryEffect) GetType() string {
	return "scry"
}

// GetDescription returns the effect description
func (e *ScryEffect) GetDescription() string {
	return fmt.Sprintf("Look at top %d cards and discard any", e.scryCount)
}

// DrawToHandSizeEffect draws cards up to a certain hand size
type DrawToHandSizeEffect struct {
	targetHandSize int
}

// NewDrawToHandSizeEffect creates a draw-to-hand-size effect
func NewDrawToHandSizeEffect(size int) *DrawToHandSizeEffect {
	return &DrawToHandSizeEffect{
		targetHandSize: size,
	}
}

// Execute draws cards until hand reaches target size
func (e *DrawToHandSizeEffect) Execute(ctx *EffectContext) (*EffectResult, error) {
	result := &EffectResult{
		Success:    true,
		Messages:   []string{},
		CardsDrawn: []string{},
	}

	// Calculate how many cards to draw
	cardsToDraw := e.targetHandSize - len(ctx.PlayerState.Hand)
	if cardsToDraw <= 0 {
		result.Messages = append(result.Messages, "Hand already at target size")
		return result, nil
	}

	// Use regular draw effect
	drawEffect := NewDrawEffect(cardsToDraw)
	return drawEffect.Execute(ctx)
}

// CanExecute checks if cards can be drawn
func (e *DrawToHandSizeEffect) CanExecute(ctx *EffectContext) (bool, string) {
	if len(ctx.PlayerState.Hand) >= e.targetHandSize {
		return false, "hand already at target size"
	}
	return true, ""
}

// GetType returns the effect type
func (e *DrawToHandSizeEffect) GetType() string {
	return "draw_to_hand_size"
}

// GetDescription returns the effect description
func (e *DrawToHandSizeEffect) GetDescription() string {
	return fmt.Sprintf("Draw cards until you have %d in hand", e.targetHandSize)
}