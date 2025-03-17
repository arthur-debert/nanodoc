"""Tests for the guides functionality."""

import os
import pathlib
import subprocess
import sys
import unittest.mock as mock

# Constants for testing
PYTHON_CMD = sys.executable
NANODOC_MODULE = "src.nanodoc"

# Get the parent directory of the current module
MODULE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))


def test_get_available_guides():
    """Test that get_available_guides returns a dictionary of available guides."""
    # Create a temporary test directory
    test_dir = pathlib.Path(__file__).parent / "fixtures" / "guides"
    test_dir.mkdir(parents=True, exist_ok=True)

    # Create manifesto guide
    manifesto_path = test_dir / "manifesto.txt"
    with open(manifesto_path, "w") as f:
        f.write("Nanodoc Manifesto\n\nTest manifesto content")

    # Create quickstart guide
    quickstart_path = test_dir / "quickstart.txt"
    with open(quickstart_path, "w") as f:
        f.write("# Quickstart Guide\n\nTest quickstart content")

    try:
        # Mock the _get_guides_dir function to return our test directory
        with mock.patch("nanodoc.help._get_guides_dir", return_value=test_dir):
            from nanodoc.help import get_available_guides

            guides = get_available_guides()

            # Check that we have the guides we created
            assert "manifesto" in guides
            assert "quickstart" in guides

            # Check that the descriptions are correct
            assert guides["manifesto"] == "Nanodoc Manifesto"
            assert guides["quickstart"] == "Quickstart Guide"
    finally:
        # Clean up the temporary files
        if manifesto_path.exists():
            manifesto_path.unlink()
        if quickstart_path.exists():
            quickstart_path.unlink()


def test_get_guide_content_existing():
    """Test that get_guide_content returns the content of an existing guide."""
    # Create a temporary test directory
    test_dir = pathlib.Path(__file__).parent / "fixtures" / "guides"
    test_dir.mkdir(parents=True, exist_ok=True)

    # Create manifesto guide
    manifesto_path = test_dir / "manifesto.txt"
    manifesto_content = "Nanodoc Manifesto\n\nTest manifesto content"
    with open(manifesto_path, "w") as f:
        f.write(manifesto_content)

    # Create quickstart guide
    quickstart_path = test_dir / "quickstart.txt"
    quickstart_content = "# Quickstart Guide\n\nTest quickstart content"
    with open(quickstart_path, "w") as f:
        f.write(quickstart_content)

    try:
        # Mock the _get_guides_dir function to return our test directory
        with mock.patch("nanodoc.help._get_guides_dir", return_value=test_dir):
            from nanodoc.help import get_guide_content

            # Test with the manifesto guide
            found, content = get_guide_content("manifesto")
            assert found is True
            assert content == manifesto_content

            # Test with the quickstart guide
            found, content = get_guide_content("quickstart")
            assert found is True
            assert content == quickstart_content
    finally:
        # Clean up the temporary files
        if manifesto_path.exists():
            manifesto_path.unlink()
        if quickstart_path.exists():
            quickstart_path.unlink()


def test_get_guide_content_nonexistent():
    """Test that get_guide_content returns an error message for a non-existent guide."""
    # Create a temporary test directory
    test_dir = pathlib.Path(__file__).parent / "fixtures" / "guides"
    test_dir.mkdir(parents=True, exist_ok=True)

    # Create manifesto guide
    manifesto_path = test_dir / "manifesto.txt"
    with open(manifesto_path, "w") as f:
        f.write("Nanodoc Manifesto\n\nTest manifesto content")

    try:
        # Mock the _get_guides_dir function to return our test directory
        with mock.patch("nanodoc.help._get_guides_dir", return_value=test_dir):
            from nanodoc.help import get_guide_content

            found, content = get_guide_content("nonexistent")
            assert found is False
            assert "Guide 'nonexistent' not found" in content
            assert "Available guides:" in content
            assert "manifesto" in content
    finally:
        # Clean up the temporary files
        if manifesto_path.exists():
            manifesto_path.unlink()


def test_help_with_guide():
    """Test the help command with a guide parameter."""
    # Create a temporary test directory
    test_dir = pathlib.Path(__file__).parent / "fixtures" / "guides"
    test_dir.mkdir(parents=True, exist_ok=True)

    # Create manifesto guide
    manifesto_path = test_dir / "manifesto.txt"
    with open(manifesto_path, "w") as f:
        f.write("Test manifesto content")

    # Create quickstart guide
    quickstart_path = test_dir / "quickstart.txt"
    with open(quickstart_path, "w") as f:
        f.write("Test quickstart content")

    try:
        # Create a symlink from the actual guides directory to our test directory
        guides_dir = (
            pathlib.Path(__file__).parent.parent / "src" / "nanodoc" / "docs" / "guides"
        )
        guides_dir.parent.mkdir(parents=True, exist_ok=True)

        # If the guides directory exists, rename it temporarily
        backup_dir = None
        if guides_dir.exists():
            backup_dir = guides_dir.parent / "guides_backup"
            guides_dir.rename(backup_dir)

        # Create a symlink to our test directory
        os.symlink(test_dir, guides_dir)

        try:
            # Test manifesto guide
            result = subprocess.run(
                [PYTHON_CMD, "-m", NANODOC_MODULE, "help", "manifesto"],
                capture_output=True,
                text=True,
            )
            assert result.returncode == 0
            assert "MANIFESTO GUIDE" in result.stdout
            assert "Test manifesto content" in result.stdout

            # Test quickstart guide
            result = subprocess.run(
                [PYTHON_CMD, "-m", NANODOC_MODULE, "help", "quickstart"],
                capture_output=True,
                text=True,
            )
            assert result.returncode == 0
            assert "QUICKSTART GUIDE" in result.stdout
            assert "Test quickstart content" in result.stdout
        finally:
            # Remove the symlink
            if guides_dir.exists():
                if os.path.islink(guides_dir):
                    os.unlink(guides_dir)
                else:
                    import shutil

                    shutil.rmtree(guides_dir)

            # Restore the original guides directory if it existed
            if backup_dir is not None:
                backup_dir.rename(guides_dir)
    finally:
        # Clean up the temporary files
        if manifesto_path.exists():
            manifesto_path.unlink()
        if quickstart_path.exists():
            quickstart_path.unlink()


def test_help_with_nonexistent_guide():
    """Test the help command with a nonexistent guide parameter."""
    # Create a temporary test directory
    test_dir = pathlib.Path(__file__).parent / "fixtures" / "guides"
    test_dir.mkdir(parents=True, exist_ok=True)

    try:
        # Create a symlink from the actual guides directory to our test directory
        guides_dir = (
            pathlib.Path(__file__).parent.parent / "src" / "nanodoc" / "docs" / "guides"
        )
        guides_dir.parent.mkdir(parents=True, exist_ok=True)

        # If the guides directory exists, rename it temporarily
        backup_dir = None
        if guides_dir.exists():
            backup_dir = guides_dir.parent / "guides_backup"
            guides_dir.rename(backup_dir)

        # Create a symlink to our test directory
        os.symlink(test_dir, guides_dir)

        try:
            # Run the command with a nonexistent guide
            result = subprocess.run(
                [PYTHON_CMD, "-m", NANODOC_MODULE, "help", "nonexistent"],
                capture_output=True,
                text=True,
            )

            # The command should fail with a non-zero exit code
            assert result.returncode != 0

            # The error message should indicate that the guide was not found
            assert "Guide 'nonexistent' not found" in result.stderr

            # It should list available guides
            assert "Available guides" in result.stderr
        finally:
            # Remove the symlink
            if guides_dir.exists():
                if os.path.islink(guides_dir):
                    os.unlink(guides_dir)
                else:
                    import shutil

                    shutil.rmtree(guides_dir)

            # Restore the original guides directory if it existed
            if backup_dir is not None:
                backup_dir.rename(guides_dir)
    finally:
        # If the test directory exists, remove it
        if test_dir.exists():
            import shutil

            shutil.rmtree(test_dir)
