"""Tests for directory handling functionality."""

from pathlib import Path

from nanodoc.v1.files import expand_directory


def test_directory_expansion(fixture_content_item, tmp_path):
    """Test expanding a directory with various file types."""
    # Create a test directory with different file types
    dir_path = tmp_path / "test_dir"
    dir_path.mkdir()

    # Create test files of different types
    txt_file = dir_path / "test.txt"
    md_file = dir_path / "test.md"
    other_file = dir_path / "test.other"

    txt_file.write_text("test content")
    md_file.write_text("test content")
    other_file.write_text("test content")

    # Test expanding the directory
    expanded_files = expand_directory(str(dir_path))

    # Should include .txt and .md files, but not .other
    assert str(txt_file) in expanded_files
    assert str(md_file) in expanded_files
    assert str(other_file) not in expanded_files

    # Test with our fixture file
    fixture_dir = Path(fixture_content_item.file_path).parent
    fixture_files = expand_directory(str(fixture_dir))
    assert fixture_content_item.file_path in fixture_files


def test_expand_directory(tmpdir):
    # Create directory structure
    dir_path = tmpdir.mkdir("test_dir")
    nested_dir = dir_path.mkdir("nested_dir")
    test_file_txt = dir_path.join("test_file.txt")
    test_file_txt.write("test")
    test_file_md = dir_path.join("test_file.md")
    test_file_md.write("test")
    test_file_other = dir_path.join("test_file.other")
    test_file_other.write("test")
    nested_file_txt = nested_dir.join("nested_file.txt")
    nested_file_txt.write("test")

    # Call expand_directory
    expanded_files = expand_directory(str(dir_path))

    # Assert that only .txt and .md files are included, and nested files too
    assert str(test_file_txt) in expanded_files
    assert str(test_file_md) in expanded_files
    assert str(nested_file_txt) in expanded_files
    assert str(test_file_other) not in expanded_files


def test_expand_directory_empty(tmpdir):
    dir_path = tmpdir.mkdir("empty_dir")
    expanded_files = expand_directory(str(dir_path))
    assert expanded_files == []


def test_expand_args_file(tmpdir):
    # Create a test file
    test_file = tmpdir.join("test_file.txt")
    test_file.write("test")

    # Call expand_args with a file path
    expanded_files = expand_args([str(test_file)])

    # Assert that the file path is returned as a single-item list
    assert expanded_files == [str(test_file)]


def test_expand_args_directory(tmpdir):
    # Create directory structure
    dir_path = tmpdir.mkdir("test_dir")
    test_file_txt = dir_path.join("test_file.txt")
    test_file_txt.write("test")

    # Call expand_args with a directory path
    expanded_files = expand_args([str(dir_path)])

    # Assert that the directory is expanded to include the file
    assert str(test_file_txt) in expanded_files


def test_expand_args_bundle(tmpdir):
    # Create a bundle file
    bundle_file = tmpdir.join("test_bundle.txt")
    test_file = tmpdir.join("test_file.txt")
    test_file.write("test")
    bundle_file.write(str(test_file))

    # Call expand_args with a bundle file path
    expanded_files = expand_args([str(bundle_file)])

    # Assert that the bundle is expanded to include the file
    assert str(test_file) in expanded_files


def test_expand_directory_with_extensions(tmpdir):
    # Create directory structure
    dir_path = tmpdir.mkdir("test_dir")
    test_file_txt = dir_path.join("test_file.txt")
    test_file_txt.write("test")
    test_file_md = dir_path.join("test_file.md")
    test_file_md.write("test")
    test_file_other = dir_path.join("test_file.other")
    test_file_other.write("test")

    # Call expand_directory with specific extensions
    expanded_files = expand_directory(str(dir_path), extensions=[".txt"])

    # Assert that only .txt files are included
    assert str(test_file_txt) in expanded_files
    assert str(test_file_md) not in expanded_files
    assert str(test_file_other) not in expanded_files
