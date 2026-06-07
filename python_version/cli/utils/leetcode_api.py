import os
import requests
import re
import time
from typing import Optional, Dict, Any

def get_problem_details(title_slug: str) -> Optional[Dict[str, Any]]:
    """Fetch problem details from LeetCode GraphQL API."""
    url = "https://leetcode.com/graphql"
    query = """
    query getQuestionDetail($titleSlug: String!) {
      question(titleSlug: $titleSlug) {
        questionId
        questionFrontendId
        title
        difficulty
        content
        hints
        similarQuestions
        topicTags {
          name
        }
      }
    }
    """
    variables = {"titleSlug": title_slug}
    headers = {
        "Content-Type": "application/json",
        "Referer": "https://leetcode.com"
    }
    payload = {"query": query, "variables": variables}
    
    try:
        response = requests.post(url, json=payload, headers=headers, timeout=10)
        response.raise_for_status()
        data = response.json()
        if "errors" in data or not data.get("data", {}).get("question"):
            return None
        return data["data"]["question"]
    except requests.exceptions.RequestException:
        return None

def get_daily_challenge() -> Optional[Dict[str, Any]]:
    """Fetch the current daily challenge from LeetCode."""
    url = "https://leetcode.com/graphql"
    query = """
    query questionOfToday {
      activeDailyCodingChallengeQuestion {
        date
        userStatus
        link
        question {
          questionId
          questionFrontendId
          title
          titleSlug
          difficulty
          topicTags {
            name
          }
        }
      }
    }
    """
    headers = {
        "Content-Type": "application/json",
        "Referer": "https://leetcode.com"
    }
    payload = {"query": query}
    
    try:
        response = requests.post(url, json=payload, headers=headers, timeout=10)
        response.raise_for_status()
        data = response.json()
        if "errors" in data or not data.get("data", {}).get("activeDailyCodingChallengeQuestion"):
            return None
        return data["data"]["activeDailyCodingChallengeQuestion"]["question"]
    except requests.exceptions.RequestException:
        return None

def slugify(text: str) -> str:
    """Convert text to a URL-friendly slug."""
    text = text.lower()
    text = re.sub(r'[^a-z0-9]+', '-', text)
    return text.strip('-')

def get_problem_by_id(frontend_id: str) -> Optional[Dict[str, str]]:
    """Fetch problem title and slug by frontend ID."""
    url = "https://leetcode.com/api/problems/all/"
    try:
        response = requests.get(url, timeout=10)
        response.raise_for_status()
        data = response.json()
        
        # Normalize the input ID (e.g., '0001' -> '1')
        clean_id = str(int(frontend_id)) if frontend_id.isdigit() else frontend_id
        
        for p in data.get("stat_status_pairs", []):
            api_id = str(p.get("stat", {}).get("frontend_question_id"))
            clean_api_id = str(int(api_id)) if api_id.isdigit() else api_id
            
            if clean_api_id == clean_id:
                return {
                    "title": p["stat"]["question__title"],
                    "slug": p["stat"]["question__title_slug"]
                }
    except requests.exceptions.RequestException:
        pass
    return None

def get_all_problems() -> Optional[Dict[str, Any]]:
    """Fetch the list of all problems from LeetCode."""
    url = "https://leetcode.com/api/problems/all/"
    try:
        response = requests.get(url, timeout=10)
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException:
        return None

def get_user_profile(username: str) -> Optional[Dict[str, Any]]:
    """Fetch user profile stats from LeetCode GraphQL API."""
    url = "https://leetcode.com/graphql"
    query = """
    query getUserProfile($username: String!) {
      matchedUser(username: $username) {
        username
        submitStats: submitStatsGlobal {
          acSubmissionNum {
            difficulty
            count
          }
        }
        profile {
          ranking
          reputation
          starRating
        }
      }
    }
    """
    variables = {"username": username}
    headers = {
        "Content-Type": "application/json",
        "Referer": "https://leetcode.com"
    }
    payload = {"query": query, "variables": variables}
    
    try:
        response = requests.post(url, json=payload, headers=headers, timeout=10)
        response.raise_for_status()
        data = response.json()
        if "errors" in data or not data.get("data", {}).get("matchedUser"):
            return None
        return data["data"]["matchedUser"]
    except requests.exceptions.RequestException:
        return None

def get_upcoming_contests() -> Optional[list]:
    """Fetch upcoming contests from LeetCode GraphQL API."""
    url = "https://leetcode.com/graphql"
    query = """
    query {
      topTwoContests {
        title
        titleSlug
        startTime
        duration
      }
    }
    """
    headers = {
        "Content-Type": "application/json",
        "Referer": "https://leetcode.com"
    }
    payload = {"query": query}
    
    try:
        response = requests.post(url, json=payload, headers=headers, timeout=10)
        response.raise_for_status()
        data = response.json()
        if "errors" in data or not data.get("data", {}).get("topTwoContests"):
            return None
        return data["data"]["topTwoContests"]
    except requests.exceptions.RequestException:
        return None


def get_leetcode_auth(config: Optional[Dict[str, Any]] = None) -> Dict[str, str]:
    session = None
    csrf = None
    if config:
        session = config.get("leetcode_session")
        csrf = config.get("leetcode_csrf")

    session = session or os.environ.get("LEETCODE_SESSION")
    csrf = csrf or os.environ.get("LEETCODE_CSRF")

    return {"session": session or "", "csrf": csrf or ""}


def build_leetcode_headers(session: str, csrf: str, referer: str) -> Dict[str, str]:
    return {
        "Content-Type": "application/json",
        "Referer": referer,
        "x-csrftoken": csrf,
        "Cookie": f"LEETCODE_SESSION={session}; csrftoken={csrf}",
    }


def get_problem_editor_data(title_slug: str) -> Optional[Dict[str, Any]]:
    url = "https://leetcode.com/graphql"
    query = """
    query questionData($titleSlug: String!) {
      question(titleSlug: $titleSlug) {
        questionId
        questionFrontendId
        title
        titleSlug
        exampleTestcases
        sampleTestCase
      }
    }
    """
    variables = {"titleSlug": title_slug}
    headers = {
        "Content-Type": "application/json",
        "Referer": "https://leetcode.com"
    }
    payload = {"query": query, "variables": variables}

    try:
        response = requests.post(url, json=payload, headers=headers, timeout=10)
        response.raise_for_status()
        data = response.json()
        if "errors" in data or not data.get("data", {}).get("question"):
            return None
        return data["data"]["question"]
    except requests.exceptions.RequestException:
        return None


def get_problem_testcases(title_slug: str) -> Optional[str]:
    data = get_problem_editor_data(title_slug)
    if not data:
        return None

    example_cases = data.get("exampleTestcases")
    if isinstance(example_cases, list) and example_cases:
        return "\n".join(example_cases).strip()

    sample = data.get("sampleTestCase")
    if isinstance(sample, str) and sample.strip():
        return sample.strip()

    return None


def poll_submission_result(
    submission_id: str,
    headers: Dict[str, str],
    max_wait_seconds: float = 30.0,
    interval_seconds: float = 0.6,
) -> Optional[Dict[str, Any]]:
    url = f"https://leetcode.com/submissions/detail/{submission_id}/check/"
    start = time.time()
    while time.time() - start < max_wait_seconds:
        try:
            response = requests.get(url, headers=headers, timeout=10)
            response.raise_for_status()
            payload = response.json()
            state = payload.get("state")
            if state and state != "PENDING":
                return payload
        except requests.exceptions.RequestException:
            return None
        time.sleep(interval_seconds)
    return None


def interpret_solution(
    title_slug: str,
    question_id: str,
    lang_slug: str,
    code: str,
    testcases: str,
    session: str,
    csrf: str,
) -> Optional[Dict[str, Any]]:
    url = f"https://leetcode.com/problems/{title_slug}/interpret_solution/"
    headers = build_leetcode_headers(session, csrf, referer=f"https://leetcode.com/problems/{title_slug}/")
    payload = {
        "lang": lang_slug,
        "question_id": str(question_id),
        "typed_code": code,
        "data_input": testcases,
        "test_mode": False,
    }

    try:
        response = requests.post(url, json=payload, headers=headers, timeout=15)
        response.raise_for_status()
        data = response.json()
    except requests.exceptions.RequestException:
        return None

    interpret_id = data.get("interpret_id") or data.get("interpretId")
    if not interpret_id:
        return None
    return poll_submission_result(str(interpret_id), headers)


def submit_solution(
    title_slug: str,
    question_id: str,
    lang_slug: str,
    code: str,
    session: str,
    csrf: str,
) -> Optional[Dict[str, Any]]:
    url = f"https://leetcode.com/problems/{title_slug}/submit/"
    headers = build_leetcode_headers(session, csrf, referer=f"https://leetcode.com/problems/{title_slug}/")
    payload = {
        "lang": lang_slug,
        "question_id": str(question_id),
        "typed_code": code,
    }

    try:
        response = requests.post(url, json=payload, headers=headers, timeout=15)
        response.raise_for_status()
        data = response.json()
    except requests.exceptions.RequestException:
        return None

    submission_id = data.get("submission_id") or data.get("submissionId")
    if not submission_id:
        return None
    return poll_submission_result(str(submission_id), headers)
