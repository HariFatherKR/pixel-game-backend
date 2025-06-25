-- Remove the unique constraint that prevents users from having multiple copies of the same card
ALTER TABLE user_cards DROP CONSTRAINT IF EXISTS user_cards_user_id_card_id_key;