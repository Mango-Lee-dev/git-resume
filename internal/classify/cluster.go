// Package classify provides commit classification functionality for grouping
// similar work patterns across multiple projects.
package classify

import (
	"regexp"
	"strings"
)

// WorkCluster represents a category of work that can span multiple projects.
type WorkCluster string

const (
	ClusterAuth        WorkCluster = "auth"
	ClusterAPI         WorkCluster = "api"
	ClusterDatabase    WorkCluster = "database"
	ClusterUI          WorkCluster = "ui"
	ClusterTesting     WorkCluster = "testing"
	ClusterPerformance WorkCluster = "performance"
	ClusterSecurity    WorkCluster = "security"
	ClusterInfra       WorkCluster = "infra"
	ClusterRefactor    WorkCluster = "refactor"
	ClusterDocs        WorkCluster = "docs"
	ClusterOther       WorkCluster = "other"
)

// AllClusters returns all defined work clusters (excluding Other).
func AllClusters() []WorkCluster {
	return []WorkCluster{
		ClusterAuth,
		ClusterAPI,
		ClusterDatabase,
		ClusterUI,
		ClusterTesting,
		ClusterPerformance,
		ClusterSecurity,
		ClusterInfra,
		ClusterRefactor,
		ClusterDocs,
	}
}

// ClusterClassifier classifies commits into work clusters based on
// commit messages and file paths.
type ClusterClassifier struct {
	patterns map[WorkCluster][]*regexp.Regexp
}

// NewClusterClassifier creates a new classifier with predefined patterns.
func NewClusterClassifier() *ClusterClassifier {
	c := &ClusterClassifier{
		patterns: make(map[WorkCluster][]*regexp.Regexp),
	}
	c.initPatterns()
	return c
}

func (c *ClusterClassifier) initPatterns() {
	// Auth patterns
	// Note: \b is ASCII-only in Go's RE2 and does not match around Hangul,
	// so Korean keywords are matched as plain substrings.
	c.patterns[ClusterAuth] = compilePatterns(
		`(?i)\bauth\b`,
		`(?i)\blogin\b`,
		`(?i)\blogout\b`,
		`(?i)\bjwt\b`,
		`(?i)\boauth\b`,
		`(?i)\bsession\b`,
		`(?i)\btoken\b`,
		`(?i)\bpassword\b`,
		`(?i)\bcredential`,
		`(?i)\bsign[_-]?(in|up|out)\b`,
		`(?i)\bpermission`,
		`(?i)\baccess[_-]?control`,
		`(?i)\brole\b`,
		`(?i)\brbac\b`,
		// Korean
		`인증`,
		`로그인`,
		`로그아웃`,
		`회원가입`,
		`권한`,
		`세션`,
		`토큰`,
		`비밀번호`,
		`패스워드`,
		`자격증명`,
		`역할`,
	)

	// API patterns
	c.patterns[ClusterAPI] = compilePatterns(
		`(?i)\bapi\b`,
		`(?i)\bendpoint`,
		`(?i)\brest\b`,
		`(?i)\bgraphql\b`,
		`(?i)\bgrpc\b`,
		`(?i)\broute\b`,
		`(?i)\bhandler\b`,
		`(?i)\bmiddleware\b`,
		`(?i)\brequest\b`,
		`(?i)\bresponse\b`,
		`(?i)\bhttp\b`,
		`(?i)\bwebhook\b`,
		`(?i)/api/`,
		`(?i)/handlers/`,
		`(?i)/routes/`,
		// Korean
		`엔드포인트`,
		`라우트`,
		`라우터`,
		`핸들러`,
		`미들웨어`,
		`웹훅`,
	)

	// Database patterns
	c.patterns[ClusterDatabase] = compilePatterns(
		`(?i)\bdatabase\b`,
		`(?i)\bdb\b`,
		`(?i)\bsql\b`,
		`(?i)\bquery\b`,
		`(?i)\bmigration`,
		`(?i)\bschema\b`,
		`(?i)\bmodel\b`,
		`(?i)\btable\b`,
		`(?i)\borm\b`,
		`(?i)\bpostgres`,
		`(?i)\bmysql\b`,
		`(?i)\bmongo`,
		`(?i)\bredis\b`,
		`(?i)\bsqlite\b`,
		`(?i)\.sql$`,
		`(?i)/migrations?/`,
		`(?i)/models?/`,
		// Korean
		`데이터베이스`,
		`디비`,
		`쿼리`,
		`마이그레이션`,
		`스키마`,
		`테이블`,
	)

	// UI patterns
	c.patterns[ClusterUI] = compilePatterns(
		`(?i)\bui\b`,
		`(?i)\bcomponent`,
		`(?i)\bstyle`,
		`(?i)\bcss\b`,
		`(?i)\bscss\b`,
		`(?i)\blayout\b`,
		`(?i)\bresponsive\b`,
		`(?i)\bdesign\b`,
		`(?i)\bfrontend\b`,
		`(?i)\bfront[_-]?end\b`,
		`(?i)\breact\b`,
		`(?i)\bvue\b`,
		`(?i)\bsvelte\b`,
		`(?i)\bbutton\b`,
		`(?i)\bform\b`,
		`(?i)\bmodal\b`,
		`(?i)\bpage\b`,
		`(?i)\.tsx$`,
		`(?i)\.jsx$`,
		`(?i)\.vue$`,
		`(?i)\.css$`,
		`(?i)\.scss$`,
		`(?i)/components?/`,
		`(?i)/pages?/`,
		`(?i)/views?/`,
		`(?i)/styles?/`,
		// Korean
		`화면`,
		`페이지`,
		`컴포넌트`,
		`레이아웃`,
		`디자인`,
		`스타일`,
		`버튼`,
		`모달`,
		`프론트엔드`,
		`반응형`,
	)

	// Testing patterns
	c.patterns[ClusterTesting] = compilePatterns(
		`(?i)\btest\b`,
		`(?i)\bspec\b`,
		`(?i)\bcoverage\b`,
		`(?i)\bmock\b`,
		`(?i)\bfixture`,
		`(?i)\be2e\b`,
		`(?i)\bunit\b`,
		`(?i)\bintegration[_-]?test`,
		`(?i)\bassert`,
		`(?i)\bexpect\b`,
		`(?i)_test\.go$`,
		`(?i)\.test\.(ts|js|tsx|jsx)$`,
		`(?i)\.spec\.(ts|js|tsx|jsx)$`,
		`(?i)/__tests__/`,
		`(?i)/test/`,
		`(?i)/tests/`,
		// Korean
		`테스트`,
		`단위테스트`,
		`통합테스트`,
		`커버리지`,
		`모킹`,
		`픽스처`,
	)

	// Performance patterns
	c.patterns[ClusterPerformance] = compilePatterns(
		`(?i)\bperformance\b`,
		`(?i)\boptimize`,
		`(?i)\bcache\b`,
		`(?i)\bcaching\b`,
		`(?i)\bspeed\b`,
		`(?i)\blatency\b`,
		`(?i)\bmemory\b`,
		`(?i)\bbenchmark`,
		`(?i)\bprofile`,
		`(?i)\bindex\b`,
		`(?i)\bparallel`,
		`(?i)\bconcurrent`,
		`(?i)\basync\b`,
		`(?i)\bbuffering\b`,
		// Korean
		`성능`,
		`최적화`,
		`캐시`,
		`캐싱`,
		`지연시간`,
		`메모리`,
		`벤치마크`,
		`프로파일링`,
		`병렬`,
		`동시성`,
		`비동기`,
	)

	// Security patterns
	c.patterns[ClusterSecurity] = compilePatterns(
		`(?i)\bsecurity\b`,
		`(?i)\bvulnerability`,
		`(?i)\bxss\b`,
		`(?i)\bcsrf\b`,
		`(?i)\bsanitize`,
		`(?i)\bencrypt`,
		`(?i)\bdecrypt`,
		`(?i)\bhash\b`,
		`(?i)\bssl\b`,
		`(?i)\btls\b`,
		`(?i)\bcert\b`,
		`(?i)\binjection\b`,
		`(?i)\baudit\b`,
		`(?i)\bcve\b`,
		// Korean
		`보안`,
		`취약점`,
		`암호화`,
		`복호화`,
		`인증서`,
		`인젝션`,
		`감사`,
		`해시`,
	)

	// Infra patterns
	c.patterns[ClusterInfra] = compilePatterns(
		`(?i)\bdeploy`,
		`(?i)\bdocker`,
		`(?i)\bkubernetes\b`,
		`(?i)\bk8s\b`,
		`(?i)\bci\b`,
		`(?i)\bcd\b`,
		`(?i)\bpipeline\b`,
		`(?i)\bterraform\b`,
		`(?i)\bansible\b`,
		`(?i)\bhelm\b`,
		`(?i)\baws\b`,
		`(?i)\bgcp\b`,
		`(?i)\bazure\b`,
		`(?i)\bcloud\b`,
		`(?i)\binfra`,
		`(?i)\bconfig`,
		`(?i)\benv\b`,
		`(?i)Dockerfile`,
		`(?i)docker-compose`,
		`(?i)\.tf$`,
		`(?i)\.yaml$`,
		`(?i)\.yml$`,
		`(?i)/\.github/`,
		`(?i)/kubernetes/`,
		`(?i)/k8s/`,
		`(?i)/deploy/`,
		`(?i)/infra/`,
		// Korean
		`배포`,
		`도커`,
		`쿠버네티스`,
		`파이프라인`,
		`클라우드`,
		`인프라`,
		`환경변수`,
		`환경설정`,
	)

	// Refactor patterns
	c.patterns[ClusterRefactor] = compilePatterns(
		`(?i)\brefactor`,
		`(?i)\bclean\b`,
		`(?i)\brestructure`,
		`(?i)\breorganize`,
		`(?i)\bsimplify`,
		`(?i)\bextract\b`,
		`(?i)\brename\b`,
		`(?i)\bmove\b`,
		`(?i)\bsplit\b`,
		`(?i)\bmerge\b`,
		`(?i)\bconsolidate`,
		`(?i)\bmodularize`,
		`(?i)\bdecouple`,
		// Korean
		`리팩토링`,
		`리팩터링`,
		`재구성`,
		`단순화`,
		`이름변경`,
		`모듈화`,
		`코드정리`,
	)

	// Docs patterns
	c.patterns[ClusterDocs] = compilePatterns(
		`(?i)\bdocument`,
		`(?i)\breadme\b`,
		`(?i)\bguide\b`,
		`(?i)\bcomment`,
		`(?i)\bjsdoc\b`,
		`(?i)\bgodoc\b`,
		`(?i)\bchangelog\b`,
		`(?i)\bapi[_-]?doc`,
		`(?i)\bswagger\b`,
		`(?i)\bopenapi\b`,
		`(?i)\.md$`,
		`(?i)\.rst$`,
		`(?i)\.adoc$`,
		`(?i)/docs?/`,
		`(?i)/documentation/`,
		// Korean
		`문서`,
		`가이드`,
		`주석`,
		`변경로그`,
		`변경기록`,
		`명세`,
	)
}

func compilePatterns(patterns ...string) []*regexp.Regexp {
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		if re, err := regexp.Compile(p); err == nil {
			compiled = append(compiled, re)
		}
	}
	return compiled
}

// Classify analyzes a commit message and file paths to determine
// which work clusters the commit belongs to.
// Returns a slice of clusters (may contain multiple or empty if no match).
func (c *ClusterClassifier) Classify(message string, files []string) []WorkCluster {
	clusters := make(map[WorkCluster]bool)

	// Combine message and file paths for analysis
	textToAnalyze := strings.ToLower(message)
	for _, f := range files {
		textToAnalyze += " " + strings.ToLower(f)
	}

	// Check each cluster's patterns
	for cluster, patterns := range c.patterns {
		for _, re := range patterns {
			if re.MatchString(textToAnalyze) {
				clusters[cluster] = true
				break
			}
		}
	}

	// Convert map to slice
	result := make([]WorkCluster, 0, len(clusters))
	for cluster := range clusters {
		result = append(result, cluster)
	}

	// If no clusters matched, assign to "other"
	if len(result) == 0 {
		result = append(result, ClusterOther)
	}

	return result
}

// ClassifyToStrings is a convenience method that returns cluster names as strings.
func (c *ClusterClassifier) ClassifyToStrings(message string, files []string) []string {
	clusters := c.Classify(message, files)
	result := make([]string, len(clusters))
	for i, cluster := range clusters {
		result[i] = string(cluster)
	}
	return result
}

// ClusterDisplayName returns a human-readable name for a cluster.
func ClusterDisplayName(cluster WorkCluster) string {
	names := map[WorkCluster]string{
		ClusterAuth:        "Authentication",
		ClusterAPI:         "API Development",
		ClusterDatabase:    "Database",
		ClusterUI:          "UI/Frontend",
		ClusterTesting:     "Testing",
		ClusterPerformance: "Performance",
		ClusterSecurity:    "Security",
		ClusterInfra:       "Infrastructure",
		ClusterRefactor:    "Refactoring",
		ClusterDocs:        "Documentation",
		ClusterOther:       "Other",
	}
	if name, ok := names[cluster]; ok {
		return name
	}
	return string(cluster)
}
