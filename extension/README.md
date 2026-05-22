# LeetCode CLI Exporter (Browser Extension)

This folder contains a Chrome/Edge MV3 extension that saves LeetCode problems to GitHub using the same folder layout as the CLI.

## Install (unpacked)

1. Open Chrome or Edge.
2. Go to chrome://extensions.
3. Enable Developer mode.
4. Click "Load unpacked" and select this folder.

## Configure

Open the extension Settings and fill:
- GitHub token (repo scope).
- Repo owner/name.
- Branch and base dir (optional).
- Data structures and languages JSON if you want to match your CLI config.

## Use

1. Open a LeetCode problem page.
2. Click the extension icon.
3. Pick data structure and language.
4. Click "Save to GitHub".

It creates the same structure as the CLI:

base_dir/<structure>/<id>-<slug>/<id>_<title>.<ext>

If "Include README" is enabled, it also creates README.md with the problem description.
