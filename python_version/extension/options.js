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

const dom = {
  githubToken: document.getElementById("githubToken"),
  repoOwner: document.getElementById("repoOwner"),
  repoName: document.getElementById("repoName"),
  branch: document.getElementById("branch"),
  baseDir: document.getElementById("baseDir"),
  dataStructures: document.getElementById("dataStructures"),
  languages: document.getElementById("languages"),
  defaultLanguage: document.getElementById("defaultLanguage"),
  save: document.getElementById("save"),
  message: document.getElementById("message")
};

function setMessage(text, type = "") {
  dom.message.textContent = text;
  dom.message.className = "muted" + (type ? ` ${type}` : "");
}

function safeParseJson(text, fallback) {
  try {
    return text.trim() ? JSON.parse(text) : fallback;
  } catch {
    return null;
  }
}

async function loadSettings() {
  const settings = await chrome.storage.sync.get(null);
  dom.githubToken.value = settings.githubToken || "";
  dom.repoOwner.value = settings.repoOwner || "";
  dom.repoName.value = settings.repoName || "";
  dom.branch.value = settings.branch || "main";
  dom.baseDir.value = settings.baseDir || "";
  dom.dataStructures.value = JSON.stringify(settings.dataStructures || DEFAULT_STRUCTURES, null, 2);
  dom.languages.value = JSON.stringify(settings.languages || DEFAULT_LANGUAGES, null, 2);
  dom.defaultLanguage.value = settings.defaultLanguage || "cpp";
}

async function saveSettings() {
  setMessage("");
  const dataStructures = safeParseJson(dom.dataStructures.value, DEFAULT_STRUCTURES);
  const languages = safeParseJson(dom.languages.value, DEFAULT_LANGUAGES);

  if (!dataStructures || !languages) {
    setMessage("Invalid JSON for data structures or languages", "error");
    return;
  }

  await chrome.storage.sync.set({
    githubToken: dom.githubToken.value.trim(),
    repoOwner: dom.repoOwner.value.trim(),
    repoName: dom.repoName.value.trim(),
    branch: dom.branch.value.trim() || "main",
    baseDir: dom.baseDir.value.trim(),
    dataStructures,
    languages,
    defaultLanguage: dom.defaultLanguage.value.trim() || "cpp"
  });

  setMessage("Settings saved", "success");
}

dom.save.addEventListener("click", saveSettings);
loadSettings();
