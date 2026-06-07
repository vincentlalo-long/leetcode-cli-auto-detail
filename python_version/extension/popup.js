const DEFAULT_LANGUAGES = {
  cpp: { label: "C++", ext: "cpp" },
  c: { label: "C", ext: "c" },
  python: { label: "Python", ext: "py" },
  go: { label: "Go", ext: "go" }
};

const DEFAULT_STRUCTURES = {
  array: "array",
  string: "string",
  linked_list: "linked_list",
  stack: "stack",
  graph: "graph",
  queue: "queue"
};

const STORAGE_KEYS = [
  "githubToken",
  "repoOwner",
  "repoName",
  "branch",
  "baseDir",
  "languages",
  "dataStructures",
  "defaultLanguage"
];

const STATUS = {
  idle: "Open a LeetCode problem page.",
  loading: "Loading problem details...",
  ready: "Ready to save.",
  error: "Could not load problem data."
};

const dom = {
  status: document.getElementById("status"),
  title: document.getElementById("problem-title"),
  id: document.getElementById("problem-id"),
  difficulty: document.getElementById("problem-difficulty"),
  tags: document.getElementById("problem-tags"),
  structure: document.getElementById("structure"),
  language: document.getElementById("language"),
  includeReadme: document.getElementById("include-readme"),
  save: document.getElementById("save"),
  message: document.getElementById("message")
};

let problemDetails = null;

function setStatus(text) {
  dom.status.textContent = text;
}

function setMessage(text, type = "") {
  dom.message.textContent = text;
  dom.message.className = "muted" + (type ? ` ${type}` : "");
}

function slugify(text) {
  return text
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/(^-|-$)/g, "");
}

function sanitizeFileName(text) {
  return text.replace(/[\\/:*?"<>|]/g, "-").trim();
}

function normalizeLanguages(languages) {
  if (!languages || Object.keys(languages).length === 0) {
    return { ...DEFAULT_LANGUAGES };
  }
  const normalized = {};
  Object.entries(languages).forEach(([key, value]) => {
    if (typeof value === "string") {
      normalized[key] = { label: key.toUpperCase(), ext: value.replace(/^\./, "") };
      return;
    }
    const label = value.label || key.toUpperCase();
    const ext = value.ext || value.extension;
    if (!ext) {
      return;
    }
    normalized[key] = { label, ext: String(ext).replace(/^\./, "") };
  });
  return Object.keys(normalized).length ? normalized : { ...DEFAULT_LANGUAGES };
}

function buildProblemHeader({ languageKey, problemNum, title, link, difficulty, tags, structure }) {
  const lines = [
    `LeetCode Problem ${problemNum}: ${title}`,
    `Link: ${link}`,
    `Difficulty: ${difficulty}`,
    `Tags: ${tags}`,
    `Data Structure: ${structure}`
  ];

  if (languageKey === "python") {
    return `"""\n${lines.join("\n")}\n"""\n\n`;
  }

  return `/*\n${lines.join("\n")}\n*/\n\n`;
}

function buildProblemTemplate({ languageKey, problemNum, title, link, difficulty, tags, structure }) {
  const header = buildProblemHeader({
    languageKey,
    problemNum,
    title,
    link,
    difficulty,
    tags,
    structure
  });

  if (languageKey === "python") {
    return (
      header +
      "from typing import List\n\n" +
      "class Solution:\n" +
      "    pass\n\n" +
      "if __name__ == \"__main__\":\n" +
      "    print(\"Test cases go here!\")\n"
    );
  }

  if (languageKey === "go") {
    return (
      header +
      "package main\n\n" +
      "import \"fmt\"\n\n" +
      "func main() {\n" +
      "    fmt.Println(\"Test cases go here!\")\n" +
      "}\n"
    );
  }

  if (languageKey === "c") {
    return (
      header +
      "#include <stdio.h>\n\n" +
      "int main() {\n" +
      "    printf(\"Test cases go here!\\n\");\n" +
      "    return 0;\n" +
      "}\n"
    );
  }

  return (
    header +
    "#include <iostream>\n" +
    "#include <vector>\n" +
    "#include <string>\n\n" +
    "using namespace std;\n\n" +
    "// class Solution {\n" +
    "// public:\n" +
    "//     \n" +
    "// };\n\n" +
    "int main() {\n" +
    "    // Solution sol;\n" +
    "    cout << \"Test cases go here!\" << endl;\n" +
    "    return 0;\n" +
    "}\n"
  );
}

function buildReadme({ problemNum, title, link, difficulty, tags, content }) {
  return (
    `# [${problemNum}. ${title}](${link})\n\n` +
    `- **Difficulty:** ${difficulty}\n` +
    `- **Tags:** ${tags}\n\n` +
    "## Description\n\n" +
    (content || "")
  );
}

async function getActiveTab() {
  const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
  return tab;
}

function getSlugFromUrl(url) {
  try {
    const parsed = new URL(url);
    const match = parsed.pathname.match(/\/problems\/([^/]+)\/?/);
    if (!match) {
      return null;
    }
    return match[1];
  } catch {
    return null;
  }
}

async function fetchProblemDetails(slug) {
  const query = `
    query getQuestionDetail($titleSlug: String!) {
      question(titleSlug: $titleSlug) {
        questionId
        questionFrontendId
        title
        titleSlug
        difficulty
        content
        topicTags {
          name
        }
      }
    }
  `;

  const response = await fetch("https://leetcode.com/graphql", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Referer": "https://leetcode.com"
    },
    body: JSON.stringify({ query, variables: { titleSlug: slug } })
  });

  if (!response.ok) {
    throw new Error("GraphQL request failed");
  }

  const data = await response.json();
  if (data.errors || !data.data || !data.data.question) {
    throw new Error("Problem not found");
  }

  return data.data.question;
}

function populateSelect(select, items, selectedKey) {
  select.innerHTML = "";
  items.forEach((item) => {
    const option = document.createElement("option");
    option.value = item.value;
    option.textContent = item.label;
    if (item.value === selectedKey) {
      option.selected = true;
    }
    select.appendChild(option);
  });
}

function getDefaultSettings(settings) {
  return {
    githubToken: settings.githubToken || "",
    repoOwner: settings.repoOwner || "",
    repoName: settings.repoName || "",
    branch: settings.branch || "main",
    baseDir: settings.baseDir || "",
    languages: normalizeLanguages(settings.languages || DEFAULT_LANGUAGES),
    dataStructures: settings.dataStructures || DEFAULT_STRUCTURES,
    defaultLanguage: settings.defaultLanguage || "cpp"
  };
}

function validateSettings(settings) {
  const missing = [];
  if (!settings.githubToken) missing.push("token");
  if (!settings.repoOwner) missing.push("owner");
  if (!settings.repoName) missing.push("repo");
  return missing;
}

function buildRepoPath(baseDir, ...parts) {
  const cleanParts = [baseDir, ...parts]
    .filter(Boolean)
    .map((part) => part.replace(/^\/+|\/+$/g, ""));
  return cleanParts.join("/");
}

async function checkFileExists({ owner, repo, path, token, branch }) {
  const url = `https://api.github.com/repos/${owner}/${repo}/contents/${path}?ref=${encodeURIComponent(branch)}`;
  const response = await fetch(url, {
    headers: {
      Authorization: `token ${token}`,
      Accept: "application/vnd.github+json"
    }
  });

  if (response.status === 404) {
    return null;
  }

  if (!response.ok) {
    throw new Error(`GitHub check failed: ${response.status}`);
  }

  return response.json();
}

async function createOrUpdateFile({ owner, repo, path, message, content, token, branch, sha }) {
  const url = `https://api.github.com/repos/${owner}/${repo}/contents/${path}`;
  const encoded = btoa(unescape(encodeURIComponent(content)));
  const body = {
    message,
    content: encoded,
    branch
  };
  if (sha) {
    body.sha = sha;
  }

  const response = await fetch(url, {
    method: "PUT",
    headers: {
      Authorization: `token ${token}`,
      Accept: "application/vnd.github+json"
    },
    body: JSON.stringify(body)
  });

  if (!response.ok) {
    const payload = await response.json().catch(() => ({}));
    const messageText = payload && payload.message ? payload.message : "Unknown error";
    throw new Error(`GitHub write failed: ${messageText}`);
  }

  return response.json();
}

async function saveToGitHub({ settings, problem, structureKey, languageKey, includeReadme }) {
  const languages = normalizeLanguages(settings.languages);
  const language = languages[languageKey] || Object.values(languages)[0];
  const ext = language.ext;

  const title = sanitizeFileName(problem.title || "Unknown");
  const problemNum = problem.questionFrontendId || problem.questionId || "0";
  const link = `https://leetcode.com/problems/${problem.titleSlug || slugify(problem.title || "")}/`;
  const tags = (problem.topicTags || []).map((tag) => tag.name).join(", ") || "None";
  const structureName = structureKey === "[Uncategorized]" ? "Uncategorized" : structureKey;
  const dsFolder = structureKey === "[Uncategorized]"
    ? "uncategorized"
    : settings.dataStructures[structureKey];

  const folderSlug = problem.titleSlug || slugify(problem.title || "problem");
  const folderName = `${problemNum}-${folderSlug}`;
  const fileName = `${problemNum}_${title}.${ext}`;

  const problemPath = buildRepoPath(settings.baseDir, dsFolder, folderName, fileName);
  const readmePath = buildRepoPath(settings.baseDir, dsFolder, folderName, "README.md");

  const problemContent = buildProblemTemplate({
    languageKey,
    problemNum,
    title: problem.title,
    link,
    difficulty: problem.difficulty || "Unknown",
    tags,
    structure: structureName
  });

  const readmeContent = buildReadme({
    problemNum,
    title: problem.title,
    link,
    difficulty: problem.difficulty || "Unknown",
    tags,
    content: problem.content || ""
  });

  const messagePrefix = `Add ${problemNum} ${problem.title}`;

  const problemExists = await checkFileExists({
    owner: settings.repoOwner,
    repo: settings.repoName,
    path: problemPath,
    token: settings.githubToken,
    branch: settings.branch
  });

  if (problemExists) {
    throw new Error("Problem file already exists in repo");
  }

  await createOrUpdateFile({
    owner: settings.repoOwner,
    repo: settings.repoName,
    path: problemPath,
    message: `${messagePrefix} (${ext})`,
    content: problemContent,
    token: settings.githubToken,
    branch: settings.branch
  });

  if (includeReadme) {
    const readmeExists = await checkFileExists({
      owner: settings.repoOwner,
      repo: settings.repoName,
      path: readmePath,
      token: settings.githubToken,
      branch: settings.branch
    });

    if (!readmeExists) {
      await createOrUpdateFile({
        owner: settings.repoOwner,
        repo: settings.repoName,
        path: readmePath,
        message: `${messagePrefix} (README)`,
        content: readmeContent,
        token: settings.githubToken,
        branch: settings.branch
      });
    }
  }

  return { problemPath, readmePath: includeReadme ? readmePath : null };
}

async function init() {
  const settings = await chrome.storage.sync.get(STORAGE_KEYS);
  const config = getDefaultSettings(settings);

  const structures = [
    { value: "[Uncategorized]", label: "[Uncategorized]" },
    ...Object.keys(config.dataStructures).map((key) => ({ value: key, label: key }))
  ];

  const languages = normalizeLanguages(config.languages);
  const languageItems = Object.entries(languages).map(([key, info]) => ({
    value: key,
    label: `${info.label} (${info.ext})`
  }));

  populateSelect(dom.structure, structures, "[Uncategorized]");
  populateSelect(dom.language, languageItems, config.defaultLanguage);

  const tab = await getActiveTab();
  const slug = tab && tab.url ? getSlugFromUrl(tab.url) : null;

  if (!slug) {
    setStatus(STATUS.idle);
    dom.save.disabled = true;
    return;
  }

  setStatus(STATUS.loading);

  try {
    problemDetails = await fetchProblemDetails(slug);
    dom.title.textContent = problemDetails.title || "-";
    dom.id.textContent = problemDetails.questionFrontendId || "-";
    dom.difficulty.textContent = problemDetails.difficulty || "-";
    dom.tags.textContent = (problemDetails.topicTags || []).map((tag) => tag.name).join(", ") || "-";
    setStatus(STATUS.ready);
    dom.save.disabled = false;
  } catch (error) {
    setStatus(STATUS.error);
    dom.save.disabled = true;
  }

  dom.save.addEventListener("click", async () => {
    setMessage("");
    dom.save.disabled = true;
    const latestSettings = getDefaultSettings(await chrome.storage.sync.get(STORAGE_KEYS));
    const missing = validateSettings(latestSettings);
    if (missing.length) {
      setMessage("Missing settings: " + missing.join(", "), "error");
      dom.save.disabled = false;
      return;
    }

    try {
      const result = await saveToGitHub({
        settings: latestSettings,
        problem: problemDetails,
        structureKey: dom.structure.value,
        languageKey: dom.language.value,
        includeReadme: dom.includeReadme.checked
      });

      const saved = result.readmePath ? `Saved: ${result.problemPath} and README` : `Saved: ${result.problemPath}`;
      setMessage(saved, "success");
    } catch (error) {
      setMessage(error.message || "Save failed", "error");
    } finally {
      dom.save.disabled = false;
    }
  });
}

init();
