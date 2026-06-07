import os
import re
from typing import Any, Dict

from cli.utils.config_manager import ConfigManager
from cli.utils.file_utils import get_all_solution_files
from cli.utils.language_support import get_language_by_extension
from cli.utils.leetcode_api import (
    get_problem_details,
    get_leetcode_auth,
    submit_solution,
)
from cli.utils.ui import (
    print_command_banner, print_success, print_error, print_info,
    styled_text_input, styled_select, console
)

LEETCODE_LANG_MAP = {
    "cpp": "cpp",
    "c": "c",
    "python": "python3",
    "go": "golang",
}


def _matches_problem_number(filename: str, problem_num: str) -> bool:
    if filename.startswith(f"{problem_num}_"):
        return True
    if problem_num.isdigit():
        num = int(problem_num)
        return filename.startswith(f"{num:03d}_") or filename.startswith(f"{num:04d}_")
    return False


def _infer_slug_from_path(file_path: str, content: str) -> str:
    parent_name = os.path.basename(os.path.dirname(file_path))
    match = re.match(r"^\d+-(.+)$", parent_name)
    if match:
        return match.group(1)

    link_match = re.search(r"Link:\s*(https?://leetcode.com/problems/([^/]+)/)", content)
    if link_match:
        return link_match.group(2)

    return ""


def _render_benchmark(result: Dict[str, Any]):
    status_msg = result.get("status_msg") or "Unknown"
    status_code = result.get("status_code")
    success = status_code == 10 or status_msg.lower() == "accepted"

    if success:
        print_success(f"{status_msg}")
    else:
        print_error(f"{status_msg}")

    runtime = result.get("runtime")
    memory = result.get("memory")
    runtime_percentile = result.get("runtime_percentile")
    memory_percentile = result.get("memory_percentile")

    if runtime:
        msg = f"Runtime: {runtime}"
        if runtime_percentile is not None:
            msg += f" (beats {runtime_percentile:.2f}%)"
        print_info(msg)

    if memory:
        msg = f"Memory: {memory}"
        if memory_percentile is not None:
            msg += f" (beats {memory_percentile:.2f}%)"
        print_info(msg)

    total_correct = result.get("total_correct")
    total_testcases = result.get("total_testcases")
    if total_testcases is not None:
        print_info(f"Passed: {total_correct}/{total_testcases}")

    last_testcase = result.get("last_testcase")
    expected_output = result.get("expected_output")
    code_output = result.get("code_output")

    if last_testcase and not success:
        console.print("\n[bold cyan]--- Last Testcase ---[/bold cyan]\n")
        console.print(last_testcase)
        console.print("\n[bold cyan]--- Expected Output ---[/bold cyan]\n")
        console.print(expected_output or "")
        console.print("\n[bold cyan]--- Your Output ---[/bold cyan]\n")
        console.print(code_output or "")


def main(config: Dict[str, Any]):
    """Submit solution to LeetCode and show benchmark stats."""
    print_command_banner("LeetCode Submit")

    problem_num = styled_text_input("Enter problem number to submit")
    base_dir = config.get("base_dir", "")

    if not base_dir or not os.path.isdir(base_dir):
        print_error(f"Invalid base directory in config: {base_dir}")
        return

    print_info("Searching for problem...")
    config_manager = ConfigManager()
    normalized_languages = config_manager.get_languages()
    extensions = [info["ext"] for info in normalized_languages.values()]
    all_files = get_all_solution_files(base_dir, extensions)

    matches = []
    for file_path in all_files:
        basename = os.path.basename(file_path)
        if _matches_problem_number(basename, problem_num):
            matches.append(file_path)

    if not matches:
        print_error(f"Could not find local file for problem {problem_num}.")
        return

    if len(matches) == 1:
        target_file = matches[0]
    else:
        choices = [os.path.basename(path) for path in matches]
        selected = styled_select("Multiple matches found. Select file to submit", choices)
        target_file = matches[choices.index(selected)]

    with open(target_file, "r", encoding="utf-8") as rf:
        content = rf.read()

    _, ext = os.path.splitext(target_file)
    language_key = get_language_by_extension(normalized_languages, ext)
    if not language_key:
        print_error("Unsupported file extension for submission.")
        return

    lang_slug = LEETCODE_LANG_MAP.get(language_key)
    if not lang_slug:
        print_error(f"Language '{language_key}' is not supported by LeetCode submit API.")
        return

    slug = _infer_slug_from_path(target_file, content)
    if not slug:
        slug = styled_text_input("Enter LeetCode problem slug (e.g., two-sum)")

    details = get_problem_details(slug)
    if not details:
        print_error("Could not fetch problem details from LeetCode.")
        return

    auth = get_leetcode_auth(config)
    if not auth["session"] or not auth["csrf"]:
        print_error("LeetCode session is missing. Set 'leetcode_session' and 'leetcode_csrf' in config.json or env vars.")
        return

    print_info("Submitting to LeetCode...")
    result = submit_solution(
        slug,
        details.get("questionId"),
        lang_slug,
        content,
        auth["session"],
        auth["csrf"],
    )

    if not result:
        print_error("LeetCode did not return a result. Please try again.")
        return

    _render_benchmark(result)
    print_success("Submit completed.")
