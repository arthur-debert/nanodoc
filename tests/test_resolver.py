"""Tests for path resolution in Nanodoc v2."""

import os

import pytest

from nanodoc.resolver import get_files_from_directory, resolve_paths


def test_resolve_single_file(sample_file):
    """Test resolving a single file path."""
    result = resolve_paths([sample_file])
    assert len(result) == 1
    assert result[0] == os.path.abspath(sample_file)


def test_resolve_missing_file():
    """Test resolving a non-existent file."""
    non_existent_file = "/path/to/non_existent_file.txt"
    with pytest.raises(FileNotFoundError):
        resolve_paths([non_existent_file])


def test_resolve_directory(nested_directory_structure):
    """Test resolving a directory (non-recursive)."""
    dir1 = nested_directory_structure["dir1"]
    result = resolve_paths([dir1], recursive=False)

    # Should return all files in dir1 (2 files)
    assert len(result) == 2

    # All returned paths should be absolute
    for path in result:
        assert os.path.isabs(path)

    # Check that all files are from dir1
    for path in result:
        assert os.path.dirname(path) == dir1


def test_resolve_directory_recursive(nested_directory_structure):
    """Test resolving a directory recursively."""
    base_dir = nested_directory_structure["base_dir"]
    result = resolve_paths([base_dir], recursive=True, include_hidden=False)

    # Count the non-hidden files in the test structure: 4 files (excluding hidden)
    expected_count = 4
    assert len(result) == expected_count

    # All returned paths should be absolute
    for path in result:
        assert os.path.isabs(path)

    # Check that no hidden files are included
    for path in result:
        assert not os.path.basename(path).startswith(".")


def test_resolve_directory_with_hidden(nested_directory_structure):
    """Test resolving a directory including hidden files."""
    base_dir = nested_directory_structure["base_dir"]
    result = resolve_paths([base_dir], recursive=True, include_hidden=True)

    # Count all files in the test structure: 6 files (including hidden)
    expected_count = 6
    assert len(result) == expected_count


def test_resolve_glob(nested_directory_structure):
    """Test resolving paths with glob patterns."""
    base_dir = nested_directory_structure["base_dir"]

    # Find all .txt files
    txt_glob = os.path.join(base_dir, "**", "*.txt")
    result = resolve_paths([txt_glob], recursive=True, include_hidden=False)

    # We should find 3 non-hidden .txt files
    expected_count = 3
    assert len(result) == expected_count

    # All files should have .txt extension
    for path in result:
        assert path.endswith(".txt")

    # No hidden files should be included
    for path in result:
        assert not os.path.basename(path).startswith(".")


def test_resolve_glob_with_hidden(nested_directory_structure):
    """Test resolving paths with glob patterns including hidden files."""
    base_dir = nested_directory_structure["base_dir"]

    # Find all .txt files including hidden
    txt_glob = os.path.join(base_dir, "**", "*.txt")
    result = resolve_paths([txt_glob], recursive=True, include_hidden=True)

    # We should find 5 .txt files (including hidden)
    expected_count = 5
    assert len(result) == expected_count

    # All files should have .txt extension
    for path in result:
        assert path.endswith(".txt")


def test_resolve_multiple_inputs(nested_directory_structure):
    """Test resolving multiple input paths."""
    dir1 = nested_directory_structure["dir1"]
    file3 = next(
        f
        for f in nested_directory_structure["files"]
        if os.path.basename(f) == "file3.txt"
    )

    result = resolve_paths([dir1, file3], recursive=False)

    # Should return 2 files from dir1 + 1 individual file
    expected_count = 3
    assert len(result) == expected_count

    # Check that file3 is in the result
    assert file3 in result


def test_get_files_from_directory(nested_directory_structure):
    """Test getting files from a directory."""
    dir1 = nested_directory_structure["dir1"]
    result = get_files_from_directory(dir1, recursive=False)

    # Should return 2 files from dir1
    assert len(result) == 2

    # All paths should be absolute
    for path in result:
        assert os.path.isabs(path)

    # All paths should be from dir1
    for path in result:
        assert os.path.dirname(path) == dir1


def test_get_files_from_directory_recursive(nested_directory_structure):
    """Test getting files from a directory recursively."""
    dir2 = nested_directory_structure["dir2"]
    result = get_files_from_directory(dir2, recursive=True, include_hidden=False)

    # Should return 2 files: file3.txt and subdir/file4.txt
    assert len(result) == 2

    # Check we have the expected files
    basenames = [os.path.basename(path) for path in result]
    assert "file3.txt" in basenames
    assert "file4.txt" in basenames


def test_resolve_paths_deduplication():
    """Test that resolve_paths deduplicates the results."""
    with pytest.raises(FileNotFoundError):
        # This should fail, but we're checking deduplication logic
        resolve_paths(["non-existent-path"])

    # Instead, let's create a temporary file to test deduplication
    with open("temp_test_file.txt", "w") as f:
        f.write("Test content")

    try:
        # Try to resolve the same file multiple times
        paths = [
            "temp_test_file.txt",
            "temp_test_file.txt",
            os.path.abspath("temp_test_file.txt"),
        ]
        result = resolve_paths(paths)

        # Even though we provided the same file in different ways,
        # we should only get one result
        assert len(result) == 1
        assert result[0] == os.path.abspath("temp_test_file.txt")
    finally:
        # Clean up
        os.remove("temp_test_file.txt")
