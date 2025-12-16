## ZAGE (Zero-Knowledge Age Verification)

Go로 작성한 영지식 기반 성인 인증 실험 프로젝트입니다. `(현재 연도 - 생년) >= 기준 나이` 관계만 공개하고, 실제 생년은 브라우저/클라이언트에만 남겨둔 채 Groth16 증명을 생성합니다.

### 디렉터리 구조

- `internal/age`: AgeCircuit 정의와 공통 헬퍼(컴파일, witness 생성 등)
- `main.go`: CLI 유틸리티. 회로 컴파일 → PK/VK 생성 → 샘플 증명/검증을 수행하고 `zage.pk`, `zage.vk` 파일을 만듭니다.
- `web/`: 브라우저(WebAssembly)용 프로버. `index.html`, `main.wasm`, `wasm_exec.js`, `zage.pk`를 함께 서빙하면 됩니다.

### CLI 실행

```bash
go run ./...
```

실행하면 루트 디렉터리에 `zage.pk`, `zage.vk`가 생성되고, 예시 입력(2025년/만 19세/2005년생)에 대해 증명과 검증이 한 번씩 수행됩니다.

### WebAssembly 빌드 & 실행

1. Wasm 바이너리 생성
   ```bash
   GOENV=$(pwd)/go.env GOOS=js GOARCH=wasm go build -o web/main.wasm ./web
   ```
2. 런타임 스크립트 복사
   ```bash
   cp "$(GOENV=$(pwd)/go.env go env GOROOT)/lib/wasm/wasm_exec.js" web/
   ```
3. 정적 서버에서 `web/` 디렉터리 서빙
   ```bash
   python3 -m http.server 8080 --directory web
   # 브라우저에서 http://localhost:8080/index.html 접속
   ```

페이지에서 연도를 입력하고 **증명서 생성하기** 버튼을 누르면 `verifyAgeZKP` Go 함수가 실행되어 Hex 인코딩된 proof가 표시됩니다.

### 개발 팁

- `go.env`에는 프로젝트 로컬 캐시 경로(`.cache/go-build`)가 정의되어 있으므로, Go 관련 명령을 실행할 때 `GOENV=$(pwd)/go.env`를 붙이면 권한 문제를 피할 수 있습니다.
- WebAssembly 빌드는 Go 표준라이브러리 Wasm 지원이 필요하므로 반드시 `wasm_exec.js`를 최신 GOROOT에서 복사해 주세요.
