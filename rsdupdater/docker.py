import subprocess
from pathlib import Path


def login(user, password):
    try:
        subprocess.run(args=f"docker login -u={user} --password-stdin".split(' '), input=password.encode('utf-8'))
    except subprocess.CalledProcessError:
        print("docker login failed")
        return False
    return True


def logout():
    try:
        subprocess.run(args="docker logout".split(' '))
    except subprocess.CalledProcessError:
        print("docker logout failed")
        return False
    return True


def run_cmd(cmd):
    print(f"run: {' '.join(cmd)}")
    try:
        process = subprocess.run(args=cmd, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, check=True)
    except subprocess.CalledProcessError as ex:
        print("FAILED")
        print(ex.stdout.decode('utf-8'))
        return False
    print(process.stdout.decode("utf-8"))
    return True


def build(work_dir: Path, dockerfile: str, branch: str, re_commit: str, arch: str):
    repo = "iceflower/redeclipse-server"

    success = run_cmd(
        ["docker", "build", "--build-arg", f"BRANCH={branch}", "--build-arg", f"RECOMMIT={re_commit}", "-t",
         f"{repo}:{arch}-{branch}", "-f", f"{dockerfile}", f"{str(work_dir)}"])
    if not success:
        return False
    success = run_cmd(["docker", "push", f"{repo}:{arch}-{branch}"])
    if not success:
        return False
    success = run_cmd(
        ["docker", "manifest", "create", f"{repo}:{branch}", f"{repo}:amd64-{branch}", f"{repo}:arm64-{branch}"])
    if not success:
        print("Manifest create failed.")
        #return False
    success = run_cmd(["docker", "manifest", "annotate", f"{repo}:{branch}", f"{repo}:arm64-{branch}", "--variant", "v8"])
    if not success:
        print("Manifest annotate failed.")
        # return False
    success = run_cmd(["docker",  "manifest",  "push", "--purge", f"{repo}:{branch}"])
    if not success:
        print("Manifest push failed.")
        # return False
    return True
