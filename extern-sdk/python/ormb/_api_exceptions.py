class ORMBLoginError(Exception):
    def __init__(self, hostname: str, username: str, stderr: str) -> None:
        msg = (
            f"Login to {hostname} as '{username}' failed with stderr:\n"
            f"{stderr}"
        )
        super().__init__(msg)


class ORMBSaveError(Exception):
    def __init__(self, src: str, reference: str, stderr: str) -> None:
        msg = (
            f"Packaging and saving artifact from '{src}' to local cache with "
            f"reference '{reference}' failed with stderr:\n{stderr}"
        )
        super().__init__(msg)


class ORMBPushError(Exception):
    def __init__(self, reference: str, stderr: str) -> None:
        msg = (
            "Pushing artifact from local cache to remote registry with "
            f"reference '{reference}' failed with stderr:\n{stderr}"
        )
        super().__init__(msg)


class ORMBPullError(Exception):
    def __init__(self, reference: str, stderr: str) -> None:
        msg = (
            "Pulling artifact from remote registry with reference "
            f"'{reference}' failed with stderr:\n{stderr}"
        )
        super().__init__(msg)


class ORMBExportError(Exception):
    def __init__(self, reference: str, destination: str, stderr: str) -> None:
        msg = (
            f"Exporting artifact with reference '{reference}' from local "
            f"cache to '{destination}' failed with stderr:\n{stderr}"
        )
        super().__init__(msg)


class ORMBRemoveError(Exception):
    def __init__(self, reference: str, stderr: str) -> None:
        msg = (
            f"Removing artifact with reference '{reference}' from local cache "
            f"failed with stderr:\n{stderr}"
        )
        super().__init__(msg)
