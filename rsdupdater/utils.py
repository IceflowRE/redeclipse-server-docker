from pathlib import Path


def update_image():
    pass


def check_paths(work_dir: Path, hash_store: Path, config_file: Path):
    if not hash_store.exists():
        print(f"hash storage '{hash_store}' does not exist, will use empty one")
    if not work_dir.exists():
        print(f"working directory '{work_dir}' does not exist")
        return False
    if config_file is not None and not config_file.exists():
        print(f"config file '{config_file}' does not exist.")
        return False
    return True
