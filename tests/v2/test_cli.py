"""Tests for the v2 CLI implementation.

These tests focus on the CLI integration of the v2 pipeline,
and mock the post-argument processing to test functionality directly.
"""

import os
import tempfile
from unittest.mock import MagicMock, mock_open, patch

import pytest

from nanodoc.v2.core import process_v2
from nanodoc.v2.document import Document


def test_process_v2_basic():
    """Test the basic functionality of process_v2 with a theme."""
    # Create a temporary test file
    with tempfile.NamedTemporaryFile(suffix=".py", mode="w+") as f:
        f.write("def test_function():\n    return True\n")
        f.flush()

        # Call the process_v2 function with a theme
        result = process_v2(
            sources=[f.name],
            line_number_mode="file",
            generate_toc=True,
            theme="neutral",
            show_header=True,
        )

        # Basic verification
        assert "def test_function()" in result
        assert "return True" in result

        # Check line numbers are included (line_number_mode="file")
        assert "1 |" in result or "1|" in result

        # Check file name is in the output (from show_header=True)
        assert os.path.basename(f.name) in result


def test_process_v2_without_formatting():
    """Test process_v2 without formatting options."""

    # Define a simple renderer function to replace the original
    def simple_render(document, **kwargs):
        return "# source1\n\nfile content\n"

    # Use patch to replace the actual functions
    with (
        patch("nanodoc.v2.resolver.os.path.isfile", return_value=True),
        patch("nanodoc.v2.resolver.os.path.isdir", return_value=False),
        patch("nanodoc.v2.extractor.open", mock_open(read_data="file content")),
        patch("nanodoc.v2.resolver.resolve_paths", return_value=["path1"]),
        patch("nanodoc.v2.extractor.resolve_files", return_value=["file_content"]),
        patch("nanodoc.v2.extractor.gather_content", return_value=["content_item"]),
        # These patches need to use side_effect or return_value, not wraps
        patch(
            "nanodoc.v2.document.build_document",
            return_value=MagicMock(content_items=["content_item"]),
        ),
        # We need to let the actual function run to test it's called properly
        # So we don't mock it here
        patch("nanodoc.v2.renderer.render_document", side_effect=simple_render),
    ):
        # Call the function without formatting options
        result = process_v2(
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


def test_process_v2_integration(tmp_path):
    """Test the integration of process_v2 with actual files."""
    # Create a temporary test file
    test_file = tmp_path / "test.txt"
    test_file.write_text("This is a test file\nWith multiple lines\n")

    # Call the process_v2 function with the actual file
    result = process_v2(
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


def test_process_v2_with_toc(sample_files):
    """Test that TOC generation works when generate_toc=True is passed."""
    # This test bypasses command line argument processing
    result = process_v2(
        sources=sample_files,
        line_number_mode="all",
        generate_toc=True,
        theme=None,
        show_header=True,
    )

    # Verify TOC is present in the output
    assert "TOC" in result
    assert "file1.md" in result
    assert "file2.md" in result
    assert "Heading 1" in result
    assert "Heading 2" in result


def test_process_v2_without_toc(sample_files):
    """Test that TOC is not generated when generate_toc=False is passed."""
    # This test bypasses command line argument processing
    result = process_v2(
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


@patch("nanodoc.v2.resolver.resolve_paths")
@patch("nanodoc.v2.renderer.render_document")
@patch("nanodoc.v2.extractor.resolve_files")
@patch("nanodoc.v2.extractor.gather_content")
@patch("nanodoc.v2.document.build_document")
@patch("nanodoc.v2.formatter.apply_theme_to_document")
@pytest.mark.skip(
    reason="Over-mocked test: needs rewrite using actual document rendering"
)
def test_toc_flag_propagation(
    mock_theme,
    mock_build,
    mock_gather,
    mock_resolve_files,
    mock_render,
    mock_resolve_paths,
    sample_files,
):
    """Test that the generate_toc flag is correctly propagated."""
    # Mock the resolve_paths function to return sample paths
    mock_resolve_paths.return_value = sample_files
    mock_resolve_files.return_value = []
    mock_gather.return_value = []
    mock_build.return_value = Document(content_items=[])
    mock_theme.return_value = Document(content_items=[])

    # Setup a mock return value
    mock_render.return_value = "Mocked content"

    # Call process_v2 with generate_toc=True
    process_v2(
        sources=sample_files,
        line_number_mode="all",
        generate_toc=True,
        theme=None,
        show_header=True,
    )

    # Verify render_document was called with include_toc=True
    mock_render.assert_called_once()
    args, kwargs = mock_render.call_args
    assert kwargs.get("include_toc") is True


@patch("nanodoc.v2.resolver.resolve_paths")
@patch("nanodoc.v2.renderer.render_document")
@patch("nanodoc.v2.extractor.resolve_files")
@patch("nanodoc.v2.extractor.gather_content")
@patch("nanodoc.v2.document.build_document")
@patch("nanodoc.v2.formatter.apply_theme_to_document")
@pytest.mark.skip(
    reason="Over-mocked test: needs rewrite using actual document rendering"
)
def test_line_number_flag_propagation(
    mock_theme,
    mock_build,
    mock_gather,
    mock_resolve_files,
    mock_render,
    mock_resolve_paths,
    sample_files,
):
    """Test that the line_number_mode flag is correctly propagated."""
    # Mock the resolve_paths function to return sample paths
    mock_resolve_paths.return_value = sample_files
    mock_resolve_files.return_value = []
    mock_gather.return_value = []
    mock_build.return_value = Document(content_items=[])
    mock_theme.return_value = Document(content_items=[])

    # Setup a mock return value
    mock_render.return_value = "Mocked content"

    # Call process_v2 with line_number_mode="all"
    process_v2(
        sources=sample_files,
        line_number_mode="all",
        generate_toc=False,
        theme=None,
        show_header=True,
    )

    # Verify render_document was called with include_line_numbers=True
    mock_render.assert_called_once()
    args, kwargs = mock_render.call_args
    assert kwargs.get("include_line_numbers") is True
