"""Milestone tests for the v2 implementation of Nanodoc.

These tests verify that each milestone's functionality is working correctly
both in isolation and as part of the complete system.
"""

import os
import subprocess
import tempfile

import pytest

from nanodoc.v2.document import build_document
from nanodoc.v2.extractor import gather_content, resolve_files
from nanodoc.v2.formatter import apply_theme_to_document
from nanodoc.v2.renderer import render_document
from nanodoc.v2.resolver import resolve_paths


def test_milestone_1():
    """Test core data structures and path resolution."""
    # Test with mix of valid and invalid paths
    valid_path = "src/nanodoc/v2/cli.py"
    glob_path = "src/nanodoc/v2/*.py"

    # Valid paths should resolve correctly
    resolved_paths = resolve_paths([valid_path, glob_path])
    assert len(resolved_paths) > 0
    assert os.path.isabs(resolved_paths[0])
    assert all(os.path.isabs(path) for path in resolved_paths)
    assert all(os.path.exists(path) for path in resolved_paths)

    # Test glob resolution
    glob_files = [
        path
        for path in resolved_paths
        if os.path.dirname(path).endswith("v2") and path.endswith(".py")
    ]
    assert len(glob_files) > 1  # Should resolve to multiple files

    # Invalid paths should raise appropriate errors
    with pytest.raises(FileNotFoundError):
        resolve_paths(["nonexistent_file.py"])


def test_milestone_2():
    """Test file resolving and content gathering with ranges."""
    # Create a test file with numbered lines
    with tempfile.NamedTemporaryFile(suffix=".py", delete=False, mode="w") as temp:
        for i in range(1, 11):
            temp.write(f"Line {i}\n")
        test_file = temp.name

    try:
        # Test with full file
        paths = resolve_paths([test_file])
        file_contents = resolve_files(paths)
        content_items = gather_content(file_contents)

        assert len(content_items) == 1
        assert len(content_items[0].content.splitlines()) == 10

        # Test with range specification - the correct way to handle ranges
        # is to first resolve the real path, then add range specifier
        abs_path = resolve_paths([test_file])[0]
        # Note: Need to use 3-8 range to include lines 3 through 7 inclusively
        # because the end index is exclusive in the implementation
        range_path = f"{abs_path}:3-8"

        # Skip resolve_paths since it doesn't handle ranges
        file_contents = resolve_files([range_path])
        content_items = gather_content(file_contents)

        assert len(content_items) == 1
        assert len(content_items[0].content.splitlines()) == 5  # Lines 3-7
        assert "Line 3" in content_items[0].content
        assert "Line 7" in content_items[0].content
    finally:
        os.unlink(test_file)


def test_milestone_3():
    """Test document building with bundle handling."""
    # Create test files for bundling
    with tempfile.TemporaryDirectory() as tmpdirname:
        # Create a main bundle file with the correct bundle extension
        main_file = os.path.join(tmpdirname, "main.ndoc")
        with open(main_file, "w") as f:
            f.write("# Main file\n")
            f.write("@include included.py\n")
            f.write("def main_function():\n")
            f.write("    pass\n")
            f.write("@inline inlined.py\n")

        # Create an included file
        included_file = os.path.join(tmpdirname, "included.py")
        with open(included_file, "w") as f:
            f.write("# Included file\n")
            f.write("def included_function():\n")
            f.write("    pass\n")

        # Create an inlined file
        inlined_file = os.path.join(tmpdirname, "inlined.py")
        with open(inlined_file, "w") as f:
            f.write("# Inlined file\n")
            f.write("def inlined_function():\n")
            f.write("    pass\n")

        # Process the files
        paths = resolve_paths([main_file])
        # Explicitly specify bundle_extensions to include .ndoc files
        file_contents = resolve_files(paths, bundle_extensions=[".ndoc"])
        content_items = gather_content(file_contents)
        document = build_document(content_items)

        # Verify document structure
        assert len(document.content_items) > 1  # Should include all files
        combined_content = "".join(item.content for item in document.content_items)

        # Check that all function definitions are present
        assert "def main_function" in combined_content
        assert "def included_function" in combined_content
        assert "def inlined_function" in combined_content

        # Check file headers
        assert "Main file" in combined_content
        assert "Included file" in combined_content
        assert "Inlined file" in combined_content


def test_milestone_4():
    """Test rendering and TOC generation."""
    # Create multiple test files
    test_files = []
    try:
        for i in range(3):
            with tempfile.NamedTemporaryFile(
                suffix=".py", delete=False, mode="w"
            ) as temp:
                temp.write(f"# File {i+1}\n\ndef function_{i+1}():\n    pass\n")
                test_files.append(temp.name)

        # Process the files through the pipeline
        paths = resolve_paths(test_files)
        file_contents = resolve_files(paths)
        content_items = gather_content(file_contents)
        document = build_document(content_items)

        # Render with TOC
        rendered = render_document(document, include_toc=True)

        # Verify TOC
        assert "TOC" in rendered or "Table of Contents" in rendered

        # Check for file names in TOC
        for file_path in test_files:
            file_name = os.path.basename(file_path)
            assert file_name in rendered

        # Check that content is included
        for i in range(3):
            assert f"File {i+1}" in rendered
            assert f"function_{i+1}" in rendered
    finally:
        for file_path in test_files:
            if os.path.exists(file_path):
                os.unlink(file_path)


def test_milestone_5():
    """Test formatting, theming, and options."""
    # Create a test file
    with tempfile.NamedTemporaryFile(suffix=".py", delete=False, mode="w") as temp:
        temp.write("# Test Header\n\ndef test_function():\n    pass\n")
        test_file = temp.name

    try:
        # Process through the pipeline
        paths = resolve_paths([test_file])
        file_contents = resolve_files(paths)
        content_items = gather_content(file_contents)
        document = build_document(content_items)

        # Test line numbers
        document_with_lines = apply_theme_to_document(
            document, theme_name=None, use_rich_formatting=False
        )
        rendered_with_lines = render_document(
            document_with_lines, include_toc=False, include_line_numbers=True
        )

        # Verify line numbers are present
        assert "1" in rendered_with_lines
        assert "2" in rendered_with_lines

        # Test with rich formatting and theming (now should work)
        themed_document = apply_theme_to_document(
            document, theme_name="neutral", use_rich_formatting=True
        )
        themed_rendered = render_document(
            themed_document, include_toc=False, include_line_numbers=False
        )

        # Verify content is present (cannot verify colors in text output)
        assert "Test Header" in themed_rendered
        assert "test_function" in themed_rendered
    finally:
        os.unlink(test_file)


def test_milestone_6():
    """Test CLI integration with all features."""
    # Create test files
    test_files = []
    try:
        for i in range(3):
            with tempfile.NamedTemporaryFile(
                suffix=".py", delete=False, mode="w"
            ) as temp:
                temp.write(f"# File {i+1}\n\ndef function_{i+1}():\n    pass\n")
                test_files.append(temp.name)

        # Test basic output
        result = subprocess.run(
            ["python", "-m", "nanodoc", "--use-v2"] + test_files,
            capture_output=True,
            text=True,
            check=True,
        )
        for i in range(3):
            assert f"File {i+1}" in result.stdout
            assert f"function_{i+1}" in result.stdout

        # Test with TOC
        result = subprocess.run(
            ["python", "-m", "nanodoc", "--use-v2", "--toc"] + test_files,
            capture_output=True,
            text=True,
            check=True,
        )
        assert "TOC" in result.stdout or "Table of Contents" in result.stdout

        # Test with line numbers
        result = subprocess.run(
            ["python", "-m", "nanodoc", "--use-v2", "-n"] + test_files,
            capture_output=True,
            text=True,
            check=True,
        )
        assert "1" in result.stdout

        # Test with theme
        result = subprocess.run(
            ["python", "-m", "nanodoc", "--use-v2", "--theme", "neutral"] + test_files,
            capture_output=True,
            text=True,
            check=True,
        )

        # Test with all options together
        result = subprocess.run(
            ["python", "-m", "nanodoc", "--use-v2", "--toc", "-n", "--theme", "neutral"]
            + test_files,
            capture_output=True,
            text=True,
            check=True,
        )
        assert "TOC" in result.stdout or "Table of Contents" in result.stdout
        assert "1" in result.stdout

        # Test invalid inputs
        result = subprocess.run(
            ["python", "-m", "nanodoc", "--use-v2", "nonexistent_file.py"],
            capture_output=True,
            text=True,
        )
        assert result.returncode != 0
    finally:
        for file_path in test_files:
            if os.path.exists(file_path):
                os.unlink(file_path)


def test_milestone_7():
    """Complete end-to-end test of all functionality."""
    # Create a complex test setup with bundles and directives
    with tempfile.TemporaryDirectory() as tmpdirname:
        # Create a main bundle file with .ndoc extension
        main_file = os.path.join(tmpdirname, "main.ndoc")
        with open(main_file, "w") as f:
            f.write("# Main bundle file\n")
            f.write("@include module1.py\n")
            f.write("def main_function():\n")
            f.write("    pass\n")
            f.write("@inline module2.py\n")

        # Create included and inlined modules
        module1_file = os.path.join(tmpdirname, "module1.py")
        with open(module1_file, "w") as f:
            f.write("# Module 1\n")
            f.write("def module1_function():\n")
            f.write("    return 'Module 1'\n")

        module2_file = os.path.join(tmpdirname, "module2.py")
        with open(module2_file, "w") as f:
            f.write("# Module 2\n")
            f.write("def module2_function():\n")
            f.write("    return 'Module 2'\n")

        # Test with all options
        result = subprocess.run(
            [
                "python",
                "-m",
                "nanodoc",
                "--use-v2",
                "--toc",
                "-n",
                "--theme",
                "neutral",
                main_file,
            ],
            capture_output=True,
            text=True,
            check=True,
        )

        # Verify TOC
        assert "TOC" in result.stdout or "Table of Contents" in result.stdout

        # Verify all module content is included
        assert "Main bundle file" in result.stdout
        assert "Module 1" in result.stdout
        assert "Module 2" in result.stdout
        assert "main_function" in result.stdout
        assert "module1_function" in result.stdout
        assert "module2_function" in result.stdout

        # Verify line numbers
        assert "1" in result.stdout

        # Test error handling with circular dependencies
        circular1_file = os.path.join(tmpdirname, "circular1.ndoc")
        with open(circular1_file, "w") as f:
            f.write("# Circular 1\n")
            f.write("@include circular2.ndoc\n")

        circular2_file = os.path.join(tmpdirname, "circular2.ndoc")
        with open(circular2_file, "w") as f:
            f.write("# Circular 2\n")
            f.write("@include circular1.ndoc\n")

        result = subprocess.run(
            ["python", "-m", "nanodoc", "--use-v2", circular1_file],
            capture_output=True,
            text=True,
        )

        # Should detect circular dependency
        assert result.returncode != 0
        assert (
            "circular" in result.stderr.lower() or "dependency" in result.stderr.lower()
        )
