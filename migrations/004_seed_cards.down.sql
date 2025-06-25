-- Remove all seeded cards
DELETE FROM cards WHERE id IN (
    'card_001', 'card_002', 'card_003', 'card_004', 'card_005',
    'card_006', 'card_007', 'card_008', 'card_009', 'card_010',
    'card_011', 'card_012', 'card_013', 'card_014', 'card_015',
    'card_016', 'card_017', 'card_018', 'card_019', 'card_020'
);