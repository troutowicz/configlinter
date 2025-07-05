package configlinter

import (
	"go/ast"
	"go/token"
	"strconv"

	"github.com/spf13/viper"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "configlinter",
	Doc:      "Check that all config keys used in the codebase are defined in the config",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func init() {
	viper.SetConfigName("config.json.template")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	_ = viper.ReadInConfig() // Ignore error as this is best-effort
}

func run(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		// Check if this is a config function call
		if !isConfigCall(call) {
			return
		}

		// Extract the config key from the first argument
		if len(call.Args) == 0 {
			return
		}

		key, lit := extractStringLiteral(call.Args[0])
		if lit == nil {
			// Report non-literal config keys as potential issues
			pass.Report(analysis.Diagnostic{
				Pos:     call.Args[0].Pos(),
				End:     call.Args[0].End(),
				Message: "config key should be a string literal for static analysis",
			})
			return
		}

		// Check if the key is valid using viper.IsSet
		if !viper.IsSet(key) {
			pass.Report(analysis.Diagnostic{
				Pos:     lit.Pos(),
				End:     lit.End(),
				Message: "config key " + `"` + key + `"` + " is not defined in config",
			})
		}
	})

	return nil, nil
}

// isConfigCall checks if the function call is a config.Get* or viper.Get* call
func isConfigCall(call *ast.CallExpr) bool {
	switch fun := call.Fun.(type) {
	case *ast.SelectorExpr:
		// Check for config.GetString, config.GetBool, config.GetStringSlice
		if ident, ok := fun.X.(*ast.Ident); ok {
			if ident.Name == "config" {
				switch fun.Sel.Name {
				case "GetString", "GetBool", "GetStringSlice":
					return true
				}
			}
			// Check for viper.GetString, viper.GetBool, viper.GetStringSlice
			if ident.Name == "viper" {
				switch fun.Sel.Name {
				case "GetString", "GetBool", "GetStringSlice":
					return true
				}
			}
		}
	}
	return false
}

// extractStringLiteral extracts the string value and literal node from a string literal expression
func extractStringLiteral(expr ast.Expr) (string, *ast.BasicLit) {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", nil
	}
	unquoted, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", nil
	}
	return unquoted, lit
}
