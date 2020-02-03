import json
from pathlib import Path


class HashStore:
    def __init__(self, store):
        self.store = store

    @classmethod
    def from_file(cls, file: Path):
        if not file.exists():
            return cls({})
        with file.open('r') as reader:
            j = json.loads(reader.read())
        return cls(j)

    def to_file(self, file: Path):
        with file.open('w') as writer:
            writer.write(json.dumps(self.store, indent=4))

    def create(self, branch: str, arch: str, os: str = 'linux'):
        self.store[branch] = {arch: {os: {'alpine': "", 'dockerfile': "", 're-commit': ""}}}

    def get(self, branch: str, arch: str, os: str = 'linux'):
        if branch in self.store and arch in self.store[branch] and os in self.store[branch][arch]:
            return self.store[branch][arch][os]
        return None

    def update(self, value: dict, branch: str, arch: str, os: str = 'linux'):
        self.store[branch][arch][os] = value
