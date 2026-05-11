package classify

import (
	"testing"
)

func TestNewClusterClassifier(t *testing.T) {
	c := NewClusterClassifier()
	if c == nil {
		t.Fatal("expected non-nil classifier")
	}
	if len(c.patterns) == 0 {
		t.Fatal("expected patterns to be initialized")
	}
}

func TestClassify_Auth(t *testing.T) {
	c := NewClusterClassifier()

	tests := []struct {
		name    string
		message string
		files   []string
		want    WorkCluster
	}{
		{
			name:    "jwt implementation",
			message: "implement JWT authentication",
			files:   []string{"internal/auth/jwt.go"},
			want:    ClusterAuth,
		},
		{
			name:    "login feature",
			message: "add login page",
			files:   []string{"src/pages/login.tsx"},
			want:    ClusterAuth,
		},
		{
			name:    "oauth integration",
			message: "integrate OAuth2 provider",
			files:   []string{"pkg/oauth/google.go"},
			want:    ClusterAuth,
		},
		{
			name:    "session management",
			message: "fix session expiry bug",
			files:   []string{"internal/session/store.go"},
			want:    ClusterAuth,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clusters := c.Classify(tt.message, tt.files)
			if !containsCluster(clusters, tt.want) {
				t.Errorf("expected %v to be in clusters %v", tt.want, clusters)
			}
		})
	}
}

func TestClassify_API(t *testing.T) {
	c := NewClusterClassifier()

	tests := []struct {
		name    string
		message string
		files   []string
		want    WorkCluster
	}{
		{
			name:    "rest endpoint",
			message: "add REST endpoint for users",
			files:   []string{"internal/api/users.go"},
			want:    ClusterAPI,
		},
		{
			name:    "graphql schema",
			message: "update GraphQL schema",
			files:   []string{"schema.graphql"},
			want:    ClusterAPI,
		},
		{
			name:    "handler implementation",
			message: "implement user handler",
			files:   []string{"internal/handlers/user_handler.go"},
			want:    ClusterAPI,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clusters := c.Classify(tt.message, tt.files)
			if !containsCluster(clusters, tt.want) {
				t.Errorf("expected %v to be in clusters %v", tt.want, clusters)
			}
		})
	}
}

func TestClassify_Database(t *testing.T) {
	c := NewClusterClassifier()

	tests := []struct {
		name    string
		message string
		files   []string
		want    WorkCluster
	}{
		{
			name:    "sql migration",
			message: "add users table migration",
			files:   []string{"migrations/001_create_users.sql"},
			want:    ClusterDatabase,
		},
		{
			name:    "query optimization",
			message: "optimize database query",
			files:   []string{"internal/db/queries.go"},
			want:    ClusterDatabase,
		},
		{
			name:    "model update",
			message: "update user model",
			files:   []string{"internal/models/user.go"},
			want:    ClusterDatabase,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clusters := c.Classify(tt.message, tt.files)
			if !containsCluster(clusters, tt.want) {
				t.Errorf("expected %v to be in clusters %v", tt.want, clusters)
			}
		})
	}
}

func TestClassify_UI(t *testing.T) {
	c := NewClusterClassifier()

	tests := []struct {
		name    string
		message string
		files   []string
		want    WorkCluster
	}{
		{
			name:    "react component",
			message: "add button component",
			files:   []string{"src/components/Button.tsx"},
			want:    ClusterUI,
		},
		{
			name:    "css styles",
			message: "update styles",
			files:   []string{"src/styles/main.css"},
			want:    ClusterUI,
		},
		{
			name:    "page layout",
			message: "responsive layout fix",
			files:   []string{"src/pages/Home.tsx"},
			want:    ClusterUI,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clusters := c.Classify(tt.message, tt.files)
			if !containsCluster(clusters, tt.want) {
				t.Errorf("expected %v to be in clusters %v", tt.want, clusters)
			}
		})
	}
}

func TestClassify_Testing(t *testing.T) {
	c := NewClusterClassifier()

	tests := []struct {
		name    string
		message string
		files   []string
		want    WorkCluster
	}{
		{
			name:    "go test file",
			message: "add unit tests",
			files:   []string{"internal/service/user_test.go"},
			want:    ClusterTesting,
		},
		{
			name:    "jest test",
			message: "add component tests",
			files:   []string{"src/components/Button.test.tsx"},
			want:    ClusterTesting,
		},
		{
			name:    "e2e tests",
			message: "add e2e tests for checkout",
			files:   []string{"tests/e2e/checkout.spec.ts"},
			want:    ClusterTesting,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clusters := c.Classify(tt.message, tt.files)
			if !containsCluster(clusters, tt.want) {
				t.Errorf("expected %v to be in clusters %v", tt.want, clusters)
			}
		})
	}
}

func TestClassify_Infra(t *testing.T) {
	c := NewClusterClassifier()

	tests := []struct {
		name    string
		message string
		files   []string
		want    WorkCluster
	}{
		{
			name:    "dockerfile",
			message: "update Docker image",
			files:   []string{"Dockerfile"},
			want:    ClusterInfra,
		},
		{
			name:    "kubernetes manifest",
			message: "update k8s deployment",
			files:   []string{"kubernetes/deployment.yaml"},
			want:    ClusterInfra,
		},
		{
			name:    "terraform config",
			message: "add terraform module",
			files:   []string{"infra/main.tf"},
			want:    ClusterInfra,
		},
		{
			name:    "github actions",
			message: "update CI pipeline",
			files:   []string{".github/workflows/ci.yml"},
			want:    ClusterInfra,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clusters := c.Classify(tt.message, tt.files)
			if !containsCluster(clusters, tt.want) {
				t.Errorf("expected %v to be in clusters %v", tt.want, clusters)
			}
		})
	}
}

func TestClassify_MultipleClusters(t *testing.T) {
	c := NewClusterClassifier()

	// A commit that touches both API and Database
	message := "add user API endpoint with database migration"
	files := []string{
		"internal/api/users.go",
		"migrations/001_users.sql",
	}

	clusters := c.Classify(message, files)

	if !containsCluster(clusters, ClusterAPI) {
		t.Error("expected API cluster")
	}
	if !containsCluster(clusters, ClusterDatabase) {
		t.Error("expected Database cluster")
	}
}

func TestClassify_NoMatch(t *testing.T) {
	c := NewClusterClassifier()

	message := "misc changes"
	files := []string{"random.txt"}

	clusters := c.Classify(message, files)

	if len(clusters) != 1 || clusters[0] != ClusterOther {
		t.Errorf("expected [other], got %v", clusters)
	}
}

func TestClassifyToStrings(t *testing.T) {
	c := NewClusterClassifier()

	message := "implement JWT auth"
	files := []string{"internal/auth/jwt.go"}

	strings := c.ClassifyToStrings(message, files)

	found := false
	for _, s := range strings {
		if s == "auth" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'auth' in strings %v", strings)
	}
}

func TestClusterDisplayName(t *testing.T) {
	tests := []struct {
		cluster WorkCluster
		want    string
	}{
		{ClusterAuth, "Authentication"},
		{ClusterAPI, "API Development"},
		{ClusterDatabase, "Database"},
		{ClusterUI, "UI/Frontend"},
		{ClusterTesting, "Testing"},
		{ClusterPerformance, "Performance"},
		{ClusterSecurity, "Security"},
		{ClusterInfra, "Infrastructure"},
		{ClusterRefactor, "Refactoring"},
		{ClusterDocs, "Documentation"},
		{ClusterOther, "Other"},
	}

	for _, tt := range tests {
		t.Run(string(tt.cluster), func(t *testing.T) {
			got := ClusterDisplayName(tt.cluster)
			if got != tt.want {
				t.Errorf("ClusterDisplayName(%v) = %v, want %v", tt.cluster, got, tt.want)
			}
		})
	}
}

func TestAllClusters(t *testing.T) {
	clusters := AllClusters()
	if len(clusters) != 10 {
		t.Errorf("expected 10 clusters, got %d", len(clusters))
	}

	// Verify Other is not included
	for _, c := range clusters {
		if c == ClusterOther {
			t.Error("AllClusters should not include ClusterOther")
		}
	}
}

func containsCluster(clusters []WorkCluster, target WorkCluster) bool {
	for _, c := range clusters {
		if c == target {
			return true
		}
	}
	return false
}
