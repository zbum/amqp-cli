# 배포 가이드

## 빠른 배포

```bash
# 1. 코드 수정 후 릴리스 빌드
make release VERSION=0.2.0

# 2. Formula 업데이트 (SHA256 복사해서 수정)
vi Formula/amqp-cli.rb

# 3. 커밋 & 배포
git add . && git commit -m "Release v0.2.0" && git push
make gh-release VERSION=0.2.0
make update-tap VERSION=0.2.0
```

## 단계별 상세 가이드

### 1. 릴리스 빌드

```bash
make release VERSION=0.2.0
```

출력 예시:
```
=== SHA256 checksums ===
abc123...  amqp-cli-darwin-amd64.tar.gz
def456...  amqp-cli-darwin-arm64.tar.gz
ghi789...  amqp-cli-linux-amd64.tar.gz
jkl012...  amqp-cli-linux-arm64.tar.gz
mno345...  amqp-cli-windows-amd64.zip
```

### 2. Formula 업데이트

`Formula/amqp-cli.rb` 파일 수정:

```ruby
version "0.2.0"  # 버전 업데이트

on_macos do
  if Hardware::CPU.arm?
    sha256 "def456..."  # darwin-arm64 SHA256
  else
    sha256 "abc123..."  # darwin-amd64 SHA256
  end
end

on_linux do
  if Hardware::CPU.arm?
    sha256 "jkl012..."  # linux-arm64 SHA256
  else
    sha256 "ghi789..."  # linux-amd64 SHA256
  end
end
```

### 3. 커밋 & Push

```bash
git add .
git commit -m "Release v0.2.0"
git push
```

### 4. GitHub 릴리스 생성

```bash
make gh-release VERSION=0.2.0
```

또는 수동으로:
```bash
gh release create v0.2.0 dist/release/* --title "v0.2.0" --notes "Release notes"
```

### 5. Homebrew Tap 업데이트

```bash
make update-tap VERSION=0.2.0
```

또는 수동으로:
```bash
git clone https://github.com/zbum/homebrew-tap.git /tmp/homebrew-tap
cp Formula/amqp-cli.rb /tmp/homebrew-tap/Formula/
cd /tmp/homebrew-tap
git add . && git commit -m "Update amqp-cli to v0.2.0" && git push
```

## Makefile 타겟 요약

| 타겟 | 설명 |
|------|------|
| `make release VERSION=x.x.x` | 빌드 + tarball 생성 + SHA256 출력 |
| `make gh-release VERSION=x.x.x` | GitHub 릴리스 생성 |
| `make update-tap VERSION=x.x.x` | homebrew-tap 레포 업데이트 |
| `make publish VERSION=x.x.x` | 위 3개 한번에 실행 (Formula 수정 필요) |

## 사용자 업그레이드

```bash
brew update
brew upgrade amqp-cli
```
