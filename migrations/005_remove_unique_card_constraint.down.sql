-- Re-add the unique constraint
ALTER TABLE user_cards ADD CONSTRAINT user_cards_user_id_card_id_key UNIQUE (user_id, card_id);