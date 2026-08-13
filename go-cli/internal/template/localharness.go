package template

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Param describes one function parameter of the solution method.
type Param struct {
	Name string
	Type string
}

// MethodSig describes the entry-point method of a LeetCode solution.
type MethodSig struct {
	Name       string
	Params     []Param
	ReturnType string
}

// Example holds one parsed Input/Output pair from the problem statement.
type Example struct {
	Input    string
	Expected string
}

// TestCase holds the raw (JSON) argument values and expected (JSON text) output.
type TestCase struct {
	Args     []string
	Expected string
}

// ---------------------------------------------------------------- signature parsing

// ExtractMethodSignature locates the entry method of a LeetCode solution.
// Returns ok=false when the language is not supported or no method was found.
func ExtractMethodSignature(langKey, code string) (MethodSig, bool) {
	switch langKey {
	case "python":
		return extractPythonMethod(code)
	case "cpp", "c":
		return extractCxxMethod(code)
	case "go":
		return extractGoMethod(code)
	case "java":
		return extractJavaMethod(code)
	case "javascript", "typescript":
		return extractJsMethod(code)
	}
	return MethodSig{}, false
}

func extractPythonMethod(code string) (MethodSig, bool) {
	classRe := regexp.MustCompile(`class\s+Solution\s*:`)
	idx := classRe.FindStringIndex(code)
	if idx == nil {
		return MethodSig{}, false
	}
	body := code[idx[1]:]
	defRe := regexp.MustCompile(`(?m)^\s+def\s+([A-Za-z_]\w*)\s*\((.*?)\)\s*(->[^:\n]*)?:`)
	m := defRe.FindStringSubmatch(body)
	if m == nil {
		return MethodSig{}, false
	}
	name := m[1]
	if strings.HasPrefix(name, "_") {
		return MethodSig{}, false
	}
	params := pythonParams(m[2])
	return MethodSig{Name: name, Params: params}, true
}

func pythonParams(s string) []Param {
	var params []Param
	for _, tok := range splitFields(s) {
		tok = strings.TrimSpace(tok)
		if tok == "self" || tok == "cls" || tok == "" {
			continue
		}
		if eq := strings.Index(tok, "="); eq >= 0 {
			tok = tok[:eq]
		}
		if idx := strings.Index(tok, ":"); idx >= 0 {
			tok = tok[:idx]
		}
		tok = strings.TrimSpace(strings.TrimLeft(tok, "*"))
		if tok != "" {
			params = append(params, Param{Name: tok})
		}
	}
	return params
}

func extractCxxMethod(code string) (MethodSig, bool) {
	classRe := regexp.MustCompile(`class\s+Solution\s*\{`)
	idx := classRe.FindStringIndex(code)
	if idx == nil {
		return MethodSig{}, false
	}
	body := code[idx[1]:]
	methodRe := regexp.MustCompile(`(?s)([\w<>{},\[\]\*&\s]+?)\s+([A-Za-z_]\w*)\s*\(([^)]*)\)\s*\{`)
	m := methodRe.FindStringSubmatch(body)
	if m == nil {
		return MethodSig{}, false
	}
	ret := strings.TrimSpace(m[1])
	if strings.Contains(ret, "}") || strings.Contains(ret, "#") {
		return MethodSig{}, false
	}
	var params []Param
	for _, tok := range splitFields(m[3]) {
		tok = strings.TrimSpace(tok)
		if tok == "" {
			continue
		}
		identRe := regexp.MustCompile(`([A-Za-z_]\w*)\s*$`)
		if im := identRe.FindStringSubmatch(tok); im != nil {
			params = append(params, Param{Name: im[1], Type: strings.TrimSpace(tok[:len(tok)-len(im[1])])})
		}
	}
	return MethodSig{Name: m[2], Params: params, ReturnType: ret}, true
}

func extractGoMethod(code string) (MethodSig, bool) {
	re := regexp.MustCompile(`(?m)^\s*func\s+([A-Za-z_]\w*)\s*\((.*?)\)(?:\s+([A-Za-z_\[\]\*\w]+))?\s*\{`)
	ms := re.FindAllStringSubmatch(code, -1)
	for _, m := range ms {
		name := m[1]
		if name == "main" {
			continue
		}
		var params []Param
		for _, tok := range splitFields(m[2]) {
			tok = strings.TrimSpace(tok)
			if tok == "" {
				continue
			}
			fields := strings.Fields(tok)
			if len(fields) < 2 {
				continue
			}
			params = append(params, Param{Name: fields[0], Type: strings.Join(fields[1:], " ")})
		}
		return MethodSig{Name: name, Params: params, ReturnType: strings.TrimSpace(m[3])}, true
	}
	return MethodSig{}, false
}

func extractJavaMethod(code string) (MethodSig, bool) {
	re := regexp.MustCompile(`(?s)public\s+([\w<>\[\],\s]+?)\s+([A-Za-z_]\w*)\s*\(([^)]*)\)\s*\{`)
	m := re.FindStringSubmatch(code)
	if m == nil {
		return MethodSig{}, false
	}
	ret := strings.TrimSpace(m[1])
	if strings.Contains(ret, "{") || strings.Contains(ret, "class") {
		return MethodSig{}, false
	}
	var params []Param
	for _, tok := range splitFields(m[3]) {
		tok = strings.TrimSpace(tok)
		if tok == "" {
			continue
		}
		fields := strings.Fields(tok)
		if len(fields) == 0 {
			continue
		}
		params = append(params, Param{Name: fields[len(fields)-1], Type: strings.Join(fields[:len(fields)-1], " ")})
	}
	return MethodSig{Name: m[2], Params: params, ReturnType: ret}, true
}

func extractJsMethod(code string) (MethodSig, bool) {
	re := regexp.MustCompile(`(?:var|let|const)\s+([A-Za-z_]\w*)\s*=\s*function\s*\(([^)]*)\)|function\s+([A-Za-z_]\w*)\s*\(([^)]*)\)|(?:var|let|const)\s+([A-Za-z_]\w*)\s*=\s*\(([^)]*)\)\s*=>`)
	m := re.FindStringSubmatch(code)
	if m == nil {
		classRe := regexp.MustCompile(`class\s+Solution\s*\{[^}]*?(\w+)\s*\(([^)]*)\)\s*\{`)
		cm := classRe.FindStringSubmatch(code)
		if cm != nil {
			return MethodSig{Name: cm[1], Params: jsParams(cm[2])}, true
		}
		return MethodSig{}, false
	}
	name := firstNonEmpty(m[1], m[3], m[5])
	paramStr := firstNonEmpty(m[2], m[4], m[6])
	if name == "" {
		return MethodSig{}, false
	}
	return MethodSig{Name: name, Params: jsParams(paramStr)}, true
}

func jsParams(s string) []Param {
	var params []Param
	for _, tok := range splitFields(s) {
		tok = strings.TrimSpace(tok)
		if tok == "" || strings.HasPrefix(tok, "...") {
			continue
		}
		if colon := strings.Index(tok, ":"); colon >= 0 {
			tok = tok[:colon]
		}
		if eq := strings.Index(tok, "="); eq >= 0 {
			tok = tok[:eq]
		}
		tok = strings.TrimSpace(strings.TrimLeft(tok, "."))
		if tok != "" {
			params = append(params, Param{Name: tok})
		}
	}
	return params
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

// splitFields splits s on top-level commas, respecting nesting and strings.
func splitFields(s string) []string {
	var parts []string
	depth := 0
	inStr := false
	start := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"':
			inStr = !inStr
		case '[', '{', '(':
			if !inStr {
				depth++
			}
		case ']', '}', ')':
			if !inStr {
				depth--
			}
		default:
			if s[i] == ',' && depth == 0 && !inStr {
				parts = append(parts, s[start:i])
				start = i + 1
			}
		}
	}
	parts = append(parts, s[start:])
	return parts
}

// CleanType returns a normalized, dereferenced type token.
func CleanType(t string) string {
	t = strings.TrimSpace(t)
	t = strings.TrimPrefix(t, "const ")
	t = strings.TrimSuffix(strings.TrimSpace(t), "*")
	t = strings.TrimSuffix(strings.TrimSpace(t), "&")
	t = strings.TrimSpace(t)
	return t
}

var stripTagsRe = regexp.MustCompile(`<[^>]+>`)

// ---------------------------------------------------------------- example parsing

var preBlockRe = regexp.MustCompile(`(?s)<pre[^>]*>(.*?)</pre>`)
var strongTagRe = regexp.MustCompile(`(?s)<strong[^>]*>(.*?)</strong>`)

// ExtractExamples pulls (Input, Output) pairs from the problem description HTML.
func ExtractExamples(contentHTML string) []Example {
	var examples []Example
	for _, block := range preBlockRe.FindAllString(contentHTML, -1) {
		lower := strings.ToLower(block)
		inIdx := strings.Index(lower, "input:")
		if inIdx < 0 {
			continue
		}
		outIdx := strings.Index(lower[inIdx:], "output:")
		if outIdx < 0 {
			continue
		}
		outIdx += inIdx
		inText := block[inIdx+len("input:"):outIdx]
		outText := block[outIdx+len("output:"):]
		if m := strongTagRe.FindStringIndex(outText); m != nil {
			outText = outText[:m[0]]
		}
		inText = cleanHTMLText(inText)
		outText = cleanHTMLText(outText)
		if inText == "" || outText == "" {
			continue
		}
		examples = append(examples, Example{Input: inText, Expected: outText})
	}
	return examples
}

func cleanHTMLText(s string) string {
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = stripTagsRe.ReplaceAllString(s, "")
	return strings.TrimSpace(s)
}

// ExtractExamplesMarkdown pulls (Input, Output) pairs from a plain-text /
// markdown problem description — the form saved in each problem's README.md by
// `add`. Unlike ExtractExamples it does not require the original HTML <pre>
// blocks, so example test cases can be built fully offline.
func ExtractExamplesMarkdown(text string) []Example {
	var examples []Example
	lines := strings.Split(text, "\n")

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if !isInputLine(line) {
			continue
		}
		inVal := valueAfterColon(line, "input:")
		j := i
		for j+1 < len(lines) {
			next := strings.TrimSpace(lines[j+1])
			if startsExampleSection(next) {
				break
			}
			j++
			inVal += " " + next
		}
		inVal = strings.TrimSpace(inVal)

		outVal := ""
		last := j
		for k := j + 1; k < len(lines); k++ {
			l := strings.TrimSpace(lines[k])
			if isOutputLine(l) {
				outVal = valueAfterColon(l, "output:")
				outVal = strings.TrimSpace(outVal)
				kk := k
				for kk+1 < len(lines) {
					next := strings.TrimSpace(lines[kk+1])
					if startsExampleSection(next) {
						break
					}
					kk++
					outVal += " " + next
				}
				outVal = strings.TrimSpace(outVal)
				last = kk
				break
			}
			if startsExampleSection(l) {
				break
			}
		}
		i = last

		inVal = cleanHTMLText(inVal)
		outVal = cleanHTMLText(outVal)
		if inVal != "" && outVal != "" {
			examples = append(examples, Example{Input: inVal, Expected: outVal})
		}
	}
	return examples
}

func isInputLine(s string) bool {
	return strings.HasPrefix(strings.ToLower(s), "input:")
}

func isOutputLine(s string) bool {
	return strings.HasPrefix(strings.ToLower(s), "output:")
}

func valueAfterColon(s, key string) string {
	lower := strings.ToLower(s)
	idx := strings.Index(lower, key)
	if idx < 0 {
		return ""
	}
	return s[idx+len(key):]
}

// startsExampleSection reports whether a line begins a section that must not be
// merged into an Input/Output value (e.g. the next example, explanation, or a
// blank / &nbsp; separator line).
func startsExampleSection(s string) bool {
	if strings.TrimSpace(strings.ReplaceAll(s, "&nbsp;", " ")) == "" {
		return true
	}
	lower := strings.ToLower(s)
	if isInputLine(s) || isOutputLine(s) {
		return true
	}
	for _, p := range []string{"explanation:", "constraints:", "example ", "note:", "follow up:", "follow-up:"} {
		if strings.HasPrefix(lower, p) {
			return true
		}
	}
	return false
}

// ParseExampleInput converts "nums = [2,7,11,15], target = 9" into name -> raw value.
func ParseExampleInput(inputText string) map[string]string {
	result := make(map[string]string)
	for _, part := range splitFields(inputText) {
		eq := strings.Index(part, "=")
		if eq < 0 {
			continue
		}
		name := strings.TrimSpace(part[:eq])
		val := strings.TrimSpace(part[eq+1:])
		if name != "" && val != "" {
			result[name] = val
		}
	}
	return result
}

// BuildTestCases aligns parsed examples to the method parameter order.
func BuildTestCases(sig MethodSig, examples []Example) []TestCase {
	var cases []TestCase
	for _, ex := range examples {
		raw := ParseExampleInput(ex.Input)
		if len(raw) == 0 {
			continue
		}
		args := make([]string, 0, len(sig.Params))
		if len(sig.Params) > 0 {
			matched := true
			for _, p := range sig.Params {
				v, ok := raw[p.Name]
				if !ok {
					matched = false
					break
				}
				args = append(args, v)
			}
			if !matched {
				args = args[:0]
				for _, v := range raw {
					args = append(args, v)
				}
			}
		} else {
			args = append(args, ex.Input)
		}
		cases = append(cases, TestCase{Args: args, Expected: ex.Expected})
	}
	return cases
}

// ---------------------------------------------------------------- harness generation

// BuildLocalHarness generates a self-contained program that embeds the user's
// solution, reads the testcases JSON file (argv[1]) and prints one
// "RESULT\t<json>" line (or "ERROR\t<msg>") per test case.
func BuildLocalHarness(langKey, solutionCode string, sig MethodSig, cases []TestCase) (string, bool) {
	switch langKey {
	case "python":
		return buildPythonHarness(solutionCode, sig)
	case "javascript", "typescript":
		return buildJsHarness(langKey, solutionCode, sig)
	case "go":
		return buildGoHarness(solutionCode, sig)
	case "cpp", "c":
		return buildCppHarness(solutionCode, sig, cases)
	case "java":
		return buildJavaHarness(solutionCode, sig, cases)
	}
	return "", false
}

func buildPythonHarness(code string, sig MethodSig) (string, bool) {
	sb := &strings.Builder{}
	sb.WriteString("import json\nimport sys\n\n")
	sb.WriteString(strings.TrimSpace(code))
	sb.WriteString("\n\n_sol = Solution()\n")
	sb.WriteString("def _main():\n")
	sb.WriteString("    with open(sys.argv[1], 'r', encoding='utf-8') as _f:\n")
	sb.WriteString("        _tests = json.load(_f)\n")
	sb.WriteString("    for _t in _tests:\n")
	sb.WriteString("        try:\n")
	sb.WriteString("            _args = [json.loads(_a) for _a in _t['args']]\n")
	fmt.Fprintf(sb, "            _res = _sol.%s(*_args)\n", sig.Name)
	sb.WriteString("            print('RESULT\\t' + json.dumps(_res, default=str))\n")
	sb.WriteString("        except Exception as _e:\n")
	sb.WriteString("            print('ERROR\\t' + repr(_e))\n")
	sb.WriteString("\nif __name__ == '__main__':\n    _main()\n")
	return sb.String(), true
}

func buildJsHarness(langKey, code string, sig MethodSig) (string, bool) {
	useClass := strings.Contains(code, "class Solution")
	sb := &strings.Builder{}
	sb.WriteString("const fs = require('fs');\n")
	if langKey == "typescript" {
		sb.WriteString("// @ts-ignore\n")
	}
	sb.WriteString(strings.TrimSpace(code))
	sb.WriteString("\n\nconst _tests = JSON.parse(fs.readFileSync(process.argv[2], 'utf8'));\n")
	sb.WriteString("for (const _t of _tests) {\n")
	sb.WriteString("  try {\n")
	sb.WriteString("    const _args = _t.args.map((s) => JSON.parse(s));\n")
	invoke := sig.Name + "(..._args)"
	if useClass {
		invoke = "new Solution()." + sig.Name + "(..._args)"
	}
	fmt.Fprintf(sb, "    const _res = %s;\n", invoke)
	sb.WriteString("    console.log('RESULT\\t' + JSON.stringify(_res));\n")
	sb.WriteString("  } catch (e) {\n")
	sb.WriteString("    console.log('ERROR\\t' + String((e && e.stack) || e));\n")
	sb.WriteString("  }\n")
	sb.WriteString("}\n")
	return sb.String(), true
}

func buildGoHarness(code string, sig MethodSig) (string, bool) {
	stripped := removeGoMain(code)
	sb := &strings.Builder{}
	sb.WriteString(stripped)
	sb.WriteString("\n\nfunc leetMain() {\n")
	sb.WriteString("\tvar tests []struct {\n")
	sb.WriteString("\t\tArgs []string `json:\"args\"`\n")
	sb.WriteString("\t}\n")
	sb.WriteString("\tf, _ := os.Open(os.Args[1])\n")
	sb.WriteString("\td, _ := io.ReadAll(f)\n")
	sb.WriteString("\t_ = json.Unmarshal(d, &tests)\n")
	sb.WriteString("\tfor _, t := range tests {\n")
	for i, p := range sig.Params {
		goT := goTypeFor(p.Type)
		if goT == "" {
			return "", false
		}
		fmt.Fprintf(sb, "\t\tvar %s %s\n", p.Name, goT)
		fmt.Fprintf(sb, "\t\t_ = json.Unmarshal([]byte(t.Args[%d]), &%s)\n", i, p.Name)
	}
	sb.WriteString("\t\t")
	fmt.Fprintf(sb, "res := %s(", sig.Name)
	for i, p := range sig.Params {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(p.Name)
	}
	sb.WriteString(")\n")
	sb.WriteString("\t\tb, _ := json.Marshal(res)\n")
	sb.WriteString("\t\tfmt.Println(\"RESULT\\t\" + string(b))\n")
	sb.WriteString("\t}\n")
	sb.WriteString("}\n\n")
	sb.WriteString("func main() {\n\tleetMain()\n}\n")
	return sb.String(), true
}

func goTypeFor(t string) string {
	t = CleanType(t)
	switch t {
	case "int", "int32":
		return "int"
	case "int64", "long":
		return "int64"
	case "float64", "double":
		return "float64"
	case "bool":
		return "bool"
	case "string":
		return "string"
	case "byte":
		return "byte"
	case "List[int]", "list[int]":
		return "[]int"
	case "List[string]", "list[string]":
		return "[]string"
	case "List[float]", "list[float]":
		return "[]float64"
	case "List[bool]", "list[bool]":
		return "[]bool"
	case "List[List[int]]", "list[list[int]]":
		return "[][]int"
	case "List[List[string]]", "list[list[string]]":
		return "[][]string"
	}
	if strings.HasPrefix(t, "[]") || strings.HasPrefix(t, "map[") {
		return t
	}
	return ""
}

func removeGoMain(code string) string {
	idx := strings.Index(code, "func main(")
	if idx < 0 {
		return code
	}
	open := strings.Index(code[idx:], "{")
	if open < 0 {
		return code
	}
	start := idx + open
	depth := 0
	for i := start; i < len(code); i++ {
		switch code[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				head := code[:idx]
				tail := strings.TrimLeft(code[i+1:], "\n")
				return head + "\n" + tail
			}
		}
	}
	return code[:idx]
}

// ---------------------------------------------------------------- literal helpers

func jsonStringToGo(token string) (string, bool) {
	token = strings.TrimSpace(token)
	if !strings.HasPrefix(token, `"`) {
		return "", false
	}
	var out string
	if err := json.Unmarshal([]byte(token), &out); err != nil {
		return "", false
	}
	return out, true
}

// cppElement converts a JSON token into a C++ braced-init / literal expression.
func cppElement(token string) (string, bool) {
	token = strings.TrimSpace(token)
	switch {
	case token == "":
		return "", false
	case token == "true", token == "false":
		return token, true
	case token == "null":
		return "nullptr", true
	case strings.HasPrefix(token, `"`):
		s, ok := jsonStringToGo(token)
		if !ok {
			return "", false
		}
		return strconv.Quote(s), true
	case strings.HasPrefix(token, "["):
		elements, ok := splitJSONArray(token)
		if !ok {
			return "", false
		}
		out := make([]string, 0, len(elements))
		for _, e := range elements {
			lit, ok := cppElement(e)
			if !ok {
				return "", false
			}
			out = append(out, lit)
		}
		return "{" + strings.Join(out, ", ") + "}", true
	default:
		return token, true
	}
}

func splitJSONArray(token string) ([]string, bool) {
	token = strings.TrimSpace(token)
	if !strings.HasPrefix(token, "[") || !strings.HasSuffix(token, "]") {
		return nil, false
	}
	inner := token[1 : len(token)-1]
	if strings.TrimSpace(inner) == "" {
		return []string{}, true
	}
	return splitFields(inner), true
}

func buildCppHarness(code string, sig MethodSig, cases []TestCase) (string, bool) {
	hasListNode := strings.Contains(code, "ListNode")
	hasTreeNode := strings.Contains(code, "TreeNode")

	sb := &strings.Builder{}
	sb.WriteString("#include <bits/stdc++.h>\nusing namespace std;\n\n")
	sb.WriteString(cppPrinterSrc)
	if hasListNode {
		sb.WriteString(cppListNodePrinter)
	}
	if hasTreeNode {
		sb.WriteString(cppTreeNodePrinter)
	}
	sb.WriteString("\n")
	sb.WriteString(strings.TrimSpace(code))
	sb.WriteString("\n\nint main() {\n    Solution __sol;\n")
	for _, tc := range cases {
		sb.WriteString("    {\n")
		varNames := make([]string, 0, len(tc.Args))
		for i, arg := range tc.Args {
			lit, ok := cppElement(arg)
			if !ok {
				return "", false
			}
			typ := cppTypeName(sig.Params[i].Type, lit)
			if typ == "" {
				return "", false
			}
			name := fmt.Sprintf("__a%d", i)
			fmt.Fprintf(sb, "        %s %s = %s;\n", typ, name, lit)
			varNames = append(varNames, name)
		}
		fmt.Fprintf(sb, "        auto __r = __sol.%s(%s);\n", sig.Name, strings.Join(varNames, ", "))
		sb.WriteString("        cout << \"RESULT\\t\" << __j(__r) << \"\\n\";\n")
		sb.WriteString("    }\n")
	}
	sb.WriteString("    return 0;\n}\n")
	return sb.String(), true
}

// cppTypeName maps a parameter type to a C++ type name used in local variable
// declarations. Strips references/value-qualifiers and known containers get
// normalized so braced-init initializers bind correctly.
func cppTypeName(declType, literal string) string {
	t := CleanType(declType)
	lower := strings.ToLower(t)
	switch {
	case strings.Contains(lower, "listnode"):
		return strings.TrimSpace(t)
	case strings.Contains(lower, "treenode"):
		return strings.TrimSpace(t)
	case strings.HasPrefix(lower, "vector<"):
		// e.g. "vector<int>&" -> "vector<int>"
		return t
	case strings.HasPrefix(lower, "string"):
		return "string"
	case lower == "int", lower == "long", lower == "long long", lower == "int64_t",
		lower == "double", lower == "float", lower == "bool", lower == "char":
		return t
	}
	// fallback: if literal starts with '{' it's a container/array; use declared text
	return t
}

const cppPrinterSrc = `
string __j(const string& s) {
    string o = "\"";
    for (char c : s) {
        if (c == '"' || c == '\\') o += '\\';
        o += c;
    }
    o += "\"";
    return o;
}
string __j(bool b) { return b ? "true" : "false"; }
string __j(int v) { return to_string(v); }
string __j(long long v) { return to_string(v); }
string __j(char v) { return string("\"") + v + "\""; }
string __j(double v) {
    long long i = (long long)v;
    if ((double)i == v) return to_string(i);
    string s = to_string(v);
    return s;
}
string __j(float v) { return __j((double)v); }
template<class T> string __j(const vector<T>& v) {
    string o = "[";
    for (size_t i = 0; i < v.size(); ++i) {
        if (i) o += ",";
        o += __j(v[i]);
    }
    return o + "]";
}
`

const cppListNodePrinter = `
string __j(const ListNode* head) {
    string o = "[";
    const ListNode* p = head;
    int i = 0;
    while (p) {
        if (i++) o += ",";
        o += to_string(p->val);
        p = p->next;
    }
    return o + "]";
}
`

const cppTreeNodePrinter = `
string __j(const TreeNode* root) {
    string o = "[";
    bool first = true;
    function<void(const TreeNode*)> walk = [&](const TreeNode* n) {
        if (!first) o += ",";
        first = false;
        if (!n) { o += "null"; return; }
        o += to_string(n->val);
        walk(n->left);
        walk(n->right);
    };
    if (root) walk(root);
    o += "]";
    return o;
}
`

func buildJavaHarness(code string, sig MethodSig, cases []TestCase) (string, bool) {
	sb := &strings.Builder{}
	sb.WriteString("import java.util.*;\nimport java.lang.reflect.*;\n\npublic class Harness {\n")
	sb.WriteString(javaPrinterSrc)
	sb.WriteString("  public static void main(String[] args) throws Exception {\n")
	sb.WriteString("    Solution sol = new Solution();\n")
	for _, tc := range cases {
		sb.WriteString("    {\n")
		names := make([]string, 0, len(tc.Args))
		for i, arg := range tc.Args {
			typ := sig.Params[i].Type
			expr, ok := javaArgLiteral(typ, arg)
			if !ok {
				return "", false
			}
			name := fmt.Sprintf("__a%d", i)
			fmt.Fprintf(sb, "      %s %s = %s;\n", javaTypeName(typ), name, expr)
			names = append(names, name)
		}
		fmt.Fprintf(sb, "      printResult(__j(sol.%s(%s)));\n", sig.Name, strings.Join(names, ", "))
		sb.WriteString("    }\n")
	}
	sb.WriteString("  }\n\n")
	sb.WriteString(javaPrintResultSrc)
	sb.WriteString("}\n")
	return sb.String(), true
}

const javaPrinterSrc = `
  static String __j(Object v) {
    if (v == null) return "null";
    if (v instanceof Integer || v instanceof Long || v instanceof Short) return String.valueOf(v);
    if (v instanceof Double || v instanceof Float) {
      double d = ((Number) v).doubleValue();
      if (d == Math.floor(d) && !Double.isInfinite(d)) return String.valueOf((long) d);
      return String.valueOf(d);
    }
    if (v instanceof Boolean) return String.valueOf(v);
    if (v instanceof Character) return "\"" + v + "\"";
    if (v instanceof String) {
      String s = (String) v;
      String q = s.replace("\\", "\\\\").replace("\"", "\\\"");
      return "\"" + q + "\"";
    }
    if (v.getClass().isArray()) {
      int n = Array.getLength(v);
      StringBuilder o = new StringBuilder("[");
      for (int i = 0; i < n; i++) {
        if (i > 0) o.append(',');
        o.append(__j(Array.get(v, i)));
      }
      return o.append(']').toString();
    }
    if (v instanceof List) {
      List<?> l = (List<?>) v;
      StringBuilder o = new StringBuilder("[");
      for (int i = 0; i < l.size(); i++) {
        if (i > 0) o.append(',');
        o.append(__j(l.get(i)));
      }
      return o.append(']').toString();
    }
    return String.valueOf(v);
  }
`

const javaPrintResultSrc = `
  static void printResult(String line) {
    System.out.println("RESULT\\t" + line);
  }
`

// javaTypeName maps a Java LeetCode lambda type to a concrete Java type.
func javaTypeName(t string) string {
	t = strings.TrimSpace(t)
	switch t {
	case "int", "long", "double", "boolean", "char", "float", "short", "byte", "String":
		return t
	case "int[]":
		return "int[]"
	case "long[]":
		return "long[]"
	case "double[]":
		return "double[]"
	case "boolean[]":
		return "boolean[]"
	case "char[]":
		return "char[]"
	case "String[]":
		return "String[]"
	}
	if strings.HasPrefix(t, "List<") && strings.HasSuffix(t, ">") {
		return t
	}
	return "Object"
}

// javaArgLiteral returns the right-hand-side expression to build a Java value
// from a JSON token, matched to the parameter's declared type.
func javaArgLiteral(paramType, jsonToken string) (string, bool) {
	typ := strings.TrimSpace(paramType)
	switch typ {
	case "int":
		return "int", true // caller prepends variable name
	case "String":
		s, ok := jsonStringToGo(jsonToken)
		if !ok {
			return "", false
		}
		return strconv.Quote(s), true
	case "boolean":
		return strings.TrimSpace(jsonToken), true
	case "long":
		return strings.TrimSpace(jsonToken), true
	case "double":
		return strings.TrimSpace(jsonToken), true
	}

	// arrays
	if strings.HasSuffix(typ, "[]") {
		base := strings.TrimSuffix(typ, "[]")
		elements, ok := splitJSONArray(jsonToken)
		if !ok {
			return "", false
		}
		inner := make([]string, 0, len(elements))
		for _, e := range elements {
			if base == "String" {
				s, ok := jsonStringToGo(e)
				if !ok {
					return "", false
				}
				inner = append(inner, strconv.Quote(s))
			} else if base == "List" {
				// nested List handled below; treat base recursively as List element
				lit, ok := javaArgLiteral("List", e)
				if !ok {
					return "", false
				}
				inner = append(inner, lit)
			} else {
				inner = append(inner, strings.TrimSpace(e))
			}
		}
		return "new " + typ + "{" + strings.Join(inner, ", ") + "}", true
	}

	if typ == "List" || strings.HasPrefix(typ, "List<") {
		return javaListLiteral(jsonToken)
	}

	// default fallback: raw number or expression
	if strings.HasPrefix(strings.TrimSpace(jsonToken), `"`) {
		s, ok := jsonStringToGo(jsonToken)
		if !ok {
			return "", false
		}
		return strconv.Quote(s), true
	}
	return strings.TrimSpace(jsonToken), true
}

func javaListLiteral(jsonToken string) (string, bool) {
	elements, ok := splitJSONArray(jsonToken)
	if !ok {
		return "", false
	}
	inner := make([]string, 0, len(elements))
	for _, e := range elements {
		lit, ok := javaArgLiteral("", e)
		if !ok {
			return "", false
		}
		inner = append(inner, lit)
	}
	return "Arrays.asList(" + strings.Join(inner, ", ") + ")", true
}