import os
import subprocess

from . import _api_exceptions


ormb_relpath = os.path.dirname(__file__)
ormb_abspath = os.path.abspath(ormb_relpath)
BIN_PATHNAME = os.path.join(ormb_abspath, "bin/ormb")


def login(
    hostname: str,
    username: str,
    password: str,
    insecure: str = "True",
) -> subprocess.CompletedProcess:
    """
    Runs ORMB CLI `login` command in a subprocess.

    Logs in an image registry.

    Args:
        hostname: Remote registry to authenticate to.
        username: Username to authenticate with.
        password: Password to authenticate with.
        insecure: Whether or not to allow connections to TLS registry without
            certificates. Either "True" or "False".

    Returns:
        The completed process.

    Raises:
        ORMBLoginError: ORMB login command process exited with a non-zero exit
            code.
    """
    args = [
        BIN_PATHNAME,
        "login",
        hostname,
        "--username",
        username,
        "--password-stdin",
        "--insecure",
        insecure,
    ]

    try:
        return subprocess.run(
            args=args,
            capture_output=True,
            check=True,
            input=password,
            text=True,
        )
    except subprocess.CalledProcessError as e:
        raise _api_exceptions.ORMBLoginError(
            hostname=hostname,
            username=username,
            stderr=e.stderr,
        )


def save(src: str, ref: str) -> subprocess.CompletedProcess:
    """
    Runs ORMB CLI `save` command in a subprocess.

    Packages an artifact and saves it to cache in the local file system.

    Args:
        src: ORMB formatted artifact dirname (containing the `model/`
            directory and `ormbfile.yaml`).
        ref: Image reference of the artifact to save, formatted as
            `<registry_address>/<repository>:<tag>`.

    Returns:
        The completed process.

    Raises:
        ORMBSaveError: ORMB save command process exited with a non-zero exit
            code.
    """
    try:
        return subprocess.run(
            args=[BIN_PATHNAME, "save", src, ref],
            capture_output=True,
            check=True,
            text=True,
        )
    except subprocess.CalledProcessError as e:
        raise _api_exceptions.ORMBSaveError(
            src=src,
            reference=ref,
            stderr=e.stderr,
        )


def push(ref: str) -> subprocess.CompletedProcess:
    """
    Runs ORMB CLI `push` command in a subprocess.

    Pushes an artifact image to a remote registry.

    Args:
        ref: Image reference of the artifact to push, formatted as
            `<registry_address>/<repository>:<tag>`.

    Returns:
        The completed process.

    Raises:
        ORMBPushError: ORMB push command process exited with a non-zero exit
            code.
    """
    try:
        return subprocess.run(
            args=[BIN_PATHNAME, "push", ref],
            capture_output=True,
            check=True,
            text=True,
        )
    except subprocess.CalledProcessError as e:
        raise _api_exceptions.ORMBPushError(reference=ref, stderr=e.stderr)


def pull(ref: str) -> subprocess.CompletedProcess:
    """
    Runs ORMB CLI `pull` command in a subprocess.

    Pulls an image from a remote registry to the cache of the local file
    system.

    Args:
        ref: Image reference of the artifact to pull, formatted as
            `<registry_address>/<repository>:<tag>`.

    Returns:
        The completed process.

    Raises:
        ORMBPullError: ORMB pull command process exited with a non-zero exit
            code.
    """
    try:
        return subprocess.run(
            args=[BIN_PATHNAME, "pull", ref],
            capture_output=True,
            check=True,
            text=True,
        )
    except subprocess.CalledProcessError as e:
        raise _api_exceptions.ORMBPullError(reference=ref, stderr=e.stderr)


def export(ref: str, dst: str) -> subprocess.CompletedProcess:
    """
    Runs ORMB CLI `export` command in a subprocess.

    Exports the artifact from the local cache to the destination directory.

    Args:
        ref: Image reference of the artifact to export, formatted as
            `<registry_address>/<repository>:<tag>`.
        dst: Destination directory path to export the `model/` directory
            (containing the artifact) and `ormbfile.yaml` to.

    Returns:
        The completed process.

    Raises:
        ORMBExportError: ORMB export command process exited with a non-zero
            exit code.
    """
    try:
        return subprocess.run(
            args=[BIN_PATHNAME, "export", ref, "-d", dst],
            capture_output=True,
            check=True,
            text=True,
        )
    except subprocess.CalledProcessError as e:
        raise _api_exceptions.ORMBExportError(
            reference=ref,
            destination=dst,
            stderr=e.stderr,
        )


def remove(ref: str) -> subprocess.CompletedProcess:
    """
    Runs ORMB CLI `remove` command in a subprocess.

    Args:
        ref: Image reference of the artifact to export, formatted as
            `<registry_address>/<repository>:<tag>`.

    Returns:
        The completed process.

    Raises:
        ORMBRemoveError: ORMB remove command process exited with a non-zero
            exit code.
    """
    try:
        return subprocess.run(
            args=[BIN_PATHNAME, "remove", ref],
            capture_output=True,
            check=True,
            text=True,
        )
    except subprocess.CalledProcessError as e:
        raise _api_exceptions.ORMBRemoveError(
            reference=ref,
            stderr=e.stderr,
        )
