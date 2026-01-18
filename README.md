# amqp-cli

RabbitMQ 메시지를 전송하고 소비할 수 있는 CLI 도구입니다.

## 설치

### Homebrew (macOS/Linux)

```bash
brew tap zbum/homebrew
brew install amqp-cli
```

### 소스에서 빌드

```bash
git clone https://github.com/zbum/amqp-cli.git
cd amqp-cli
make build
```

빌드된 바이너리는 `dist/amqp-cli`에 생성됩니다.

### 크로스 플랫폼 빌드

모든 플랫폼용 바이너리를 한번에 빌드합니다:

```bash
make build-all
```

빌드 결과:

| 파일명 | OS | Architecture |
|--------|-----|--------------|
| `amqp-cli-linux-amd64` | Linux | x86_64 |
| `amqp-cli-linux-arm64` | Linux | ARM64 |
| `amqp-cli-darwin-amd64` | macOS | Intel |
| `amqp-cli-darwin-arm64` | macOS | Apple Silicon |
| `amqp-cli-windows-amd64.exe` | Windows | x86_64 |

### Go install

```bash
go install github.com/zbum/amqp-cli@latest
```

## 사용법

### 공통 옵션

| 옵션 | 단축 | 기본값 | 설명 |
|------|------|--------|------|
| `--host` | `-H` | localhost | RabbitMQ 호스트 |
| `--port` | `-P` | 5672 | RabbitMQ 포트 |
| `--username` | `-u` | guest | RabbitMQ 사용자명 |
| `--password` | `-p` | guest | RabbitMQ 비밀번호 |
| `--vhost` | `-v` | (empty) | RabbitMQ 가상 호스트 |

### 메시지 발행 (publish)

큐 또는 익스체인지에 메시지를 발행합니다.

```bash
# 큐에 직접 발행
amqp-cli publish -q myqueue -m "Hello World"

# 익스체인지에 발행 (라우팅 키 지정)
amqp-cli publish -e myexchange -r mykey -m "Hello World"

# 원격 서버에 연결
amqp-cli publish -H rabbitmq.example.com -u admin -p secret -q myqueue -m "Hello"

# stdin에서 메시지 읽기
echo "Hello World" | amqp-cli publish -q myqueue
cat message.json | amqp-cli publish -q myqueue
```

#### Publish 옵션

| 옵션 | 단축 | 설명 |
|------|------|------|
| `--queue` | `-q` | 발행할 큐 이름 |
| `--exchange` | `-e` | 익스체인지 이름 |
| `--routing-key` | `-r` | 라우팅 키 |
| `--message` | `-m` | 메시지 본문 |

### 메시지 소비 (consume)

큐에서 메시지를 소비합니다.

```bash
# 큐에서 메시지 소비 (지속적)
amqp-cli consume -q myqueue

# 자동 확인 모드로 소비
amqp-cli consume -q myqueue --auto-ack

# 지정된 개수만 소비
amqp-cli consume -q myqueue -n 10

# Verbose 모드 (Method/Header/Body 프레임 상세 출력)
amqp-cli consume -q myqueue -V

# Hex dump 모드 (바이너리 메시지 디버깅)
amqp-cli consume -q myqueue --hex

# Verbose + Hex dump 조합
amqp-cli consume -q myqueue -V --hex

# 원격 서버에서 소비
amqp-cli consume -H rabbitmq.example.com -u admin -p secret -q myqueue
```

#### Consume 옵션

| 옵션 | 단축 | 기본값 | 설명 |
|------|------|--------|------|
| `--queue` | `-q` | (필수) | 소비할 큐 이름 |
| `--auto-ack` | | false | 자동 확인 모드 |
| `--count` | `-n` | 0 | 소비할 메시지 수 (0=무제한) |
| `--verbose` | `-V` | false | Method/Header/Body 프레임 상세 출력 |
| `--hex` | | false | Body를 hex dump로 출력 |

## 예제

### Docker로 RabbitMQ 실행

```bash
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:management
```

### 기본 발행/소비 테스트

터미널 1 (소비자):
```bash
amqp-cli consume -q test-queue
```

터미널 2 (발행자):
```bash
amqp-cli publish -q test-queue -m "Hello RabbitMQ!"
```

### JSON 메시지 발행

```bash
amqp-cli publish -q events -m '{"event":"user_created","user_id":123}'
```

## 개발

### 빌드

```bash
make build       # 현재 플랫폼용 빌드
make build-all   # 모든 플랫폼용 빌드
```

### 테스트

```bash
make test
```

### 코드 포맷팅

```bash
make fmt
```

### 린트

```bash
make lint
```

## 프로젝트 구조

```
.
├── cmd/
│   ├── root.go      # 루트 명령어 및 공통 플래그
│   ├── publish.go   # publish 명령어
│   └── consume.go   # consume 명령어
├── internal/
│   └── rabbitmq/
│       └── client.go # RabbitMQ 클라이언트
├── main.go
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## 라이선스

MIT License
