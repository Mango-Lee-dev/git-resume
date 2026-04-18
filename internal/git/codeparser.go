package git

import (
	"fmt"
	"regexp"
	"strings"
)

// CodeParser extracts meaningful code changes from diffs
type CodeParser struct {
	parsers map[string]LanguageParser
}

// LanguageParser defines how to parse a specific language
type LanguageParser interface {
	// ExtractFunctions extracts function/method definitions from code
	ExtractFunctions(content string) []FunctionInfo
	// ExtractChanges analyzes diff to find meaningful changes
	ExtractChanges(diff string) []CodeChange
}

// FunctionInfo represents a function or method
type FunctionInfo struct {
	Name       string
	Signature  string
	StartLine  int
	EndLine    int
	Visibility string // public, private, etc.
}

// CodeChange represents a meaningful code change
type CodeChange struct {
	Type        ChangeType
	Name        string
	Description string
	Impact      string
}

// ChangeType categorizes the type of code change
type ChangeType string

const (
	ChangeTypeNewFunction    ChangeType = "new_function"
	ChangeTypeModifyFunction ChangeType = "modify_function"
	ChangeTypeDeleteFunction ChangeType = "delete_function"
	ChangeTypeNewClass       ChangeType = "new_class"
	ChangeTypeNewInterface   ChangeType = "new_interface"
	ChangeTypeNewEndpoint    ChangeType = "new_endpoint"
	ChangeTypeNewTest        ChangeType = "new_test"
	ChangeTypeRefactor       ChangeType = "refactor"
)

// NewCodeParser creates a new code parser with all supported languages
func NewCodeParser() *CodeParser {
	return &CodeParser{
		parsers: map[string]LanguageParser{
			".go":   &GoParser{},
			".py":   &PythonParser{},
			".js":   &JavaScriptParser{},
			".ts":   &TypeScriptParser{},
			".tsx":  &TypeScriptParser{},
			".java": &JavaParser{},
			".rs":   &RustParser{},
		},
	}
}

// GetParser returns the appropriate parser for a file extension
func (cp *CodeParser) GetParser(ext string) LanguageParser {
	if parser, ok := cp.parsers[strings.ToLower(ext)]; ok {
		return parser
	}
	return nil
}

// AnalyzeDiff analyzes a diff and returns meaningful changes
func (cp *CodeParser) AnalyzeDiff(filename string, diff string) []CodeChange {
	ext := getFileExtension(filename)
	parser := cp.GetParser(ext)
	if parser == nil {
		return nil
	}
	return parser.ExtractChanges(diff)
}

// GoParser parses Go code
type GoParser struct{}

func (p *GoParser) ExtractFunctions(content string) []FunctionInfo {
	var functions []FunctionInfo

	// Match function declarations
	funcPattern := regexp.MustCompile(`(?m)^func\s+(?:\([\w\s*]+\)\s+)?(\w+)\s*\([^)]*\)`)
	matches := funcPattern.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) >= 2 {
			functions = append(functions, FunctionInfo{
				Name:       match[1],
				Signature:  match[0],
				Visibility: getGoVisibility(match[1]),
			})
		}
	}

	return functions
}

func (p *GoParser) ExtractChanges(diff string) []CodeChange {
	var changes []CodeChange

	// Look for new functions
	newFuncPattern := regexp.MustCompile(`(?m)^\+func\s+(?:\([\w\s*]+\)\s+)?(\w+)\s*\(([^)]*)\)`)
	newFuncs := newFuncPattern.FindAllStringSubmatch(diff, -1)
	for _, match := range newFuncs {
		changes = append(changes, CodeChange{
			Type:        ChangeTypeNewFunction,
			Name:        match[1],
			Description: fmt.Sprintf("Added function %s", match[1]),
		})
	}

	// Look for new structs
	newStructPattern := regexp.MustCompile(`(?m)^\+type\s+(\w+)\s+struct`)
	newStructs := newStructPattern.FindAllStringSubmatch(diff, -1)
	for _, match := range newStructs {
		changes = append(changes, CodeChange{
			Type:        ChangeTypeNewClass,
			Name:        match[1],
			Description: fmt.Sprintf("Added struct %s", match[1]),
		})
	}

	// Look for new interfaces
	newInterfacePattern := regexp.MustCompile(`(?m)^\+type\s+(\w+)\s+interface`)
	newInterfaces := newInterfacePattern.FindAllStringSubmatch(diff, -1)
	for _, match := range newInterfaces {
		changes = append(changes, CodeChange{
			Type:        ChangeTypeNewInterface,
			Name:        match[1],
			Description: fmt.Sprintf("Added interface %s", match[1]),
		})
	}

	// Look for HTTP handlers (common Go patterns)
	handlerPattern := regexp.MustCompile(`(?m)^\+.*\.(Handle|HandleFunc|Get|Post|Put|Delete|Patch)\s*\(\s*["']([^"']+)["']`)
	handlers := handlerPattern.FindAllStringSubmatch(diff, -1)
	for _, match := range handlers {
		changes = append(changes, CodeChange{
			Type:        ChangeTypeNewEndpoint,
			Name:        match[2],
			Description: fmt.Sprintf("Added %s endpoint %s", match[1], match[2]),
		})
	}

	return changes
}

func getGoVisibility(name string) string {
	if len(name) > 0 && name[0] >= 'A' && name[0] <= 'Z' {
		return "public"
	}
	return "private"
}

// PythonParser parses Python code
type PythonParser struct{}

func (p *PythonParser) ExtractFunctions(content string) []FunctionInfo {
	var functions []FunctionInfo

	// Match function definitions
	funcPattern := regexp.MustCompile(`(?m)^(?:async\s+)?def\s+(\w+)\s*\(([^)]*)\)`)
	matches := funcPattern.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) >= 2 {
			visibility := "public"
			if strings.HasPrefix(match[1], "_") {
				visibility = "private"
			}
			functions = append(functions, FunctionInfo{
				Name:       match[1],
				Signature:  match[0],
				Visibility: visibility,
			})
		}
	}

	return functions
}

func (p *PythonParser) ExtractChanges(diff string) []CodeChange {
	var changes []CodeChange

	// New functions
	newFuncPattern := regexp.MustCompile(`(?m)^\+(?:async\s+)?def\s+(\w+)\s*\(`)
	newFuncs := newFuncPattern.FindAllStringSubmatch(diff, -1)
	for _, match := range newFuncs {
		changes = append(changes, CodeChange{
			Type:        ChangeTypeNewFunction,
			Name:        match[1],
			Description: fmt.Sprintf("Added function %s", match[1]),
		})
	}

	// New classes
	newClassPattern := regexp.MustCompile(`(?m)^\+class\s+(\w+)`)
	newClasses := newClassPattern.FindAllStringSubmatch(diff, -1)
	for _, match := range newClasses {
		changes = append(changes, CodeChange{
			Type:        ChangeTypeNewClass,
			Name:        match[1],
			Description: fmt.Sprintf("Added class %s", match[1]),
		})
	}

	// FastAPI/Flask routes
	routePattern := regexp.MustCompile(`(?m)^\+@(app|router)\.(get|post|put|delete|patch)\s*\(\s*["']([^"']+)["']`)
	routes := routePattern.FindAllStringSubmatch(diff, -1)
	for _, match := range routes {
		changes = append(changes, CodeChange{
			Type:        ChangeTypeNewEndpoint,
			Name:        match[3],
			Description: fmt.Sprintf("Added %s endpoint %s", strings.ToUpper(match[2]), match[3]),
		})
	}

	return changes
}

// JavaScriptParser parses JavaScript code
type JavaScriptParser struct{}

func (p *JavaScriptParser) ExtractFunctions(content string) []FunctionInfo {
	var functions []FunctionInfo

	// Match various function patterns
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?m)function\s+(\w+)\s*\(`),
		regexp.MustCompile(`(?m)const\s+(\w+)\s*=\s*(?:async\s*)?\(`),
		regexp.MustCompile(`(?m)(\w+)\s*:\s*(?:async\s*)?function\s*\(`),
	}

	for _, pattern := range patterns {
		matches := pattern.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			if len(match) >= 2 {
				functions = append(functions, FunctionInfo{
					Name:      match[1],
					Signature: match[0],
				})
			}
		}
	}

	return functions
}

func (p *JavaScriptParser) ExtractChanges(diff string) []CodeChange {
	var changes []CodeChange

	// New functions/arrow functions
	funcPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?m)^\+\s*(?:export\s+)?(?:async\s+)?function\s+(\w+)`),
		regexp.MustCompile(`(?m)^\+\s*(?:export\s+)?const\s+(\w+)\s*=\s*(?:async\s*)?\(`),
	}

	for _, pattern := range funcPatterns {
		matches := pattern.FindAllStringSubmatch(diff, -1)
		for _, match := range matches {
			changes = append(changes, CodeChange{
				Type:        ChangeTypeNewFunction,
				Name:        match[1],
				Description: fmt.Sprintf("Added function %s", match[1]),
			})
		}
	}

	// New classes
	classPattern := regexp.MustCompile(`(?m)^\+\s*(?:export\s+)?class\s+(\w+)`)
	classes := classPattern.FindAllStringSubmatch(diff, -1)
	for _, match := range classes {
		changes = append(changes, CodeChange{
			Type:        ChangeTypeNewClass,
			Name:        match[1],
			Description: fmt.Sprintf("Added class %s", match[1]),
		})
	}

	// Express routes
	routePattern := regexp.MustCompile(`(?m)^\+.*\.(get|post|put|delete|patch)\s*\(\s*["']([^"']+)["']`)
	routes := routePattern.FindAllStringSubmatch(diff, -1)
	for _, match := range routes {
		changes = append(changes, CodeChange{
			Type:        ChangeTypeNewEndpoint,
			Name:        match[2],
			Description: fmt.Sprintf("Added %s endpoint %s", strings.ToUpper(match[1]), match[2]),
		})
	}

	// React components
	componentPattern := regexp.MustCompile(`(?m)^\+\s*(?:export\s+)?(?:const|function)\s+([A-Z]\w+)`)
	components := componentPattern.FindAllStringSubmatch(diff, -1)
	for _, match := range components {
		changes = append(changes, CodeChange{
			Type:        ChangeTypeNewFunction,
			Name:        match[1],
			Description: fmt.Sprintf("Added React component %s", match[1]),
		})
	}

	return changes
}

// TypeScriptParser extends JavaScript parser with TS-specific patterns
type TypeScriptParser struct {
	JavaScriptParser
}

func (p *TypeScriptParser) ExtractChanges(diff string) []CodeChange {
	changes := p.JavaScriptParser.ExtractChanges(diff)

	// TypeScript interfaces
	interfacePattern := regexp.MustCompile(`(?m)^\+\s*(?:export\s+)?interface\s+(\w+)`)
	interfaces := interfacePattern.FindAllStringSubmatch(diff, -1)
	for _, match := range interfaces {
		changes = append(changes, CodeChange{
			Type:        ChangeTypeNewInterface,
			Name:        match[1],
			Description: fmt.Sprintf("Added interface %s", match[1]),
		})
	}

	// Type aliases
	typePattern := regexp.MustCompile(`(?m)^\+\s*(?:export\s+)?type\s+(\w+)\s*=`)
	types := typePattern.FindAllStringSubmatch(diff, -1)
	for _, match := range types {
		changes = append(changes, CodeChange{
			Type:        ChangeTypeNewInterface,
			Name:        match[1],
			Description: fmt.Sprintf("Added type %s", match[1]),
		})
	}

	return changes
}

// JavaParser parses Java code
type JavaParser struct{}

func (p *JavaParser) ExtractFunctions(content string) []FunctionInfo {
	var functions []FunctionInfo

	// Match method declarations
	methodPattern := regexp.MustCompile(`(?m)(public|private|protected)?\s*(?:static\s+)?(?:\w+\s+)+(\w+)\s*\([^)]*\)\s*(?:throws\s+[\w,\s]+)?\s*\{`)
	matches := methodPattern.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			visibility := "package"
			if match[1] != "" {
				visibility = match[1]
			}
			functions = append(functions, FunctionInfo{
				Name:       match[2],
				Signature:  match[0],
				Visibility: visibility,
			})
		}
	}

	return functions
}

func (p *JavaParser) ExtractChanges(diff string) []CodeChange {
	var changes []CodeChange

	// New methods
	methodPattern := regexp.MustCompile(`(?m)^\+\s*(public|private|protected)?\s*(?:static\s+)?(?:\w+\s+)+(\w+)\s*\(`)
	methods := methodPattern.FindAllStringSubmatch(diff, -1)
	for _, match := range methods {
		if len(match) >= 3 && !isJavaKeyword(match[2]) {
			changes = append(changes, CodeChange{
				Type:        ChangeTypeNewFunction,
				Name:        match[2],
				Description: fmt.Sprintf("Added method %s", match[2]),
			})
		}
	}

	// New classes
	classPattern := regexp.MustCompile(`(?m)^\+\s*(?:public\s+)?(?:abstract\s+)?class\s+(\w+)`)
	classes := classPattern.FindAllStringSubmatch(diff, -1)
	for _, match := range classes {
		changes = append(changes, CodeChange{
			Type:        ChangeTypeNewClass,
			Name:        match[1],
			Description: fmt.Sprintf("Added class %s", match[1]),
		})
	}

	// New interfaces
	interfacePattern := regexp.MustCompile(`(?m)^\+\s*(?:public\s+)?interface\s+(\w+)`)
	interfaces := interfacePattern.FindAllStringSubmatch(diff, -1)
	for _, match := range interfaces {
		changes = append(changes, CodeChange{
			Type:        ChangeTypeNewInterface,
			Name:        match[1],
			Description: fmt.Sprintf("Added interface %s", match[1]),
		})
	}

	// Spring endpoints
	endpointPattern := regexp.MustCompile(`(?m)^\+\s*@(GetMapping|PostMapping|PutMapping|DeleteMapping|RequestMapping)\s*\(\s*(?:value\s*=\s*)?["']([^"']+)["']`)
	endpoints := endpointPattern.FindAllStringSubmatch(diff, -1)
	for _, match := range endpoints {
		changes = append(changes, CodeChange{
			Type:        ChangeTypeNewEndpoint,
			Name:        match[2],
			Description: fmt.Sprintf("Added %s endpoint %s", match[1], match[2]),
		})
	}

	return changes
}

func isJavaKeyword(s string) bool {
	keywords := map[string]bool{
		"if": true, "else": true, "for": true, "while": true,
		"class": true, "interface": true, "return": true, "new": true,
	}
	return keywords[s]
}

// RustParser parses Rust code
type RustParser struct{}

func (p *RustParser) ExtractFunctions(content string) []FunctionInfo {
	var functions []FunctionInfo

	// Match function declarations
	funcPattern := regexp.MustCompile(`(?m)(?:pub\s+)?(?:async\s+)?fn\s+(\w+)\s*[<(]`)
	matches := funcPattern.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) >= 2 {
			visibility := "private"
			if strings.Contains(match[0], "pub ") {
				visibility = "public"
			}
			functions = append(functions, FunctionInfo{
				Name:       match[1],
				Signature:  match[0],
				Visibility: visibility,
			})
		}
	}

	return functions
}

func (p *RustParser) ExtractChanges(diff string) []CodeChange {
	var changes []CodeChange

	// New functions
	funcPattern := regexp.MustCompile(`(?m)^\+\s*(?:pub\s+)?(?:async\s+)?fn\s+(\w+)`)
	funcs := funcPattern.FindAllStringSubmatch(diff, -1)
	for _, match := range funcs {
		changes = append(changes, CodeChange{
			Type:        ChangeTypeNewFunction,
			Name:        match[1],
			Description: fmt.Sprintf("Added function %s", match[1]),
		})
	}

	// New structs
	structPattern := regexp.MustCompile(`(?m)^\+\s*(?:pub\s+)?struct\s+(\w+)`)
	structs := structPattern.FindAllStringSubmatch(diff, -1)
	for _, match := range structs {
		changes = append(changes, CodeChange{
			Type:        ChangeTypeNewClass,
			Name:        match[1],
			Description: fmt.Sprintf("Added struct %s", match[1]),
		})
	}

	// New traits
	traitPattern := regexp.MustCompile(`(?m)^\+\s*(?:pub\s+)?trait\s+(\w+)`)
	traits := traitPattern.FindAllStringSubmatch(diff, -1)
	for _, match := range traits {
		changes = append(changes, CodeChange{
			Type:        ChangeTypeNewInterface,
			Name:        match[1],
			Description: fmt.Sprintf("Added trait %s", match[1]),
		})
	}

	// New impls
	implPattern := regexp.MustCompile(`(?m)^\+\s*impl\s+(?:(\w+)\s+for\s+)?(\w+)`)
	impls := implPattern.FindAllStringSubmatch(diff, -1)
	for _, match := range impls {
		if match[1] != "" {
			changes = append(changes, CodeChange{
				Type:        ChangeTypeNewClass,
				Name:        match[2],
				Description: fmt.Sprintf("Implemented %s for %s", match[1], match[2]),
			})
		}
	}

	return changes
}

func getFileExtension(filename string) string {
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			return filename[i:]
		}
	}
	return ""
}

