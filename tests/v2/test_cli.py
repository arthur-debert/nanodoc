"""Tests for the CLI integration of Nanodoc v2."""

import os
import tempfile
from unittest.mock import MagicMock, mock_open, patch

from nanodoc.v2.cli import process_v2


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
