# Git Resume Analyzer (Go + Claude)

**Go 언어와 LLM을 활용한 이력서 성과 자동 추출기**

## 1. 프로젝트 목적

- 정기적으로 Git 커밋 내역을 분석하여 이력서용 '성과 중심 문구' 생성
- **Go**의 성능과 **Claude**의 문장력을 결합한 효율적인 자동화 도구 구축
- 불필요한 비용 발생을 막기 위한 **토큰 최적화** 프로세스 정립

## 2. 핵심 기술 스택

| 구분         | 선택 기술              | 이유                                                                   |
| :----------- | :--------------------- | :--------------------------------------------------------------------- |
| **Language** | **Go (Golang)**        | 병렬 처리(Goroutine)를 통한 빠른 API 호출 및 단일 바이너리 관리 편의성 |
| **Git SDK**  | `go-git/go-git`        | 로컬 리포지토리의 커밋 객체에 대한 정교한 제어 가능                    |
| **CLI**      | `spf13/cobra`          | `analyze --month=4` 와 같은 직관적인 명령어 인터페이스 구축            |
| **Database** | `SQLite (modernc.org)` | 처리된 커밋 해시를 저장하여 중복 분석 방지 (CGO-free)                  |
| **LLM**      | `Claude API`           | 문맥 이해도가 높고 이력서용 톤앤매너 구현에 탁월                       |
| **Config**   | `spf13/viper`          | 환경변수 및 설정 파일 통합 관리                                        |

## 3. 토큰 최적화 & 단계별 실행 계획

### [Step 1] 로컬 전처리 (Zero Token Cost)

- **필터링**: `go.sum`, `*.lock`, `vendor/`, 이미지 파일 등 로직과 무관한 변경사항 자동 제외
- **메시지 분석**: "Fix typo", "Merge branch" 등 영양가 없는 커밋은 LLM에 던지지 않고 폐기
- **Diff 압축**: 코드 전체 대신 `git show --stat` 기반의 파일 변경 리스트와 핵심 함수 변경점만 추출
- **중요도 스코어링**: 변경된 파일 수, 라인 수, 파일 유형 기반으로 커밋 중요도 사전 평가

### [Step 2] 증분 분석 시스템 (Caching)

- **SQLite 연동**: 분석이 완료된 커밋 해시는 DB에 기록
- **중복 제거**: 실행 시 DB에 없는 '신규 커밋'만 추출하여 API 호출 횟수 최소화
- **토큰 사용량 추적**: 각 분석 실행 시 예상/실제 토큰 사용량 로깅

### [Step 3] LLM 요약 전략 (Prompt Engineering)

- **Batching**: 사소한 커밋 여러 개를 묶어서(Batch) 하나의 프롬프트로 처리
- **Role Playing**: "IT 전문 테크니컬 라이터" 페르소나 부여
- **Constraint**:
  - STAR(Situation, Task, Action, Result) 기법 적용 요청
  - 출력 포맷을 **JSON**으로 강제하여 파싱 오류 방지 및 출력 토큰 절약
- **Rate Limiting**: API 호출 제한 대응 (429 에러 방지)
- **재시도 로직**: Exponential backoff로 API 실패 시 안정적 재시도

### [Step 4] 최종 결과물 생성

- **다중 포맷 지원**:
  - CSV: 스프레드시트 분석용
  - Markdown: 이력서에 바로 복사 가능
  - JSON: 다른 도구와 연동용
- **Columns**: `Date`, `Project`, `Category(Feature/Refactor/Fix)`, `Impact_Summary`, `Commit_Hash`

## 4. CLI 명령어 인터페이스

```bash
# 기본 분석 (이번 달)
git-resume analyze

# 특정 월 분석
git-resume analyze --month=4

# 기간 범위 지정
git-resume analyze --from=2024-01-01 --to=2024-03-31

# 다중 레포지토리 분석
git-resume analyze --repos=/path/to/repo1,/path/to/repo2

# 템플릿 지정 (회사/직군별 톤 조절)
git-resume analyze --template=startup
git-resume analyze --template=backend

# 드라이런 모드 (API 호출 없이 테스트)
git-resume analyze --dry-run

# Slack 알림과 함께 분석
git-resume analyze --notify

# 출력 포맷 지정
git-resume export --format=markdown
git-resume export --format=json
git-resume export --format=csv --output=resume.csv

# 비용 추정
git-resume estimate --month=4

# 사용 가능한 템플릿 목록
git-resume templates
```

## 5. 디렉토리 구조

```text
.
├── .github/
│   └── workflows/
│       └── monthly-resume.yml  # GitHub Actions 자동화
├── cmd/
│   ├── root.go           # Cobra root 명령어
│   ├── analyze.go        # analyze 서브커맨드
│   ├── export.go         # export 서브커맨드
│   ├── estimate.go       # 비용 추정 커맨드
│   └── templates.go      # 템플릿 목록 커맨드
├── internal/
│   ├── git/
│   │   ├── parser.go     # 커밋 파싱
│   │   ├── filter.go     # 노이즈 커밋 필터링
│   │   ├── scorer.go     # 커밋 중요도 스코어링
│   │   └── codeparser.go # 언어별 코드 파서
│   ├── llm/
│   │   ├── client.go     # Claude API 클라이언트
│   │   ├── prompt.go     # 프롬프트 템플릿
│   │   ├── retry.go      # 재시도 로직 (exponential backoff)
│   │   └── templates.go  # 커스텀 템플릿 시스템
│   ├── db/
│   │   ├── sqlite.go     # SQLite 연결 관리
│   │   └── cache.go      # 캐시 CRUD
│   ├── export/
│   │   ├── csv.go        # CSV 출력
│   │   ├── markdown.go   # Markdown 출력
│   │   ├── json.go       # JSON 출력
│   │   ├── common.go     # 공통 유틸리티
│   │   └── exporter.go   # Exporter 인터페이스
│   ├── errors/
│   │   └── errors.go     # 커스텀 에러 타입
│   ├── notify/
│   │   └── slack.go      # Slack 알림
│   └── ui/
│       └── progress.go   # 프로그레스 바
├── pkg/
│   └── models/
│       ├── commit.go     # 커밋 데이터 구조체
│       └── result.go     # 분석 결과 구조체
├── scripts/
│   └── migrate.sql       # DB 스키마 마이그레이션
├── .env.example          # 환경변수 템플릿
├── .gitignore
├── Makefile              # 빌드/테스트/린트 명령어
├── main.go               # 엔트리 포인트
├── go.mod
└── README.md
```

## 6. 구현 로드맵

### Phase 1: MVP (핵심 기능) ✅

- [x] 프로젝트 초기 설정 (go.mod, Makefile, .env.example)
- [x] Git 커밋 파싱 및 필터링 (`internal/git/`)
- [x] SQLite 캐싱 시스템 (`internal/db/`)
- [x] Claude API 연동 - 단일 호출 (`internal/llm/`)
- [x] CSV 출력 (`internal/export/csv.go`)
- [x] 기본 CLI 구현 (`cmd/analyze.go`)

### Phase 2: 안정화 ✅

- [x] Rate limiting + 재시도 로직
- [x] 드라이런 모드 (`--dry-run`)
- [x] 에러 핸들링 강화 (`internal/errors/`)
- [x] 프로그레스 바 (장시간 분석 시 UX)
- [x] 단위 테스트 작성

### Phase 3: 확장 ✅

- [x] Markdown/JSON 출력 포맷 추가
- [x] 다중 레포지토리 지원 (`--repos`)
- [x] 기간 범위 지정 (`--from`, `--to`)
- [x] 비용 추정 커맨드 (`estimate`)
- [x] Export 서브커맨드 (`export`)

### Phase 4: 고도화 ✅

- [x] 템플릿 커스터마이징 (회사/직군별 톤 조절)
- [x] GitHub Actions 연동 (월말 자동 실행)
- [x] Slack 알림 연동
- [x] 언어별 코드 파서 (Go/Python/JS/TS/Java/Rust 함수 변경점 추출)

## 7. 환경 설정

```bash
# .env.example
CLAUDE_API_KEY=your_api_key_here
DEFAULT_REPO_PATH=/path/to/your/repo
DB_PATH=./data/cache.db
LOG_LEVEL=info
OUTPUT_DIR=./output
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/WEBHOOK/URL
```

## 8. 예상 비용 (Claude API)

| 항목 | 예상치 |
|------|--------|
| 평균 커밋/월 | ~50개 |
| 필터링 후 유의미한 커밋 | ~20개 |
| 배치 처리 (5개씩) | 4회 API 호출 |
| 예상 입력 토큰 | ~2,000 tokens/batch |
| 예상 출력 토큰 | ~500 tokens/batch |
| **월 예상 비용** | **~$0.50 이하** |

## 9. 템플릿 시스템

7개의 내장 템플릿이 제공되며, JSON 파일로 커스텀 템플릿 추가 가능:

| 템플릿 | 설명 | 톤 |
|--------|------|-----|
| `default` | 범용 (대부분의 기술직에 적합) | 전문적, 간결, 성과 중심 |
| `startup` | 스타트업 환경 강조 | 역동적, 결과 지향, 기업가적 |
| `enterprise` | 대기업/엔터프라이즈 | 공식적, 프로세스 중심, 컴플라이언스 |
| `backend` | 백엔드/인프라 엔지니어 | 기술적, 정밀, 시스템 중심 |
| `frontend` | 프론트엔드/UI 엔지니어 | 사용자 중심, 시각적, 접근성 |
| `devops` | DevOps/SRE | 운영적, 메트릭 중심, 자동화 |
| `data` | 데이터 엔지니어 | 분석적, 데이터 중심, 정밀 |

## 10. GitHub Actions 자동화

`.github/workflows/monthly-resume.yml`을 통해:
- 매월 마지막 날 자동 실행 (cron)
- 수동 트리거 지원 (workflow_dispatch)
- 템플릿 선택 가능
- Slack 알림 연동
- 결과물 아티팩트로 저장 (90일 보관)
