# Git Resume Analyzer API Specification

Base URL: `http://localhost:8080`

## Overview

Git Resume Analyzer API는 Git 커밋 히스토리를 분석하여 STAR 형식의 이력서 bullet point를 생성합니다.

---

## Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | 헬스 체크 |
| GET | `/ready` | 준비 상태 체크 |
| POST | `/api/analyze` | 분석 작업 시작 |
| GET | `/api/jobs` | 작업 목록 조회 |
| GET | `/api/jobs/:id` | 작업 상태 조회 |
| DELETE | `/api/jobs/:id` | 작업 취소 |
| GET | `/api/results` | 결과 목록 조회 |
| GET | `/api/results/:id` | 단일 결과 조회 |
| GET | `/api/export` | 결과 내보내기 |
| GET | `/api/templates` | 템플릿 목록 |
| GET | `/api/stats` | 통계 조회 |

---

## Health Check

### GET /health

서버 상태 확인

**Response**
```json
{
  "status": "ok",
  "timestamp": "2024-04-25T08:00:00Z"
}
```

### GET /ready

데이터베이스 연결 상태 확인

**Response**
```json
{
  "status": "ready",
  "timestamp": "2024-04-25T08:00:00Z"
}
```

---

## Analysis

### POST /api/analyze

비동기 분석 작업을 시작합니다.

**Request Body**
```json
{
  "repos": ["/path/to/repo1", "/path/to/repo2"],
  "from_date": "2024-01-01",
  "to_date": "2024-03-31",
  "month": 4,
  "year": 2024,
  "template": "default",
  "batch_size": 5,
  "dry_run": false
}
```

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `repos` | string[] | ✅ | - | 분석할 Git 저장소 경로 목록 |
| `from_date` | string | - | - | 시작 날짜 (YYYY-MM-DD) |
| `to_date` | string | - | - | 종료 날짜 (YYYY-MM-DD) |
| `month` | int | - | 현재 월 | 분석할 월 (1-12) |
| `year` | int | - | 현재 연도 | 분석할 연도 |
| `template` | string | - | "default" | 사용할 템플릿 |
| `batch_size` | int | - | 5 | 배치당 커밋 수 (1-20) |
| `dry_run` | bool | - | false | API 호출 없이 미리보기 |

> **Note**: `from_date`/`to_date` 또는 `month`/`year` 중 하나를 지정. 둘 다 없으면 현재 월 기준.

**Response** `202 Accepted`
```json
{
  "job_id": "550e8400-e29b-41d4-a716-446655440000",
  "message": "Analysis job started"
}
```

**Errors**
| Status | Code | Description |
|--------|------|-------------|
| 400 | `VALIDATION_ERROR` | 입력값 검증 실패 |
| 400 | `BAD_REQUEST` | 잘못된 요청 |

---

## Jobs

### GET /api/jobs

모든 작업 목록을 조회합니다.

**Response**
```json
{
  "jobs": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "status": "completed",
      "progress": 100,
      "phase": "complete",
      "message": "Generated 15 bullet points",
      "created_at": "2024-04-25T08:00:00Z",
      "started_at": "2024-04-25T08:00:01Z",
      "completed_at": "2024-04-25T08:00:30Z",
      "result_count": 15
    }
  ],
  "total": 1
}
```

### GET /api/jobs/:id

특정 작업의 상태를 조회합니다.

**Path Parameters**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string | Job UUID |

**Response**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "running",
  "progress": 45,
  "phase": "analyzing",
  "message": "Processing batch 3/7",
  "created_at": "2024-04-25T08:00:00Z",
  "started_at": "2024-04-25T08:00:01Z",
  "result_count": 6
}
```

**Job Status Values**
| Status | Description |
|--------|-------------|
| `pending` | 대기 중 |
| `running` | 실행 중 |
| `completed` | 완료 |
| `failed` | 실패 |
| `cancelled` | 취소됨 |

**Job Phase Values**
| Phase | Description |
|-------|-------------|
| `scanning` | Git 저장소 스캔 중 |
| `filtering` | 처리된 커밋 필터링 중 |
| `analyzing` | LLM 분석 중 |
| `saving` | 결과 저장 중 |
| `complete` | 완료 |
| `error` | 오류 발생 |

**Errors**
| Status | Code | Description |
|--------|------|-------------|
| 404 | `JOB_NOT_FOUND` | 작업을 찾을 수 없음 |

### DELETE /api/jobs/:id

실행 중인 작업을 취소합니다.

**Response**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "cancelled",
  "message": "Job cancelled"
}
```

---

## Results

### GET /api/results

분석 결과 목록을 조회합니다. 페이지네이션을 지원합니다.

**Query Parameters**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | int | 1 | 페이지 번호 |
| `page_size` | int | 20 | 페이지 크기 (max: 100) |
| `project` | string | - | 프로젝트 필터 |
| `category` | string | - | 카테고리 필터 |
| `from` | string | - | 시작 날짜 (YYYY-MM-DD) |
| `to` | string | - | 종료 날짜 (YYYY-MM-DD) |

**Response**
```json
{
  "results": [
    {
      "id": 1,
      "commit_hash": "abc1234def5678",
      "date": "2024-04-15T10:30:00Z",
      "project": "my-project",
      "category": "Feature",
      "impact_summary": "Engineered real-time notification system serving 10K concurrent users, reducing message latency by 40%",
      "created_at": "2024-04-25T08:00:30Z"
    }
  ],
  "total": 94,
  "page": 1,
  "page_size": 20,
  "total_pages": 5
}
```

**Category Values**
| Category | Description |
|----------|-------------|
| `Feature` | 새 기능 |
| `Fix` | 버그 수정 |
| `Refactor` | 리팩토링 |
| `Test` | 테스트 |
| `Docs` | 문서화 |
| `Chore` | 유지보수 |

### GET /api/results/:id

단일 결과를 조회합니다.

**Path Parameters**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | int | Result ID |

**Response**
```json
{
  "id": 1,
  "commit_hash": "abc1234def5678",
  "date": "2024-04-15T10:30:00Z",
  "project": "my-project",
  "category": "Feature",
  "impact_summary": "Engineered real-time notification system serving 10K concurrent users",
  "created_at": "2024-04-25T08:00:30Z"
}
```

**Errors**
| Status | Code | Description |
|--------|------|-------------|
| 404 | `NOT_FOUND` | 결과를 찾을 수 없음 |

---

## Export

### GET /api/export

결과를 파일 형식으로 내보냅니다.

**Query Parameters**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `format` | string | "json" | 출력 형식: `json`, `csv`, `markdown`, `md` |
| `project` | string | - | 프로젝트 필터 |
| `from` | string | - | 시작 날짜 (YYYY-MM-DD) |
| `to` | string | - | 종료 날짜 (YYYY-MM-DD) |

**Response Headers**
```
Content-Type: application/json | text/csv | text/markdown
Content-Disposition: attachment; filename=resume-export.{format}
```

**JSON Response**
```json
{
  "metadata": {
    "generated_at": "2024-04-25T08:00:00Z",
    "from_date": "2024-01-01",
    "to_date": "2024-03-31",
    "total_count": 94,
    "format": "json"
  },
  "achievements": [
    {
      "id": 1,
      "commit_hash": "abc1234",
      "date": "2024-04-15T10:30:00Z",
      "project": "my-project",
      "category": "Feature",
      "impact_summary": "..."
    }
  ],
  "summary": {
    "by_category": {
      "Feature": 60,
      "Fix": 21,
      "Refactor": 10,
      "Chore": 2,
      "Docs": 1
    },
    "by_project": {
      "my-project": 50,
      "another-project": 44
    }
  }
}
```

**CSV Response**
```csv
Date,Project,Category,Impact Summary,Commit Hash
2024-04-15,my-project,Feature,"Engineered real-time...",abc1234
```

**Markdown Response**
```markdown
# Resume Achievements

*Generated on April 25, 2024*

## ✨ Feature

- Engineered real-time notification system (my-project, 2024-04-15) `abc1234`

## 🐛 Fix

- Resolved critical memory leak... (my-project, 2024-04-10) `def5678`
```

---

## Templates

### GET /api/templates

사용 가능한 프롬프트 템플릿 목록을 조회합니다.

**Response**
```json
{
  "templates": [
    {
      "name": "Default",
      "description": "Balanced template suitable for most tech roles",
      "tone_style": "professional, concise, achievement-focused",
      "focus": ["technical impact", "quantifiable results", "problem-solving"]
    },
    {
      "name": "Startup",
      "description": "Fast-paced startup environment emphasis",
      "tone_style": "dynamic, results-driven, entrepreneurial",
      "focus": ["rapid delivery", "cross-functional impact", "innovation", "scalability"]
    },
    {
      "name": "Enterprise",
      "description": "Large corporation and enterprise focus",
      "tone_style": "formal, process-oriented, compliance-aware",
      "focus": ["reliability", "security", "compliance", "stakeholder management"]
    },
    {
      "name": "Backend Engineer",
      "description": "Backend/infrastructure engineering focus",
      "tone_style": "technical, precise, systems-focused",
      "focus": ["performance", "scalability", "reliability", "data integrity"]
    },
    {
      "name": "Frontend Engineer",
      "description": "Frontend/UI engineering focus",
      "tone_style": "user-centric, visual, accessibility-aware",
      "focus": ["user experience", "performance", "accessibility", "design systems"]
    },
    {
      "name": "DevOps/SRE",
      "description": "DevOps and Site Reliability Engineering focus",
      "tone_style": "operational, metrics-driven, automation-focused",
      "focus": ["reliability", "automation", "observability", "incident response"]
    },
    {
      "name": "Data Engineer",
      "description": "Data engineering and analytics focus",
      "tone_style": "analytical, data-driven, precision-focused",
      "focus": ["data quality", "pipeline efficiency", "analytics enablement", "data governance"]
    }
  ]
}
```

---

## Statistics

### GET /api/stats

대시보드용 통계 데이터를 조회합니다.

**Response**
```json
{
  "total_results": 94,
  "total_commits": 0,
  "tokens_used": {
    "input_tokens": 6244,
    "output_tokens": 6244,
    "total_cost": 0.15
  },
  "category_breakdown": {
    "Feature": 60,
    "Fix": 21,
    "Refactor": 10,
    "Chore": 2,
    "Docs": 1
  },
  "project_breakdown": {
    "my-project": 50,
    "another-project": 44
  },
  "recent_activity": [
    {"date": "2024-04-25", "count": 15},
    {"date": "2024-04-24", "count": 8}
  ]
}
```

---

## Error Responses

모든 에러는 다음 형식으로 반환됩니다:

```json
{
  "code": "ERROR_CODE",
  "message": "Human readable error message",
  "details": {
    "field": "error description"
  },
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Error Codes**
| Code | HTTP Status | Description |
|------|-------------|-------------|
| `BAD_REQUEST` | 400 | 잘못된 요청 |
| `VALIDATION_ERROR` | 400 | 입력값 검증 실패 |
| `NOT_FOUND` | 404 | 리소스 없음 |
| `JOB_NOT_FOUND` | 404 | 작업 없음 |
| `REPOSITORY_ERROR` | 400 | Git 저장소 오류 |
| `ANALYSIS_ERROR` | 500 | 분석 처리 오류 |
| `EXPORT_ERROR` | 500 | 내보내기 오류 |
| `INTERNAL_ERROR` | 500 | 서버 내부 오류 |

---

## Examples

### cURL Examples

```bash
# 헬스 체크
curl http://localhost:8080/health

# 분석 시작
curl -X POST http://localhost:8080/api/analyze \
  -H "Content-Type: application/json" \
  -d '{"repos": ["/path/to/repo"], "month": 4, "year": 2024}'

# 작업 상태 확인
curl http://localhost:8080/api/jobs/550e8400-e29b-41d4-a716-446655440000

# 결과 조회 (페이지네이션)
curl "http://localhost:8080/api/results?page=1&page_size=10"

# 결과 필터링
curl "http://localhost:8080/api/results?category=Feature&project=my-project"

# CSV 내보내기
curl "http://localhost:8080/api/export?format=csv" -o results.csv

# 통계 조회
curl http://localhost:8080/api/stats

# 템플릿 목록
curl http://localhost:8080/api/templates
```

---

## Rate Limiting

현재 Rate limiting은 구현되어 있지 않습니다. 프로덕션 환경에서는 reverse proxy (nginx, etc.)를 통해 구현하는 것을 권장합니다.

---

## CORS

기본적으로 모든 origin에서 접근 가능합니다 (`Access-Control-Allow-Origin: *`).

서버 시작 시 `--cors-origins` 플래그로 제한할 수 있습니다:

```bash
./bin/git-resume serve --cors-origins=http://localhost:3000,https://myapp.com
```
