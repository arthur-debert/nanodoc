"""Global pytest fixtures."""

import os
from pathlib import Path

import pytest

from nanodoc.v2.structures import FileContent

# Import test utilities from the same directory
utils_path = os.path.join(os.path.dirname(__file__), "utils.py")
with open(utils_path) as f:
    exec(f.read(), globals())


@pytest.fixture
def fixtures_dir() -> Path:
    """Get the fixtures directory path."""
    return Path(__file__).parent / "fixtures"


@pytest.fixture(
    params=[
        "cake.txt",
        "incident.txt",
        "new-telephone.txt",
        "test_file1.py",
        "test_file2.py",
        "test_bundle.ndoc",
    ]
)
def fixture_file(request) -> str:
    """A parametrized fixture that provides access to test fixture files.

    Usage:
        def test_something(fixture_file):
            # fixture_file will be each of the files in sequence
            content = read_fixture(fixture_file)
            ...
    """
    return request.param


@pytest.fixture
def fixture_content(fixture_file: str) -> str:
    """Get the content of a fixture file.

    This fixture depends on fixture_file, so it will also be parametrized.

    Usage:
        def test_something(fixture_file, fixture_content):
            # fixture_file is the name
            # fixture_content is the actual content
            assert "some text" in fixture_content
    """
    return read_fixture(fixture_file)


@pytest.fixture
def fixture_path(fixture_file: str) -> Path:
    """Get the path to a fixture file.

    This fixture depends on fixture_file, so it will also be parametrized.

    Usage:
        def test_something(fixture_file, fixture_path):
            # fixture_file is the name
            # fixture_path is the full path
            assert fixture_path.exists()
    """
    return get_fixture_path(fixture_file)


FIXTURE_FILES = [
    "cake.txt",
    "incident.txt",
    "new-telephone.txt",
    "test_file1.py",
    "test_file2.py",
    "test_bundle.ndoc",
]


@pytest.fixture(params=FIXTURE_FILES)
def fixture_content_item(request) -> FileContent:
    """A parametrized fixture that provides FileContent instances for test files.

    This is the preferred way to access fixture files in tests. It provides
    a complete FileContent object that can be used with the nanodoc functions.

    The fixture is parametrized, so tests using it will run once for each
    fixture file.

    Usage:
        def test_something(fixture_content_item):
            # fixture_content_item is a FileContent instance
            result = run_all([fixture_content_item], ...)
            assert result ...

    Returns:
        FileContent instance for the current fixture file
    """
    return create_fixture_content_item(request.param)
