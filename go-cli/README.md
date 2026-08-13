# Leet CLI — LeetCode Workspace Manager

A terminal (TUI) + command-line tool to organize, solve, test, submit, and review
LeetCode problems entirely from your terminal. Problems are stored as plain files
on disk organized by data structure, so your whole history is version-controlled
with Git.

## ✨ Features

- **TUI** (interactive) and **non-interactive** CLI mode: `leet list --difficulty Medium`
- **Add problems** from LeetCode by number, daily challenge, or random pick
- **Local run** with your installed compilers (C, C++, Python, Go, Java, Rust, JS, TS, C#)
- **Local example tests without cookies**: `leet test <num> --local` builds a harness
  from the problem statement's example Input/Output blocks (read from the local
  `README.md` first, falling back to LeetCode) and compiles/runs your
  solution entirely on your machine (Python, JS/TS, Go, C/C++, Java)
- **Test & Submit** to the LeetCode API (requires cookies) with progress tracking
- **Progress tracker** — submissions and local passes are recorded in `.leet/progress.json`
- **Review queue** — spaced-repetition (1/3/7/15/30/60 days) for revision
- **Difficulty filtering** in `list`
- **Root README index** auto-generated on `sync` / `readme`
- **`help` / `man` in both modes** — `leet help` lists commands, `leet help <cmd>` or `leet <cmd> --help` show usage
- **`sync --progress`** — also backs up `.leet/progress.json` (solving history + review schedule)
- **`version` command** — prints the current CLI version

## 📦 Install / Build

Requires [Go 1.25+](https://go.dev/dl/).

```bash
cd go-cli
go build -o leet.exe .
```

Move `leet.exe` somewhere on your `PATH` (e.g. `C:\Users\you\bin`) so you can run
`leet` from anywhere, or run it from the repo:

```bash
# on Windows
go-cli\leet.exe
# on macOS / Linux
./go-cli/leet
```

## ⚙️ Configuration

The tool looks for `config.json` in the current directory, then walks up parent
directories. A committed **`config.json`** holds the safe defaults; a git-ignored
**`config.local.json`** holds your machine-specific settings (workspace path,
editor, LeetCode cookies). Settings can also be overridden by environment
variables.

| Setting | config key | env var |
|---|---|---|
| Workspace directory | `base_dir` | – |
| Default language | `default_language` | – |
| Text editor | `editor` | – |
| LeetCode session cookie | `leetcode_session` | `LEETCODE_SESSION` |
| LeetCode CSRF token | `leetcode_csrf` | `LEETCODE_CSRF` |

`base_dir` is auto-detected to the nearest Git repo root when empty or invalid, so
clones work on any machine without editing the file.

### Local testing without cookies (`test <num> --local`)

The API test/submit commands need your LeetCode cookies. To verify a solution
**without** going online, pass `--local`:

```bash
leet test 1 --local
```

This:
1. parses the entry method of your `class Solution` from the solution file
2. extracts the example **Input/Output** blocks (from the problem's local
   `README.md` when present, otherwise from the LeetCode description)
3. generates a small harness that calls your method with each example
4. compiles/runs it with your local toolchain and compares the results

Supported languages: Python, JavaScript, TypeScript, Go, C/C++, Java. When all
example cases pass, the problem is automatically marked solved in the tracker.

### Setting up LeetCode cookies (for `test` / `submit`)

1. Log in to leetcode.com in your browser.
2. Open DevTools → Application/Storage → Cookies.
3. Copy the values of `LEETCODE_SESSION` and `csrftoken`.

Then either export them per shell:

```powershell
$env:LEETCODE_SESSION = "your-session-cookie"
$env:LEETCODE_CSRF    = "your-csrf-token"
```

or store them in the git-ignored local config (keeps secrets out of your repo):

```json
{
  "base_dir": "D:/my-workspace",
  "leetcode_session": "your-session-cookie",
  "leetcode_csrf": "your-csrf-token"
}
```

## 🚀 Quick Start

```bash
leet                    # launch interactive TUI
leet help               # list all commands (also in the TUI)
leet help add           # usage for a single command
leet add 1              # add Two Sum (fetches from LeetCode API)
leet list --difficulty Easy
leet run 1              # compile & run locally
leet test 1 --local     # run example test cases locally (no cookies needed)
leet test 1             # run LeetCode testcases (needs cookies)
leet submit 1           # submit & track progress
leet daily --add        # add today's challenge non-interactively
leet review             # review queue (spaced repetition)
leet review --solve 1   # mark problem solved
leet readme             # regenerate root README index
leet sync               # commit & push workspace to GitHub
leet sync --progress    # ... and back up .leet/progress.json too
leet version            # print CLI version
```

## 📁 Workspace Layout

```
<base_dir>/
├── README.md                  # auto-generated index table
├── config.json                # committed safe defaults
├── config.local.json          # git-ignored local settings (cookies, base_dir)
├── .leet/progress.json        # git-ignored solving progress + review schedule
├── array/
│   └── 1-two-sum/
│       ├── 1_Two Sum.cpp
│       └── README.md          # problem description (auto-fetched)
├── string/
└── ...
```

## 🧪 Tests

```bash
cd go-cli
go test ./...
```

## 📖 Command Manual

Type `help` (or `man <cmd>`) inside the TUI, or run `leet help` / `leet help <cmd>`
/ `leet <cmd> --help` from the shell to see usage for any command.
