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

### 🚧 Phase 5: 카드 시스템 구현 (다음 목표)
- [ ] 카드 도메인 모델 및 Repository 구현
- [ ] 카드 목록 조회 API (`GET /api/v1/cards`)
- [ ] 카드 상세 조회 API (`GET /api/v1/cards/:id`)
- [ ] 사용자 카드 컬렉션 시스템 구현
- [ ] 카드 효과 처리 엔진 구현
- [ ] 카드 업그레이드 시스템
- [ ] 카드 시너지 계산 로직

### 📅 Phase 6: 게임 플레이 시스템
- [ ] 게임 시작 API (`POST /api/v1/games/start`)
- [ ] 게임 상태 조회 API (`GET /api/v1/games/:id`)
- [ ] 카드 플레이 API (`POST /api/v1/games/:id/actions`)
- [ ] 게임 종료 API (`POST /api/v1/games/:id/end`)
- [ ] 턴제 전투 시스템 구현
- [ ] 적 AI 행동 패턴 구현

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

### Phase 5: 카드 시스템 구현
실제 게임 카드 데이터와 사용자 카드 컬렉션을 관리하는 시스템을 구현합니다.

**다음 작업:**
1. Card 도메인 모델과 Repository 패턴 구현
2. 카드 마스터 데이터 관리 시스템
3. 사용자별 카드 컬렉션 시스템
4. 카드 효과와 코드 실행 엔진 기초 구현

## 📊 진행률

- **전체 진행률**: 40% (4/10 Phase 완료)
- **현재 Phase 진행률**: 0% (Phase 5 시작)

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

### 다음 마일스톤
- 카드 시스템 및 컬렉션 관리 구현
- 게임 플레이 세션 및 실시간 통신 시스템
- 실제 "Vibe 코딩" 카드 효과 실행 엔진 구현

---

*이 문서는 개발 진행 상황에 따라 지속적으로 업데이트됩니다.*