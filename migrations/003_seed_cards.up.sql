-- Insert initial card data

-- Action Cards
INSERT INTO cards (id, name, type, rarity, cost, description, code_snippet, effects, visual_effects) VALUES
('code_slash', '코드 슬래시', 'ACTION', 'COMMON', 1, '적 하나에 8의 물리 대미지를 주고, 50% 확률로 취약 디버프를 건다.', 
'function codeSlash(enemy) {
  enemy.hp -= 8;
  if (Math.random() < 0.5) {
    enemy.addDebuff("vulnerable");
  }
}', 
'[{"type": "DAMAGE", "target": "ENEMY", "value": 8}, {"type": "DEBUFF", "target": "ENEMY", "value": 1, "conditions": {"chance": 0.5}}]'::jsonb,
'{"action": "create", "element": "div", "class": "slash-effect", "target": "enemy"}'::jsonb),

('firewall_up', '방화벽 업', 'ACTION', 'COMMON', 1, '자신의 방어력을 10 올린다.',
'function firewallUp(player) {
  player.shield += 10;
  player.element.style.outline = "3px solid cyan";
}',
'[{"type": "SHIELD", "target": "SELF", "value": 10}]'::jsonb,
'{"action": "modify", "selector": ".player", "style": {"outline": "3px solid cyan"}}'::jsonb),

('debug_punch', '디버그 펀치', 'ACTION', 'COMMON', 1, '적 하나에 6의 대미지를 준다. 업그레이드 시 대미지 증가 또는 코스트 0.',
'function debugPunch(enemy) {
  enemy.hp -= 6;
  console.log("Debug: Enemy HP =", enemy.hp);
}',
'[{"type": "DAMAGE", "target": "ENEMY", "value": 6}]'::jsonb,
'{"action": "console", "message": "Debug: Enemy HP = ${enemy.hp}"}'::jsonb),

('virus_upload', '바이러스 업로드', 'ACTION', 'RARE', 2, '모든 적에게 독 5를 부여한다.',
'function virusUpload(enemies) {
  enemies.forEach(enemy => {
    enemy.addDebuff("poison", 5);
  });
}',
'[{"type": "DEBUFF", "target": "ALL_ENEMIES", "value": 5}]'::jsonb,
'{"action": "create", "element": "particle", "class": "virus-particle", "target": "all-enemies"}'::jsonb),

('system_restore', '시스템 복원', 'ACTION', 'RARE', 2, '체력을 15 회복하고 모든 디버프를 제거한다.',
'function systemRestore(player) {
  player.hp = Math.min(player.hp + 15, player.maxHp);
  player.clearDebuffs();
}',
'[{"type": "HEAL", "target": "SELF", "value": 15}, {"type": "CLEANSE", "target": "SELF"}]'::jsonb,
'{"action": "animate", "selector": ".player", "animation": "restore-glow"}'::jsonb),

('memory_leak', '메모리 누수', 'ACTION', 'EPIC', 3, '적 하나에 현재 턴 수 x 5의 대미지를 준다.',
'function memoryLeak(enemy, gameState) {
  const damage = gameState.turn * 5;
  enemy.hp -= damage;
}',
'[{"type": "DAMAGE", "target": "ENEMY", "value": 0, "conditions": {"multiply": "turn", "factor": 5}}]'::jsonb,
'{"action": "create", "element": "div", "class": "memory-leak", "animate": "grow"}'::jsonb),

('infinite_loop', '무한 루프', 'ACTION', 'LEGENDARY', 3, '이번 턴에 사용한 카드 수만큼 카드를 뽑는다.',
'function infiniteLoop(player, gameState) {
  const cardsToDraw = gameState.cardsPlayedThisTurn;
  player.drawCards(cardsToDraw);
}',
'[{"type": "DRAW", "target": "SELF", "value": 0, "conditions": {"equals": "cardsPlayedThisTurn"}}]'::jsonb,
'{"action": "animate", "selector": ".deck", "animation": "infinite-spin"}'::jsonb),

-- Event Cards
('bug_found', '버그 발견', 'EVENT', 'COMMON', 0, '현재 층의 모든 함정이 무력화되고, 무작위 카드 한 장을 획득한다.',
'function bugFound(gameState) {
  gameState.traps = [];
  console.log("Bug found! Traps deactivated.");
  return gameState.addRandomCard();
}',
'[{"type": "DISABLE_TRAPS", "target": "FLOOR"}, {"type": "ADD_CARD", "target": "SELF", "value": 1}]'::jsonb,
'{"action": "console", "message": "Bug found! Traps deactivated.", "fadeOut": ".trap"}'::jsonb),

('glitch_out', '글리치 아웃', 'EVENT', 'RARE', 0, '3턴간 피해를 받지 않으면 모든 적의 방어력 30% 감소.',
'function glitchOut(gameState) {
  if (gameState.turnsWithoutDamage >= 3) {
    document.body.classList.add("glitch");
    gameState.enemies.forEach(e => e.defense *= 0.7);
  }
}',
'[{"type": "CONDITIONAL_DEBUFF", "target": "ALL_ENEMIES", "value": 30, "conditions": {"turnsWithoutDamage": 3}}]'::jsonb,
'{"action": "glitch", "selector": "body", "duration": 2000}'::jsonb),

('cache_overflow', '캐시 오버플로우', 'EVENT', 'EPIC', 0, '다음 3장의 카드 비용이 0이 된다.',
'function cacheOverflow(player) {
  player.addBuff("zeroCost", 3);
  console.log("Cache overflow! Next 3 cards cost 0.");
}',
'[{"type": "BUFF", "target": "SELF", "value": 3, "conditions": {"zeroCost": true}}]'::jsonb,
'{"action": "highlight", "selector": ".hand", "effect": "zero-cost-glow"}'::jsonb),

('kernel_panic', '커널 패닉', 'EVENT', 'LEGENDARY', 0, '모든 적에게 99의 대미지. 사용 후 덱에서 제거.',
'function kernelPanic(enemies) {
  document.body.style.backgroundColor = "blue";
  enemies.forEach(enemy => enemy.hp -= 99);
  setTimeout(() => location.reload(), 3000);
}',
'[{"type": "DAMAGE", "target": "ALL_ENEMIES", "value": 99}, {"type": "EXHAUST", "target": "SELF"}]'::jsonb,
'{"action": "bsod", "duration": 3000, "reload": true}'::jsonb);

-- Basic starter deck composition
INSERT INTO cards (id, name, type, rarity, cost, description, code_snippet, effects) VALUES
('basic_attack', '기본 공격', 'ACTION', 'COMMON', 1, '적 하나에 5의 대미지를 준다.',
'function basicAttack(enemy) { enemy.hp -= 5; }',
'[{"type": "DAMAGE", "target": "ENEMY", "value": 5}]'::jsonb),

('basic_defend', '기본 방어', 'ACTION', 'COMMON', 1, '방어력을 5 얻는다.',
'function basicDefend(player) { player.shield += 5; }',
'[{"type": "SHIELD", "target": "SELF", "value": 5}]'::jsonb);