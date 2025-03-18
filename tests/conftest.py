"""Test configuration for Nanodoc v2."""

import os
import tempfile

import pytest

from .utils import create_fixture_content_item


@pytest.fixture
def temp_dir():
    """Provide a temporary directory for tests."""
    with tempfile.TemporaryDirectory() as tmpdirname:
        yield tmpdirname


@pytest.fixture
def fixtures_dir():
    """Provide the path to the test fixtures directory."""
    return os.path.abspath(os.path.join(os.path.dirname(__file__), "fixtures"))


@pytest.fixture
def sample_file(temp_dir):
    """Create a sample text file in the temporary directory."""
    file_path = os.path.join(temp_dir, "sample.txt")
    content = "Line 1\nLine 2\nLine 3\nLine 4\nLine 5\n"

    with open(file_path, "w") as f:
        f.write(content)

    return file_path


@pytest.fixture
def sample_hidden_file(temp_dir):
    """Create a sample hidden file in the temporary directory."""
    file_path = os.path.join(temp_dir, ".hidden.txt")
    content = "Hidden content\n"

    with open(file_path, "w") as f:
        f.write(content)

    return file_path


@pytest.fixture
def nested_directory_structure(temp_dir):
    """Create a nested directory structure for testing."""
    # Create directories
    dir1 = os.path.join(temp_dir, "dir1")
    dir2 = os.path.join(temp_dir, "dir2")
    subdir = os.path.join(dir2, "subdir")
    hidden_dir = os.path.join(temp_dir, ".hidden_dir")

    os.makedirs(dir1, exist_ok=True)
    os.makedirs(dir2, exist_ok=True)
    os.makedirs(subdir, exist_ok=True)
    os.makedirs(hidden_dir, exist_ok=True)

    # Create files
    file1 = os.path.join(dir1, "file1.txt")
    file2 = os.path.join(dir1, "file2.md")
    file3 = os.path.join(dir2, "file3.txt")
    file4 = os.path.join(subdir, "file4.txt")
    hidden_file = os.path.join(dir2, ".hidden.txt")
    hidden_dir_file = os.path.join(hidden_dir, "file.txt")

    # Write content to files
    for file_path in [file1, file2, file3, file4, hidden_file, hidden_dir_file]:
        with open(file_path, "w") as f:
            f.write(f"Content of {os.path.basename(file_path)}\n")

    return {
        "base_dir": temp_dir,
        "dir1": dir1,
        "dir2": dir2,
        "subdir": subdir,
        "hidden_dir": hidden_dir,
        "files": [file1, file2, file3, file4, hidden_file, hidden_dir_file],
    }


# Add fixture_content_item parametrized fixture
test_files = [
    "test_file1.py",
    "test_file2.py",
    "test_bundle.ndoc",
    "cake.txt",
    "incident.txt",
    "new-telephone.txt",
]


@pytest.fixture(params=test_files)
def fixture_content_item(request):
    """Provide test fixture content items for parametrized tests.

    This fixture will create a FileContent object for each test file,
    allowing tests to run multiple times with different input files.
    """
    fixture_name = request.param
    return create_fixture_content_item(fixture_name)
