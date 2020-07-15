import json
import requests
import platform
import sys
import tarfile
import os


REPOS = "kleveross/ormb"
VERSION = "latest"
BIN_PATH = './ormb/bin'

OS_LIST = ["Linux", "Darwin"]
ARCH_LIST = ["x86_64", "i386"]


def untar(fname, dirs):
    t = tarfile.open(fname)
    t.extractall(path=dirs)


def download():
    url = 'https://api.github.com/repos/%s/releases/%s' % (REPOS, VERSION)
    r = requests.get(url)

    if r.status_code != 200:
        raise Exception("get assets info err, ret code: %s" % r.status_code)

    json_info = json.loads(r.text)

    cur_version = json_info["tag_name"][1:]

    asset_name = "ormb_%s_Linux_x86_64.tar.gz" % cur_version
    for os_name in OS_LIST:
        for arch_name in ARCH_LIST:
            if arch_name in platform.platform().lower() and os_name.lower() in sys.platform:
                asset_name = "ormb_%s_%s_%s.tar.gz" % (cur_version, os_name, arch_name)

    asset_url = ""

    for asset in json_info["assets"]:
        if asset_name in asset["browser_download_url"]:
            asset_url = asset["url"]

    # download the url contents in binary format
    headers = {'Accept': 'application/octet-stream'}
    r = requests.get(asset_url, headers=headers)

    # open method to open a file on your system and write the contents
    with open(asset_name, "wb") as code:
        code.write(r.content)

    if not os.path.exists(BIN_PATH):
        os.mkdir(BIN_PATH)
    untar(asset_name, BIN_PATH)

    os.remove(asset_name)

