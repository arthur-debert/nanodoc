"""Integration tests for Nanodoc v2 formatting and rendering.

This module focuses on testing the theming and rendering functionality
with minimal mocking to ensure the components work together correctly.
"""

from nanodoc.formatter import (
    apply_theme_to_document,
    enhance_rendering,
    get_available_themes,
)
from nanodoc.renderer import render_document
from nanodoc.structures import Document, FileContent


def test_formatting_with_theme():
    """Test the full formatting pipeline with a theme."""
    # Create a simple document structure
    content = FileContent(
        filepath="test.py",
        content="def test_function():\n    return True\n",
        ranges=[(1, None)],
        is_bundle=False,
    )
    document = Document(content_items=[content])

    # Apply theme to document
    themed_doc = apply_theme_to_document(
        document, theme_name="classic", use_rich_formatting=True
    )

    # Render the document
    rendered = render_document(themed_doc)

    # Enhance the rendering with the theme
    final_output = enhance_rendering(
        rendered,
        theme_name=themed_doc.theme_name,
        use_rich_formatting=themed_doc.use_rich_formatting,
    )

    # Check the output contains the content
    assert "def test_function()" in final_output
    assert "return True" in final_output
    # Since we're using a real theme, this should work without mocking
    assert themed_doc.theme_name == "classic"
    assert themed_doc.use_rich_formatting is True


def test_all_themes():
    """Test rendering with all available themes."""
    # Create a simple document with headings to test theme styling
    content = FileContent(
        filepath="test.md",
        content="# Heading 1\n\nSome content\n\n## Heading 2\n\nMore content\n",
        ranges=[(1, None)],
        is_bundle=False,
    )
    document = Document(content_items=[content])

    # Get all available themes
    themes = get_available_themes()
    assert len(themes) > 0, "No themes found"

    for theme in themes:
        # Apply the current theme
        themed_doc = apply_theme_to_document(
            document, theme_name=theme, use_rich_formatting=True
        )

        # Render and enhance
        rendered = render_document(themed_doc)
        final_output = enhance_rendering(
            rendered,
            theme_name=themed_doc.theme_name,
            use_rich_formatting=themed_doc.use_rich_formatting,
        )

        # Basic validation - should contain our content
        assert "Heading 1" in final_output
        assert "Heading 2" in final_output
        assert "Some content" in final_output
        assert "More content" in final_output


def test_theme_with_real_file(tmp_path):
    """Test theming with a real file from disk."""
    # Create a temporary test file
    file_path = tmp_path / "sample.py"
    file_path.write_text(
        "# Sample Python file\n\ndef hello():\n    print('Hello world')\n\n"
        "# Another section\nclass Test:\n    def method(self):\n        pass\n"
    )

    # Create document with the real file
    content = FileContent(
        filepath=str(file_path),
        content=file_path.read_text(),
        ranges=[(1, None)],
        is_bundle=False,
    )
    document = Document(content_items=[content])

    # Apply theme and render
    themed_doc = apply_theme_to_document(
        document, theme_name="classic", use_rich_formatting=True
    )
    rendered = render_document(themed_doc)
    final_output = enhance_rendering(
        rendered,
        theme_name=themed_doc.theme_name,
        use_rich_formatting=themed_doc.use_rich_formatting,
    )

    # Check content
    assert "Sample Python file" in final_output
    assert "def hello()" in final_output
    assert "class Test" in final_output
    assert "Hello world" in final_output
