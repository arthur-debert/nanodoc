from nanodoc.core import process_all
from nanodoc.files import create_content_item
from nanodoc.toc import (
    calculate_line_numbers,
    calculate_toc_size,
    create_toc_content,
    format_filenames,
    generate_table_of_contents,
    group_content_items_by_file,
)


def test_group_content_items_by_file(tmpdir):
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("Line 1\nLine 2")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("Line 3\nLine 4")

    content_items = [
        create_content_item(str(test_file1)),
        create_content_item(str(test_file2)),
    ]

    file_groups = group_content_items_by_file(content_items)

    # Check that we have the correct number of groups
    assert len(file_groups) == 2

    # Check that each file path is a key in the dictionary
    assert str(test_file1) in file_groups
    assert str(test_file2) in file_groups

    # Check that each group contains the correct ContentItem
    assert len(file_groups[str(test_file1)]) == 1
    assert len(file_groups[str(test_file2)]) == 1
    assert file_groups[str(test_file1)][0] == content_items[0]
    assert file_groups[str(test_file2)][0] == content_items[1]


def test_calculate_toc_size(tmpdir):
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("Line 1\nLine 2")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("Line 3\nLine 4")

    content_items = [
        create_content_item(str(test_file1)),
        create_content_item(str(test_file2)),
    ]

    file_groups = group_content_items_by_file(content_items)
    toc_size = calculate_toc_size(file_groups)

    # TOC size should be:
    # 2 lines for header (header + blank line)
    # 2 lines for entries (1 per file)
    # 1 line for footer (blank line)
    # Total: 5 lines
    assert toc_size == 5

    # Test with multiple ranges for the same file
    test_file3 = tmpdir.join("test_file3.txt")
    test_file3.write("Line 5\nLine 6\nLine 7\nLine 8")

    # Create two ContentItems for the same file
    content_items = [
        create_content_item(str(test_file3)),
        create_content_item(str(test_file3)),
    ]

    file_groups = group_content_items_by_file(content_items)
    toc_size = calculate_toc_size(file_groups)

    # TOC size should be:
    # 2 lines for header (header + blank line)
    # 1 line for the file entry
    # 2 lines for subentries (1 per range)
    # 1 line for footer (blank line)
    # Total: 6 lines
    assert toc_size == 6


def test_calculate_line_numbers(tmpdir):
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("Line 1\nLine 2")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("Line 3\nLine 4")

    content_items = [
        create_content_item(str(test_file1)),
        create_content_item(str(test_file2)),
    ]

    file_groups = group_content_items_by_file(content_items)
    toc_size = calculate_toc_size(file_groups)
    line_numbers = calculate_line_numbers(file_groups, toc_size)

    # Check that we have line numbers for both files
    assert str(test_file1) in line_numbers
    assert str(test_file2) in line_numbers

    # First file should start at line toc_size + 3
    assert line_numbers[str(test_file1)] == toc_size + 3

    # Second file should start after first file
    # First file has 2 lines of content + 3 lines for header/footer
    # So second file should start at toc_size + 3 + 2 + 3
    assert line_numbers[str(test_file2)] == toc_size + 3 + 2 + 3


def test_format_filenames(tmpdir):
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("Line 1\nLine 2")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("Line 3\nLine 4")

    content_items = [
        create_content_item(str(test_file1)),
        create_content_item(str(test_file2)),
    ]

    file_groups = group_content_items_by_file(content_items)

    # Test with default style (None)
    formatted_filenames = format_filenames(file_groups)
    assert formatted_filenames[str(test_file1)] == "test_file1.txt"
    assert formatted_filenames[str(test_file2)] == "test_file2.txt"

    # Test with 'nice' style
    formatted_filenames = format_filenames(file_groups, style="nice")
    assert formatted_filenames[str(test_file1)] == ("Test File1 (test_file1.txt)")
    assert formatted_filenames[str(test_file2)] == ("Test File2 (test_file2.txt)")

    # Test with 'filename' style
    formatted_filenames = format_filenames(file_groups, style="filename")
    assert formatted_filenames[str(test_file1)] == "test_file1.txt"
    assert formatted_filenames[str(test_file2)] == "test_file2.txt"

    # Test with 'path' style
    formatted_filenames = format_filenames(file_groups, style="path")
    assert formatted_filenames[str(test_file1)] == str(test_file1)
    assert formatted_filenames[str(test_file2)] == str(test_file2)


def test_create_toc_content(tmpdir):
    test_file1 = tmpdir.join("test_file1.txt")
    test_file1.write("Line 1\nLine 2")
    test_file2 = tmpdir.join("test_file2.txt")
    test_file2.write("Line 3\nLine 4")

    content_items = [
        create_content_item(str(test_file1)),
        create_content_item(str(test_file2)),
    ]

    file_groups = group_content_items_by_file(content_items)
    toc_size = calculate_toc_size(file_groups)
    line_numbers = calculate_line_numbers(file_groups, toc_size)
    formatted_filenames = format_filenames(file_groups)

    toc = create_toc_content(file_groups, formatted_filenames, line_numbers)

    # Check that the TOC contains the expected content
    assert "TOC" in toc
    assert "test_file1.txt" in toc
    assert "test_file2.txt" in toc
    assert str(line_numbers[str(test_file1)]) in toc
    assert str(line_numbers[str(test_file2)]) in toc

    # Test with multiple ranges for the same file
    test_file3 = tmpdir.join("test_file3.txt")
    test_file3.write("Line 5\nLine 6\nLine 7\nLine 8")

    # Create two ContentItems for the same file
    content_items = [
        create_content_item(str(test_file3)),
        create_content_item(str(test_file3)),
    ]

    file_groups = group_content_items_by_file(content_items)
    toc_size = calculate_toc_size(file_groups)
    line_numbers = calculate_line_numbers(file_groups, toc_size)
    formatted_filenames = format_filenames(file_groups)

    toc = create_toc_content(file_groups, formatted_filenames, line_numbers)

    # Check that the TOC contains the expected content
    assert "TOC" in toc
    assert "test_file3.txt" in toc
    # Since we're not using line ranges, we won't have a. and b. subentries
    # Instead, check that the file name appears and the line number is present
    assert str(line_numbers[str(test_file3)]) in toc


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
