from nanodoc.v1.core import run_all
from nanodoc.v1.files import get_files_from_args
from nanodoc.v1.formatting import create_header


def test_init_multiple_paths(tmpdir):
    """Test processing multiple paths."""
    # Create test files
    file1 = tmpdir.join("file1.txt")
    file2 = tmpdir.join("file2.txt")
    file1.write("Content 1")
    file2.write("Content 2")

    # Process multiple files
    output = run_all(
        [str(file1), str(file2)],
        line_number_mode=None,
        generate_toc=False,
        theme=None,
        show_header=True,
        sequence=None,
        style=None,
        txt_ext=None,
    )

    # Verify output contains both files
    assert "file1.txt" in output
    assert "Content 1" in output
    assert "file2.txt" in output
    assert "Content 2" in output


def test_init_multiple_paths_with_nonexistent(tmpdir):
    """Test processing multiple paths with a nonexistent file."""
    # Create one test file
    file1 = tmpdir.join("file1.txt")
    file1.write("Content 1")
    nonexistent = str(tmpdir.join("nonexistent.txt"))

    # Process multiple files including nonexistent
    output = run_all(
        [str(file1), nonexistent],
        line_number_mode=None,
        generate_toc=False,
        theme=None,
        show_header=True,
        sequence=None,
        style=None,
        txt_ext=None,
    )

    # Verify output contains existing file and error for nonexistent
    assert "file1.txt" in output
    assert "Content 1" in output
    assert "Error: File not found" in output
    assert "nonexistent.txt" in output


def test_init_multiple_paths_file_line_numbers(tmpdir):
    # Create test files
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("Line 1\nLine 2")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("Line 3\nLine 4")
    file_paths = [str(test_file1), str(test_file2)]

    # Call init with multiple test files and file line numbers
    # Get verified sources and process them with file line numbers
    verified_sources = get_files_from_args(file_paths)
    result = run_all(verified_sources, line_number_mode="file")

    # Assert that the file content is printed with file line numbers
    assert "1: Line 1" in result
    assert "2: Line 2" in result
    assert "1: Line 3" in result
    assert "2: Line 4" in result


def test_init_multiple_paths_all_line_numbers(tmpdir):
    # Create test files
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("Line 1\nLine 2")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("Line 3\nLine 4")
    file_paths = [str(test_file1), str(test_file2)]

    # Call init with multiple test files and all line numbers
    # Get verified sources and process them with all line numbers
    verified_sources = get_files_from_args(file_paths)
    result = run_all(verified_sources, line_number_mode="all")

    # Assert that the file content is printed with all line numbers
    assert "1: Line 1" in result
    assert "2: Line 2" in result
    assert "3: Line 3" in result
    assert "4: Line 4" in result


def test_init_multiple_paths_toc(tmpdir):
    # Create test files
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("Line 1\nLine 2")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("Line 3\nLine 4")
    file_paths = [str(test_file1), str(test_file2)]

    # Call init with multiple test files and TOC generation
    # Get verified sources and process them with TOC generation
    verified_sources = get_files_from_args(file_paths)
    result = run_all(verified_sources, generate_toc=True)

    # Assert that the TOC is generated and the file content is printed
    assert create_header("TOC") in result
    assert "test_file1.txt" in result
    assert "test_file2.txt" in result
    assert "Line 1" in result
    assert "Line 2" in result
    assert "Line 3" in result
    assert "Line 4" in result
