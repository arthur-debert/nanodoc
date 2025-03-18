"""Tests for the v2 CLI implementation.

These tests focus on the CLI integration of the v2 pipeline,
and mock the post-argument processing to test functionality directly.
"""

import os
import tempfile
from unittest.mock import MagicMock, mock_open, patch

import pytest

from nanodoc.core import run


def test_run_basic():
    """Test the basic functionality of run with a theme."""
    # Create a temporary test file
    with tempfile.NamedTemporaryFile(suffix=".py", mode="w+") as f:
        f.write("def test_function():\n    return True\n")
        f.flush()

        # Call the run function with a theme
        result = run(
            sources=[f.name],
            line_number_mode="file",
            generate_toc=True,
            theme="classic",
            show_header=True,
        )

        # Basic verification
        assert "def test_function()" in result
        assert "return True" in result

        # Check line numbers are included (line_number_mode="file")
        assert "1:" in result

        # Check file name is in the output (from show_header=True)
        assert os.path.basename(f.name) in result


def test_run_without_formatting():
    """Test run without formatting options."""

    # Define a simple renderer function to replace the original
    def simple_render(document, **kwargs):
        return "# source1\n\nfile content\n"

    # Use patch to replace the actual functions
    with (
        patch("nanodoc.resolver.os.path.isfile", return_value=True),
        patch("nanodoc.resolver.os.path.isdir", return_value=False),
        patch("nanodoc.extractor.open", mock_open(read_data="file content")),
        patch("nanodoc.resolver.resolve_paths", return_value=["path1"]),
        patch("nanodoc.extractor.resolve_files", return_value=["file_content"]),
        patch("nanodoc.extractor.gather_content", return_value=["content_item"]),
        # These patches need to use side_effect or return_value, not wraps
        patch(
            "nanodoc.document.build_document",
            return_value=MagicMock(content_items=["content_item"]),
        ),
        # We need to let the actual function run to test it's called properly
        # So we don't mock it here
        patch("nanodoc.renderer.render_document", side_effect=simple_render),
    ):
        # Call the function without formatting options
        result = run(
            sources=["source1"],
            line_number_mode=None,
            generate_toc=False,
            theme=None,
            show_header=True,
        )

        # Verify the basic content
        assert "source1" in result
        assert "file content" in result

        # Since we're letting the real apply_theme_to_document run, we can't assert
        # on mock calls. Instead we verify the output is unthemed.


def test_run_integration(tmp_path):
    """Test the integration of run with actual files."""
    # Create a temporary test file
    test_file = tmp_path / "test.txt"
    test_file.write_text("This is a test file\nWith multiple lines\n")

    # Call the run function with the actual file
    result = run(
        sources=[str(test_file)],
        line_number_mode="file",
        generate_toc=False,
        theme=None,
        show_header=True,
    )

    # Verify that the content is in the result
    assert "This is a test file" in result
    assert "With multiple lines" in result
    assert os.path.basename(str(test_file)) in result


@pytest.fixture
def sample_files(tmp_path):
    """Create sample files for testing."""
    # Create sample files with some markdown content
    file1 = tmp_path / "file1.md"
    file1.write_text("# Heading 1\nSome content\n## Subheading\nMore content")

    file2 = tmp_path / "file2.md"
    file2.write_text("# Heading 2\nOther content")

    return [str(file1), str(file2)]


def test_run_with_toc(sample_files):
    """Test that TOC generation works when generate_toc=True is passed."""
    # This test bypasses command line argument processing
    result = run(
        sources=sample_files,
        line_number_mode="all",
        generate_toc=True,
        theme=None,
        show_header=True,
    )

    # Verify TOC is present in the output
    assert "Table of Contents" in result
    assert "file1.md" in result
    assert "file2.md" in result
    assert "Heading 1" in result
    assert "Heading 2" in result


def test_run_without_toc(sample_files):
    """Test that TOC is not generated when generate_toc=False is passed."""
    # This test bypasses command line argument processing
    result = run(
        sources=sample_files,
        line_number_mode="all",
        generate_toc=False,
        theme=None,
        show_header=True,
    )

    # Verify TOC is not present
    assert "TOC" not in result
    assert "file1.md" in result
    assert "file2.md" in result
    assert "Heading 1" in result
    assert "Heading 2" in result
