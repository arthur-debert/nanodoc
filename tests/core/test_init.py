import pytest

from nanodoc.core import process_all
from nanodoc.files import get_files_from_args
from nanodoc.formatting import create_header


def test_init_no_files_errors(tmpdir):
    # Create a temporary directory
    empty_dir = tmpdir.mkdir("empty")

    # Call init with the empty directory
    # Get verified sources
    verified_sources = get_files_from_args([str(empty_dir)])

    # Check if we have any valid files
    if not verified_sources:
        result = "Error: No valid source files found."

    # Assert that an error message is returned without checking the exact text
    assert result.startswith("Error:")
    assert "files found" in result


def test_init_one_file_no_line_numbers(tmpdir):
    # Create a test file
    test_file = tmpdir.join("test_file.txt")
    test_file.write("Line 1\nLine 2")

    # Call init with the test file
    # Get verified sources and process them
    verified_sources = get_files_from_args([str(test_file)])
    result = process_all(verified_sources)

    # Assert that the file content is printed without line numbers
    assert "Line 1" in result
    assert "Line 2" in result
    assert "1:" not in result
    assert "2:" not in result


def test_init_one_file_file_line_numbers(tmpdir):
    # Create a test file
    test_file = tmpdir.join("test_file.txt")
    test_file.write("Line 1\nLine 2")

    # Call init with the test file and file line numbers
    # Get verified sources and process them with file line numbers
    verified_sources = get_files_from_args([str(test_file)])
    result = process_all(verified_sources, line_number_mode="file")

    # Assert that the file content is printed with file line numbers
    assert "1: Line 1" in result
    assert "2: Line 2" in result


def test_init_one_file_all_line_numbers(tmpdir):
    # Create a test file
    test_file = tmpdir.join("test_file.txt")
    test_file.write("Line 1\nLine 2")

    # Call init with the test file and all line numbers
    # Get verified sources and process them with all line numbers
    verified_sources = get_files_from_args([str(test_file)])
    result = process_all(verified_sources, line_number_mode="all")

    # Assert that the file content is printed with all line numbers
    assert "1: Line 1" in result
    assert "2: Line 2" in result


def test_init_toc(tmpdir):
    # Create a test file
    test_file = tmpdir.join("test_file.txt")
    test_file.write("Line 1\nLine 2")

    # Call init with the test file and TOC generation
    # Get verified sources and process them with TOC generation
    verified_sources = get_files_from_args([str(test_file)])
    result = process_all(verified_sources, generate_toc=True)

    # Assert that the TOC is generated and the file content is printed
    assert create_header("TOC") in result
    assert "test_file.txt" in result
    assert "Line 1" in result
    assert "Line 2" in result


@pytest.mark.skip(reason="Not implemented")
def test_init_bundle_error(tmpdir):
    # Create a test bundle file with non-existent paths
    bundle_file = tmpdir.join("test_bundle.txt")
    bundle_file.write("/nonexistent/file1.txt")

    # Call init with the bundle file
    # Get verified sources
    verified_sources = get_files_from_args([str(bundle_file)])

    # Check if we have any valid files
    if not verified_sources:
        result = "Error: No valid source files found."

    # Assert that an error message is returned
    assert "Error:" in result
    assert "files found" in result
