from nanodoc.core import process_all
from nanodoc.files import get_files_from_args
from nanodoc.formatting import create_header


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
    result = process_all(verified_sources)

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
    result = process_all(verified_sources, line_number_mode="file")

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
    result = process_all(verified_sources, line_number_mode="all")

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
    result = process_all(verified_sources, generate_toc=True)

    # Assert that the TOC is generated and the file content is printed
    assert create_header("TOC") in result
    assert "test_file.txt" in result
    assert "test_file.md" in result
    assert "Line 1" in result
    assert "Line 2" in result
    assert "Line 3" in result
    assert "Line 4" in result
