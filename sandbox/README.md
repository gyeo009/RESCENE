# Sandbox Backend

이 프로젝트는 Docker 기반의 격리된 Sandbox 환경을 관리하기 위한 백엔드 애플리케이션의 초기 스켈레톤 및 핵심 레이어 설계본입니다. Gin 프레임워크와 Swagger 문서화, 그리고 호스트의 Docker 데몬 제어를 위한 Docker SDK 연동이 완료되어 있습니다.

---

## 1. 프로젝트 구조 (Directory Structure)

Go의 일반적인 레이어드 아키텍처 관례를 준수하며, 불필요한 추상화 없이 계층별 책임을 명확히 설계했습니다.

```text
sandbox/
├── cmd/
│   └── server/
│       └── main.go       # 애플리케이션 엔트리 포인트 (설정 및 의존성 주입 수행)
│
├── internal/
│   ├── config/
│   │   └── config.go     # 포트 등 환경설정 로드 레이어
│   ├── router/
│   │   └── router.go     # Gin 라우팅 설정 및 Swagger 연동
│   ├── handler/
│   │   └── sandbox.go    # HTTP 요청 수신, 바인딩 및 응답 처리 (Controller)
│   ├── service/
│   │   └── sandbox.go    # 비즈니스 로직 조율 레이어
│   ├── docker/
│   │   ├── client.go     # Docker 인터페이스 정의
│   │   ├── types.go      # 패키지 공통 데이터 타입 구조체 정의
│   │   ├── mock.go       # 로컬/테스트용 모의 Docker 클라이언트
│   │   └── docker.go     # 실제 Docker SDK 연동 모듈 (ContainerCreate -> ContainerStart)
│   ├── dto/
│   │   └── sandbox.go    # API Request 구조체 정의
│   └── response/
│       └── sandbox.go    # API Response 구조체 정의
│
├── docs/                 # Swaggo로 자동 생성된 Swagger API Spec 파일들
│
├── Dockerfile            # Golang 기반 개발 컨테이너 이미지 사양서
├── docker-compose.yml    # Docker 소켓 볼륨 마운트 및 포트 맵핑이 적용된 로컬 멀티 컨테이너 정의서
└── .air.toml             # 파일 변경 시 빌드 및 컨테이너 리로드를 제공하는 Hot Reload 설정 파일
```

---

## 2. API 사용 방법 (API Usage Guide)

### 2.1 샌드박스 구동 API

* **메서드 및 경로**: `POST /sandbox/run`
* **요청 헤더**: `Content-Type: application/json`

**요청 예시 (Request Body)**
```json
{
    "scenario": "sandbox-app:latest",
    "user": "test-user-2"
}
```
> **주의**: Docker SDK 실제 모드가 활성화되어 있으므로, `scenario` 값은 호스트 컴퓨터의 Docker Engine에 실제로 존재하는 이미지명(예: `sandbox-app:latest`)이어야 컨테이너 기동에 성공합니다.

**성공 응답 예시 (Response Body - 200 OK)**
```json
{
    "container_id": "56dc3cf5965b8876fb1797760f6f0a1ef329787979b4e698fa129446de68475e",
    "status": "running"
}
```

---

### 2.2 Swagger API 문서화 접속

Swagger UI를 지원하며 아래 URL을 통해 전체 API 스펙과 파라미터 구조를 직관적으로 확인하고 직접 테스트해 볼 수 있습니다.

* **Swagger 접속 주소**: `http://127.0.0.1:8080/swagger/index.html`

---

## 3. 개발환경 구동 및 빌드 방법

### 3.1 Docker Compose를 통한 기동 (Air 실시간 리로드 포함)

프로젝트 루트 디렉토리에서 아래 명령어를 실행하면, 호스트의 Docker 소켓 마운트와 실시간 소스코드 감지 빌드가 적용된 백엔드가 구동됩니다.

```bash
# 컨테이너 빌드 및 백그라운드 구동
docker-compose up -d --build

# 실행 로그 및 빌드 진행 상태 실시간 모니터링
docker logs -f sandbox-app-1
```

### 3.2 로컬 수동 빌드 검증

```bash
# 로컬 바이너리 컴파일 검증
go build -o tmp/main ./cmd/server

# Swagger 스펙 수동 갱신 (API 변경 시)
go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/server/main.go -o docs
```

---

## 4. 실제 API 실행 테스트 (Execution Test History)

서버가 띄워진 상태에서 아래 PowerShell 또는 curl 명령어로 실행 내역을 확인할 수 있습니다.

### PowerShell을 이용한 테스트 요청
```powershell
Invoke-RestMethod -Method Post -Uri http://127.0.0.1:8080/sandbox/run -ContentType "application/json" -Body '{"scenario":"sandbox-app:latest","user":"test-user-2"}' | ConvertTo-Json
```

### curl을 이용한 테스트 요청
```bash
curl -X POST http://127.0.0.1:8080/sandbox/run \
  -H "Content-Type: application/json" \
  -d '{"scenario":"sandbox-app:latest","user":"test-user-2"}'
```

**호스트 Docker Engine 컨테이너 감지 검증**
요청 직후 터미널에 `docker ps -a`를 입력하면 백엔드 서버에서 실행한 `sandbox-test-user-2` 컨테이너가 성공적으로 부팅되어 작동하고 있는 것을 확인할 수 있습니다.

```text
CONTAINER ID   IMAGE                COMMAND   CREATED         STATUS         PORTS   NAMES
56dc3cf5965b   sandbox-app:latest   "air"     6 seconds ago   Up 5 seconds           sandbox-test-user-2
```
