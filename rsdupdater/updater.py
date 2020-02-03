import hashlib
import json
import subprocess
from http import client
from pathlib import Path


def get_file_hash(file: Path):
    if not file.exists():
        return None
    h = hashlib.sha256()
    with file.open('rb') as reader:
        data = reader.read()
    h.update(data)
    return 'sha256:' + h.hexdigest()


def get_alpine_hash(arch: str, os: str = 'linux'):
    conn = client.HTTPSConnection("registry.hub.docker.com")
    conn.request('GET', "/v2/repositories/library/alpine/tags/latest")
    req = conn.getresponse()
    if req.status != 200 and req.status != 404:
        return None, "Docker request failed."
    j = json.loads(req.read())
    if req.status == 404:
        return None, j['message']
    for image in j['images']:
        if image['architecture'] == arch and image['os'] == os:
            return image['digest'], None
    return None, f"base image for {arch}, {os} not found"


def get_commit_hash(branch: str):
    try:
        process = subprocess.run(
            args=f"git ls-remote https://github.com/redeclipse/base.git refs/heads/{branch}".split(' '),
            capture_output=True
        )
    except subprocess.CalledProcessError:
        return None
    return 'sha1:' + process.stdout.decode('utf-8').split('\t')[0]


def get_hashs(dockerfile: Path, branch: str, arch: str, os: str = 'linux'):
    re_commit = get_commit_hash(branch)
    if re_commit is None:
        print("failed to get newest git commit hash")
        return None
    docker_hash = get_file_hash(dockerfile)
    if docker_hash is None:
        print(f"dockerfile hash failed {dockerfile}")
        return None
    alpine_hash, err = get_alpine_hash(arch, os)
    if err is not None:
        print("alpine hash failed: ", err)
        return None
    return {'alpine': alpine_hash, 'dockerfile': docker_hash, 're-commit': re_commit}


def need_update(cur_hash: dict, new_hash: dict) -> bool:
    return cur_hash['re-commit'] != new_hash['re-commit'] or cur_hash['dockerfile'] != new_hash['dockerfile'] or \
           cur_hash['alpine'] != new_hash['alpine']
