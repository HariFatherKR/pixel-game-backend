# Pixel Game Backend Development Tasks

## 📋 개발 진행 상황

### ✅ Phase 1: 프로젝트 초기 설정 (완료)
- [x] Go 프로젝트 구조 설정 (DDD 패턴)
- [x] 데이터베이스 스키마 설계
- [x] PostgreSQL 마이그레이션 파일 작성
- [x] 기본 도메인 엔티티 구현 (User, Card, Game)
- [x] 환경 설정 시스템 구축
- [x] Docker 환경 구성
- [x] 기본 헬스체크 API 구현

### ✅ Phase 2: 개발 환경 강화 (완료)
- [x] Swagger/OpenAPI 문서화 설정
- [x] TypeScript 타입 정의 파일 생성
- [x] API 클라이언트 예제 구현
- [x] Postman 컬렉션 생성
- [x] Docker 재빌드 자동화 스크립트
- [x] PostgreSQL 연결 및 데이터 검증
- [x] Redis 연결 및 동작 확인

### ✅ Phase 3: 인증 시스템 구현 (완료)
- [x] JWT 토큰 기반 인증 구현
- [x] 회원가입 API (`POST /api/v1/auth/register`)
- [x] 로그인 API (`POST /api/v1/auth/login`)
- [x] 로그아웃 API (`POST /api/v1/auth/logout`)
- [x] 토큰 갱신 API (`POST /api/v1/auth/refresh`)
- [x] 인증 미들웨어 구현
- [x] 비밀번호 해싱 (bcrypt)
- [x] 사용자 프로필 조회 API (`GET /api/v1/auth/profile`)

### ✅ Phase 4: 사용자 관리 시스템 (완료)
- [x] PostgreSQL 연동 Repository 패턴 구현
- [x] User 도메인 모델 및 데이터베이스 스키마 설계
- [x] 실제 데이터베이스 기반 인증 시스템 업그레이드
- [x] 사용자 프로필 수정 API (`PUT /api/v1/users/profile`)
- [x] 사용자 통계 조회 API (`GET /api/v1/users/stats`)
- [x] 사용자 카드 컬렉션 조회 API (`GET /api/v1/users/collection`)
- [x] 플랫폼별 사용자 구분 (Android/iOS/Web)
- [x] 게임 통계 업데이트 API 구현

### ✅ Phase 5: 카드 시스템 구현 (완료)
- [x] 카드 도메인 모델 및 Repository 구현
- [x] 카드 목록 조회 API (`GET /api/v1/cards`)
- [x] 카드 상세 조회 API (`GET /api/v1/cards/:id`)
- [x] 사용자 카드 컬렉션 시스템 구현
  - [x] 사용자 카드 컬렉션 조회 (`GET /api/v1/cards/my-collection`)
  - [x] 신규 사용자 초기 카드 지급 시스템
  - [x] 카드 중복 소유 가능하도록 DB 스키마 수정
- [x] 덱 관리 시스템 구현
  - [x] 덱 생성 API (`POST /api/v1/cards/decks`)
  - [x] 덱 목록 조회 API (`GET /api/v1/cards/decks`)
  - [x] 덱 상세 조회 API (`GET /api/v1/cards/decks/:id`)
  - [x] 덱 수정 API (`PUT /api/v1/cards/decks/:id`)
  - [x] 덱 삭제 API (`DELETE /api/v1/cards/decks/:id`)
  - [x] 덱 활성화 API (`PUT /api/v1/cards/decks/:id/activate`)
  - [x] 활성 덱 조회 API (`GET /api/v1/cards/decks/active`)
- [x] 카드 마스터 데이터 시드 (20개 카드)
- [ ] 카드 효과 처리 엔진 구현
- [ ] 카드 업그레이드 시스템
- [ ] 카드 시너지 계산 로직

### ✅ Phase 6: 게임 플레이 시스템 (완료)
- [x] 게임 세션 도메인 모델 설계 (GameSession, PlayerState, EnemyState)
- [x] PostgreSQL 기반 Game Repository 구현
- [x] 게임 시작 API (`POST /api/v1/games/start`)
- [x] 게임 상태 조회 API (`GET /api/v1/games/:id`)
- [x] 카드 플레이 API (`POST /api/v1/games/:id/actions`)
- [x] 턴 종료 API (`POST /api/v1/games/:id/end-turn`)
- [x] 게임 포기 API (`POST /api/v1/games/:id/surrender`)
- [x] 게임 통계 API (`GET /api/v1/games/stats`)
- [x] 기본 턴제 전투 시스템 구현
- [x] 게임 핸들러 서버 통합

### 📅 Phase 7: 실시간 통신
- [ ] WebSocket 연결 관리
- [ ] 게임 상태 실시간 동기화
- [ ] 카드 플레이 애니메이션 이벤트
- [ ] 연결 끊김 처리 및 재연결 로직

### 📅 Phase 8: 게임 진행 시스템
- [ ] 던전 맵 생성 (절차적 생성)
- [ ] 층별 이벤트 시스템
- [ ] 상점 시스템
- [ ] 휴식 포인트 시스템
- [ ] 보스 전투 시스템

### 📅 Phase 9: 메타 시스템
- [ ] 리더보드 API (`GET /api/v1/leaderboard`)
- [ ] 일일 챌린지 시스템 (`GET /api/v1/challenges/daily`)
- [ ] 업적 시스템
- [ ] 카드 해금 시스템
- [ ] 메타 프로그레션 (영구 업그레이드)

### 📅 Phase 10: 최적화 및 보안
- [ ] API Rate Limiting
- [ ] 입력 검증 강화
- [ ] SQL 인젝션 방지
- [ ] 서버 사이드 치팅 방지
- [ ] 캐싱 전략 구현 (Redis)
- [ ] 데이터베이스 쿼리 최적화

## 🎯 현재 작업 중인 태스크

### Phase 7: 카드 효과 실행 엔진
카드의 실제 효과를 처리하고 게임 상태를 업데이트하는 시스템을 구현합니다.

**다음 작업:**
- [ ] 카드 효과 처리 엔진 구현
- [ ] 카드 타입별 효과 실행 로직
- [ ] 버프/디버프 시스템 구현
- [ ] 카드 시너지 계산 로직
- [ ] 적 AI 행동 패턴 구현

## 📊 진행률

- **전체 진행률**: 60% (6/10 Phase 완료)
- **현재 Phase 진행률**: 0% (시작 전)

## 🔗 관련 문서

- [백엔드 역할 가이드](./BACKEND_ROLE.md)
- [게임 기획 문서](./PRD.md)
- [API 통합 가이드](./API_INTEGRATION_GUIDE.md)

## 📝 개발 노트

### 2025-06-25
- 프로젝트 초기 설정 완료
- Docker 환경 구성 및 자동화 스크립트 작성
- PostgreSQL, Redis 연결 확인
  - PostgreSQL: 테이블 생성 및 초기 데이터 마이그레이션 완료
  - Redis: 연결 테스트 및 기본 동작 확인
- API 문서화 도구 설정 (Swagger)
  - Swagger UI 정상 작동 확인: http://localhost:8080/swagger/index.html
  - API 문서 자동 생성 설정 완료
  - 프로젝트 맞춤형 문서 업데이트 (한국어 설명, Vibe 코딩 개념 반영)
- 프론트엔드 통합을 위한 TypeScript 타입 정의
- 개발 환경 자동화
  - `/rebuild` 명령어로 전체 시스템 재빌드
  - `/quick-rebuild` 명령어로 백엔드만 재빌드
- API 경로 표준화
  - 모든 API 엔드포인트를 `/api/v1` 하위로 통일
  - Health API: http://localhost:8080/api/v1/health
- CORS 설정 완료
  - 프론트엔드 개발 서버 (localhost:3000, localhost:5173) 지원
  - OPTIONS preflight 요청 처리
  - 인증 헤더 및 쿠키 지원 설정
- JWT 기반 인증 시스템 구현 완료
  - JWT 토큰 생성/검증 유틸리티 (internal/auth/jwt.go)
  - 비밀번호 해싱 (bcrypt) 유틸리티 (internal/auth/password.go)
  - 인증 미들웨어 (internal/middleware/auth.go)
  - 인증 API 엔드포인트 (internal/handlers/auth.go)
    - POST /api/v1/auth/register (회원가입)
    - POST /api/v1/auth/login (로그인)
    - POST /api/v1/auth/logout (로그아웃)
    - POST /api/v1/auth/refresh (토큰 갱신)
    - GET /api/v1/auth/profile (프로필 조회)
  - JWT 시크릿 키 환경변수 설정 (JWT_SECRET_KEY)
- PostgreSQL 기반 사용자 관리 시스템 구현 완료
  - DDD 패턴을 활용한 User 도메인 모델 (internal/domain/user.go)
  - Repository 패턴 기반 데이터 액세스 (internal/repository/postgres/user.go)
  - 데이터베이스 연결 및 환경 설정 (internal/database/connection.go)
  - 실제 DB 기반 사용자 인증 (중복 검사, 비밀번호 검증)
  - 사용자 프로필 및 통계 관리 API (internal/handlers/user.go)
    - PUT /api/v1/users/profile (프로필 수정)
    - GET /api/v1/users/stats (통계 조회)
    - GET /api/v1/users/collection (카드 컬렉션)
    - POST /api/v1/users/stats/* (통계 업데이트)
  - 플랫폼별 사용자 구분 (Web, Android, iOS)
  - 데이터베이스 마이그레이션 (002_user_system_update)
- 카드 시스템 구현 완료
  - Card 도메인 모델 (internal/domain/card.go)
    - 카드 타입: ACTION (액션), EVENT (이벤트), POWER (지속 효과)
    - 카드 희귀도: COMMON, RARE, EPIC, LEGENDARY
    - 카드 효과 및 시각 효과 JSON 데이터 구조
  - Card Repository 구현 (internal/repository/postgres/card.go)
    - 카드 마스터 데이터 CRUD
    - 사용자 카드 컬렉션 관리
    - 덱 시스템 (생성, 수정, 삭제, 활성화)
  - 카드 API 핸들러 (internal/handlers/card.go)
    - GET /api/v1/cards (카드 목록 - 필터, 페이지네이션 지원)
    - GET /api/v1/cards/:id (카드 상세)
    - GET /api/v1/cards/my-collection (내 컬렉션)
    - POST /api/v1/cards/decks (덱 생성)
    - GET /api/v1/cards/decks (내 덱 목록)
    - GET /api/v1/cards/decks/:id (덱 상세)
    - PUT /api/v1/cards/decks/:id (덱 수정)
    - DELETE /api/v1/cards/decks/:id (덱 삭제)
    - PUT /api/v1/cards/decks/:id/activate (덱 활성화)
    - GET /api/v1/cards/decks/active (활성 덱 조회)
  - 데이터베이스 마이그레이션
    - 003_card_system: 카드, 사용자 카드, 덱 테이블 생성
    - 004_seed_cards: 초기 카드 데이터 20개 시드
    - 005_remove_unique_card_constraint: 카드 중복 소유 허용
  - 신규 사용자 초기 카드 지급 시스템
    - 회원가입 시 13장의 스타터 카드 자동 지급
- 게임 플레이 시스템 구현 (Phase 6)
  - Game 도메인 모델 (internal/domain/game.go)
    - GameSession: 게임 세션 관리
    - PlayerState: 플레이어 상태 (체력, 에너지, 카드 등)
    - EnemyState: 적 상태 및 행동 의도
    - GameState: 전체 게임 진행 상태
  - Game Repository 구현 (internal/repository/postgres/game.go)
    - 게임 세션 CRUD 및 상태 관리
    - 게임 액션 기록 및 통계
    - JSONB를 활용한 복잡한 게임 상태 저장
  - Game Handler 구현 (internal/handlers/game.go)
    - POST /api/v1/games/start (게임 시작)
    - GET /api/v1/games/current (현재 게임 조회)
    - GET /api/v1/games/:id (특정 게임 조회)
    - POST /api/v1/games/:id/actions (카드 플레이 등 액션)
    - POST /api/v1/games/:id/end-turn (턴 종료)
    - POST /api/v1/games/:id/surrender (게임 포기)
    - GET /api/v1/games/stats (게임 통계)
  - 데이터베이스 마이그레이션
    - 006_game_system: 게임 세션 및 액션 테이블 생성
  - Docker 빌드 프로세스 개선
    - Swagger 문서 자동 생성 통합
    - 빌드 중 에러 발생 시에도 계속 진행하도록 설정
  - API 문서 업데이트
    - 모든 Phase 5, 6 API 엔드포인트 문서화
    - API_INTEGRATION_GUIDE.md 업데이트
    - Swagger docs.go 수동 작성 (자동 생성 이슈 해결 필요)

### 다음 마일스톤
- 카드 효과 실행 엔진 (진행 예정)
- 적 AI 행동 패턴 구현
- 실시간 통신 시스템 (WebSocket)
- 게임 밸런싱 시스템
- 실제 "Vibe 코딩" 카드 효과 실행 엔진 구현

---

*이 문서는 개발 진행 상황에 따라 지속적으로 업데이트됩니다.*