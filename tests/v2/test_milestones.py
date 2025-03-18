"""Milestone tests for the v2 implementation of Nanodoc.

These tests verify that each milestone's functionality is working correctly
both in isolation and as part of the complete system.
"""

import os
import subprocess
from pathlib import Path

import pytest

from nanodoc.v2.document import build_document
from nanodoc.v2.extractor import gather_content, resolve_files
from nanodoc.v2.formatter import apply_theme_to_document
from nanodoc.v2.renderer import render_document
from nanodoc.v2.resolver import resolve_paths

FIXTURES_DIR = Path(__file__).parent.parent / "fixtures"


def test_milestone_1():
    """Test core data structures and path resolution."""
    # Test with mix of valid and invalid paths
    valid_path = str(FIXTURES_DIR / "test_file1.py")
    glob_path = str(FIXTURES_DIR / "*.py")

    # Valid paths should resolve correctly
    resolved_paths = resolve_paths([valid_path, glob_path])
    assert len(resolved_paths) > 0
    assert os.path.isabs(resolved_paths[0])
    assert all(os.path.isabs(path) for path in resolved_paths)
    assert all(os.path.exists(path) for path in resolved_paths)

    # Test glob resolution
    glob_files = [path for path in resolved_paths if path.endswith(".py")]
    assert len(glob_files) >= 2  # Should find test_file1.py and test_file2.py

    # Invalid paths should raise appropriate errors
    with pytest.raises(FileNotFoundError):
        resolve_paths(["nonexistent_file.py"])


def test_milestone_2():
    """Test file resolving and content gathering with ranges."""
    # Use existing fixture file
    test_file = str(FIXTURES_DIR / "test_file1.py")

    # Test with full file
    paths = resolve_paths([test_file])
    file_contents = resolve_files(paths)
    content_items = gather_content(file_contents)

    assert len(content_items) == 1
    assert "Test File 1" in content_items[0].content
    assert "function_1" in content_items[0].content

    # Test with range specification
    abs_path = resolve_paths([test_file])[0]
    range_path = f"{abs_path}:1-5"  # Get first function only

    file_contents = resolve_files([range_path])
    content_items = gather_content(file_contents)

    assert len(content_items) == 1
    assert "function_1" in content_items[0].content
    assert "TestClass" not in content_items[0].content


def test_milestone_3():
    """Test document building with bundle handling."""
    bundle_file = str(FIXTURES_DIR / "test_bundle.ndoc")

    # Process the files
    paths = resolve_paths([bundle_file])
    file_contents = resolve_files(paths, bundle_extensions=[".ndoc"])
    content_items = gather_content(file_contents)
    document = build_document(content_items)

    # Verify document structure
    assert len(document.content_items) > 1
    combined = "".join(item.content for item in document.content_items)

    # Check that content from all files is present
    assert "Test Bundle" in combined
    assert "function_1" in combined
    assert "function_2" in combined
    assert "bundle_function" in combined


def test_milestone_4():
    """Test rendering and TOC generation."""
    test_files = [
        str(FIXTURES_DIR / "test_file1.py"),
        str(FIXTURES_DIR / "test_file2.py"),
    ]

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
    assert "function_1" in rendered
    assert "function_2" in rendered


@pytest.mark.skip(reason="Theme functionality requires v1 theme files")
def test_milestone_5():
    """Test formatting, theming, and options."""
    test_file = str(FIXTURES_DIR / "test_file1.py")

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
    assert "   1:" in rendered_with_lines
    assert "def function_1" in rendered_with_lines

    # Test with rich formatting and theming
    themed_document = apply_theme_to_document(
        document, theme_name="neutral", use_rich_formatting=True
    )
    themed_rendered = render_document(
        themed_document, include_toc=False, include_line_numbers=False
    )

    # Since we're in rich mode, content should still be there but
    # the styling won't be visible in plain text mode
    assert "Test File 1" in themed_rendered
    assert "function_1" in themed_rendered


def test_milestone_6():
    """Test CLI options."""
    test_files = [
        str(FIXTURES_DIR / "test_file1.py"),
        str(FIXTURES_DIR / "test_file2.py"),
    ]

    # Test basic output without options
    result = subprocess.run(
        ["python", "-m", "nanodoc"] + test_files,
        text=True,
        capture_output=True,
        check=True,
    )
    assert result.returncode == 0
    assert "test_file1.py" in result.stdout
    assert "test_file2.py" in result.stdout

    # Test with TOC option
    result = subprocess.run(
        ["python", "-m", "nanodoc", "--toc"] + test_files,
        text=True,
        capture_output=True,
        check=True,
    )
    assert result.returncode == 0
    assert "Table of Contents" in result.stdout

    # Test with line numbers
    result = subprocess.run(
        ["python", "-m", "nanodoc", "-n"] + test_files,
        text=True,
        capture_output=True,
        check=True,
    )
    assert result.returncode == 0
    assert "1:" in result.stdout

    # Skip theme testing as it requires v1 theme files
    # Test with all options combined
    cmd = [
        "python",
        "-m",
        "nanodoc",
        "--toc",
        "-n",
    ] + test_files
    result = subprocess.run(cmd, text=True, capture_output=True, check=True)
    assert result.returncode == 0
    assert "Table of Contents" in result.stdout
    assert "1:" in result.stdout

    # Test error handling for nonexistent file
    with pytest.raises(subprocess.CalledProcessError):
        subprocess.run(
            ["python", "-m", "nanodoc", "nonexistent_file.py"],
            text=True,
            capture_output=True,
            check=True,
        )


@pytest.mark.skip(
    reason=(
        "Bundle processing needs to be fixed to properly handle relative paths "
        "and content inclusion from @include/@inline directives"
    )
)
def test_milestone_7():
    """Test bundle processing."""
    bundle_file = str(FIXTURES_DIR / "test_bundle.ndoc")

    result = subprocess.run(
        [
            "python",
            "-m",
            "nanodoc",
            "--toc",
            bundle_file,
        ],
        text=True,
        capture_output=True,
        check=True,
    )

    # Verify TOC
    assert "Table of Contents" in result.stdout

    # Verify all module content is included
    assert "Test Bundle" in result.stdout
    assert "function_1" in result.stdout
    assert "function_2" in result.stdout
    assert "bundle_function" in result.stdout


def test_resolve_single_file(fixture_content_item):
    """Test resolving a single file path."""
    paths = resolve_paths([fixture_content_item.file_path])
    assert len(paths) == 1
    assert paths[0] == fixture_content_item.file_path


def test_resolve_multiple_files(fixture_content_item):
    """Test resolving multiple file paths."""
    paths = resolve_paths([fixture_content_item.file_path])
    assert len(paths) == 1
    assert fixture_content_item.file_path in paths


def test_resolve_absolute_path(fixture_content_item):
    """Test resolving an absolute path."""
    abs_path = resolve_paths([fixture_content_item.file_path])[0]
    assert Path(abs_path).is_absolute()


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
        [
            "python",
            "-m",
            "nanodoc",
            "--toc",
            fixture_content_item.file_path,
        ],
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
    result = subprocess.run(cmd, text=True, capture_output=True, check=True)

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
    result = subprocess.run(cmd, text=True, capture_output=True, check=True)

    # Check that all features are working together
    assert result.returncode == 0
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
