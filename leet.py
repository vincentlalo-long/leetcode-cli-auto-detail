import json
import os
import sys

from add import main as add_main
from add_sol import main as add_sol_main

CONFIG_PATH = os.path.join(os.path.dirname(__file__), "config.json")

def load_config():
    with open(CONFIG_PATH) as f:
        return json.load(f)

def main():
    config = load_config()

    if len(sys.argv) < 2:
        print("Usage: leet [add | add-sol]")
        return

    cmd = sys.argv[1]

    if cmd == "add":
        add_main(config)
    elif cmd == "add-sol":
        add_sol_main(config)
    else:
        print("Unknown command")

if __name__ == "__main__":
    main()