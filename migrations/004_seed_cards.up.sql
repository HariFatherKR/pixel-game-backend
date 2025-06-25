-- Insert initial cards into the cards table

-- ACTION Cards (사이버 공격 카드)
INSERT INTO cards (id, name, type, rarity, cost, description, code_snippet, effects, visual_effects) VALUES
('card_001', '해킹 스트라이크', 'ACTION', 'COMMON', 2, '적에게 5 데미지를 입힙니다.', 
'damage(target, 5);', 
'[{"type": "damage", "target": "enemy", "value": 5}]'::jsonb,
'{"action": "shake", "target": ".enemy", "duration": 500}'::jsonb),

('card_002', '코드 인젝션', 'ACTION', 'COMMON', 3, '적에게 7 데미지를 입히고 1턴 동안 취약 상태로 만듭니다.', 
'damage(target, 7);\napplyDebuff(target, "vulnerable", 1);', 
'[{"type": "damage", "target": "enemy", "value": 7}, {"type": "debuff", "target": "enemy", "effect": "vulnerable", "duration": 1}]'::jsonb,
'{"action": "flash", "target": ".enemy", "color": "#ff0000", "duration": 300}'::jsonb),

('card_003', 'DDoS 공격', 'ACTION', 'RARE', 4, '모든 적에게 4 데미지를 입힙니다.', 
'enemies.forEach(e => damage(e, 4));', 
'[{"type": "damage", "target": "all_enemies", "value": 4}]'::jsonb,
'{"action": "pulse", "target": ".enemy", "count": 3, "duration": 1000}'::jsonb),

('card_004', '백도어', 'ACTION', 'RARE', 1, '카드를 2장 뽑고 이번 턴에 코스트가 1 감소합니다.', 
'drawCards(2);\nreduceCostThisTurn(1);', 
'[{"type": "draw", "value": 2}, {"type": "cost_reduction", "value": 1, "duration": "this_turn"}]'::jsonb,
'{"action": "glow", "target": ".hand", "color": "#00ff00", "duration": 500}'::jsonb),

('card_005', '시스템 크래시', 'ACTION', 'EPIC', 6, '적에게 15 데미지를 입힙니다. 이 카드를 사용하면 다음 턴을 건너뜁니다.', 
'damage(target, 15);\nskipNextTurn();', 
'[{"type": "damage", "target": "enemy", "value": 15}, {"type": "skip_turn", "target": "self"}]'::jsonb,
'{"action": "explode", "target": ".enemy", "particles": 50, "duration": 1500}'::jsonb),

('card_006', '바이러스 확산', 'ACTION', 'EPIC', 5, '적에게 3 데미지를 입히고, 3턴 동안 매 턴 2 데미지를 입힙니다.', 
'damage(target, 3);\napplyPoison(target, 2, 3);', 
'[{"type": "damage", "target": "enemy", "value": 3}, {"type": "poison", "target": "enemy", "value": 2, "duration": 3}]'::jsonb,
'{"action": "spread", "target": ".enemy", "color": "#9400d3", "duration": 1000}'::jsonb),

('card_007', '제로데이 익스플로잇', 'ACTION', 'LEGENDARY', 8, '적에게 20 데미지를 입히고 모든 방어막을 무시합니다.', 
'piercing_damage(target, 20);', 
'[{"type": "piercing_damage", "target": "enemy", "value": 20}]'::jsonb,
'{"action": "lightning", "target": ".enemy", "color": "#ffff00", "duration": 2000}'::jsonb),

-- EVENT Cards (이벤트 카드)
('card_008', '방화벽', 'EVENT', 'COMMON', 1, '5 방어막을 얻습니다.', 
'gainShield(5);', 
'[{"type": "shield", "target": "self", "value": 5}]'::jsonb,
'{"action": "shield_up", "target": ".player", "color": "#0080ff", "duration": 500}'::jsonb),

('card_009', '백업', 'EVENT', 'COMMON', 2, '체력을 7 회복합니다.', 
'heal(7);', 
'[{"type": "heal", "target": "self", "value": 7}]'::jsonb,
'{"action": "heal_effect", "target": ".player", "color": "#00ff00", "duration": 700}'::jsonb),

('card_010', '시스템 리부트', 'EVENT', 'RARE', 3, '모든 디버프를 제거하고 카드를 1장 뽑습니다.', 
'clearDebuffs();\ndrawCards(1);', 
'[{"type": "cleanse", "target": "self"}, {"type": "draw", "value": 1}]'::jsonb,
'{"action": "refresh", "target": ".player", "duration": 1000}'::jsonb),

('card_011', '오버클럭', 'EVENT', 'RARE', 2, '이번 턴에 추가로 에너지를 2 얻습니다.', 
'gainEnergy(2);', 
'[{"type": "energy", "value": 2, "duration": "this_turn"}]'::jsonb,
'{"action": "energy_surge", "target": ".energy-bar", "color": "#ff8000", "duration": 500}'::jsonb),

('card_012', '퀀텀 점프', 'EVENT', 'EPIC', 4, '다음 3장의 카드 코스트가 0이 됩니다.', 
'setNextCardsCost(3, 0);', 
'[{"type": "cost_reduction", "value": "all", "count": 3}]'::jsonb,
'{"action": "quantum_effect", "target": ".hand", "duration": 1500}'::jsonb),

('card_013', '타임 루프', 'EVENT', 'LEGENDARY', 7, '이번 턴에 사용한 모든 카드를 손으로 되돌립니다.', 
'returnPlayedCards();', 
'[{"type": "return_played_cards", "duration": "this_turn"}]'::jsonb,
'{"action": "time_rewind", "target": ".discard-pile", "duration": 2000}'::jsonb),

-- POWER Cards (지속 효과 카드)
('card_014', '알고리즘 최적화', 'POWER', 'COMMON', 2, '매 턴 시작 시 카드를 1장 더 뽑습니다.', 
'onTurnStart(() => drawCards(1));', 
'[{"type": "passive", "trigger": "turn_start", "effect": "draw", "value": 1}]'::jsonb,
'{"action": "permanent_glow", "target": ".draw-pile", "color": "#00ffff", "duration": -1}'::jsonb),

('card_015', 'AI 어시스턴트', 'POWER', 'RARE', 3, '카드를 사용할 때마다 1 방어막을 얻습니다.', 
'onCardPlay(() => gainShield(1));', 
'[{"type": "passive", "trigger": "card_play", "effect": "shield", "value": 1}]'::jsonb,
'{"action": "ai_presence", "target": ".player", "duration": -1}'::jsonb),

('card_016', '머신러닝', 'POWER', 'EPIC', 5, '매 턴 종료 시 손에 있는 카드 수만큼 데미지를 입힙니다.', 
'onTurnEnd(() => damage(target, hand.length));', 
'[{"type": "passive", "trigger": "turn_end", "effect": "damage", "value": "hand_size"}]'::jsonb,
'{"action": "data_flow", "target": ".enemy", "duration": -1}'::jsonb),

('card_017', '양자 컴퓨팅', 'POWER', 'LEGENDARY', 6, '매 턴 시작 시 에너지가 1 증가합니다.', 
'onTurnStart(() => increaseMaxEnergy(1));', 
'[{"type": "passive", "trigger": "turn_start", "effect": "max_energy", "value": 1}]'::jsonb,
'{"action": "quantum_field", "target": ".energy-bar", "duration": -1}'::jsonb),

-- Additional ACTION Cards
('card_018', '메모리 누수', 'ACTION', 'COMMON', 2, '적에게 4 데미지를 입히고 다음 턴에 추가로 4 데미지를 입힙니다.', 
'damage(target, 4);\ndelayedDamage(target, 4, 1);', 
'[{"type": "damage", "target": "enemy", "value": 4}, {"type": "delayed_damage", "target": "enemy", "value": 4, "delay": 1}]'::jsonb,
'{"action": "leak_effect", "target": ".enemy", "duration": 800}'::jsonb),

('card_019', '버퍼 오버플로우', 'ACTION', 'RARE', 3, '적의 방어막을 모두 제거하고 그만큼 데미지를 입힙니다.', 
'let shield = target.shield;\nremoveShield(target);\ndamage(target, shield);', 
'[{"type": "shield_break", "target": "enemy"}, {"type": "damage", "target": "enemy", "value": "shield_amount"}]'::jsonb,
'{"action": "shatter", "target": ".enemy", "duration": 1000}'::jsonb),

('card_020', '트로이 목마', 'ACTION', 'EPIC', 4, '적에게 10 데미지를 입힙니다. 적이 다음 턴에 카드를 사용하면 5 추가 데미지를 입습니다.', 
'damage(target, 10);\napplyTrap(target, "card_play", 5);', 
'[{"type": "damage", "target": "enemy", "value": 10}, {"type": "trap", "target": "enemy", "trigger": "card_play", "value": 5}]'::jsonb,
'{"action": "trojan_install", "target": ".enemy", "duration": 1200}'::jsonb);

-- Grant initial cards to new users (this will be handled by the application logic when users register)