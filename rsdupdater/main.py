import sys
from argparse import ArgumentParser
from pathlib import Path

import rsdupdater.config as config
import rsdupdater.docker as docker
import rsdupdater.updater as updater
import rsdupdater.utils as utils
from rsdupdater.hash_store import HashStore


def main(argv=None):
    if argv is None:
        argv = sys.argv[1:]

    parser = ArgumentParser(prog="RSD-Updater", description="Updates the Red Eclipse server docker images.")
    parser.add_argument('-v', '--version', action='version', version="RSD-Updater 2.0.0")
    parser.add_argument('work_dir', type=str, metavar='<execution folder>',
                        help='working directory')
    parser.add_argument('--dry', dest='dry_run', action='store_true', default=False,
                        help='dry run, just shows which images are outdated')

    root_group = parser.add_mutually_exclusive_group()
    root_group.add_argument('--config', dest='config', action='store_true', help='use of config file')

    single_group = root_group.add_argument_group()
    single_group.add_argument('--dockerfile', dest='dockerfile', type=str, metavar='<dockerfile name>',
                              help='name of the dockerfile in the working directory, optional, default to "Dockerfile_<branch>"')
    single_group.add_argument('--branch', dest='branch', type=str, metavar='<branch>',
                              help='red eclipse branch')
    single_group.add_argument('--arch', dest='arch', type=str, metavar='<architecture>',
                              help='cpu architecture')
    single_group.add_argument('--user', dest='user', type=str, metavar='<docker user>',
                              help='docker user')
    single_group.add_argument('--password', dest='password', type=str, metavar='<docker password>',
                              help='docker password')

    args = parser.parse_args(argv)

    # get all arguments and check paths
    dry_run = args.dry_run
    working_dir = Path(args.work_dir)
    print(f"Use '{str(working_dir.absolute())}' as working directory")
    config_file = None
    hash_storage_file = working_dir.joinpath('hash-storage.json')
    if args.config:
        config_file = working_dir.joinpath('config.json')
    if not utils.check_paths(working_dir, hash_storage_file, config_file):
        exit(1)

    # load or create config
    if args.config:
        conf = config.load_config(config_file)
    else:
        conf = config.create_config(args.user, args.password, args.dockerfile, args.branch, args.arch)
    hash_storage = HashStore.from_file(hash_storage_file)

    for build in conf['build']:
        print(f"Update step: {build}")
        # get current hash
        cur_hash = hash_storage.get(build['branch'], build['arch'])
        if cur_hash is None:
            hash_storage.create(build['branch'], build['arch'])
            cur_hash = hash_storage.get(build['branch'], build['arch'])

        # get new hashs, stop on error
        new_hash = updater.get_hashs(working_dir.joinpath(build['dockerfile']), build['branch'], build['arch'])
        if new_hash is None:
            continue

        print("Current: ", cur_hash)
        print("New:     ", new_hash)

        # check for updates
        if not updater.need_update(cur_hash, new_hash):
            print("No update required.")
            continue
        print("Update required.")

        # stop on dry run
        if dry_run:
            print("Dry run, stop here.")
            continue

        # update image
        if not docker.login(conf['docker']['user'], conf['docker']['password']):
            exit(2)
        try:
            success = docker.build(working_dir, build['dockerfile'], build['branch'], new_hash['re-commit'], build['arch'])
            if not success:
                continue
        except Exception:
            docker.logout()
            exit(2)
        finally:
            docker.logout()
        # update hash file
        hash_storage.update(new_hash, build['branch'], build['arch'])
        print("Save new hash values")
        hash_storage.to_file(hash_storage_file)
