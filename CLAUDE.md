# Claude Code Instructions for Pixel Game Backend

## 🚀 Quick Commands

### Docker Rebuild & Restart
- **Full Rebuild** (모든 서비스 재시작): `/rebuild`
  ```bash
  make rebuild
  ```
  
- **Quick Rebuild** (백엔드만 재빌드): `/quick-rebuild`
  ```bash
  make quick-rebuild
  ```

### Common Development Commands
- **Start Services**: `make docker-up`
- **Stop Services**: `make docker-down`
- **View Logs**: `make docker-logs`
- **Check Health**: `curl http://localhost:8080/health`

## 📝 Development Workflow

1. **코드 수정 후 재빌드가 필요한 경우**:
   - 간단한 백엔드 코드 변경: `/quick-rebuild`
   - 전체 시스템 재시작 필요: `/rebuild`

2. **API 테스트**:
   - Swagger UI: http://localhost:8080/swagger/index.html
   - Health Check: `curl http://localhost:8080/health | jq`

3. **데이터베이스 작업**:
   - 마이그레이션 실행: `make migrate-up`
   - 마이그레이션 롤백: `make migrate-down`

## 🛠️ Troubleshooting

### Docker 관련 문제
- Docker Desktop이 실행 중인지 확인
- 포트 충돌 확인 (8080, 5432, 6379)
- 컨테이너 상태 확인: `docker compose ps`

### 빌드 실패 시
1. go.mod 의존성 확인
2. Docker 캐시 정리: `docker compose build --no-cache backend`
3. 로그 확인: `docker compose logs backend`

## 📁 Important Files
- `/cmd/server/main.go` - 메인 서버 엔트리포인트
- `/api/types/api.types.ts` - TypeScript 타입 정의
- `/docker-compose.yml` - Docker 서비스 설정
- `/Makefile` - 개발 명령어 모음

## 🔍 API Endpoints
- GET `/health` - 헬스 체크
- GET `/swagger/*` - Swagger 문서
- GET `/api/v1/cards` - 카드 목록
- GET `/api/v1/version` - 버전 정보

## 💡 Tips
- 코드 변경 후 항상 `/quick-rebuild` 실행
- API 변경 시 Swagger 문서 업데이트 필요
- TypeScript 타입 동기화 유지 중요