import json
from pathlib import Path


def load_config(file: Path):
    with file.open('r') as reader:
        j = json.loads(reader.read())
    for build in j['build']:
        if 'dockerfile' not in build:
            build['dockerfile'] = f"Dockerfile_{build['branch']}"
    return j


def create_config(docker_user: str, docker_password: str, dockerfile: str, branch: str, arch: str):
    return {
        "docker": {
            "user": docker_user,
            "password": docker_password
        },
        "build": [
            {
                "branch": branch,
                "arch": arch,
                "os": "linux",
                "dockerfile": f"Dockerfile_{branch}" if dockerfile is "" else dockerfile
            }
        ]
    }
