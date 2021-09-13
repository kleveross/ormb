from subprocess import Popen
from os.path import abspath, join, dirname

BIN_PATH = join(abspath(dirname(__file__)), 'bin')


def login(hostname: str, username: str, password: str, insecure_opt: str):
    ex = Popen([join(BIN_PATH, "ormb"), "login", hostname, "--username", username, "--password", password, "--insecure",
                insecure_opt])
    status = ex.wait()
    return status


def push(ref: str):
    ex = Popen([join(BIN_PATH, "ormb"), "push", ref])
    status = ex.wait()
    return status


def pull(ref: str):
    ex = Popen([join(BIN_PATH, "ormb"), "pull", ref])
    status = ex.wait()
    return status


def export(ref: str, dst: str):
    ex = Popen([join(BIN_PATH, "ormb"), "export", ref, "-d", dst])
    status = ex.wait()
    return status


def save(src: str, ref: str):
    ex = Popen([join(BIN_PATH, "ormb"), "save", src, ref])
    status = ex.wait()
    return status


def remove(ref: str):
    ex = Popen([join(BIN_PATH, "ormb"), "remove", ref])
    status = ex.wait()
    return status
