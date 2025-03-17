"""Tests for the guides functionality."""

import os
import subprocess
import sys

# Constants for testing
PYTHON_CMD = sys.executable
NANODOC_MODULE = "src.nanodoc"

# Get the parent directory of the current module
MODULE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))


def test_help_command():
    """Test the basic help command."""
    result = subprocess.run(
        [PYTHON_CMD, "-m", NANODOC_MODULE, "help"],
        capture_output=True,
        text=True,
    )
    # Check for successful execution
    assert result.returncode == 0

    # Check for key content sections that should be present
    # These are less likely to be affected by Rich formatting
    assert "nanodoc" in result.stdout
    assert "file1.txt" in result.stdout  # From examples
    assert "file2.txt" in result.stdout  # From examples
    assert "manifesto" in result.stdout  # From help topics
    assert "quickstart" in result.stdout  # From help topics


def test_help_with_guide():
    """Test the help command with a guide parameter."""
    # Test with the manifesto guide (which should exist in the project)
    result = subprocess.run(
        [PYTHON_CMD, "-m", NANODOC_MODULE, "help", "manifesto"],
        capture_output=True,
        text=True,
    )
    # Check for successful execution
    assert result.returncode == 0

    # Check for key content that should be present in the manifesto guide
    # These are less likely to be affected by Rich formatting
    assert "manifesto" in result.stdout.lower()
    assert "less clutter" in result.stdout.lower()


def test_help_with_nonexistent_guide():
    """Test the help command with a nonexistent guide parameter."""
    # Run the command with a nonexistent guide
    nonexistent_guide = "nonexistent_guide_xyz"
    result = subprocess.run(
        [PYTHON_CMD, "-m", NANODOC_MODULE, "help", nonexistent_guide],
        capture_output=True,
        text=True,
    )

    # The command should fail with a non-zero exit code
    assert result.returncode != 0

    # The error message should indicate that the guide was not found
    # This is in stderr and less likely to be affected by Rich formatting
    assert f"Guide '{nonexistent_guide}' not found" in result.stderr

    # It should list available guides
    assert "Available guides" in result.stderr
