# Pixel Game Backend

사이버펑크 덱 빌딩 카드 게임의 Go 백엔드 서버입니다.

## 기술 스택

- **언어**: Go 1.21+
- **웹 프레임워크**: Gin
- **데이터베이스**: PostgreSQL
- **캐시**: Redis
- **컨테이너**: Docker & Docker Compose

## 시작하기

### 필수 요구사항

- Docker Desktop 설치
- Docker Compose 설치
- Make (선택사항)

### 개발 환경 실행

1. Docker Desktop이 실행 중인지 확인합니다.

2. 프로젝트 디렉토리로 이동합니다:
```bash
cd pixel-game/backend
```

3. Docker 컨테이너를 시작합니다:
```bash
make dev
```

또는 make가 없는 경우:
```bash
docker compose up -d
docker compose run --rm migrate
```

4. 서버가 정상적으로 실행되는지 확인합니다:
```bash
curl http://localhost:8080/health
```

### 개발 명령어

```bash
# 모든 서비스 시작
make docker-up

# 모든 서비스 중지
make docker-down

# 로그 확인
make docker-logs

# 데이터베이스 마이그레이션 실행
make migrate-up

# 마이그레이션 롤백
make migrate-down
```

### 환경 변수

`.env` 파일을 생성하여 환경 변수를 설정합니다:

```env
DATABASE_URL=postgres://pixelgame:pixelgame123@localhost:5432/pixelgame_db?sslmode=disable
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-secret-key-here
PORT=8080
```

## API 엔드포인트

- `GET /health` - 헬스 체크
- 추가 엔드포인트는 개발 진행에 따라 업데이트 예정

## 프로젝트 구조

```
backend/
├── cmd/
│   └── server/main.go      # 서버 진입점
├── internal/               # 내부 패키지
│   ├── config/             # 환경 설정
│   └── domain/             # 도메인 계층
│       └── entity/         # 엔티티
├── migrations/             # DB 마이그레이션
├── docker-compose.yml      # Docker Compose 설정
├── Dockerfile              # Docker 이미지 빌드
└── Makefile                # 개발 명령어
```

## 문제 해결

### Docker 연결 오류
```
Cannot connect to the Docker daemon
```
→ Docker Desktop이 실행 중인지 확인하세요.

### 포트 충돌
```
bind: address already in use
```
→ 다른 서비스가 같은 포트를 사용 중입니다. docker-compose.yml에서 포트를 변경하세요.

## 참고 문서

- [백엔드 개발자 역할 가이드](docs/BACKEND_ROLE.md)
- [게임 기획 문서](docs/PRD.md)