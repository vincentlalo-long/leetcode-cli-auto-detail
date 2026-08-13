package template

import (
	"strings"
	"testing"
)

func TestExtractMethodSignaturePython(t *testing.T) {
	code := `from typing import List

class Solution:
    def twoSum(self, nums: List[int], target: int) -> List[int]:
        pass
`
	sig, ok := ExtractMethodSignature("python", code)
	if !ok {
		t.Fatal("expected signature")
	}
	if sig.Name != "twoSum" {
		t.Fatalf("got %s", sig.Name)
	}
	if len(sig.Params) != 2 || sig.Params[0].Name != "nums" || sig.Params[1].Name != "target" {
		t.Fatalf("unexpected params: %+v", sig.Params)
	}
}

func TestExtractMethodSignatureCpp(t *testing.T) {
	code := `class Solution {
public:
    vector<int> twoSum(vector<int>& nums, int target) {
    }
};`
	sig, ok := ExtractMethodSignature("cpp", code)
	if !ok {
		t.Fatal("expected signature")
	}
	if sig.Name != "twoSum" {
		t.Fatalf("got %s", sig.Name)
	}
	if sig.ReturnType != "vector<int>" {
		t.Fatalf("return %s", sig.ReturnType)
	}
	if len(sig.Params) != 2 {
		t.Fatalf("params: %+v", sig.Params)
	}
	if sig.Params[0].Name != "nums" || strings.TrimSpace(sig.Params[0].Type) != "vector<int>&" {
		t.Fatalf("param0: %+v", sig.Params[0])
	}
}

func TestExtractMethodSignatureGo(t *testing.T) {
	code := `package main

func twoSum(nums []int, target int) []int {
	return nil
}`
	sig, ok := ExtractMethodSignature("go", code)
	if !ok {
		t.Fatal("expected signature")
	}
	if sig.Name != "twoSum" || len(sig.Params) != 2 {
		t.Fatalf("unexpected: %+v", sig)
	}
	if sig.Params[0].Type != "[]int" {
		t.Fatalf("param0 type: %s", sig.Params[0].Type)
	}
}

func TestExtractExamplesTwoSum(t *testing.T) {
	content := `<pre>
<strong>Input:</strong> nums = [2,7,11,15], target = 9
<strong>Output:</strong> [0,1]
<strong>Explanation:</strong> Because nums[0] + nums[1] == 9, we return [0, 1].
</pre>

<pre>
<strong>Input:</strong> nums = [3,2,4], target = 6
<strong>Output:</strong> [1,2]
</pre>`
	examples := ExtractExamples(content)
	if len(examples) != 2 {
		t.Fatalf("got %d examples: %+v", len(examples), examples)
	}
	if !strings.Contains(examples[0].Input, "nums") || examples[0].Expected != "[0,1]" {
		t.Fatalf("example0: %+v", examples[0])
	}
	if examples[1].Expected != "[1,2]" {
		t.Fatalf("example1: %+v", examples[1])
	}
}

func TestExtractExamplesMarkdown(t *testing.T) {
	content := `# [1. Two Sum](https://leetcode.com/problems/two-sum/)

- **Difficulty:** Easy
- **Tags:** Array, Hash Table

## Description

You are given an array of integers.

Example 1:

Input: nums = [2,7,11,15], target = 9
Output: [0,1]
Explanation: Because nums[0] + nums[1] == 9, we return [0, 1].

Example 2:

Input: nums = [3,2,4], target = 6
Output: [1,2]

Example 3:

Input: s = "abc"
Output: 3

Constraints:

2 &lt;= nums.length &lt;= 104
`
	examples := ExtractExamplesMarkdown(content)
	if len(examples) != 3 {
		t.Fatalf("got %d examples: %+v", len(examples), examples)
	}
	if examples[0].Input != "nums = [2,7,11,15], target = 9" || examples[0].Expected != "[0,1]" {
		t.Fatalf("example0: %+v", examples[0])
	}
	if examples[1].Input != "nums = [3,2,4], target = 6" || examples[1].Expected != "[1,2]" {
		t.Fatalf("example1: %+v", examples[1])
	}
	if examples[2].Input != `s = "abc"` || examples[2].Expected != "3" {
		t.Fatalf("example2: %+v", examples[2])
	}
}

func TestExtractExamplesMarkdownNoInput(t *testing.T) {
	examples := ExtractExamplesMarkdown("Just a description without any examples.")
	if len(examples) != 0 {
		t.Fatalf("expected no examples, got %+v", examples)
	}
}

func TestBuildTestCasesTwoSum(t *testing.T) {
	sig := MethodSig{Name: "twoSum", Params: []Param{{Name: "nums"}, {Name: "target"}}}
	examples := []Example{
		{Input: "nums = [2,7,11,15], target = 9", Expected: "[0,1]"},
		{Input: "nums = [3,3], target = 6", Expected: "[0,1]"},
	}
	tt := BuildTestCases(sig, examples)
	if len(tt) != 2 {
		t.Fatalf("got %d cases", len(tt))
	}
	if tt[0].Args[0] != "[2,7,11,15]" || tt[0].Args[1] != "9" || tt[0].Expected != "[0,1]" {
		t.Fatalf("case0: %+v", tt[0])
	}
}

func TestBuildPythonHarness(t *testing.T) {
	code := `class Solution:
    def add(self, a: int, b: int) -> int:
        return a + b`
	sig := MethodSig{Name: "add", Params: []Param{{Name: "a"}, {Name: "b"}}}
	h, ok := BuildLocalHarness("python", code, sig, nil)
	if !ok {
		t.Fatal("harness build failed")
	}
	if !strings.Contains(h, "_sol.add(*_args)") {
		t.Fatalf("missing call in harness:\n%s", h)
	}
}

func TestGoTypeFor(t *testing.T) {
	cases := map[string]string{
		"int":             "int",
		"[]int":           "[]int",
		"List[int]":       "[]int",
		"List[List[int]]": "[][]int",
		"string":          "string",
		"ListNode":        "",
	}
	for in, want := range cases {
		if got := goTypeFor(in); got != want {
			t.Errorf("goTypeFor(%q) = %q, want %q", in, got, want)
		}
	}
}