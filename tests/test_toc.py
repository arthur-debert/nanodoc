from nanodoc.core import process_all
from nanodoc.files import create_content_item
from nanodoc.toc import generate_table_of_contents


def test_generate_table_of_contents(tmpdir):
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("Line 1\nLine 2")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("Line 3\nLine 4")
    file_paths = [
        create_content_item(str(test_file1)),
        create_content_item(str(test_file2)),
    ]

    toc, line_numbers = generate_table_of_contents(file_paths)

    # Check TOC content
    assert "TOC" in toc
    assert "test_file1.txt" in toc
    assert "test_file2.txt" in toc

    # Check line numbers
    assert isinstance(line_numbers, dict)
    file1_path = str(test_file1)
    file2_path = str(test_file2)
    assert line_numbers[file1_path] > 0
    assert line_numbers[file2_path] > line_numbers[file1_path]


def test_generate_table_of_contents_with_style(tmpdir):
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("Line 1\nLine 2")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("Line 3\nLine 4")
    file_paths = [
        create_content_item(str(test_file1)),
        create_content_item(str(test_file2)),
    ]

    # Test with 'nice' style
    toc, line_numbers = generate_table_of_contents(file_paths, style="nice")

    # Check that styled filenames are in the TOC
    assert "Test File1 (test_file1.txt)" in toc
    assert "Test File2 (test_file2.txt)" in toc

    # Test with 'filename' style
    toc, _ = generate_table_of_contents(file_paths, style="filename")

    # Check that plain filenames are in the TOC
    assert "test_file1.txt" in toc
    assert "test_file2.txt" in toc

    # Test with 'path' style
    toc, _ = generate_table_of_contents(file_paths, style="path")

    # Check that full paths are in the TOC
    assert str(test_file1) in toc
    assert str(test_file2) in toc


# The following tests use process_all but are kept for integration testing
# They verify that the TOC generation works correctly when integrated with
# the rest of the code


def test_process_all_toc_generation(tmpdir):
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("Line 1")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("Line 2")
    file_paths = [
        create_content_item(str(test_file1)),
        create_content_item(str(test_file2)),
    ]
    output = process_all(file_paths, None, True)
    assert "TOC" in output
    assert "test_file1.txt" in output
    assert "test_file2.txt" in output
    assert output.startswith("\nTOC\n\n")


def test_process_all_with_no_header(tmpdir):
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("Line 1")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("Line 2")
    file_paths = [
        create_content_item(str(test_file1)),
        create_content_item(str(test_file2)),
    ]

    output = process_all(file_paths, None, False, show_header=False)

    assert "Line 1" in output
    assert "Line 2" in output
    assert "test_file1.txt" not in output
    assert "test_file2.txt" not in output


def test_process_all_with_header_sequence(tmpdir):
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("Line 1")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("Line 2")
    file_paths = [
        create_content_item(str(test_file1)),
        create_content_item(str(test_file2)),
    ]

    # Test numerical sequence
    output = process_all(file_paths, None, False, sequence="numerical")
    assert "1. test_file1.txt" in output
    assert "2. test_file2.txt" in output

    # Test letter sequence
    output = process_all(file_paths, None, False, sequence="letter")
    assert "a. test_file1.txt" in output
    assert "b. test_file2.txt" in output


def test_process_all_with_header_style(tmpdir):
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("Line 1")
    file_paths = [create_content_item(str(test_file1))]

    output = process_all(file_paths, None, False, style="nice")
    assert "Test File1 (test_file1.txt)" in output
