import os
import questionary
from rich import print

def main(config):
    base_dir = config["base_dir"]
    
    problem_num = questionary.text("Problem number:").ask()
    problem_name = questionary.text("Problem name:").ask()
    data_structure = questionary.select(
        "Data structure:",
        choices=list(config["data_structures"].keys())
    ).ask()
    
    ds_folder = config["data_structures"][data_structure]
    problem_dir = os.path.join(base_dir, ds_folder, problem_num)
    
    if os.path.exists(problem_dir):
        print("[red]Problem directory already exists![/red]")
        return
    
    os.makedirs(problem_dir, exist_ok=True)
    
    problem_file = os.path.join(problem_dir, f"{problem_num}_{problem_name}.cpp")
    
    with open(problem_file, "w") as f:
        f.write(f"""/*
LeetCode Problem {problem_num}: {problem_name}
Data Structure: {data_structure}
*/

// Solution 1
/*
Method: 
Time Complexity: 
Space Complexity: 
*/
""")
    
    print(f"\n[green]✔ Created problem directory and file[/green]\n")
    print(f"[cyan]Path: {problem_file}[/cyan]\n")