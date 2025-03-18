from nanodoc.v1.core import run_all
from nanodoc.v1.files import get_files_from_args
from nanodoc.v1.formatting import create_header


def test_directory_integration(tmpdir):
    """Test integration of directory processing."""
    # Create test files in a directory
    test_dir = tmpdir.mkdir("test_dir")
    file1 = test_dir.join("file1.txt")
    file2 = test_dir.join("file2.txt")
    file1.write("Content 1")
    file2.write("Content 2")

    # Process directory
    verified_sources = get_files_from_args([str(test_dir)])
    result = run_all(
        verified_sources,
        line_number_mode=None,
        generate_toc=False,
        theme=None,
        show_header=True,
        sequence=None,
        style=None,
        txt_ext=None,
    )

    # Verify output contains both files
    assert "file1.txt" in result
    assert "Content 1" in result
    assert "file2.txt" in result
    assert "Content 2" in result


def test_directory_with_toc(tmpdir):
    """Test directory processing with table of contents."""
    # Create test files in a directory
    test_dir = tmpdir.mkdir("test_dir")
    file1 = test_dir.join("file1.txt")
    file2 = test_dir.join("file2.txt")
    file1.write("Content 1")
    file2.write("Content 2")

    # Process directory with TOC
    verified_sources = get_files_from_args([str(test_dir)])
    result = run_all(
        verified_sources,
        line_number_mode=None,
        generate_toc=True,
        theme=None,
        show_header=True,
        sequence=None,
        style=None,
        txt_ext=None,
    )

    # Verify output contains TOC and file contents
    assert "Table of Contents" in result
    assert "file1.txt" in result
    assert "file2.txt" in result
    assert "Content 1" in result
    assert "Content 2" in result


def test_directory_with_subdirectories(tmpdir):
    """Test processing of directory with subdirectories."""
    # Create test files in nested directories
    test_dir = tmpdir.mkdir("test_dir")
    subdir = test_dir.mkdir("subdir")
    file1 = test_dir.join("file1.txt")
    file2 = subdir.join("file2.txt")
    file1.write("Content 1")
    file2.write("Content 2")

    # Process directory recursively
    verified_sources = get_files_from_args([str(test_dir)])
    result = run_all(
        verified_sources,
        line_number_mode=None,
        generate_toc=False,
        theme=None,
        show_header=True,
        sequence=None,
        style=None,
        txt_ext=None,
    )

    # Verify output contains files from all levels
    assert "file1.txt" in result
    assert "Content 1" in result
    assert "subdir/file2.txt" in result
    assert "Content 2" in result


def test_init_directory_no_line_numbers(tmpdir):
    # Create directory structure
    dir_path = tmpdir.mkdir("test_dir")
    test_file_txt = dir_path.join("test_file.txt")
    test_file_txt.write("Line 1\nLine 2")
    test_file_md = dir_path.join("test_file.md")
    test_file_md.write("Line 3\nLine 4")

    # Call init with the directory
    # Get verified sources and process them
    verified_sources = get_files_from_args([str(dir_path)])
    result = run_all(
        verified_sources,
        line_number_mode=None,
        generate_toc=False,
        theme=None,
        show_header=True,
        sequence=None,
        style=None,
        txt_ext=None,
    )

    # Assert that the file content is printed without line numbers
    assert "Line 1" in result
    assert "Line 2" in result
    assert "Line 3" in result
    assert "Line 4" in result
    assert "1:" not in result
    assert "2:" not in result
    assert "3:" not in result
    assert "4:" not in result


def test_init_directory_file_line_numbers(tmpdir):
    # Create directory structure
    dir_path = tmpdir.mkdir("test_dir")
    test_file_txt = dir_path.join("test_file.txt")
    test_file_txt.write("Line 1\nLine 2")
    test_file_md = dir_path.join("test_file.md")
    test_file_md.write("Line 3\nLine 4")

    # Call init with the directory and file line numbers
    # Get verified sources and process them with file line numbers
    verified_sources = get_files_from_args([str(dir_path)])
    result = run_all(
        verified_sources,
        line_number_mode="file",
        generate_toc=False,
        theme=None,
        show_header=True,
        sequence=None,
        style=None,
        txt_ext=None,
    )

    # Assert that the file content is printed with file line numbers
    assert "1: Line 1" in result
    assert "2: Line 2" in result
    assert "1: Line 3" in result
    assert "2: Line 4" in result


def test_init_directory_all_line_numbers(tmpdir):
    # Create directory structure
    dir_path = tmpdir.mkdir("test_dir")
    test_file_txt = dir_path.join("test_file.txt")
    test_file_txt.write("Line 1\nLine 2")
    test_file_md = dir_path.join("test_file.md")
    test_file_md.write("Line 3\nLine 4")

    # Call init with the directory and all line numbers
    # Get verified sources and process them with all line numbers
    verified_sources = get_files_from_args([str(dir_path)])
    result = run_all(
        verified_sources,
        line_number_mode="all",
        generate_toc=False,
        theme=None,
        show_header=True,
        sequence=None,
        style=None,
        txt_ext=None,
    )

    # Assert that the file content is printed with all line numbers
    assert "1: Line 1" in result
    assert "2: Line 2" in result
    assert "3: Line 3" in result
    assert "4: Line 4" in result


def test_init_directory_toc(tmpdir):
    # Create directory structure
    dir_path = tmpdir.mkdir("test_dir")
    test_file_txt = dir_path.join("test_file.txt")
    test_file_txt.write("Line 1\nLine 2")
    test_file_md = dir_path.join("test_file.md")
    test_file_md.write("Line 3\nLine 4")

    # Call init with the directory and TOC generation
    # Get verified sources and process them with TOC generation
    verified_sources = get_files_from_args([str(dir_path)])
    result = run_all(
        verified_sources,
        line_number_mode=None,
        generate_toc=True,
        theme=None,
        show_header=True,
        sequence=None,
        style=None,
        txt_ext=None,
    )

    # Assert that the TOC is generated and the file content is printed
    assert create_header("TOC") in result
    assert "test_file.txt" in result
    assert "test_file.md" in result
    assert "Line 1" in result
    assert "Line 2" in result
    assert "Line 3" in result
    assert "Line 4" in result
