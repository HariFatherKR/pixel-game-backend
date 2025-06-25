# 백엔드 개발자 역할 가이드

## 🎯 역할 개요
사이버펑크 덱 빌딩 카드 게임의 Go 백엔드 서버 개발을 담당합니다. 게임 로직, API 설계, 데이터베이스 관리, 실시간 통신 구현이 주요 업무입니다.

## 🛠️ 기술 스택
- **언어**: Go 1.21+
- **웹 프레임워크**: Gin
- **데이터베이스**: PostgreSQL
- **캐시**: Redis
- **실시간 통신**: Gorilla WebSocket
- **인증**: JWT
- **컨테이너**: Docker

## 📁 프로젝트 구조

```
backend/
├── cmd/
│   ├── server/main.go      # 서버 진입점
│   └── migrate/            # 마이그레이션 도구
├── internal/               # 내부 패키지 (DDD 패턴)
│   ├── config/             # 환경 설정
│   ├── domain/             # 도메인 계층
│   │   ├── entity/         # 엔티티 (User, Card, Game)
│   │   ├── valueobject/    # 값 객체
│   │   └── repository/     # 리포지토리 인터페이스
│   ├── application/        # 애플리케이션 계층
│   │   ├── usecase/        # 유스케이스
│   │   └── service/        # 도메인 서비스
│   ├── infrastructure/     # 인프라 계층
│   │   ├── persistence/    # DB 구현
│   │   ├── cache/          # Redis 구현
│   │   └── messaging/      # WebSocket
│   └── interfaces/         # 인터페이스 계층
│       ├── http/           # REST API
│       └── websocket/      # 실시간 통신
├── pkg/                    # 공개 패키지
├── migrations/             # DB 마이그레이션
└── deployments/            # 배포 설정
```

## 🗂️ 현재 구현 상태

### ✅ 완료된 작업 (Phase 1 - 2024.06.25)

1. **프로젝트 구조 설정**
   - DDD(Domain-Driven Design) 패턴 적용
   - 마이크로서비스 전환 가능한 구조 설계
   - 2024년 Go 게임 서버 모범 사례 반영

2. **데이터베이스 스키마**
   - PostgreSQL 마이그레이션 파일 작성
   - 주요 테이블: users, cards, user_cards, game_sessions, leaderboard
   - 크로스 플랫폼(Android/iOS/Web) 지원 고려

3. **도메인 모델**
   - Card 엔티티: 카드 타입, 효과, 시각 효과 정의
   - User 엔티티: 플랫폼별 사용자 관리
   - Game 엔티티: 게임 세션 및 상태 관리

4. **개발 환경**
   - 환경 설정 시스템 (config/config.go)
   - Makefile 명령어 추가
   - .env.example 제공

### 🚧 진행 중인 작업 (Phase 2)

1. **API 구현**
   - JWT 인증 시스템
   - RESTful API 엔드포인트
   - 미들웨어 (CORS, 로깅, 에러 핸들링)

2. **게임 로직**
   - 카드 효과 처리 엔진
   - 턴제 전투 시스템
   - 덱 빌딩 메커니즘

## 📋 주요 책임 사항

### 1. API 설계 및 구현
- RESTful API 설계 원칙 준수
- 명확한 엔드포인트 네이밍
- 적절한 HTTP 상태 코드 사용
- 요청/응답 검증

### 2. 게임 로직 구현
- 카드 효과 시스템 설계
- 전투 메커니즘 구현
- 게임 상태 관리
- 서버 사이드 검증 (치팅 방지)

### 3. 데이터베이스 관리
- 스키마 설계 및 최적화
- 마이그레이션 관리
- 인덱스 최적화
- 쿼리 성능 튜닝

### 4. 실시간 통신
- WebSocket 연결 관리
- 게임 상태 동기화
- 메시지 프로토콜 설계
- 연결 오류 처리

### 5. 보안
- JWT 기반 인증
- API 권한 관리
- 입력 검증 및 sanitization
- SQL 인젝션 방지

## 🔧 개발 가이드라인

### 코드 컨벤션
```go
// 패키지명: 소문자, 단수형
package entity

// 구조체: PascalCase
type GameSession struct {
    ID     uuid.UUID
    UserID uuid.UUID
}

// 인터페이스: ~er 접미사
type CardHandler interface {
    PlayCard(cardID string) error
}

// 에러 처리: errors.Wrap 사용
if err != nil {
    return errors.Wrap(err, "failed to play card")
}
```

### 프로젝트 규칙
1. **단일 책임 원칙**: 각 함수/구조체는 하나의 책임만
2. **의존성 주입**: 인터페이스를 통한 의존성 관리
3. **테스트 우선**: 유닛 테스트 작성 필수
4. **문서화**: 공개 API는 반드시 문서화

## 🔗 협업 방식

### 프론트엔드 팀과의 협업
- API 스펙 문서 제공 (OpenAPI/Swagger)
- WebSocket 메시지 프로토콜 문서화
- CORS 설정 및 개발 환경 지원

### PM과의 소통
- 일일 진행 상황 업데이트
- 기술적 이슈 및 해결 방안 공유
- 일정 관리 및 우선순위 조정

### 코드 리뷰
- PR 생성 시 상세한 설명 포함
- 테스트 코드 필수
- 성능 영향 검토

## 📚 참고 자료

### 필수 문서
- `/docs/PRD.md` - 게임 기획 문서
- `/docs/BACKEND_TASKS.md` - 개발 태스크 목록
- `/docs/GAME_DESIGN_DATA.md` - 게임 데이터 정의

### 외부 참고
- [Gin Web Framework](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [Gorilla WebSocket](https://github.com/gorilla/websocket)
- [JWT-Go](https://github.com/golang-jwt/jwt)

## 🚀 시작하기

### 개발 환경 설정
```bash
# 1. 환경 변수 설정
cp backend/.env.example backend/.env

# 2. Docker 서비스 시작
make dev

# 3. 데이터베이스 마이그레이션
make backend-migrate
make backend-seed

# 4. 서버 실행
make backend-run
```

### 주요 Make 명령어
- `make backend-run` - 서버 실행
- `make backend-build` - 빌드
- `make backend-test` - 테스트 실행
- `make backend-migrate` - DB 마이그레이션
- `make backend-migrate-down` - 마이그레이션 롤백

## ⚠️ 주의사항

1. **보안 우선**: 모든 사용자 입력은 검증
2. **성능 고려**: 쿼리 최적화 및 캐싱 활용
3. **에러 처리**: 명확한 에러 메시지와 로깅
4. **버전 관리**: 마이그레이션 순서 주의
5. **테스트**: 새 기능은 반드시 테스트 작성

## 📈 향후 계획

### Phase 2 (진행 중)
- JWT 인증 시스템
- 기본 API 엔드포인트
- 카드 관리 시스템

### Phase 3
- 게임 로직 구현
- 전투 시스템
- 카드 효과 엔진

### Phase 4
- WebSocket 실시간 통신
- 게임 상태 동기화
- 관전 모드

### Phase 5
- 리더보드
- 일일 챌린지
- 이벤트 시스템

---

이 문서는 백엔드 개발자가 프로젝트를 이해하고 개발을 진행하는 데 필요한 모든 정보를 담고 있습니다. 새로운 개발자가 합류하거나 인수인계가 필요한 경우 이 문서를 참고하여 빠르게 프로젝트에 적응할 수 있습니다.