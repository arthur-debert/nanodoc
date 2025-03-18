"""Tests for argument expansion functionality."""

from nanodoc.v1.files import expand_args


def test_single_file(fixture_content_item):
    """Test expanding a single file argument."""
    expanded_files = expand_args([fixture_content_item.file_path])
    assert len(expanded_files) == 1
    assert expanded_files[0] == fixture_content_item.file_path


def test_multiple_files(fixture_content_item):
    """Test expanding multiple file arguments."""
    expanded_files = expand_args([fixture_content_item.file_path])
    assert len(expanded_files) == 1
    assert fixture_content_item.file_path in expanded_files


def test_directory_expansion(fixture_content_item, tmp_path):
    """Test expanding a directory argument."""
    # Create a test directory with files
    dir_path = tmp_path / "test_dir"
    dir_path.mkdir()
    test_file = dir_path / "test.txt"
    test_file.write_text("test content")

    expanded_files = expand_args([str(dir_path)])
    assert len(expanded_files) >= 1
    assert str(test_file) in expanded_files


def test_bundle_expansion(fixture_content_item, tmp_path):
    """Test expanding a bundle file argument."""
    # Create a bundle file that references our fixture
    bundle_file = tmp_path / "bundle.ndoc"
    bundle_file.write_text(fixture_content_item.file_path)

    expanded_files = expand_args([str(bundle_file)])
    assert len(expanded_files) >= 1
    assert fixture_content_item.file_path in expanded_files


def test_mixed_expansion(fixture_content_item, tmp_path):
    """Test expanding mixed arguments (files, directories, bundles)."""
    # Create a test directory with files
    dir_path = tmp_path / "test_dir"
    dir_path.mkdir()
    test_file = dir_path / "test.txt"
    test_file.write_text("test content")

    # Create a bundle file that references our fixture
    bundle_file = tmp_path / "bundle.ndoc"
    bundle_file.write_text(fixture_content_item.file_path)

    expanded_files = expand_args(
        [fixture_content_item.file_path, str(dir_path), str(bundle_file)]
    )

    assert len(expanded_files) >= 3
    assert fixture_content_item.file_path in expanded_files
    assert str(test_file) in expanded_files
