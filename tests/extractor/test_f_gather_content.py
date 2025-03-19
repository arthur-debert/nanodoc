"""Tests for extractor.gather_content function."""

import pytest

from nanodoc.extractor import gather_content
from nanodoc.structures import FileContent


def test_gather_content_basic(tmp_path):
    """Test gathering content from files."""
    # Create temporary files for testing
    file1 = tmp_path / "file1.txt"
    file1.write_text("Line 1\nLine 2\nLine 3\n")

    file2 = tmp_path / "file2.txt"
    file2.write_text("Line A\nLine B\nLine C\n")

    # Create FileContent objects
    file_contents = [
        FileContent(filepath=str(file1), ranges=[(1, None)], is_bundle=False),
        FileContent(filepath=str(file2), ranges=[(1, None)], is_bundle=False),
    ]

    # Gather content
    result = gather_content(file_contents)

    # Check result
    assert len(result) == 2
    assert result[0].content == "Line 1\nLine 2\nLine 3\n"
    assert result[1].content == "Line A\nLine B\nLine C\n"


def test_gather_content_with_ranges(tmp_path):
    """Test gathering content from files with ranges."""
    # Create temporary file for testing
    file1 = tmp_path / "file1.txt"
    file1.write_text("Line 1\nLine 2\nLine 3\nLine 4\nLine 5\n")

    # Create FileContent objects with ranges
    file_contents = [
        FileContent(filepath=str(file1), ranges=[(2, 4)], is_bundle=False),
    ]

    # Gather content
    result = gather_content(file_contents)

    # Check result
    assert len(result) == 1
    assert result[0].content == "Line 2\nLine 3\n"


def test_gather_content_file_not_found():
    """Test gathering content from non-existent files."""
    # Create FileContent object for non-existent file
    file_contents = [
        FileContent(
            filepath="/path/to/nonexistent/file.txt",
            ranges=[(1, None)],
            is_bundle=False,
        ),
    ]

    # Should raise FileNotFoundError
    with pytest.raises(FileNotFoundError):
        gather_content(file_contents)
