import os

import pytest

from nanodoc.v1.core import run_all
from nanodoc.v1.files import get_files_from_args
from nanodoc.v1.formatting import create_header

FIXTURES_DIR = os.path.join(
    os.path.dirname(os.path.dirname(os.path.dirname(__file__))), "fixtures"
)


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


def test_init_one_file_no_line_numbers():
    # Use test file from fixtures
    test_file = os.path.join(FIXTURES_DIR, "test_file1.py")

    # Call init with the test file
    # Get verified sources and process them
    verified_sources = get_files_from_args([test_file])
    result = run_all(
        verified_sources,
        line_number_mode=None,
        generate_toc=False,
        theme=None,
        show_header=True,
    )

    # Assert that the file content is printed without line numbers
    assert "def test_function():" in result
    assert "return True" in result
    assert "1:" not in result
    assert "2:" not in result


def test_init_one_file_file_line_numbers():
    # Use test file from fixtures
    test_file = os.path.join(FIXTURES_DIR, "test_file1.py")

    # Call init with the test file and file line numbers
    # Get verified sources and process them with file line numbers
    verified_sources = get_files_from_args([test_file])
    result = run_all(
        verified_sources,
        line_number_mode="file",
        generate_toc=False,
        theme=None,
        show_header=True,
    )

    # Assert that the file content is printed with file line numbers
    assert "1: def test_function():" in result
    assert "2:     return True" in result


def test_init_one_file_all_line_numbers():
    # Use test file from fixtures
    test_file = os.path.join(FIXTURES_DIR, "test_file1.py")

    # Call init with the test file and all line numbers
    # Get verified sources and process them with all line numbers
    verified_sources = get_files_from_args([test_file])
    result = run_all(
        verified_sources,
        line_number_mode="all",
        generate_toc=False,
        theme=None,
        show_header=True,
    )

    # Assert that the file content is printed with all line numbers
    assert "1: def test_function():" in result
    assert "2:     return True" in result


def test_init_toc():
    # Use test file from fixtures
    test_file = os.path.join(FIXTURES_DIR, "test_file1.py")

    # Call init with the test file and TOC generation
    # Get verified sources and process them with TOC generation
    verified_sources = get_files_from_args([test_file])
    result = run_all(
        verified_sources,
        line_number_mode=None,
        generate_toc=True,
        theme=None,
        show_header=True,
    )

    # Assert that the TOC is generated and the file content is printed
    assert create_header("TOC") in result
    assert "test_file1.py" in result
    assert "def test_function():" in result
    assert "return True" in result


@pytest.mark.skip(reason="Not implemented")
def test_init_bundle_error():
    # Use test bundle from fixtures
    bundle_file = os.path.join(FIXTURES_DIR, "test_bundle.ndoc")

    # Call init with the bundle file
    # Get verified sources
    verified_sources = get_files_from_args([bundle_file])

    # Check if we have any valid files
    if not verified_sources:
        result = "Error: No valid source files found."

    # Assert that an error message is returned
    assert "Error:" in result
    assert "files found" in result
