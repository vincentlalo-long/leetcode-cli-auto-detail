import os
import questionary
from rich import print

def choose_file(base_dir):
    files = []

    for root, _, filenames in os.walk(base_dir):
        for f in filenames:
            if f.endswith(".cpp"):
                files.append(os.path.join(root, f))

    return questionary.select(
        "Select file:",
        choices=files
    ).ask()

def count_solutions(content):
    return content.count("Solution")

def main(config):
    base_dir = config["base_dir"]

    file_path = choose_file(base_dir)

    if not file_path:
        print("[red]No file selected[/red]")
        return

    with open(file_path, "r") as f:
        content = f.read()

    sol_num = count_solutions(content) + 1

    print(f"\n[cyan]Adding Solution {sol_num}[/cyan]\n")

    method = questionary.text("Method:").ask()
    time = questionary.text("Time complexity:").ask()
    space = questionary.text("Space complexity:").ask()

    print("\n[yellow]Paste code (end with EOF):[/yellow]")

    lines = []
    while True:
        line = input()
        if line.strip() == "EOF":
            break
        lines.append(line)

    code = "\n".join(lines)

    new_block = f"""

/// ================== Solution {sol_num} ==================
/*
Method: {method}
Time Complexity: {time}
Space Complexity: {space}
*/

{code}
"""

    with open(file_path, "a") as f:
        f.write(new_block)

    print(f"\n[green]✔ Added Solution {sol_num}[/green]\n")