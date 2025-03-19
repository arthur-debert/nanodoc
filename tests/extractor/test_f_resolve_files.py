"""Tests for extractor.resolve_files function."""

from nanodoc.extractor import resolve_files


def test_resolve_files_basic():
    """Test resolving files with default settings."""
    file_paths = ["/path/to/file1.txt", "/path/to/file2.txt"]

    # Resolve files
    result = resolve_files(file_paths)

    # Check result
    assert len(result) == 2
    assert result[0].filepath == "/path/to/file1.txt"
    assert result[0].ranges == [(1, None)]
    assert result[0].is_bundle is False
    assert result[1].filepath == "/path/to/file2.txt"
    assert result[1].ranges == [(1, None)]
    assert result[1].is_bundle is False


def test_resolve_files_with_ranges():
    """Test resolving files with range specifiers."""
    file_paths = ["/path/to/file1.txt:10-20", "/path/to/file2.txt:5"]

    # Resolve files
    result = resolve_files(file_paths)

    # Check result
    assert len(result) == 2
    assert result[0].filepath == "/path/to/file1.txt"
    assert result[0].ranges == [(10, 20)]
    assert result[0].is_bundle is False
    assert result[1].filepath == "/path/to/file2.txt"
    assert result[1].ranges == [(5, 5)]
    assert result[1].is_bundle is False


def test_resolve_files_with_bundle_extensions():
    """Test resolving files with various bundle extensions."""
    file_paths = [
        "/path/to/file1.ndoc",  # Default bundle extension
        "/path/to/file2.bundle",  # Direct .bundle extension
        "/path/to/file3.bundle.txt",  # .bundle.* pattern
        "/path/to/file4.txt",  # Regular file
    ]

    # Use default bundle extensions
    result = resolve_files(file_paths)

    assert len(result) == 4
    assert result[0].is_bundle is True  # .ndoc file
    assert result[1].is_bundle is True  # .bundle file
    assert result[2].is_bundle is True  # .bundle.txt file
    assert result[3].is_bundle is False  # .txt file


def test_resolve_files_with_custom_bundle_extensions():
    """Test resolving files with custom bundle extensions."""
    file_paths = [
        "/path/to/file1.ndoc",  # Standard extension
        "/path/to/file2.bundle",  # Standard extension
        "/path/to/file3.custom",  # Custom extension
        "/path/to/file4.txt",  # Regular file
    ]

    # Use custom bundle extensions
    result = resolve_files(file_paths, bundle_extensions=[".custom"])

    assert len(result) == 4
    assert result[0].is_bundle is False  # .ndoc file (not in custom list)
    assert result[1].is_bundle is False  # .bundle file (not in custom list)
    assert result[2].is_bundle is True  # .custom file (in custom list)
    assert result[3].is_bundle is False  # .txt file
