"""End-to-end tests for the nanodoc v2 CLI."""

import os
import subprocess
import tempfile
from pathlib import Path

import pytest


@pytest.fixture
def test_file():
    """Create a temporary test file."""
    with tempfile.NamedTemporaryFile(suffix=".py", delete=False, mode="w") as temp:
        temp.write("# Test Header\n\ndef test_function():\n    pass\n")
        temp_name = temp.name

    yield temp_name
    os.unlink(temp_name)


@pytest.fixture
def test_files():
    """Create multiple temporary test files."""
    temp_files = []

    # Create 3 test files
    for i in range(3):
        with tempfile.NamedTemporaryFile(suffix=".py", delete=False, mode="w") as temp:
            temp.write(f"# File {i+1}\n\ndef function_{i+1}():\n    pass\n")
            temp_files.append(temp.name)

    yield temp_files

    # Clean up the files
    for filename in temp_files:
        os.unlink(filename)


def test_basic_output(test_file):
    """Test basic CLI output."""
    # Run the command
    result = subprocess.run(
        ["python", "-m", "nanodoc", "--use-v2", test_file],
        capture_output=True,
        text=True,
        check=True,
    )

    # Check the output contains the file content
    assert "Test Header" in result.stdout
    assert "test_function" in result.stdout
    # The file name should appear in the output
    assert Path(test_file).name in result.stdout


def test_toc_generation(test_files):
    """Test TOC generation."""
    # Run the command with TOC enabled
    result = subprocess.run(
        ["python", "-m", "nanodoc", "--use-v2", "--toc"] + test_files,
        capture_output=True,
        text=True,
        check=True,
    )

    # Check that TOC is generated
    assert "TOC" in result.stdout or "Table of Contents" in result.stdout

    # Check that all file names appear in the output
    for filename in test_files:
        assert Path(filename).name in result.stdout

    # Check that content from all files is included
    assert "File 1" in result.stdout
    assert "File 2" in result.stdout
    assert "File 3" in result.stdout


def test_line_number_mode(test_file):
    """Test line number mode."""
    # Run the command with line numbers
    result = subprocess.run(
        ["python", "-m", "nanodoc", "--use-v2", "-n", test_file],
        capture_output=True,
        text=True,
        check=True,
    )

    # Check that line numbers are included
    assert "1" in result.stdout and "2" in result.stdout
    assert "def test_function" in result.stdout


def test_theme_option(test_file):
    """Test theme option."""
    # Run the command with a theme
    result = subprocess.run(
        ["python", "-m", "nanodoc", "--use-v2", "--theme", "neutral", test_file],
        capture_output=True,
        text=True,
        check=True,
    )

    # Check that the theme is applied (can only verify it runs without errors)
    assert "Test Header" in result.stdout
    assert "test_function" in result.stdout


def test_multiple_options(test_files):
    """Test multiple options together."""
    # Run the command with multiple options
    result = subprocess.run(
        ["python", "-m", "nanodoc", "--use-v2", "--toc", "-n", "--theme", "neutral"]
        + test_files,
        capture_output=True,
        text=True,
        check=True,
    )

    # Check that all expected elements are present
    assert "TOC" in result.stdout or "Table of Contents" in result.stdout
    assert "1" in result.stdout  # Line numbers

    # Check content from all files
    for i in range(3):
        assert f"File {i+1}" in result.stdout
        assert f"function_{i+1}" in result.stdout


def test_invalid_arguments():
    """Test with invalid arguments."""
    # Run the command with invalid sources
    result = subprocess.run(
        ["python", "-m", "nanodoc", "--use-v2", "--nonexistent-file"],
        capture_output=True,
        text=True,
    )

    # Should exit with non-zero status
    assert result.returncode != 0
    # Should contain an error message
    assert "Error" in result.stderr
