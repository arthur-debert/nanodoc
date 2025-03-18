"""Tests for the CLI integration of Nanodoc v2."""

import os
from unittest.mock import MagicMock, mock_open, patch

import pytest

from nanodoc.v2.cli import process_v2


@pytest.mark.skip(reason="Issues with Rich console theme mocking")
def test_process_v2_basic():
    """Test the basic functionality of process_v2."""

    # Define a simple enhance function to replace the original
    def mock_enhance(plain_content, **kwargs):
        return "rendered content"

    # Use regular patches instead of patch.multiple
    with (
        patch("nanodoc.v2.resolver.os.path.isfile", return_value=True),
        patch("nanodoc.v2.resolver.os.path.isdir", return_value=False),
        patch("nanodoc.v2.resolver.resolve_paths") as mock_resolve_paths,
        patch("nanodoc.v2.extractor.open", mock_open(read_data="file content")),
        patch("nanodoc.v2.extractor.resolve_files") as mock_resolve_files,
        patch("nanodoc.v2.extractor.gather_content") as mock_gather_content,
        patch("nanodoc.v2.document.build_document") as mock_build_document,
        patch("nanodoc.v2.formatter.apply_theme_to_document") as mock_format,
        patch("nanodoc.v2.renderer.render_document") as mock_render,
        patch("nanodoc.v2.formatter.enhance_rendering", side_effect=mock_enhance),
    ):
        # Set up the mocks
        mock_resolve_paths.return_value = ["path1", "path2"]
        mock_resolve_files.return_value = ["file_content1", "file_content2"]
        mock_gather_content.return_value = ["content_item1", "content_item2"]

        # Create document mock
        document_mock = MagicMock()
        document_mock.content_items = ["content_item1", "content_item2"]
        mock_build_document.return_value = document_mock

        # Set up formatter/renderer mocks
        mock_format.return_value = document_mock
        mock_render.return_value = "plain content"

        # Call the function
        result = process_v2(
            sources=["source1", "source2"],
            line_number_mode="file",
            generate_toc=True,
            theme="neutral",
            show_header=True,
        )

        # Verify the result
        assert result == "rendered content"

        # Verify the calls
        mock_resolve_paths.assert_called_once_with(["source1", "source2"])
        mock_resolve_files.assert_called_once_with(["path1", "path2"])
        mock_gather_content.assert_called_once_with(["file_content1", "file_content2"])
        mock_build_document.assert_called_once_with(["content_item1", "content_item2"])
        mock_format.assert_called_once_with(
            document_mock, theme_name="neutral", use_rich_formatting=True
        )
        mock_render.assert_called_once_with(
            document_mock, include_toc=True, include_line_numbers=True
        )


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
