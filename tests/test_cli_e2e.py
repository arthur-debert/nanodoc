"""End-to-end tests for the nanodoc v2 CLI."""

import os
import subprocess
import tempfile

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
            temp.write(f"# File {i + 1}\n\ndef function_{i + 1}():\n    pass\n")
            temp_files.append(temp.name)

    yield temp_files

    # Clean up the files
    for filename in temp_files:
        os.unlink(filename)


def run_command(args: list[str]) -> subprocess.CompletedProcess:
    """Run a command and return the CompletedProcess object."""
    return subprocess.run(
        args,
        text=True,
        capture_output=True,
        check=True,
    )


def test_basic_output(fixture_content_item):
    """Test basic CLI output."""
    result = subprocess.run(
        ["python", "-m", "nanodoc", fixture_content_item.file_path],
        text=True,
        capture_output=True,
        check=True,
    )

    # Check that the output contains the content
    assert result.returncode == 0

    # Check that file name is in output
    filename = os.path.basename(fixture_content_item.file_path)
    if not fixture_content_item.original_source.endswith(".ndoc"):
        assert filename in result.stdout

    # Check for key content by looking for smaller distinct snippets
    if "cake.txt" in fixture_content_item.file_path:
        assert "appalling" in result.stdout
        assert "cake consumption" in result.stdout
    elif "incident.txt" in fixture_content_item.file_path:
        assert "Team" in result.stdout
        assert "Palmer" in result.stdout
    elif "new-telephone.txt" in fixture_content_item.file_path:
        assert "All Employees" in result.stdout
        assert "funny meme" in result.stdout
    elif "test_file1.py" in fixture_content_item.file_path:
        assert "Test File 1" in result.stdout
        assert "function_1" in result.stdout
    elif "test_file2.py" in fixture_content_item.file_path:
        assert "Test File 2" in result.stdout
        assert "CONSTANT" in result.stdout
    elif "test_bundle.ndoc" in fixture_content_item.file_path:
        assert "Test Bundle" in result.stdout
        assert "bundle_function" in result.stdout
        # For bundles, check that included content is present
        assert "function_1" in result.stdout


def test_toc_generation(fixture_content_item):
    """Test table of contents generation."""
    result = subprocess.run(
        ["python", "-m", "nanodoc", "--toc", fixture_content_item.file_path],
        text=True,
        capture_output=True,
        check=True,
    )

    # Check that the output contains TOC and content
    assert result.returncode == 0
    assert "Table of Contents" in result.stdout

    # Check that file name is in output
    filename = os.path.basename(fixture_content_item.file_path)
    if not fixture_content_item.original_source.endswith(".ndoc"):
        assert filename in result.stdout

    # Check for key content by looking for smaller distinct snippets
    if "cake.txt" in fixture_content_item.file_path:
        assert "appalling" in result.stdout
        assert "cake consumption" in result.stdout
    elif "incident.txt" in fixture_content_item.file_path:
        assert "Team" in result.stdout
        assert "Palmer" in result.stdout
    elif "new-telephone.txt" in fixture_content_item.file_path:
        assert "All Employees" in result.stdout
        assert "funny meme" in result.stdout
    elif "test_file1.py" in fixture_content_item.file_path:
        assert "Test File 1" in result.stdout
        assert "function_1" in result.stdout
    elif "test_file2.py" in fixture_content_item.file_path:
        assert "Test File 2" in result.stdout
        assert "CONSTANT" in result.stdout
    elif "test_bundle.ndoc" in fixture_content_item.file_path:
        assert "Test Bundle" in result.stdout
        assert "bundle_function" in result.stdout
        # For bundles, check that included content is present
        assert "function_1" in result.stdout


def test_line_number_mode(fixture_content_item):
    """Test line number display modes."""
    result = subprocess.run(
        ["python", "-m", "nanodoc", "-n", fixture_content_item.file_path],
        text=True,
        capture_output=True,
        check=True,
    )

    # Check that line numbers are included
    assert result.returncode == 0
    assert "1:" in result.stdout

    # Check that file name is in output
    filename = os.path.basename(fixture_content_item.file_path)
    if not fixture_content_item.original_source.endswith(".ndoc"):
        assert filename in result.stdout

    # Check for key content by looking for smaller distinct snippets
    if "cake.txt" in fixture_content_item.file_path:
        assert "appalling" in result.stdout
        assert "cake consumption" in result.stdout
    elif "incident.txt" in fixture_content_item.file_path:
        assert "Team" in result.stdout
        assert "Palmer" in result.stdout
    elif "new-telephone.txt" in fixture_content_item.file_path:
        assert "All Employees" in result.stdout
        assert "funny meme" in result.stdout
    elif "test_file1.py" in fixture_content_item.file_path:
        assert "Test File 1" in result.stdout
        assert "function_1" in result.stdout
    elif "test_file2.py" in fixture_content_item.file_path:
        assert "Test File 2" in result.stdout
        assert "CONSTANT" in result.stdout
    elif "test_bundle.ndoc" in fixture_content_item.file_path:
        assert "Test Bundle" in result.stdout
        assert "bundle_function" in result.stdout
        # For bundles, check that included content is present
        assert "function_1" in result.stdout


@pytest.mark.skip(reason="Theme functionality requires v1 theme files")
def test_theme_option(fixture_content_item):
    """Test theme application."""
    cmd = [
        "python",
        "-m",
        "nanodoc",
        "--theme",
        "neutral",
        fixture_content_item.file_path,
    ]

    # Don't use check=True since some themes might be missing
    result = subprocess.run(cmd, text=True, capture_output=True)

    # Check that file name is in output
    filename = os.path.basename(fixture_content_item.file_path)
    if not fixture_content_item.original_source.endswith(".ndoc"):
        assert filename in result.stdout

    # Check for key content by looking for smaller distinct snippets
    if "cake.txt" in fixture_content_item.file_path:
        assert "appalling" in result.stdout
        assert "cake consumption" in result.stdout
    elif "incident.txt" in fixture_content_item.file_path:
        assert "Team" in result.stdout
        assert "Palmer" in result.stdout
    elif "new-telephone.txt" in fixture_content_item.file_path:
        assert "All Employees" in result.stdout
        assert "funny meme" in result.stdout
    elif "test_file1.py" in fixture_content_item.file_path:
        assert "Test File 1" in result.stdout
        assert "function_1" in result.stdout
    elif "test_file2.py" in fixture_content_item.file_path:
        assert "Test File 2" in result.stdout
        assert "CONSTANT" in result.stdout
    elif "test_bundle.ndoc" in fixture_content_item.file_path:
        assert "Test Bundle" in result.stdout
        assert "bundle_function" in result.stdout
        # For bundles, check that included content is present
        assert "function_1" in result.stdout


@pytest.mark.skip(reason="Theme functionality requires v1 theme files")
def test_multiple_options(fixture_content_item):
    """Test multiple options together."""
    cmd = [
        "python",
        "-m",
        "nanodoc",
        "--toc",
        "-n",
        "--theme",
        "neutral",
        fixture_content_item.file_path,
    ]

    # Don't use check=True since some themes might be missing
    result = subprocess.run(cmd, text=True, capture_output=True)

    assert "Table of Contents" in result.stdout
    assert "1:" in result.stdout

    # Check that file name is in output
    filename = os.path.basename(fixture_content_item.file_path)
    if not fixture_content_item.original_source.endswith(".ndoc"):
        assert filename in result.stdout

    # Check for key content by looking for smaller distinct snippets
    if "cake.txt" in fixture_content_item.file_path:
        assert "appalling" in result.stdout
        assert "cake consumption" in result.stdout
    elif "incident.txt" in fixture_content_item.file_path:
        assert "Team" in result.stdout
        assert "Palmer" in result.stdout
    elif "new-telephone.txt" in fixture_content_item.file_path:
        assert "All Employees" in result.stdout
        assert "funny meme" in result.stdout
    elif "test_file1.py" in fixture_content_item.file_path:
        assert "Test File 1" in result.stdout
        assert "function_1" in result.stdout
    elif "test_file2.py" in fixture_content_item.file_path:
        assert "Test File 2" in result.stdout
        assert "CONSTANT" in result.stdout
    elif "test_bundle.ndoc" in fixture_content_item.file_path:
        assert "Test Bundle" in result.stdout
        assert "bundle_function" in result.stdout
        # For bundles, check that included content is present
        assert "function_1" in result.stdout


def test_error_handling():
    """Test error handling for invalid file paths."""
    with pytest.raises(subprocess.CalledProcessError):
        subprocess.run(
            ["python", "-m", "nanodoc", "--nonexistent-file"],
            text=True,
            capture_output=True,
            check=True,
        )
