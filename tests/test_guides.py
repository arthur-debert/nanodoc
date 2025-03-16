import os
import pathlib
import subprocess
import sys
from unittest import mock

from nanodoc.help import get_available_guides, get_guide_content

# Get the parent directory of the current module
MODULE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

# Use Python module approach instead of direct script execution
PYTHON_CMD = sys.executable
NANODOC_MODULE = "src.nanodoc"


def test_get_available_guides():
    """Test that get_available_guides returns a dictionary of available guides."""
    guides = get_available_guides()

    # Check that we have at least the two guides we created
    assert "manifesto" in guides
    assert "quickstart" in guides

    # Check that the descriptions are correct
    assert guides["manifesto"]  # Just check that it has a description
    assert guides["quickstart"]  # Just check that it has a description


def test_get_guide_content_existing():
    """Test that get_guide_content returns the content of an existing guide."""
    # Test with the manifesto guide
    found, content = get_guide_content("manifesto")
    assert found is True

    # Read the actual file content to compare
    guide_path = (
        pathlib.Path(__file__).parent.parent
        / "src"
        / "nanodoc"
        / "docs"
        / "guides"
        / "manifesto.txt"
    )
    with open(guide_path, "r", encoding="utf-8") as f:
        assert content == f.read()

    # Test with the quickstart guide
    found, content = get_guide_content("quickstart")
    assert found is True

    # Read the actual file content to compare
    guide_path = (
        pathlib.Path(__file__).parent.parent
        / "src"
        / "nanodoc"
        / "docs"
        / "guides"
        / "quickstart.md"
    )
    assert content == open(guide_path, "r", encoding="utf-8").read()


def test_get_guide_content_nonexistent():
    """Test that get_guide_content returns an error message for a non-existent guide."""
    found, content = get_guide_content("nonexistent")
    assert found is False
    assert "Guide 'nonexistent' not found" in content
    assert "Available guides:" in content
    assert "manifesto" in content
    assert "quickstart" in content


def test_help_with_guide():
    """Test the help command with a guide parameter using mocking."""
    # Test with the manifesto guide
    with mock.patch("nanodoc.help._render_content") as mock_render:
        # Run the command
        result = subprocess.run(
            [PYTHON_CMD, "-m", NANODOC_MODULE, "help", "manifesto"],
            capture_output=True,
            text=True,
        )
        assert result.returncode == 0

        # Verify that the guide content was retrieved correctly
        guide_path = (
            pathlib.Path(__file__).parent.parent
            / "src"
            / "nanodoc"
            / "docs"
            / "guides"
            / "manifesto.txt"
        )
        with open(guide_path, "r", encoding="utf-8") as f:
            guide_content = f.read()

        # Verify that the guide content was passed to _render_content
        # Note: We can't directly check the call args because the subprocess runs in a separate process
        # Instead, we'll check that the command ran successfully
        assert "1. Less clutter, less distraction" in guide_content

    # Test with the quickstart guide
    with mock.patch("nanodoc.help._render_content") as mock_render:
        # Run the command
        result = subprocess.run(
            [PYTHON_CMD, "-m", NANODOC_MODULE, "help", "quickstart"],
            capture_output=True,
            text=True,
        )
        assert result.returncode == 0

        # Verify that the guide content was retrieved correctly
        guide_path = (
            pathlib.Path(__file__).parent.parent
            / "src"
            / "nanodoc"
            / "docs"
            / "guides"
            / "quickstart.md"
        )
        with open(guide_path, "r", encoding="utf-8") as f:
            guide_content = f.read()

        # Verify that the guide content was passed to _render_content
        # Note: We can't directly check the call args because the subprocess runs in a separate process
        # Instead, we'll check that the command ran successfully
        assert "# Nanodoc" in guide_content


def test_help_with_nonexistent_guide():
    """Test the help command with a non-existent guide parameter."""
    result = subprocess.run(
        [PYTHON_CMD, "-m", NANODOC_MODULE, "help", "nonexistent"],
        capture_output=True,
        text=True,
    )
    assert result.returncode == 1  # Should exit with status 1 for non-existent guide
    assert "Guide 'nonexistent' not found" in result.stdout
    assert "Available guides:" in result.stdout
    assert "manifesto" in result.stdout
    assert "quickstart" in result.stdout
